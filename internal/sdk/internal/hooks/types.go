package hooks

import (
	"time"
	"sync"

	"golang.org/x/sync/singleflight"
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

type CriblConfig struct {
	ClientID       string `json:"client_id" ini:"client_id"`
	ClientSecret   string `json:"client_secret" ini:"client_secret"`
	OrganizationID string `json:"organization_id" ini:"organization_id"`
	Workspace      string `json:"workspace" ini:"workspace"`
	CloudDomain    string `json:"cloud_domain" ini:"cloud_domain"`

	// On-prem configuration fields
	OnpremServerURL string `json:"onprem_server_url" ini:"onprem_server_url"`
	OnpremUsername  string `json:"onprem_username" ini:"onprem_username"`
	OnpremPassword  string `json:"onprem_password" ini:"onprem_password"`
}

type CriblConfigFile struct {
	Profiles map[string]CriblConfig `json:"profiles"`
}
