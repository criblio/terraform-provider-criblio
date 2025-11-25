package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	_ sdkInitHook       = (*CriblTerraformHook)(nil)
	_ beforeRequestHook = (*CriblTerraformHook)(nil)
	_ afterErrorHook    = (*CriblTerraformHook)(nil)
)

func NewCriblTerraformHook() *CriblTerraformHook {
	log.Printf("[DEBUG] Creating new CriblTerraformHook")
	return &CriblTerraformHook{}
}

func (o *CriblTerraformHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	log.Printf("[DEBUG] Initializing SDK with baseURL: %s", baseURL)
	o.client = client

	// Check for on-prem configuration in environment variables first
	onpremServerURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	if onpremServerURL != "" {
		o.baseURL = onpremServerURL
		log.Printf("[DEBUG] On-prem configuration detected in environment, using server URL: %s", onpremServerURL)
		log.Printf("[DEBUG] Initialization complete")
		return onpremServerURL, client
	}

	// Get credentials from config or environment
	log.Printf("[DEBUG] Attempting to get credentials")
	config, err := GetCredentials()
	if err != nil {
		log.Printf("[ERROR] Failed to get credentials: %v", err)
		o.baseURL = baseURL
		return baseURL, client
	}

	// Check for on-prem configuration in credentials file
	if config != nil && config.OnpremServerURL != "" {
		o.baseURL = config.OnpremServerURL
		log.Printf("[DEBUG] On-prem configuration detected in credentials file, using server URL: %s", config.OnpremServerURL)
		log.Printf("[DEBUG] Initialization complete")
		return config.OnpremServerURL, client
	}

	// Set orgID and workspaceID from config
	if config != nil {
		log.Printf("[DEBUG] Setting orgID: %s and workspaceID: %s", config.OrganizationID, config.Workspace)
		o.orgID = config.OrganizationID
		o.workspaceID = config.Workspace

		// If baseURL is not provided or is a template, construct it from credentials
		input := ConstructBaseUrlInput{
			BaseURL: baseURL,
		}
		finalBaseURL := constructBaseURL(input, config)
		o.baseURL = finalBaseURL
		log.Printf("[DEBUG] Final baseURL: %s", finalBaseURL)
		log.Printf("[DEBUG] Initialization complete")
		return finalBaseURL, client
	} else {
		log.Printf("[DEBUG] No credentials found")
		o.baseURL = baseURL
	}

	log.Printf("[DEBUG] Initialization complete")
	return baseURL, client
}

