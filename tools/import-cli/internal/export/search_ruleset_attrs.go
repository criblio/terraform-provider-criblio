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
// missing required `rules`. When the model has API data under computed `items` but top-level `rules`
// is still nil (same as datatype), copy rules from items here. Dataset uses datasetRulesForExport
// to mirror provider pickDatasetRulesTFAfterRefresh.
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
		rules := datasetRulesForExport(pm)
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

// datasetRulesForExport mirrors pickDatasetRulesTFAfterRefresh: prefer top-level Rules; if empty, copy from computed Items (default id or sole element).
func datasetRulesForExport(pm *provider.SearchDatasetRulesetResourceModel) []tfTypes.DatasetRule {
	if pm == nil {
		return []tfTypes.DatasetRule{}
	}
	if len(pm.Rules) > 0 {
		return normalizeDatasetRuleSliceForExport(pm.Rules)
	}
	if len(pm.Items) == 0 {
		return []tfTypes.DatasetRule{}
	}
	if len(pm.Items) == 1 {
		return normalizeDatasetRuleSliceForExport(pm.Items[0].Rules)
	}
	for _, it := range pm.Items {
		if it.ID.ValueString() == string(shared.DatasetRulesetIDDefault) {
			return normalizeDatasetRuleSliceForExport(it.Rules)
		}
	}
	return normalizeDatasetRuleSliceForExport(pm.Items[0].Rules)
}

func normalizeDatasetRuleSliceForExport(in []tfTypes.DatasetRule) []tfTypes.DatasetRule {
	if in == nil {
		return []tfTypes.DatasetRule{}
	}
	return in
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
