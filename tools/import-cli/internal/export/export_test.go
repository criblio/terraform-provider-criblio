package export

import (
	"context"
	"errors"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToResourceItems_empty_results(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	results := []discovery.Result{
		{TypeName: "criblio_source", Count: 0},
		{TypeName: "criblio_pipeline", Count: 0},
	}
	result, err := ToResourceItems(ctx, nil, reg, results, []string{"default"}, nil, 1, false, IncludeOverride{}, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.Items)
}

func TestToResourceItems_nil_client_list_skipped(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	results := []discovery.Result{
		{TypeName: "criblio_source", Count: 1},
	}
	result, err := ToResourceItems(ctx, nil, reg, results, []string{"default"}, nil, 1, false, IncludeOverride{}, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.Items)
	// List fails (e.g. nil client) → recorded as list skip, not as returned error.
	assert.Len(t, result.ListSkipped, 1, "list failure should be recorded in ListSkipped")
	assert.Equal(t, "criblio_source", result.ListSkipped[0].TypeName)
}

func buildTestRegistry(t *testing.T) *registry.Registry {
	t.Helper()
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)
	reg, err := registry.NewFromResources(ctx, constructors, registry.MetadataFromProvider(), nil, converter.OneOfBlockNamesFromModel)
	require.NoError(t, err)
	return reg
}

func TestResolveNestedDiscriminator(t *testing.T) {
	t.Run("extracts nested field", func(t *testing.T) {
		itemMap := map[string]string{
			"type":      `"collection"`,
			"collector": `{"type":"rest","conf":{"url":"https://example.com"}}`,
		}
		got := resolveNestedDiscriminator(itemMap, "collector.type")
		assert.Equal(t, `"rest"`, got)
	})
	t.Run("returns empty for missing parent", func(t *testing.T) {
		itemMap := map[string]string{"type": `"collection"`}
		got := resolveNestedDiscriminator(itemMap, "collector.type")
		assert.Empty(t, got)
	})
	t.Run("returns empty for missing inner field", func(t *testing.T) {
		itemMap := map[string]string{
			"collector": `{"conf":{"url":"https://example.com"}}`,
		}
		got := resolveNestedDiscriminator(itemMap, "collector.type")
		assert.Empty(t, got)
	})
	t.Run("returns empty for invalid path (no dot)", func(t *testing.T) {
		itemMap := map[string]string{"type": `"collection"`}
		got := resolveNestedDiscriminator(itemMap, "type")
		assert.Empty(t, got)
	})
	t.Run("returns empty for invalid JSON parent", func(t *testing.T) {
		itemMap := map[string]string{"collector": `not-json`}
		got := resolveNestedDiscriminator(itemMap, "collector.type")
		assert.Empty(t, got)
	})
}

// testModel is a minimal struct matching the shape expected by firstItemMapFromModel.
type testModel struct {
	Items []map[string]jsontypes.Normalized
}

