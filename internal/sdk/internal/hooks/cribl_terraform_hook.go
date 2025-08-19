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
	o.baseURL = baseURL

	// Get credentials from config or environment
	log.Printf("[DEBUG] Attempting to get credentials")
	config, err := GetCredentials()
	if err != nil {
		log.Printf("[ERROR] Failed to get credentials: %v", err)
		return baseURL, client
	}

	if config != nil {
		log.Printf("[DEBUG] Setting orgID: %s and workspaceID: %s", config.OrganizationID, config.Workspace)
		o.orgID = config.OrganizationID
		o.workspaceID = config.Workspace
	} else {
		log.Printf("[DEBUG] No credentials found")
	}

	log.Printf("[DEBUG] Initialization complete")
	return baseURL, client
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
			audience = strings.Replace(o.baseURL, "app.", "api.", 1)
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
	serverURL := strings.TrimRight(o.baseURL, "/")
	path := strings.TrimLeft(req.URL.Path, "/")

	// If we have org and workspace IDs from security context or config, use them
	if orgID != "" && workspaceID != "" && !strings.Contains(req.URL.String(), "app/api/v1") {
		newURL := fmt.Sprintf("%s/organizations/%s/workspaces/%s/app/api/v1/%s",
			serverURL,
			orgID,
			workspaceID,
			path,
		)

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
		authURL = strings.Replace(o.baseURL, "app.", "login.", 1)
		authURL = strings.TrimSuffix(authURL, "/")
		authURL = authURL + "/oauth/token"
	} else {
		return nil, fmt.Errorf("no base URL provided")
	}

	// Create form data
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)

	// Ensure audience is set
	if audience == "" {
		audience = strings.Replace(o.baseURL, "app.", "api.", 1)
	}
	formData.Set("audience", audience)

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
				audience = strings.Replace(o.baseURL, "app.", "api.", 1)
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
