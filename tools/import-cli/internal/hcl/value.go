// Package hcl converts Terraform model values into HCL-compatible values,
// preserving structure, null vs empty, and masking sensitive fields.
package hcl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

// SecretRefPattern matches Cribl stored secret references (e.g. #42:base64...).
var SecretRefPattern = regexp.MustCompile(`^#\d+:`)

// Kind indicates the type of an HCL Value.
type Kind int

const (
	KindNull        Kind = iota
	KindString
	KindNumber
	KindBool
	KindList
	KindMap
	KindSensitive
	// KindVariableRef is a plain variable reference (var.name) for string attributes like criblio_secret.value.
	// Emitted as var.name in HCL; use for sensitive values that are plain strings, not JSON.
	KindVariableRef
)

// Value represents an HCL-compatible value that preserves null vs empty
// and can represent sensitive placeholders.
type Value struct {
	Kind     Kind
	String   string
	Number   float64
	Bool     bool
	List     []Value
	Map      map[string]Value
	// Sensitive: placeholder text when Kind == KindSensitive (e.g. "(sensitive)" or variable name for jsonencode(var.x)).
	// VarName: variable name when Kind == KindVariableRef (plain var.x, not jsonencode).
	Sensitive string
	VarName   string
}

// IsNull returns true if the value is explicitly null.
func (v Value) IsNull() bool {
	return v.Kind == KindNull
}

// IsSensitive returns true if the value should be emitted as a placeholder.
func (v Value) IsSensitive() bool {
	return v.Kind == KindSensitive
}

// PruneEmptyLists returns a copy of v with map entries whose value is an empty list
// removed recursively. Use when SkipEmptyListAttributes is true to avoid emitting
// nested empty lists that violate schema (e.g. list must contain at least 1 element).
func PruneEmptyLists(v Value) Value {
	if v.Kind == KindMap {
		out := make(map[string]Value, len(v.Map))
		for k, val := range v.Map {
			val = PruneEmptyLists(val)
			if val.Kind == KindList && len(val.List) == 0 {
				continue
			}
			out[k] = val
		}
		return Value{Kind: KindMap, Map: out}
	}
	if v.Kind == KindList {
		list := make([]Value, len(v.List))
		for i, el := range v.List {
			list[i] = PruneEmptyLists(el)
		}
		return Value{Kind: KindList, List: list}
	}
	return v
}

// PruneNulls returns a copy of v with null map entries removed recursively.
// Optional+Computed attributes written as null in config cause Terraform to show
// "(known after apply)" and a perpetual diff; omitting them lets Terraform use state.
func PruneNulls(v Value) Value {
	if v.Kind != KindMap && v.Kind != KindList {
		return v
	}
	if v.Kind == KindMap {
		out := make(map[string]Value, len(v.Map))
		for k, val := range v.Map {
			if val.Kind == KindNull {
				continue
			}
			out[k] = PruneNulls(val)
		}
		return Value{Kind: KindMap, Map: out}
	}
	// KindList: prune each element, keep list length (do not drop null elements).
	list := make([]Value, len(v.List))
	for i, el := range v.List {
		list[i] = PruneNulls(el)
	}
	return Value{Kind: KindList, List: list}
}

// Options configures conversion from Terraform models to HCL values.
type Options struct {
	// SensitivePaths is a set of attribute paths (e.g. "password", "conf.auth_token")
	// or path prefixes that should be emitted as placeholders. Keys are path strings.
	SensitivePaths map[string]bool
	// SensitivePlaceholder is the string emitted for sensitive values (default "(sensitive)").
	SensitivePlaceholder string
	// SkipAttributes is a set of top-level attribute names to omit from output (e.g. "items" when
	// the provider marks them as read-only/Computed only). Generated config must not set these.
	SkipAttributes map[string]bool
	// SkipAttributesNested is a set of attribute names to omit at any nesting level (e.g.
	// "additional_properties" in route items). Use when the provider marks nested attrs as read-only.
	SkipAttributesNested map[string]bool
}

func (o *Options) placeholder() string {
	if o != nil && o.SensitivePlaceholder != "" {
		return o.SensitivePlaceholder
	}
	return "(sensitive)"
}

// ModelToValue converts a Terraform resource model (struct with tfsdk tags)
// into a map of HCL-compatible values keyed by attribute name. Nested objects
// become nested maps; lists and maps are preserved. Null and empty are distinguished.
func ModelToValue(model interface{}, opts *Options) (map[string]Value, error) {
	if model == nil {
		return nil, nil
	}
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be struct or *struct, got %s", val.Kind())
	}
	return structToValue(val, "", opts)
}

