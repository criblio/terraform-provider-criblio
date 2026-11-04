package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscoverSystemLimitsSingleton(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/m/test-fleet/system/limits" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{map[string]any{"minFreeSpace": "2 GB"}}, "count": 1}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()
	ctx := context.Background()
	entry, ok := mustBuildRegistry(t, ctx).ByTypeName("criblio_system_limits")
	require.True(t, ok)
	require.Equal(t, "group_id", entry.ImportIDFormat)
	require.Equal(t, "/m/{group_id}/system/limits", entry.RESTGetPath)
	ids, perGroup, err := listRESTIdentifiers(ctx, criblMockClient(server), entry, []string{"test-fleet"})
	require.NoError(t, err)
	require.Equal(t, []map[string]string{{"group_id": "test-fleet", "id": "test-fleet"}}, ids)
	require.Equal(t, map[string]int{"test-fleet": 1}, perGroup)
}
