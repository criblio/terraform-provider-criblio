package provider

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// priorSourcePlainAuthTokens holds copies of plain string auth_tokens before a GET refresh.
// Some source APIs omit or truncate these tokens on read; restore keeps prior state when
// the API returns fewer tokens than Terraform already knows about.
type priorSourcePlainAuthTokens struct {
	http          []types.String
	httpRaw       []types.String
	firehose      []types.String
	elastic       []types.String
	criblLakeHTTP []types.String
	wizWebhook    []types.String
}

func snapshotSourcePlainAuthTokensPriorRead(r *SourceModel) priorSourcePlainAuthTokens {
	var prior priorSourcePlainAuthTokens
	if r.InputHttp != nil {
		prior.http = cloneSourceStringValues(sourceStringListValues(r.InputHttp.AuthTokens))
	}
	if r.InputHttpRaw != nil {
		prior.httpRaw = cloneSourceStringValues(sourceStringListValues(r.InputHttpRaw.AuthTokens))
	}
	if r.InputFirehose != nil {
		prior.firehose = cloneSourceStringValues(sourceStringListValues(r.InputFirehose.AuthTokens))
	}
	if r.InputElastic != nil {
		prior.elastic = cloneSourceStringValues(sourceStringListValues(r.InputElastic.AuthTokens))
	}
	if r.InputCriblLakeHttp != nil {
		prior.criblLakeHTTP = cloneSourceStringValues(sourceStringListValues(r.InputCriblLakeHttp.AuthTokens))
	}
	if r.InputWizWebhook != nil {
		prior.wizWebhook = cloneSourceStringValues(sourceStringListValues(r.InputWizWebhook.AuthTokens))
	}
	return prior
}

func restoreSourcePlainAuthTokensIfAPIShrank(r *SourceModel, prior priorSourcePlainAuthTokens) {
	restore := func(dst *types.List, prev []types.String) {
		if dst == nil || prev == nil {
			return
		}
		// Preserve legacy behavior: restore only when len(prior) > len(api_returned).
		if len(prev) > len(sourceStringListValues(*dst)) {
			*dst = sourceStringListFromValues(prev)
		}
	}

	if r.InputHttp != nil {
		restore(&r.InputHttp.AuthTokens, prior.http)
	}
	if r.InputHttpRaw != nil {
		restore(&r.InputHttpRaw.AuthTokens, prior.httpRaw)
	}
	if r.InputFirehose != nil {
		restore(&r.InputFirehose.AuthTokens, prior.firehose)
	}
	if r.InputElastic != nil {
		restore(&r.InputElastic.AuthTokens, prior.elastic)
	}
	if r.InputCriblLakeHttp != nil {
		restore(&r.InputCriblLakeHttp.AuthTokens, prior.criblLakeHTTP)
	}
	if r.InputWizWebhook != nil {
		restore(&r.InputWizWebhook.AuthTokens, prior.wizWebhook)
	}
}

func snapshotPackSourcePlainAuthTokensPriorRead(r *PackSourceModel) priorSourcePlainAuthTokens {
	var prior priorSourcePlainAuthTokens
	if r.InputHttp != nil {
		prior.http = cloneSourceStringValues(sourceStringListValues(r.InputHttp.AuthTokens))
	}
	if r.InputHttpRaw != nil {
		prior.httpRaw = cloneSourceStringValues(sourceStringListValues(r.InputHttpRaw.AuthTokens))
	}
	if r.InputFirehose != nil {
		prior.firehose = cloneSourceStringValues(sourceStringListValues(r.InputFirehose.AuthTokens))
	}
	if r.InputElastic != nil {
		prior.elastic = cloneSourceStringValues(sourceStringListValues(r.InputElastic.AuthTokens))
	}
	if r.InputCriblLakeHttp != nil {
		prior.criblLakeHTTP = cloneSourceStringValues(sourceStringListValues(r.InputCriblLakeHttp.AuthTokens))
	}
	if r.InputWizWebhook != nil {
		prior.wizWebhook = cloneSourceStringValues(sourceStringListValues(r.InputWizWebhook.AuthTokens))
	}
	return prior
}

