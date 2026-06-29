// Package converter converts REST API responses into provider ResourceModel
// instances using reflection and generated model metadata.
package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"unicode"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// RefreshFromMethodName returns the RefreshFrom* method name for a legacy Get method name.
// Convention: RefreshFromOperations + GetMethod + "ResponseBody" (e.g. GetInputByID -> RefreshFromOperationsGetInputByIDResponseBody).
func RefreshFromMethodName(getMethod string) string {
	if getMethod == "" {
		return ""
	}
	return "RefreshFromOperations" + getMethod + "ResponseBody"
}

// Convert fetches a single resource via REST and converts the response into a
// provider ResourceModel.
func Convert(ctx context.Context, client *importclient.Client, e registry.Entry, requestParams map[string]string) (model interface{}, err error) {
	if e.RESTGetPath == "" {
		return nil, fmt.Errorf("%s: no RESTGetPath in registry", e.TypeName)
	}
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	modelTypes := ResourceModelTypes()
	modelType, ok := modelTypes[e.ModelTypeName]
	if !ok {
		return nil, fmt.Errorf("%s: unknown model type %q", e.TypeName, e.ModelTypeName)
	}

	raw, err := callRESTGetByID(ctx, client, e, requestParams)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", e.TypeName, err)
	}
	converted, convErr := convertGeneratedModelFromRawItem(e, modelType, raw)
	if convErr != nil {
		return nil, convErr
	}
	if err = InjectRequiredIdentifiers(converted, requestParams); err != nil {
		return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, err)
	}
	return converted, nil
}

// ConvertRawItemWithIdentifiers converts a raw API item into a provider resource
// model and injects Terraform identifier fields.
func ConvertRawItemWithIdentifiers(e registry.Entry, itemJSON json.RawMessage, identifiers map[string]string) (interface{}, error) {
	modelTypes := ResourceModelTypes()
	modelType, ok := modelTypes[e.ModelTypeName]
	if !ok {
		return nil, fmt.Errorf("%s: unknown model type %q", e.TypeName, e.ModelTypeName)
	}
	model, err := convertGeneratedModelFromRawItem(e, modelType, itemJSON)
	if err != nil {
		return nil, err
	}
	if err := InjectRequiredIdentifiers(model, identifiers); err != nil {
		return nil, fmt.Errorf("%s: inject identifiers: %w", e.TypeName, err)
	}
	return model, nil
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
		if _, ok := GeneratedModelTypes()[e.ModelTypeName]; ok {
			return convertGeneratedModelFromResponseBody(e, modelType, responseBody)
		}
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
	return convertGeneratedModelFromRawItem(e, modelType, itemJSON)
}

func convertGeneratedModelFromRawItem(e registry.Entry, modelType reflect.Type, itemJSON json.RawMessage) (interface{}, error) {
	var values map[string]json.RawMessage
	if err := json.Unmarshal(itemJSON, &values); err != nil {
		return nil, fmt.Errorf("%s: decode response item: %w", e.TypeName, err)
	}

	modelVal := reflect.New(modelType)
	if generatedType, ok := GeneratedModelTypes()[e.ModelTypeName]; ok {
		generatedVal := reflect.New(generatedType)
		if err := json.Unmarshal(itemJSON, generatedVal.Interface()); err == nil {
			copyMatchingFields(modelVal.Elem(), generatedVal.Elem())
			return modelVal.Interface(), nil
		}
	}
	if err := populateGeneratedModel(modelVal.Elem(), values); err != nil {
		return nil, fmt.Errorf("%s: populate generated model: %w", e.TypeName, err)
	}
	return modelVal.Interface(), nil
}

func callRESTGetByID(ctx context.Context, client *importclient.Client, e registry.Entry, requestParams map[string]string) (json.RawMessage, error) {
	if client == nil || client.REST == nil {
		return nil, fmt.Errorf("REST client is nil")
	}
	values := restPathValues(requestParams)
	path := renderRESTPath(e.RESTGetPath, values)
	items, err := restclient.Get[[]json.RawMessage](ctx, client.REST, path)
	if err == nil {
		if items == nil || len(*items) == 0 {
			return nil, fmt.Errorf("%s: empty REST response", e.TypeName)
		}
		if id := requestParams["ID"]; id != "" {
			for _, item := range *items {
				var object map[string]any
				if json.Unmarshal(item, &object) == nil && rawString(object, "id", "ID", "Id", "keyID", "keyId", "name") == id {
					return item, nil
				}
			}
		}
		return (*items)[0], nil
	}
	item, itemErr := restclient.Get[map[string]json.RawMessage](ctx, client.REST, path)
	if itemErr != nil {
		return nil, err
	}
	return json.Marshal(item)
}

func restPathValues(requestParams map[string]string) map[string]string {
	return map[string]string{
		"group_id": requestParams["GroupID"],
		"id":       requestParams["ID"],
		"pack":     requestParams["Pack"],
		"lake_id":  requestParams["LakeID"],
	}
}

func renderRESTPath(path string, values map[string]string) string {
	rendered := path
	for key, value := range values {
		rendered = strings.ReplaceAll(rendered, "{"+key+"}", url.PathEscape(value))
	}
	return rendered
}

func rawString(item map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := item[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case string:
			return strings.TrimSpace(typed)
		case fmt.Stringer:
			return strings.TrimSpace(typed.String())
		default:
			if s := fmt.Sprint(typed); s != "" && s != "<nil>" {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func copyMatchingFields(dst, src reflect.Value) {
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		if !dstField.CanSet() {
			continue
		}
		srcField := src.FieldByName(dst.Type().Field(i).Name)
		if !srcField.IsValid() {
			continue
		}
		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
			continue
		}
		if srcField.Type().ConvertibleTo(dstField.Type()) {
			dstField.Set(srcField.Convert(dstField.Type()))
		}
	}
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
