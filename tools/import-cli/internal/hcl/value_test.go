package hcl

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
		SensitivePaths:       map[string]bool{"id": true},
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
		Conf:    emptyPipelineConf(),
	}
	out, err := ModelToValue(model, nil)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, KindMap, out["conf"].Kind)
	assert.NotNil(t, out["conf"].Map)
}

func TestModelToValue_pipelineNestedNormalizedFunctionConf(t *testing.T) {
	function := types.ObjectValueMust(provider.PipelineConfFunctionsAttrTypes(), map[string]attr.Value{
		"filter":      types.StringValue("true"),
		"id":          types.StringValue("eval"),
		"description": types.StringNull(),
		"disabled":    types.BoolValue(false),
		"final":       types.BoolValue(false),
		"conf":        jsontypes.NewNormalizedValue(`{"add":[{"name":"foo","value":"bar"}]}`),
		"group_id":    types.StringNull(),
	})
	functions := types.ListValueMust(
		types.ObjectType{AttrTypes: provider.PipelineConfFunctionsAttrTypes()},
		[]attr.Value{function},
	)
	model := &provider.PipelineResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("pipe-1"),
		Conf: types.ObjectValueMust(provider.PipelineConfAttrTypes(), map[string]attr.Value{
			"async_func_timeout": types.Int64Null(),
			"output":             types.StringNull(),
			"description":        types.StringNull(),
			"streamtags":         types.ListNull(types.StringType),
			"functions":          functions,
			"groups":             types.MapNull(types.ObjectType{AttrTypes: provider.PipelineConfGroupsAttrTypes()}),
		}),
	}

	out, err := ModelToValue(model, nil)
	require.NoError(t, err)

	functionValues := out["conf"].Map["functions"].List
	require.Len(t, functionValues, 1)
	conf := functionValues[0].Map["conf"]
	assert.Equal(t, KindString, conf.Kind)
	assert.JSONEq(t, `{"add":[{"name":"foo","value":"bar"}]}`, conf.String)
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

func TestReplaceSecretValuesWithVariableRefs_token_timeout_secs_not_secret(t *testing.T) {
	// token_timeout_secs must not be treated as .token (substring match); export should not create tfvars for it.
	attrs := map[string]Value{
		"output_open_telemetry": {
			Kind: KindMap,
			Map: map[string]Value{
				"token_timeout_secs": {Kind: KindString, String: "3600"},
			},
		},
	}
	used := ReplaceSecretValuesWithVariableRefs(attrs, "destination_default_elastic_otel")
	assert.Empty(t, used)
	assert.Equal(t, KindString, attrs["output_open_telemetry"].Map["token_timeout_secs"].Kind)
	assert.Equal(t, "3600", attrs["output_open_telemetry"].Map["token_timeout_secs"].String)
}

func TestReplaceSecretValuesWithVariableRefs_password(t *testing.T) {
	// password fields (e.g. in Redis functions) should be treated as secrets.
	attrs := map[string]Value{
		"conf": {
			Kind: KindMap,
			Map: map[string]Value{
				"functions": {
					Kind: KindList,
					List: []Value{
						{
							Kind: KindMap,
							Map: map[string]Value{
								"id":       {Kind: KindString, String: "redis-1"},
								"password": {Kind: KindString, String: "test-password"},
							},
						},
					},
				},
			},
		},
	}
	used := ReplaceSecretValuesWithVariableRefs(attrs, "pipeline_default_redis_populater")
	require.Len(t, used, 1)
	assert.Contains(t, used[0], "password")
	// password should now be KindSensitive with variable name
	pwVal := attrs["conf"].Map["functions"].List[0].Map["password"]
	assert.Equal(t, KindSensitive, pwVal.Kind)
	assert.True(t, isVariableName(pwVal.Sensitive))
}

