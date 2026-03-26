// Mirrors items[0] one-of into top-level input_* fields so state matches config after
// import/refresh (same idea as pack_resource_sdk RefreshFrom). The GET response only
// carries items; root input_* were left unset, so Terraform planned perpetual adds.

package provider

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"

	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
)

// priorPlainAuthTokens holds copies of []types.String auth_tokens from state before a GET refresh.
// The API often omits or redacts tokens; sync copies a shorter list onto root and Terraform then
// plans spurious adds. We restore the prior slice when the API returned fewer elements.
type priorPlainAuthTokens struct {
	http          []types.String
	httpRaw       []types.String
	firehose      []types.String
	elastic       []types.String
	criblLakeHTTP []types.String
	wizWebhook    []types.String
}

func snapshotPlainAuthTokensPriorRead(r *SourceResourceModel) priorPlainAuthTokens {
	var p priorPlainAuthTokens
	if r.InputHTTP != nil {
		p.http = cloneTypesStringSlice(r.InputHTTP.AuthTokens)
	}
	if r.InputHTTPRaw != nil {
		p.httpRaw = cloneTypesStringSlice(r.InputHTTPRaw.AuthTokens)
	}
	if r.InputFirehose != nil {
		p.firehose = cloneTypesStringSlice(r.InputFirehose.AuthTokens)
	}
	if r.InputElastic != nil {
		p.elastic = cloneTypesStringSlice(r.InputElastic.AuthTokens)
	}
	if r.InputCriblLakeHTTP != nil {
		p.criblLakeHTTP = cloneTypesStringSlice(r.InputCriblLakeHTTP.AuthTokens)
	}
	if r.InputWizWebhook != nil {
		p.wizWebhook = cloneTypesStringSlice(r.InputWizWebhook.AuthTokens)
	}
	return p
}

func restorePlainAuthTokensIfAPIShrank(r *SourceResourceModel, prior priorPlainAuthTokens) {
	restore := func(dst **tfTypes.InputHTTP, prev []types.String) {
		if *dst == nil || prev == nil {
			return
		}
		if len(prev) > len((*dst).AuthTokens) {
			(*dst).AuthTokens = cloneTypesStringSlice(prev)
		}
	}
	restoreHTTPRaw := func(dst **tfTypes.InputHTTPRaw, prev []types.String) {
		if *dst == nil || prev == nil {
			return
		}
		if len(prev) > len((*dst).AuthTokens) {
			(*dst).AuthTokens = cloneTypesStringSlice(prev)
		}
	}
	restoreFirehose := func(dst **tfTypes.InputFirehose, prev []types.String) {
		if *dst == nil || prev == nil {
			return
		}
		if len(prev) > len((*dst).AuthTokens) {
			(*dst).AuthTokens = cloneTypesStringSlice(prev)
		}
	}
	restoreElastic := func(dst **tfTypes.InputElastic, prev []types.String) {
		if *dst == nil || prev == nil {
			return
		}
		if len(prev) > len((*dst).AuthTokens) {
			(*dst).AuthTokens = cloneTypesStringSlice(prev)
		}
	}
	restoreCriblLake := func(dst **tfTypes.InputCriblLakeHTTP, prev []types.String) {
		if *dst == nil || prev == nil {
			return
		}
		if len(prev) > len((*dst).AuthTokens) {
			(*dst).AuthTokens = cloneTypesStringSlice(prev)
		}
	}
	restoreWiz := func(dst **tfTypes.InputWizWebhook, prev []types.String) {
		if *dst == nil || prev == nil {
			return
		}
		if len(prev) > len((*dst).AuthTokens) {
			(*dst).AuthTokens = cloneTypesStringSlice(prev)
		}
	}

	restore(&r.InputHTTP, prior.http)
	restoreHTTPRaw(&r.InputHTTPRaw, prior.httpRaw)
	restoreFirehose(&r.InputFirehose, prior.firehose)
	restoreElastic(&r.InputElastic, prior.elastic)
	restoreCriblLake(&r.InputCriblLakeHTTP, prior.criblLakeHTTP)
	restoreWiz(&r.InputWizWebhook, prior.wizWebhook)
}

func cloneTypesStringSlice(s []types.String) []types.String {
	if s == nil {
		return nil
	}
	out := make([]types.String, len(s))
	copy(out, s)
	return out
}

func (r *SourceResourceModel) syncRootInputFromFirstItem() {
	if len(r.Items) == 0 {
		return
	}
	u := r.Items[0]
	uv := reflect.ValueOf(u)
	ut := uv.Type()
	rv := reflect.ValueOf(r).Elem()

	for i := 0; i < ut.NumField(); i++ {
		name := ut.Field(i).Name
		if name == "InputOpenai" {
			continue
		}
		src := uv.Field(i)
		dst := rv.FieldByName(name)
		if !dst.IsValid() || !dst.CanSet() {
			continue
		}
		if src.Type().AssignableTo(dst.Type()) {
			dst.Set(src)
		}
	}

	if u.InputOpenai != nil {
		r.InputOpenai = openai1ToOpenai(u.InputOpenai)
	} else {
		r.InputOpenai = nil
	}

	// Root input_* blocks use Optional list attributes; empty API slices encode as [] on the
	// shared struct while Terraform treats omitted optional lists as null. Nil out empty
	// slices on the active root input so plan matches typical HCL (same struct backs items[0]).
	normalizeRootInputEmptySlices(r)
}