func (o *CriblTerraformHook) BeforeRequest(ctx BeforeRequestContext, req *http.Request) (*http.Request, error) {

	// Check for on-prem configuration first
	onpremServerURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	onPrem := false
	if onpremServerURL != "" {
		onPrem = true
	}

	// First try to get credentials from security context
	var clientID, clientSecret, orgID, workspaceID, cloudDomain string

	//this should get moved into GetCredentials since we're getting creds from the securityCtx
	if !onPrem {
		if ctx.SecuritySource != nil {
			if security, err := ctx.SecuritySource(ctx.Context); err == nil {
				if s, ok := security.(shared.Security); ok {
					// Get OAuth credentials
					if s.ClientOauth != nil {
						clientID = s.ClientOauth.ClientID
						clientSecret = s.ClientOauth.ClientSecret
					}

					// Get org and workspace IDs from provider config (higher precedence than credentials file)
					if s.OrganizationID != nil {
						orgID = *s.OrganizationID
						o.orgID = orgID
					}
					if s.WorkspaceID != nil {
						workspaceID = *s.WorkspaceID
						o.workspaceID = workspaceID
					}
					if s.CloudDomain != nil {
						cloudDomain = *s.CloudDomain
					}
				}
			} else {
				log.Printf("[ERROR] Failed to get security info: %v", err)
			}
		}
	}

	var config *CriblConfig
	// Get credentials file config for fallback values
	// this function respects our ONPREM scheme and returns criblconfig with vars set correctly
	config, err := GetCredentials()
	if err != nil {
		log.Printf("[ERROR] Failed to get credentials from config: %v", err)
	}

	// If we don't have credentials from security context, use config file
	//then this can get nixed too
	if !onPrem {
		if clientID == "" || clientSecret == "" {
			if config != nil {
				clientID = config.ClientID
				clientSecret = config.ClientSecret
			}
		}

		// Reconstruct baseURL with proper precedence: Provider Config > Environment > Credentials File > Default
		input := ConstructBaseUrlInput{
			BaseURL:             o.baseURL,
			ProviderOrgID:       orgID,
			ProviderWorkspaceID: workspaceID,
			ProviderCloudDomain: cloudDomain,
		}
		o.baseURL = constructBaseURL(input, config)
	} else {
		if clientID == "" || clientSecret == "" {
			if config != nil {
				clientID = config.OnpremUsername
				clientSecret = config.OnpremPassword
			}
		}
	}

	// Handle authentication
	if bearerToken := os.Getenv("CRIBL_BEARER_TOKEN"); bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	} else if clientID != "" && clientSecret != "" {
		var authToken string
		if !onPrem {
			// Get audience from base URL
			audience := ""
			if o.baseURL != "" {
				// Extract domain from workspace URL (e.g., from https://main-org.cribl.cloud)
				parsedURL, err := url.Parse(o.baseURL)
				if err != nil {
					return req, fmt.Errorf("failed to parse base URL for audience: %v", err)
				}

				host := parsedURL.Host

				// Handle test/localhost URLs differently
				if strings.Contains(host, "127.0.0.1") || strings.Contains(host, "localhost") {
					// For test URLs, use the same URL as audience
					audience = o.baseURL
				} else {
					// Extract domain part after the first dash (remove workspace-org prefix)
					parts := strings.SplitN(host, ".", 2)
					if len(parts) < 2 {
						return req, fmt.Errorf("invalid workspace URL format for audience: %s", host)
					}
					domain := parts[1] // e.g., "cribl.cloud"

					audience = fmt.Sprintf("https://api.%s", domain)
				}
			} else if os.Getenv("CRIBL_AUDIENCE") != "" {
				audience = os.Getenv("CRIBL_AUDIENCE")
			} else {
				return req, fmt.Errorf("no base URL or audience provided")
			}

			// Get or create session
			sessionKey := fmt.Sprintf("%s:%s", clientID, clientSecret)
			var tokenInfo *TokenInfo

			if cachedTokenInfo, ok := o.sessions.Load(sessionKey); ok {
				tokenInfo = cachedTokenInfo.(*TokenInfo)
				if time.Until(tokenInfo.ExpiresAt) < 60*time.Minute {
					tokenInfo = nil
				}
			}

			if tokenInfo == nil {
				newTokenInfo, err := o.getBearerToken(ctx.Context, clientID, clientSecret, audience)
				if err != nil {
					return req, err
				}
				tokenInfo = newTokenInfo
				o.sessions.Store(sessionKey, tokenInfo)
			}
			authToken = tokenInfo.Token
		} else {
			// Check for cached token
			sessionKey := fmt.Sprintf("onprem:%s:%s:%s", config.OnpremServerURL, config.OnpremUsername, config.OnpremPassword)
			var tokenInfo *TokenInfo

			if cachedTokenInfo, ok := o.sessions.Load(sessionKey); ok {
				tokenInfo = cachedTokenInfo.(*TokenInfo)
				if time.Until(tokenInfo.ExpiresAt) < 30*time.Minute {
					tokenInfo = nil
				}
			}

			if tokenInfo == nil {
				log.Printf("[DEBUG] Retrieving bearer token using username/password for on-prem authentication")
				newTokenInfo, err := o.getOnPremBearerToken(ctx.Context, config.OnpremServerURL, config.OnpremUsername, config.OnpremPassword)
				if err != nil {
					return req, fmt.Errorf("failed to get on-prem bearer token: %v", err)
				}
				tokenInfo = newTokenInfo
				o.sessions.Store(sessionKey, tokenInfo)
			}

			authToken = tokenInfo.Token
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
	} else {
		return req, fmt.Errorf("authentication requires either CRIBL_BEARER_TOKEN OR CRIBL_ONPREM_USERNAME and CRIBL_ONPREM_PASSWORD OR Cloud stuff")
	}

	if !onPrem {
		// Handle URL routing
		path := strings.TrimLeft(req.URL.Path, "/")

		// Check if this is a gateway path (management endpoints)
		if isGatewayPath(path) || strings.Contains(req.URL.Host, "gateway.") {
			// Construct gateway URL
			gatewayURL := constructGatewayURL(cloudDomain, config)

			// Parse gateway URL to get the host
			parsedGatewayURL, err := url.Parse(gatewayURL)
			if err != nil {
				return req, fmt.Errorf("failed to parse gateway URL: %v", err)
			}

			// For gateway requests, don't add /api/v1 prefix - use path as-is
			newURL := fmt.Sprintf("%s/%s", strings.TrimRight(gatewayURL, "/"), path)

			parsedURL, err := url.Parse(newURL)
			if err != nil {
				return req, err
			}

			// Set both URL host and explicit Host header for gateway requests
			req.URL = parsedURL
			req.Host = parsedGatewayURL.Host
		} else {
			// Handle regular workspace API routing
			trimmedBaseURL := strings.TrimRight(o.baseURL, "/")

			// For workspace API, add /api/v1 prefix if not already present
			if !strings.Contains(req.URL.String(), "/api/v1") && !strings.HasPrefix(path, "api/v1") {
				newURL := fmt.Sprintf("%s/api/v1/%s", trimmedBaseURL, path)

				parsedURL, err := url.Parse(newURL)
				if err != nil {
					return req, err
				}

				req.URL = parsedURL
			}
		}
	} else {

		path := trimPath(req.URL.Path)

		// Check if this is a restricted endpoint for on-prem
		if isRestrictedOnPremEndpoint(path) {
			return req, fmt.Errorf("endpoint '%s' is not supported for on-prem deployments. On-prem deployments only support workspace resources (sources, destinations, routes, pipelines, etc.)", path)
		}

		// Handle URL routing for on-prem (always use serverURL/api/v1 path)
		baseURL := strings.TrimRight(config.OnpremServerURL, "/")

		// Construct full URL
		if path == "" {
			newURL := fmt.Sprintf("%s/api/v1", baseURL)
			parsedURL, err := url.Parse(newURL)
			if err != nil {
				return req, fmt.Errorf("failed to parse on-prem URL: %v", err)
			}
			req.URL = parsedURL
		} else {
			newURL := fmt.Sprintf("%s/api/v1/%s", baseURL, path)
			parsedURL, err := url.Parse(newURL)
			if err != nil {
				return req, fmt.Errorf("failed to parse on-prem URL: %v", err)
			}
			req.URL = parsedURL
		}

		log.Printf("[DEBUG] On-prem request URL: %s", req.URL.String())
	}

	return req, nil
}

