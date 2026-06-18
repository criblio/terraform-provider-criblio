// Package converter converts SDK API responses into provider ResourceModel instances
// using reflection, provider RefreshFrom* methods, and generated model metadata.
package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// RefreshFromMethodName returns the RefreshFrom* method name for the given SDK Get method.
// Convention: RefreshFromOperations + GetMethod + "ResponseBody" (e.g. GetInputByID -> RefreshFromOperationsGetInputByIDResponseBody).
func RefreshFromMethodName(getMethod string) string {
	if getMethod == "" {
		return ""
	}
	return "RefreshFromOperations" + getMethod + "ResponseBody"
}

// Convert fetches a single resource via the SDK Get* method and converts the response
// into a provider ResourceModel by instantiating the model and calling its RefreshFrom* method.
// requestParams is used to build the Get request (e.g. GroupID, ID, Pack). Keys must match
// the request struct field names (e.g. "GroupID", "ID").
// For criblio_search_dataset and criblio_search_dataset_provider, if the list response was
// cached (SDK union unmarshal failed for cribl_lake), we build the model from the cached
// list item instead of calling Get—same pattern as reusing list response for usage group.
func Convert(ctx context.Context, client *sdk.CriblIo, e registry.Entry, requestParams map[string]string) (model interface{}, err error) {
	if e.GetMethod == "" {
		return nil, fmt.Errorf("%s: no GetMethod in registry", e.TypeName)
	}
	modelTypes := ResourceModelTypes()
	modelType, ok := modelTypes[e.ModelTypeName]
	if !ok {
		return nil, fmt.Errorf("%s: unknown model type %q", e.TypeName, e.ModelTypeName)
	}

	// Call SDK Get* first so we get full response (type-specific blocks like s3_provider, apihttp_provider).
	// For search_dataset/search_dataset_provider, only fall back to cached list item when Get fails
	// (e.g. SDK unmarshal error for cribl_lake); list item has only id/description/type so HCL would be empty.
	respBody, err := callGetByID(ctx, client, e, requestParams)
	// criblio_pack_destination: SDK GetPackOutputByIDResponseBody is empty; use captured raw response (same shape as GetOutputByID).
	// Store raw items[0] so export can emit the oneOf block (e.g. output_cribl_lake) when the model has no Items field.
	if err == nil && e.TypeName == "criblio_pack_destination" {
		if captured := custom.GetAndClearPackOutputGetBody(requestParams["GroupID"], requestParams["Pack"], requestParams["ID"]); len(captured) > 0 {
			var rawResp struct {
				Items []json.RawMessage `json:"items"`
			}
			if jsonErr := json.Unmarshal(captured, &rawResp); jsonErr == nil && len(rawResp.Items) > 0 {
				custom.StorePackOutputFirstItem(requestParams["GroupID"], requestParams["Pack"], requestParams["ID"], rawResp.Items[0])
			}
			var getOut operations.GetOutputByIDResponseBody
			if jsonErr := json.Unmarshal(captured, &getOut); jsonErr == nil {
				respBody = &getOut
			}
		}
	}
	if err == nil && respBody != nil {
		convModelType := modelType
		// criblio_pack_destination: PackDestinationResourceModel has no RefreshFrom for GetOutputByID; use DestinationResourceModel (same response shape).
		if e.TypeName == "criblio_pack_destination" {
			if destType, ok := modelTypes["DestinationResourceModel"]; ok {
				convModelType = destType
			}
		}
		converted, convErr := convertFromResponseBody(ctx, e, convModelType, respBody)
		if convErr == nil {
			if injErr := InjectRequiredIdentifiers(converted, requestParams); injErr != nil {
				return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, injErr)
			}
			// Pack pipeline GET returns Routes (items) only; try to fill conf from lib pipeline or set minimal conf.
			if e.TypeName == "criblio_pack_pipeline" {
				fillPackPipelineConf(ctx, client, converted, requestParams)
			}
			return converted, nil
		}
	}

	// Fallback for search_dataset/search_dataset_provider: build from cached list item when Get failed
	// or conversion failed (e.g. SDK union unmarshal for cribl_lake).
	id := requestParams["ID"]
	if id != "" && (e.TypeName == "criblio_search_dataset" || e.TypeName == "criblio_search_dataset_provider") {
		var pathKey string
		if e.TypeName == "criblio_search_dataset" {
			pathKey = custom.PathSearchDatasets
		} else {
			pathKey = custom.PathSearchDatasetProviders
		}
		if itemJSON := custom.GetCachedSearchListItem(pathKey, id); len(itemJSON) > 0 {
			// For search_dataset, try to get full model from captured Get body or list item (SDK GenericDataset).
			// Get body is captured when Get returned 200 but SDK unmarshal failed (e.g. cribl_lake); list item
			// may have full or summary payload. Try Get body first (full), then list item.
			if e.TypeName == "criblio_search_dataset" {
				tryGenericDataset := func(data []byte) (interface{}, bool) {
					var g shared.GenericDataset
					if err := json.Unmarshal(data, &g); err != nil {
						return nil, false
					}
					modelVal := reflect.New(modelType)
					method := modelVal.MethodByName("RefreshFromSharedGenericDataset")
					if !method.IsValid() {
						return nil, false
					}
					ctx := context.Background()
					outs := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(&g)})
					if len(outs) == 1 && !outs[0].IsNil() {
						if DiagnosticsToError(outs[0].Interface().(diag.Diagnostics), e.TypeName, "") != nil {
							return nil, false
						}
					}
					return modelVal.Interface(), true
				}
				// 1) Try captured Get response body (full single-dataset response from API).
				if getBody := custom.GetAndClearSearchDatasetGetBody(id); len(getBody) > 0 {
					var getResp struct {
						Items []json.RawMessage `json:"items"`
					}
					if err := json.Unmarshal(getBody, &getResp); err == nil && len(getResp.Items) > 0 {
						if model, ok := tryGenericDataset(getResp.Items[0]); ok {
							if injErr := InjectRequiredIdentifiers(model, requestParams); injErr != nil {
								return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, injErr)
							}
							return model, nil
						}
					}
				}
				// 2) Try list item JSON (may be full or summary).
				if model, ok := tryGenericDataset(itemJSON); ok {
					if injErr := InjectRequiredIdentifiers(model, requestParams); injErr != nil {
						return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, injErr)
					}
					return model, nil
				}
			}
			model, buildErr := modelFromSearchItemJSON(modelType, itemJSON)
			if buildErr == nil {
				if injErr := InjectRequiredIdentifiers(model, requestParams); injErr != nil {
					return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, injErr)
				}
				return model, nil
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", e.TypeName, err)
	}
	if respBody == nil {
		return nil, fmt.Errorf("%s: Get response body is nil", e.TypeName)
	}
	// Get succeeded but convertFromResponseBody failed and we had no cache fallback
	converted, convErr := convertFromResponseBody(ctx, e, modelType, respBody)
	if convErr != nil {
		return nil, convErr
	}
	if err = InjectRequiredIdentifiers(converted, requestParams); err != nil {
		return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, err)
	}
	return converted, nil
}

