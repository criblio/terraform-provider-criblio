// Package config loads Cribl auth and connection settings from flags, environment,
// and ~/.cribl/credentials. Env var names and credentials file format match the SDK's CriblTerraformHook.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

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

// LoadCredentialsFile reads ~/.cribl/credentials (or legacy ~/.cribl) and sets Viper keys
// for the profile given by CRIBL_PROFILE (default "default"). File format is INI.
func LoadCredentialsFile(v *viper.Viper) error {
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
	cfg, err := ini.Load(data)
	if err != nil {
		return fmt.Errorf("failed to parse credentials file %s: %w", path, err)
	}
	profile := os.Getenv(EnvProfile)
	if profile == "" {
		profile = "default"
	}
	sec, err := cfg.GetSection(profile)
	if err != nil {
		return nil
	}
	setFromIniIfUnset(v, sec, "client_id", KeyClientID)
	setFromIniIfUnset(v, sec, "client_secret", KeyClientSecret)
	setFromIniIfUnset(v, sec, "organization_id", KeyOrganizationID)
	setFromIniIfUnset(v, sec, "workspace", KeyWorkspaceID)
	setFromIniIfUnset(v, sec, "cloud_domain", KeyCloudDomain)
	setFromIniIfUnset(v, sec, "onprem_server_url", KeyOnpremServerURL)
	setFromIniIfUnset(v, sec, "onprem_username", KeyOnpremUsername)
	setFromIniIfUnset(v, sec, "onprem_password", KeyOnpremPassword)
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

// BindEnv binds all supported env vars to Viper keys.
func BindEnv(v *viper.Viper) {
	_ = v.BindEnv(KeyOnpremServerURL, EnvOnpremServerURL)
	_ = v.BindEnv(KeyBearerToken, EnvBearerToken)
	_ = v.BindEnv(KeyClientID, EnvClientID)
	_ = v.BindEnv(KeyClientSecret, EnvClientSecret)
	_ = v.BindEnv(KeyOrganizationID, EnvOrganizationID)
	_ = v.BindEnv(KeyWorkspaceID, EnvWorkspaceID)
	_ = v.BindEnv(KeyCloudDomain, EnvCloudDomain)
	_ = v.BindEnv(KeyOnpremUsername, EnvOnpremUsername)
	_ = v.BindEnv(KeyOnpremPassword, EnvOnpremPassword)
}

// Get returns the resolved value for key (flag > env > file).
func Get(v *viper.Viper, key string) string {
	return strings.TrimSpace(v.GetString(key))
}

// ValidateRequired returns an actionable error if required configuration is missing.
// On-prem: need onprem_server_url and (bearer_token or username/password).
// Cloud: need (bearer_token or client_id/client_secret) and organization_id and workspace_id.
func ValidateRequired(v *viper.Viper) error {
	serverURL := Get(v, KeyOnpremServerURL)
	token := Get(v, KeyBearerToken)
	clientID := Get(v, KeyClientID)
	clientSecret := Get(v, KeyClientSecret)
	orgID := Get(v, KeyOrganizationID)
	workspaceID := Get(v, KeyWorkspaceID)
	username := Get(v, KeyOnpremUsername)
	password := Get(v, KeyOnpremPassword)

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
