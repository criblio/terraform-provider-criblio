package hooks

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestGetCredentialsEnvConfig(t *testing.T) {
	client := "test-client"
	secret := "test-secret"
	org := "test-org"
	workspace := "test-workspace"

	os.Setenv("CRIBL_CLIENT_ID", client)
	os.Setenv("CRIBL_CLIENT_SECRET", secret)
	os.Setenv("CRIBL_ORGANIZATION_ID", org)
	os.Setenv("CRIBL_WORKSPACE_ID", workspace)

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}

	if cfg.ClientID != client {
		t.Errorf("GetCredentials returned incorrect ClientID, expected %s got %s", client, cfg.ClientID)
	} else if cfg.ClientSecret != secret {
		t.Errorf("GetCredentials returned incorrect ClientSecret, expected %s got %s", secret, cfg.ClientSecret)
	} else if cfg.OrganizationID != org {
		t.Errorf("GetCredentials returned incorrect OrganizationID, expected %s got %s", org, cfg.OrganizationID)
	} else if cfg.Workspace != workspace {
		t.Errorf("GetCredentials returned incorrect Workspace, expected %s got %s", workspace, cfg.Workspace)
	}
}

func TestGetCredentialsNoFile(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/tmp")

	_, err := checkLocalConfigDir()
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("GetCredentials error os.ErrNotExist was expected, but returned \"%s\"", err)
	}
}

func TestGetCredentialsIniFile(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/var/tmp")

	creds := `[default]
		  client_id = your-client-id
		  client_secret = your-client-secret
		  organization_id = your-organization-id
		  workspace = your-workspace-id`

	path := "/var/tmp/.cribl"

	err := os.Mkdir(path, 0777)
	if err != nil {
		t.Errorf("Could not write temporary config directory: %s", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/credentials", path), []byte(creds), 0644)
	if err != nil {
		t.Errorf("Could not write temporary config file: %s", err)
	}

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}

	if cfg.ClientID != "your-client-id" {
		t.Errorf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID)
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret)
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID)
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace)
	}

	err = os.RemoveAll(path)
	if err != nil {
		t.Errorf("Could not remove temporary config directory: %s", err)
	}
}

func TestGetCredentialsJSONFile(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/var/tmp")

	creds := `{"client_id": "your-client-id", 
	   "client_secret": "your-client-secret",
	   "organization_id": "your-organization-id",
	   "workspace": "your-workspace-id"}`

	path := "/var/tmp/.cribl"

	err := os.Mkdir(path, 0777)
	if err != nil {
		t.Errorf("Could not write temporary config directory: %s", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/credentials", path), []byte(creds), 0644)
	if err != nil {
		t.Errorf("Could not write temporary config file: %s", err)
	}

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}

	if cfg.ClientID != "your-client-id" {
		t.Errorf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID)
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret)
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID)
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace)
	}

	err = os.RemoveAll(path)
	if err != nil {
		t.Errorf("Could not remove temporary config directory: %s", err)
	}
}

