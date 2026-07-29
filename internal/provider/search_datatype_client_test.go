package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSearchDatatypeV2Payload(t *testing.T) {
	automaticExtraction := types.ObjectValueMust(SearchDatatypeAutomaticExtractionAttrTypes(), map[string]attr.Value{
		"extraction_type": types.StringValue("json"),
	})
	timestampExtraction := types.ObjectValueMust(SearchDatatypeTimestampExtractionAttrTypes(), map[string]attr.Value{
		"anchor_regex": types.StringValue("/^/"),
		"earliest":     types.StringValue("0"),
		"latest":       types.StringValue("+10years"),
		"scan_depth":   types.Int64Value(150),
		"source_field": types.StringNull(),
		"timezone":     types.StringValue("UTC"),
		"type":         types.StringValue("auto"),
	})

	payload, err := json.Marshal(SearchDatatypeModel{
		AutomaticExtraction: automaticExtraction,
		DataFormat:          types.StringValue("ndjson"),
		ID:                  types.StringValue("test"),
		Lib:                 types.StringValue("custom"),
		MaxEventBytes:       types.Int64Value(65536),
		MinRawLength:        types.Int64Value(0),
		SearchVersion:       types.StringValue("v2"),
		Tags:                types.StringValue(""),
		TimestampExtraction: timestampExtraction,
	})
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(payload, &got))
	require.Equal(t, map[string]any{
		"automaticExtraction": map[string]any{"extractionType": "json"},
		"dataFormat":          "ndjson",
		"id":                  "test",
		"lib":                 "custom",
		"maxEventBytes":       float64(65536),
		"minRawLength":        float64(0),
		"searchVersion":       "v2",
		"tags":                "",
		"timestampExtraction": map[string]any{
			"anchorRegex": "/^/",
			"earliest":    "0",
			"latest":      "+10years",
			"scanDepth":   float64(150),
			"timezone":    "UTC",
			"type":        "auto",
		},
	}, got)
}
