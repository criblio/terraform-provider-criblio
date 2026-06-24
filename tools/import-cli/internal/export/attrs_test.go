package export

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHclOptionsForType_criblioDestination_skipsHoistedRootAttrs(t *testing.T) {
	e := registry.Entry{
		OneOf: &registry.OneOfConfig{ReadOnlyAttr: "items"},
	}
	opts := hclOptionsForType("criblio_destination", e)
	require.NotNil(t, opts)
	require.NotNil(t, opts.SkipAttributes)
	assert.True(t, opts.SkipAttributes["items"])
	assert.True(t, opts.SkipAttributes["environment"])
	assert.True(t, opts.SkipAttributes["pipeline"])
	assert.True(t, opts.SkipAttributes["type"])
}

func TestHclOptionsForType_searchMacroSkipsComputedOnlyAttrs(t *testing.T) {
	opts := hclOptionsForType("criblio_search_macro", registry.Entry{})
	require.NotNil(t, opts)

	model := &provider.SearchMacroResourceModel{
		Created:     types.Int64Value(1771983246699),
		CreatedBy:   types.StringValue("user@example.com"),
		Description: types.StringValue("Filters to high-severity events."),
		GroupID:     types.StringValue("default_search"),
		ID:          types.StringValue("test_macro_2"),
		Modified:    types.Int64Value(1771983246699),
		Replacement: types.StringValue(`severity >= "Error"`),
		Tags:        types.StringValue("errors,prod"),
	}

	attrs, err := hcl.ModelToValue(model, opts)
	require.NoError(t, err)

	assert.NotContains(t, attrs, "created")
	assert.NotContains(t, attrs, "created_by")
	assert.NotContains(t, attrs, "modified")
	assert.Equal(t, "default_search", attrs["group_id"].String)
	assert.Equal(t, "test_macro_2", attrs["id"].String)
	assert.Equal(t, `severity >= "Error"`, attrs["replacement"].String)
}

func TestHclOptionsForType_searchEngineSkipsComputedOnlyAttrs(t *testing.T) {
	opts := hclOptionsForType("criblio_search_engine", registry.Entry{})
	require.NotNil(t, opts)

	model := &provider.SearchEngineResourceModel{
		ActiveWorkflow:         types.StringValue("provision"),
		Datasets:               []types.String{types.StringValue("main")},
		DeletionStartedAt:      types.Int64Value(1771983246699),
		Description:            types.StringValue("Search engine"),
		EffectiveStatus:        types.StringValue("ready"),
		EngineType:             types.StringValue("local"),
		GroupID:                types.StringValue("default_search"),
		HasMain:                types.BoolValue(true),
		ID:                     types.StringValue("local_ingest_primary"),
		IsComputeDeprovisioned: types.BoolValue(false),
		IsStorageDeprovisioned: types.BoolValue(false),
		LastProvisionedMs:      types.Int64Value(1771983246699),
		MetricsLastPublishedAt: types.Int64Value(1771983246699),
		Status:                 types.StringValue("ready"),
		TierSize:               types.StringValue("small"),
	}

	attrs, err := hcl.ModelToValue(model, opts)
	require.NoError(t, err)

	for _, name := range []string{
		"active_workflow",
		"datasets",
		"deletion_started_at",
		"effective_status",
		"engine_type",
		"has_main",
		"is_compute_deprovisioned",
		"is_storage_deprovisioned",
		"last_provisioned_ms",
		"metrics_last_published_at",
		"status",
	} {
		assert.NotContains(t, attrs, name)
	}
	assert.Equal(t, "default_search", attrs["group_id"].String)
	assert.Equal(t, "local_ingest_primary", attrs["id"].String)
	assert.Equal(t, "Search engine", attrs["description"].String)
	assert.Equal(t, "small", attrs["tier_size"].String)
}

func TestHclOptionsForType_criblLakeHouseSkipsStatus(t *testing.T) {
	opts := hclOptionsForType("criblio_cribl_lake_house", registry.Entry{})
	require.NotNil(t, opts)

	model := &provider.CriblLakeHouseResourceModel{
		Description: types.StringValue("Lakehouse"),
		ID:          types.StringValue("test-lakehouse"),
		Status:      types.StringValue("provisioning"),
		TierSize:    types.StringValue("medium"),
	}

	attrs, err := hcl.ModelToValue(model, opts)
	require.NoError(t, err)

	assert.NotContains(t, attrs, "status")
	assert.Equal(t, "test-lakehouse", attrs["id"].String)
	assert.Equal(t, "Lakehouse", attrs["description"].String)
	assert.Equal(t, "medium", attrs["tier_size"].String)
}
