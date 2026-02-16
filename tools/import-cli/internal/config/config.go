// Package config loads Cribl auth and connection settings from flags, environment,
// and credentials file. Uses internal/sdk/credentials for file/env so SDK and CLI share one implementation.
package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/sdk/credentials"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds Cribl config loaded from flags, env, and credentials file.
type Config struct {
	v *viper.Viper
}

// NewConfig returns a Config that uses the given Viper for storage.
func NewConfig(v *viper.Viper) *Config {
	return &Config{v: v}
}

// Viper returns the underlying Viper (e.g. for client.NewFromViper).
func (c *Config) Viper() *viper.Viper {
	return c.v
}

// Key names used in Viper (and credentials file profile section).
const (
	KeyOnpremServerURL = "onprem_server_url"
	KeyOrganizationID  = "organization_id"
	KeyWorkspaceID     = "workspace_id"
	KeyCloudDomain     = "cloud_domain"
	KeyBearerToken     = "bearer_token"
	KeyClientID        = "client_id"
	KeyClientSecret    = "client_secret"
	KeyOnpremUsername  = "onprem_username"
	KeyOnpremPassword  = "onprem_password"
)

// Env variable names. Must match internal/sdk/internal/hooks (credentials.go, CriblTerraformHook.go)
const (
	EnvOnpremServerURL = "CRIBL_ONPREM_SERVER_URL"
	EnvBearerToken     = "CRIBL_BEARER_TOKEN"
	EnvClientID        = "CRIBL_CLIENT_ID"
	EnvClientSecret    = "CRIBL_CLIENT_SECRET"
	EnvOrganizationID  = "CRIBL_ORGANIZATION_ID"
	EnvWorkspaceID     = "CRIBL_WORKSPACE_ID"
	EnvCloudDomain     = "CRIBL_CLOUD_DOMAIN"
	EnvProfile         = "CRIBL_PROFILE"
	EnvOnpremUsername  = "CRIBL_ONPREM_USERNAME"
	EnvOnpremPassword  = "CRIBL_ONPREM_PASSWORD"
)

// LoadCredentialsFile merges credentials from internal/sdk/credentials (env + ~/.cribl/credentials)
// into the config. Only keys not already set (e.g. by flags) are filled.
func (c *Config) LoadCredentialsFile() error {
	creds, err := credentials.GetCredentials()
	if err != nil {
		return err
	}
	if creds == nil {
		return nil
	}
	setIfUnset(c.v, KeyClientID, creds.ClientID)
	setIfUnset(c.v, KeyClientSecret, creds.ClientSecret)
	setIfUnset(c.v, KeyOrganizationID, creds.OrganizationID)
	setIfUnset(c.v, KeyWorkspaceID, creds.Workspace)
	setIfUnset(c.v, KeyCloudDomain, creds.CloudDomain)
	setIfUnset(c.v, KeyOnpremServerURL, creds.OnpremServerURL)
	setIfUnset(c.v, KeyOnpremUsername, creds.OnpremUsername)
	setIfUnset(c.v, KeyOnpremPassword, creds.OnpremPassword)
	return nil
}

func setIfUnset(v *viper.Viper, key, value string) {
	if strings.TrimSpace(v.GetString(key)) != "" {
		return
	}
	if value != "" {
		v.Set(key, strings.TrimSpace(value))
	}
}

// BindEnv binds all supported env vars to the config's Viper keys.
func (c *Config) BindEnv() {
	_ = c.v.BindEnv(KeyOnpremServerURL, EnvOnpremServerURL)
	_ = c.v.BindEnv(KeyBearerToken, EnvBearerToken)
	_ = c.v.BindEnv(KeyClientID, EnvClientID)
	_ = c.v.BindEnv(KeyClientSecret, EnvClientSecret)
	_ = c.v.BindEnv(KeyOrganizationID, EnvOrganizationID)
	_ = c.v.BindEnv(KeyWorkspaceID, EnvWorkspaceID)
	_ = c.v.BindEnv(KeyCloudDomain, EnvCloudDomain)
	_ = c.v.BindEnv(KeyOnpremUsername, EnvOnpremUsername)
	_ = c.v.BindEnv(KeyOnpremPassword, EnvOnpremPassword)
}

// BindPFlag binds a flag to a config key (for Cobra flag override).
func (c *Config) BindPFlag(key string, flag *pflag.Flag) error {
	return c.v.BindPFlag(key, flag)
}

// Get returns the resolved value for key (flag > env > file).
func (c *Config) Get(key string) string {
	return strings.TrimSpace(c.v.GetString(key))
}

// validateOnPremServerURL returns an error if the given string is not a valid
// on-prem server URL (must be http or https with a non-empty host).
func validateOnPremServerURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid on-prem server URL %q: %w", s, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid on-prem server URL %q: scheme must be http or https", s)
	}
	if u.Host == "" {
		return fmt.Errorf("invalid on-prem server URL %q: missing host", s)
	}
	return nil
}

// ValidateRequired returns an actionable error if required configuration is missing.
// On-prem: need onprem_server_url and (bearer_token or username/password).
// Cloud: need (bearer_token or client_id/client_secret) and organization_id and workspace_id.
func (c *Config) ValidateRequired() error {
	serverURL := c.Get(KeyOnpremServerURL)
	token := c.Get(KeyBearerToken)
	clientID := c.Get(KeyClientID)
	clientSecret := c.Get(KeyClientSecret)
	orgID := c.Get(KeyOrganizationID)
	workspaceID := c.Get(KeyWorkspaceID)
	username := c.Get(KeyOnpremUsername)
	password := c.Get(KeyOnpremPassword)

	if serverURL != "" {
		if token == "" && (username == "" || password == "") {
			return fmt.Errorf("on-prem server URL is set but authentication is missing: set %s or set both %s and %s", EnvBearerToken, EnvOnpremUsername, EnvOnpremPassword)
		}
		if err := validateOnPremServerURL(serverURL); err != nil {
			return err
		}
		return nil
	}

	if token != "" {
		if orgID == "" || workspaceID == "" {
			return fmt.Errorf("cloud auth with token requires %s and %s", EnvOrganizationID, EnvWorkspaceID)
		}
		return nil
	}
	if clientID != "" && clientSecret != "" {
		if orgID == "" || workspaceID == "" {
			return fmt.Errorf("cloud auth with client credentials requires %s and %s", EnvOrganizationID, EnvWorkspaceID)
		}
		return nil
	}
	return fmt.Errorf("no valid configuration: set on-prem (%s and auth) or cloud (%s and %s, plus %s and %s)", EnvOnpremServerURL, EnvClientID, EnvClientSecret, EnvOrganizationID, EnvWorkspaceID)
}
