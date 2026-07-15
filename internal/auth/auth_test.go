package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/useragent"
)

func TestGetTokenCloud(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	var receivedContentType string
	var receivedUserAgent string
	var receivedBody struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Audience     string `json:"audience"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		receivedContentType = r.Header.Get("Content-Type")
		receivedUserAgent = r.Header.Get("User-Agent")
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"cloud-token","expires_in":3600}`))
	}))
	defer server.Close()

	token, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "cloud-token" {
		t.Fatalf("GetToken returned %q, expected %q", token, "cloud-token")
	}
	if receivedContentType != "application/json" {
		t.Errorf("content type = %q, expected application/json", receivedContentType)
	}
	if receivedUserAgent != useragent.TerraformProvider {
		t.Errorf("User-Agent = %q, expected %q", receivedUserAgent, useragent.TerraformProvider)
	}
	if receivedBody.GrantType != "client_credentials" {
		t.Errorf("grant_type = %q, expected client_credentials", receivedBody.GrantType)
	}
	if receivedBody.ClientID != "client-id" {
		t.Errorf("client_id = %q, expected client-id", receivedBody.ClientID)
	}
	if receivedBody.ClientSecret != "client-secret" {
		t.Errorf("client_secret = %q, expected client-secret", receivedBody.ClientSecret)
	}
	if receivedBody.Audience != server.URL {
		t.Errorf("audience = %q, expected %q", receivedBody.Audience, server.URL)
	}
}

func TestGetTokenLocalhostURL(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"localhost-token","expires_in":3600}`))
	}))
	defer server.Close()

	parsedURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}
	localhostURL := "http://localhost" + parsedURL.Port()
	if parsedURL.Port() != "" {
		localhostURL = "http://localhost:" + parsedURL.Port()
	}

	token, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  localhostURL,
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "localhost-token" {
		t.Fatalf("GetToken returned %q, expected %q", token, "localhost-token")
	}
}

func TestGetTokenLoopbackIP(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"loopback-token","expires_in":3600}`))
	}))
	defer server.Close()

	token, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "loopback-token" {
		t.Fatalf("GetToken returned %q, expected %q", token, "loopback-token")
	}
}

func TestCloudAuthSettingsStandardDomains(t *testing.T) {
	testCases := []struct {
		name         string
		domain       string
		expectedURL  string
		expectedAud  string
		expectedForm bool
	}{
		{
			name:        "cloud",
			domain:      "cribl.cloud",
			expectedURL: "https://login.cribl.cloud/oauth/token",
			expectedAud: "https://api.cribl.cloud",
		},
		{
			name:        "playground",
			domain:      "cribl-playground.cloud",
			expectedURL: "https://login.cribl-playground.cloud/oauth/token",
			expectedAud: "https://api.cribl-playground.cloud",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			authURL, audience, useFormEncoded, err := cloudAuthSettings(test.domain)
			if err != nil {
				t.Fatalf("cloudAuthSettings returned error: %v", err)
			}
			if authURL != test.expectedURL {
				t.Fatalf("auth URL = %q, expected %q", authURL, test.expectedURL)
			}
			if audience != test.expectedAud {
				t.Fatalf("audience = %q, expected %q", audience, test.expectedAud)
			}
			if useFormEncoded != test.expectedForm {
				t.Fatalf("useFormEncoded = %v, expected %v", useFormEncoded, test.expectedForm)
			}
		})
	}
}