// fillPackPipelineConf fills the pack pipeline model's Conf. GetPipelinesByPackWithID returns the
// pipeline definition for pack pipelines. We need Pack, GroupID, and ID to fetch the correct content.
func fillPackPipelineConf(ctx context.Context, client *sdk.CriblIo, converted interface{}, requestParams map[string]string) {
	pm, ok := converted.(*provider.PackPipelineResourceModel)
	if !ok || client == nil || client.Routes == nil {
		return
	}
	groupID := requestParams["GroupID"]
	id := requestParams["ID"]
	pack := requestParams["Pack"]
	if groupID == "" {
		groupID = "default"
	}
	if id == "" || pack == "" {
		ensureMinimalPackPipelineConf(pm)
		return
	}
	resp, err := client.Routes.GetPipelinesByPackWithID(ctx, operations.GetPipelinesByPackWithIDRequest{
		GroupID: groupID,
		Pack:    pack,
		ID:      id,
	})
	if err != nil || resp == nil || resp.Object == nil {
		ensureMinimalPackPipelineConf(pm)
		return
	}
	items := resp.Object.GetItems()
	if len(items) != 1 {
		ensureMinimalPackPipelineConf(pm)
		return
	}
	diags := pm.RefreshFromSharedPipeline(ctx, &items[0])
	if diags != nil && diags.HasError() {
		ensureMinimalPackPipelineConf(pm)
		return
	}
}

