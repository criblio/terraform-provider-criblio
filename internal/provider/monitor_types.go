// Hand-written: do not regenerate (listed in .codegen-ignore).
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ = context.Background
var _ = jsontypes.NormalizedType{}

// MonitorModel is the shared Terraform state model for the criblio_monitor resource
// and data source. It maps to the IAetosMonitorConf API schema returned by
// GET/POST/PATCH /products/aetos/monitors[/{id}].
type MonitorModel struct {
	// Identifying fields
	ID   types.String `tfsdk:"id" json:"id,omitempty"`
	Name types.String `tfsdk:"name" json:"name,omitempty"`

	// Lifecycle control
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,omitempty"`
	Type    types.String `tfsdk:"type" json:"type,omitempty"`

	// JSON blob fields — stored as jsontypes.Normalized so Terraform normalises
	// whitespace and key order before diffing.
	Priority        jsontypes.Normalized `tfsdk:"priority" json:"priority,omitempty"`
	Team            jsontypes.Normalized `tfsdk:"team" json:"team,omitempty"`
	Query           jsontypes.Normalized `tfsdk:"query" json:"query,omitempty"`
	Expr            jsontypes.Normalized `tfsdk:"expr" json:"expr,omitempty"`
	FiringCondition jsontypes.Normalized `tfsdk:"firing_condition" json:"firing_condition,omitempty"`
	FiringRule      jsontypes.Normalized `tfsdk:"firing_rule" json:"firing_rule,omitempty"`
	Metadata        jsontypes.Normalized `tfsdk:"metadata" json:"metadata,omitempty"`
	Notification    jsontypes.Normalized `tfsdk:"notification" json:"notification,omitempty"`

	// Scalar optional fields
	Description     types.String `tfsdk:"description" json:"description,omitempty"`
	DatasetID       types.String `tfsdk:"dataset_id" json:"dataset_id,omitempty"`
	DetectionConfig types.String `tfsdk:"detection_config" json:"detection_config,omitempty"`
	Unit            types.String `tfsdk:"unit" json:"unit,omitempty"`

	// List field
	Silence types.List `tfsdk:"silence" json:"silence,omitempty"` // List(String)

	// Computed-only — stamped by the backend when the User-Agent contains the
	// terraform SDK marker; not settable by the user.
	ManagedBy types.String `tfsdk:"managed_by" json:"managed_by,omitempty"`
}

// MonitorAPIModel is the plain-Go representation used for JSON
// marshal/unmarshal to/from the Aetos API.
type MonitorAPIModel struct {
	ID              *string `json:"id,omitempty"`
	Name            *string `json:"name,omitempty"`
	Enabled         *bool   `json:"enabled,omitempty"`
	Type            *string `json:"type,omitempty"`
	Priority        any     `json:"priority,omitempty"`
	Team            any     `json:"team,omitempty"`
	Query           any     `json:"query,omitempty"`
	Expr            any     `json:"expr,omitempty"`
	FiringCondition any     `json:"firing_condition,omitempty"`
	FiringRule      any     `json:"firing_rule,omitempty"`
	Metadata        any     `json:"metadata,omitempty"`
	Notification    any     `json:"notification,omitempty"`
	Description     *string `json:"description,omitempty"`
	DatasetID       *string `json:"dataset_id,omitempty"`
	DetectionConfig *string `json:"detection_config,omitempty"`
	Unit            *string `json:"unit,omitempty"`
	Silence         []string `json:"silence,omitempty"`
	ManagedBy       *string `json:"managed_by,omitempty"`
}

