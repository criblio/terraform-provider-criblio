package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPackLookupsReadDownloadsContent(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/m/default/packs":
			writePackLookupJSON(t, w, []map[string]any{
				{"id": "knowledge"},
			})
		case "/api/v1/m/default/p/knowledge/system/lookups/pack_lookup.csv":
			writePackLookupJSON(t, w, map[string]any{
				"items": []map[string]any{
					{
						"id":   "pack_lookup.csv",
						"mode": "memory",
					},
				},
			})
		case "/api/v1/m/default/p/knowledge/system/lookups/pack_lookup.csv/content":
			if got := r.URL.Query().Get("raw"); got != "true" {
				t.Fatalf("raw = %q, want true", got)
			}
			w.Header().Set("Content-Type", "text/csv")
			_, _ = w.Write([]byte("key,value\nalpha,beta\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	api := newPackLookupsAPI(restclient.New(restclient.Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	}))

	model, err := api.Read(context.Background(), PackLookupsModel{
		GroupID: types.StringValue("default"),
		Pack:    types.StringValue("knowledge"),
		ID:      types.StringValue("pack_lookup.csv"),
	})
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if got := model.ID.ValueString(); got != "pack_lookup.csv" {
		t.Fatalf("ID = %q, want configured ID", got)
	}
	if got := model.Content.ValueString(); got != "key,value\nalpha,beta\n" {
		t.Fatalf("Content = %q, want downloaded content", got)
	}
	wantPaths := []string{
		"/api/v1/m/default/packs",
		"/api/v1/m/default/p/knowledge/system/lookups/pack_lookup.csv",
		"/api/v1/m/default/p/knowledge/system/lookups/pack_lookup.csv/content",
	}
	if len(requestedPaths) != len(wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", requestedPaths, wantPaths)
	}
	for i := range wantPaths {
		if requestedPaths[i] != wantPaths[i] {
			t.Fatalf("requested paths = %#v, want %#v", requestedPaths, wantPaths)
		}
	}
}

func writePackLookupJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("write response: %v", err)
	}
}