func TestAddOneOfBlockFromFirstItem_nestedDiscriminator(t *testing.T) {
	supportedBlocks := []string{
		"input_collector_azure_blob", "input_collector_cribl_lake",
		"input_collector_database", "input_collector_gcs",
		"input_collector_health_check", "input_collector_rest",
		"input_collector_s3", "input_collector_script", "input_collector_splunk",
	}
	cfg := &registry.OneOfConfig{
		ReadOnlyAttr:                   "items",
		DiscriminatorField:             "type",
		BlockNamePrefix:                "input_collector_",
		KeysToSkip:                     []string{"status"},
		UnsupportedDiscriminatorValues: []string{"scheduledSearch", "executor"},
		NestedDiscriminatorField:       "collector.type",
		SupportedBlockNames:            supportedBlocks,
	}

	t.Run("collection type resolves via nested collector.type=rest", func(t *testing.T) {
		model := testModel{Items: []map[string]jsontypes.Normalized{{
			"type":      jsontypes.NewNormalizedValue(`"collection"`),
			"id":        jsontypes.NewNormalizedValue(`"crowdstrike_ngsiem_api"`),
			"collector": jsontypes.NewNormalizedValue(`{"type":"rest","conf":{"url":"https://api.crowdstrike.com"}}`),
		}}}
		attrs := make(map[string]hcl.Value)
		err := addOneOfBlockFromFirstItem(model, attrs, cfg)
		require.NoError(t, err)
		assert.Contains(t, attrs, "input_collector_rest")
	})

	t.Run("collection type resolves via nested collector.type=s3", func(t *testing.T) {
		model := testModel{Items: []map[string]jsontypes.Normalized{{
			"type":      jsontypes.NewNormalizedValue(`"collection"`),
			"id":        jsontypes.NewNormalizedValue(`"my_s3_collector"`),
			"collector": jsontypes.NewNormalizedValue(`{"type":"s3","conf":{"bucket":"my-bucket"}}`),
		}}}
		attrs := make(map[string]hcl.Value)
		err := addOneOfBlockFromFirstItem(model, attrs, cfg)
		require.NoError(t, err)
		assert.Contains(t, attrs, "input_collector_s3")
	})

	t.Run("scheduledSearch type returns ErrUnsupportedOneOfType", func(t *testing.T) {
		model := testModel{Items: []map[string]jsontypes.Normalized{{
			"type":         jsontypes.NewNormalizedValue(`"scheduledSearch"`),
			"id":           jsontypes.NewNormalizedValue(`"scheduledSearch_test"`),
			"savedQueryId": jsontypes.NewNormalizedValue(`"my_saved_query"`),
		}}}
		attrs := make(map[string]hcl.Value)
		err := addOneOfBlockFromFirstItem(model, attrs, cfg)
		assert.ErrorIs(t, err, ErrUnsupportedOneOfType)
	})

	t.Run("executor type returns ErrUnsupportedOneOfType", func(t *testing.T) {
		model := testModel{Items: []map[string]jsontypes.Normalized{{
			"type": jsontypes.NewNormalizedValue(`"executor"`),
			"id":   jsontypes.NewNormalizedValue(`"my_executor"`),
		}}}
		attrs := make(map[string]hcl.Value)
		err := addOneOfBlockFromFirstItem(model, attrs, cfg)
		assert.ErrorIs(t, err, ErrUnsupportedOneOfType)
	})

	t.Run("empty items returns nil", func(t *testing.T) {
		model := testModel{Items: nil}
		attrs := make(map[string]hcl.Value)
		err := addOneOfBlockFromFirstItem(model, attrs, cfg)
		assert.NoError(t, err)
		assert.Empty(t, attrs)
	})

	t.Run("unsupported nested collector type returns error", func(t *testing.T) {
		model := testModel{Items: []map[string]jsontypes.Normalized{{
			"type":      jsontypes.NewNormalizedValue(`"collection"`),
			"id":        jsontypes.NewNormalizedValue(`"unknown"`),
			"collector": jsontypes.NewNormalizedValue(`{"type":"kafka","conf":{}}`),
		}}}
		attrs := make(map[string]hcl.Value)
		err := addOneOfBlockFromFirstItem(model, attrs, cfg)
		assert.ErrorIs(t, err, ErrUnsupportedOneOfType)
	})
}

func TestAddOneOfBlockFromFirstItem_generatedSourceInputBlock(t *testing.T) {
	cfg := &registry.OneOfConfig{
		ReadOnlyAttr:        "items",
		DiscriminatorField:  "type",
		BlockNamePrefix:     "input_",
		KeysToSkip:          []string{"status"},
		SupportedBlockNames: []string{"input_cribl_http"},
	}
	model := &provider.SourceResourceModel{
		InputCriblHttp: &provider.InputCriblHttpModel{
			Type: types.StringValue("cribl_http"),
			ID:   types.StringValue("in_test"),
			Host: types.StringValue("0.0.0.0"),
			Port: types.Float64Value(10080),
		},
	}
	attrs := make(map[string]hcl.Value)
	err := addOneOfBlockFromFirstItem(model, attrs, cfg)
	require.NoError(t, err)
	v, ok := attrs["input_cribl_http"]
	require.True(t, ok, "expected input_cribl_http block")
	assert.Equal(t, hcl.KindMap, v.Kind)
	assert.Equal(t, "in_test", v.Map["id"].String)
}

