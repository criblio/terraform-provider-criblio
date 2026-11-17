package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestListPagination(t *testing.T) {
	for _, tc := range []struct {
		name       string
		total, cap int
		metadata   bool
	}{
		{"total count and server cap", 7, 2, true},
		{"count only and server cap", 7, 2, false},
		{"empty", 0, 2, true},
		{"exact page", 100, 100, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if r.URL.Query().Get("engine") != "handlebars" || r.URL.Query().Get("limit") != "100" {
					t.Errorf("query: %s", r.URL.RawQuery)
				}
				offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
				if err != nil {
					t.Error(err)
					return
				}
				items := []map[string]string{}
				for i := offset; i < tc.total && i < offset+tc.cap; i++ {
					items = append(items, map[string]string{"id": fmt.Sprintf("item-%d", i)})
				}
				envelope := map[string]any{"items": items, "count": len(items)}
				if tc.metadata {
					envelope["offset"] = offset
					envelope["totalCount"] = tc.total
					envelope["limit"] = tc.cap
				}
				if err := json.NewEncoder(w).Encode(envelope); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			result, err := List[map[string]string](context.Background(), New(Config{BaseURL: server.URL, BearerToken: "test"}), "/collection?engine=handlebars", true)
			if err != nil {
				t.Fatal(err)
			}
			if result == nil || len(*result) != tc.total {
				t.Fatalf("result=%v, want %d items", result, tc.total)
			}
			for i, item := range *result {
				if item["id"] != fmt.Sprintf("item-%d", i) {
					t.Fatalf("item %d: %v", i, item)
				}
			}
			if requests > tc.total+1 {
				t.Fatalf("too many requests: %d", requests)
			}
		})
	}
}

func TestListResponseShapes(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		paginated  bool
		requests   int
	}{
		{"array", `[{"id":"one"}]`, false, 1},
		{"array ignores pagination", `[{"id":"one"}]`, true, 1},
		{"envelope", `{"items":[{"id":"one"}],"count":1}`, false, 1},
		{"older envelope ignores pagination", `{"items":[{"id":"one"}],"count":1}`, true, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if _, err := fmt.Fprint(w, tc.body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			result, err := List[map[string]string](context.Background(), New(Config{BaseURL: server.URL, BearerToken: "test"}), "/collection", tc.paginated)
			if err != nil {
				t.Fatal(err)
			}
			if len(*result) != 1 || requests != tc.requests {
				t.Fatalf("items=%v requests=%d", result, requests)
			}
		})
	}
}

func TestListFailsOnIncompletePagination(t *testing.T) {
	for _, body := range []string{`{"items":[{"id":"one"}],"offset":0,"totalCount":2}`, `invalid`, `{"items":[{"id":"one"}],"offset":-1}`} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, err := fmt.Fprint(w, body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			if _, err := List[map[string]string](context.Background(), New(Config{BaseURL: server.URL, BearerToken: "test"}), "/collection", true); err == nil {
				t.Fatal("expected pagination/decode error")
			}
		})
	}
}

func TestListDoesNotReturnPartialResults(t *testing.T) {
	for _, failure := range []string{"http", "empty page"} {
		t.Run(failure, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("offset") == "0" {
					if _, err := fmt.Fprint(w, `{"items":[{"id":"first"}],"offset":0,"totalCount":2}`); err != nil {
						t.Error(err)
					}
					return
				}
				if failure == "http" {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				if _, err := fmt.Fprint(w, `{"items":[],"offset":1,"totalCount":2}`); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			result, err := List[map[string]string](context.Background(), New(Config{BaseURL: server.URL, BearerToken: "test"}), "/collection", true)
			if err == nil || result != nil {
				t.Fatalf("result=%v err=%v; want an error and no partial result", result, err)
			}
		})
	}
}

func TestListBoundsVolatilePages(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		// This API ignores the offset and changes a non-identity field each time.
		if _, err := fmt.Fprintf(w, `{"items":[{"id":"same","counter":%d}]}`, requests); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	result, err := List[map[string]any](context.Background(), New(Config{BaseURL: server.URL, BearerToken: "test"}), "/collection", true)
	if err == nil || result != nil || requests != listMaxPages {
		t.Fatalf("result=%v error=%v requests=%d; want bounded failure", result, err, requests)
	}
}
