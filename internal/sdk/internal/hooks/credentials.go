package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CriblConfig struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	OrganizationID string `json:"organization_id"`
	Workspace      string `json:"workspace"`
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
	// Try to parse as INI format first
	config := &CriblConfigFile{
		Profiles: make(map[string]CriblConfig),
	}

	// Read the file line by line
	lines := strings.Split(string(file), "\n")
	currentProfile := "default"

	log.Printf("[DEBUG] Parsing credentials file, starting with profile: %s", currentProfile)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentProfile = strings.Trim(line, "[]")
			log.Printf("[DEBUG] Found profile section: %s", currentProfile)
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		profile := config.Profiles[currentProfile]
		switch key {
		case "client_id":
			profile.ClientID = value
		case "client_secret":
			profile.ClientSecret = value
		case "organization_id":
			profile.OrganizationID = value
			log.Printf("[DEBUG] Set organization_id=%s for profile %s", value, currentProfile)
		case "workspace":
			profile.Workspace = value
			log.Printf("[DEBUG] Set workspace=%s for profile %s", value, currentProfile)
		}
		config.Profiles[currentProfile] = profile
	}

	// Get the profile from environment variable or use default
	profileName := os.Getenv("CRIBL_PROFILE")
	if profileName == "" {
		profileName = "default"
	}
	log.Printf("[DEBUG] Using profile: %s", profileName)

	// Get the specified profile
	profile, ok := config.Profiles[profileName]
	if !ok {
		if profileName == "default" {
			// If default profile is not found, use the first available profile
			for _, p := range config.Profiles {
				profile = p
				break
			}
			log.Printf("[DEBUG] Default profile not found, using first available profile")
		} else {
			// If a specific profile was requested but not found, return an error
			log.Printf("[ERROR] Profile '%s' not found in credentials file", profileName)
			return nil, fmt.Errorf("profile '%s' not found in credentials file", profileName)
		}
	}

	log.Printf("[DEBUG] Selected profile values - clientID=%s, orgID=%s, workspace=%s",
		profile.ClientID, profile.OrganizationID, profile.Workspace)
	return &profile, nil

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