func structToValue(val reflect.Value, pathPrefix string, opts *Options) (map[string]Value, error) {
	out := make(map[string]Value)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		sf := typ.Field(i)
		tfsdk := sf.Tag.Get("tfsdk")
		if tfsdk == "" || tfsdk == "-" {
			continue
		}
		path := pathPrefix
		if path != "" {
			path += "."
		}
		path += tfsdk
		// Skip top-level read-only attrs (id, items); skip nested read-only attrs (e.g. additional_properties, id in route items).
		if opts != nil {
			if opts.SkipAttributes != nil && opts.SkipAttributes[tfsdk] && pathPrefix == "" {
				continue
			}
			if opts.SkipAttributesNested != nil && opts.SkipAttributesNested[tfsdk] {
				continue
			}
		}
		sensitive := opts != nil && opts.SensitivePaths != nil && (opts.SensitivePaths[path] || opts.SensitivePaths[tfsdk])
		val, err := fieldToValue(field, path, sensitive, opts)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		out[tfsdk] = val
	}
	return out, nil
}

func fieldToValue(field reflect.Value, path string, sensitive bool, opts *Options) (Value, error) {
	if sensitive && opts != nil {
		return Value{Kind: KindSensitive, Sensitive: opts.placeholder()}, nil
	}
	// Unwrap pointer for optional nested structs
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return Value{Kind: KindNull}, nil
		}
		field = field.Elem()
	}
	if !field.IsValid() {
		return Value{Kind: KindNull}, nil
	}

	// types.String and other framework types
	if v, ok := field.Interface().(types.String); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindString, String: v.ValueString()}, nil
	}
	if v, ok := field.Interface().(types.Bool); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindBool, Bool: v.ValueBool()}, nil
	}
	if v, ok := field.Interface().(types.Int64); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindNumber, Number: float64(v.ValueInt64())}, nil
	}
	if v, ok := field.Interface().(types.Float64); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindNumber, Number: v.ValueFloat64()}, nil
	}
	if v, ok := field.Interface().(types.List); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		elems := v.Elements()
		list := make([]Value, 0, len(elems))
		for i, el := range elems {
			ev, err := attrToValue(el, path+fmt.Sprintf("[%d]", i), false, opts)
			if err != nil {
				return Value{}, err
			}
			list = append(list, ev)
		}
		return Value{Kind: KindList, List: list}, nil
	}
	if v, ok := field.Interface().(types.Map); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		elemMap := v.Elements()
		m := make(map[string]Value, len(elemMap))
		for k, el := range elemMap {
			ev, err := attrToValue(el, path+"."+k, false, opts)
			if err != nil {
				return Value{}, err
			}
			m[k] = ev
		}
		return Value{Kind: KindMap, Map: m}, nil
	}
	if v, ok := field.Interface().(types.Object); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		attrs := v.Attributes()
		m := make(map[string]Value, len(attrs))
		for k, av := range attrs {
			ev, err := attrToValue(av, path+"."+k, false, opts)
			if err != nil {
				return Value{}, err
			}
			m[k] = ev
		}
		return Value{Kind: KindMap, Map: m}, nil
	}

	// jsontypes.Normalized
	if v, ok := field.Interface().(jsontypes.Normalized); ok {
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindString, String: v.ValueString()}, nil
	}

	// Nested struct
	if field.Kind() == reflect.Struct {
		nested, err := structToValue(field, path, opts)
		if err != nil {
			return Value{}, err
		}
		return Value{Kind: KindMap, Map: nested}, nil
	}

	// Slice of map[string]jsontypes.Normalized (e.g. Items)
	if field.Kind() == reflect.Slice {
		if field.IsNil() {
			return Value{Kind: KindNull}, nil
		}
		n := field.Len()
		list := make([]Value, 0, n)
		for i := 0; i < n; i++ {
			el := field.Index(i)
			if el.Kind() == reflect.Map {
				m := make(map[string]Value)
				for _, key := range el.MapKeys() {
					mv := el.MapIndex(key)
					if nv, ok := mv.Interface().(jsontypes.Normalized); ok {
						if nv.IsNull() || nv.IsUnknown() {
							m[key.String()] = Value{Kind: KindNull}
						} else {
							m[key.String()] = Value{Kind: KindString, String: nv.ValueString()}
						}
					}
				}
				list = append(list, Value{Kind: KindMap, Map: m})
			} else {
				ev, err := fieldToValue(el, path+fmt.Sprintf("[%d]", i), false, opts)
				if err != nil {
					return Value{}, err
				}
				list = append(list, ev)
			}
		}
		return Value{Kind: KindList, List: list}, nil
	}

	// map[string]T (e.g. map[string]PipelineGroups, map[string]types.String)
	if field.Kind() == reflect.Map {
		if field.IsNil() {
			return Value{Kind: KindNull}, nil
		}
		m := make(map[string]Value)
		for _, key := range field.MapKeys() {
			mv := field.MapIndex(key)
			ev, err := fieldToValue(mv, path+"."+key.String(), false, opts)
			if err != nil {
				return Value{}, err
			}
			m[key.String()] = ev
		}
		return Value{Kind: KindMap, Map: m}, nil
	}

	return Value{}, fmt.Errorf("unsupported field type %s", field.Type())
}