func TestGetTokenGovDomain(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	var receivedContentType string
	var receivedBody string
	var requestedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		receivedContentType = r.Header.Get("Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"gov-token","expires_in":3600}`))
	}))
	defer server.Close()

	t.Setenv("CRIBL_OKTA_DOMAIN", server.URL)
	t.Setenv("CRIBL_OKTA_AUTH_SERVER_ID", "auth-server-id")
	t.Setenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", "")

	token, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "gov-client-id",
		ClientSecret: "gov-client-secret",
		CloudDomain:  "cribl-gov-staging.cloud",
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "gov-token" {
		t.Fatalf("GetToken returned %q, expected %q", token, "gov-token")
	}
	if requestedPath != "/oauth2/auth-server-id/v1/token" {
		t.Errorf("requested path = %q, expected /oauth2/auth-server-id/v1/token", requestedPath)
	}
	if receivedContentType != "application/x-www-form-urlencoded" {
		t.Errorf("content type = %q, expected application/x-www-form-urlencoded", receivedContentType)
	}
	if !strings.Contains(receivedBody, "grant_type=client_credentials") {
		t.Errorf("form body missing grant_type: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "client_id=gov-client-id") {
		t.Errorf("form body missing client_id: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "client_secret=gov-client-secret") {
		t.Errorf("form body missing client_secret: %s", receivedBody)
	}
}

func TestGetTokenGovDomainMissingAuthServerID(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")
	t.Setenv("CRIBL_OKTA_DOMAIN", "")
	t.Setenv("CRIBL_OKTA_AUTH_SERVER_ID", "")
	t.Setenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", "")

	_, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "gov-client-id",
		ClientSecret: "gov-client-secret",
		CloudDomain:  "cribl-gov-staging.cloud",
	})
	if err == nil {
		t.Fatal("GetToken returned nil error, expected missing auth server ID error")
	}
	if !strings.Contains(err.Error(), "CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID") {
		t.Fatalf("error = %q, expected CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", err.Error())
	}
}

func TestGetTokenOnPrem(t *testing.T) {
	ClearTokenCache()

	var receivedBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/login" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"onprem-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	token, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "onprem-token" {
		t.Fatalf("GetToken returned %q, expected %q", token, "onprem-token")
	}
	if receivedBody.Username != "admin" {
		t.Errorf("username = %q, expected admin", receivedBody.Username)
	}
	if receivedBody.Password != "secret" {
		t.Errorf("password = %q, expected secret", receivedBody.Password)
	}
}

func TestOnPremTokenStripsBearer(t *testing.T) {
	ClearTokenCache()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"Bearer stripped-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	token, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "stripped-token" {
		t.Fatalf("GetToken returned %q, expected %q", token, "stripped-token")
	}
}

func TestTokenCacheHit(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"cached-token","expires_in":3600}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	}

	for range 2 {
		token, err := GetToken(context.Background(), config)
		if err != nil {
			t.Fatalf("GetToken returned error: %v", err)
		}
		if token != "cached-token" {
			t.Fatalf("GetToken returned %q, expected %q", token, "cached-token")
		}
	}

	if requestCount != 1 {
		t.Fatalf("request count = %d, expected 1", requestCount)
	}
}

func TestGetTokenOnPremConcurrentSingleflight(t *testing.T) {
	ClearTokenCache()

	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"concurrent-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}

	var wg sync.WaitGroup
	errs := make(chan error, 25)
	for range 25 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := GetToken(context.Background(), config)
			if err != nil {
				errs <- err
				return
			}
			if token != "concurrent-token" {
				errs <- fmt.Errorf("token = %q, expected concurrent-token", token)
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent GetToken returned error: %v", err)
		}
	}
	if got := atomic.LoadInt32(&requestCount); got != 1 {
		t.Fatalf("login request count = %d, expected 1", got)
	}
}

func TestOnPremCacheKeyNormalizesTrailingSlash(t *testing.T) {
	ClearTokenCache()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"normalized-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}
	slashedConfig := &CriblConfig{
		OnpremServerURL: server.URL + "/",
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}

	if _, err := GetToken(context.Background(), config); err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if _, err := GetToken(context.Background(), slashedConfig); err != nil {
		t.Fatalf("GetToken with trailing slash returned error: %v", err)
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected 1", requestCount)
	}
}

