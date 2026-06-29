package converter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshFromMethodName(t *testing.T) {
	tests := []struct {
		getMethod string
		want      string
	}{
		{"GetInputByID", "RefreshFromOperationsGetInputByIDResponseBody"},
		{"GetPipelineByID", "RefreshFromOperationsGetPipelineByIDResponseBody"},
		{"GetRoutesByGroupID", "RefreshFromOperationsGetRoutesByGroupIDResponseBody"},
		{"GetSavedJobByID", "RefreshFromOperationsGetSavedJobByIDResponseBody"},
	}
	for _, tt := range tests {
		t.Run(tt.getMethod, func(t *testing.T) {
			got := RefreshFromMethodName(tt.getMethod)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertUsesRESTGetPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/m/default/pipelines/p1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"id":"p1","conf":{}}]}`))
	}))
	defer server.Close()

	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_pipeline")
	require.True(t, ok)
	require.NotEmpty(t, e.RESTGetPath)

	client := &importclient.Client{
		REST: restclient.New(restclient.Config{
			BaseURL:     server.URL,
			BearerToken: "mock",
			HTTPClient:  server.Client(),
		}),
	}
	model, err := Convert(context.Background(), client, e, map[string]string{
		"GroupID": "default",
		"ID":      "p1",
	})
	require.NoError(t, err)
	pipeline, ok := model.(*provider.PipelineResourceModel)
	require.True(t, ok)
	assert.Equal(t, "p1", pipeline.ID.ValueString())
	assert.Equal(t, "default", pipeline.GroupID.ValueString())
}

func TestConvertFromResponseBody_source(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok, "registry must contain criblio_source")

	body := &struct {
		Items []map[string]any
	}{
		Items: []map[string]any{
			{
				"id":   "test",
				"type": "cribl_http",
				"host": "0.0.0.0",
				"port": 10200,
			},
		},
	}
	model, err := ConvertFromResponseBody(ctx, e, body)
	require.NoError(t, err)
	require.NotNil(t, model)
	// Model should be *SourceResourceModel
	_, ok = model.(*provider.SourceResourceModel)
	assert.True(t, ok, "model should be *SourceResourceModel")
}

func TestConvertFromResponseBody_pipeline(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_pipeline")
	require.True(t, ok, "registry must contain criblio_pipeline")

	body := &struct {
		Items []map[string]any
	}{}
	_, err := ConvertFromResponseBody(ctx, e, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "response body has no items")
}

func TestConvertFromResponseBody_unknownModelType(t *testing.T) {
	ctx := context.Background()
	e := registry.Entry{
		TypeName:       "criblio_fake",
		ModelTypeName:  "FakeResourceModel",
		SDKService:     "Inputs",
		GetMethod:      "GetInputByID",
		ImportIDFormat: "id",
	}
	body := &struct {
		Items []map[string]any
	}{}
	_, err := ConvertFromResponseBody(ctx, e, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown model type")
	assert.Contains(t, err.Error(), "FakeResourceModel")
}

func TestConvertFromResponseBody_noGetMethod(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, _ := reg.ByTypeName("criblio_source")
	e.GetMethod = ""

	_, err := ConvertFromResponseBody(ctx, e, &struct {
		Items []map[string]any
	}{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no GetMethod")
}

func TestResourceModelTypes_hasEntriesForSupportedTypes(t *testing.T) {
	types := ResourceModelTypes()
	require.NotEmpty(t, types)
	assert.Contains(t, types, "SourceResourceModel")
	assert.Contains(t, types, "PipelineResourceModel")
	assert.Contains(t, types, "RoutesResourceModel")
}

func TestGeneratedModelTypes_hasOneOfMigrationEntries(t *testing.T) {
	types := GeneratedModelTypes()
	require.NotEmpty(t, types)
	assert.Contains(t, types, "NotificationTargetResourceModel")
	assert.Contains(t, types, "SearchDatasetProviderResourceModel")
}

func TestConvertGeneratedNotificationTargetPopulatesOneOfBlock(t *testing.T) {
	e := registry.Entry{
		TypeName:      "criblio_notification_target",
		ModelTypeName: "NotificationTargetResourceModel",
		GetMethod:     "GetNotificationTargetByID",
	}
	responseBody := struct {
		Items []map[string]any
	}{
		Items: []map[string]any{
			{
				"id":   "slack-1",
				"type": "slack",
			},
		},
	}

	model, err := convertGeneratedModelFromResponseBody(e, reflect.TypeOf((*provider.NotificationTargetResourceModel)(nil)).Elem(), responseBody)
	require.NoError(t, err)

	target, ok := model.(*provider.NotificationTargetResourceModel)
	require.True(t, ok)
	require.NotNil(t, target.SlackTarget)
	assert.Equal(t, "slack-1", target.SlackTarget.ID.ValueString())
	assert.True(t, target.SlackTarget.URL.IsNull())
}

func TestObjectListValue(t *testing.T) {
	value, err := objectListValue(json.RawMessage(`[{"name":"val","type":"number"}]`))
	require.NoError(t, err)
	require.False(t, value.IsNull())
	require.False(t, value.IsUnknown())

	var args []provider.GlobalVarArgsModel
	diags := value.ElementsAs(context.Background(), &args, false)
	require.False(t, diags.HasError(), diags)
	require.Len(t, args, 1)
	assert.Equal(t, "val", args[0].Name.ValueString())
	assert.Equal(t, "number", args[0].Type.ValueString())
	assert.Equal(t, types.ObjectType{AttrTypes: provider.GlobalVarArgsAttrTypes()}, value.ElementType(context.Background()))
}

func TestObjectListValueNormalizesAPIKeysToTerraformNames(t *testing.T) {
	value, err := objectListValue(json.RawMessage(`[{
		"name": "test",
		"eventBreakerRegex": "/[\\n\\r]+/",
		"maxEventBytes": 51200,
		"parserEnabled": false,
		"shouldUseDataRaw": false,
		"timestampAnchorRegex": "/^/",
		"timestamp": {
			"type": "auto",
			"length": 150
		}
	}]`))
	require.NoError(t, err)

	elementType, ok := value.ElementType(context.Background()).(types.ObjectType)
	require.True(t, ok)
	assert.Contains(t, elementType.AttrTypes, "event_breaker_regex")
	assert.Contains(t, elementType.AttrTypes, "max_event_bytes")
	assert.Contains(t, elementType.AttrTypes, "parser_enabled")
	assert.Contains(t, elementType.AttrTypes, "should_use_data_raw")
	assert.Contains(t, elementType.AttrTypes, "timestamp_anchor_regex")
	assert.NotContains(t, elementType.AttrTypes, "eventBreakerRegex")
	assert.NotContains(t, elementType.AttrTypes, "timestampAnchorRegex")

	first := value.Elements()[0].(types.Object)
	attrs := first.Attributes()
	assert.Contains(t, attrs, "timestamp_anchor_regex")
	timestamp := attrs["timestamp"].(types.Object)
	assert.Contains(t, timestamp.Attributes(), "type")
	assert.Contains(t, timestamp.Attributes(), "length")
}

func TestConvertGeneratedEventBreakerRulesetKeepsTypedEmptyRules(t *testing.T) {
	e := registry.Entry{
		TypeName:      "criblio_event_breaker_ruleset",
		ModelTypeName: "EventBreakerRulesetResourceModel",
		GetMethod:     "GetEventBreakerRulesetByID",
	}
	responseBody := struct {
		Items []map[string]any
	}{
		Items: []map[string]any{
			{
				"id":      "empty-rules",
				"groupId": "default",
				"rules":   []any{},
			},
		},
	}

	model, err := convertGeneratedModelFromResponseBody(e, reflect.TypeOf((*provider.EventBreakerRulesetResourceModel)(nil)).Elem(), responseBody)
	require.NoError(t, err)

	eventBreaker, ok := model.(*provider.EventBreakerRulesetResourceModel)
	require.True(t, ok)
	require.False(t, eventBreaker.Rules.IsNull())
	require.False(t, eventBreaker.Rules.IsUnknown())
	assert.Equal(t, types.ObjectType{AttrTypes: provider.EventBreakerRulesetRulesAttrTypes()}, eventBreaker.Rules.ElementType(context.Background()))
	assert.Empty(t, eventBreaker.Rules.Elements())
}

// TestConvertFromResponseBody_destination verifies the correct RefreshFrom* method is invoked
// for another resource type (GetOutputByID -> RefreshFromOperationsGetOutputByIDResponseBody).
func TestConvertFromResponseBody_destination(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_destination")
	require.True(t, ok, "registry must contain criblio_destination")

	body := &struct {
		Items []map[string]any
	}{
		Items: []map[string]any{
			{
				"id":   "test",
				"type": "devnull",
			},
		},
	}
	model, err := ConvertFromResponseBody(ctx, e, body)
	require.NoError(t, err)
	require.NotNil(t, model)
	_, ok = model.(*provider.DestinationResourceModel)
	assert.True(t, ok, "model should be *DestinationResourceModel")
}

func TestConvertFromResponseBody_certificateGeneratedModel(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_certificate")
	require.True(t, ok, "registry must contain criblio_certificate")

	body := &struct {
		Items []struct {
			ID          string   `json:"id,omitempty"`
			Cert        string   `json:"cert,omitempty"`
			Description string   `json:"description,omitempty"`
			InUse       []string `json:"inUse,omitempty"`
		} `json:"items"`
	}{
		Items: []struct {
			ID          string   `json:"id,omitempty"`
			Cert        string   `json:"cert,omitempty"`
			Description string   `json:"description,omitempty"`
			InUse       []string `json:"inUse,omitempty"`
		}{
			{
				ID:          "my-cert",
				Cert:        "cert-body",
				Description: "generated certificate",
				InUse:       []string{},
			},
		},
	}
	model, err := ConvertFromResponseBodyWithIdentifiers(ctx, e, body, map[string]string{
		"GroupID": "default",
		"ID":      "my-cert",
	})
	require.NoError(t, err)
	require.NotNil(t, model)

	cert, ok := model.(*provider.CertificateResourceModel)
	require.True(t, ok, "model should be *CertificateResourceModel")
	assert.Equal(t, "default", cert.GroupID.ValueString())
	assert.Equal(t, "my-cert", cert.ID.ValueString())
	assert.Equal(t, "generated certificate", cert.Description.ValueString())
	assert.Empty(t, cert.InUse)
}

func TestConvertFromResponseBody_schemaGeneratedModelNormalizedJSON(t *testing.T) {
	ctx := context.Background()
	e := registry.Entry{
		TypeName:       "criblio_schema",
		ModelTypeName:  "SchemaResourceModel",
		GetMethod:      "GetLibSchemasByID",
		ImportIDFormat: "json:group_id,id",
	}
	body := &struct {
		Items []struct {
			Description string         `json:"description,omitempty"`
			ID          string         `json:"id,omitempty"`
			Schema      map[string]any `json:"schema,omitempty"`
		} `json:"items"`
	}{
		Items: []struct {
			Description string         `json:"description,omitempty"`
			ID          string         `json:"id,omitempty"`
			Schema      map[string]any `json:"schema,omitempty"`
		}{
			{
				Description: "schema from API",
				ID:          "my-schema",
				Schema: map[string]any{
					"type":       "object",
					"properties": map[string]any{"message": map[string]any{"type": "string"}},
				},
			},
		},
	}

	model, err := ConvertFromResponseBodyWithIdentifiers(ctx, e, body, map[string]string{
		"GroupID": "default",
		"ID":      "my-schema",
	})
	require.NoError(t, err)

	schemaModel, ok := model.(*provider.SchemaResourceModel)
	require.True(t, ok, "model should be *SchemaResourceModel")
	assert.Equal(t, "default", schemaModel.GroupID.ValueString())
	assert.Equal(t, "my-schema", schemaModel.ID.ValueString())
	require.False(t, schemaModel.Schema.IsNull())
	assert.Contains(t, schemaModel.Schema.ValueString(), `"type":"object"`)
}

func TestConvertFromResponseBodyWithIdentifiers_injects_required_fields(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok)

	body := &struct {
		Items []map[string]any
	}{
		Items: []map[string]any{
			{
				"id":   "test",
				"type": "cribl_http",
				"host": "0.0.0.0",
				"port": 10200,
			},
		},
	}
	identifiers := map[string]string{"GroupID": "default", "ID": "input-hec-1"}
	model, err := ConvertFromResponseBodyWithIdentifiers(ctx, e, body, identifiers)
	require.NoError(t, err)
	require.NotNil(t, model)
	src, ok := model.(*provider.SourceResourceModel)
	require.True(t, ok)
	assert.Equal(t, "default", src.GroupID.ValueString(), "group_id must be set for valid Terraform model")
	assert.Equal(t, "input-hec-1", src.ID.ValueString(), "id must be set for valid Terraform model")
}

func TestConvertFromResponseBodyWithIdentifiers_nil_identifiers_ok(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, _ := reg.ByTypeName("criblio_source")
	body := &struct {
		Items []map[string]any
	}{
		Items: []map[string]any{
			{
				"id":   "test",
				"type": "cribl_http",
				"host": "0.0.0.0",
				"port": 10200,
			},
		},
	}
	model, err := ConvertFromResponseBodyWithIdentifiers(ctx, e, body, nil)
	require.NoError(t, err)
	require.NotNil(t, model)
	// No panic; identifiers simply not set
}

func buildTestRegistry(t *testing.T) *registry.Registry {
	t.Helper()
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)
	reg, err := registry.NewFromResources(ctx, constructors, registry.MetadataFromProvider(), nil, OneOfBlockNamesFromModel)
	require.NoError(t, err)
	return reg
}