func TestSanitizeConvertError(t *testing.T) {
	t.Run("nil returns empty string", func(t *testing.T) {
		assert.Empty(t, sanitizeConvertError(nil))
	})
	t.Run("truncates long error", func(t *testing.T) {
		long := errors.New("a" + string(make([]byte, 500)))
		got := sanitizeConvertError(long)
		assert.LessOrEqual(t, len(got), 120+3) // 120 + "..."
		assert.Contains(t, got, "...")
	})
	t.Run("short error unchanged", func(t *testing.T) {
		err := errors.New("short error")
		assert.Equal(t, "short error", sanitizeConvertError(err))
	})
	t.Run("unmarshal error sanitized", func(t *testing.T) {
		err := errors.New("could not unmarshal json: {\"password\":\"secret\"}")
		assert.Equal(t, "unsupported type (SDK unmarshal failed)", sanitizeConvertError(err))
	})
}

func TestAppendResourceItemFromModel_skipsSearchWorkerGroup(t *testing.T) {
	// criblio_group export used to bypass skipResourceByID (only convertOneResource applied it).
	out := &ExportResult{}
	e := registry.Entry{
		TypeName:       "criblio_group",
		ModelTypeName:  "GroupResourceModel",
		ImportIDFormat: "group_id",
	}
	idMap := map[string]string{"group_id": "search", "product": "stream"}
	err := appendResourceItemFromModel(out, "criblio_group", e, idMap, nil, nil, nil, false, IncludeOverride{})
	require.NoError(t, err)
	assert.Empty(t, out.Items)
}

func TestLifecycleIgnoreChangesForConvertedResource_certificateUsesProviderStatePreservation(t *testing.T) {
	ignored := lifecycleIgnoreChangesForConvertedResource("criblio_certificate", map[string]hcl.Value{
		"cert":     {Kind: hcl.KindVariableRef, VarName: "certificate_cert"},
		"priv_key": {Kind: hcl.KindVariableRef, VarName: "certificate_priv_key"},
	})

	assert.Empty(t, ignored, "certificate should not need import-cli lifecycle ignores with generated REST resource state preservation")
}

func TestHCLOptionsForType_certificateKeepsConfigurableCA(t *testing.T) {
	opts := hclOptionsForType("criblio_certificate", registry.Entry{})
	require.NotNil(t, opts)

	assert.True(t, opts.SkipAttributes["in_use"])
	assert.True(t, opts.SkipAttributes["passphrase"])
	assert.False(t, opts.SkipAttributes["ca"])
}

func TestSkipResourceByID(t *testing.T) {
	t.Run("skip by exclusions.SkipExportIDs", func(t *testing.T) {
		assert.True(t, skipResourceByID("criblio_notification_target", map[string]string{"id": "system_email"}))
		assert.True(t, skipResourceByID("criblio_source", map[string]string{"id": "in_syslog"}))
		assert.True(t, skipResourceByID("criblio_group", map[string]string{"group_id": "search", "product": "stream"}))
	})
	t.Run("skip when id equals group_id", func(t *testing.T) {
		idMap := map[string]string{"group_id": "default", "id": "default"}
		assert.True(t, skipResourceByID("criblio_group", idMap))
	})
	t.Run("skip criblio_pack_lookups when id starts with cribl.", func(t *testing.T) {
		assert.True(t, skipResourceByID("criblio_pack_lookups", map[string]string{"id": "cribl.something"}))
	})
	t.Run("skip criblio_pack_vars when id contains dots", func(t *testing.T) {
		assert.True(t, skipResourceByID("criblio_pack_vars", map[string]string{"id": "foo.bar.baz"}))
	})
	t.Run("skip pack in SkipPacks", func(t *testing.T) {
		for pack := range custom.SkipPacks {
			assert.True(t, skipResourceByID("criblio_pack", map[string]string{"id": pack}))
			break
		}
	})
	t.Run("not skipped for normal resource", func(t *testing.T) {
		assert.False(t, skipResourceByID("criblio_source", map[string]string{"id": "my_custom_source"}))
	})
	t.Run("skip criblio_source in default_search group", func(t *testing.T) {
		assert.True(t, skipResourceByID("criblio_source", map[string]string{"group_id": "default_search", "id": "in_open_telemetry"}))
	})
	t.Run("skip criblio_routes in default_search group", func(t *testing.T) {
		assert.True(t, skipResourceByID("criblio_routes", map[string]string{"group_id": "default_search", "id": "default_search"}))
	})
	t.Run("skip criblio_routes in search group", func(t *testing.T) {
		assert.True(t, skipResourceByID("criblio_routes", map[string]string{"group_id": "search", "id": "search"}))
	})
	t.Run("not skip criblio_source in other groups when same id", func(t *testing.T) {
		assert.False(t, skipResourceByID("criblio_source", map[string]string{"group_id": "default", "id": "in_open_telemetry"}))
	})
	t.Run("not skip criblio_routes in worker group", func(t *testing.T) {
		assert.False(t, skipResourceByID("criblio_routes", map[string]string{"group_id": "default", "id": "default"}))
	})
}