func TestOnPremCacheKeySeparatesCredentials(t *testing.T) {
	ClearTokenCache()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		var requestBody struct {
			Username string `json:"username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"` + requestBody.Username + `-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	first, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "first",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("first GetToken returned error: %v", err)
	}
	second, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "second",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("second GetToken returned error: %v", err)
	}

	if first != "first-token" {
		t.Fatalf("first token = %q, expected first-token", first)
	}
	if second != "second-token" {
		t.Fatalf("second token = %q, expected second-token", second)
	}
	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2", requestCount)
	}
}

func TestCloudCacheKeySeparatesDomains(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	firstCount := 0
	firstServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"first-domain-token","expires_in":3600}`))
	}))
	defer firstServer.Close()

	secondCount := 0
	secondServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"second-domain-token","expires_in":3600}`))
	}))
	defer secondServer.Close()

	first, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  firstServer.URL,
	})
	if err != nil {
		t.Fatalf("first GetToken returned error: %v", err)
	}
	second, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  secondServer.URL,
	})
	if err != nil {
		t.Fatalf("second GetToken returned error: %v", err)
	}

	if first != "first-domain-token" {
		t.Fatalf("first token = %q, expected first-domain-token", first)
	}
	if second != "second-domain-token" {
		t.Fatalf("second token = %q, expected second-domain-token", second)
	}
	if firstCount != 1 || secondCount != 1 {
		t.Fatalf("request counts = %d, %d; expected 1, 1", firstCount, secondCount)
	}
}

func TestTokenCacheExpiry(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"short-token","expires_in":10}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	}

	for range 2 {
		if _, err := GetToken(context.Background(), config); err != nil {
			t.Fatalf("GetToken returned error: %v", err)
		}
	}

	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2 because token expires inside buffer", requestCount)
	}
}

func TestTokenCacheDistinctKeys(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		var requestBody struct {
			ClientID string `json:"client_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"` + requestBody.ClientID + `-token","expires_in":3600}`))
	}))
	defer server.Close()

	firstToken, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "first",
		ClientSecret: "secret",
		CloudDomain:  server.URL,
	})
	if err != nil {
		t.Fatalf("first GetToken returned error: %v", err)
	}

	secondToken, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "second",
		ClientSecret: "secret",
		CloudDomain:  server.URL,
	})
	if err != nil {
		t.Fatalf("second GetToken returned error: %v", err)
	}

	if firstToken != "first-token" {
		t.Fatalf("first token = %q, expected first-token", firstToken)
	}
	if secondToken != "second-token" {
		t.Fatalf("second token = %q, expected second-token", secondToken)
	}
	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2", requestCount)
	}
}

func TestTokenCacheCloudKeyMatchesHook(t *testing.T) {
	key, err := tokenCacheKey(&CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  "cribl-playground.cloud",
	})
	if err != nil {
		t.Fatalf("tokenCacheKey returned error: %v", err)
	}
	if strings.Contains(key, "client-id") || strings.Contains(key, "client-secret") {
		t.Fatalf("cloud cache key = %q, expected no raw credential material", key)
	}
	if !strings.HasPrefix(key, "cloud:") {
		t.Fatalf("cloud cache key = %q, expected cloud prefix", key)
	}
}

func TestTokenCacheOnPremKeyDoesNotExposeCredentials(t *testing.T) {
	key, err := tokenCacheKey(&CriblConfig{
		OnpremServerURL: "https://example.local",
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("tokenCacheKey returned error: %v", err)
	}
	if strings.Contains(key, "admin") || strings.Contains(key, "secret") || strings.Contains(key, "example.local") {
		t.Fatalf("on-prem cache key = %q, expected no raw credential material", key)
	}
	if !strings.HasPrefix(key, "onprem:") {
		t.Fatalf("on-prem cache key = %q, expected onprem prefix", key)
	}
}

func TestInvalidateToken(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"token","expires_in":3600}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	}

	if _, err := GetToken(context.Background(), config); err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	InvalidateToken(config)
	if _, err := GetToken(context.Background(), config); err != nil {
		t.Fatalf("GetToken after invalidation returned error: %v", err)
	}

	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2", requestCount)
	}
}

func TestInvalidateTokenValueOnlyRemovesMatchingToken(t *testing.T) {
	ClearTokenCache()

	key, err := tokenCacheKey(&CriblConfig{
		OnpremServerURL: "https://example.local",
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("tokenCacheKey returned error: %v", err)
	}
	config := &CriblConfig{
		OnpremServerURL: "https://example.local",
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}

	tokenCache.Store(key, &TokenInfo{Token: "fresh-token", ExpiresAt: timeNow().Add(time.Hour)})
	InvalidateTokenValue(config, "stale-token")
	if _, ok := loadCachedToken(key); !ok {
		t.Fatal("InvalidateTokenValue removed non-matching token")
	}

	InvalidateTokenValue(config, "fresh-token")
	if _, ok := loadCachedToken(key); ok {
		t.Fatal("InvalidateTokenValue kept matching token")
	}
}

func TestRefreshTokenInvalidatesCache(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"token","expires_in":3600}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	}

	if _, err := GetToken(context.Background(), config); err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if _, err := RefreshToken(context.Background(), config); err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}

	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2", requestCount)
	}
}

