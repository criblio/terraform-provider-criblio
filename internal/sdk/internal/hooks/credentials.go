package hooks

import (
	"encoding/json"
	"fmt"
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

	// Then try ~/.cribl/credentials file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[ERROR] Failed to get home directory: %v", err)
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Check for credentials in ~/.cribl/credentials
	configDir := filepath.Join(homeDir, ".cribl")
	configPath := filepath.Join(configDir, "credentials")

	log.Printf("[DEBUG] Reading credentials from: %s", configPath)
	file, err := os.ReadFile(configPath)
	if err != nil {
		// Fallback to legacy ~/.cribl file if credentials file doesn't exist
		legacyPath := filepath.Join(homeDir, ".cribl")
		log.Printf("[DEBUG] Credentials file not found, trying legacy path: %s", legacyPath)
		file, err = os.ReadFile(legacyPath)
		if err != nil {
			log.Printf("[ERROR] Failed to read both credentials files: %v", err)
			return nil, fmt.Errorf("failed to read credentials file: %v", err)
		}
	}

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

	// If no profiles were found, try parsing as JSON
	if len(config.Profiles) == 0 {
		log.Printf("[DEBUG] No profiles found in INI format, trying JSON")
		var legacyConfig CriblConfig
		if err := json.Unmarshal(file, &legacyConfig); err != nil {
			log.Printf("[ERROR] Failed to parse config file: %v", err)
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
		return &legacyConfig, nil
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
