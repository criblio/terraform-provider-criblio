package hooks

import (
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
	if err == nil {
		t.Errorf("GetCredentials error was expected, but not returned")
	}
}
