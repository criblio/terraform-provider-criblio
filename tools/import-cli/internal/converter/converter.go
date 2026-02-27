// Package converter converts SDK API responses into provider ResourceModel instances
// using reflection and provider RefreshFrom* methods.
package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// fillPackPipelineConf fills the pack pipeline model's Conf. GetPipelinesByPackWithID returns only
// Routes (items), not the pipeline definition. We try Pipelines.GetPipelineByID (lib) for the same
// group/id; if that returns one pipeline we use it. Otherwise we set minimal conf so HCL is valid.
func fillPackPipelineConf(ctx context.Context, client *sdk.CriblIo, converted interface{}, requestParams map[string]string) {
	pm, ok := converted.(*provider.PackPipelineResourceModel)
	if !ok || client == nil || client.Pipelines == nil {
		return
	}
	groupID := requestParams["GroupID"]
	id := requestParams["ID"]
	if groupID == "" {
		groupID = "default"
	}
	if id == "" {
		ensureMinimalPackPipelineConf(pm)
		return
	}
	resp, err := client.Pipelines.GetPipelineByID(ctx, operations.GetPipelineByIDRequest{GroupID: groupID, ID: id})
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
	// Pass only fixed args (ctx, request). Variadic opts ...Option are left empty.
	args := []reflect.Value{reflect.ValueOf(ctx), reqArg}
	outs := method.Call(args)
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
	if !objectField.IsValid() {
		return nil, fmt.Errorf("response has no Object field")
	}
	if objectField.Kind() == reflect.Ptr && objectField.IsNil() {
		return nil, nil
	}
	return objectField.Interface(), nil
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