// ensureMinimalPackPipelineConf sets a minimal conf (empty functions, default output) so exported
// HCL has a valid required conf block when the API did not return pipeline definition.
func ensureMinimalPackPipelineConf(pm *provider.PackPipelineResourceModel) {
	if pm == nil {
		return
	}
	// Check if conf is effectively empty (no functions).
	if len(pm.Conf.Functions) == 0 {
		if pm.Conf.Output.ValueString() == "" {
			pm.Conf.Output = types.StringValue("default")
		}
		// Ensure non-nil slice for HCL
		if pm.Conf.Functions == nil {
			pm.Conf.Functions = []tfTypes.PipelineFunctionConf{}
		}
	}
}

// fillPackLookupsContent fetches the pack lookup file (GET /m/{groupID}/p/{pack}/system/lookups/{id},
// ConvertFromResponseBody instantiates the resource model and calls its RefreshFrom* method
// with the given response body. Used when the response body is already available (e.g. tests or
// when the caller has fetched the response). Supports nested and complex structures because
// the provider's RefreshFrom* methods handle the full response body type.
// It does not set required identifiers (id, group_id); use ConvertFromResponseBodyWithIdentifiers
// or call InjectRequiredIdentifiers after conversion when the model is used for HCL generation.
func ConvertFromResponseBody(ctx context.Context, e registry.Entry, responseBody interface{}) (interface{}, error) {
	if e.GetMethod == "" {
		return nil, fmt.Errorf("%s: no GetMethod", e.TypeName)
	}
	modelTypes := ResourceModelTypes()
	modelType, ok := modelTypes[e.ModelTypeName]
	if !ok {
		return nil, fmt.Errorf("%s: unknown model type %q", e.TypeName, e.ModelTypeName)
	}
	return convertFromResponseBody(ctx, e, modelType, responseBody)
}

// ConvertFromResponseBodyWithIdentifiers is like ConvertFromResponseBody but also injects
// required Terraform identifier fields (id, group_id, pack, etc.) from the given map.
// Keys should match request param names (e.g. "ID", "GroupID", "Pack"). Use this when
// the converted model will be used for HCL generation so that required fields are set.
func ConvertFromResponseBodyWithIdentifiers(ctx context.Context, e registry.Entry, responseBody interface{}, identifiers map[string]string) (interface{}, error) {
	model, err := ConvertFromResponseBody(ctx, e, responseBody)
	if err != nil {
		return nil, err
	}
	if err := InjectRequiredIdentifiers(model, identifiers); err != nil {
		return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, err)
	}
	return model, nil
}

// convertFromResponseBody allocates the model and invokes the RefreshFrom* method.
func convertFromResponseBody(ctx context.Context, e registry.Entry, modelType reflect.Type, responseBody interface{}) (interface{}, error) {
	modelVal := reflect.New(modelType)
	refreshMethodName := e.RefreshFromMethod
	if refreshMethodName == "" {
		refreshMethodName = RefreshFromMethodName(e.GetMethod)
	}
	method := modelVal.MethodByName(refreshMethodName)
	if !method.IsValid() {
		converted, err := convertGeneratedModelFromResponseBody(e, modelType, responseBody)
		if err == nil {
			return converted, nil
		}
		return nil, fmt.Errorf("%s: model %s has no method %s", e.TypeName, e.ModelTypeName, refreshMethodName)
	}
	mt := method.Type()
	if mt.NumIn() != 2 || mt.NumOut() != 1 {
		return nil, fmt.Errorf("%s: %s has unexpected signature", e.TypeName, refreshMethodName)
	}
	// Second parameter must accept the response body type (e.g. *operations.GetInputByIDResponseBody).
	respType := mt.In(1)
	respVal := reflect.ValueOf(responseBody)
	if responseBody != nil && !respVal.Type().AssignableTo(respType) {
		return nil, fmt.Errorf("%s: response body type %v not assignable to %s parameter %v", e.TypeName, respVal.Type(), refreshMethodName, respType)
	}
	outs := method.Call([]reflect.Value{reflect.ValueOf(ctx), respVal})
	diags := outs[0].Interface().(diag.Diagnostics)
	if err := DiagnosticsToError(diags, e.TypeName, ""); err != nil {
		return nil, err
	}
	return modelVal.Interface(), nil
}

