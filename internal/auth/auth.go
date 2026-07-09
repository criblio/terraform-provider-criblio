package auth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/useragent"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/sync/singleflight"
)

const (
	defaultCloudDomain = "cribl.cloud"
	tokenExpiryBuffer  = 30 * time.Second
	onPremTokenTTL     = time.Hour
	tokenRetryMax      = 5
	maxErrorBodyLength = 512
)

var (
	sensitiveTextPattern = regexp.MustCompile(`(?i)\b(client[_-]?secret|password|access[_-]?token|refresh[_-]?token|token)(\s*[:=]\s*)([^&\s,;]+)`)
	tokenFetchTimeout    = 30 * time.Second
	tokenRetryWaitMin    = 500 * time.Millisecond
	tokenRetryWaitMax    = 1500 * time.Millisecond
	tokenBackoffCap      = 5 * time.Second
	tokenCache           sync.Map
	tokenFetchGroup      singleflight.Group
	timeNow              = time.Now
)

// GetToken returns an OAuth2 or on-prem auth token for the given config.
func GetToken(ctx context.Context, config *CriblConfig) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context is required for authentication")
	}
	if config == nil {
		return "", fmt.Errorf("config is required for authentication")
	}

	key, err := tokenCacheKey(config)
	if err != nil {
		return "", err
	}

	if tokenInfo, ok := loadCachedToken(key); ok {
		return tokenInfo.Token, nil
	}

	ch := tokenFetchGroup.DoChan(key, func() (interface{}, error) {
		if tokenInfo, ok := loadCachedToken(key); ok {
			return tokenInfo, nil
		}

		fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), tokenFetchTimeout)
		defer cancel()

		tokenInfo, err := fetchToken(fetchCtx, config)
		if err != nil {
			return nil, err
		}

		tokenCache.Store(key, tokenInfo)
		return tokenInfo, nil
	})

	var result interface{}
	select {
	case call := <-ch:
		if call.Err != nil {
			return "", call.Err
		}
		result = call.Val
	case <-ctx.Done():
		return "", ctx.Err()
	}
	tokenInfo, ok := result.(*TokenInfo)
	if !ok || tokenInfo == nil {
		return "", fmt.Errorf("failed to fetch token: unexpected cache value")
	}
	return tokenInfo.Token, nil
}

func validCachedToken(cached interface{}) (*TokenInfo, bool) {
	tokenInfo, ok := cached.(*TokenInfo)
	if !ok || tokenInfo == nil {
		return nil, false
	}
	if tokenInfo.Token == "" {
		return nil, false
	}
	return tokenInfo, tokenInfo.ExpiresAt.After(timeNow().Add(tokenExpiryBuffer))
}

func loadCachedToken(key string) (*TokenInfo, bool) {
	cached, ok := tokenCache.Load(key)
	if !ok {
		return nil, false
	}
	return validCachedToken(cached)
}

func fetchToken(ctx context.Context, config *CriblConfig) (*TokenInfo, error) {
	var (
		tokenInfo *TokenInfo
		err       error
	)
	if IsOnPrem(config) {
		tokenInfo, err = getOnPremBearerToken(ctx, config)
	} else {
		tokenInfo, err = getCloudBearerToken(ctx, config)
	}
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
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

// InvalidateTokenValue removes the cached token for config only when it matches token.
func InvalidateTokenValue(config *CriblConfig, token string) {
	key, err := tokenCacheKey(config)
	if err != nil {
		return
	}
	cached, ok := tokenCache.Load(key)
	if !ok {
		return
	}
	tokenInfo, ok := cached.(*TokenInfo)
	if !ok || tokenInfo == nil || tokenInfo.Token != token {
		return
	}
	tokenCache.CompareAndDelete(key, tokenInfo)
}

// RefreshToken invalidates the cached token for config and fetches a fresh token.
func RefreshToken(ctx context.Context, config *CriblConfig) (string, error) {
	InvalidateToken(config)
	return GetToken(ctx, config)
}

func tokenCacheKey(config *CriblConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config is required for authentication")
	}
	if IsOnPrem(config) {
		if config.OnpremUsername == "" || config.OnpremPassword == "" {
			return "", fmt.Errorf("on-prem authentication requires username and password")
		}
		return "onprem:" + cacheDigest(normalizedServerURL(config.OnpremServerURL), config.OnpremUsername, config.OnpremPassword), nil
	}
	if config.ClientID == "" || config.ClientSecret == "" {
		return "", fmt.Errorf("cloud authentication requires client ID and client secret")
	}
	authURL, _, _, err := cloudAuthSettings(cloudDomain(config))
	if err != nil {
		return "", err
	}
	return "cloud:" + cacheDigest(authURL, config.ClientID, config.ClientSecret), nil
}

