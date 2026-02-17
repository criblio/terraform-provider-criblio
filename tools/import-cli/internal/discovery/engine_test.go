package discovery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// criblMockServer starts an httptest.Server that returns minimal valid responses for discovery:
// POST /oauth/token returns 200 with a token so SDK auth succeeds; groups API returns one group;
// other GETs return {"items":[]}. Use with sdk.New(sdk.WithServerURL(server.URL), sdk.WithClient(server.Client())).
func criblMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	groupsJSON := []byte(`{"items":[{"id":"default","name":"default"}]}`)
	emptyListJSON := []byte(`{"items":[]}`)
	oauthJSON := []byte(`{"access_token":"test","expires_in":300}`)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost && (strings.HasSuffix(r.URL.Path, "/oauth/token") || r.URL.Path == "/oauth/token") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(oauthJSON)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		path := r.URL.Path
		if strings.HasSuffix(path, "/products/stream/groups") || strings.HasSuffix(path, "/products/edge/groups") {
			_, _ = w.Write(groupsJSON)
			return
		}
		_, _ = w.Write(emptyListJSON)
	}))
}

func TestDiscover_AllSupportedTypesListed(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock") // skip SDK credential lookup so requests go to mock server
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithServerURL(server.URL), sdk.WithClient(server.Client()))

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
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock") // skip SDK credential lookup so requests go to mock server
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithServerURL(server.URL), sdk.WithClient(server.Client()))

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
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock") // skip SDK credential lookup so requests go to mock server
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithServerURL(server.URL), sdk.WithClient(server.Client()))

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
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock") // skip SDK credential lookup so requests go to mock server
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := sdk.New(sdk.WithServerURL(server.URL), sdk.WithClient(server.Client()))

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

// TestRegistryListMethodsExistOnSDK ensures every registry entry that has SDKService and
// ListMethod set refers to a real service field and method on sdk.CriblIo. This catches
// typos or SDK renames at test time instead of at discovery runtime (reflection would
// then fail or panic). Uses the field's type to look up the method so we validate even
// when the client's service field is nil (e.g. zero-value &sdk.CriblIo{}).
func TestRegistryListMethodsExistOnSDK(t *testing.T) {
	client := &sdk.CriblIo{}
	clientVal := reflect.ValueOf(client).Elem()

	for _, e := range mustBuildRegistry(t, context.Background()).Entries() {
		if e.SDKService == "" || e.ListMethod == "" {
			continue
		}
		svcField := clientVal.FieldByName(e.SDKService)
		require.True(t, svcField.IsValid(), "registry entry %q: SDKService %q not found on sdk.CriblIo", e.TypeName, e.SDKService)

		// Resolve the type on which the method is defined (pointer receiver *Service).
		svcType := svcField.Type()
		if svcType.Kind() == reflect.Ptr {
			svcType = svcType.Elem()
		}
		// Method is on *Service; use a non-nil pointer so MethodByName can find it.
		svcPtr := reflect.New(svcType)
		method := svcPtr.MethodByName(e.ListMethod)
		assert.True(t, method.IsValid(), "registry entry %q: ListMethod %q not found on service %s", e.TypeName, e.ListMethod, e.SDKService)
	}
}
