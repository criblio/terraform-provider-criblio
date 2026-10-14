package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func mustPackTags(t *testing.T, raw string) types.Dynamic {
	t.Helper()
	tags, err := PackTagsFromAPIJSON([]byte(raw))
	require.NoError(t, err)
	return tags
}

func TestPackCustomTagsSchemaAndValidation(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	(&PackResource{}).Schema(ctx, resource.SchemaRequest{}, &response)
	attribute := response.Schema.Attributes["tags"].(schema.DynamicAttribute)
	require.True(t, attribute.Optional)
	require.True(t, attribute.Computed)
	for _, test := range []struct {
		raw  string
		fail bool
	}{
		{`{"data_type":["logs"],"domain":["observability"],"streamtags":["processing"],"technology":["aws"],"ManagedByTFE":"true","product_id":"123456","used_for":"non_prod"}`, false},
		{`{"Anything":true,"count":42,"metadata":{"OwnerName":"team"},"labels":["one","two"]}`, false},
		{`{"domain":"not-a-list"}`, true},
		{`{"data_type":[1]}`, true},
		{`{"streamtags":[null]}`, true},
		{`{"":"empty-key"}`, true},
		{`{"product_id":null}`, true},
	} {
		var response validator.DynamicResponse
		packTagsValidator{}.ValidateDynamic(ctx, validator.DynamicRequest{
			Path: path.Root("tags"), ConfigValue: mustPackTags(t, test.raw),
		}, &response)
		require.Equal(t, test.fail, response.Diagnostics.HasError(), "%s: %v", test.raw, response.Diagnostics)
	}
	var scalar validator.DynamicResponse
	packTagsValidator{}.ValidateDynamic(ctx, validator.DynamicRequest{
		Path: path.Root("tags"), ConfigValue: types.DynamicValue(types.StringValue("invalid")),
	}, &scalar)
	require.True(t, scalar.Diagnostics.HasError())
	var apiName validator.DynamicResponse
	packTagsValidator{}.ValidateDynamic(ctx, validator.DynamicRequest{
		Path: path.Root("tags"), ConfigValue: types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"dataType": types.ListType{ElemType: types.StringType}}, map[string]attr.Value{"dataType": types.ListValueMust(types.StringType, nil)})),
	}, &apiName)
	require.True(t, apiName.Diagnostics.HasError())
	unknownTags := types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"ManagedByTFE": types.StringType}, map[string]attr.Value{"ManagedByTFE": types.StringUnknown()}))
	var unknown validator.DynamicResponse
	packTagsValidator{}.ValidateDynamic(ctx, validator.DynamicRequest{Path: path.Root("tags"), ConfigValue: unknownTags}, &unknown)
	require.False(t, unknown.Diagnostics.HasError())
	_, err := tagsAPIMap(ctx, unknownTags)
	require.Error(t, err)
	_, err = packTagsAPIFromTF(ctx, types.DynamicUnknown())
	require.NoError(t, err)
	missing := preservePackTags(ctx, mustPackTags(t, `{}`), mustPackTags(t, `{"domain":[],"ManagedByTFE":"true"}`))
	missingValues := missing.UnderlyingValue().(types.Object).Attributes()
	require.True(t, missingValues["ManagedByTFE"].IsNull())
	require.True(t, missingValues["domain"].Equal(types.ListValueMust(types.StringType, nil)))
}

func TestMergePackSettingsCustomTags(t *testing.T) {
	remote := map[string]any{"streamtags": []any{"old-label"}, "ManagedByTFE": "true", "removed": "old", "useCase": []any{"analytics"}}
	for _, test := range []struct {
		name       string
		configured map[string]any
		want       map[string]any
	}{
		{"omitted", nil, remote},
		{"replace", map[string]any{"ManagedByTFE": "TRUE", "product_id": "001234"}, map[string]any{"streamtags": []any{"old-label"}, "ManagedByTFE": "TRUE", "product_id": "001234"}},
		{"clear custom", map[string]any{}, map[string]any{"streamtags": []any{"old-label"}}},
		{"clear built-in", map[string]any{"streamtags": []any{}}, map[string]any{"streamtags": []any{}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, mergePackSettingsTags(remote, test.configured))
			require.Equal(t, "old", remote["removed"])
		})
	}
	_, err := tagsAPIMap(context.Background(), types.DynamicUnknown())
	require.Error(t, err)
}

