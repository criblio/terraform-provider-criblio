// Package credentials provides shared Cribl credential loading from environment
// and ~/.cribl/credentials. Used by the SDK hooks and the import CLI.
package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// CriblConfig holds Cribl auth and connection settings (env or file).
type CriblConfig struct {
	ClientID       string `json:"client_id" ini:"client_id"`
	ClientSecret   string `json:"client_secret" ini:"client_secret"`
	OrganizationID string `json:"organization_id" ini:"organization_id"`
	Workspace      string `json:"workspace" ini:"workspace"`
	CloudDomain    string `json:"cloud_domain" ini:"cloud_domain"`

	OnpremServerURL string `json:"onprem_server_url" ini:"onprem_server_url"`
	OnpremUsername  string `json:"onprem_username" ini:"onprem_username"`
	OnpremPassword  string `json:"onprem_password" ini:"onprem_password"`
}

func checkLocalConfigDir() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[ERROR] Failed to get home directory: %v", err)
		return []byte{}, fmt.Errorf("failed to get home directory: %v", err)
	}
	configDir := filepath.Join(homeDir, ".cribl")
	configPath := filepath.Join(configDir, "credentials")
	var filePath string

	_, err = os.Stat(configPath)
	if err != nil {
		log.Printf("[DEBUG] No config file found %s", configPath)
		legacyPath := filepath.Join(homeDir, ".cribl")
		_, err := os.Stat(legacyPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Printf("[DEBUG] No config file found %s", legacyPath)
				return []byte{}, err
			}
			return []byte{}, err
		}
		filePath = legacyPath
	} else {
		filePath = configPath
	}

	log.Printf("[DEBUG] Reading credentials from: %s", filePath)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read credentials file: %v", err)
	}
	return file, nil
}

func checkConfigFileFormat(input []byte) (string, error) {
	var data interface{}
	if err := json.Unmarshal(input, &data); err == nil {
		return "json", nil
	}
	if _, err := ini.Load(input); err == nil {
		return "ini", nil
	}
	return "", errors.New("config file type not recognized")
}

func parseJSONConfig(file []byte) (*CriblConfig, error) {
	log.Printf("[DEBUG] parsing JSON config")
	var config CriblConfig
	if err := json.Unmarshal(file, &config); err != nil {
		log.Printf("[ERROR] Failed to parse config file: %v", err)
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	return &config, nil
}

func parseIniConfig(file []byte) (*CriblConfig, error) {
	config := CriblConfig{}
	profileName := os.Getenv("CRIBL_PROFILE")
	if profileName == "" {
		profileName = "default"
	}
	log.Printf("[DEBUG] Using profile: %s", profileName)

	cfg, err := ini.Load(file)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config file: %v", err)
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	if err := cfg.Section(profileName).MapTo(&config); err != nil {
		log.Printf("[ERROR] Failed to parse config file profile: %v", err)
		return nil, fmt.Errorf("failed to parse config file profile: %v", err)
	}
	log.Printf("[DEBUG] Selected profile values - clientID=%s, orgID=%s, workspace=%s, domain=%s",
		config.ClientID, config.OrganizationID, config.Workspace, config.CloudDomain)
	return &config, nil
}

// GetCredentials reads credentials from environment variables or ~/.cribl/credentials.
// Env takes precedence; file is used when env is not set.
func GetCredentials() (*CriblConfig, error) {
	clientID := os.Getenv("CRIBL_CLIENT_ID")
	clientSecret := os.Getenv("CRIBL_CLIENT_SECRET")
	organizationID := os.Getenv("CRIBL_ORGANIZATION_ID")
	workspace := os.Getenv("CRIBL_WORKSPACE_ID")
	cloudDomain := os.Getenv("CRIBL_CLOUD_DOMAIN")
	onpremServerURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	onpremUsername := os.Getenv("CRIBL_ONPREM_USERNAME")
	onpremPassword := os.Getenv("CRIBL_ONPREM_PASSWORD")

	log.Printf("[DEBUG] Environment variables - clientID=%s, orgID=%s, workspace=%s, domain=%s, onpremURL=%s",
		clientID, organizationID, workspace, cloudDomain, onpremServerURL)

	if clientID != "" && clientSecret != "" {
		log.Printf("[DEBUG] Using cloud credentials from environment variables")
		return &CriblConfig{
			ClientID:       clientID,
			ClientSecret:   clientSecret,
			OrganizationID: organizationID,
			Workspace:      workspace,
			CloudDomain:    cloudDomain,
		}, nil
	}
	if onpremServerURL != "" {
		log.Printf("[DEBUG] Using on-prem credentials from environment variables")
		return &CriblConfig{
			OnpremServerURL: onpremServerURL,
			OnpremUsername:  onpremUsername,
			OnpremPassword:  onpremPassword,
		}, nil
	}

	file, err := checkLocalConfigDir()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("[DEBUG] No configuration file found, continuing")
			return &CriblConfig{
				ClientID:       clientID,
				ClientSecret:   clientSecret,
				OrganizationID: organizationID,
				Workspace:      workspace,
				CloudDomain:    cloudDomain,
			}, nil
		}
		return nil, err
	}

	format, err := checkConfigFileFormat(file)
	if err != nil {
		log.Printf("[DEBUG] No configuration file found, continuing - error: %v", err)
		return nil, err
	}
	switch format {
	case "json":
		return parseJSONConfig(file)
	case "ini":
		return parseIniConfig(file)
	}
	return nil, nil
}
