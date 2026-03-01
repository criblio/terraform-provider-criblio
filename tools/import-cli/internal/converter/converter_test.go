package converter

import (
	"context"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshFromMethodName(t *testing.T) {
	tests := []struct {
		getMethod string
		want     string
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

	// RefreshFrom* requires at least one item; use minimal Input.
	typ := shared.InputCriblHTTPTypeCriblHTTP
	body := &operations.GetInputByIDResponseBody{
		Items: []shared.Input{{
			InputCriblHTTP: &shared.InputCriblHTTP{
				ID:   stringPtr("test"),
				Type: &typ,
			},
		}},
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

// TestConvertFromResponseBody_destination verifies the correct RefreshFrom* method is invoked
// for another resource type (GetOutputByID -> RefreshFromOperationsGetOutputByIDResponseBody).
func TestConvertFromResponseBody_destination(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_destination")
	require.True(t, ok, "registry must contain criblio_destination")

	// RefreshFrom* requires at least one item; use minimal Output.
	body := &operations.GetOutputByIDResponseBody{
		Items: []shared.Output{{
			OutputDevnull: &shared.OutputDevnull{
				ID:   "test",
				Type: shared.OutputDevnullTypeDevnull,
			},
		}},
	}
	model, err := ConvertFromResponseBody(ctx, e, body)
	require.NoError(t, err)
	require.NotNil(t, model)
	_, ok = model.(*provider.DestinationResourceModel)
	assert.True(t, ok, "model should be *DestinationResourceModel")
}

func TestConvertFromResponseBodyWithIdentifiers_injects_required_fields(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok)

	typ := shared.InputCriblHTTPTypeCriblHTTP
	body := &operations.GetInputByIDResponseBody{
		Items: []shared.Input{{
			InputCriblHTTP: &shared.InputCriblHTTP{
				ID:   stringPtr("test"),
				Type: &typ,
			},
		}},
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
	typ := shared.InputCriblHTTPTypeCriblHTTP
	body := &operations.GetInputByIDResponseBody{
		Items: []shared.Input{{
			InputCriblHTTP: &shared.InputCriblHTTP{
				ID:   stringPtr("test"),
				Type: &typ,
			},
		}},
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