func TestCheckLocalConfigDirRegularFilePath(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/var/tmp")

	creds := []byte("hello\ngo\n")
	path := "/var/tmp/.cribl"

	err := os.Mkdir(path, 0777)
	if err != nil {
		t.Errorf("Could not write temporary config directory: %s", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/credentials", path), creds, 0644)
	if err != nil {
		t.Errorf("Could not write temporary config file: %s", err)
	}

	readCreds, err := checkLocalConfigDir()
	if err != nil {
		t.Errorf("GetCredentials error was not expected, but returned: %s", err)
	}

	if string(creds) != string(readCreds) {
		t.Errorf("GetCredentials returned \"%s\", expected \"%s\"", creds, readCreds)
	}

	err = os.RemoveAll(path)
	if err != nil {
		t.Errorf("Could not remove temporary config directory: %s", err)
	}
}

func TestCheckLocalConfigDirLegacyFilePath(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/var/tmp")

	creds := []byte("hello\ngo\n")
	path := "/var/tmp/.cribl"

	err := os.WriteFile(path, creds, 0644)
	if err != nil {
		t.Errorf("Could not write temporary config file: %s", err)
	}

	readCreds, err := checkLocalConfigDir()
	if err != nil {
		t.Errorf("GetCredentials error was not expected, but returned: %s", err)
	}

	if string(creds) != string(readCreds) {
		t.Errorf("GetCredentials returned '%s', expected '%s'", creds, readCreds)
	}

	err = os.RemoveAll(path)
	if err != nil {
		t.Errorf("Could not remove temporary config directory: %s", err)
	}
}

func TestCheckConfigFileFormatIni(t *testing.T) {
	creds := `
	                [default]
			client_id = your-client-id
			client_secret = your-client-secret
			organization_id = your-organization-id
			workspace = your-workspace-id`

	fileType, err := checkConfigFileFormat([]byte(creds))
	if err != nil {
		t.Errorf("checkConfigFileFormat error was not expected, but returned: %s", err)
	}

	if fileType != "ini" {
		t.Errorf("checkConfigFileFormat returned '%s', expected 'ini'", fileType)
	}
}

func TestCheckConfigFileFormatJSON(t *testing.T) {
	creds := `{"client_id": "your-client-id", 
		   "client_secret": "your-client-secret",
		   "organization_id": "your-organization-id",
		   "workspace": "your-workspace-id"}`

	fileType, err := checkConfigFileFormat([]byte(creds))
	if err != nil {
		t.Errorf("checkConfigFileFormat error was not expected, but returned: %s", err)
	}

	if fileType != "json" {
		t.Errorf("checkConfigFileFormat returned '%s', expected 'json'", fileType)
	}
}

func TestCheckConfigFileFormatBusted(t *testing.T) {
	creds := `"client_id"-=---"your-client-id", 
		   "client_secret""your-client-secret",
		   "organization_id"& "your-organization-id",
		   "workspace": "your-workspace-id"}`

	fileType, err := checkConfigFileFormat([]byte(creds))
	if err == nil {
		t.Errorf("checkConfigFileFormat error was expected, but not returned")
	}

	if fileType != "" {
		t.Errorf("checkConfigFileFormat returned '%s', expected ''", fileType)
	}
}

func TestParseJSONConfig(t *testing.T) {
	creds := `{"client_id": "your-client-id", 
	   "client_secret": "your-client-secret",
	   "organization_id": "your-organization-id",
	   "workspace": "your-workspace-id"}`

	cfg, err := parseJSONConfig([]byte(creds))
	if err != nil {
		t.Errorf("parseJSONConfig threw an error in operation: %s", err)
	}

	if cfg.ClientID != "your-client-id" {
		t.Errorf("parseJSONConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID)
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf("parseJSONConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret)
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf("parseJSONConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID)
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf("parseJSONConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace)
	}
}

func TestParseJSONConfigBusted(t *testing.T) {
	creds := `{"client_id":fdsffds "your-client-id", 
	   "client_secret": "your-client-secret",
	   "organization_id": "your-organization-id",
	   "workspace": "your-workspace-id"}`

	_, err := parseJSONConfig([]byte(creds))
	if err == nil {
		t.Errorf("parseJSONConfig did not throw expected error")
	}
}

func TestParseIniConfig(t *testing.T) {
	creds := `[default]
		  client_id = your-client-id
		  client_secret = your-client-secret
		  organization_id = your-organization-id
		  workspace = your-workspace-id`

	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Errorf("parseIniConfig threw an error in operation: %s", err)
	}

	if cfg.ClientID != "your-client-id" {
		t.Errorf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID)
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret)
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID)
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace)
	}
}

func TestParseIniConfigMultiProfile(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "secondary")
	creds := `[default]
		  client_id = your-client-id
		  client_secret = your-client-secret
		  organization_id = your-organization-id
		  workspace = your-workspace-id
                  [secondary]
		  client_id = your-secondary-id
		  client_secret = your-secondary-secret
		  organization_id = your-secondary-id
		  workspace = your-secondary-id
		  `

	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Errorf("parseIniConfig threw an error in operation: %s", err)
	}

	if cfg.ClientID != "your-secondary-id" {
		t.Errorf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-secondary-id", cfg.ClientID)
	} else if cfg.ClientSecret != "your-secondary-secret" {
		t.Errorf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-secondary-secret", cfg.ClientSecret)
	} else if cfg.OrganizationID != "your-secondary-id" {
		t.Errorf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-secondary-id", cfg.OrganizationID)
	} else if cfg.Workspace != "your-secondary-id" {
		t.Errorf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-secondary-id", cfg.Workspace)
	}
}

