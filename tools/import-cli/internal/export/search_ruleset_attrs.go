package export

import (
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// injectSearchRulesetRulesForExport ensures required `rules` appears in attrs for HCL export.
// ModelToValue omits null lists, then ResourceBlock skips nulls — Terraform then reports
// missing required `rules`. Generated search ruleset resources keep the API rules directly
// on the top-level rules attribute, so inject an empty typed list when the API returned none.
func injectSearchRulesetRulesForExport(typeName string, model interface{}, attrs map[string]hcl.Value, e registry.Entry) error {
	opts := hclOptionsForType(typeName, e)
	switch typeName {
	case "criblio_search_datatype_ruleset":
		pm, ok := model.(*provider.SearchDatatypeRulesetResourceModel)
		if !ok {
			return nil
		}
		wrap := struct {
			Rules types.List `tfsdk:"rules"`
		}{Rules: searchDatatypeRulesForExport(pm)}
		rulesAttrs, err := hcl.ModelToValue(&wrap, opts)
		if err != nil {
			return fmt.Errorf("search datatype ruleset rules to value: %w", err)
		}
		if v, ok := rulesAttrs["rules"]; ok {
			attrs["rules"] = v
		}
	case "criblio_search_dataset_ruleset":
		pm, ok := model.(*provider.SearchDatasetRulesetResourceModel)
		if !ok {
			return nil
		}
		wrap := struct {
			Rules types.List `tfsdk:"rules"`
		}{Rules: searchDatasetRulesForExport(pm)}
		rulesAttrs, err := hcl.ModelToValue(&wrap, opts)
		if err != nil {
			return fmt.Errorf("search dataset ruleset rules to value: %w", err)
		}
		if v, ok := rulesAttrs["rules"]; ok {
			attrs["rules"] = v
		}
	default:
		return nil
	}
	return nil
}

func searchDatasetRulesForExport(pm *provider.SearchDatasetRulesetResourceModel) types.List {
	if pm != nil && !pm.Rules.IsNull() && !pm.Rules.IsUnknown() {
		return pm.Rules
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: provider.SearchDatasetRulesetRulesAttrTypes()}, nil)
}

func searchDatatypeRulesForExport(pm *provider.SearchDatatypeRulesetResourceModel) types.List {
	if pm != nil && !pm.Rules.IsNull() && !pm.Rules.IsUnknown() {
		return pm.Rules
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: provider.SearchDatatypeRulesetRulesAttrTypes()}, nil)
}
