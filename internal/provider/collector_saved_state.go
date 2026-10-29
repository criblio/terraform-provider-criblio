package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Decode checkpoint numbers independently so database IDs retain their precision.
func decodeCollectorSavedState(data []byte, raw map[string]any) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	state, ok := fields["savedState"]
	if !ok {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(state))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return fmt.Errorf("decode savedState: %w", err)
	}
	raw["savedState"] = value
	return nil
}

func dynamicValueFromAPI(value any) (types.Dynamic, error) {
	if value == nil {
		return types.DynamicNull(), nil
	}
	underlying, err := dynamicUnderlyingValue(value)
	if err != nil {
		return types.DynamicNull(), err
	}
	return types.DynamicValue(underlying), nil
}

func dynamicCollectionValue(values []attr.Value) (types.Dynamic, error) {
	raw, err := dynamicArrayToJSON(values)
	if err != nil {
		return types.DynamicNull(), err
	}
	return dynamicValueFromAPI(raw)
}

func dynamicUnderlyingValue(value any) (attr.Value, error) {
	switch value := value.(type) {
	case nil:
		return types.DynamicNull(), nil
	case string:
		return types.StringValue(value), nil
	case bool:
		return types.BoolValue(value), nil
	case json.Number:
		number, _, err := big.ParseFloat(string(value), 10, 512, big.ToNearestEven)
		if err != nil {
			return nil, fmt.Errorf("decode checkpoint number: %w", err)
		}
		return types.NumberValue(number), nil
	case float64:
		return types.Float64Value(value), nil
	case int64:
		return types.Int64Value(value), nil
	case int:
		return types.Int64Value(int64(value)), nil
	case map[string]any:
		values := make(map[string]attr.Value, len(value))
		attributeTypes := make(map[string]attr.Type, len(value))
		for key, item := range value {
			nested, err := dynamicUnderlyingValue(item)
			if err != nil {
				return nil, fmt.Errorf("decode checkpoint key %q: %w", key, err)
			}
			values[key] = nested
			attributeTypes[key] = nested.Type(context.Background())
		}
		result, diags := types.ObjectValue(attributeTypes, values)
		if diags.HasError() {
			return nil, fmt.Errorf("decode checkpoint object: %v", diags)
		}
		return result, nil
	case []any:
		values := make([]attr.Value, len(value))
		elementTypes := make([]attr.Type, len(value))
		for index, item := range value {
			nested, err := dynamicUnderlyingValue(item)
			if err != nil {
				return nil, fmt.Errorf("decode checkpoint index %d: %w", index, err)
			}
			values[index] = nested
			elementTypes[index] = nested.Type(context.Background())
		}
		result, diags := types.TupleValue(elementTypes, values)
		if diags.HasError() {
			return nil, fmt.Errorf("decode checkpoint array: %v", diags)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported checkpoint value %T", value)
	}
}

func dynamicValueToJSON(value attr.Value) (any, error) {
	if value == nil || value.IsNull() {
		return nil, nil
	}
	if value.IsUnknown() {
		return nil, fmt.Errorf("checkpoint contains an unknown value")
	}
	switch value := value.(type) {
	case types.Dynamic:
		return dynamicValueToJSON(value.UnderlyingValue())
	case types.String:
		return value.ValueString(), nil
	case types.Bool:
		return value.ValueBool(), nil
	case types.Number:
		return json.Number(value.ValueBigFloat().Text('f', -1)), nil
	case types.Int64:
		return value.ValueInt64(), nil
	case types.Float64:
		return value.ValueFloat64(), nil
	case types.Object:
		return dynamicObjectToJSON(value.Attributes())
	case types.Map:
		return dynamicObjectToJSON(value.Elements())
	case types.Tuple:
		return dynamicArrayToJSON(value.Elements())
	case types.List:
		return dynamicArrayToJSON(value.Elements())
	case types.Set:
		return dynamicArrayToJSON(value.Elements())
	default:
		return nil, fmt.Errorf("unsupported checkpoint type %T", value)
	}
}

func dynamicObjectToJSON(values map[string]attr.Value) (map[string]any, error) {
	result := make(map[string]any, len(values))
	for key, value := range values {
		converted, err := dynamicValueToJSON(value)
		if err != nil {
			return nil, err
		}
		result[key] = converted
	}
	return result, nil
}

func dynamicArrayToJSON(values []attr.Value) ([]any, error) {
	result := make([]any, len(values))
	for index, value := range values {
		converted, err := dynamicValueToJSON(value)
		if err != nil {
			return nil, err
		}
		result[index] = converted
	}
	return result, nil
}

// Shared variants retain the same saved_state semantics for group and pack collectors.
func collectorSavedStates(model any) map[string]types.Dynamic {
	result := make(map[string]types.Dynamic)
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	for index := 0; index < value.NumField(); index++ {
		name := value.Type().Field(index).Tag.Get("tfsdk")
		if !strings.HasPrefix(name, "input_collector_") || value.Field(index).IsNil() {
			continue
		}
		result[name] = value.Field(index).Elem().FieldByName("SavedState").Interface().(types.Dynamic)
	}
	return result
}

func restoreCollectorConfiguredSavedState(model any, prior map[string]types.Dynamic, configured map[string]bool) {
	value := reflect.ValueOf(model).Elem()
	for index := 0; index < value.NumField(); index++ {
		name := value.Type().Field(index).Tag.Get("tfsdk")
		if !configured[name] || value.Field(index).IsNil() {
			continue
		}
		if state, ok := prior[name]; ok {
			value.Field(index).Elem().FieldByName("SavedState").Set(reflect.ValueOf(state))
		}
	}
}

// Clone variant pointers before stripping computed checkpoints from the request.
func collectorRequestSavedState[T any](model T, configured map[string]bool) T {
	value := reflect.ValueOf(&model).Elem()
	for index := 0; index < value.NumField(); index++ {
		name := value.Type().Field(index).Tag.Get("tfsdk")
		if !strings.HasPrefix(name, "input_collector_") || value.Field(index).IsNil() || configured[name] {
			continue
		}
		variant := reflect.New(value.Field(index).Type().Elem())
		variant.Elem().Set(value.Field(index).Elem())
		variant.Elem().FieldByName("SavedState").Set(reflect.ValueOf(types.DynamicNull()))
		value.Field(index).Set(variant)
	}
	return model
}

func collectorConfiguredSavedState(ctx context.Context, config tfsdk.Config) (map[string]bool, diag.Diagnostics) {
	configured := make(map[string]bool)
	var object types.Object
	diags := config.Get(ctx, &object)
	if diags.HasError() || object.IsNull() || object.IsUnknown() {
		return configured, diags
	}
	for name, value := range object.Attributes() {
		if !strings.HasPrefix(name, "input_collector_") || value.IsNull() {
			continue
		}
		if value.IsUnknown() {
			configured[name] = true
			continue
		}
		state := value.(types.Object).Attributes()["saved_state"].(types.Dynamic)
		configured[name] = !state.IsNull() && !state.IsUnderlyingValueNull()
		if !configured[name] || state.IsUnknown() || state.IsUnderlyingValueUnknown() {
			continue
		}
		switch state.UnderlyingValue().(type) {
		case types.Map, types.Object:

		default:
			diags.AddError("Invalid saved_state", name+".saved_state must be a map or object.")
		}
	}
	return configured, diags
}

type collectorSavedStatePrivate interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
	SetKey(context.Context, string, []byte) diag.Diagnostics
}