func cacheDigest(parts ...string) string {
	hash := sha256.New()
	for _, part := range parts {
		_, _ = hash.Write([]byte(part))
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func normalizedServerURL(serverURL string) string {
	return strings.TrimRight(serverURL, "/")
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
	authURL := fmt.Sprintf("%s/api/v1/auth/login", normalizedServerURL(config.OnpremServerURL))
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
		Token string `json:"token"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	token := strings.TrimPrefix(result.Token, "Bearer ")
	if token == "" {
		return nil, fmt.Errorf("on-prem authentication response did not include a token")
	}
	return &TokenInfo{
		Token:     token,
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
	if result.AccessToken == "" {
		return nil, fmt.Errorf("cloud authentication response did not include an access token")
	}
	if result.ExpiresIn <= 0 {
		return nil, fmt.Errorf("cloud authentication response included invalid expires_in=%d", result.ExpiresIn)
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
	attempts := 0
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = tokenRetryMax
	retryClient.RetryWaitMin = tokenRetryWaitMin
	retryClient.RetryWaitMax = tokenRetryWaitMax
	retryClient.Backoff = tokenBackoff
	retryClient.Logger = &retryableLogger{}
	retryClient.CheckRetry = tokenRetryPolicy
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, _ *http.Request, attempt int) {
		attempts = attempt + 1
	}

	req, err := retryablehttp.NewRequest(method, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create retryable request: %v", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", useragent.TerraformProvider)

	resp, err := retryClient.Do(req)
	if attempts == 0 {
		attempts = 1
	}
	if err != nil {
		return nil, fmt.Errorf("failed to make token request after %d attempt(s): %s %s: %v", attempts, method, sanitizedRequestURL(requestURL), err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		return responseBody, nil
	}

	return nil, fmt.Errorf("token request failed after %d attempt(s): %s %s: status=%d body=%s", attempts, method, sanitizedRequestURL(requestURL), resp.StatusCode, conciseBody(responseBody))
}

func tokenRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		shouldRetry, retryErr := retryablehttp.ErrorPropagatedRetryPolicy(ctx, resp, err)
		if retryErr != nil {
			return false, retryErr
		}
		return shouldRetry, nil
	}
	if resp == nil {
		return false, nil
	}
	switch resp.StatusCode {
	case http.StatusRequestTimeout, http.StatusTooManyRequests:
		return true, nil
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return false, nil
	}
	if resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented {
		return true, nil
	}
	return false, nil
}

func tokenBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	wait := retryablehttp.RateLimitLinearJitterBackoff(min, max, attemptNum, resp)
	if wait > tokenBackoffCap {
		return tokenBackoffCap
	}
	return wait
}

func sanitizedRequestURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "<invalid-url>"
	}
	parsedURL.RawQuery = ""
	parsedURL.User = nil
	return parsedURL.String()
}

func conciseBody(body []byte) string {
	value := strings.TrimSpace(string(body))
	if redacted, ok := redactedJSONBody([]byte(value)); ok {
		value = redacted
	} else {
		value = redactSensitiveText(value)
	}
	if len(value) > maxErrorBodyLength {
		return value[:maxErrorBodyLength] + "...(truncated)"
	}
	if value == "" {
		return "<empty>"
	}
	return value
}

func redactSensitiveText(value string) string {
	return sensitiveTextPattern.ReplaceAllString(value, `${1}${2}(sensitive)`)
}

func redactedJSONBody(body []byte) (string, bool) {
	if len(body) == 0 {
		return "", false
	}
	var value interface{}
	if err := json.Unmarshal(body, &value); err != nil {
		return "", false
	}
	redactJSONValue(value)
	redacted, err := json.Marshal(value)
	if err != nil {
		return "", false
	}
	return string(redacted), true
}

func redactJSONValue(value interface{}) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, child := range typed {
			if isSensitiveResponseKey(key) {
				typed[key] = "(sensitive)"
				continue
			}
			redactJSONValue(child)
		}
	case []interface{}:
		for _, child := range typed {
			redactJSONValue(child)
		}
	}
}

func isSensitiveResponseKey(key string) bool {
	lowerKey := strings.ToLower(key)
	return strings.Contains(lowerKey, "token") ||
		strings.Contains(lowerKey, "secret") ||
		strings.Contains(lowerKey, "password")
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