func convertGeneratedModelFromResponseBody(e registry.Entry, modelType reflect.Type, responseBody interface{}) (interface{}, error) {
	item, err := firstResponseItem(responseBody)
	if err != nil {
		return nil, err
	}
	itemJSON, err := json.Marshal(item.Interface())
	if err != nil {
		return nil, fmt.Errorf("%s: marshal response item: %w", e.TypeName, err)
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(itemJSON, &values); err != nil {
		return nil, fmt.Errorf("%s: decode response item: %w", e.TypeName, err)
	}

	modelVal := reflect.New(modelType)
	if err := populateGeneratedModel(modelVal.Elem(), values); err != nil {
		return nil, fmt.Errorf("%s: populate generated model: %w", e.TypeName, err)
	}
	return modelVal.Interface(), nil
}

func firstResponseItem(responseBody interface{}) (reflect.Value, error) {
	if responseBody == nil {
		return reflect.Value{}, fmt.Errorf("response body is nil")
	}
	respVal := reflect.ValueOf(responseBody)
	if respVal.Kind() == reflect.Ptr {
		if respVal.IsNil() {
			return reflect.Value{}, fmt.Errorf("response body is nil")
		}
		respVal = respVal.Elem()
	}
	if respVal.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("response body %T is not a struct", responseBody)
	}
	items := respVal.FieldByName("Items")
	if !items.IsValid() || items.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("response body has no Items slice")
	}
	if items.Len() == 0 {
		return reflect.Value{}, fmt.Errorf("response body has no items")
	}
	return items.Index(0), nil
}

func populateGeneratedModel(model reflect.Value, values map[string]json.RawMessage) error {
	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)
		if !field.CanSet() {
			continue
		}
		jsonName := model.Type().Field(i).Tag.Get("json")
		jsonName = strings.Split(jsonName, ",")[0]
		if jsonName == "" || jsonName == "-" {
			continue
		}
		raw, ok := values[jsonName]
		if !ok || string(raw) == "null" {
			continue
		}
		if err := setGeneratedModelField(field, raw); err != nil {
			return fmt.Errorf("%s: %w", model.Type().Field(i).Name, err)
		}
	}
	return nil
}

func setGeneratedModelField(field reflect.Value, raw json.RawMessage) error {
	switch field.Type() {
	case reflect.TypeOf(types.String{}):
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		field.Set(reflect.ValueOf(types.StringValue(value)))
		return nil
	case reflect.TypeOf(types.Bool{}):
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		field.Set(reflect.ValueOf(types.BoolValue(value)))
		return nil
	case reflect.TypeOf(types.Int64{}):
		var value int64
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		field.Set(reflect.ValueOf(types.Int64Value(value)))
		return nil
	case reflect.TypeOf(types.Float64{}):
		var value float64
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		field.Set(reflect.ValueOf(types.Float64Value(value)))
		return nil
	case reflect.TypeOf(jsontypes.Normalized{}):
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			value = string(raw)
		}
		field.Set(reflect.ValueOf(jsontypes.NewNormalizedValue(value)))
		return nil
	case reflect.TypeOf(types.List{}):
		var values []string
		if err := json.Unmarshal(raw, &values); err == nil {
			field.Set(reflect.ValueOf(types.ListValueMust(types.StringType, stringValues(values))))
			return nil
		}
		value, err := objectListValue(raw)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(value))
		return nil
	case reflect.TypeOf(types.Object{}):
		value, err := objectValue(raw)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(value))
		return nil
	}

	if field.Type() == reflect.TypeOf([]types.String{}) {
		var values []string
		if err := json.Unmarshal(raw, &values); err != nil {
			return err
		}
		tfValues := make([]types.String, 0, len(values))
		for _, value := range values {
			tfValues = append(tfValues, types.StringValue(value))
		}
		field.Set(reflect.ValueOf(tfValues))
	}
	return nil
}