func TestSkipResourceWhenLibCribl(t *testing.T) {
	t.Run("skip when lib is cribl", func(t *testing.T) {
		attrs := map[string]hcl.Value{"lib": {Kind: hcl.KindString, String: "cribl"}}
		assert.True(t, skipResourceWhenLibCribl(attrs))
	})
	t.Run("not skip when lib is other", func(t *testing.T) {
		attrs := map[string]hcl.Value{"lib": {Kind: hcl.KindString, String: "user"}}
		assert.False(t, skipResourceWhenLibCribl(attrs))
	})
	t.Run("not skip when lib missing", func(t *testing.T) {
		attrs := map[string]hcl.Value{}
		assert.False(t, skipResourceWhenLibCribl(attrs))
	})
}

func TestDefaultResource(t *testing.T) {
	noOverride := IncludeOverride{}
	t.Run("skip when lib is cribl", func(t *testing.T) {
		attrs := map[string]hcl.Value{"lib": {Kind: hcl.KindString, String: "cribl"}}
		assert.True(t, DefaultResource("criblio_grok", map[string]string{"id": "custom"}, attrs, noOverride))
	})
	t.Run("skip when tags equals cribl:default", func(t *testing.T) {
		attrs := map[string]hcl.Value{"tags": {Kind: hcl.KindString, String: "cribl:default"}}
		assert.True(t, DefaultResource("criblio_search_dataset", map[string]string{"id": "custom"}, attrs, noOverride))
	})
	t.Run("not skip when tags is partial cribl match", func(t *testing.T) {
		attrs := map[string]hcl.Value{"tags": {Kind: hcl.KindString, String: "my-cribl-app"}}
		assert.False(t, DefaultResource("criblio_appscope_config", map[string]string{"id": "custom"}, attrs, noOverride))
	})
	t.Run("skip when tags list contains cribl:default", func(t *testing.T) {
		attrs := map[string]hcl.Value{"tags": {Kind: hcl.KindList, List: []hcl.Value{
			{Kind: hcl.KindString, String: "other"},
			{Kind: hcl.KindString, String: "cribl:default"},
		}}}
		assert.True(t, DefaultResource("criblio_search_dataset", map[string]string{"id": "custom"}, attrs, noOverride))
	})
	t.Run("skip default pipeline IDs", func(t *testing.T) {
		attrs := map[string]hcl.Value{}
		assert.True(t, DefaultResource("criblio_pipeline", map[string]string{"id": "devnull"}, attrs, noOverride))
		assert.True(t, DefaultResource("criblio_pipeline", map[string]string{"id": "main"}, attrs, noOverride))
		assert.True(t, DefaultResource("criblio_pipeline", map[string]string{"id": "passthru"}, attrs, noOverride))
	})
	t.Run("skip default grok IDs", func(t *testing.T) {
		attrs := map[string]hcl.Value{}
		assert.True(t, DefaultResource("criblio_grok", map[string]string{"id": "aws"}, attrs, noOverride))
		assert.True(t, DefaultResource("criblio_grok", map[string]string{"id": "core-patterns"}, attrs, noOverride))
	})
	t.Run("skip default source IDs", func(t *testing.T) {
		attrs := map[string]hcl.Value{}
		assert.True(t, DefaultResource("criblio_source", map[string]string{"id": "CriblLogs"}, attrs, noOverride))
		assert.True(t, DefaultResource("criblio_source", map[string]string{"id": "CriblMetrics"}, attrs, noOverride))
	})
	t.Run("not skip user-created resources", func(t *testing.T) {
		attrs := map[string]hcl.Value{}
		assert.False(t, DefaultResource("criblio_pipeline", map[string]string{"id": "my_custom_pipeline"}, attrs, noOverride))
		assert.False(t, DefaultResource("criblio_grok", map[string]string{"id": "my_grok"}, attrs, noOverride))
		assert.False(t, DefaultResource("criblio_source", map[string]string{"id": "my_source"}, attrs, noOverride))
	})
	t.Run("not skip when tags is empty list", func(t *testing.T) {
		attrs := map[string]hcl.Value{"tags": {Kind: hcl.KindList, List: []hcl.Value{}}}
		assert.False(t, DefaultResource("criblio_search_dataset", map[string]string{"id": "custom"}, attrs, noOverride))
	})
	t.Run("include override bypasses lib=cribl", func(t *testing.T) {
		attrs := map[string]hcl.Value{"lib": {Kind: hcl.KindString, String: "cribl"}}
		override := ParseIncludeDefaultIDs([]string{"in_system_metrics"})
		assert.False(t, DefaultResource("criblio_source", map[string]string{"id": "in_system_metrics"}, attrs, override))
	})
	t.Run("include override bypasses cribl:default tag", func(t *testing.T) {
		attrs := map[string]hcl.Value{"tags": {Kind: hcl.KindString, String: "cribl:default"}}
		override := ParseIncludeDefaultIDs([]string{"criblio_source:CriblLogs"})
		assert.False(t, DefaultResource("criblio_source", map[string]string{"id": "CriblLogs"}, attrs, override))
		assert.True(t, DefaultResource("criblio_destination", map[string]string{"id": "CriblLogs"}, attrs, override))
	})
	t.Run("include override unqualified matches any type", func(t *testing.T) {
		attrs := map[string]hcl.Value{}
		override := ParseIncludeDefaultIDs([]string{"devnull"})
		assert.False(t, DefaultResource("criblio_pipeline", map[string]string{"id": "devnull"}, attrs, override))
		assert.False(t, DefaultResource("criblio_destination", map[string]string{"id": "devnull"}, attrs, override))
	})
}

