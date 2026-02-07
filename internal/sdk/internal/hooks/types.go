package hooks

import (
	"sync"
	"time"

	credpkg "github.com/criblio/terraform-provider-criblio/internal/sdk/credentials"
)

type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

type ConstructBaseUrlInput struct {
	BaseURL             string
	ProviderOrgID       string
	ProviderWorkspaceID string
	ProviderCloudDomain string
}

type LoadTokenInfoInput struct {
	Context    HookContext
	Config     *CriblConfig
	Audience   string
	SessionKey string
}

// CriblTerraformHook implements both authentication and URL routing for Cribl Terraform API
type CriblTerraformHook struct {
	client      HTTPClient
	sessions    sync.Map
	baseURL     string
	orgID       string
	workspaceID string
}

// CriblConfig is the shared credential config from internal/sdk/credentials.
type CriblConfig = credpkg.CriblConfig

type CriblConfigFile struct {
	Profiles map[string]CriblConfig `json:"profiles"`
}
