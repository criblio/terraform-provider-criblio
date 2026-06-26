package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// deleteSearchDatatypeRuleset removes only rules whose ids appear in Terraform state,
// then PATCHes the rest. Product-managed rules (e.g. default auto-datatype) stay on the server.
func deleteSearchDatatypeRuleset(ctx context.Context, client *restclient.Client, model SearchDatatypeRulesetModel) error {
	managed, err := searchDatatypeRuleIDs(model.Rules)
	if err != nil {
		return err
	}
	if len(managed) == 0 {
		return nil
	}

	rulesetID := searchDatatypeRulesetID(model)
	path := fmt.Sprintf("/m/default_search/search/local_search/datatype-rulesets/%s", url.PathEscape(rulesetID))
	current, err := restclient.Get[SearchDatatypeRulesetModel](ctx, client, path)
	if err != nil {
		if restclient.IsNotFound(err) {
			return nil
		}
		return err
	}
	if current == nil {
		return nil
	}

	currentRules, err := searchDatatypeRules(current.Rules)
	if err != nil {
		return err
	}
	kept := make([]any, 0, len(currentRules))
	for _, rule := range currentRules {
		id, _ := rule["id"].(string)
		if _, remove := managed[id]; remove {
			continue
		}
		kept = append(kept, rule)
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

func searchDatatypeRulesetID(model SearchDatatypeRulesetModel) string {
	if !model.ID.IsNull() && !model.ID.IsUnknown() && model.ID.ValueString() != "" {
		return model.ID.ValueString()
	}
	return "default"
}

func searchDatatypeRuleIDs(value attr.Value) (map[string]struct{}, error) {
	rules, err := searchDatatypeRules(value)
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

func searchDatatypeRules(value attr.Value) ([]map[string]any, error) {
	raw, err := SearchDatatypeRulesetTerraformValueToJSON(value)
	if err != nil {
		return nil, fmt.Errorf("convert datatype rules to API value: %v", err)
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