func objectListValue(raw json.RawMessage) (types.List, error) {
	var objects []map[string]any
	if err := json.Unmarshal(raw, &objects); err != nil {
		return types.List{}, err
	}
	for index := range objects {
		objects[index] = terraformNameMap(objects[index])
	}
	attrTypes := map[string]attr.Type{}
	for _, object := range objects {
		for key, value := range object {
			attrTypes[key] = inferredAttrType(value)
		}
	}
	objectType := types.ObjectType{AttrTypes: attrTypes}
	values := make([]attr.Value, 0, len(objects))
	for _, object := range objects {
		attrs := make(map[string]attr.Value, len(attrTypes))
		for key, typ := range attrTypes {
			if value, ok := object[key]; ok {
				attrValue, err := attrValueFromAny(value, typ)
				if err != nil {
					return types.List{}, err
				}
				attrs[key] = attrValue
			} else {
				attrs[key] = nullValue(typ)
			}
		}
		objectValue, diags := types.ObjectValue(attrTypes, attrs)
		if diags.HasError() {
			return types.List{}, DiagnosticsToError(diags, "types.List", "")
		}
		values = append(values, objectValue)
	}
	listValue, diags := types.ListValue(objectType, values)
	if diags.HasError() {
		return types.List{}, DiagnosticsToError(diags, "types.List", "")
	}
	return listValue, nil
}

func objectValue(raw json.RawMessage) (types.Object, error) {
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		return types.Object{}, err
	}
	object = terraformNameMap(object)
	attrTypes := make(map[string]attr.Type, len(object))
	for key, value := range object {
		attrTypes[key] = inferredAttrType(value)
	}
	attrs := make(map[string]attr.Value, len(object))
	for key, value := range object {
		attrValue, err := attrValueFromAny(value, attrTypes[key])
		if err != nil {
			return types.Object{}, err
		}
		attrs[key] = attrValue
	}
	result, diags := types.ObjectValue(attrTypes, attrs)
	if diags.HasError() {
		return types.Object{}, DiagnosticsToError(diags, "types.Object", "")
	}
	return result, nil
}

func terraformNameMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[apiKeyToTerraformName(key)] = terraformNameValue(value)
	}
	return output
}

func terraformNameValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return terraformNameMap(typed)
	case []any:
		output := make([]any, len(typed))
		for index, item := range typed {
			output[index] = terraformNameValue(item)
		}
		return output
	default:
		return value
	}
}

func apiKeyToTerraformName(key string) string {
	prefix := ""
	if strings.HasPrefix(key, "__template_") {
		prefix = "__template_"
		key = strings.TrimPrefix(key, prefix)
	}
	var output strings.Builder
	var previous rune
	for index, char := range key {
		if char == '_' {
			output.WriteRune(char)
			previous = char
			continue
		}
		if unicode.IsUpper(char) {
			if index > 0 && previous != '_' {
				output.WriteByte('_')
			}
			output.WriteRune(unicode.ToLower(char))
		} else {
			output.WriteRune(char)
		}
		previous = char
	}
	return prefix + output.String()
}

func inferredAttrType(value any) attr.Type {
	switch typed := value.(type) {
	case bool:
		return types.BoolType
	case float64:
		return types.Float64Type
	case []any:
		if len(typed) > 0 {
			return types.ListType{ElemType: inferredAttrType(typed[0])}
		}
		return types.ListType{ElemType: types.StringType}
	case map[string]any:
		attrTypes := make(map[string]attr.Type, len(typed))
		for key, nested := range typed {
			attrTypes[key] = inferredAttrType(nested)
		}
		return types.ObjectType{AttrTypes: attrTypes}
	default:
		return types.StringType
	}
}

