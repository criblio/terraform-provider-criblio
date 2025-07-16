package hooks

import (
	"fmt"
	"os"
	"testing"
	"errors"
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
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect ClientID, expected %s got %s", secret, cfg.ClientSecret))
	} else if cfg.OrganizationID != org {
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect ClientID, expected %s got %s", org, cfg.OrganizationID))
	} else if cfg.Workspace != workspace {
		t.Errorf(fmt.Sprintf("GetCredentials returned incorrect ClientID, expected %s got %s", workspace, cfg.Workspace))
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

func TestGetCredentialsRegularFilePath(t *testing.T) {
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

func TestGetCredentialsLegacyFilePath(t *testing.T) {
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
