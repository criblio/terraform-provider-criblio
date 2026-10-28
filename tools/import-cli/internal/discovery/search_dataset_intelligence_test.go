package discovery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchDatasetIntelligenceDiscovery(t *testing.T) {
	ctx := context.Background()
	entry, ok := mustBuildRegistry(t, ctx).ByTypeName("criblio_search_dataset_intelligence")
	require.True(t, ok)
	for _, denied := range []bool{false, true} {
		t.Run(map[bool]string{false: "existing_only", true: "access_error"}[denied], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet {
					t.Errorf("discovery must not mutate intelligence: %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				var body string
				switch req.URL.Path {
				case "/api/v1/m/default_search/search/datasets":
					body = `{"items":[{"id":"logs"},{"id":"disabled"},{"id":"cribl_internal_logs","tags":["cribl:default"]}],"count":3}`
				case "/api/v1/ai/settings/dataset-intelligence/disabled":
					w.WriteHeader(http.StatusNotFound)
					return
				case "/api/v1/ai/settings/dataset-intelligence/logs":
					if denied {
						w.WriteHeader(http.StatusForbidden)
						return
					}
					body = `{"datasetId":"logs","status":"processing","generatedAt":100}`
				case "/api/v1/ai/settings/dataset-intelligence/cribl_internal_logs":
					body = `{"datasetId":"cribl_internal_logs","status":"failed","generatedAt":100}`
				default:
					t.Errorf("unexpected discovery route: %s", req.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(body)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()
			ids, err := ListItemIdentifiers(ctx, criblMockClient(server), entry, []string{"default_search"})
			if denied {
				if err == nil || !strings.Contains(err.Error(), "403") {
					t.Fatalf("access error was swallowed: %v", err)
				}
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, []map[string]string{{"dataset_id": "logs"}, {"dataset_id": "cribl_internal_logs"}}, ids)
		})
	}
}
