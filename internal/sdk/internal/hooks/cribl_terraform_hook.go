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

	// Get credentials from config or environment
	log.Printf("[DEBUG] Attempting to get credentials")
	config, err := GetCredentials()
	if err != nil {
		log.Printf("[ERROR] Failed to get credentials: %v", err)
		o.baseURL = baseURL
		return baseURL, client
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
	// If we have a concrete URL (not a template), use it as-is
	if baseURL != "" && !strings.Contains(baseURL, "{workspaceName}") && !strings.Contains(baseURL, "{organizationId}") {
		return baseURL
	}

	// Get values with proper precedence: Environment > Config > Default
	workspace := os.Getenv("CRIBL_WORKSPACE_ID")
	if workspace == "" {
		workspace = config.Workspace
	}
	if workspace == "" {
		workspace = "main" // Default workspace name
	}

	organizationID := os.Getenv("CRIBL_ORGANIZATION_ID")
	if organizationID == "" {
		organizationID = config.OrganizationID
	}
	if organizationID == "" {
		organizationID = "ian" // Default organization ID
	}

	// Get cloud domain with proper precedence: Environment > Config > Default
	cloudDomain := os.Getenv("CRIBL_CLOUD_DOMAIN")
	if cloudDomain == "" {
		cloudDomain = config.CloudDomain
	}
	if cloudDomain == "" {
		cloudDomain = "cribl.cloud" // Default domain
	}

	// Construct the workspace URL: https://{workspace}-{organizationId}.{cloudDomain}
	constructedURL := fmt.Sprintf("https://%s-%s.%s", workspace, organizationID, cloudDomain)
	log.Printf("[DEBUG] Constructed URL: %s from workspace=%s, org=%s, domain=%s",
		constructedURL, workspace, organizationID, cloudDomain)

	return constructedURL
}

func (o *CriblTerraformHook) BeforeRequest(ctx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	// First try to get credentials from security context
	var clientID, clientSecret, orgID, workspaceID string

	if ctx.SecuritySource != nil {
		if security, err := ctx.SecuritySource(ctx.Context); err == nil {
			if s, ok := security.(shared.Security); ok {
				// Get OAuth credentials
				if s.ClientOauth != nil {
					clientID = s.ClientOauth.ClientID
					clientSecret = s.ClientOauth.ClientSecret
				}

				// Get org and workspace IDs
				if s.OrganizationID != nil {
					orgID = *s.OrganizationID
					o.orgID = orgID
				}
				if s.WorkspaceID != nil {
					workspaceID = *s.WorkspaceID
					o.workspaceID = workspaceID
				}
			}
		} else {
			log.Printf("[ERROR] Failed to get security info: %v", err)
		}
	}

	// If we don't have credentials from security context, try to get them from config file
	if clientID == "" || clientSecret == "" {
		config, err := GetCredentials()
		if err != nil {
			log.Printf("[ERROR] Failed to get credentials from config: %v", err)
		} else if config != nil {
			clientID = config.ClientID
			clientSecret = config.ClientSecret
			if orgID == "" {
				orgID = config.OrganizationID
				o.orgID = orgID
			}
			if workspaceID == "" {
				workspaceID = config.Workspace
				o.workspaceID = workspaceID
			}
		}
	}

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

	// Handle URL routing - for new workspace format, requests go directly to workspace URL
	trimmedBaseURL := strings.TrimRight(o.baseURL, "/")
	path := strings.TrimLeft(req.URL.Path, "/")

	// For new workspace format, construct API URLs directly (e.g., https://main-org.domain/api/v1/...)
	if !strings.Contains(req.URL.String(), "/api/v1") && !strings.HasPrefix(path, "api/v1") {
		newURL := fmt.Sprintf("%s/api/v1/%s", trimmedBaseURL, path)

		parsedURL, err := url.Parse(newURL)
		if err != nil {
			return req, err
		}

		req.URL = parsedURL
	}

	return req, nil
}

func (o *CriblTerraformHook) getBearerToken(ctx context.Context, clientID, clientSecret, audience string) (*TokenInfo, error) {
	// Get auth URL from base URL
	authURL := ""
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

			authURL = fmt.Sprintf("https://login.%s/oauth/token", domain)
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

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get token: %s", string(body))
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
	// If we get an authentication error, try to handle it with our custom auth
	if res != nil && res.StatusCode == http.StatusUnauthorized {
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
	}
	return res, err
}