func TestReplaceSecretValuesWithVariableRefs_sensitive_attributes(t *testing.T) {
	// Various sensitive attribute names should be treated as secrets.
	testCases := []struct {
		attrName string
		value    string
	}{
		{"password", "secret-pass"},
		{"api_key", "my-api-key"},
		{"secret", "my-secret"},
		{"auth_token", "my-auth-token"},
		{"client_secret", "my-client-secret"},
		{"access_key", "my-access-key"},
		{"secret_key", "my-secret-key"},
	}
	for _, tc := range testCases {
		t.Run(tc.attrName, func(t *testing.T) {
			attrs := map[string]Value{
				"config": {
					Kind: KindMap,
					Map: map[string]Value{
						tc.attrName: {Kind: KindString, String: tc.value},
					},
				},
			}
			used := ReplaceSecretValuesWithVariableRefs(attrs, "resource_test")
			require.Len(t, used, 1)
			assert.Contains(t, used[0], tc.attrName)
			val := attrs["config"].Map[tc.attrName]
			assert.Equal(t, KindSensitive, val.Kind)
		})
	}
}

func TestReplaceSecretValuesWithVariableRefs_json_with_password(t *testing.T) {
	// JSON string containing password (e.g. pipeline function conf with redis password).
	// The password should be replaced with a variable reference within the JSON.
	attrs := map[string]Value{
		"conf": {
			Kind: KindMap,
			Map: map[string]Value{
				"functions": {
					Kind: KindList,
					List: []Value{
						{
							Kind: KindMap,
							Map: map[string]Value{
								"id":   {Kind: KindString, String: "redis"},
								"conf": {Kind: KindString, String: `{"authType":"manual","commands":[{"command":"get","keyExpr":"test"}],"password":"test-password","url":"localhost:6379"}`},
							},
						},
					},
				},
			},
		},
	}
	used := ReplaceSecretValuesWithVariableRefs(attrs, "pipeline_default_redis_test")
	require.Len(t, used, 1)
	// Variable name should include the sensitive field name (password)
	assert.Contains(t, used[0], "password")
	// The conf JSON string should remain KindString but with password as variable reference
	confVal := attrs["conf"].Map["functions"].List[0].Map["conf"]
	assert.Equal(t, KindString, confVal.Kind)
	assert.Contains(t, confVal.String, `"password":"${var.`)
	assert.Contains(t, confVal.String, `"commands"`)
	assert.NotContains(t, confVal.String, "test-password")
	// MaskedVarNames should be populated
	assert.Len(t, confVal.MaskedVarNames, 1)
	assert.Contains(t, confVal.MaskedVarNames[0], "password")
}

func TestReplaceSecretValuesWithVariableRefs_json_without_password(t *testing.T) {
	// JSON string without sensitive fields should NOT be modified.
	attrs := map[string]Value{
		"conf": {
			Kind: KindMap,
			Map: map[string]Value{
				"functions": {
					Kind: KindList,
					List: []Value{
						{
							Kind: KindMap,
							Map: map[string]Value{
								"id":   {Kind: KindString, String: "eval"},
								"conf": {Kind: KindString, String: `{"add":[{"name":"test","value":"123"}]}`},
							},
						},
					},
				},
			},
		},
	}
	used := ReplaceSecretValuesWithVariableRefs(attrs, "pipeline_default_test")
	assert.Empty(t, used)
	// The conf JSON string should remain as-is
	confVal := attrs["conf"].Map["functions"].List[0].Map["conf"]
	assert.Equal(t, KindString, confVal.Kind)
	assert.Contains(t, confVal.String, "test")
}

func TestReplaceSecretValuesWithVariableRefs_nested_json_password(t *testing.T) {
	// JSON with password nested inside an object - should have variable reference.
	attrs := map[string]Value{
		"config": {
			Kind: KindMap,
			Map: map[string]Value{
				"settings": {Kind: KindString, String: `{"database":{"host":"localhost","password":"db-secret"}}`},
			},
		},
	}
	used := ReplaceSecretValuesWithVariableRefs(attrs, "resource_test")
	require.Len(t, used, 1)
	// Variable name should include nested path and sensitive field name
	assert.Contains(t, used[0], "database_password")
	val := attrs["config"].Map["settings"]
	assert.Equal(t, KindString, val.Kind)
	assert.Contains(t, val.String, `"password":"${var.`)
	assert.Contains(t, val.String, `"host":"localhost"`)
	assert.NotContains(t, val.String, "db-secret")
}

