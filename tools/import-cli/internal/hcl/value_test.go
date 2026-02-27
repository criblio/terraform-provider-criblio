package hcl

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	ptypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelToValue_source_simple(t *testing.T) {
	model := &provider.SourceResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("input-1"),
	}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "default", out["group_id"].String)
	assert.Equal(t, "input-1", out["id"].String)
	assert.Equal(t, KindString, out["group_id"].Kind)
	assert.Equal(t, KindString, out["id"].Kind)
}

func TestModelToValue_null_vs_empty(t *testing.T) {
	// Null string
	model := &provider.SourceResourceModel{
		GroupID: types.StringNull(),
		ID:      types.StringValue(""),
	}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	assert.Equal(t, KindNull, out["group_id"].Kind, "null should be preserved")
	assert.Equal(t, KindString, out["id"].Kind)
	assert.Equal(t, "", out["id"].String, "empty string should be preserved")
}

func TestModelToValue_sensitive_placeholder(t *testing.T) {
	model := &provider.SourceResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("secret-id"),
	}
	opts := &Options{
		SensitivePaths:      map[string]bool{"id": true},
		SensitivePlaceholder: "(sensitive)",
	}
	out, err := ModelToValue(model, opts)
	require.NoError(t, err)
	assert.Equal(t, KindString, out["group_id"].Kind)
	assert.Equal(t, "default", out["group_id"].String)
	assert.Equal(t, KindSensitive, out["id"].Kind)
	assert.Equal(t, "(sensitive)", out["id"].Sensitive)
}

func TestModelToValue_sensitive_by_path(t *testing.T) {
	model := &provider.PipelineResourceModel{
		GroupID: types.StringValue("g1"),
		ID:      types.StringValue("p1"),
	}
	opts := &Options{
		SensitivePaths: map[string]bool{"group_id": true},
	}
	out, err := ModelToValue(model, opts)
	require.NoError(t, err)
	assert.True(t, out["group_id"].IsSensitive())
	assert.Equal(t, "p1", out["id"].String)
}

func TestModelToValue_pipeline_nested(t *testing.T) {
	model := &provider.PipelineResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("pipe-1"),
		Conf:    ptypes.PipelineConf{}, // nested object
	}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, KindMap, out["conf"].Kind)
	assert.NotNil(t, out["conf"].Map)
}

func TestToHCLExpr_null_empty_sensitive(t *testing.T) {
	assert.Equal(t, "null", Value{Kind: KindNull}.ToHCLExpr())
	assert.Equal(t, `""`, Value{Kind: KindString, String: ""}.ToHCLExpr())
	assert.Equal(t, `"(sensitive)"`, Value{Kind: KindSensitive, Sensitive: "(sensitive)"}.ToHCLExpr())
}

func TestToHCLExpr_list_and_map(t *testing.T) {
	list := Value{
		Kind: KindList,
		List: []Value{
			{Kind: KindString, String: "a"},
			{Kind: KindNull},
			{Kind: KindString, String: "b"},
		},
	}
	assert.Equal(t, `["a", null, "b"]`, list.ToHCLExpr())

	m := Value{
		Kind: KindMap,
		Map: map[string]Value{
			"x": {Kind: KindString, String: "y"},
			"n": {Kind: KindNull},
		},
	}
	// Map keys sorted in ToHCLExpr
	assert.Contains(t, m.ToHCLExpr(), "n = null")
	assert.Contains(t, m.ToHCLExpr(), `x = "y"`)
}

func TestToHCLExpr_bool_and_number(t *testing.T) {
	assert.Equal(t, "true", Value{Kind: KindBool, Bool: true}.ToHCLExpr())
	assert.Equal(t, "false", Value{Kind: KindBool, Bool: false}.ToHCLExpr())
	assert.Equal(t, "42", Value{Kind: KindNumber, Number: 42}.ToHCLExpr())
	assert.Equal(t, "0", Value{Kind: KindNumber, Number: 0}.ToHCLExpr())
}

func TestAttrToValue_via_model_with_null_pointer(t *testing.T) {
	// SourceResourceModel uses pointer-to-struct oneOf blocks; nil pointer should produce KindNull.
	model := &provider.SourceResourceModel{
		GroupID: types.StringValue("g"),
		ID:      types.StringValue("i"),
	}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	// All oneOf block fields are nil pointers, so they should be KindNull.
	for k, v := range out {
		if k == "group_id" || k == "id" {
			continue
		}
		assert.Equal(t, KindNull, v.Kind, "nil pointer field %q should be KindNull", k)
	}
}

func TestAttrToValue_list_empty_vs_null(t *testing.T) {
	emptyList := types.ListValueMust(types.StringType, []attr.Value{})
	type modelWithList struct {
		Tags types.List `tfsdk:"tags"`
	}
	model := &modelWithList{Tags: emptyList}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	require.NotNil(t, out["tags"])
	assert.Equal(t, KindList, out["tags"].Kind)
	assert.Len(t, out["tags"].List, 0, "empty list preserved")
}

func TestReplaceSecretValuesWithVariableRefs_tokenPlaceholder(t *testing.T) {
	// API returns "yes" as placeholder for masked Splunk HEC token; should be replaced with variable ref.
	attrs := map[string]Value{
		"input_splunk_hec": {
			Kind: KindMap,
			Map: map[string]Value{
				"auth_tokens": {
					Kind: KindList,
					List: []Value{{
						Kind: KindMap,
						Map: map[string]Value{
							"auth_type":   {Kind: KindString, String: "manual"},
							"description": {Kind: KindString, String: "Default token"},
							"token":       {Kind: KindString, String: "yes"},
						},
					}},
				},
			},
		},
	}
	used := ReplaceSecretValuesWithVariableRefs(attrs, "source_default_in_splunk_hec")
	require.Len(t, used, 1)
	assert.Contains(t, used[0], "token")
	// token should now be KindSensitive with variable name
	tokenVal := attrs["input_splunk_hec"].Map["auth_tokens"].List[0].Map["token"]
	assert.Equal(t, KindSensitive, tokenVal.Kind)
	assert.True(t, isVariableName(tokenVal.Sensitive))
}

// TestModelToValue_complex_nested validates lists, maps, and nested objects.
func TestModelToValue_complex_nested(t *testing.T) {
	model := &provider.PipelineResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("pipe-1"),
		Conf:    ptypes.PipelineConf{}, // has map[string]PipelineGroups and other nested fields
	}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	require.NotNil(t, out)
	conf := out["conf"]
	require.Equal(t, KindMap, conf.Kind)
	// conf.Map may contain "groups" (map), "description" (string), etc.
	assert.NotNil(t, conf.Map)
	// Encode to HCL and ensure no panic; structure preserved
	expr := conf.ToHCLExpr()
	assert.True(t, len(expr) > 0)
	assert.Contains(t, expr, "{")
}

