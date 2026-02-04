// Package client builds the Cribl SDK client from Viper config (flags, env, credentials file).
// It applies Viper's resolved config to CRIBL_* environment variables so the SDK's CriblTerraformHook handles authentication.
package client

import (
	"net/http"
	"os"

	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/spf13/viper"
)

// NewFromViper builds a Cribl SDK client from the resolved config.
func NewFromViper(v *viper.Viper) (*sdk.CriblIo, error) {
	applyViperToEnv(v)
	return sdk.New(sdk.WithClient(&http.Client{})), nil
}

// applyViperToEnv sets CRIBL_* environment variables from Viper's resolved config.
func applyViperToEnv(v *viper.Viper) {
	setEnv(config.EnvOnpremServerURL, config.Get(v, config.KeyOnpremServerURL))
	setEnv(config.EnvBearerToken, config.Get(v, config.KeyBearerToken))
	setEnv(config.EnvClientID, config.Get(v, config.KeyClientID))
	setEnv(config.EnvClientSecret, config.Get(v, config.KeyClientSecret))
	setEnv(config.EnvOrganizationID, config.Get(v, config.KeyOrganizationID))
	setEnv(config.EnvWorkspaceID, config.Get(v, config.KeyWorkspaceID))
	setEnv(config.EnvCloudDomain, config.Get(v, config.KeyCloudDomain))
	setEnv(config.EnvOnpremUsername, config.Get(v, config.KeyOnpremUsername))
	setEnv(config.EnvOnpremPassword, config.Get(v, config.KeyOnpremPassword))
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