// normalizeRootInputEmptySlices clears len-0 slices on each non-nil root input_* pointer.
func normalizeRootInputEmptySlices(r *SourceResourceModel) {
	rv := reflect.ValueOf(r).Elem()
	ut := reflect.TypeOf(tfTypes.InputUnion1{})
	for i := 0; i < ut.NumField(); i++ {
		name := ut.Field(i).Name
		dst := rv.FieldByName(name)
		if !dst.IsValid() || dst.IsNil() || dst.Kind() != reflect.Pointer {
			continue
		}
		normalizeEmptySlicesToNil(dst)
	}
}

// normalizeEmptySlicesToNil sets empty Go slices to nil so optional TF list attrs serialize as null.
// Recurses into non-empty slices of structs (e.g. input_splunk_hec.auth_tokens[*].allowed_indexes_at_token).
func normalizeEmptySlicesToNil(v reflect.Value) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Slice:
			if f.Len() == 0 {
				f.Set(reflect.Zero(f.Type()))
				continue
			}
			elemT := f.Type().Elem()
			switch elemT.Kind() {
			case reflect.Struct:
				for j := 0; j < f.Len(); j++ {
					el := f.Index(j)
					if el.CanAddr() {
						normalizeEmptySlicesToNil(el.Addr())
					}
				}
			case reflect.Pointer:
				if elemT.Elem().Kind() == reflect.Struct {
					for j := 0; j < f.Len(); j++ {
						el := f.Index(j)
						if !el.IsNil() {
							normalizeEmptySlicesToNil(el)
						}
					}
				}
			}
		case reflect.Pointer:
			if !f.IsNil() {
				normalizeEmptySlicesToNil(f)
			}
		case reflect.Struct:
			if f.CanAddr() {
				normalizeEmptySlicesToNil(f.Addr())
			}
		}
	}
}

func openai1ToOpenai(p *tfTypes.InputOpenai1) *tfTypes.InputOpenai {
	if p == nil {
		return nil
	}
	cc := make([]tfTypes.ContentConfig, 0, len(p.ContentConfig))
	for _, row := range p.ContentConfig {
		cc = append(cc, inputOpenaiContentConfigToContentConfig(row))
	}
	return &tfTypes.InputOpenai{
		TemplateOpenaiOrganization: p.TemplateOpenaiOrganization,
		TemplateOpenaiProject:      p.TemplateOpenaiProject,
		APIKey:                     p.APIKey,
		Connections:                p.Connections,
		ContentConfig:              cc,
		Description:                p.Description,
		Disabled:                   p.Disabled,
		Environment:                p.Environment,
		ID:                         p.ID,
		IgnoreGroupJobsLimit:       p.IgnoreGroupJobsLimit,
		KeepAliveTime:              p.KeepAliveTime,
		MaxMissedKeepAlives:        p.MaxMissedKeepAlives,
		Metadata:                   p.Metadata,
		OpenaiOrganization:         p.OpenaiOrganization,
		OpenaiProject:              p.OpenaiProject,
		Pipeline:                   p.Pipeline,
		Pq:                         p.Pq,
		PqEnabled:                  p.PqEnabled,
		RequestTimeout:             p.RequestTimeout,
		RetryRules:                 p.RetryRules,
		SendToRoutes:               p.SendToRoutes,
		Streamtags:                 p.Streamtags,
		TextSecret:                 p.TextSecret,
		TTL:                        p.TTL,
		Type:                       p.Type,
	}
}

func inputOpenaiContentConfigToContentConfig(in tfTypes.InputOpenaiContentConfig) tfTypes.ContentConfig {
	return tfTypes.ContentConfig{
		CronSchedule:                    in.CronSchedule,
		Disabled:                        in.Disabled,
		Earliest:                        in.Earliest,
		EndpointMetadata:                in.EndpointMetadata,
		JobTimeout:                      in.JobTimeout,
		Latest:                          in.Latest,
		LogLevel:                        in.LogLevel,
		ManageState:                     in.ManageState,
		MaxPages:                        in.MaxPages,
		PaginationAttribute:             in.PaginationAttribute,
		PaginationCurRelationAttribute:  in.PaginationCurRelationAttribute,
		PaginationLastPageExpr:          in.PaginationLastPageExpr,
		PaginationNextRelationAttribute: in.PaginationNextRelationAttribute,
		PaginationType:                  in.PaginationType,
		RequestParams:                   in.RequestParams,
		StateMergeExpression:            in.StateMergeExpression,
		StateTracking:                   in.StateTracking,
		StateUpdateExpression:           in.StateUpdateExpression,
	}
}
