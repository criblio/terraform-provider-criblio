package discovery

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatasetACLDiscoveryIncludesSharedBuiltins(t *testing.T) {
	entry, ok := mustBuildRegistry(t, context.Background()).ByTypeName("criblio_search_dataset_acl")
	require.True(t, ok)
	for _, denied := range []bool{false, true} {
		t.Run(fmt.Sprintf("denied=%t", denied), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				require.Equal(t, http.MethodGet, req.Method)
				w.Header().Set("Content-Type", "application/json")
				body := `{"count":0,"items":[]}`
				switch req.URL.Path {
				case "/api/v1/m/default_search/search/datasets":
					body = `{"count":3,"items":[{"id":"logs"},{"id":"empty"},{"id":"cribl_edge_logs","lib":"cribl"}]}`
				case "/api/v1/m/default_search/search/datasets/logs/acl":
					body = `{"count":1,"items":[{"user":"client-id@clients","perms":[{"type":"datasets","gid":"default_search","id":"logs","policy":"DatasetRead"}]}]}`
				case "/api/v1/m/default_search/search/datasets/cribl_edge_logs/acl/teams":
					body = `{"count":1,"items":[{"team":"engineering","perms":[{"type":"datasets","gid":"default_search","id":"cribl_edge_logs","policy":"DatasetRead"}]}]}`
				default:
					require.True(t, strings.HasSuffix(req.URL.Path, "/acl") || strings.HasSuffix(req.URL.Path, "/acl/teams"))
				}
				if denied && strings.Contains(req.URL.Path, "/empty/acl") {
					w.WriteHeader(http.StatusForbidden)
					body = `{"status":"error","message":"ACL access denied"}`
				}
				if _, err := w.Write([]byte(body)); err != nil {
					t.Errorf("write fixture: %v", err)
				}
			}))
			defer server.Close()
			ids, err := ListItemIdentifiers(context.Background(), criblMockClient(server), entry, []string{"default_search"})
			if denied {
				require.ErrorContains(t, err, "ACL access denied")
				return
			}
			require.NoError(t, err)
			require.Equal(t, []map[string]string{{"dataset_id": "logs"}, {"dataset_id": "cribl_edge_logs"}}, ids)
		})
	}
}