// Tests for constructBaseURL function
func TestConstructBaseURLConcreteURL(t *testing.T) {
	// Clear environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{}

	concreteURL := "https://my-concrete-url.com"
	result := hook.constructBaseURL(concreteURL, config)

	if result != concreteURL {
		t.Errorf("constructBaseURL should return concrete URL as-is, got %s expected %s", result, concreteURL)
	}
}

func TestConstructBaseURLFromConfigValues(t *testing.T) {
	// Clear environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{
		Workspace:      "test-workspace",
		OrganizationID: "test-org",
		CloudDomain:    "cribl-playground.cloud",
	}

	expected := "https://test-workspace-test-org.cribl-playground.cloud"
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}
}

func TestConstructBaseURLEnvironmentVariableFallbacks(t *testing.T) {
	// Set environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "env-workspace")
	os.Setenv("CRIBL_ORGANIZATION_ID", "env-org")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-staging.cloud")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{} // Empty config

	expected := "https://env-workspace-env-org.cribl-staging.cloud"
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}

	// Clean up
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
}

func TestConstructBaseURLDefaultFallbacks(t *testing.T) {
	// Clear environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{} // Empty config

	expected := "https://main-ian.cribl.cloud" // Default values
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}
}

func TestConstructBaseURLConfigOverridesEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "env-workspace")
	os.Setenv("CRIBL_ORGANIZATION_ID", "env-org")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-staging.cloud")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{
		Workspace:      "config-workspace",  // Should override env
		OrganizationID: "config-org",        // Should override env
		CloudDomain:    "config-domain.com", // Should override env
	}

	expected := "https://config-workspace-config-org.config-domain.com"
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}

	// Clean up
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
}

func TestConstructBaseURLMixedConfigAndEnvironment(t *testing.T) {
	// Set some environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "env-workspace")
	os.Setenv("CRIBL_ORGANIZATION_ID", "") // Not set
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-staging.cloud")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{
		OrganizationID: "config-org", // Should be used since env is empty
		// Workspace and CloudDomain empty, should use env vars
	}

	expected := "https://env-workspace-config-org.cribl-staging.cloud"
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}

	// Clean up
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
}

func TestConstructBaseURLPlaygroundEnvironment(t *testing.T) {
	// Clear environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{
		Workspace:      "main",
		OrganizationID: "beautiful-nguyen-y8y4azd", // Playground org ID
		CloudDomain:    "cribl-playground.cloud",
	}

	expected := "https://main-beautiful-nguyen-y8y4azd.cribl-playground.cloud"
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}
}

func TestConstructBaseURLTemplateBaseURL(t *testing.T) {
	// Clear environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{
		Workspace:      "test-workspace",
		OrganizationID: "test-org",
		CloudDomain:    "test-domain.com",
	}

	// Template URLs should be ignored and constructed from components
	templateURL := "https://{workspaceName}-{organizationId}.{cloudDomain}"
	expected := "https://test-workspace-test-org.test-domain.com"
	result := hook.constructBaseURL(templateURL, config)

	if result != expected {
		t.Errorf("constructBaseURL returned %s, expected %s", result, expected)
	}
}

func TestConstructBaseURLEmptyBaseURL(t *testing.T) {
	// Clear environment variables
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	hook := NewCriblTerraformHook()
	config := &CriblConfig{
		Workspace:      "custom-workspace",
		OrganizationID: "custom-org",
		CloudDomain:    "custom.domain.com",
	}

	// Empty baseURL should trigger URL construction
	expected := "https://custom-workspace-custom-org.custom.domain.com"
	result := hook.constructBaseURL("", config)

	if result != expected {
		t.Errorf("constructBaseURL with empty baseURL returned %s, expected %s", result, expected)
	}
}