func TestGroupIDFromIDMap(t *testing.T) {
	t.Run("returns group_id when present", func(t *testing.T) {
		idMap := map[string]string{"group_id": "default"}
		assert.Equal(t, "default", groupIDFromIDMap(idMap))
	})
	t.Run("returns global when group_id missing", func(t *testing.T) {
		idMap := map[string]string{"id": "foo"}
		assert.Equal(t, "global", groupIDFromIDMap(idMap))
	})
	t.Run("returns global for nil map", func(t *testing.T) {
		assert.Equal(t, "global", groupIDFromIDMap(nil))
	})
}

func TestImportIDMapForType(t *testing.T) {
	t.Run("notification uses legacy group field for import ID", func(t *testing.T) {
		idMap := map[string]string{"group_id": "default_search", "id": "notif-1"}
		got := importIDMapForType("criblio_notification", idMap)
		assert.Equal(t, "default_search", got["group"])
		assert.Equal(t, "default_search", got["group_id"])
		assert.Equal(t, "", idMap["group"])
	})
	t.Run("key maps id to key_id", func(t *testing.T) {
		idMap := map[string]string{"group_id": "default", "id": "key-1"}
		got := importIDMapForType("criblio_key", idMap)
		assert.Equal(t, "key-1", got["key_id"])
		assert.Equal(t, "", idMap["key_id"])
	})
}

