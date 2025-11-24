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
	"sync"
	"time"
)

type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

// CriblTerraformHook implements both authentication and URL routing for Cribl Terraform API
type CriblTerraformHook struct {
	client      HTTPClient
	sessions    sync.Map
	baseURL     string
	orgID       string
	workspaceID string
}

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
		log.Printf("[DEBUG] On-prem configuration detected in environment, using server URL: %s", onpremServerURL)
		o.baseURL = onpremServerURL
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
		log.Printf("[DEBUG] On-prem configuration detected in credentials file, using server URL: %s", config.OnpremServerURL)
		o.baseURL = config.OnpremServerURL
		return config.OnpremServerURL, client
	}

	// Set orgID and workspaceID from config
	if config != nil {
		log.Printf("[DEBUG] Setting orgID: %s and workspaceID: %s", config.OrganizationID, config.Workspace)
		o.orgID = config.OrganizationID
		o.workspaceID = config.Workspace

		// If baseURL is not provided or is a template, construct it from credentials
		finalBaseURL := o.constructBaseURL(baseURL, config)
		o.baseURL = finalBaseURL
		log.Printf("[DEBUG] Final baseURL: %s", finalBaseURL)
		return finalBaseURL, client
	} else {
		log.Printf("[DEBUG] No credentials found")
		o.baseURL = baseURL
	}

	log.Printf("[DEBUG] Initialization complete")
	return baseURL, client
}

// constructBaseURL builds the workspace URL from credentials when needed
func (o *CriblTerraformHook) constructBaseURL(baseURL string, config *CriblConfig) string {
	// Always check if environment variables are set - they take highest priority
	// If any environment variables are set, reconstruct the URL from components
	workspaceEnv := os.Getenv("CRIBL_WORKSPACE_ID")
	orgEnv := os.Getenv("CRIBL_ORGANIZATION_ID")
	domainEnv := os.Getenv("CRIBL_CLOUD_DOMAIN")

	// Special case: if we have a localhost/test URL, keep it even with environment variables
	if baseURL != "" && (strings.Contains(baseURL, "127.0.0.1") || strings.Contains(baseURL, "localhost")) {
		log.Printf("[DEBUG] Localhost URL detected, keeping as-is: %s", baseURL)
		return baseURL
	}

	// If no environment variables are set and we have a concrete URL, use it as-is
	if workspaceEnv == "" && orgEnv == "" && domainEnv == "" &&
		baseURL != "" && !strings.Contains(baseURL, "{workspaceName}") && !strings.Contains(baseURL, "{organizationId}") {
		log.Printf("[DEBUG] No environment variables set, using provided concrete URL: %s", baseURL)
		return baseURL
	}

	// Get values with proper precedence: Environment > Config > Default
	workspace := workspaceEnv
	workspaceSource := "environment"
	if workspace == "" {
		workspace = config.Workspace
		workspaceSource = "config"
	}
	if workspace == "" {
		workspace = "main" // Default workspace name
		workspaceSource = "default"
	}
	log.Printf("[DEBUG] Workspace selection: env='%s', config='%s', final='%s', source='%s'",
		workspaceEnv, config.Workspace, workspace, workspaceSource)

	organizationID := orgEnv
	orgSource := "environment"
	if organizationID == "" {
		organizationID = config.OrganizationID
		orgSource = "config"
	}
	if organizationID == "" {
		organizationID = "ian" // Default organization ID
		orgSource = "default"
	}
	log.Printf("[DEBUG] Organization selection: env='%s', config='%s', final='%s', source='%s'",
		orgEnv, config.OrganizationID, organizationID, orgSource)

	// Get cloud domain with proper precedence: Environment > Config > Default
	cloudDomain := domainEnv
	domainSource := "environment"
	if cloudDomain == "" {
		cloudDomain = config.CloudDomain
		domainSource = "config"
	}
	if cloudDomain == "" {
		cloudDomain = "cribl.cloud" // Default domain
		domainSource = "default"
	}
	log.Printf("[DEBUG] Domain selection: env='%s', config='%s', final='%s', source='%s'",
		domainEnv, config.CloudDomain, cloudDomain, domainSource)

	// Construct the workspace URL: https://{workspace}-{organizationId}.{cloudDomain}
	constructedURL := fmt.Sprintf("https://%s-%s.%s", workspace, organizationID, cloudDomain)
	log.Printf("[DEBUG] Constructed URL: %s from workspace=%s, org=%s, domain=%s",
		constructedURL, workspace, organizationID, cloudDomain)

	return constructedURL
}

