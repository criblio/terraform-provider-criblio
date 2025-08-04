package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

type CriblConfig struct {
	ClientID       string `json:"client_id" ini:"client_id"`
	ClientSecret   string `json:"client_secret" ini:"client_secret"`
	OrganizationID string `json:"organization_id" ini:"organization_id"`
	Workspace      string `json:"workspace" ini:"workspace"`
}

type CriblConfigFile struct {
	Profiles map[string]CriblConfig `json:"profiles"`
}

func checkLocalConfigDir() ([]byte, error) {
	// Then try ~/.cribl/credentials file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[ERROR] Failed to get home directory: %v", err)
		return []byte{}, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Check for credentials in ~/.cribl/credentials
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
			} else {
				return []byte{}, err
			}
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
	err := json.Unmarshal(input, &data)
	if err == nil {
		return "json", nil
	}

	_, err = ini.Load(input)
	if err == nil {
		return "ini", nil
	}

	return "", errors.New("config file type not recognized")
}

func parseJSONConfig(file []byte) (*CriblConfig, error) {
	log.Printf("[DEBUG] parsing JSON config")

	var legacyConfig CriblConfig
	err := json.Unmarshal(file, &legacyConfig)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config file: %v", err)
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &legacyConfig, nil
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

	err = cfg.Section(profileName).MapTo(&config)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config file profile: %v", err)
		return nil, fmt.Errorf("failed to parse config file profile: %v", err)
	}

	log.Printf("[DEBUG] Selected profile values - clientID=%s, orgID=%s, workspace=%s",
		config.ClientID, config.OrganizationID, config.Workspace)
	return &config, nil
}

// GetCredentials reads credentials from environment variables or credentials file
func GetCredentials() (*CriblConfig, error) {
	// First try environment variables
	clientID := os.Getenv("CRIBL_CLIENT_ID")
	clientSecret := os.Getenv("CRIBL_CLIENT_SECRET")
	organizationID := os.Getenv("CRIBL_ORGANIZATION_ID")
	workspace := os.Getenv("CRIBL_WORKSPACE_ID")

	log.Printf("[DEBUG] Environment variables - clientID=%s, orgID=%s, workspace=%s",
		clientID, organizationID, workspace)

	// If we have direct credentials in environment, use them
	if clientID != "" && clientSecret != "" {
		log.Printf("[DEBUG] Using credentials from environment variables")
		return &CriblConfig{
			ClientID:       clientID,
			ClientSecret:   clientSecret,
			OrganizationID: organizationID,
			Workspace:      workspace,
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
			}, nil
		} else {
			return nil, err
		}
	}

	format, err := checkConfigFileFormat(file)
	if err != nil {
		log.Printf("[DEBUG] No configuration file found, continuing")
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