func TestSkipExportForGroupFilter(t *testing.T) {
	t.Run("no filter never skips", func(t *testing.T) {
		assert.False(t, skipExportForGroupFilter("criblio_search_dataset", map[string]string{"id": "x"}, nil, []string{"stream-leaders"}))
	})
	t.Run("skip search when only stream group", func(t *testing.T) {
		assert.True(t, skipExportForGroupFilter("criblio_search_dataset", map[string]string{"id": "x"}, []string{"stream-leaders"}, []string{"stream-leaders"}))
	})
	t.Run("include search when default_search selected", func(t *testing.T) {
		assert.False(t, skipExportForGroupFilter("criblio_search_dataset", map[string]string{"id": "x"}, []string{"default_search"}, []string{"default_search"}))
	})
	t.Run("skip global lake when group filter", func(t *testing.T) {
		assert.True(t, skipExportForGroupFilter("criblio_cribl_lake_dataset", map[string]string{"id": "x"}, []string{"stream-leaders"}, []string{"stream-leaders"}))
	})
	t.Run("include worker resource in filtered group", func(t *testing.T) {
		assert.False(t, skipExportForGroupFilter("criblio_source", map[string]string{"group_id": "stream-leaders", "id": "s1"}, []string{"stream-leaders"}, []string{"stream-leaders"}))
	})
}

func TestGroupIDForOutput(t *testing.T) {
	t.Run("criblio_cribl_lake_house returns global", func(t *testing.T) {
		assert.Equal(t, "global", groupIDForOutput("criblio_cribl_lake_house", "my-lake"))
	})
	t.Run("default_search returns search", func(t *testing.T) {
		assert.Equal(t, "search", groupIDForOutput("criblio_default_search", "default_search"))
	})
	t.Run("criblio_search_* returns search", func(t *testing.T) {
		assert.Equal(t, "search", groupIDForOutput("criblio_search_dataset", "global"))
		assert.Equal(t, "search", groupIDForOutput("criblio_search_dataset_provider", "global"))
	})
	t.Run("others return gid", func(t *testing.T) {
		assert.Equal(t, "my-group", groupIDForOutput("criblio_source", "my-group"))
	})
}

func TestToRequestParams(t *testing.T) {
	t.Run("includes GroupID and ID", func(t *testing.T) {
		idMap := map[string]string{"group_id": "default", "id": "my-source"}
		got := toRequestParams(idMap)
		assert.Equal(t, "default", got["GroupID"])
		assert.Equal(t, "my-source", got["ID"])
	})
	t.Run("includes Pack", func(t *testing.T) {
		idMap := map[string]string{"pack": "MyPack", "id": "pipeline-1"}
		got := toRequestParams(idMap)
		assert.Equal(t, "MyPack", got["Pack"])
		assert.Equal(t, "pipeline-1", got["ID"])
	})
	t.Run("includes LakeID", func(t *testing.T) {
		idMap := map[string]string{"lake_id": "lake-1", "id": "ds-1"}
		got := toRequestParams(idMap)
		assert.Equal(t, "lake-1", got["LakeID"])
	})
}

func TestAttrsHasOutputBlock(t *testing.T) {
	t.Run("true when output_ prefix present", func(t *testing.T) {
		attrs := map[string]hcl.Value{"output_s3": {Kind: hcl.KindString, String: "x"}}
		assert.True(t, attrsHasOutputBlock(attrs))
	})
	t.Run("false when no output_ prefix", func(t *testing.T) {
		attrs := map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "x"}}
		assert.False(t, attrsHasOutputBlock(attrs))
	})
	t.Run("false for empty attrs", func(t *testing.T) {
		assert.False(t, attrsHasOutputBlock(nil))
	})
}

