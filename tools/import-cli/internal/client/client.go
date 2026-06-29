// Package client builds the import CLI API clients from resolved config.
package client

import (
	"net/http"
	"os"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
)

// Client carries the REST client used by import-cli discovery and export paths.
type Client struct {
	REST *restclient.Client
}

// NewFromConfig builds Cribl API clients from the resolved config.
// Uses a capture transport for search list endpoints so the CLI can parse list
// responses when the SDK fails to unmarshal (e.g. cribl_lake union type).
func NewFromConfig(cfg *config.Config) (*Client, error) {
	applyConfigToEnv(cfg)
	transport := &custom.SearchListTransport{Base: http.DefaultTransport}
	httpClient := &http.Client{Transport: transport}
	userAgent := BulkExporterUserAgent()

	restClient := restclient.New(restclient.Config{
		Credentials: credentialsFromConfig(cfg),
		BearerToken: cfg.Get(config.KeyBearerToken),
		HTTPClient:  httpClient,
		UserAgent:   userAgent,
	})

	return &Client{REST: restClient}, nil
}

func credentialsFromConfig(cfg *config.Config) *auth.CriblConfig {
	return &auth.CriblConfig{
		ClientID:        cfg.Get(config.KeyClientID),
		ClientSecret:    cfg.Get(config.KeyClientSecret),
		OrganizationID:  cfg.Get(config.KeyOrganizationID),
		Workspace:       cfg.Get(config.KeyWorkspaceID),
		CloudDomain:     cfg.Get(config.KeyCloudDomain),
		OnpremServerURL: cfg.Get(config.KeyOnpremServerURL),
		OnpremUsername:  cfg.Get(config.KeyOnpremUsername),
		OnpremPassword:  cfg.Get(config.KeyOnpremPassword),
	}
}

// applyConfigToEnv sets CRIBL_* environment variables from the config.
func applyConfigToEnv(cfg *config.Config) {
	setEnv(config.EnvOnpremServerURL, cfg.Get(config.KeyOnpremServerURL))
	setEnv(config.EnvBearerToken, cfg.Get(config.KeyBearerToken))
	setEnv(config.EnvClientID, cfg.Get(config.KeyClientID))
	setEnv(config.EnvClientSecret, cfg.Get(config.KeyClientSecret))
	setEnv(config.EnvOrganizationID, cfg.Get(config.KeyOrganizationID))
	setEnv(config.EnvWorkspaceID, cfg.Get(config.KeyWorkspaceID))
	setEnv(config.EnvCloudDomain, cfg.Get(config.KeyCloudDomain))
	setEnv(config.EnvOnpremUsername, cfg.Get(config.KeyOnpremUsername))
	setEnv(config.EnvOnpremPassword, cfg.Get(config.KeyOnpremPassword))
}

func setEnv(key, value string) {
	if key == "" {
		return
	}
	if value == "" {
		_ = os.Unsetenv(key)
		return
	}
	_ = os.Setenv(key, value)
}