func TestMaskSensitiveValuesInJSON(t *testing.T) {
	// Test the MaskSensitiveValuesInJSON function directly.
	testCases := []struct {
		name         string
		input        string
		resourceName string
		attrPath     string
		expectVars   []string
	}{
		{
			name:         "password at top level",
			input:        `{"password":"secret","url":"localhost"}`,
			resourceName: "test_resource",
			attrPath:     "conf",
			expectVars:   []string{"test_resource_conf_password"},
		},
		{
			name:         "password nested",
			input:        `{"db":{"password":"secret","host":"localhost"}}`,
			resourceName: "test_resource",
			attrPath:     "conf",
			expectVars:   []string{"test_resource_conf_db_password"},
		},
		{
			name:         "no sensitive fields",
			input:        `{"host":"localhost","port":6379}`,
			resourceName: "test_resource",
			attrPath:     "conf",
			expectVars:   nil,
		},
		{
			name:         "multiple sensitive fields",
			input:        `{"password":"pass1","api_key":"key1","host":"localhost"}`,
			resourceName: "test_resource",
			attrPath:     "conf",
			expectVars:   []string{"test_resource_conf_api_key", "test_resource_conf_password"},
		},
		{
			name:         "authorization header in list",
			input:        `{"collect_request_headers":[{"name":"Authorization","value":"Bearer secret123"},{"name":"Content-Type","value":"application/json"}]}`,
			resourceName: "test_collector",
			attrPath:     "conf",
			expectVars:   []string{"test_collector_conf_collect_request_headers_0_value_Authorization"},
		},
		{
			name:         "x-api-key header",
			input:        `{"headers":[{"name":"X-Api-Key","value":"my-api-key"}]}`,
			resourceName: "test_resource",
			attrPath:     "conf",
			expectVars:   []string{"test_resource_conf_headers_0_value_X_Api_Key"},
		},
		{
			name:         "non-sensitive header not masked",
			input:        `{"headers":[{"name":"Content-Type","value":"application/json"}]}`,
			resourceName: "test_resource",
			attrPath:     "conf",
			expectVars:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, varNames := MaskSensitiveValuesInJSON(tc.input, tc.resourceName, tc.attrPath)
			if tc.expectVars == nil {
				assert.Nil(t, varNames)
				assert.Equal(t, tc.input, result)
			} else {
				assert.Len(t, varNames, len(tc.expectVars))
				for _, ev := range tc.expectVars {
					assert.Contains(t, varNames, ev)
				}
				// Check that variable references are in the JSON
				for _, varName := range varNames {
					assert.Contains(t, result, "${var."+varName+"}")
				}
				assert.NotContains(t, result, "secret")
				assert.NotContains(t, result, "pass1")
				assert.NotContains(t, result, "key1")
			}
		})
	}
}