func attrToValue(av attr.Value, path string, sensitive bool, opts *Options) (Value, error) {
	if av == nil {
		return Value{Kind: KindNull}, nil
	}
	if sensitive && opts != nil {
		return Value{Kind: KindSensitive, Sensitive: opts.placeholder()}, nil
	}
	switch v := av.(type) {
	case types.String:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindString, String: v.ValueString()}, nil
	case types.Bool:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindBool, Bool: v.ValueBool()}, nil
	case types.Int64:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindNumber, Number: float64(v.ValueInt64())}, nil
	case types.Float64:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		return Value{Kind: KindNumber, Number: v.ValueFloat64()}, nil
	case types.List:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		elems := v.Elements()
		list := make([]Value, 0, len(elems))
		for i, el := range elems {
			ev, err := attrToValue(el, path+fmt.Sprintf("[%d]", i), false, opts)
			if err != nil {
				return Value{}, err
			}
			list = append(list, ev)
		}
		return Value{Kind: KindList, List: list}, nil
	case types.Map:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		elemMap := v.Elements()
		m := make(map[string]Value, len(elemMap))
		for k, el := range elemMap {
			ev, err := attrToValue(el, path+"."+k, false, opts)
			if err != nil {
				return Value{}, err
			}
			m[k] = ev
		}
		return Value{Kind: KindMap, Map: m}, nil
	case types.Object:
		if v.IsNull() || v.IsUnknown() {
			return Value{Kind: KindNull}, nil
		}
		attrs := v.Attributes()
		m := make(map[string]Value, len(attrs))
		for k, el := range attrs {
			// Skip top-level read-only attrs (id, items); skip nested read-only attrs (e.g. additional_properties in route items).
			if opts != nil {
				if opts.SkipAttributes != nil && opts.SkipAttributes[k] && path == "" {
					continue
				}
				if opts.SkipAttributesNested != nil && opts.SkipAttributesNested[k] {
					continue
				}
			}
			ev, err := attrToValue(el, path+"."+k, false, opts)
			if err != nil {
				return Value{}, err
			}
			m[k] = ev
		}
		return Value{Kind: KindMap, Map: m}, nil
	default:
		return Value{}, fmt.Errorf("unsupported attr type %T", av)
	}
}

// IsSecretValue returns true if the string looks like a secret (e.g. Cribl stored secret ref #42:...).
func IsSecretValue(s string) bool {
	return s != "" && SecretRefPattern.MatchString(s)
}

// SecretValueVariableName returns the variable name used for criblio_secret.value (plain var ref).
// Example: SecretValueVariableName("secret_default_test_secret") => "secret_default_test_secret_value"
func SecretValueVariableName(resourceName string) string {
	return sanitizeVarName(resourceName, "value")
}

// CertificateCertVariableName returns the variable name for criblio_certificate.cert (sensitive, plain var ref).
func CertificateCertVariableName(resourceName string) string {
	return sanitizeVarName(resourceName, "cert")
}

// CertificatePrivKeyVariableName returns the variable name for criblio_certificate.priv_key (sensitive, plain var ref).
func CertificatePrivKeyVariableName(resourceName string) string {
	return sanitizeVarName(resourceName, "priv_key")
}

// sanitizeVarName returns a Terraform-valid variable name from a resource name and path.
func sanitizeVarName(resourceName, path string) string {
	combined := resourceName
	if path != "" {
		combined = resourceName + "_" + path
	}
	// Replace invalid chars with underscore; keep alphanumeric and underscore.
	var b strings.Builder
	for i, r := range combined {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || (r >= '0' && r <= '9' && i > 0) {
			b.WriteRune(r)
		} else if r == '.' || r == '[' || r == ']' || r == ' ' || r == '-' {
			b.WriteByte('_')
		}
	}
	s := b.String()
	// Collapse multiple underscores and trim
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	s = strings.Trim(s, "_")
	if s == "" {
		return "secret"
	}
	return s
}