// getOnPremBearerToken authenticates with on-prem server using username/password
// Reference: https://docs.cribl.io/cribl-as-code/authentication/#sdk-cust-managed-auth
func (o *CriblTerraformHook) getOnPremBearerToken(ctx context.Context, serverURL, username, password string) (*TokenInfo, error) {
	// Construct auth URL - use /api/v1/auth/login as per documentation
	baseURL := strings.TrimRight(serverURL, "/")
	authURL := fmt.Sprintf("%s/api/v1/auth/login", baseURL)

	log.Printf("[DEBUG] Getting on-prem bearer token from: %s", authURL)

	// Create request with username/password in JSON body
	requestBody := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := o.doRequestWithRetry(req)
	if err != nil {
		return nil, err
	}

	// Parse response - response contains "token" field with "Bearer " prefix
	var result struct {
		Token               string `json:"token"`
		ForcePasswordChange bool   `json:"forcePasswordChange"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Remove "Bearer " prefix if present
	token := strings.TrimPrefix(result.Token, "Bearer ")

	// Token TTL from global settings (default 3600 seconds / 1 hour)
	// Reference: Settings > Global > General Settings > API Server Settings > Advanced > Auth-token TTL
	expiresIn := 3600 // Default to 1 hour
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	log.Printf("[DEBUG] Successfully obtained on-prem bearer token (expires in %d seconds)", expiresIn)

	return &TokenInfo{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (o *CriblTerraformHook) getBearerToken(ctx context.Context, clientID, clientSecret, audience string) (*TokenInfo, error) {
	// Get auth URL from base URL
	authURL := ""
	useFormEncoded := false

	if o.baseURL != "" {
		// Extract domain from workspace URL (e.g., from https://main-org.cribl.cloud)
		parsedURL, err := url.Parse(o.baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %v", err)
		}

		host := parsedURL.Host

		// Handle test/localhost URLs differently
		if strings.Contains(host, "127.0.0.1") || strings.Contains(host, "localhost") {
			// For test URLs, use the same host but with /oauth/token path
			authURL = fmt.Sprintf("%s://%s/oauth/token", parsedURL.Scheme, host)
		} else {
			// Extract domain part after the first dash (remove workspace-org prefix)
			parts := strings.SplitN(host, ".", 2)
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid workspace URL format: %s", host)
			}
			domain := parts[1] // e.g., "cribl.cloud"

			// Check if domain contains "gov" - use Okta OAuth2 for gov domains
			if strings.Contains(domain, "gov") {
				// For gov domains, use Okta OAuth2 endpoint
				// Check for custom Okta settings from environment
				oktaDomain := os.Getenv("CRIBL_OKTA_DOMAIN")
				authServerID := os.Getenv("CRIBL_OKTA_AUTH_SERVER_ID")

				if oktaDomain == "" {
					// Default mapping: derive Okta domain from cloud domain
					// e.g., "cribl-gov-staging.cloud" -> "criblgov-stg.okta.com"
					oktaDomain = strings.ReplaceAll(domain, "cribl-gov-", "criblgov-")
					oktaDomain = strings.ReplaceAll(oktaDomain, "staging", "stg")
					oktaDomain = strings.ReplaceAll(oktaDomain, ".cloud", ".okta.com")
				}

				if authServerID == "" {
					// Authorization server ID must be configured via environment variable
					authServerID = os.Getenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID")
					if authServerID == "" {
						return nil, fmt.Errorf("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID environment variable is required for gov domains")
					}
				}

				authURL = fmt.Sprintf("https://%s/oauth2/%s/v1/token", oktaDomain, authServerID)
				useFormEncoded = true
				log.Printf("[DEBUG] Using Okta OAuth2 for gov domain: %s", authURL)
			} else {
				authURL = fmt.Sprintf("https://login.%s/oauth/token", domain)
			}
		}
	} else {
		return nil, fmt.Errorf("no base URL provided")
	}

	// Ensure audience is set
	if audience == "" {
		// Extract domain from workspace URL for audience
		parsedURL, err := url.Parse(o.baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL for fallback audience: %v", err)
		}

		host := parsedURL.Host

		// Handle test/localhost URLs differently
		if strings.Contains(host, "127.0.0.1") || strings.Contains(host, "localhost") {
			// For test URLs, use the same URL as audience
			audience = o.baseURL
		} else {
			// Extract domain part after the first dash (remove workspace-org prefix)
			parts := strings.SplitN(host, ".", 2)
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid workspace URL format for fallback audience: %s", host)
			}
			domain := parts[1] // e.g., "cribl.cloud"

			audience = fmt.Sprintf("https://api.%s", domain)
		}
	}

	var req *http.Request
	var err error

	if useFormEncoded {
		// For gov domains (Okta OAuth2), use form-encoded data
		formData := url.Values{}
		formData.Set("grant_type", "client_credentials")
		formData.Set("client_id", clientID)
		formData.Set("client_secret", clientSecret)
		formData.Set("audience", audience)

		req, err = http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(formData.Encode()))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		log.Printf("[DEBUG] Using form-encoded OAuth2 request for gov domain")
	} else {
		// Create JSON request body (matching bootstrap template format)
		requestBody := map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"audience":      audience,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err = http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(string(jsonData)))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	}

	body, err := o.doRequestWithRetry(req)
	if err != nil {
		return nil, err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return &TokenInfo{
		Token:     result.AccessToken,
		ExpiresAt: expiresAt,
	}, nil
}

func (o *CriblTerraformHook) AfterError(ctx AfterErrorContext, res *http.Response, err error) (*http.Response, error) {
	if res == nil {
		return res, err
	}

	switch res.StatusCode {
	// If we get an authentication error, try to handle it with our custom auth
	case http.StatusUnauthorized:
		// Get credentials from config or environment
		config, err := GetCredentials()
		if err != nil {
			return res, err
		}

		if config != nil {
			// Get audience from base URL
			audience := ""
			if o.baseURL != "" {
				// Extract domain from workspace URL for audience in error handler
				parsedURL, err := url.Parse(o.baseURL)
				if err != nil {
					return res, fmt.Errorf("failed to parse base URL for audience in error handler: %v", err)
				}

				host := parsedURL.Host

				// Handle test/localhost URLs differently
				if strings.Contains(host, "127.0.0.1") || strings.Contains(host, "localhost") {
					// For test URLs, use the same URL as audience
					audience = o.baseURL
				} else {
					// Extract domain part after the first dash (remove workspace-org prefix)
					parts := strings.SplitN(host, ".", 2)
					if len(parts) < 2 {
						return res, fmt.Errorf("invalid workspace URL format for audience in error handler: %s", host)
					}
					domain := parts[1] // e.g., "cribl-playground.cloud"

					audience = fmt.Sprintf("https://api.%s", domain)
				}
			} else if os.Getenv("CRIBL_AUDIENCE") != "" {
				audience = os.Getenv("CRIBL_AUDIENCE")
			} else {
				return res, fmt.Errorf("no base URL or audience provided")
			}

			// Get or create session
			sessionKey := fmt.Sprintf("%s:%s", config.ClientID, config.ClientSecret)
			var tokenInfo *TokenInfo

			if cachedTokenInfo, ok := o.sessions.Load(sessionKey); ok {
				tokenInfo = cachedTokenInfo.(*TokenInfo)
				if time.Until(tokenInfo.ExpiresAt) < 60*time.Minute {
					tokenInfo = nil
				}
			}

			if tokenInfo == nil {
				newTokenInfo, err := o.getBearerToken(ctx.Context, config.ClientID, config.ClientSecret, audience)
				if err != nil {
					return res, err
				}
				tokenInfo = newTokenInfo
				o.sessions.Store(sessionKey, tokenInfo)
			}

			// Update org and workspace IDs from config
			o.orgID = config.OrganizationID
			o.workspaceID = config.Workspace

			// Return a FailEarly error to stop other hooks from being called
			return res, &FailEarly{Cause: fmt.Errorf("authentication handled by custom hook")}
		}
	case http.StatusTooManyRequests:
		time.Sleep(1 * time.Second)
	}

	return res, err
}

func (o *CriblTerraformHook) doRequestWithRetry(req *http.Request) ([]byte, error) {
	var body []byte
	var resp *http.Response
	var err error

	success := false
	for i := range 3 {
		log.Printf("[DEBUG] http request attempt %d, doing query: %+v", i, req)

		resp, err = o.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request with retry: %v", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response with retry: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			success = true
			break
		} else if resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("[DEBUG] 429 during request, waiting to retry %d seconds", i)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}

	if !success {
		return nil, fmt.Errorf("failed to do request: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return body, nil
}
