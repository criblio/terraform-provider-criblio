// Conversion from our Value type to cty.Value for hclwrite SetAttributeValue.
package hcl

import (
	"fmt"
	"sort"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// VarRefPlaceholderPrefix is the prefix for a placeholder string that will be replaced
// with jsonencode(var.<name>) in the generated HCL. Used for nested secret refs.
const VarRefPlaceholderPrefix = "__VAR_REF__"
const VarRefPlaceholderSuffix = "__"

// PlainVarRefPlaceholderPrefix/Suffix are replaced with var.<name> (no jsonencode) for plain string secrets (e.g. criblio_secret.value).
const PlainVarRefPlaceholderPrefix = "__VAR_REF_PLAIN__"
const PlainVarRefPlaceholderSuffix = "__"

// ValueToCty converts a Value to cty.Value so it can be passed to hclwrite.Body.SetAttributeValue.
// Nested maps and lists are converted recursively; null and structure are preserved.
// KindSensitive with a variable name is emitted as a placeholder string; FileWithResources
// replaces it with jsonencode(var.xxx) so secrets are always JSON-encoded.
func ValueToCty(v Value) (cty.Value, error) {
	switch v.Kind {
	case KindNull:
		return cty.NullVal(cty.DynamicPseudoType), nil
	case KindSensitive:
		if isVariableName(v.Sensitive) {
			return cty.StringVal(VarRefPlaceholderPrefix + v.Sensitive + VarRefPlaceholderSuffix), nil
		}
		return cty.StringVal(v.Sensitive), nil
	case KindVariableRef:
		if isVariableName(v.VarName) {
			return cty.StringVal(PlainVarRefPlaceholderPrefix + v.VarName + PlainVarRefPlaceholderSuffix), nil
		}
		return cty.NullVal(cty.DynamicPseudoType), nil
	case KindString:
		return cty.StringVal(v.String), nil
	case KindBool:
		return cty.BoolVal(v.Bool), nil
	case KindNumber:
		return cty.NumberFloatVal(v.Number), nil
	case KindList:
		v = normalizeAndPruneList(v)
		elems, err := listToCtyElems(v.List)
		if err != nil {
			return cty.NilVal, err
		}
		if len(elems) == 0 {
			return cty.ListValEmpty(cty.DynamicPseudoType), nil
		}
		elems, err = unifyListElementTypes(elems)
		if err != nil {
			return cty.NilVal, err
		}
		return cty.ListVal(elems), nil
	case KindMap:
		m := make(map[string]cty.Value, len(v.Map))
		keys := make([]string, 0, len(v.Map))
		for k := range v.Map {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			c, err := ValueToCty(v.Map[k])
			if err != nil {
				return cty.NilVal, fmt.Errorf("map.%s: %w", k, err)
			}
			m[k] = c
		}
		if len(m) == 0 {
			return cty.MapValEmpty(cty.DynamicPseudoType), nil
		}
		return cty.ObjectVal(m), nil
	default:
		return cty.NullVal(cty.DynamicPseudoType), nil
	}
}

// normalizeAndPruneList normalizes list-of-maps (union keys), prunes nulls, and re-normalizes.
// cty requires all list elements to have the same type; optional/omitted attributes cause object types to differ.
func normalizeAndPruneList(v Value) Value {
	v = normalizeListOfMaps(v)
	v = PruneNulls(v)
	v = normalizeListOfMaps(v)
	return v
}

// listToCtyElems converts each list element to cty.Value, pruning nulls from nested maps.
func listToCtyElems(list []Value) ([]cty.Value, error) {
	elems := make([]cty.Value, 0, len(list))
	for i, el := range list {
		el = PruneNulls(el)
		c, err := ValueToCty(el)
		if err != nil {
			return nil, fmt.Errorf("list[%d]: %w", i, err)
		}
		elems = append(elems, c)
	}
	return elems, nil
}

// unifyListElementTypes ensures all elements have the same cty type so cty.ListVal succeeds.
// Same keys can yield different attribute types (e.g. null -> DynamicPseudoType vs string -> String).
func unifyListElementTypes(elems []cty.Value) ([]cty.Value, error) {
	if len(elems) <= 1 {
		return elems, nil
	}
	types := make([]cty.Type, len(elems))
	for i := range elems {
		types[i] = elems[i].Type()
	}
	unified, convs := convert.UnifyUnsafe(types)
	if unified != cty.NilType && convs != nil {
		for j := range elems {
			if convs[j] != nil {
				converted, err := convs[j](elems[j])
				if err != nil {
					return nil, fmt.Errorf("list[%d]: %w", j, err)
				}
				elems[j] = converted
			}
		}
	}
	if unified != cty.NilType && allSameType(elems) {
		return elems, nil
	}
	// Fallback: convert each element to DynamicPseudoType so list type is consistent.
	for j := range elems {
		c, err := convert.Convert(elems[j], cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("list[%d]: %w", j, err)
		}
		elems[j] = c
	}
	return elems, nil
}

// normalizeListOfMaps ensures every element in a list of maps has the same set of keys (union).
// Missing keys are set to null. This allows cty.ListVal to succeed (all elements same object type).
// Nested lists of maps are normalized recursively.
func normalizeListOfMaps(v Value) Value {
	if v.Kind != KindList || len(v.List) == 0 {
		return v
	}
	allMaps := true
	for _, el := range v.List {
		if el.Kind != KindMap {
			allMaps = false
			break
		}
	}
	if !allMaps {
		// Recursively normalize nested lists (e.g. list of maps containing lists).
		out := make([]Value, len(v.List))
		for i, el := range v.List {
			if el.Kind == KindList {
				out[i] = normalizeListOfMaps(el)
			} else if el.Kind == KindMap && el.Map != nil {
				m := make(map[string]Value, len(el.Map))
				for k, val := range el.Map {
					m[k] = normalizeListOfMaps(val)
				}
				out[i] = Value{Kind: KindMap, Map: m}
			} else {
				out[i] = el
			}
		}
		return Value{Kind: KindList, List: out}
	}
	// Collect union of keys.
	keySet := make(map[string]struct{})
	for _, el := range v.List {
		for k := range el.Map {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	nullVal := Value{Kind: KindNull}
	out := make([]Value, 0, len(v.List))
	for _, el := range v.List {
		m := make(map[string]Value, len(keys))
		for _, k := range keys {
			if val, ok := el.Map[k]; ok {
				m[k] = normalizeListOfMaps(val)
			} else {
				m[k] = nullVal
			}
		}
		out = append(out, Value{Kind: KindMap, Map: m})
	}
	return Value{Kind: KindList, List: out}
}

func allSameType(elems []cty.Value) bool {
	if len(elems) <= 1 {
		return true
	}
	t := elems[0].Type()
	for i := 1; i < len(elems); i++ {
		if !elems[i].Type().Equals(t) {
			return false
		}
	}
	return true
}