func TestPackCustomTagsLifecycle(t *testing.T) {
	ctx := context.Background()
	remoteTags := map[string]any{"streamtags": []any{"old-label"}}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/api/v1/m/default/packs", "/api/v1/m/default/packs/test-pack-metadata":
			require.NoError(t, json.NewEncoder(writer).Encode(map[string]any{"items": []any{map[string]any{
				"id": "test-pack-metadata", "version": "1.0.0", "tags": map[string]any{"streamtags": []string{"inventory-label"}},
			}}}))
		case "/api/v1/m/default/p/test-pack-metadata/pack/settings":
			if request.Method == http.MethodPatch {
				var body struct {
					Package struct {
						Tags map[string]any `json:"tags"`
					} `json:"package"`
				}
				decoder := json.NewDecoder(request.Body)
				decoder.UseNumber()
				require.NoError(t, decoder.Decode(&body))
				remoteTags = body.Package.Tags
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			require.Equal(t, http.MethodGet, request.Method)
			require.NoError(t, json.NewEncoder(writer).Encode(map[string]any{"items": []any{map[string]any{
				"id": "settings", "package": map[string]any{"tags": remoteTags},
			}}, "count": 1}))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})
	pack := &PackResource{client: client}
	var schemaResponse resource.SchemaResponse
	pack.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	plan := packPlanObjectWithItems(t, packInstallInfo{ID: types.StringValue("test-pack-metadata"), Version: types.StringValue("1.0.0")})
	planValues := plan.Attributes()
	planValues["items"] = types.ListNull(plan.AttributeTypes(ctx)["items"].(types.ListType).ElemType)
	configured := mustPackTags(t, `{"data_type":["logs"],"domain":["observability"],"streamtags":["processing"],"technology":["aws"],"ManagedByTFE":"true","product_id":"123456","used_for":"non_prod"}`)
	attributes := configured.UnderlyingValue().(types.Object).Attributes()
	attributeTypes := configured.UnderlyingValue().(types.Object).AttributeTypes(ctx)
	for _, key := range []string{"data_type", "domain", "streamtags", "technology"} {
		elements := attributes[key].(types.List).Elements()
		attributes[key] = types.TupleValueMust([]attr.Type{types.StringType}, elements)
		attributeTypes[key] = attributes[key].Type(ctx)
	}
	configured = types.DynamicValue(types.ObjectValueMust(attributeTypes, attributes))
	planValues["tags"] = configured
	plan = types.ObjectValueMust(plan.AttributeTypes(ctx), planValues)
	raw, err := plan.ToTerraformValue(ctx)
	require.NoError(t, err)
	var created resource.CreateResponse
	created.State = tfsdk.State{Schema: schemaResponse.Schema}
	pack.Create(ctx, resource.CreateRequest{
		Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: raw},
		Plan:   tfsdk.Plan{Schema: schemaResponse.Schema, Raw: raw},
	}, &created)
	require.False(t, created.Diagnostics.HasError(), "%v", created.Diagnostics)
	assertTags := func(state tfsdk.State, want types.Dynamic) {
		t.Helper()
		var tags types.Dynamic
		require.False(t, state.GetAttribute(ctx, path.Root("tags"), &tags).HasError())
		require.True(t, tags.Equal(want), "got %s, want %s", tags, want)
	}
	assertTags(created.State, configured)
	require.Equal(t, "true", remoteTags["ManagedByTFE"])
	require.Equal(t, []any{"logs"}, remoteTags["dataType"])
	require.NotContains(t, remoteTags, "custom")
	require.NotContains(t, remoteTags, "data_type")
	for _, next := range []types.Dynamic{
		mustPackTags(t, `{"ManagedByTFE":"TRUE","product_id":"001234","useCase":["analytics"],"metadata":{"OwnerName":"team","Enabled":true,"Count":9007199254740993}}`),
		types.DynamicNull(),
		mustPackTags(t, `{}`),
	} {
		planValues["tags"] = next
		plan = types.ObjectValueMust(plan.AttributeTypes(ctx), planValues)
		raw, err = plan.ToTerraformValue(ctx)
		require.NoError(t, err)
		var updated resource.UpdateResponse
		updated.State = tfsdk.State{Schema: schemaResponse.Schema}
		pack.Update(ctx, resource.UpdateRequest{
			Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: raw},
			Plan:   tfsdk.Plan{Schema: schemaResponse.Schema, Raw: raw}, State: created.State,
		}, &updated)
		require.False(t, updated.Diagnostics.HasError(), "%v", updated.Diagnostics)
		if !next.IsNull() {
			assertTags(updated.State, next)
		} else {
			require.Equal(t, "TRUE", remoteTags["ManagedByTFE"])
		}
		created.State = updated.State
	}
	require.NotContains(t, remoteTags, "ManagedByTFE")
	require.NotContains(t, remoteTags, "useCase")
	planValues["tags"] = mustPackTags(t, `{"domain":["observability"],"ManagedByTFE":"true"}`)
	plan = types.ObjectValueMust(plan.AttributeTypes(ctx), planValues)
	raw, err = plan.ToTerraformValue(ctx)
	require.NoError(t, err)
	priorState := tfsdk.State{Schema: schemaResponse.Schema, Raw: raw}
	remoteTags["ManagedByTFE"] = "FALSE"
	remoteTags["domain"] = []any{"changed"}
	remoteTags["Unmanaged"] = "remote-added-key"
	var read resource.ReadResponse
	read.State = priorState
	pack.Read(ctx, resource.ReadRequest{State: priorState}, &read)
	require.False(t, read.Diagnostics.HasError(), "%v", read.Diagnostics)
	assertTags(read.State, mustPackTags(t, `{"domain":["changed"],"ManagedByTFE":"FALSE","Unmanaged":"remote-added-key"}`))
	planValues["tags"] = types.DynamicNull()
	plan = types.ObjectValueMust(plan.AttributeTypes(ctx), planValues)
	raw, err = plan.ToTerraformValue(ctx)
	require.NoError(t, err)
	var imported resource.ImportStateResponse
	imported.State = tfsdk.State{Schema: schemaResponse.Schema, Raw: raw}
	pack.ImportState(ctx, resource.ImportStateRequest{ID: `{"group_id":"default","id":"test-pack-metadata"}`}, &imported)
	require.False(t, imported.Diagnostics.HasError(), "%v", imported.Diagnostics)
	read.State = imported.State
	pack.Read(ctx, resource.ReadRequest{State: imported.State}, &read)
	require.False(t, read.Diagnostics.HasError(), "%v", read.Diagnostics)
	var importedTags types.Dynamic
	require.False(t, read.State.GetAttribute(ctx, path.Root("tags"), &importedTags).HasError())
	importedMap, err := tagsAPIMap(ctx, importedTags)
	require.NoError(t, err)
	require.Equal(t, "remote-added-key", importedMap["Unmanaged"])
	dataSource := &PackDataSource{client: client}
	var dataSchema datasource.SchemaResponse
	dataSource.Schema(ctx, datasource.SchemaRequest{}, &dataSchema)
	dataState := tfsdk.State{Schema: dataSchema.Schema}
	require.False(t, dataState.Set(ctx, &PackDataSourceModel{GroupID: types.StringValue("default"), ID: types.StringValue("test-pack-metadata")}).HasError())
	var dataRead datasource.ReadResponse
	dataRead.State = dataState
	dataSource.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: dataSchema.Schema, Raw: dataState.Raw}}, &dataRead)
	require.False(t, dataRead.Diagnostics.HasError(), "%v", dataRead.Diagnostics)
	var data PackDataSourceModel
	require.False(t, dataRead.State.Get(ctx, &data).HasError())
	require.True(t, data.Tags.Equal(importedTags))
}