// ReplaceSecretValuesWithVariableRefs recursively replaces secret-looking string values with
// KindSensitive and Sensitive set to a variable name (var.<name> will be emitted in HCL).
// resourceName is used to build unique variable names. Returns the list of variable names used.
//
// Mutates attrs in place. Safe because we only modify Value structs (Kind, Sensitive, String),
// not the map keys or structure.
func ReplaceSecretValuesWithVariableRefs(attrs map[string]Value, resourceName string) []string {
	var used []string
	var replaceInValue func(*Value, string)
	replaceInValue = func(v *Value, path string) {
		switch v.Kind {
		case KindString:
			// Replace secret refs (#42:...) and any token attribute value (emit as jsonencode(var.xxx)).
			shouldReplace := IsSecretValue(v.String) ||
				(strings.Contains(path, ".token") && v.String != "")
			if shouldReplace {
				name := sanitizeVarName(resourceName, path)
				v.Kind = KindSensitive
				v.Sensitive = name
				v.String = ""
				used = append(used, name)
			}
		case KindList:
			for i := range v.List {
				replaceInValue(&v.List[i], path+fmt.Sprintf("[%d]", i))
			}
		case KindMap:
			for k, val := range v.Map {
				replaceInValue(&val, path+"."+k)
				v.Map[k] = val
			}
		}
	}
	for k, val := range attrs {
		replaceInValue(&val, k)
		attrs[k] = val
	}
	return used
}

// isValidJSON returns true if s is valid JSON (RFC 7159). Used to detect values that must be
// replaced with var refs when the provider expects jsontypes.NormalizedType.
func isValidJSON(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	return json.Valid([]byte(s))
}

// ReplaceNonJSONStringsWithVariableRefs recursively replaces string values that are not valid
// JSON with KindSensitive (emitted as jsonencode(var.<name>)). Use for criblio_source/criblio_pack_source
// where many nested attributes use jsontypes.NormalizedType and literal non-JSON strings (e.g.
// "test", "C.Misc.uuidv5(...)") cause "Invalid JSON String Value" at plan.
//
// Mutates attrs in place. Safe because we only modify Value structs, not the map structure.
func ReplaceNonJSONStringsWithVariableRefs(attrs map[string]Value, resourceName string) []string {
	var used []string
	var replaceInValue func(*Value, string)
	replaceInValue = func(v *Value, path string) {
		switch v.Kind {
		case KindString:
			if v.String != "" && !isValidJSON(v.String) {
				name := sanitizeVarName(resourceName, path)
				v.Kind = KindSensitive
				v.Sensitive = name
				v.String = ""
				used = append(used, name)
			}
		case KindList:
			for i := range v.List {
				replaceInValue(&v.List[i], path+fmt.Sprintf("[%d]", i))
			}
		case KindMap:
			for k, val := range v.Map {
				replaceInValue(&val, path+"."+k)
				v.Map[k] = val
			}
		}
	}
	for k, val := range attrs {
		replaceInValue(&val, k)
		attrs[k] = val
	}
	return used
}

// CollectPlainVariableRefNames returns variable names used as KindVariableRef (plain var.x, not jsonencode).
// Used when replacing placeholders in generated HCL.
func CollectPlainVariableRefNames(attrs map[string]Value) []string {
	var names []string
	seen := make(map[string]bool)
	var collect func(Value)
	collect = func(v Value) {
		switch v.Kind {
		case KindVariableRef:
			if v.VarName != "" && isVariableName(v.VarName) && !seen[v.VarName] {
				seen[v.VarName] = true
				names = append(names, v.VarName)
			}
		case KindList:
			for _, el := range v.List {
				collect(el)
			}
		case KindMap:
			for _, val := range v.Map {
				collect(val)
			}
		}
	}
	for _, v := range attrs {
		collect(v)
	}
	return names
}

// CollectSecretVariableNames recurses through attrs and returns all Sensitive values that are
// variable names (non-empty, valid identifier). Used when writing variables.tf.
func CollectSecretVariableNames(attrs map[string]Value) []string {
	var names []string
	seen := make(map[string]bool)
	var collect func(Value)
	collect = func(v Value) {
		switch v.Kind {
		case KindSensitive:
			if v.Sensitive != "" && isVariableName(v.Sensitive) && !seen[v.Sensitive] {
				seen[v.Sensitive] = true
				names = append(names, v.Sensitive)
			}
		case KindVariableRef:
			if v.VarName != "" && isVariableName(v.VarName) && !seen[v.VarName] {
				seen[v.VarName] = true
				names = append(names, v.VarName)
			}
		case KindList:
			for _, el := range v.List {
				collect(el)
			}
		case KindMap:
			for _, val := range v.Map {
				collect(val)
			}
		}
	}
	for _, v := range attrs {
		collect(v)
	}
	return names
}

func isVariableName(s string) bool {
	if s == "" || s == "(sensitive)" {
		return false
	}
	for i, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' {
			continue
		}
		if r >= '0' && r <= '9' && i > 0 {
			continue
		}
		return false
	}
	return true
}
