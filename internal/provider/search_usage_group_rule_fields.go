package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// searchUsageGroupRulesWithKnownAPIValues fills the API-owned attributes of each
// limit rule from the API response.
//
// Every limit type carries a fixed `description`, and some carry a fixed
// `metric`, which the API assigns when absent and overwrites when sent. Both are
// therefore computed, which leaves them unknown after an apply whose
// configuration set only `limit` and `type`. Reads preserve configured inputs,
// so `rules` comes from configuration and nothing would otherwise fill those
// nested values, and unknown-after-apply is always rejected by Terraform.
//
// This mirrors routesListWithKnownAPIValues: walk the state object, replace only
// the attributes that are still unknown, and leave configured values alone so
// the plan stays clean.
func searchUsageGroupRulesWithKnownAPIValues(apiRules types.Object, stateRules types.Object) types.Object {
	if apiRules.IsNull() || apiRules.IsUnknown() || stateRules.IsNull() || stateRules.IsUnknown() {
		return stateRules
	}

	ctx := context.Background()
	stateAttributes := stateRules.Attributes()
	apiAttributes := apiRules.Attributes()
	changed := false

	for name, stateValue := range stateAttributes {
		stateRule, ok := stateValue.(types.Object)
		if !ok || stateRule.IsNull() || stateRule.IsUnknown() {
			continue
		}
		apiRule, ok := apiAttributes[name].(types.Object)
		if !ok || apiRule.IsNull() || apiRule.IsUnknown() {
			continue
		}
		merged, ok := ruleWithKnownAPIValues(ctx, apiRule, stateRule)
		if !ok {
			continue
		}
		stateAttributes[name] = merged
		changed = true
	}

	if !changed {
		return stateRules
	}
	merged, diags := types.ObjectValue(stateRules.AttributeTypes(ctx), stateAttributes)
	if diags.HasError() {
		return stateRules
	}
	return merged
}

// ruleWithKnownAPIValues copies the API's value for every attribute the state
// still reports as unknown. It reports false when nothing changed.
func ruleWithKnownAPIValues(ctx context.Context, apiRule, stateRule types.Object) (types.Object, bool) {
	attributes := stateRule.Attributes()
	apiRuleAttributes := apiRule.Attributes()
	changed := false

	for name, value := range attributes {
		if !value.IsUnknown() {
			continue
		}
		apiValue, ok := apiRuleAttributes[name]
		if !ok || apiValue.IsUnknown() {
			continue
		}
		attributes[name] = apiValue
		changed = true
	}
	if !changed {
		return stateRule, false
	}

	merged, diags := types.ObjectValue(stateRule.AttributeTypes(ctx), attributes)
	if diags.HasError() {
		return stateRule, false
	}
	return merged, true
}

var _ = attr.Value(types.StringNull())
