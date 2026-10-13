package discovery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLocalSearchWorkersUnavailable(t *testing.T) {
	missingWorkers := &restclient.HTTPError{StatusCode: http.StatusInternalServerError, Body: `{"status":"error","message":"No workers registered"}`}
	for _, test := range []struct {
		name string
		path string
		err  error
		want bool
	}{
		{"sources", "/m/default_search/search/local_search/sources", missingWorkers, true},
		{"engines", "/m/default_search/search/local_search/engines", missingWorkers, true},
		{"future resource", "/m/default_search/search/local_search/new-resource", missingWorkers, true},
		{"wrapped error", "/m/default_search/search/local_search/sources", fmt.Errorf("list: %w", missingWorkers), true},
		{"case and whitespace", "/m/default_search/search/local_search/sources", &restclient.HTTPError{StatusCode: 500, Body: `{"message":" no WORKERS registered "}`}, true},
		{"other endpoint", "/m/default_search/search/datasets", missingWorkers, false},
		{"query string only", "/search/datasets?path=/search/local_search/sources", missingWorkers, false},
		{"similar endpoint", "/m/default_search/search/local_search_other/sources", missingWorkers, false},
		{"different failure", "/m/default_search/search/local_search/sources", &restclient.HTTPError{StatusCode: 500, Body: `{"message":"internal error"}`}, false},
		{"unauthorized", "/m/default_search/search/local_search/sources", &restclient.HTTPError{StatusCode: 401, Body: missingWorkers.Body}, false},
		{"malformed response", "/m/default_search/search/local_search/sources", &restclient.HTTPError{StatusCode: 500, Body: "No workers registered"}, false},
		{"cancelled", "/m/default_search/search/local_search/sources", context.Canceled, false},
		{"success", "/m/default_search/search/local_search/sources", nil, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isLocalSearchWorkersUnavailable(test.path, test.err))
		})
	}
}

func TestLocalSearchListIgnoresMissingWorkersOnly(t *testing.T) {
	for _, endpoint := range []string{"sources", "engines", "dataset-rulesets", "datatype-rulesets"} {
		t.Run(endpoint, func(t *testing.T) {
			for _, message := range []string{"No workers registered", "internal error"} {
				t.Run(message, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
						writer.Header().Set("Content-Type", "application/json")
						writer.WriteHeader(http.StatusInternalServerError)
						if _, err := fmt.Fprintf(writer, `{"status":"error","message":%q}`, message); err != nil {
							t.Error(err)
						}
					}))
					defer server.Close()
					ids, err := listRESTPathIdentifiers(context.Background(), criblMockClient(server), registry.Entry{
						RESTListPath: "/m/default_search/search/local_search/" + endpoint,
					}, nil)
					if message == "No workers registered" {
						require.NoError(t, err)
						assert.Empty(t, ids)
					} else {
						var httpErr *restclient.HTTPError
						require.True(t, errors.As(err, &httpErr))
						assert.Equal(t, http.StatusInternalServerError, httpErr.StatusCode)
					}
				})
			}
		})
	}
}

func TestDiscoverContinuesWithoutLocalSearchWorkers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if strings.Contains(request.URL.Path, "/search/local_search/") {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err := fmt.Fprint(writer, `{"status":"error","message":"No workers registered"}`); err != nil {
				t.Error(err)
			}
			return
		}
		response := `{"items":[]}`
		if strings.HasSuffix(request.URL.Path, "/search/dashboards") {
			response = `{"items":[{"id":"existing-dashboard"}]}`
		}
		if _, err := fmt.Fprint(writer, response); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	ctx := context.Background()
	reg := mustBuildRegistry(t, ctx)
	results, err := Discover(ctx, criblMockClient(server), reg,
		[]string{"criblio_search_source", "criblio_search_engine", "criblio_search_dashboard"}, nil, nil, false)
	require.NoError(t, err)
	require.Len(t, results, 3)
	for _, result := range results {
		require.NoError(t, result.Err)
		assert.True(t, result.InventoryComplete)
		if result.TypeName == "criblio_search_dashboard" {
			assert.Equal(t, 1, result.Count)
		} else {
			assert.Zero(t, result.Count)
			assert.Empty(t, result.Identifiers)
		}
	}
}
