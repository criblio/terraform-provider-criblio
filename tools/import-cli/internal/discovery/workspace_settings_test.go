package discovery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceSettingsDiscoveryIsWorkspaceScoped(t *testing.T) {
	ctx := context.Background()
	entry, ok := mustBuildRegistry(t, ctx).ByTypeName("criblio_workspace_system_settings")
	require.True(t, ok)
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requests++
		require.Equal(t, "/api/v1/system/settings/conf", req.URL.Path)
		require.Equal(t, http.MethodGet, req.Method)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"count":1,"items":[{"upgradeSettings":{"disableAutomaticUpgrade":true}}]}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	ids, err := ListItemIdentifiers(ctx, criblMockClient(server), entry, []string{"default", "default_search", "other-fleet"})
	require.NoError(t, err)
	require.Equal(t, []map[string]string{{"id": "conf"}}, ids)
	require.Equal(t, 1, requests)
	ids, err = ListItemIdentifiers(ctx, criblMockClient(server), entry, nil)
	require.NoError(t, err)
	require.Equal(t, []map[string]string{{"id": "conf"}}, ids)
	require.Equal(t, 2, requests)
}