func TestTokenNon200Error(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid credentials"}`))
	}))
	defer server.Close()

	_, err := GetToken(context.Background(), &CriblConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CloudDomain:  server.URL,
	})
	if err == nil {
		t.Fatal("GetToken returned nil error, expected non-200 error")
	}
	if !strings.Contains(err.Error(), "status=401") {
		t.Fatalf("error = %q, expected status=401", err.Error())
	}
}

func TestTokenRequestRetriesTransientFailures(t *testing.T) {
	ClearTokenCache()
	fastTokenRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount < 3 {
			http.Error(w, "temporary", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"retry-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	token, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err != nil {
		t.Fatalf("GetToken returned error: %v", err)
	}
	if token != "retry-token" {
		t.Fatalf("token = %q, expected retry-token", token)
	}
	if requestCount != 3 {
		t.Fatalf("request count = %d, expected 3", requestCount)
	}
}

func TestTokenRequestReportsFinalRetriedHTTPResponse(t *testing.T) {
	ClearTokenCache()
	fastTokenRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"leader is still starting"}`))
	}))
	defer server.Close()

	_, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err == nil {
		t.Fatal("GetToken returned nil error, expected retry exhaustion")
	}
	if requestCount != tokenRetryMax+1 {
		t.Fatalf("request count = %d, expected %d", requestCount, tokenRetryMax+1)
	}
	for _, want := range []string{"status=500", `body={"error":"leader is still starting"}`} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, expected %q", err.Error(), want)
		}
	}
	if strings.Contains(err.Error(), "giving up after") {
		t.Fatalf("error = %q, expected auth-layer response details instead of retryablehttp exhaustion", err.Error())
	}
}

func TestTokenRequestDoesNotRetryUnauthorized(t *testing.T) {
	ClearTokenCache()
	fastTokenRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := GetToken(context.Background(), &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err == nil {
		t.Fatal("GetToken returned nil error, expected unauthorized error")
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected 1", requestCount)
	}
	if !strings.Contains(err.Error(), "after 1 attempt") {
		t.Fatalf("error = %q, expected actual attempt count", err.Error())
	}
}

func TestInvalidTokenResponsesAreNotCached(t *testing.T) {
	testCases := []struct {
		name      string
		firstBody string
	}{
		{
			name:      "empty on-prem token",
			firstBody: `{"token":"","forcePasswordChange":false}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ClearTokenCache()

			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.WriteHeader(http.StatusOK)
				if requestCount == 1 {
					_, _ = w.Write([]byte(test.firstBody))
					return
				}
				_, _ = w.Write([]byte(`{"token":"valid-token","forcePasswordChange":false}`))
			}))
			defer server.Close()

			config := &CriblConfig{
				OnpremServerURL: server.URL,
				OnpremUsername:  "admin",
				OnpremPassword:  "secret",
			}
			if _, err := GetToken(context.Background(), config); err == nil {
				t.Fatal("first GetToken returned nil error, expected invalid response error")
			}
			token, err := GetToken(context.Background(), config)
			if err != nil {
				t.Fatalf("second GetToken returned error: %v", err)
			}
			if token != "valid-token" {
				t.Fatalf("token = %q, expected valid-token", token)
			}
			if requestCount != 2 {
				t.Fatalf("request count = %d, expected 2", requestCount)
			}
		})
	}
}

func TestOnPremTokenAllowsForcePasswordChangeWhenTokenPresent(t *testing.T) {
	ClearTokenCache()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"password-change-token","forcePasswordChange":true}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}

	for range 2 {
		token, err := GetToken(context.Background(), config)
		if err != nil {
			t.Fatalf("GetToken returned error: %v", err)
		}
		if token != "password-change-token" {
			t.Fatalf("token = %q, expected password-change-token", token)
		}
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected cached token after first login", requestCount)
	}
}