// TestModelToValue_complex_nested validates lists, maps, and nested objects.
func TestModelToValue_complex_nested(t *testing.T) {
	model := &provider.PipelineResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("pipe-1"),
		Conf:    emptyPipelineConf(),
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

func emptyPipelineConf() types.Object {
	return types.ObjectValueMust(provider.PipelineConfAttrTypes(), map[string]attr.Value{
		"async_func_timeout": types.Int64Null(),
		"output":             types.StringNull(),
		"description":        types.StringNull(),
		"streamtags":         types.ListNull(types.StringType),
		"functions":          types.ListNull(types.ObjectType{AttrTypes: provider.PipelineConfFunctionsAttrTypes()}),
		"groups":             types.MapNull(types.ObjectType{AttrTypes: provider.PipelineConfGroupsAttrTypes()}),
	})
}

func TestReplaceSecretValuesWithVariableRefs_masks_authorization_header(t *testing.T) {
	// Simulate what ModelToValue produces for collect_request_headers with Authorization header
	attrs := map[string]Value{
		"collect_request_headers": {
			Kind: KindList,
			List: []Value{
				{
					Kind: KindMap,
					Map: map[string]Value{
						"name":  {Kind: KindString, String: "Authorization"},
						"value": {Kind: KindString, String: "Bearer secret-token"},
					},
				},
			},
		},
	}

	used := ReplaceSecretValuesWithVariableRefs(attrs, "test_resource")

	// The value should be masked
	require.Len(t, used, 1, "should have one masked variable")
	header := attrs["collect_request_headers"].List[0]
	assert.Equal(t, KindString, header.Map["name"].Kind, "name should remain as string")
	assert.Equal(t, "Authorization", header.Map["name"].String)
	assert.Equal(t, KindSensitive, header.Map["value"].Kind, "value should be marked sensitive")
}

func TestReplaceSecretValuesWithVariableRefs_preserves_non_sensitive_headers(t *testing.T) {
	attrs := map[string]Value{
		"headers": {
			Kind: KindList,
			List: []Value{
				{
					Kind: KindMap,
					Map: map[string]Value{
						"name":  {Kind: KindString, String: "Content-Type"},
						"value": {Kind: KindString, String: "application/json"},
					},
				},
			},
		},
	}

	used := ReplaceSecretValuesWithVariableRefs(attrs, "test_resource")

	assert.Empty(t, used, "no variables should be masked for non-sensitive headers")
	header := attrs["headers"].List[0]
	assert.Equal(t, KindString, header.Map["value"].Kind, "value should remain as string")
	assert.Equal(t, "application/json", header.Map["value"].String)
}

func TestReplaceSecretValuesWithVariableRefs_masks_xapikey_header(t *testing.T) {
	attrs := map[string]Value{
		"extra_headers": {
			Kind: KindList,
			List: []Value{
				{
					Kind: KindMap,
					Map: map[string]Value{
						"name":  {Kind: KindString, String: "X-API-Key"},
						"value": {Kind: KindString, String: "my-api-key"},
					},
				},
			},
		},
	}

	used := ReplaceSecretValuesWithVariableRefs(attrs, "test_resource")

	require.Len(t, used, 1)
	header := attrs["extra_headers"].List[0]
	assert.Equal(t, KindSensitive, header.Map["value"].Kind)
}

func TestReplaceSecretValuesWithVariableRefs_masks_all_extra_http_headers(t *testing.T) {
	// All values in extra_http_headers should be masked, regardless of header name
	// This matches UI behavior where extra HTTP headers are treated as sensitive
	// Uses KindVariableRef (plain var.x) since header values are plain strings, not JSON.
	attrs := map[string]Value{
		"extra_http_headers": {
			Kind: KindList,
			List: []Value{
				{
					Kind: KindMap,
					Map: map[string]Value{
						"name":  {Kind: KindString, String: "X-google-api-key"},
						"value": {Kind: KindString, String: "tests"},
					},
				},
				{
					Kind: KindMap,
					Map: map[string]Value{
						"name":  {Kind: KindString, String: "X-Webhook-Access-Key"},
						"value": {Kind: KindString, String: "test"},
					},
				},
			},
		},
	}

	used := ReplaceSecretValuesWithVariableRefs(attrs, "destination_webhook")

	require.Len(t, used, 2, "should mask both header values")
	header0 := attrs["extra_http_headers"].List[0]
	header1 := attrs["extra_http_headers"].List[1]

	// Should use KindVariableRef (plain var.x) not KindSensitive (jsonencode)
	assert.Equal(t, KindVariableRef, header0.Map["value"].Kind, "first header value should use plain var ref")
	assert.Equal(t, KindVariableRef, header1.Map["value"].Kind, "second header value should use plain var ref")

	// Verify variable names include header name for uniqueness
	assert.Contains(t, header0.Map["value"].VarName, "X_google_api_key", "var name should include sanitized header name")
	assert.Contains(t, header1.Map["value"].VarName, "X_Webhook_Access_Key", "var name should include sanitized header name")

	// Verify String is cleared (value moved to VarName)
	assert.Empty(t, header0.Map["value"].String, "original string should be cleared")
	assert.Empty(t, header1.Map["value"].String, "original string should be cleared")

	// Verify "name" field is unchanged (not re-processed as sensitive)
	assert.Equal(t, KindString, header0.Map["name"].Kind, "header name should remain as string")
	assert.Equal(t, "X-google-api-key", header0.Map["name"].String, "header name value should be unchanged")
}

func TestPruneEmptyLists_filters_urls_with_empty_url(t *testing.T) {
	// Simulate what ModelToValue produces for output_webhook with empty urls
	// When loadBalanced=false, the model contains urls: [{url: "", weight: 1}]
	outputWebhook := Value{
		Kind: KindMap,
		Map: map[string]Value{
			"url": {Kind: KindString, String: "http://localhost:9000"},
			"urls": {
				Kind: KindList,
				List: []Value{
					{
						Kind: KindMap,
						Map: map[string]Value{
							"url":    {Kind: KindString, String: ""},
							"weight": {Kind: KindNumber, Number: 1},
						},
					},
				},
			},
			"type": {Kind: KindString, String: "webhook"},
		},
	}

	pruned := PruneEmptyLists(outputWebhook)

	// urls should be filtered out because the url field is empty
	_, hasUrls := pruned.Map["urls"]
	assert.False(t, hasUrls, "urls with empty url items should be filtered out")

	// url (singular) should still be present
	urlVal, hasUrl := pruned.Map["url"]
	assert.True(t, hasUrl, "url (singular) should be preserved")
	assert.Equal(t, "http://localhost:9000", urlVal.String)
}

func TestPruneEmptyLists_preserves_urls_with_valid_url(t *testing.T) {
	outputWebhook := Value{
		Kind: KindMap,
		Map: map[string]Value{
			"urls": {
				Kind: KindList,
				List: []Value{
					{
						Kind: KindMap,
						Map: map[string]Value{
							"url":    {Kind: KindString, String: "http://localhost:9000"},
							"weight": {Kind: KindNumber, Number: 1},
						},
					},
				},
			},
		},
	}

	pruned := PruneEmptyLists(outputWebhook)

	urlsVal, hasUrls := pruned.Map["urls"]
	assert.True(t, hasUrls, "urls with valid url items should be preserved")
	require.Len(t, urlsVal.List, 1)
	assert.Equal(t, "http://localhost:9000", urlsVal.List[0].Map["url"].String)
}

func TestPruneEmptyLists_filters_nested_urls_in_output_block(t *testing.T) {
	// Simulate nested structure like: output_webhook = { urls = [{url: "", weight: 1}] }
	attrs := Value{
		Kind: KindMap,
		Map: map[string]Value{
			"output_webhook": {
				Kind: KindMap,
				Map: map[string]Value{
					"url": {Kind: KindString, String: "http://localhost:9000"},
					"urls": {
						Kind: KindList,
						List: []Value{
							{
								Kind: KindMap,
								Map: map[string]Value{
									"url":    {Kind: KindString, String: ""},
									"weight": {Kind: KindNumber, Number: 1},
								},
							},
						},
					},
				},
			},
		},
	}

	pruned := PruneEmptyLists(attrs)

	// urls inside output_webhook should be filtered
	outputBlock, hasOutput := pruned.Map["output_webhook"]
	require.True(t, hasOutput, "output_webhook should exist")
	_, hasUrls := outputBlock.Map["urls"]
	assert.False(t, hasUrls, "nested urls with empty url items should be filtered out")

	// url (singular) should still be present
	urlVal, hasUrl := outputBlock.Map["url"]
	assert.True(t, hasUrl, "url (singular) should be preserved")
	assert.Equal(t, "http://localhost:9000", urlVal.String)
}
