package credentials

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
	t.Cleanup(func() {
		os.Unsetenv("CRIBL_CLIENT_ID")
		os.Unsetenv("CRIBL_CLIENT_SECRET")
		os.Unsetenv("CRIBL_ORGANIZATION_ID")
		os.Unsetenv("CRIBL_WORKSPACE_ID")
	})

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}
	if cfg != nil {
		if cfg.ClientID != client {
			t.Errorf("GetCredentials returned incorrect ClientID, expected %s got %s", client, cfg.ClientID)
		}
		if cfg.ClientSecret != secret {
			t.Errorf("GetCredentials returned incorrect ClientSecret, expected %s got %s", secret, cfg.ClientSecret)
		}
		if cfg.OrganizationID != org {
			t.Errorf("GetCredentials returned incorrect OrganizationID, expected %s got %s", org, cfg.OrganizationID)
		}
		if cfg.Workspace != workspace {
			t.Errorf("GetCredentials returned incorrect Workspace, expected %s got %s", workspace, cfg.Workspace)
		}
	}
}

func TestCheckLocalConfigDirHomeDirError(t *testing.T) {
	os.Setenv("HOME", "")
	t.Cleanup(func() { os.Unsetenv("HOME") })

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
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	t.Cleanup(func() { os.Unsetenv("HOME") })

	_, err := checkLocalConfigDir()
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("GetCredentials error os.ErrNotExist was expected, but returned \"%v\"", err)
	}
}

func TestGetCredentialsIniFile(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	t.Cleanup(func() { os.Unsetenv("HOME") })

	path := dir + "/.cribl"
	if err := os.Mkdir(path, 0777); err != nil {
		t.Fatalf("Could not write temporary config directory: %s", err)
	}
	creds := `[default]
client_id = your-client-id
client_secret = your-client-secret
organization_id = your-organization-id
workspace = your-workspace-id`
	if err := os.WriteFile(fmt.Sprintf("%s/credentials", path), []byte(creds), 0644); err != nil {
		t.Fatalf("Could not write temporary config file: %s", err)
	}

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}
	if cfg != nil {
		if cfg.ClientID != "your-client-id" {
			t.Errorf("parseIniConfig returned incorrect ClientID, expected your-client-id got %s", cfg.ClientID)
		}
		if cfg.ClientSecret != "your-client-secret" {
			t.Errorf("parseIniConfig returned incorrect ClientSecret, expected your-client-secret got %s", cfg.ClientSecret)
		}
		if cfg.OrganizationID != "your-organization-id" {
			t.Errorf("parseIniConfig returned incorrect OrganizationID, expected your-organization-id got %s", cfg.OrganizationID)
		}
		if cfg.Workspace != "your-workspace-id" {
			t.Errorf("parseIniConfig returned incorrect Workspace, expected your-workspace-id got %s", cfg.Workspace)
		}
	}
}

func TestParseIniConfigLoadFailure(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "")
	t.Cleanup(func() { os.Unsetenv("CRIBL_PROFILE") })

	creds := `default
client_id = your-client-id`
	_, err := parseIniConfig([]byte(creds))
	if err == nil {
		t.Errorf("parseIniConfig did not throw expected error in operation")
	} else if !strings.Contains(err.Error(), "failed to parse config file") {
		t.Errorf("parseIniConfig did not return expected error, got: %s", err)
	}
}

func TestGetCredentialsJSONFile(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	t.Cleanup(func() { os.Unsetenv("HOME") })

	path := dir + "/.cribl"
	if err := os.Mkdir(path, 0777); err != nil {
		t.Fatalf("Could not write temporary config directory: %s", err)
	}
	creds := `{"client_id": "your-client-id", "client_secret": "your-client-secret", "organization_id": "your-organization-id", "workspace": "your-workspace-id"}`
	if err := os.WriteFile(fmt.Sprintf("%s/credentials", path), []byte(creds), 0644); err != nil {
		t.Fatalf("Could not write temporary config file: %s", err)
	}

	cfg, err := GetCredentials()
	if err != nil {
		t.Errorf("GetCredentials threw an error in operation: %s", err)
	}
	if cfg != nil {
		if cfg.ClientID != "your-client-id" {
			t.Errorf("parseJSONConfig returned incorrect ClientID, expected your-client-id got %s", cfg.ClientID)
		}
		if cfg.Workspace != "your-workspace-id" {
			t.Errorf("parseJSONConfig returned incorrect Workspace, expected your-workspace-id got %s", cfg.Workspace)
		}
	}
}

func TestCheckLocalConfigDirRegularFilePath(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	t.Cleanup(func() { os.Unsetenv("HOME") })

	path := dir + "/.cribl"
	creds := []byte("hello\ngo\n")
	if err := os.Mkdir(path, 0777); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fmt.Sprintf("%s/credentials", path), creds, 0644); err != nil {
		t.Fatal(err)
	}

	readCreds, err := checkLocalConfigDir()
	if err != nil {
		t.Errorf("GetCredentials error was not expected, but returned: %s", err)
	}
	if string(creds) != string(readCreds) {
		t.Errorf("GetCredentials returned %q, expected %q", readCreds, creds)
	}
}