// MarshalJSON serialises MonitorModel to the Aetos API wire format.
func (m MonitorModel) MarshalJSON() ([]byte, error) {
	output := map[string]any{}

	if !m.ID.IsNull() && !m.ID.IsUnknown() {
		output["id"] = m.ID.ValueString()
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		output["name"] = m.Name.ValueString()
	}
	if !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		output["enabled"] = m.Enabled.ValueBool()
	}
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		output["type"] = m.Type.ValueString()
	}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		output["description"] = m.Description.ValueString()
	}
	if !m.DatasetID.IsNull() && !m.DatasetID.IsUnknown() {
		output["dataset_id"] = m.DatasetID.ValueString()
	}
	if !m.DetectionConfig.IsNull() && !m.DetectionConfig.IsUnknown() {
		output["detection_config"] = m.DetectionConfig.ValueString()
	}
	if !m.Unit.IsNull() && !m.Unit.IsUnknown() {
		output["unit"] = m.Unit.ValueString()
	}

	// JSON blob fields — unmarshal the normalised JSON string back to any so
	// the API receives a real object/array, not a string.
	for apiKey, field := range map[string]jsontypes.Normalized{
		"priority":         m.Priority,
		"team":             m.Team,
		"query":            m.Query,
		"expr":             m.Expr,
		"firing_condition": m.FiringCondition,
		"firing_rule":      m.FiringRule,
		"metadata":         m.Metadata,
		"notification":     m.Notification,
	} {
		if field.IsNull() || field.IsUnknown() {
			continue
		}
		raw := field.ValueString()
		if raw == "" {
			continue
		}
		var parsed any
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			return nil, fmt.Errorf("marshal %s: %v", apiKey, err)
		}
		output[apiKey] = parsed
	}

	// silence is a list of strings
	if !m.Silence.IsNull() && !m.Silence.IsUnknown() {
		elems := m.Silence.Elements()
		silenceList := make([]string, 0, len(elems))
		for _, e := range elems {
			if sv, ok := e.(types.String); ok {
				silenceList = append(silenceList, sv.ValueString())
			}
		}
		output["silence"] = silenceList
	}

	return json.Marshal(output)
}

// UnmarshalJSON deserialises a wire MonitorAPIModel into MonitorModel.
func (m *MonitorModel) UnmarshalJSON(data []byte) error {
	var input MonitorAPIModel
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if input.ID != nil {
		m.ID = types.StringValue(*input.ID)
	} else {
		m.ID = types.StringNull()
	}
	if input.Name != nil {
		m.Name = types.StringValue(*input.Name)
	} else {
		m.Name = types.StringNull()
	}
	if input.Enabled != nil {
		m.Enabled = types.BoolValue(*input.Enabled)
	} else {
		m.Enabled = types.BoolNull()
	}
	if input.Type != nil {
		m.Type = types.StringValue(*input.Type)
	} else {
		m.Type = types.StringNull()
	}
	if input.Description != nil {
		m.Description = types.StringValue(*input.Description)
	} else {
		m.Description = types.StringNull()
	}
	if input.DatasetID != nil {
		m.DatasetID = types.StringValue(*input.DatasetID)
	} else {
		m.DatasetID = types.StringNull()
	}
	if input.DetectionConfig != nil {
		m.DetectionConfig = types.StringValue(*input.DetectionConfig)
	} else {
		m.DetectionConfig = types.StringNull()
	}
	if input.Unit != nil {
		m.Unit = types.StringValue(*input.Unit)
	} else {
		m.Unit = types.StringNull()
	}
	if input.ManagedBy != nil {
		m.ManagedBy = types.StringValue(*input.ManagedBy)
	} else {
		m.ManagedBy = types.StringNull()
	}

	// JSON blob fields
	toNorm := map[string]any{
		"priority":         input.Priority,
		"team":             input.Team,
		"query":            input.Query,
		"expr":             input.Expr,
		"firing_condition": input.FiringCondition,
		"firing_rule":      input.FiringRule,
		"metadata":         input.Metadata,
		"notification":     input.Notification,
	}
	ptrs := map[string]*jsontypes.Normalized{
		"priority":         &m.Priority,
		"team":             &m.Team,
		"query":            &m.Query,
		"expr":             &m.Expr,
		"firing_condition": &m.FiringCondition,
		"firing_rule":      &m.FiringRule,
		"metadata":         &m.Metadata,
		"notification":     &m.Notification,
	}
	for key, val := range toNorm {
		ptr := ptrs[key]
		if val != nil {
			raw, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("unmarshal %s from API: %v", key, err)
			}
			*ptr = jsontypes.NewNormalizedValue(string(raw))
		} else {
			*ptr = jsontypes.NewNormalizedNull()
		}
	}

	// silence list
	if input.Silence != nil {
		val, diags := types.ListValueFrom(context.Background(), types.StringType, input.Silence)
		if diags.HasError() {
			return fmt.Errorf("convert silence from API: %v", diags)
		}
		m.Silence = val
	} else {
		m.Silence = types.ListNull(types.StringType)
	}

	return nil
}

