package converter

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
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

func TestConvertFromResponseBody_source(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok, "registry must contain criblio_source")

	// RefreshFrom* requires at least one item; use minimal InputUnion1 (Cribl HTTP).
	body := &operations.GetInputByIDResponseBody{
		Items: []shared.InputUnion1{
			shared.CreateInputUnion1InputCriblHTTP(shared.InputCriblHTTP{
				ID:   stringPtr("test"),
				Type: shared.InputCriblHTTPTypeCriblHTTP,
				Host: "0.0.0.0",
				Port: 10200,
			}),
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

	// Pipeline RefreshFrom* expects at least one item; empty items produce a diagnostic error.
	body := &operations.GetPipelineByIDResponseBody{Items: nil}
	_, err := ConvertFromResponseBody(ctx, e, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "criblio_pipeline")
	assert.Contains(t, err.Error(), "Unexpected response from API")
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
	body := &operations.GetInputByIDResponseBody{}
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

	_, err := ConvertFromResponseBody(ctx, e, &operations.GetInputByIDResponseBody{})
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

	// RefreshFrom* requires at least one item; use minimal Output.
	body := &operations.GetOutputByIDResponseBody{
		Items: []shared.Output{
			shared.CreateOutputOutputDevnull(shared.OutputDevnull{
				ID:   stringPtr("test"),
				Type: shared.OutputDevnullTypeDevnull,
			}),
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

	body := &operations.GetCertificateByIDResponseBody{
		Items: []shared.Certificate{
			{
				ID:          "my-cert",
				Cert:        "cert-body",
				Description: stringPtr("generated certificate"),
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

	body := &operations.GetInputByIDResponseBody{
		Items: []shared.InputUnion1{
			shared.CreateInputUnion1InputCriblHTTP(shared.InputCriblHTTP{
				ID:   stringPtr("test"),
				Type: shared.InputCriblHTTPTypeCriblHTTP,
				Host: "0.0.0.0",
				Port: 10200,
			}),
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
	body := &operations.GetInputByIDResponseBody{
		Items: []shared.InputUnion1{
			shared.CreateInputUnion1InputCriblHTTP(shared.InputCriblHTTP{
				ID:   stringPtr("test"),
				Type: shared.InputCriblHTTPTypeCriblHTTP,
				Host: "0.0.0.0",
				Port: 10200,
			}),
		},
	}
	model, err := ConvertFromResponseBodyWithIdentifiers(ctx, e, body, nil)
	require.NoError(t, err)
	require.NotNil(t, model)
	// No panic; identifiers simply not set
}

func stringPtr(s string) *string { return &s }

func buildTestRegistry(t *testing.T) *registry.Registry {
	t.Helper()
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)
	reg, err := registry.NewFromResources(ctx, constructors, registry.MetadataFromProvider(), nil, OneOfBlockNamesFromModel)
	require.NoError(t, err)
	return reg
}