func TestInvalidCloudTokenResponsesAreNotCached(t *testing.T) {
	testCases := []struct {
		name      string
		firstBody string
	}{
		{
			name:      "empty access token",
			firstBody: `{"access_token":"","expires_in":3600}`,
		},
		{
			name:      "invalid expiry",
			firstBody: `{"access_token":"token","expires_in":0}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ClearTokenCache()
			t.Setenv("CRIBL_CLOUD_DOMAIN", "")

			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.WriteHeader(http.StatusOK)
				if requestCount == 1 {
					_, _ = w.Write([]byte(test.firstBody))
					return
				}
				_, _ = w.Write([]byte(`{"access_token":"valid-token","expires_in":3600}`))
			}))
			defer server.Close()

			config := &CriblConfig{
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				CloudDomain:  server.URL,
			}
			if _, err := GetToken(context.Background(), config); err == nil {
				t.Fatal("first GetToken returned nil error, expected invalid response error")
			}
			token, err := GetToken(context.Background(), config)
			if err != nil {
				t.Fatalf("second GetToken returned error: %v", err)
			}
			if token != "valid-token" {
				t.Fatalf("token = %q, expected valid-token", token)
			}
			if requestCount != 2 {
				t.Fatalf("request count = %d, expected 2", requestCount)
			}
		})
	}
}

func TestTokenRequestRespectsContextCancellation(t *testing.T) {
	ClearTokenCache()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "temporary", http.StatusInternalServerError)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := GetToken(ctx, &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	})
	if err == nil {
		t.Fatal("GetToken returned nil error, expected context cancellation")
	}
	if !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Fatalf("error = %q, expected context cancellation", err.Error())
	}
}

func TestSharedTokenFetchSurvivesLeadingCallerCancellation(t *testing.T) {
	ClearTokenCache()

	var requestCount int32
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		started <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"shared-token","forcePasswordChange":false}`))
	}))
	defer server.Close()

	config := &CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}
	leaderCtx, cancelLeader := context.WithCancel(context.Background())
	leaderErr := make(chan error, 1)
	go func() {
		_, err := GetToken(leaderCtx, config)
		leaderErr <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for token request to start")
	}

	cancelLeader()
	select {
	case err := <-leaderErr:
		if err == nil || !strings.Contains(err.Error(), context.Canceled.Error()) {
			t.Fatalf("leader error = %v, expected context cancellation", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for leader cancellation")
	}

	followerDone := make(chan struct{})
	var followerToken string
	var followerErr error
	go func() {
		defer close(followerDone)
		followerToken, followerErr = GetToken(context.Background(), config)
	}()

	close(release)
	select {
	case <-followerDone:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for follower token")
	}
	if followerErr != nil {
		t.Fatalf("follower GetToken returned error: %v", followerErr)
	}
	if followerToken != "shared-token" {
		t.Fatalf("follower token = %q, expected shared-token", followerToken)
	}
	if got := atomic.LoadInt32(&requestCount); got != 1 {
		t.Fatalf("request count = %d, expected 1 shared request", got)
	}
}

func TestConciseBodyRedactsNonJSONSecrets(t *testing.T) {
	body := []byte("client_secret=super-secret&password=hunter2 token:abc123 access_token=xyz")
	got := conciseBody(body)

	for _, leaked := range []string{"super-secret", "hunter2", "abc123", "xyz"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("conciseBody leaked %q in %q", leaked, got)
		}
	}
	if !strings.Contains(got, "(sensitive)") {
		t.Fatalf("conciseBody = %q, expected redacted markers", got)
	}
}

func TestIsOnPremAndIsGov(t *testing.T) {
	if !IsOnPrem(&CriblConfig{OnpremServerURL: "http://localhost:9000"}) {
		t.Fatal("IsOnPrem returned false, expected true")
	}
	if IsOnPrem(&CriblConfig{}) {
		t.Fatal("IsOnPrem returned true, expected false")
	}
	if !IsGov("cribl-gov-staging.cloud") {
		t.Fatal("IsGov returned false, expected true")
	}
	if IsGov("cribl.cloud") {
		t.Fatal("IsGov returned true, expected false")
	}
}

func fastTokenRetry(t *testing.T) {
	t.Helper()

	oldMin := tokenRetryWaitMin
	oldMax := tokenRetryWaitMax
	oldCap := tokenBackoffCap
	tokenRetryWaitMin = time.Millisecond
	tokenRetryWaitMax = time.Millisecond
	tokenBackoffCap = time.Millisecond
	t.Cleanup(func() {
		tokenRetryWaitMin = oldMin
		tokenRetryWaitMax = oldMax
		tokenBackoffCap = oldCap
	})
}