// constructBaseURLWithProviderConfig builds URL with proper precedence: Provider Config > Environment > Credentials File > Default
func (o *CriblTerraformHook) constructBaseURLWithProviderConfig(providerOrgID, providerWorkspaceID, providerCloudDomain string, config *CriblConfig) string {
	// Special case: if we have a localhost/test URL, keep it even with provider config
	if o.baseURL != "" && (strings.Contains(o.baseURL, "127.0.0.1") || strings.Contains(o.baseURL, "localhost")) {
		log.Printf("[DEBUG] Localhost URL detected, keeping as-is: %s", o.baseURL)
		return o.baseURL
	}

	// Get environment variables (second precedence)
	workspaceEnv := os.Getenv("CRIBL_WORKSPACE_ID")
	orgEnv := os.Getenv("CRIBL_ORGANIZATION_ID")
	domainEnv := os.Getenv("CRIBL_CLOUD_DOMAIN")

	// Apply precedence: Provider Config > Environment > Credentials File > Default
	workspace := providerWorkspaceID
	workspaceSource := "provider"
	if workspace == "" {
		workspace = workspaceEnv
		workspaceSource = "environment"
	}
	if workspace == "" && config != nil {
		workspace = config.Workspace
		workspaceSource = "credentials"
	}
	if workspace == "" {
		workspace = "main" // Default workspace name
		workspaceSource = "default"
	}

	organizationID := providerOrgID
	orgSource := "provider"
	if organizationID == "" {
		organizationID = orgEnv
		orgSource = "environment"
	}
	if organizationID == "" && config != nil {
		organizationID = config.OrganizationID
		orgSource = "credentials"
	}
	if organizationID == "" {
		organizationID = "ian" // Default organization ID
		orgSource = "default"
	}

	cloudDomain := providerCloudDomain
	domainSource := "provider"
	if cloudDomain == "" {
		cloudDomain = domainEnv
		domainSource = "environment"
	}
	if cloudDomain == "" && config != nil {
		cloudDomain = config.CloudDomain
		domainSource = "credentials"
	}
	if cloudDomain == "" {
		cloudDomain = "cribl.cloud" // Default domain
		domainSource = "default"
	}

	log.Printf("[DEBUG] Final precedence - Workspace: '%s' (from %s), Org: '%s' (from %s), Domain: '%s' (from %s)",
		workspace, workspaceSource, organizationID, orgSource, cloudDomain, domainSource)

	// Construct the workspace URL: https://{workspace}-{organizationId}.{cloudDomain}
	constructedURL := fmt.Sprintf("https://%s-%s.%s", workspace, organizationID, cloudDomain)
	log.Printf("[DEBUG] Final constructed URL: %s", constructedURL)

	return constructedURL
}

// isGatewayPath determines if a request path should be routed to the gateway
func (o *CriblTerraformHook) isGatewayPath(path string) bool {
	// Remove leading slash and api prefix for consistent checking
	cleanPath := strings.TrimLeft(path, "/")
	cleanPath = strings.TrimPrefix(cleanPath, "api/")
	cleanPath = strings.TrimPrefix(cleanPath, "v1/")

	// Gateway paths are for organization and workspace management
	return strings.HasPrefix(cleanPath, "organizations/")
}

// isRestrictedOnPremEndpoint determines if a path is for a restricted endpoint that is not supported on on-prem deployments
func (o *CriblTerraformHook) isRestrictedOnPremEndpoint(path string) bool {
	// These endpoints are only available in Cribl.Cloud (gateway) and not in on-prem deployments

	// Check for search endpoints (can be in format /m/group/search/... or search/...)
	// Search endpoints contain "search/" somewhere in the path
	if strings.Contains(path, "search/") {
		// Exclude /system/config-search from restrictions as it's an admin endpoint
		if !strings.Contains(path, "/system/config-search") {
			return true
		}
	}

	// Check for lake endpoints
	if strings.Contains(path, "/lake/") || strings.Contains(path, "products/lake/") {
		return true
	}

	// Check for lakehouse endpoints
	if strings.Contains(path, "lakehouse") {
		return true
	}

	// Check for products/search endpoints
	if strings.Contains(path, "products/search/") {
		return true
	}

	// Check for gateway/management endpoints
	restrictedPrefixes := []string{
		"v1/organizations/", // Gateway management endpoints
		"organizations/",    // Workspace management endpoints
	}

	// Check for restricted prefixes
	for _, prefix := range restrictedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Check for workspace management (these are gateway paths)
	if strings.Contains(path, "workspace") && !strings.Contains(path, "/system/") {
		return true
	}

	return false
}

// constructGatewayURL builds the gateway URL using the appropriate cloud domain
func (o *CriblTerraformHook) constructGatewayURL(providerCloudDomain string, config *CriblConfig) string {
	// Get cloud domain with proper precedence: Provider > Environment > Credentials > Default
	cloudDomain := providerCloudDomain
	if cloudDomain == "" {
		cloudDomain = os.Getenv("CRIBL_CLOUD_DOMAIN")
	}
	if cloudDomain == "" && config != nil {
		cloudDomain = config.CloudDomain
	}
	if cloudDomain == "" {
		cloudDomain = "cribl.cloud" // Default domain
	}

	return fmt.Sprintf("https://gateway.%s", cloudDomain)
}