func TestPackCustomTagsSettingsReadFailures(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
		want   string
		fail   bool
	}{
		{"settings authoritative", 200, `{"items":[{"package":{"tags":{"product_id":"001234"}}}]}`, `{"product_id":"001234"}`, false},
		{"server dropped custom tags", 200, `{"items":[{"package":{"tags":{}}}]}`, `{}`, false},
		{"unsupported GET fallback", 404, `{}`, `{"InventoryKey":"inventory"}`, false},
		{"permission denied", 403, `{}`, "", true},
		{"malformed package", 200, `{"items":[{"id":"settings"}]}`, "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				if request.URL.Path == "/api/v1/m/default/p/test/pack/settings" {
					writer.WriteHeader(test.status)
					_, err := writer.Write([]byte(test.body))
					require.NoError(t, err)
					return
				}
				_, err := writer.Write([]byte(`{"items":[{"id":"test","tags":{"InventoryKey":"inventory"}}]}`))
				require.NoError(t, err)
			}))
			defer server.Close()
			pack := &PackResource{client: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}
			data := &PackResourceModel{GroupID: types.StringValue("default"), ID: types.StringValue("test")}
			err := pack.refreshPackState(context.Background(), data)
			if test.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, data.Tags.Equal(mustPackTags(t, test.want)))
			planned := mustPackTags(t, `{"PlannedKey":"must-not-be-preserved"}`)
			preserved := preservePackTags(context.Background(), data.Tags, planned)
			actual := preserved.UnderlyingValue().(types.Object).Attributes()
			require.True(t, actual["PlannedKey"].IsNull())
		})
	}
}
