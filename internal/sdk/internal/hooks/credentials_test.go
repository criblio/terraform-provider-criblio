package hooks

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

func TestCheckLocalConfigDirHomeDirError(t *testing.T) {
	os.Setenv("HOME", "")

	readCreds, err := checkLocalConfigDir()
	if err == nil {
		t.Errorf("GetCredentials error was not returned, but expected")
	}

	if string(readCreds) != string([]byte{}) {
		t.Errorf("checkLocalConfigDir unexpected return - want []byte{}, got %s", readCreds)
	}
}

func TestGetCredentialsNoFile(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/var/tmp")

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

func TestParseIniConfigLoadFailure(t *testing.T) {
	// Reset profile to default for this test
	os.Setenv("CRIBL_PROFILE", "")

	creds := `default
                 client_id = your-client-id
                 client_secret = your-client-secret
                 organization_id = your-organization-id
                 workspace = your-workspace-id
                 cloud_domain = cribl-playground.cloud`

	_, err := parseIniConfig([]byte(creds))
	if err == nil {
		t.Errorf("parseIniConfig did not throw expected error in operation")
	} else {
		if !strings.Contains(err.Error(), "failed to parse config file") {
			t.Errorf("parseIniConfig did not return expected error, got: %s", err)
		}
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

func TestParseIniConfigWithCloudDomain(t *testing.T) {
	// Reset profile to default for this test
	os.Setenv("CRIBL_PROFILE", "")

	creds := `[default]
		  client_id = your-client-id
		  client_secret = your-client-secret
		  organization_id = your-organization-id
		  workspace = your-workspace-id
		  cloud_domain = cribl-playground.cloud`

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
	} else if cfg.CloudDomain != "cribl-playground.cloud" {
		t.Errorf("parseIniConfig returned incorrect CloudDomain, expected %s got %s", "cribl-playground.cloud", cfg.CloudDomain)
	}
}

func TestParseJSONConfigWithCloudDomain(t *testing.T) {
	creds := `{"client_id": "your-client-id", 
	   "client_secret": "your-client-secret",
	   "organization_id": "your-organization-id",
	   "workspace": "your-workspace-id",
	   "cloud_domain": "cribl-staging.cloud"}`

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
	} else if cfg.CloudDomain != "cribl-staging.cloud" {
		t.Errorf("parseJSONConfig returned incorrect CloudDomain, expected %s got %s", "cribl-staging.cloud", cfg.CloudDomain)
	}
}

func TestParseIniConfigMultiProfileWithCloudDomain(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "playground")
	creds := `[default]
		  client_id = default-client-id
		  client_secret = default-client-secret
		  organization_id = default-organization-id
		  workspace = default-workspace-id
		  cloud_domain = cribl.cloud
                  [playground]
		  client_id = playground-client-id
		  client_secret = playground-client-secret
		  organization_id = playground-organization-id
		  workspace = playground-workspace-id
		  cloud_domain = cribl-playground.cloud
		  `

	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Errorf("parseIniConfig threw an error in operation: %s", err)
	}

	if cfg.ClientID != "playground-client-id" {
		t.Errorf("parseIniConfig returned incorrect ClientID, expected %s got %s", "playground-client-id", cfg.ClientID)
	} else if cfg.ClientSecret != "playground-client-secret" {
		t.Errorf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "playground-client-secret", cfg.ClientSecret)
	} else if cfg.OrganizationID != "playground-organization-id" {
		t.Errorf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "playground-organization-id", cfg.OrganizationID)
	} else if cfg.Workspace != "playground-workspace-id" {
		t.Errorf("parseIniConfig returned incorrect Workspace, expected %s got %s", "playground-workspace-id", cfg.Workspace)
	} else if cfg.CloudDomain != "cribl-playground.cloud" {
		t.Errorf("parseIniConfig returned incorrect CloudDomain, expected %s got %s", "cribl-playground.cloud", cfg.CloudDomain)
	}

	// Clean up
	os.Setenv("CRIBL_PROFILE", "")
}
