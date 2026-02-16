package discovery

import (
	"context"
	"net/http"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover_AllSupportedTypesListed(t *testing.T) {
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithClient(&http.Client{}))

	results, err := Discover(ctx, client, reg, nil, nil, nil)
	require.NoError(t, err)

	// Every registry entry gets a result (types without list method show count 0).
	assert.Len(t, results, reg.Len(), "Discover should return one result per registry entry")
	discoverable := 0
	for _, e := range reg.Entries() {
		if e.SDKService != "" && e.ListMethod != "" {
			discoverable++
		}
	}
	assert.GreaterOrEqual(t, discoverable, 5, "registry should have several types with list method")

	// Each result has the type name set
	typeNames := make(map[string]struct{})
	for _, r := range results {
		assert.NotEmpty(t, r.TypeName, "result should have TypeName")
		typeNames[r.TypeName] = struct{}{}
	}
	assert.Len(t, typeNames, len(results), "results should be unique by type name")
}

func TestDiscover_IncludeExcludeFilter(t *testing.T) {
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithClient(&http.Client{}))

	// Only criblio_source and criblio_pipeline
	results, err := Discover(ctx, client, reg, []string{"criblio_source", "criblio_pipeline"}, nil, nil)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	names := make(map[string]bool)
	for _, r := range results {
		names[r.TypeName] = true
	}
	assert.True(t, names["criblio_source"])
	assert.True(t, names["criblio_pipeline"])

	// Exclude criblio_source
	results, err = Discover(ctx, client, reg, nil, []string{"criblio_source"}, nil)
	require.NoError(t, err)
	for _, r := range results {
		assert.NotEqual(t, "criblio_source", r.TypeName, "excluded type should not appear")
	}
}

func TestDiscover_SDKErrorsSurfacedWithResourceContext(t *testing.T) {
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	// Client with no auth / invalid server will fail when we call List*
	client := sdk.New(sdk.WithClient(&http.Client{}))

	results, err := Discover(ctx, client, reg, []string{"criblio_source"}, nil, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	r := results[0]
	assert.Equal(t, "criblio_source", r.TypeName)
	// SDK should fail (e.g. connection or auth), and error must include resource context
	if r.Err != nil {
		assert.Contains(t, r.Err.Error(), "criblio_source", "error should include resource type name for context")
	}
}

func TestDiscover_EmptyIncludeNoDiscoverableTypes(t *testing.T) {
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithClient(&http.Client{}))

	// Include only a type that doesn't exist
	results, err := Discover(ctx, client, reg, []string{"criblio_nonexistent"}, nil, nil)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func mustBuildRegistry(t *testing.T, ctx context.Context) *registry.Registry {
	t.Helper()
	p := provider.New("test")()
	constructors := p.Resources(ctx)
	reg, err := registry.NewFromResources(ctx, constructors, registry.MetadataFromProvider(), nil)
	require.NoError(t, err)
	return reg
}
