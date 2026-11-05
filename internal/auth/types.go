package auth

import (
	"net/http"
	"time"
)

// CriblConfig holds Cribl auth and connection settings from the environment or credentials file.
type CriblConfig struct {
	ClientID       string `json:"client_id" ini:"client_id"`
	ClientSecret   string `json:"client_secret" ini:"client_secret"`
	OrganizationID string `json:"organization_id" ini:"organization_id"`
	Workspace      string `json:"workspace" ini:"workspace"`
	CloudDomain    string `json:"cloud_domain" ini:"cloud_domain"`

	OnpremServerURL string `json:"onprem_server_url" ini:"onprem_server_url"`
	OnpremUsername  string `json:"onprem_username" ini:"onprem_username"`
	OnpremPassword  string `json:"onprem_password" ini:"onprem_password"`

	// HTTPClient optionally supplies the authentication transport. It is runtime-only
	// and must not be changed after the config is used to fetch a token.
	HTTPClient *http.Client `json:"-" ini:"-"`

	// HTTPClientCacheScope opts into token sharing across HTTP clients with the
	// same immutable transport and TLS trust settings. Callers must derive it
	// from those settings, including CA contents, rather than a CA file path.
	// When empty, custom clients remain isolated by client identity.
	HTTPClientCacheScope string `json:"-" ini:"-"`
}

// TokenInfo holds a bearer token and its expiration time.
type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

// ConstructBaseURLInput holds inputs for constructing a Cribl workspace base URL.
type ConstructBaseURLInput struct {
	BaseURL             string
	ProviderOrgID       string
	ProviderWorkspaceID string
	ProviderCloudDomain string
}
