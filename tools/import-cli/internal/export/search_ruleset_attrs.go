package export

import (
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// injectSearchRulesetRulesForExport ensures required `rules` appears in attrs for HCL export.
// ModelToValue omits nil slices as null, then ResourceBlock skips nulls — Terraform then reports
// missing required `rules`. Refresh only fills `items` for datatype ruleset; copy rules from items
// here. For dataset ruleset, normalize nil Rules to an empty non-nil slice before conversion.
func injectSearchRulesetRulesForExport(typeName string, model interface{}, attrs map[string]hcl.Value, e registry.Entry) error {
	opts := hclOptionsForType(typeName, e)
	switch typeName {
	case "criblio_search_datatype_ruleset":
		pm, ok := model.(*provider.SearchDatatypeRulesetResourceModel)
		if !ok {
			return nil
		}
		rules := datatypeRulesForExport(pm)
		wrap := struct {
			Rules []tfTypes.DatatypeRule `tfsdk:"rules"`
		}{Rules: rules}
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
		rules := pm.Rules
		if rules == nil {
			rules = []tfTypes.DatasetRule{}
		}
		wrap := struct {
			Rules []tfTypes.DatasetRule `tfsdk:"rules"`
		}{Rules: rules}
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

func datatypeRulesForExport(pm *provider.SearchDatatypeRulesetResourceModel) []tfTypes.DatatypeRule {
	if len(pm.Rules) > 0 {
		return pm.Rules
	}
	if len(pm.Items) == 0 {
		return []tfTypes.DatatypeRule{}
	}
	if len(pm.Items) == 1 {
		return normalizeDatatypeRuleSlice(pm.Items[0].Rules)
	}
	for _, it := range pm.Items {
		if it.ID.ValueString() == string(shared.DatatypeRulesetIDDefault) {
			return normalizeDatatypeRuleSlice(it.Rules)
		}
	}
	return normalizeDatatypeRuleSlice(pm.Items[0].Rules)
}

func normalizeDatatypeRuleSlice(in []tfTypes.DatatypeRule) []tfTypes.DatatypeRule {
	if in == nil {
		return []tfTypes.DatatypeRule{}
	}
	return in
}