func TestFlattenItemsListToTopLevel(t *testing.T) {
	t.Run("merges first list element into attrs", func(t *testing.T) {
		attrs := map[string]hcl.Value{
			"items": {
				Kind: hcl.KindList,
				List: []hcl.Value{{
					Kind: hcl.KindMap,
					Map: map[string]hcl.Value{
						"id":   {Kind: hcl.KindString, String: "a"},
						"name": {Kind: hcl.KindString, String: "n"},
					},
				}},
			},
		}
		flattenItemsListToTopLevel(attrs)
		assert.Contains(t, attrs, "id")
		assert.Equal(t, "a", attrs["id"].String)
		assert.Contains(t, attrs, "name")
		assert.Equal(t, "n", attrs["name"].String)
		assert.NotContains(t, attrs, "items")
	})
	t.Run("no items key leaves attrs unchanged", func(t *testing.T) {
		attrs := map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "x"}}
		flattenItemsListToTopLevel(attrs)
		assert.Equal(t, "x", attrs["id"].String)
	})
	t.Run("empty list removes items key", func(t *testing.T) {
		attrs := map[string]hcl.Value{"items": {Kind: hcl.KindList, List: nil}}
		flattenItemsListToTopLevel(attrs)
		assert.NotContains(t, attrs, "items")
	})
}

func TestFilterAttrsBySchema(t *testing.T) {
	t.Run("allowed attrs kept", func(t *testing.T) {
		allowed := converter.AllAttributeNamesFromModel("SourceResourceModel")
		require.NotEmpty(t, allowed)
		attrs := map[string]hcl.Value{}
		for _, a := range allowed[:3] {
			attrs[a] = hcl.Value{Kind: hcl.KindString, String: "v"}
		}
		filterAttrsBySchema(attrs, "SourceResourceModel")
		assert.Len(t, attrs, 3)
	})
	t.Run("disallowed attrs removed", func(t *testing.T) {
		attrs := map[string]hcl.Value{
			"id":   {Kind: hcl.KindString, String: "x"},
			"fake": {Kind: hcl.KindString, String: "y"},
		}
		filterAttrsBySchema(attrs, "SourceResourceModel")
		assert.Contains(t, attrs, "id")
		assert.NotContains(t, attrs, "fake")
	})
	t.Run("unknown model leaves attrs unchanged", func(t *testing.T) {
		attrs := map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "x"}}
		filterAttrsBySchema(attrs, "unknown_type")
		assert.Contains(t, attrs, "id")
	})
}

func TestHclOptionsForType_skipsGrokComputedAttrs(t *testing.T) {
	opts := hclOptionsForType("criblio_grok", registry.Entry{})
	require.NotNil(t, opts)
	assert.True(t, opts.SkipAttributes["size"])
	assert.True(t, opts.SkipAttributes["tags"])
}

func TestRawJSONToItemMap(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		raw := []byte(`{"id":"x","name":"y"}`)
		got := rawJSONToItemMap(raw)
		require.NotNil(t, got)
		assert.Equal(t, `"x"`, got["id"]) // values are JSON-encoded
		assert.Equal(t, `"y"`, got["name"])
	})
	t.Run("invalid JSON returns nil", func(t *testing.T) {
		got := rawJSONToItemMap([]byte(`{invalid`))
		assert.Nil(t, got)
	})
	t.Run("empty object returns nil", func(t *testing.T) {
		got := rawJSONToItemMap([]byte(`{}`))
		assert.Nil(t, got)
	})
}

func TestTfsdkNameToGoFieldName(t *testing.T) {
	t.Run("capitalizes first letter", func(t *testing.T) {
		assert.Equal(t, "Input_collector", tfsdkNameToGoFieldName("input_collector"))
	})
	t.Run("empty string", func(t *testing.T) {
		assert.Empty(t, tfsdkNameToGoFieldName(""))
	})
}

func TestHclOptionsForType(t *testing.T) {
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_routes")
	require.True(t, ok, "criblio_routes should exist in registry")
	opts := hclOptionsForType("criblio_routes", e)
	require.NotNil(t, opts)
	assert.True(t, opts.SkipAttributes["id"])
	assert.True(t, opts.SkipAttributes["additional_properties"])
}
