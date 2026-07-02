// Schema aligned with IAetosMonitorConf (packages/metrics-types/src/shared/types.ts).
// NOT code-generated — see monitor_resource.go for context.
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

// MonitorModel is the Terraform state representation of an IAetosMonitorConf.
// Field names use snake_case matching the schema attribute names in monitor_resource.go.
type MonitorModel struct {
	ID              types.String         `tfsdk:"id"`
	Name            types.String         `tfsdk:"name"`
	Enabled         types.Bool           `tfsdk:"enabled"`
	Type            types.String         `tfsdk:"type"`
	Description     types.String         `tfsdk:"description"`
	Priority        jsontypes.Normalized `tfsdk:"priority"`
	Team            jsontypes.Normalized `tfsdk:"team"`
	DatasetId       types.String         `tfsdk:"dataset_id"`
	Query           jsontypes.Normalized `tfsdk:"query"`
	Expr            jsontypes.Normalized `tfsdk:"expr"`
	FiringCondition jsontypes.Normalized `tfsdk:"firing_condition"`
	FiringRule      jsontypes.Normalized `tfsdk:"firing_rule"`
	Metadata        jsontypes.Normalized `tfsdk:"metadata"`
	Notification    jsontypes.Normalized `tfsdk:"notification"`
	Silence         types.List           `tfsdk:"silence"`
	DetectionConfig jsontypes.Normalized `tfsdk:"detection_config"`
	Unit            types.String         `tfsdk:"unit"`
	ManagedBy       types.String         `tfsdk:"managed_by"`
}

// MonitorResourceModel is an alias used by codegen-adjacent utilities.
type MonitorResourceModel = MonitorModel

// MonitorDataSourceModel is an alias — data sources share the same struct.
type MonitorDataSourceModel = MonitorModel

// MonitorAPIModel is the raw JSON shape of IAetosMonitorConf as returned/accepted
// by POST/PATCH/GET /products/aetos/monitors.
type MonitorAPIModel struct {
	ID              *string  `json:"id,omitempty"`
	Name            *string  `json:"name,omitempty"`
	Enabled         *bool    `json:"enabled,omitempty"`
	Type            *string  `json:"type,omitempty"`
	Description     *string  `json:"description,omitempty"`
	Priority        any      `json:"priority,omitempty"`
	Team            any      `json:"team,omitempty"`
	DatasetId       *string  `json:"datasetId,omitempty"`
	Query           any      `json:"query,omitempty"`
	Expr            any      `json:"expr,omitempty"`
	FiringCondition any      `json:"firingCondition,omitempty"`
	FiringRule      any      `json:"firingRule,omitempty"`
	Metadata        any      `json:"metadata,omitempty"`
	Notification    any      `json:"notification,omitempty"`
	Silence         []string `json:"silence,omitempty"`
	DetectionConfig any      `json:"detectionConfig,omitempty"`
	Unit            *string  `json:"unit,omitempty"`
	ManagedBy       *string  `json:"managedBy,omitempty"`
}

// ── JSON serialization helpers ───────────────────────────────────────────────

// MonitorTerraformValueToJSON converts a Terraform attr.Value to a Go value
// suitable for json.Marshal.
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
			value, err := MonitorTerraformValueToJSON(element)
			if err != nil {
				return nil, err
			}
			output = append(output, value)
		}
		return output, nil
	case types.Map:
		output := make(map[string]any, len(typed.Elements()))
		for key, element := range typed.Elements() {
			value, err := MonitorTerraformValueToJSON(element)
			if err != nil {
				return nil, err
			}
			if value == nil {
				continue
			}
			output[key] = value
		}
		return output, nil
	case types.Object:
		output := make(map[string]any, len(typed.Attributes()))
		attributeTypes := typed.AttributeTypes(context.Background())
		for key, attribute := range typed.Attributes() {
			var value any
			var err error
			if attributeType, ok := attributeTypes[key]; ok && attributeType.Equal(jsontypes.NormalizedType{}) {
				value, err = MonitorObjectJSONFromTerraformValue(attribute)
			} else {
				value, err = MonitorTerraformValueToJSON(attribute)
			}
			if err != nil {
				return nil, err
			}
			if value == nil {
				continue
			}
			output[MonitorTerraformNameToAPIName(key)] = value
		}
		return output, nil
	case interface{ ValueString() string }:
		return typed.ValueString(), nil
	default:
		return nil, fmt.Errorf("unsupported Terraform value %T", value)
	}
}

// MonitorObjectJSONFromTerraformValue parses a jsontypes.Normalized string field
// into a Go value, suitable for embedding as a nested object in the API payload.
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

// MonitorTerraformNameToAPIName converts a snake_case TF attribute name to the
// camelCase JSON key the API expects (e.g. "firing_condition" → "firingCondition").
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

// MonitorAPIValueToTerraformValue converts a JSON-decoded Go value to the
// Terraform attr.Value type described by typ.
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
		value, diags := types.ListValue(typed.ElemType, output)
		if diags.HasError() {
			return nil, fmt.Errorf("%v", diags)
		}
		return value, nil
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
		value, diags := types.MapValue(typed.ElemType, output)
		if diags.HasError() {
			return nil, fmt.Errorf("%v", diags)
		}
		return value, nil
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
		value, diags := types.ObjectValue(typed.AttrTypes, output)
		if diags.HasError() {
			return nil, fmt.Errorf("%v", diags)
		}
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported Terraform type %T", typ)
	}
}

// MonitorTerraformNullValue returns the appropriate null value for a given type.
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

// ── MonitorModel JSON Marshal/Unmarshal ──────────────────────────────────────

