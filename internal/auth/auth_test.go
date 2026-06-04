package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetTokenCloud(t *testing.T) {
	ClearTokenCache()
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	var receivedContentType string
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
