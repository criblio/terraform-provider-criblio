package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultCloudDomain = "cribl.cloud"
	tokenExpiryBuffer  = 30 * time.Second
	onPremTokenTTL     = time.Hour
)

var (
	tokenCache sync.Map
	timeNow    = time.Now
)

// GetToken returns an OAuth2 or on-prem auth token for the given config.
func GetToken(ctx context.Context, config *CriblConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config is required for authentication")
	}

	key, err := tokenCacheKey(config)
	if err != nil {
		return "", err
	}

	if cached, ok := tokenCache.Load(key); ok {
		tokenInfo, ok := cached.(*TokenInfo)
		if ok && tokenInfo.ExpiresAt.After(timeNow().Add(tokenExpiryBuffer)) {
			return tokenInfo.Token, nil
		}
	}

	var tokenInfo *TokenInfo
	if IsOnPrem(config) {
		tokenInfo, err = getOnPremBearerToken(ctx, config)
	} else {
		tokenInfo, err = getCloudBearerToken(ctx, config)
	}
	if err != nil {
		return "", err
	}

	tokenCache.Store(key, tokenInfo)
	return tokenInfo.Token, nil
}

// IsOnPrem reports whether config indicates an on-prem deployment.
func IsOnPrem(config *CriblConfig) bool {
	return config != nil && config.OnpremServerURL != ""
}

// IsGov reports whether domain requires the gov Okta OAuth2 flow.
func IsGov(domain string) bool {
	return strings.Contains(domain, "gov")
}

// ClearTokenCache clears all cached tokens.
func ClearTokenCache() {
	tokenCache.Range(func(key, _ interface{}) bool {
		tokenCache.Delete(key)
		return true
	})
}

// InvalidateToken removes the cached token for config.
func InvalidateToken(config *CriblConfig) {
	key, err := tokenCacheKey(config)
	if err != nil {
		return
	}
	tokenCache.Delete(key)
}

func tokenCacheKey(config *CriblConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config is required for authentication")
	}
	if IsOnPrem(config) {
		if config.OnpremUsername == "" || config.OnpremPassword == "" {
			return "", fmt.Errorf("on-prem authentication requires username and password")
		}
		return fmt.Sprintf("onprem:%s:%s:%s", config.OnpremServerURL, config.OnpremUsername, config.OnpremPassword), nil
	}
	if config.ClientID == "" || config.ClientSecret == "" {
		return "", fmt.Errorf("cloud authentication requires client ID and client secret")
	}
	return fmt.Sprintf("%s:%s:%s", config.ClientID, config.ClientSecret, cloudDomain(config)), nil
}

func cloudDomain(config *CriblConfig) string {
	switch {
	case config != nil && config.CloudDomain != "":
		return config.CloudDomain
	case os.Getenv("CRIBL_CLOUD_DOMAIN") != "":
		return os.Getenv("CRIBL_CLOUD_DOMAIN")
	default:
		return defaultCloudDomain
	}
}