const collectorSavedStateConfigKey = "collector_saved_state_configured"

// Removing an Optional+Computed checkpoint otherwise retains its previous plan
// value and skips Update, leaving the old configuration ownership in private state.
func modifyCollectorSavedStatePlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}
	configured, diags := collectorConfiguredSavedState(ctx, req.Config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var priorObject types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &priorObject)...)
	if resp.Diagnostics.HasError() {
		return
	}
	prior := make(map[string]types.Dynamic)
	for name, value := range priorObject.Attributes() {
		if !strings.HasPrefix(name, "input_collector_") || value.IsNull() || value.IsUnknown() {
			continue
		}
		prior[name] = value.(types.Object).Attributes()["saved_state"].(types.Dynamic)
	}
	previouslyConfigured, diags := readCollectorSavedStateConfig(ctx, req.Private, prior)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	for name, wasConfigured := range previouslyConfigured {
		if !wasConfigured || configured[name] {
			continue
		}
		var block types.Object
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(name), &block)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if block.IsNull() || block.IsUnknown() {
			continue
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(name).AtName("saved_state"), types.DynamicUnknown())...)
	}
}

func writeCollectorSavedStateConfig(ctx context.Context, private collectorSavedStatePrivate, configured map[string]bool) diag.Diagnostics {
	data, err := json.Marshal(configured)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Unable to record saved_state configuration", err.Error())}
	}
	return private.SetKey(ctx, collectorSavedStateConfigKey, data)
}

func readCollectorSavedStateConfig(ctx context.Context, private collectorSavedStatePrivate, prior map[string]types.Dynamic) (map[string]bool, diag.Diagnostics) {
	configured := make(map[string]bool)
	data, diags := private.GetKey(ctx, collectorSavedStateConfigKey)
	if diags.HasError() {
		return configured, diags
	}
	if len(data) != 0 {
		if err := json.Unmarshal(data, &configured); err != nil {
			diags.AddError("Unable to read saved_state configuration", err.Error())
		}
		return configured, diags
	}
	// Older providers did not record configuration provenance. Retain legacy
	// string maps until a create/update records whether they were configured.
	for name, value := range prior {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		if legacy, ok := value.UnderlyingValue().(types.Map); ok {
			configured[name] = legacy.ElementType(ctx).Equal(types.StringType)
		}
	}
	return configured, diags
}
