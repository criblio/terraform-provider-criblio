package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLookupFileCreateUploadsContentBeforeMetadata(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/m/default/system/lookups":
			if got := r.URL.Query().Get("filename"); got != "my_id.csv" {
				t.Fatalf("filename = %q, want my_id.csv", got)
			}
			if got := r.Header.Get("Content-Type"); got != "text/csv" {
				t.Fatalf("Content-Type = %q, want text/csv", got)
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if got := string(body); got != "column1,column2\nvalue1,value2\n" {
				t.Fatalf("upload body = %q", got)
			}
			writeLookupFileJSON(t, w, map[string]any{"items": []any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/m/default/system/lookups":
			writeLookupFileJSON(t, w, map[string]any{
				"items": []map[string]any{
					{
						"id":   "my_id.csv",
						"mode": "memory",
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	api := newLookupFileAPI(restclient.New(restclient.Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	}))

	model, err := api.Create(context.Background(), LookupFileModel{
		Content: types.StringValue("column1,column2\nvalue1,value2\n"),
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("my_id"),
		Mode:    types.StringValue("memory"),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got := model.ID.ValueString(); got != "my_id" {
		t.Fatalf("ID = %q, want configured ID", got)
	}
	wantCalls := []string{
		"PUT /api/v1/m/default/system/lookups?filename=my_id.csv",
		"POST /api/v1/m/default/system/lookups",
	}
	if len(calls) != len(wantCalls) {
		t.Fatalf("calls = %#v, want %#v", calls, wantCalls)
	}
	for i := range wantCalls {
		if calls[i] != wantCalls[i] {
			t.Fatalf("calls = %#v, want %#v", calls, wantCalls)
		}
	}
}

func TestLookupFileReadFallsBackToCSVExtension(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/m/default/system/lookups/my_id":
			writeLookupFileJSON(t, w, map[string]any{"items": []any{}})
		case "/api/v1/m/default/system/lookups/my_id.csv":
			writeLookupFileJSON(t, w, map[string]any{
				"items": []map[string]any{
					{
						"id":   "my_id.csv",
						"mode": "memory",
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	api := newLookupFileAPI(restclient.New(restclient.Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	}))

	model, err := api.Read(context.Background(), LookupFileModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("my_id"),
	})
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if got := model.ID.ValueString(); got != "my_id" {
		t.Fatalf("ID = %q, want configured ID", got)
	}
	if len(requestedPaths) != 2 {
		t.Fatalf("requested paths = %#v, want raw ID then .csv fallback", requestedPaths)
	}
}

func TestLookupFileAPIIDsDoesNotFallbackForKnownExtensions(t *testing.T) {
	for _, id := range []string{"lookup.csv", "lookup.gz", "lookup.csv.gz", "GeoIP.mmdb"} {
		ids := lookupFileAPIIDs(id)
		if len(ids) != 1 || ids[0] != id {
			t.Fatalf("lookupFileAPIIDs(%q) = %#v, want only original ID", id, ids)
		}
	}
}

func writeLookupFileJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("write response: %v", err)
	}
}