func getOnPremBearerToken(ctx context.Context, config *CriblConfig) (*TokenInfo, error) {
	authURL := fmt.Sprintf("%s/api/v1/auth/login", strings.TrimRight(config.OnpremServerURL, "/"))
	log.Printf("[DEBUG] Getting on-prem bearer token from: %s", authURL)

	requestBody := map[string]string{
		"username": config.OnpremUsername,
		"password": config.OnpremPassword,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	responseBody, err := doTokenRequest(ctx, http.MethodPost, authURL, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("failed to get on-prem bearer token: %v", err)
	}

	var result struct {
		Token               string `json:"token"`
		ForcePasswordChange bool   `json:"forcePasswordChange"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &TokenInfo{
		Token:     strings.TrimPrefix(result.Token, "Bearer "),
		ExpiresAt: timeNow().Add(onPremTokenTTL),
	}, nil
}

func getCloudBearerToken(ctx context.Context, config *CriblConfig) (*TokenInfo, error) {
	domain := cloudDomain(config)
	authURL, audience, useFormEncoded, err := cloudAuthSettings(domain)
	if err != nil {
		return nil, err
	}

	var body []byte
	var contentType string
	if useFormEncoded {
		formData := url.Values{}
		formData.Set("grant_type", "client_credentials")
		formData.Set("client_id", config.ClientID)
		formData.Set("client_secret", config.ClientSecret)
		formData.Set("audience", audience)

		body = []byte(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else {
		requestBody := map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     config.ClientID,
			"client_secret": config.ClientSecret,
			"audience":      audience,
		}

		body, err = json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		contentType = "application/json"
	}

	responseBody, err := doTokenRequest(ctx, http.MethodPost, authURL, contentType, body)
	if err != nil {
		return nil, err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &TokenInfo{
		Token:     result.AccessToken,
		ExpiresAt: timeNow().Add(time.Duration(result.ExpiresIn) * time.Second),
	}, nil
}

func cloudAuthSettings(domain string) (string, string, bool, error) {
	if domain == "" {
		domain = defaultCloudDomain
	}

	if baseURL, ok, err := localBaseURL(domain); err != nil {
		return "", "", false, err
	} else if ok {
		return fmt.Sprintf("%s/oauth/token", strings.TrimRight(baseURL, "/")), baseURL, false, nil
	}

	audience := fmt.Sprintf("https://api.%s", domain)
	if !IsGov(domain) {
		return fmt.Sprintf("https://login.%s/oauth/token", domain), audience, false, nil
	}

	oktaDomain := os.Getenv("CRIBL_OKTA_DOMAIN")
	if oktaDomain == "" {
		oktaDomain = strings.ReplaceAll(domain, "cribl-gov-", "criblgov-")
		oktaDomain = strings.ReplaceAll(oktaDomain, "staging", "stg")
		oktaDomain = strings.ReplaceAll(oktaDomain, ".cloud", ".okta.com")
	}

	authServerID := os.Getenv("CRIBL_OKTA_AUTH_SERVER_ID")
	if authServerID == "" {
		authServerID = os.Getenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID")
	}
	if authServerID == "" {
		return "", "", false, fmt.Errorf("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID environment variable is required for gov domains")
	}

	oktaBaseURL, ok, err := localBaseURL(oktaDomain)
	if err != nil {
		return "", "", false, err
	}
	if ok {
		return fmt.Sprintf("%s/oauth2/%s/v1/token", strings.TrimRight(oktaBaseURL, "/"), authServerID), audience, true, nil
	}

	return fmt.Sprintf("https://%s/oauth2/%s/v1/token", oktaDomain, authServerID), audience, true, nil
}

func localBaseURL(input string) (string, bool, error) {
	parsedURL, err := url.Parse(input)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		if IsLocalHost(parsedURL.Host) {
			return strings.TrimRight(parsedURL.String(), "/"), true, nil
		}
		return "", false, nil
	}
	if err != nil && strings.Contains(input, "://") {
		return "", false, fmt.Errorf("failed to parse local URL: %v", err)
	}
	if IsLocalHost(input) {
		return fmt.Sprintf("http://%s", input), true, nil
	}
	return "", false, nil
}

func doTokenRequest(ctx context.Context, method, requestURL, contentType string, body []byte) ([]byte, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 500 * time.Millisecond
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.Backoff = retryablehttp.DefaultBackoff
	retryClient.Logger = &retryableLogger{}
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("[DEBUG] 429 Too Many Requests, will retry")
			return true, nil
		}
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	req, err := retryablehttp.NewRequest(method, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create retryable request: %v", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)

	resp, err := retryClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request with retry: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		return responseBody, nil
	}

	return nil, fmt.Errorf("failed to do request: status=%d, body=%s", resp.StatusCode, string(responseBody))
}

type retryableLogger struct{}

func (l *retryableLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] retryablehttp: %s %v", msg, keysAndValues)
}

func (l *retryableLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] retryablehttp: %s %v", msg, keysAndValues)
}

func (l *retryableLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Printf("[DEBUG] retryablehttp: %s %v", msg, keysAndValues)
}

func (l *retryableLogger) Warn(msg string, keysAndValues ...interface{}) {
	log.Printf("[WARN] retryablehttp: %s %v", msg, keysAndValues)
}