// MonitorTerraformValueToJSON converts a Terraform attr.Value to a JSON-compatible Go value.
// Kept for potential reuse by generated helpers; not used by hand-written monitor code.
func MonitorTerraformValueToJSON(value attr.Value) (any, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	switch typed := value.(type) {
	case types.Bool:
		return typed.ValueBool(), nil
	case types.Int64:
		return typed.ValueInt64(), nil
	case types.Float64:
		return typed.ValueFloat64(), nil
	case jsontypes.Normalized:
		raw := typed.ValueString()
		if raw == "" {
			return map[string]any{}, nil
		}
		var output any
		if err := json.Unmarshal([]byte(raw), &output); err != nil {
			return nil, err
		}
		return output, nil
	case types.String:
		return typed.ValueString(), nil
	case types.List:
		output := make([]any, 0, len(typed.Elements()))
		for _, element := range typed.Elements() {
			val, err := MonitorTerraformValueToJSON(element)
			if err != nil {
				return nil, err
			}
			output = append(output, val)
		}
		return output, nil
	case types.Map:
		output := make(map[string]any, len(typed.Elements()))
		for key, element := range typed.Elements() {
			val, err := MonitorTerraformValueToJSON(element)
			if err != nil {
				return nil, err
			}
			if val == nil {
				continue
			}
			output[key] = val
		}
		return output, nil
	case types.Object:
		output := make(map[string]any, len(typed.Attributes()))
		attributeTypes := typed.AttributeTypes(context.Background())
		for key, attribute := range typed.Attributes() {
			var val any
			var err error
			if attributeType, ok := attributeTypes[key]; ok && attributeType.Equal(jsontypes.NormalizedType{}) {
				val, err = MonitorObjectJSONFromTerraformValue(attribute)
			} else {
				val, err = MonitorTerraformValueToJSON(attribute)
			}
			if err != nil {
				return nil, err
			}
			if val == nil {
				continue
			}
			output[MonitorTerraformNameToAPIName(key)] = val
		}
		return output, nil
	case interface{ ValueString() string }:
		return typed.ValueString(), nil
	default:
		return nil, fmt.Errorf("unsupported Terraform value %T", value)
	}
}

// MonitorObjectJSONFromTerraformValue unmarshals a jsontypes.Normalized string into any.
func MonitorObjectJSONFromTerraformValue(value attr.Value) (any, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	typed, ok := value.(interface{ ValueString() string })
	if !ok {
		return nil, fmt.Errorf("expected normalized JSON string, got %T", value)
	}
	raw := typed.ValueString()
	if raw == "" {
		return map[string]any{}, nil
	}
	var output any
	if err := json.Unmarshal([]byte(raw), &output); err != nil {
		return nil, err
	}
	return output, nil
}

// MonitorTerraformNameToAPIName converts snake_case to camelCase for API field names.
func MonitorTerraformNameToAPIName(name string) string {
	prefix := ""
	if strings.HasPrefix(name, "__template_") {
		prefix = "__template_"
		name = strings.TrimPrefix(name, prefix)
	}
	var output strings.Builder
	upperNext := false
	for _, char := range name {
		if char == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			output.WriteRune(unicode.ToUpper(char))
			upperNext = false
			continue
		}
		output.WriteRune(char)
	}
	return prefix + output.String()
}

