// Package client builds the Cribl SDK client from Viper config (flags, env, credentials file).
// It applies Viper's resolved config to CRIBL_* environment variables so the SDK's CriblTerraformHook handles authentication.
package client

import (
	"net/http"
	"os"

	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
)

// NewFromConfig builds a Cribl SDK client from the resolved config.
// Uses a capture transport for search list endpoints so the CLI can parse list
// responses when the SDK fails to unmarshal (e.g. cribl_lake union type).
func NewFromConfig(cfg *config.Config) (*sdk.CriblIo, error) {
	applyConfigToEnv(cfg)
	transport := &custom.SearchListTransport{Base: http.DefaultTransport}
	return sdk.New(sdk.WithClient(&http.Client{Transport: transport})), nil
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
