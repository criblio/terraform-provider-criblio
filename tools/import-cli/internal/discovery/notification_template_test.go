package discovery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotificationTemplateDiscoverySelectsHandlebars(t *testing.T) {
	ctx := context.Background()
	entry, ok := mustBuildRegistry(t, ctx).ByTypeName("criblio_notification_template")
	require.True(t, ok)
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requests++
		if req.Method != http.MethodGet || req.URL.Path != "/api/v1/notification-templates" || req.URL.Query().Get("engine") != "handlebars" {
			t.Errorf("unexpected template discovery request: %s %s", req.Method, req.URL)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"items":[{"id":"custom-email","type":"smtp"},{"id":"default-email","type":"smtp"}],"count":2}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	ids, err := ListItemIdentifiers(ctx, criblMockClient(server), entry, []string{"default", "default_search"})
	require.NoError(t, err)
	require.ElementsMatch(t, []map[string]string{{"id": "custom-email"}, {"id": "default-email"}}, ids)
	// Templates are workspace-global, so discovering multiple groups must not
	// duplicate either requests or identifiers. Export filters built-in IDs.
	require.Equal(t, 1, requests)
}
