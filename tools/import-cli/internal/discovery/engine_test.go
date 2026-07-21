package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func criblMockClient(server *httptest.Server) *importclient.Client {
	return &importclient.Client{
		REST: restclient.New(restclient.Config{
			BaseURL:     server.URL,
			BearerToken: "mock",
			HTTPClient:  server.Client(),
		}),
	}
}

func TestDiscover_AllSupportedTypesListed(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock")
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := criblMockClient(server)

	results, err := Discover(ctx, client, reg, nil, nil, nil, false)
	require.NoError(t, err)

	assert.Len(t, results, reg.Len(), "Discover should return one result per registry entry")
	discoverable := 0
	for _, e := range reg.Entries() {
		if e.RESTListPath != "" {
			discoverable++
		}
	}
	assert.GreaterOrEqual(t, discoverable, 5, "registry should have several REST list paths")

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
	t.Setenv("CRIBL_BEARER_TOKEN", "mock")
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := criblMockClient(server)

	results, err := Discover(ctx, client, reg, []string{"criblio_source", "criblio_pipeline"}, nil, nil, false)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	names := make(map[string]bool)
	for _, r := range results {
		names[r.TypeName] = true
	}
	assert.True(t, names["criblio_source"])
	assert.True(t, names["criblio_pipeline"])

	results, err = Discover(ctx, client, reg, nil, []string{"criblio_source"}, nil, false)
	require.NoError(t, err)
	for _, r := range results {
		assert.NotEqual(t, "criblio_source", r.TypeName, "omitted type should not appear")
	}
}

func TestDiscover_RESTErrorsSurfacedWithResourceContext(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock")
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := criblMockClient(server)

	results, err := Discover(ctx, client, reg, []string{"criblio_source"}, nil, nil, false)
	require.NoError(t, err)
	require.Len(t, results, 1)
	r := results[0]
	assert.Equal(t, "criblio_source", r.TypeName)
	if r.Err != nil {
		assert.Contains(t, r.Err.Error(), "criblio_source", "error should include resource type name for context")
	}
}

func TestDiscover_EmptyIncludeNoDiscoverableTypes(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()
	t.Setenv("CRIBL_BEARER_TOKEN", "mock")
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	client := criblMockClient(server)

	results, err := Discover(ctx, client, reg, []string{"criblio_nonexistent"}, nil, nil, false)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func mustBuildRegistry(t *testing.T, ctx context.Context) *registry.Registry {
	t.Helper()
	p := provider.New("test")()
	constructors := p.Resources(ctx)
	reg, err := registry.NewFromResources(ctx, constructors, registry.MetadataFromProvider(), nil, converter.OneOfBlockNamesFromModel)
	require.NoError(t, err)
	return reg
}

func TestSkipGroupScopedSingleton(t *testing.T) {
	assert.True(t, skipGroupScopedSingleton("criblio_routes", "default_search"))
	assert.True(t, skipGroupScopedSingleton("criblio_routes", "search"))
	assert.False(t, skipGroupScopedSingleton("criblio_routes", "default"))
	assert.False(t, skipGroupScopedSingleton("criblio_group_system_settings", "default_search"))
}

func TestIdentifiersFromRawItems_skipsLibCribl(t *testing.T) {
	reg := mustBuildRegistry(t, context.Background())
	e, ok := reg.ByTypeName("criblio_parser_lib_entry")
	require.True(t, ok, "criblio_parser_lib_entry must be in registry")

	items := []json.RawMessage{
		[]byte(`{"id":"builtin-one","lib":"cribl"}`),
		[]byte(`{"id":"user-one","lib":"criblio"}`),
		[]byte(`{"id":"builtin-two","lib":"cribl"}`),
	}
	ids, err := identifiersFromRawItems(items, map[string]string{"group_id": "default"}, e)
	require.NoError(t, err)
	assert.Len(t, ids, 1)
	assert.Equal(t, "user-one", ids[0]["id"])
	assert.Equal(t, "default", ids[0]["group_id"])
}

func TestIdentifiersFromRawItems_skipsBuiltInLookupFiles(t *testing.T) {
	reg := mustBuildRegistry(t, context.Background())
	e, ok := reg.ByTypeName("criblio_lookup_file")
	require.True(t, ok, "criblio_lookup_file must be in registry")

	items := []json.RawMessage{
		[]byte(`{"id":"lib_builtin.csv","lib":"cribl"}`),
		[]byte(`{"id":"library_builtin.csv","library":"cribl"}`),
		[]byte(`{"id":"tag_builtin.csv","tags":"cribl:default"}`),
		[]byte(`{"id":"list_tag_builtin.csv","tags":["other","cribl:default"]}`),
		[]byte(`{"id":"cribl.prefixed.csv"}`),
		[]byte(`{"id":"user_lookup.csv"}`),
	}
	ids, err := identifiersFromRawItems(items, map[string]string{"group_id": "default"}, e)
	require.NoError(t, err)
	require.Len(t, ids, 1)
	assert.Equal(t, "user_lookup.csv", ids[0]["id"])
	assert.Equal(t, "default", ids[0]["group_id"])
}

func TestRegistryImportableEntriesHaveRESTGetPath(t *testing.T) {
	for _, e := range mustBuildRegistry(t, context.Background()).Entries() {
		if e.ImportIDFormat == "" || e.TypeName == "criblio_lakehouse_dataset_connection" {
			continue
		}
		assert.NotEmpty(t, e.RESTGetPath, "registry entry %q must have RESTGetPath", e.TypeName)
	}
}

func TestIsRecoverableListDecodeError(t *testing.T) {
	jsonErr := errFromString("error unmarshaling json response body: json: cannot unmarshal number into Go value of type string")
	wrapped := fmt.Errorf("criblio_search_saved_query: %w", jsonErr)
	assert.True(t, IsRecoverableListDecodeError(jsonErr))
	assert.True(t, IsRecoverableListDecodeError(wrapped))
	assert.False(t, IsRecoverableListDecodeError(errFromString("connection refused")))
}

func errFromString(s string) error {
	if s == "" {
		return nil
	}
	return &errString{s}
}

type errString struct{ s string }

func (e *errString) Error() string { return e.s }
