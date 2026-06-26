package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// reservedDatasetRuleIDs are rule ids that must stay on the server when destroying
// the resource (product catch-all / system rules). Terraform state may still list them.
var reservedDatasetRuleIDs = map[string]struct{}{
	"default": {},
}

// deleteSearchDatasetRuleset removes rules whose ids appear in Terraform state, then PATCHes
// the rest. Rules in reservedDatasetRuleIDs are never removed. Other product-managed rules
// not listed in state are unchanged.
func deleteSearchDatasetRuleset(ctx context.Context, client *restclient.Client, model SearchDatasetRulesetModel) error {
	managed, err := searchDatasetRuleIDs(model.Rules)
	if err != nil {
		return err
	}
	if len(managed) == 0 {
		return nil
	}

	rulesetID := searchDatasetRulesetID(model)
	path := fmt.Sprintf("/m/default_search/search/local_search/dataset-rulesets/%s", url.PathEscape(rulesetID))
	current, err := restclient.Get[SearchDatasetRulesetModel](ctx, client, path)
	if err != nil {
		if restclient.IsNotFound(err) {
			return nil
		}
		return err
	}
	if current == nil {
		return nil
	}

	currentRules, err := searchDatasetRules(current.Rules)
	if err != nil {
		return err
	}
	kept := make([]any, 0, len(currentRules))
	for _, rule := range currentRules {
		id, _ := rule["id"].(string)
		if _, inManaged := managed[id]; !inManaged {
			kept = append(kept, rule)
			continue
		}
		if _, reserved := reservedDatasetRuleIDs[id]; reserved {
			kept = append(kept, rule)
		}
	}

	patch := map[string]any{
		"id":    rulesetID,
		"rules": kept,
	}
	if err := restclient.PatchNoResponse(ctx, client, path, patch); err != nil && !restclient.IsNotFound(err) {
		return err
	}
	return nil
}

func searchDatasetRulesetID(model SearchDatasetRulesetModel) string {
	if !model.ID.IsNull() && !model.ID.IsUnknown() && model.ID.ValueString() != "" {
		return model.ID.ValueString()
	}
	return "default"
}

func searchDatasetRuleIDs(value attr.Value) (map[string]struct{}, error) {
	rules, err := searchDatasetRules(value)
	if err != nil {
		return nil, err
	}
	managed := make(map[string]struct{}, len(rules))
	for _, rule := range rules {
		id, _ := rule["id"].(string)
		if id != "" {
			managed[id] = struct{}{}
		}
	}
	return managed, nil
}

func searchDatasetRules(value attr.Value) ([]map[string]any, error) {
	raw, err := SearchDatasetRulesetTerraformValueToJSON(value)
	if err != nil {
		return nil, fmt.Errorf("convert dataset rules to API value: %v", err)
	}
	items, _ := raw.([]any)
	rules := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rule, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
