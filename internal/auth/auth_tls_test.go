package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestCustomHTTPClientTokenCacheIsolation(t *testing.T) {
	for _, mode := range []string{"onprem", "cloud"} {
		t.Run(mode, func(t *testing.T) {
			ClearTokenCache()
			t.Cleanup(ClearTokenCache)
			var requests atomic.Int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if err := json.NewEncoder(w).Encode(map[string]any{
					"token": "test-token", "access_token": "test-token", "expires_in": 3600,
				}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			config := CriblConfig{HTTPClient: server.Client()}
			if mode == "onprem" {
				config.OnpremServerURL = server.URL
				config.OnpremUsername = "test-user"
				config.OnpremPassword = "test-password"
			} else {
				config.CloudDomain = server.URL
				config.ClientID = "test-client"
				config.ClientSecret = "test-secret"
			}
			ctx := context.Background()
			for range 2 {
				if token, err := GetToken(ctx, &config); err != nil || token != "test-token" {
					t.Fatalf("GetToken = %q, %v", token, err)
				}
			}
			if requests.Load() != 1 {
				t.Fatalf("requests = %d, want 1 cached token request", requests.Load())
			}
			for _, client := range []*http.Client{nil, {}} {
				strict := config
				strict.HTTPClient = client
				if _, err := GetToken(ctx, &strict); err == nil || !strings.Contains(err.Error(), "certificate signed by unknown authority") {
					t.Fatalf("strict client reused cached authentication: %v", err)
				}
			}
			if _, err := RefreshToken(ctx, &config); err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 2 {
				t.Fatalf("requests after refresh = %d, want 2", requests.Load())
			}
			InvalidateTokenValue(&config, "test-token")
			if _, err := GetToken(ctx, &config); err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 3 {
				t.Fatalf("requests after invalidation = %d, want 3", requests.Load())
			}
		})
	}
}