// MarshalJSON serializes MonitorModel to the JSON payload expected by
// /products/aetos/monitors. managed_by is intentionally omitted — it is
// computed (read-only) and must not be sent to the API.
func (m MonitorModel) MarshalJSON() ([]byte, error) {
	output := map[string]any{}

	if !m.ID.IsNull() && !m.ID.IsUnknown() {
		v, err := MonitorTerraformValueToJSON(m.ID)
		if err != nil {
			return nil, fmt.Errorf("convert id: %v", err)
		}
		output["id"] = v
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		v, err := MonitorTerraformValueToJSON(m.Name)
		if err != nil {
			return nil, fmt.Errorf("convert name: %v", err)
		}
		output["name"] = v
	}
	if !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		v, err := MonitorTerraformValueToJSON(m.Enabled)
		if err != nil {
			return nil, fmt.Errorf("convert enabled: %v", err)
		}
		output["enabled"] = v
	}
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		v, err := MonitorTerraformValueToJSON(m.Type)
		if err != nil {
			return nil, fmt.Errorf("convert type: %v", err)
		}
		output["type"] = v
	}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		v, err := MonitorTerraformValueToJSON(m.Description)
		if err != nil {
			return nil, fmt.Errorf("convert description: %v", err)
		}
		output["description"] = v
	}
	// JSON object fields — parse the stored JSON string → native Go value
	for tfName, jsName := range map[string]string{
		"priority":         "priority",
		"team":             "team",
		"query":            "query",
		"expr":             "expr",
		"firingCondition":  "firingCondition",
		"firingRule":       "firingRule",
		"metadata":         "metadata",
		"notification":     "notification",
		"detectionConfig":  "detectionConfig",
	} {
		var field jsontypes.Normalized
		switch tfName {
		case "priority":
			field = m.Priority
		case "team":
			field = m.Team
		case "query":
			field = m.Query
		case "expr":
			field = m.Expr
		case "firingCondition":
			field = m.FiringCondition
		case "firingRule":
			field = m.FiringRule
		case "metadata":
			field = m.Metadata
		case "notification":
			field = m.Notification
		case "detectionConfig":
			field = m.DetectionConfig
		}
		if !field.IsNull() && !field.IsUnknown() {
			v, err := MonitorObjectJSONFromTerraformValue(field)
			if err != nil {
				return nil, fmt.Errorf("convert %s: %v", jsName, err)
			}
			if v != nil {
				output[jsName] = v
			}
		}
	}
	if !m.DatasetId.IsNull() && !m.DatasetId.IsUnknown() {
		output["datasetId"] = m.DatasetId.ValueString()
	}
	if !m.Unit.IsNull() && !m.Unit.IsUnknown() {
		output["unit"] = m.Unit.ValueString()
	}
	if !m.Silence.IsNull() && !m.Silence.IsUnknown() {
		v, err := MonitorTerraformValueToJSON(m.Silence)
		if err != nil {
			return nil, fmt.Errorf("convert silence: %v", err)
		}
		output["silence"] = v
	}
	// managed_by is NOT sent — it is computed/read-only.

	return json.Marshal(output)
}

// UnmarshalJSON reads an IAetosMonitorConf API response into a MonitorModel.
func (m *MonitorModel) UnmarshalJSON(data []byte) error {
	var input MonitorAPIModel
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	// Simple string / bool fields
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
	if input.DatasetId != nil {
		m.DatasetId = types.StringValue(*input.DatasetId)
	} else {
		m.DatasetId = types.StringNull()
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

	// JSON object fields — marshal the Go value back to a JSON string for TF state
	jsonField := func(v any) (jsontypes.Normalized, error) {
		if v == nil {
			return jsontypes.NewNormalizedNull(), nil
		}
		raw, err := json.Marshal(v)
		if err != nil {
			return jsontypes.NewNormalizedNull(), err
		}
		return jsontypes.NewNormalizedValue(string(raw)), nil
	}

	var err error
	if m.Priority, err = jsonField(input.Priority); err != nil {
		return fmt.Errorf("unmarshal priority: %v", err)
	}
	if m.Team, err = jsonField(input.Team); err != nil {
		return fmt.Errorf("unmarshal team: %v", err)
	}
	if m.Query, err = jsonField(input.Query); err != nil {
		return fmt.Errorf("unmarshal query: %v", err)
	}
	if m.Expr, err = jsonField(input.Expr); err != nil {
		return fmt.Errorf("unmarshal expr: %v", err)
	}
	if m.FiringCondition, err = jsonField(input.FiringCondition); err != nil {
		return fmt.Errorf("unmarshal firingCondition: %v", err)
	}
	if m.FiringRule, err = jsonField(input.FiringRule); err != nil {
		return fmt.Errorf("unmarshal firingRule: %v", err)
	}
	if m.Metadata, err = jsonField(input.Metadata); err != nil {
		return fmt.Errorf("unmarshal metadata: %v", err)
	}
	if m.Notification, err = jsonField(input.Notification); err != nil {
		return fmt.Errorf("unmarshal notification: %v", err)
	}
	if m.DetectionConfig, err = jsonField(input.DetectionConfig); err != nil {
		return fmt.Errorf("unmarshal detectionConfig: %v", err)
	}

	// silence: []string → types.List
	if input.Silence != nil {
		value, diags := types.ListValueFrom(context.Background(), types.StringType, input.Silence)
		if diags.HasError() {
			return fmt.Errorf("unmarshal silence: %v", diags)
		}
		m.Silence = value
	} else {
		m.Silence = types.ListNull(types.StringType)
	}

	return nil
}
