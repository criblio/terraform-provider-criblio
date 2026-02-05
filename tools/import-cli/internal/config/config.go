// Package config loads Cribl auth and connection settings from flags, environment,
// and ~/.cribl/credentials. Env var names and credentials file format match the SDK's CriblTerraformHook.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
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

// DefaultCredentialsPath returns ~/.cribl/credentials.
func DefaultCredentialsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".cribl", "credentials"), nil
}

// LoadCredentialsFile reads ~/.cribl/credentials (or legacy ~/.cribl) and sets keys
// for the profile given by CRIBL_PROFILE (default "default"). File format is INI.
func (c *Config) LoadCredentialsFile() error {
	path, err := DefaultCredentialsPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		legacy := filepath.Join(filepath.Dir(filepath.Dir(path)), ".cribl")
		if _, err := os.Stat(legacy); os.IsNotExist(err) {
			return nil
		}
		path = legacy
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read credentials file %s: %w", path, err)
	}
	iniCfg, err := ini.Load(data)
	if err != nil {
		return fmt.Errorf("failed to parse credentials file %s: %w", path, err)
	}
	profile := os.Getenv(EnvProfile)
	if profile == "" {
		profile = "default"
	}
	sec, err := iniCfg.GetSection(profile)
	if err != nil {
		return nil
	}
	setFromIniIfUnset(c.v, sec, "client_id", KeyClientID)
	setFromIniIfUnset(c.v, sec, "client_secret", KeyClientSecret)
	setFromIniIfUnset(c.v, sec, "organization_id", KeyOrganizationID)
	setFromIniIfUnset(c.v, sec, "workspace", KeyWorkspaceID)
	setFromIniIfUnset(c.v, sec, "cloud_domain", KeyCloudDomain)
	setFromIniIfUnset(c.v, sec, "onprem_server_url", KeyOnpremServerURL)
	setFromIniIfUnset(c.v, sec, "onprem_username", KeyOnpremUsername)
	setFromIniIfUnset(c.v, sec, "onprem_password", KeyOnpremPassword)
	return nil
}

func setFromIniIfUnset(v *viper.Viper, sec *ini.Section, iniKey, viperKey string) {
	if v.GetString(viperKey) != "" {
		return
	}
	k, err := sec.GetKey(iniKey)
	if err != nil {
		return
	}
	v.Set(viperKey, strings.TrimSpace(k.String()))
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