// MonitorAPIValueToTerraformValue converts a JSON-decoded Go value to a Terraform attr.Value.
func MonitorAPIValueToTerraformValue(value any, typ attr.Type) (attr.Value, error) {
	if value == nil {
		return MonitorTerraformNullValue(typ)
	}
	if typ.Equal(types.BoolType) {
		typed, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", value)
		}
		return types.BoolValue(typed), nil
	}
	if typ.Equal(types.Int64Type) {
		typed, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("expected number, got %T", value)
		}
		return types.Int64Value(int64(typed)), nil
	}
	if typ.Equal(types.Float64Type) {
		typed, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("expected number, got %T", value)
		}
		return types.Float64Value(typed), nil
	}
	if typ.Equal(types.StringType) {
		typed, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", value)
		}
		return types.StringValue(typed), nil
	}
	if typ.Equal(jsontypes.NormalizedType{}) {
		if typed, ok := value.(string); ok {
			return jsontypes.NewNormalizedValue(typed), nil
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		return jsontypes.NewNormalizedValue(string(raw)), nil
	}
	switch typed := typ.(type) {
	case types.ListType:
		input, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("expected list, got %T", value)
		}
		output := make([]attr.Value, 0, len(input))
		for _, item := range input {
			nested, err := MonitorAPIValueToTerraformValue(item, typed.ElemType)
			if err != nil {
				return nil, err
			}
			output = append(output, nested)
		}
		val, diags := types.ListValue(typed.ElemType, output)
		if diags.HasError() {
			return nil, fmt.Errorf("%v", diags)
		}
		return val, nil
	case types.MapType:
		input, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map, got %T", value)
		}
		output := make(map[string]attr.Value, len(input))
		for key, item := range input {
			nested, err := MonitorAPIValueToTerraformValue(item, typed.ElemType)
			if err != nil {
				return nil, err
			}
			output[key] = nested
		}
		val, diags := types.MapValue(typed.ElemType, output)
		if diags.HasError() {
			return nil, fmt.Errorf("%v", diags)
		}
		return val, nil
	case types.ObjectType:
		input, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object, got %T", value)
		}
		output := make(map[string]attr.Value, len(typed.AttrTypes))
		for key, attrType := range typed.AttrTypes {
			apiKey := MonitorTerraformNameToAPIName(key)
			item, ok := input[apiKey]
			if !ok {
				item, ok = input[key]
			}
			if !ok {
				nested, err := MonitorTerraformNullValue(attrType)
				if err != nil {
					return nil, err
				}
				output[key] = nested
				continue
			}
			nested, err := MonitorAPIValueToTerraformValue(item, attrType)
			if err != nil {
				return nil, err
			}
			output[key] = nested
		}
		val, diags := types.ObjectValue(typed.AttrTypes, output)
		if diags.HasError() {
			return nil, fmt.Errorf("%v", diags)
		}
		return val, nil
	default:
		return nil, fmt.Errorf("unsupported Terraform type %T", typ)
	}
}

// MonitorTerraformNullValue returns the null value for the given Terraform type.
func MonitorTerraformNullValue(typ attr.Type) (attr.Value, error) {
	if typ.Equal(types.BoolType) {
		return types.BoolNull(), nil
	}
	if typ.Equal(types.Int64Type) {
		return types.Int64Null(), nil
	}
	if typ.Equal(types.Float64Type) {
		return types.Float64Null(), nil
	}
	if typ.Equal(types.StringType) {
		return types.StringNull(), nil
	}
	if typ.Equal(jsontypes.NormalizedType{}) {
		return jsontypes.NewNormalizedNull(), nil
	}
	switch typed := typ.(type) {
	case types.ListType:
		return types.ListNull(typed.ElemType), nil
	case types.MapType:
		return types.MapNull(typed.ElemType), nil
	case types.ObjectType:
		return types.ObjectNull(typed.AttrTypes), nil
	default:
		return nil, fmt.Errorf("unsupported Terraform type %T", typ)
	}
}

// MonitorDebug formats a value for debug output.
func MonitorDebug(value any) string {
	return fmt.Sprintf("%v", value)
}