func attrValueFromAny(value any, typ attr.Type) (attr.Value, error) {
	if value == nil {
		return nullValue(typ), nil
	}
	switch typ := typ.(type) {
	case basetypes.BoolType:
		value, _ := value.(bool)
		return types.BoolValue(value), nil
	case basetypes.Float64Type:
		value, _ := value.(float64)
		return types.Float64Value(value), nil
	case basetypes.StringType:
		return types.StringValue(fmt.Sprintf("%v", value)), nil
	case types.ListType:
		values, _ := value.([]any)
		elements := make([]attr.Value, 0, len(values))
		for _, item := range values {
			element, err := attrValueFromAny(item, typ.ElementType())
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
		}
		result, diags := types.ListValue(typ.ElementType(), elements)
		if diags.HasError() {
			return nil, DiagnosticsToError(diags, "types.List", "")
		}
		return result, nil
	case types.ObjectType:
		object, _ := value.(map[string]any)
		attrs := make(map[string]attr.Value, len(typ.AttrTypes))
		for key, attrType := range typ.AttrTypes {
			attrValue, err := attrValueFromAny(object[key], attrType)
			if err != nil {
				return nil, err
			}
			attrs[key] = attrValue
		}
		result, diags := types.ObjectValue(typ.AttrTypes, attrs)
		if diags.HasError() {
			return nil, DiagnosticsToError(diags, "types.Object", "")
		}
		return result, nil
	default:
		return types.StringValue(fmt.Sprintf("%v", value)), nil
	}
}

func nullValue(typ attr.Type) attr.Value {
	switch typ := typ.(type) {
	case basetypes.BoolType:
		return types.BoolNull()
	case basetypes.Float64Type:
		return types.Float64Null()
	case basetypes.StringType:
		return types.StringNull()
	case types.ListType:
		return types.ListNull(typ.ElementType())
	case types.ObjectType:
		return types.ObjectNull(typ.AttrTypes)
	default:
		return types.StringNull()
	}
}

func stringValues(values []string) []attr.Value {
	out := make([]attr.Value, 0, len(values))
	for _, value := range values {
		out = append(out, types.StringValue(value))
	}
	return out
}

// modelFromSearchItemJSON builds a SearchDataset or SearchDatasetProvider model from a list-item
// JSON (e.g. when SDK cannot unmarshal cribl_lake). Sets ID, Description, Type, ProviderID from
// the item so HCL generation and import work without calling Get.
func modelFromSearchItemJSON(modelType reflect.Type, itemJSON []byte) (interface{}, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(itemJSON, &parsed); err != nil {
		return nil, err
	}
	getStr := func(m map[string]interface{}, keys ...string) string {
		for _, k := range keys {
			if v, ok := m[k]; ok && v != nil {
				if s, ok := v.(string); ok {
					return s
				}
			}
		}
		return ""
	}
	modelVal := reflect.New(modelType)
	val := modelVal.Elem()
	fieldValues := map[string]string{
		"ID":          getStr(parsed, "id"),
		"Description": getStr(parsed, "description"),
		"Type":        getStr(parsed, "type"),
		"ProviderID":  getStr(parsed, "provider_id", "provider"),
	}
	for fieldName, s := range fieldValues {
		f := val.FieldByName(fieldName)
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		if f.Type() != reflect.TypeOf(types.String{}) {
			continue
		}
		f.Set(reflect.ValueOf(types.StringValue(s)))
	}
	return modelVal.Interface(), nil
}

// bodyFromSearchRulesetGetResponse returns the Counted* payload for Local Search ruleset GET responses.
func bodyFromSearchRulesetGetResponse(resp interface{}) (interface{}, error) {
	switch r := resp.(type) {
	case *operations.GetDatasetRuleByIDResponse:
		if r == nil {
			return nil, nil
		}
		if r.CountedDatasetRuleset != nil {
			return r.CountedDatasetRuleset, nil
		}
	case *operations.GetDatatypeRuleByIDResponse:
		if r == nil {
			return nil, nil
		}
		if r.CountedDatatypeRuleset != nil {
			return r.CountedDatatypeRuleset, nil
		}
	}
	return nil, fmt.Errorf("search ruleset response has no Counted* body")
}