// handleOnPremRequest handles authentication and routing for on-prem deployments
func (o *CriblTerraformHook) handleOnPremRequest(ctx BeforeRequestContext, req *http.Request, serverURL string) (*http.Request, error) {
	log.Printf("[DEBUG] Handling on-prem request to: %s", serverURL)

	// Get username and password from environment
	username := os.Getenv("CRIBL_ONPREM_USERNAME")
	password := os.Getenv("CRIBL_ONPREM_PASSWORD")
	bearerToken := os.Getenv("CRIBL_BEARER_TOKEN")

	// Get credentials from config file as fallback
	config, err := GetCredentials()
	if err != nil {
		log.Printf("[WARN] Failed to get credentials from config: %v", err)
	}

	// Use config as fallback if not in environment
	if username == "" && config != nil {
		username = config.OnpremUsername
	}
	if password == "" && config != nil {
		password = config.OnpremPassword
	}

	var authToken string

	// Try direct bearer token first
	if bearerToken != "" {
		log.Printf("[DEBUG] Using direct bearer token for on-prem authentication")
		authToken = bearerToken
	} else if username != "" && password != "" {
		// Check for cached token
		sessionKey := fmt.Sprintf("onprem:%s:%s:%s", serverURL, username, password)
		var tokenInfo *TokenInfo

		if cachedTokenInfo, ok := o.sessions.Load(sessionKey); ok {
			tokenInfo = cachedTokenInfo.(*TokenInfo)
			if time.Until(tokenInfo.ExpiresAt) < 30*time.Minute {
				tokenInfo = nil
			}
		}

		if tokenInfo == nil {
			log.Printf("[DEBUG] Retrieving bearer token using username/password for on-prem authentication")
			newTokenInfo, err := o.getOnPremBearerToken(ctx.Context, serverURL, username, password)
			if err != nil {
				return req, fmt.Errorf("failed to get on-prem bearer token: %v", err)
			}
			tokenInfo = newTokenInfo
			o.sessions.Store(sessionKey, tokenInfo)
		}

		authToken = tokenInfo.Token
	} else {
		return req, fmt.Errorf("on-prem authentication requires either CRIBL_BEARER_TOKEN or both CRIBL_ONPREM_USERNAME and CRIBL_ONPREM_PASSWORD")
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+authToken)

	// Validate that the requested endpoint is supported for on-prem deployments
	path := strings.TrimLeft(req.URL.Path, "/")

	// Remove /api/v1 if already present in path
	path = strings.TrimPrefix(path, "api/v1/")
	path = strings.TrimPrefix(path, "api/v1")

	// Check if this is a restricted endpoint for on-prem
	if o.isRestrictedOnPremEndpoint(path) {
		return req, fmt.Errorf("endpoint '%s' is not supported for on-prem deployments. On-prem deployments only support workspace resources (sources, destinations, routes, pipelines, etc.)", path)
	}

	// Handle URL routing for on-prem (always use serverURL/api/v1 path)
	baseURL := strings.TrimRight(serverURL, "/")

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

	var body []byte
	var resp *http.Response
	success := false
	for i := range 3 {
		resp, err = o.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			success = true
			break
		} else if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("[DEBUG] 429 getting on-prem bearer token, waiting to retry %d seconds", i)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}

	if !success {
		return nil, fmt.Errorf("failed to get token: status=%d, body=%s", resp.StatusCode, string(body))
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

func (o *CriblTerraformHook) BeforeRequest(ctx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	// Check for on-prem configuration first
	onpremServerURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")

	// Handle on-prem authentication
	if onpremServerURL != "" {
		return o.handleOnPremRequest(ctx, req, onpremServerURL)
	}

	// First try to get credentials from security context
	var clientID, clientSecret, orgID, workspaceID, cloudDomain string

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

	// Get credentials file config for fallback values
	config, err := GetCredentials()
	if err != nil {
		log.Printf("[ERROR] Failed to get credentials from config: %v", err)
	}

	// If we don't have credentials from security context, use config file
	if clientID == "" || clientSecret == "" {
		if config != nil {
			clientID = config.ClientID
			clientSecret = config.ClientSecret
		}
	}

	// Reconstruct baseURL with proper precedence: Provider Config > Environment > Credentials File > Default
	o.baseURL = o.constructBaseURLWithProviderConfig(orgID, workspaceID, cloudDomain, config)

	// Handle authentication
	if bearerToken := os.Getenv("CRIBL_BEARER_TOKEN"); bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	} else if clientID != "" && clientSecret != "" {
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

		req.Header.Set("Authorization", "Bearer "+tokenInfo.Token)
	}

	// Handle URL routing
	path := strings.TrimLeft(req.URL.Path, "/")

	// Check if this is a gateway path (management endpoints)
	if o.isGatewayPath(path) || strings.Contains(req.URL.Host, "gateway.") {
		// Construct gateway URL
		gatewayURL := o.constructGatewayURL(cloudDomain, config)

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

	return req, nil
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

	var body []byte
	var resp *http.Response
	success := false
	for i := range 3 {
		resp, err = o.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			success = true
			break
		} else if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("[DEBUG] 429 getting on-prem bearer token, waiting to retry %d seconds", i)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}

	if !success {
		return nil, fmt.Errorf("failed to get token: status=%d, body=%s", resp.StatusCode, string(body))
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

	switch {
	// If we get an authentication error, try to handle it with our custom auth
	case res.StatusCode == http.StatusUnauthorized:
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