func restorePackSourcePlainAuthTokensIfAPIShrank(r *PackSourceModel, prior priorSourcePlainAuthTokens) {
	restore := func(dst *types.List, prev []types.String) {
		if dst == nil || prev == nil {
			return
		}
		// Preserve legacy behavior: restore only when len(prior) > len(api_returned).
		if len(prev) > len(sourceStringListValues(*dst)) {
			*dst = sourceStringListFromValues(prev)
		}
	}

	if r.InputHttp != nil {
		restore(&r.InputHttp.AuthTokens, prior.http)
	}
	if r.InputHttpRaw != nil {
		restore(&r.InputHttpRaw.AuthTokens, prior.httpRaw)
	}
	if r.InputFirehose != nil {
		restore(&r.InputFirehose.AuthTokens, prior.firehose)
	}
	if r.InputElastic != nil {
		restore(&r.InputElastic.AuthTokens, prior.elastic)
	}
	if r.InputCriblLakeHttp != nil {
		restore(&r.InputCriblLakeHttp.AuthTokens, prior.criblLakeHTTP)
	}
	if r.InputWizWebhook != nil {
		restore(&r.InputWizWebhook.AuthTokens, prior.wizWebhook)
	}
}

func sourceStringListValues(list types.List) []types.String {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	elements := list.Elements()
	values := make([]types.String, 0, len(elements))
	for _, element := range elements {
		value, ok := element.(types.String)
		if !ok {
			return nil
		}
		values = append(values, value)
	}
	return values
}

func sourceStringListFromValues(values []types.String) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, value)
	}
	return types.ListValueMust(types.StringType, elements)
}

func cloneSourceStringValues(values []types.String) []types.String {
	if values == nil {
		return nil
	}
	out := make([]types.String, len(values))
	copy(out, values)
	return out
}

func sourceRequestModelWithHoistedIdentity(model SourceModel) SourceModel {
	return sourceLikeRequestModelWithHoistedIdentity(model)
}

func packSourceRequestModelWithHoistedIdentity(model PackSourceModel) PackSourceModel {
	return sourceLikeRequestModelWithHoistedIdentity(model)
}

func sourceLikeRequestModelWithHoistedIdentity[T any](model T) T {
	request := model
	rv := reflect.ValueOf(&request).Elem()
	idField := rv.FieldByName("ID")
	if !idField.IsValid() || idField.Type() != reflect.TypeOf(types.String{}) {
		return request
	}
	modelID := idField.Interface().(types.String)
	if modelID.IsNull() || modelID.IsUnknown() {
		return request
	}
	for i := 0; i < rv.NumField(); i++ {
		fieldInfo := rv.Type().Field(i)
		if !isSourceInputField(fieldInfo.Name) {
			continue
		}
		field := rv.Field(i)
		if field.Kind() != reflect.Pointer || field.IsNil() {
			continue
		}
		clone := reflect.New(field.Type().Elem())
		clone.Elem().Set(field.Elem())
		field.Set(clone)
		input := field.Elem()
		idField := input.FieldByName("ID")
		if !idField.IsValid() || !idField.CanSet() || idField.Type() != reflect.TypeOf(types.String{}) {
			return request
		}
		id := idField.Interface().(types.String)
		if id.IsNull() || id.IsUnknown() {
			idField.Set(reflect.ValueOf(modelID))
		}
		return request
	}
	return request
}

// normalizeSourceRootInputEmptyLists converts empty API list values on imported
// active source blocks back to null so omitted optional list fields do not drift.
// Normal create/read/update paths preserve configured [] values for compatibility.
func normalizeSourceRootInputEmptyLists(r *SourceModel) {
	normalizeSourceLikeRootInputEmptyLists(r)
}

func normalizePackSourceRootInputEmptyLists(r *PackSourceModel) {
	normalizeSourceLikeRootInputEmptyLists(r)
}

func normalizeSourceLikeRootInputEmptyLists[T any](r *T) {
	if r == nil {
		return
	}
	rv := reflect.ValueOf(r).Elem()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.Kind() != reflect.Pointer || field.IsNil() {
			continue
		}
		if !rv.Type().Field(i).Anonymous && !isSourceInputField(rv.Type().Field(i).Name) {
			continue
		}
		normalizeSourceEmptyListsToNull(field)
	}
}

func isSourceInputField(name string) bool {
	return len(name) > len("Input") && name[:len("Input")] == "Input"
}

func normalizeSourceEmptyListsToNull(v reflect.Value) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	listType := reflect.TypeOf(types.List{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		if field.Type() == listType {
			list := field.Interface().(types.List)
			if !list.IsNull() && !list.IsUnknown() && len(list.Elements()) == 0 {
				field.Set(reflect.ValueOf(types.ListNull(list.ElementType(context.Background()))))
			}
			continue
		}
		switch field.Kind() {
		case reflect.Pointer:
			normalizeSourceEmptyListsToNull(field)
		case reflect.Struct:
			if field.CanAddr() {
				normalizeSourceEmptyListsToNull(field.Addr())
			}
		}
	}
}
