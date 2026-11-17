package restclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestResolvePackIDSharesListingsAndExpires(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.URL.Path != "/api/v1/m/group/packs" && r.URL.Path != "/api/v1/m/other/packs" {
			t.Errorf("incorrect path: %s", r.URL.Path)
		}
		if _, err := fmt.Fprint(w, `{"items":[{"id":"LegacyPack"},{"id":"another"}],"totalCount":2,"offset":0}`); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := New(Config{BaseURL: server.URL, BearerToken: "test"})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := client.ResolvePackID(context.Background(), "group", "legacypack"); got != "LegacyPack" {
				t.Errorf("resolved ID = %s", got)
			}
		}()
	}
	wg.Wait()
	if got := client.ResolvePackID(context.Background(), "group", "ANOTHER"); got != "another" {
		t.Errorf("another ID = %s", got)
	}
	if requests.Load() != 1 {
		t.Fatalf("group list requests = %d, want 1", requests.Load())
	}
	client.ResolvePackID(context.Background(), "other", "LegacyPack")
	if requests.Load() != 2 {
		t.Fatal("cache was shared across groups")
	}
	client.packIDs.mu.Lock()
	entry := client.packIDs.groups["group"]
	entry.expires = time.Now().Add(-time.Second)
	client.packIDs.groups["group"] = entry
	client.packIDs.mu.Unlock()
	client.ResolvePackID(context.Background(), "group", "LegacyPack")
	if requests.Load() != 3 {
		t.Fatal("expired lookup did not refresh")
	}
	otherClient := New(Config{BaseURL: server.URL, BearerToken: "test"})
	otherClient.ResolvePackID(context.Background(), "group", "LegacyPack")
	if requests.Load() != 4 {
		t.Fatal("cache was shared across clients")
	}
}

func TestResolvePackIDDoesNotCacheFailures(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if _, err := fmt.Fprint(w, `{"items":[{"id":"LegacyPack"}],"totalCount":1}`); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := New(Config{BaseURL: server.URL, BearerToken: "test"})
	if got := client.ResolvePackID(context.Background(), "group", "LegacyPack"); got != "legacypack" {
		t.Errorf("fallback = %s", got)
	}
	if got := client.ResolvePackID(context.Background(), "group", "LegacyPack"); got != "LegacyPack" {
		t.Errorf("recovered ID = %s", got)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
}