func TestCheckLocalConfigDirLegacyFilePath(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	t.Cleanup(func() { os.Unsetenv("HOME") })

	legacyPath := dir + "/.cribl"
	creds := []byte("hello\ngo\n")
	if err := os.WriteFile(legacyPath, creds, 0644); err != nil {
		t.Fatal(err)
	}

	readCreds, err := checkLocalConfigDir()
	if err != nil {
		t.Errorf("GetCredentials error was not expected, but returned: %s", err)
	}
	if string(creds) != string(readCreds) {
		t.Errorf("GetCredentials returned %q, expected %q", readCreds, creds)
	}
}

func TestCheckConfigFileFormatIni(t *testing.T) {
	creds := `[default]
client_id = your-client-id`
	fileType, err := checkConfigFileFormat([]byte(creds))
	if err != nil {
		t.Errorf("checkConfigFileFormat error was not expected, but returned: %s", err)
	}
	if fileType != "ini" {
		t.Errorf("checkConfigFileFormat returned %q, expected 'ini'", fileType)
	}
}

func TestCheckConfigFileFormatJSON(t *testing.T) {
	creds := `{"client_id": "your-client-id"}`
	fileType, err := checkConfigFileFormat([]byte(creds))
	if err != nil {
		t.Errorf("checkConfigFileFormat error was not expected, but returned: %s", err)
	}
	if fileType != "json" {
		t.Errorf("checkConfigFileFormat returned %q, expected 'json'", fileType)
	}
}

func TestCheckConfigFileFormatBusted(t *testing.T) {
	creds := `"client_id"-=---"your-client-id"`
	fileType, err := checkConfigFileFormat([]byte(creds))
	if err == nil {
		t.Errorf("checkConfigFileFormat error was expected, but not returned")
	}
	if fileType != "" {
		t.Errorf("checkConfigFileFormat returned %q, expected ''", fileType)
	}
}

func TestParseJSONConfig(t *testing.T) {
	creds := `{"client_id": "your-client-id", "client_secret": "your-client-secret", "organization_id": "your-organization-id", "workspace": "your-workspace-id"}`
	cfg, err := parseJSONConfig([]byte(creds))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientID != "your-client-id" || cfg.ClientSecret != "your-client-secret" || cfg.OrganizationID != "your-organization-id" || cfg.Workspace != "your-workspace-id" {
		t.Errorf("parseJSONConfig returned incorrect values: %+v", cfg)
	}
}

func TestParseJSONConfigBusted(t *testing.T) {
	creds := `{"client_id":fdsffds "your-client-id"}`
	_, err := parseJSONConfig([]byte(creds))
	if err == nil {
		t.Errorf("parseJSONConfig did not throw expected error")
	}
}

func TestParseIniConfig(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "")
	t.Cleanup(func() { os.Unsetenv("CRIBL_PROFILE") })

	creds := `[default]
client_id = your-client-id
client_secret = your-client-secret
organization_id = your-organization-id
workspace = your-workspace-id`
	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientID != "your-client-id" || cfg.Workspace != "your-workspace-id" {
		t.Errorf("parseIniConfig returned incorrect values: %+v", cfg)
	}
}

func TestParseIniConfigMultiProfile(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "secondary")
	t.Cleanup(func() { os.Setenv("CRIBL_PROFILE", "") })

	creds := `[default]
client_id = your-client-id
[secondary]
client_id = your-secondary-id
client_secret = your-secondary-secret
organization_id = your-secondary-id
workspace = your-secondary-id`
	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientID != "your-secondary-id" || cfg.OrganizationID != "your-secondary-id" {
		t.Errorf("parseIniConfig returned incorrect values: %+v", cfg)
	}
}

func TestParseIniConfigWithCloudDomain(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "")
	t.Cleanup(func() { os.Unsetenv("CRIBL_PROFILE") })

	creds := `[default]
client_id = your-client-id
client_secret = your-client-secret
organization_id = your-organization-id
workspace = your-workspace-id
cloud_domain = cribl-playground.cloud`
	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CloudDomain != "cribl-playground.cloud" {
		t.Errorf("parseIniConfig returned CloudDomain %q, expected cribl-playground.cloud", cfg.CloudDomain)
	}
}

func TestParseJSONConfigWithCloudDomain(t *testing.T) {
	creds := `{"client_id": "x", "client_secret": "y", "organization_id": "o", "workspace": "w", "cloud_domain": "cribl-staging.cloud"}`
	cfg, err := parseJSONConfig([]byte(creds))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CloudDomain != "cribl-staging.cloud" {
		t.Errorf("parseJSONConfig returned CloudDomain %q, expected cribl-staging.cloud", cfg.CloudDomain)
	}
}

func TestParseIniConfigMultiProfileWithCloudDomain(t *testing.T) {
	os.Setenv("CRIBL_PROFILE", "playground")
	t.Cleanup(func() { os.Setenv("CRIBL_PROFILE", "") })

	creds := `[default]
cloud_domain = cribl.cloud
[playground]
client_id = playground-client-id
client_secret = playground-client-secret
organization_id = playground-organization-id
workspace = playground-workspace-id
cloud_domain = cribl-playground.cloud`
	cfg, err := parseIniConfig([]byte(creds))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientID != "playground-client-id" || cfg.CloudDomain != "cribl-playground.cloud" {
		t.Errorf("parseIniConfig returned incorrect values: %+v", cfg)
	}
}
