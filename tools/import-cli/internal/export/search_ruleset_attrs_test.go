package export

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	ptypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestInjectSearchRulesetRulesForExport_datatypeFillsRulesFromItems(t *testing.T) {
	e := registryEntryForSearchDatatypeRuleset(t)
	model := &provider.SearchDatatypeRulesetResourceModel{
		ID: types.StringValue("default"),
		Items: []ptypes.DatatypeRuleset{
			{
				ID: types.StringValue("default"),
				Rules: []ptypes.DatatypeRule{
					{
						ID:              types.StringValue("r1"),
						Name:            types.StringValue("n"),
						KustoExpression: types.StringValue("*"),
						Datatype:        types.StringValue("generic_ndjson"),
					},
				},
			},
		},
		// Rules unset (nil) — same as after Refresh without provider-side rules copy.
		Rules: nil,
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_datatype_ruleset", e))
	require.NoError(t, err)
	br, ok := attrs["rules"]
	require.True(t, ok)
	require.Equal(t, hcl.KindNull, br.Kind, "nil Rules slice serializes as null before inject")

	require.NoError(t, injectSearchRulesetRulesForExport("criblio_search_datatype_ruleset", model, attrs, e))
	rv, ok := attrs["rules"]
	require.True(t, ok)
	require.Equal(t, hcl.KindList, rv.Kind)
	require.Len(t, rv.List, 1)
}

func TestSearchDatasetRulesetExport_skipsReadOnlyItems(t *testing.T) {
	e := registryEntryForSearchDatasetRuleset(t)
	model := &provider.SearchDatasetRulesetResourceModel{
		ID: types.StringValue("default"),
		Items: []ptypes.DatasetRuleset{
			{ID: types.StringValue("default")},
		},
		Rules: nil,
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_dataset_ruleset", e))
	require.NoError(t, err)
	_, hasItems := attrs["items"]
	require.False(t, hasItems, "items is Computed-only; export must not emit it")
}

func TestInjectSearchRulesetRulesForExport_datasetFillsRulesFromItems(t *testing.T) {
	e := registryEntryForSearchDatasetRuleset(t)
	model := &provider.SearchDatasetRulesetResourceModel{
		ID: types.StringValue("default"),
		Items: []ptypes.DatasetRuleset{
			{
				ID: types.StringValue("default"),
				Rules: []ptypes.DatasetRule{
					{
						ID:              types.StringValue("rule_1"),
						Name:            types.StringValue("n"),
						KustoExpression: types.StringValue("*"),
						SendDataTo:      types.StringValue("destinationDataset"),
					},
				},
			},
		},
		Rules: nil,
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
		Rules: nil,
	}
	attrs, err := hcl.ModelToValue(model, hclOptionsForType("criblio_search_dataset_ruleset", e))
	require.NoError(t, err)
	require.NoError(t, injectSearchRulesetRulesForExport("criblio_search_dataset_ruleset", model, attrs, e))
	rv, ok := attrs["rules"]
	require.True(t, ok)
	require.Equal(t, hcl.KindList, rv.Kind)
	require.Len(t, rv.List, 0)
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
