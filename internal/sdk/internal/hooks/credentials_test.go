package hooks

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestGetCredentialsEnvConfig(t *testing.T) {
	client := "yo"
	secret := "AllAmericanRejects"
	org := "lung"
	workspace := "desk"

	os.Setenv("CRIBL_CLIENT_ID", client)
	os.Setenv("CRIBL_CLIENT_SECRET", secret)
	os.Setenv("CRIBL_ORGANIZATION_ID", org)
	os.Setenv("CRIBL_WORKSPACE_ID", workspace)

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}

	if cfg.ClientID != client {
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect ClientID, expected %s got %s", client, cfg.ClientID))
	} else if cfg.ClientSecret != secret {
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect ClientSecret, expected %s got %s", secret, cfg.ClientSecret))
	} else if cfg.OrganizationID != org {
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect OrganizationID, expected %s got %s", org, cfg.OrganizationID))
	} else if cfg.Workspace != workspace {
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect Workspace, expected %s got %s", workspace, cfg.Workspace))
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
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID))
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret))
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID))
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace))
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
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID))
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret))
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID))
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace))
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
		t.Errorf(fmt.Sprintf("parseJSONConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID))
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf(fmt.Sprintf("parseJSONConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret))
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf(fmt.Sprintf("parseJSONConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID))
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf(fmt.Sprintf("parseJSONConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace))
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
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect ClientID, expected %s got %s", "your-client-id", cfg.ClientID))
	} else if cfg.ClientSecret != "your-client-secret" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect ClientSecret, expected %s got %s", "your-client-secret", cfg.ClientSecret))
	} else if cfg.OrganizationID != "your-organization-id" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect OrganizationID, expected %s got %s", "your-organization-id", cfg.OrganizationID))
	} else if cfg.Workspace != "your-workspace-id" {
		t.Errorf(fmt.Sprintf("parseIniConfig returned incorrect Workspace, expected %s got %s", "your-workspace-id", cfg.Workspace))
	}
}
