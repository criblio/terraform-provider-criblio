package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSearchDatasetIDFallsBackToActiveVariant(t *testing.T) {
	model := SearchDatasetModel{
		ID: types.StringNull(),
		DatasetS3: &DatasetS3Model{
			ID: types.StringValue("s3-dataset-001"),
		},
	}

	if got := searchDatasetID(model); got != "s3-dataset-001" {
		t.Fatalf("searchDatasetID() = %q, want %q", got, "s3-dataset-001")
	}
}

func TestSearchDatasetV2PathsConvertToAPIPayload(t *testing.T) {
	filterTypes := DatasetS3PathsFiltersAttrTypes()
	filter := types.ObjectValueMust(filterTypes, map[string]attr.Value{
		"data_path_format":      types.StringNull(),
		"data_type_id":          types.StringValue("generic_ndjson"),
		"filter":                types.StringValue("**"),
		"preprocess_outer_json": types.BoolNull(),
	})
	pathTypes := DatasetS3PathsAttrTypes()
	path := types.ObjectValueMust(pathTypes, map[string]attr.Value{
		"auto_detect_region":  types.BoolValue(true),
		"bucket":              types.StringValue("logs"),
		"filters":             types.ListValueMust(types.ObjectType{AttrTypes: filterTypes}, []attr.Value{filter}),
		"partitioning_scheme": types.StringValue("none"),
		"region":              types.StringValue("us-east-1"),
	})

	payload, err := DatasetS3Model{
		Type:          types.StringValue("s3"),
		ID:            types.StringValue("search_dataset_v2"),
		ProviderID:    types.StringValue("S3"),
		SearchVersion: types.StringValue("v2"),
		Paths:         types.ListValueMust(types.ObjectType{AttrTypes: pathTypes}, []attr.Value{path}),
	}.terraformPayload()
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"type":          "s3",
		"id":            "search_dataset_v2",
		"provider":      "S3",
		"searchVersion": "v2",
		"paths": []any{map[string]any{
			"autoDetectRegion":   true,
			"bucket":             "logs",
			"filters":            []any{map[string]any{"dataTypeId": "generic_ndjson", "filter": "**"}},
			"partitioningScheme": "none",
			"region":             "us-east-1",
		}},
	}, payload)
}

func TestSearchDatasetLakehousePayloadUsesCriblSearchVariant(t *testing.T) {
	const response = `{
		"breakerRulesets":["Cribl Search"],
		"staleChannelFlushMs":10000,
		"metadata":{"enableAcceleration":false,"latestRunInfo":{}},
		"retentionPeriod":3650,
		"expectedRelativeTimeRange":{"latest":"1d"},
		"skipEventTimeFilter":false,
		"storageClasses":[],
		"partitioningScheme":"none",
		"searchVersion":"v1",
		"filter":"true",
		"autoDetectRegion":true,
		"provider":"lakehouse",
		"engine":"tf_testing",
		"eventStorageSchemaVersion":1,
		"id":"test_engine_dataset",
		"engineDeleted":false,
		"favorites":[],
		"favoriteCount":0,
		"isFavorited":false
	}`

	var model SearchDatasetModel
	require.NoError(t, json.Unmarshal([]byte(response), &model))
	require.NotNil(t, model.DatasetCriblSearch)
	require.Equal(t, "test_engine_dataset", model.DatasetCriblSearch.ID.ValueString())
	require.Equal(t, "lakehouse", model.DatasetCriblSearch.ProviderID.ValueString())
	require.Equal(t, "tf_testing", model.DatasetCriblSearch.Engine.ValueString())
	require.Equal(t, int64(3650), model.DatasetCriblSearch.RetentionPeriod.ValueInt64())

	payload, err := model.DatasetCriblSearch.terraformPayload()
	require.NoError(t, err)
	require.Equal(t, "lakehouse", payload["provider"])
	require.Equal(t, "tf_testing", payload["engine"])
	require.Equal(t, int64(10000), payload["staleChannelFlushMs"])
	require.Equal(t, int64(3650), payload["retentionPeriod"])
	require.Equal(t, int64(1), payload["eventStorageSchemaVersion"])
	require.Equal(t, map[string]any{"latest": "1d"}, payload["expectedRelativeTimeRange"])
	require.Equal(t, false, payload["skipEventTimeFilter"])
	require.Equal(t, []any{}, payload["storageClasses"])
	require.Equal(t, "none", payload["partitioningScheme"])
	require.Equal(t, "v1", payload["searchVersion"])
	require.Equal(t, "true", payload["filter"])
	require.Equal(t, true, payload["autoDetectRegion"])
}
