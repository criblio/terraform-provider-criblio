package provider

import (
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
