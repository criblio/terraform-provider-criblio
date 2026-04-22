package provider

import (
	"context"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// deleteSearchDatatypeRuleset removes only rules whose ids appear in Terraform state,
// then PATCHes the rest. Product-managed rules (e.g. default auto-datatype) stay on the server.
func deleteSearchDatatypeRuleset(ctx context.Context, r *SearchDatatypeRulesetResource, data *SearchDatatypeRulesetResourceModel, resp *resource.DeleteResponse) {
	managed := make(map[string]struct{})
	for _, rule := range data.Rules {
		id := rule.ID.ValueString()
		if id != "" {
			managed[id] = struct{}{}
		}
	}
	if len(managed) == 0 {
		return
	}

	wantRulesetID := shared.DatatypeRulesetID(data.ID.ValueString())
	getRes, err := r.client.Search.GetDatatypeRuleByID(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		if getRes != nil && getRes.RawResponse != nil {
			resp.Diagnostics.AddError("unexpected http request/response", debugResponse(getRes.RawResponse))
		}
		return
	}
	if getRes == nil {
		resp.Diagnostics.AddError("unexpected response from API", fmt.Sprintf("%v", getRes))
		return
	}
	if getRes.StatusCode == 404 {
		return
	}
	if getRes.StatusCode != 200 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", getRes.StatusCode),
			debugResponse(getRes.RawResponse),
		)
		return
	}
	if getRes.CountedDatatypeRuleset == nil {
		resp.Diagnostics.AddError("unexpected response from API", "missing counted datatype ruleset body")
		return
	}

	var current *shared.DatatypeRuleset
	for i := range getRes.CountedDatatypeRuleset.Items {
		item := &getRes.CountedDatatypeRuleset.Items[i]
		if item.ID == wantRulesetID {
			current = item
			break
		}
	}
	if current == nil && len(getRes.CountedDatatypeRuleset.Items) == 1 {
		current = &getRes.CountedDatatypeRuleset.Items[0]
	}
	if current == nil {
		resp.Diagnostics.AddError(
			"unexpected response from API",
			fmt.Sprintf("could not find datatype ruleset with id %q in GET response", wantRulesetID),
		)
		return
	}

	kept := make([]shared.DatatypeRule, 0, len(current.Rules))
	for _, rule := range current.Rules {
		if _, remove := managed[rule.ID]; remove {
			continue
		}
		kept = append(kept, rule)
	}

	patch := shared.DatatypeRuleset{
		ID:    wantRulesetID,
		Rules: kept,
	}
	res, err := r.client.Search.UpdateDatatypeRuleByID(ctx, patch)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		if res != nil && res.RawResponse != nil {
			resp.Diagnostics.AddError("unexpected http request/response", debugResponse(res.RawResponse))
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("unexpected response from API", fmt.Sprintf("%v", res))
		return
	}
	switch res.StatusCode {
	case 200, 404:
		return
	default:
		resp.Diagnostics.AddError(
			fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", res.StatusCode),
			debugResponse(res.RawResponse),
		)
	}
}
