package provider

import (
	"context"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// reservedDatasetRuleIDs are rule ids that must stay on the server when destroying
// the resource (product catch-all / system rules). Terraform state may still list them.
var reservedDatasetRuleIDs = map[string]struct{}{
	"default": {},
}

// deleteSearchDatasetRuleset removes rules whose ids appear in Terraform state, then PATCHes
// the rest. Rules in reservedDatasetRuleIDs are never removed. Other product-managed rules
// not listed in state are unchanged.
func deleteSearchDatasetRuleset(ctx context.Context, r *SearchDatasetRulesetResource, data *SearchDatasetRulesetResourceModel, resp *resource.DeleteResponse) {
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

	wantRulesetID := shared.DatasetRulesetID(data.ID.ValueString())
	getRes, err := r.client.Search.GetDatasetRuleByID(ctx)
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
	if getRes.CountedDatasetRuleset == nil {
		resp.Diagnostics.AddError("unexpected response from API", "missing counted dataset ruleset body")
		return
	}

	var current *shared.DatasetRuleset
	for i := range getRes.CountedDatasetRuleset.Items {
		item := &getRes.CountedDatasetRuleset.Items[i]
		if item.ID == wantRulesetID {
			current = item
			break
		}
	}
	if current == nil && len(getRes.CountedDatasetRuleset.Items) == 1 {
		current = &getRes.CountedDatasetRuleset.Items[0]
	}
	if current == nil {
		resp.Diagnostics.AddError(
			"unexpected response from API",
			fmt.Sprintf("could not find dataset ruleset with id %q in GET response", wantRulesetID),
		)
		return
	}

	kept := make([]shared.DatasetRule, 0, len(current.Rules))
	for _, rule := range current.Rules {
		if _, inManaged := managed[rule.ID]; !inManaged {
			kept = append(kept, rule)
			continue
		}
		if _, reserved := reservedDatasetRuleIDs[rule.ID]; reserved {
			kept = append(kept, rule)
			continue
		}
	}

	patch := shared.DatasetRuleset{
		ID:    wantRulesetID,
		Rules: kept,
	}
	res, err := r.client.Search.UpdateDatasetRuleByID(ctx, patch)
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