// callGetByID invokes the SDK Get* method for the entry and returns the response Object (ResponseBody).
func callGetByID(ctx context.Context, client *sdk.CriblIo, e registry.Entry, requestParams map[string]string) (interface{}, error) {
	clientVal := reflect.ValueOf(client)
	if clientVal.Kind() == reflect.Ptr {
		clientVal = clientVal.Elem()
	}
	svcField := clientVal.FieldByName(e.SDKService)
	if !svcField.IsValid() {
		return nil, fmt.Errorf("SDK service %q not found", e.SDKService)
	}
	if svcField.Kind() == reflect.Ptr && svcField.IsNil() {
		return nil, fmt.Errorf("SDK service %q is nil", e.SDKService)
	}
	// Search ruleset GETs are (ctx, opts ...operations.Option). MethodByName+reflect.Call/CallSlice is unreliable
	// across how the method value is obtained; call the SDK directly.
	if e.SDKService == "Search" {
		switch e.GetMethod {
		case "GetDatasetRuleByID":
			resp, err := client.Search.GetDatasetRuleByID(ctx)
			if err != nil {
				return nil, err
			}
			return bodyFromSearchRulesetGetResponse(resp)
		case "GetDatatypeRuleByID":
			resp, err := client.Search.GetDatatypeRuleByID(ctx)
			if err != nil {
				return nil, err
			}
			return bodyFromSearchRulesetGetResponse(resp)
		}
	}
	svc := svcField
	if svc.Kind() == reflect.Ptr {
		svc = svc.Elem()
	}
	method := svc.Addr().MethodByName(e.GetMethod)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %q not found on service %s", e.GetMethod, e.SDKService)
	}
	mt := method.Type()
	if mt.NumIn() < 2 || mt.NumOut() != 2 {
		return nil, fmt.Errorf("get method %s has unexpected signature", e.GetMethod)
	}
	// Variadic (ctx, opts ...operations.Option): must use CallSlice with a []T final arg, never Call with a slice.
	if mt.IsVariadic() && mt.NumIn() == 2 && mt.In(1).Kind() == reflect.Slice {
		optSlice := reflect.MakeSlice(mt.In(1), 0, 0)
		args := []reflect.Value{reflect.ValueOf(ctx), optSlice}
		outs := method.CallSlice(args)
		return unwrapGetResponse(outs)
	}
	paramType := mt.In(1)
	reqType := paramType
	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
	}
	reqVal := reflect.New(reqType)
	setRequestFields(reqVal, requestParams)
	// Method may take request by value (GetOutputByIDRequest) or by pointer (*GetXRequest).
	reqArg := reqVal
	if paramType.Kind() == reflect.Struct {
		reqArg = reqVal.Elem()
	}
	// Pass only fixed args (ctx, request). Variadic opts ...Option after the request are left empty.
	args := []reflect.Value{reflect.ValueOf(ctx), reqArg}
	outs := method.Call(args)
	return unwrapGetResponse(outs)
}

func unwrapGetResponse(outs []reflect.Value) (interface{}, error) {
	if len(outs) != 2 {
		return nil, fmt.Errorf("unexpected get method return count")
	}
	respVal := outs[0]
	errVal := outs[1]
	if !errVal.IsNil() {
		return nil, errVal.Interface().(error)
	}
	if respVal.IsNil() {
		return nil, nil
	}
	// response.Object; respVal may be *GetXResponse.
	if respVal.Kind() == reflect.Ptr {
		respVal = respVal.Elem()
	}
	objectField := respVal.FieldByName("Object")
	if objectField.IsValid() && !(objectField.Kind() == reflect.Ptr && objectField.IsNil()) {
		return objectField.Interface(), nil
	}
	// Search (and similar) SDK responses use Counted* payloads instead of Object.
	for _, name := range []string{"CountedLocalSearchEngine", "CountedLocalSearchSource", "CountedDatasetRuleset", "CountedDatatypeRuleset"} {
		f := respVal.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.Ptr && !f.IsNil() {
			return f.Interface(), nil
		}
	}
	return nil, fmt.Errorf("response has no Object or supported Counted* body field")
}

func setRequestFields(reqVal reflect.Value, params map[string]string) {
	if reqVal.Kind() == reflect.Ptr {
		reqVal = reqVal.Elem()
	}
	if !reqVal.IsValid() || reqVal.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < reqVal.NumField(); i++ {
		f := reqVal.Field(i)
		if !f.CanSet() {
			continue
		}
		name := reqVal.Type().Field(i).Name
		if v, ok := params[name]; ok && v != "" {
			switch f.Kind() {
			case reflect.String:
				f.SetString(v)
			}
		}
	}
}
