package export

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestInjectSearchRulesetRulesForExport_datatypeKeepsRules(t *testing.T) {
	e := registryEntryForSearchDatatypeRuleset(t)
	model := &provider.SearchDatatypeRulesetResourceModel{
		ID:    types.StringValue("default"),
		Rules: datatypeRuleList(t),
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_datatype_ruleset", e))
	require.NoError(t, err)

	require.NoError(t, injectSearchRulesetRulesForExport("criblio_search_datatype_ruleset", model, attrs, e))
	rv, ok := attrs["rules"]
	require.True(t, ok)
	require.Equal(t, hcl.KindList, rv.Kind)
	require.Len(t, rv.List, 1)
}

func TestSearchDatasetRulesetExport_skipsReadOnlyItems(t *testing.T) {
	e := registryEntryForSearchDatasetRuleset(t)
	model := &provider.SearchDatasetRulesetResourceModel{
		ID:    types.StringValue("default"),
		Rules: datasetRuleList(t),
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_dataset_ruleset", e))
	require.NoError(t, err)
	_, hasItems := attrs["items"]
	require.False(t, hasItems, "items is Computed-only; export must not emit it")
}

func TestInjectSearchRulesetRulesForExport_datasetKeepsRules(t *testing.T) {
	e := registryEntryForSearchDatasetRuleset(t)
	model := &provider.SearchDatasetRulesetResourceModel{
		ID:    types.StringValue("default"),
		Rules: datasetRuleList(t),
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_dataset_ruleset", e))
	require.NoError(t, err)
	require.NoError(t, injectSearchRulesetRulesForExport("criblio_search_dataset_ruleset", model, attrs, e))
	rv, ok := attrs["rules"]
	require.True(t, ok)
	require.Equal(t, hcl.KindList, rv.Kind)
	require.Len(t, rv.List, 1)
}

func TestInjectSearchRulesetRulesForExport_datasetNilRulesBecomesEmptyList(t *testing.T) {
	e := registryEntryForSearchDatasetRuleset(t)
	model := &provider.SearchDatasetRulesetResourceModel{
		ID:    types.StringValue("default"),
		Rules: types.ListNull(types.ObjectType{AttrTypes: provider.SearchDatasetRulesetRulesAttrTypes()}),
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_dataset_ruleset", e))
	require.NoError(t, err)
	require.NoError(t, injectSearchRulesetRulesForExport("criblio_search_dataset_ruleset", model, attrs, e))
	rv, ok := attrs["rules"]
	require.True(t, ok)
	require.Equal(t, hcl.KindList, rv.Kind)
	require.Len(t, rv.List, 0)
}

func datasetRuleList(t *testing.T) types.List {
	t.Helper()
	value, diags := types.ObjectValue(provider.SearchDatasetRulesetRulesAttrTypes(), map[string]attr.Value{
		"dataset":                   types.StringNull(),
		"description":               types.StringNull(),
		"disabled":                  types.BoolNull(),
		"extend_expression":         types.StringNull(),
		"extend_expression_enabled": types.BoolNull(),
		"id":                        types.StringValue("rule_1"),
		"kusto_expression":          types.StringValue("*"),
		"name":                      types.StringValue("n"),
		"send_data_to":              types.StringValue("destinationDataset"),
	})
	require.False(t, diags.HasError(), "%v", diags)
	list, diags := types.ListValue(types.ObjectType{AttrTypes: provider.SearchDatasetRulesetRulesAttrTypes()}, []attr.Value{value})
	require.False(t, diags.HasError(), "%v", diags)
	return list
}

func datatypeRuleList(t *testing.T) types.List {
	t.Helper()
	value, diags := types.ObjectValue(provider.SearchDatatypeRulesetRulesAttrTypes(), map[string]attr.Value{
		"datatype":         types.StringValue("generic_ndjson"),
		"description":      types.StringNull(),
		"disabled":         types.BoolNull(),
		"id":               types.StringValue("r1"),
		"kusto_expression": types.StringValue("*"),
		"name":             types.StringValue("n"),
	})
	require.False(t, diags.HasError(), "%v", diags)
	list, diags := types.ListValue(types.ObjectType{AttrTypes: provider.SearchDatatypeRulesetRulesAttrTypes()}, []attr.Value{value})
	require.False(t, diags.HasError(), "%v", diags)
	return list
}

func registryEntryForSearchDatatypeRuleset(t *testing.T) registry.Entry {
	t.Helper()
	return registry.Entry{
		TypeName:          "criblio_search_datatype_ruleset",
		ModelTypeName:     "SearchDatatypeRulesetResourceModel",
		GetMethod:         "GetDatatypeRuleByID",
		RefreshFromMethod: "RefreshFromSharedCountedDatatypeRuleset",
		ImportIDFormat:    "id",
	}
}

func registryEntryForSearchDatasetRuleset(t *testing.T) registry.Entry {
	t.Helper()
	return registry.Entry{
		TypeName:          "criblio_search_dataset_ruleset",
		ModelTypeName:     "SearchDatasetRulesetResourceModel",
		GetMethod:         "GetDatasetRuleByID",
		RefreshFromMethod: "RefreshFromSharedCountedDatasetRuleset",
		ImportIDFormat:    "id",
	}
}
