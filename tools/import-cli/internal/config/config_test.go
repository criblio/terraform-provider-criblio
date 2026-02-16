package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestGet(t *testing.T) {
	v := viper.New()
	v.Set(KeyOrganizationID, "  org1  ")
	cfg := NewConfig(v)
	if got := cfg.Get(KeyOrganizationID); got != "org1" {
		t.Errorf("Get trim space: got %q", got)
	}
	v.Set(KeyWorkspaceID, "")
	if got := cfg.Get(KeyWorkspaceID); got != "" {
		t.Errorf("Get empty: got %q", got)
	}
}

func TestBindEnv(t *testing.T) {
	const val = "env-org-1"
	t.Setenv(EnvOrganizationID, val)
	t.Cleanup(func() { _ = os.Unsetenv(EnvOrganizationID) })

	cfg := NewConfig(viper.New())
	cfg.BindEnv()
	if got := cfg.Get(KeyOrganizationID); got != val {
		t.Errorf("after BindEnv, Get(KeyOrganizationID): got %q want %q", got, val)
	}
}

func TestValidateRequired_onPremMissingAuth(t *testing.T) {
	v := viper.New()
	v.Set(KeyOnpremServerURL, "https://cribl.local")
	cfg := NewConfig(v)
	cfg.BindEnv()
	err := cfg.ValidateRequired()
	if err == nil {
		t.Fatal("expected error for on-prem without auth")
	}
	if err.Error() == "" || len(err.Error()) < 10 {
		t.Errorf("expected actionable message, got %q", err.Error())
	}
}

func TestValidateRequired_onPremWithToken(t *testing.T) {
	v := viper.New()
	v.Set(KeyOnpremServerURL, "https://cribl.local")
	v.Set(KeyBearerToken, "token")
	cfg := NewConfig(v)
	err := cfg.ValidateRequired()
	if err != nil {
		t.Errorf("on-prem with token should be valid: %v", err)
	}
}

func TestValidateRequired_onPremInvalidURL(t *testing.T) {
	tests := []struct {
		name    string
		serverURL string
		wantErr string
	}{
		{"missing scheme", "cribl.local", "scheme must be http or https"},
		{"empty host", "https://", "missing host"},
		{"not a URL", "://bad", "invalid on-prem server URL"},
		{"invalid scheme", "ftp://cribl.local", "scheme must be http or https"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.Set(KeyOnpremServerURL, tt.serverURL)
			v.Set(KeyBearerToken, "token")
			cfg := NewConfig(v)
			err := cfg.ValidateRequired()
			if err == nil {
				t.Fatal("expected error for invalid on-prem URL")
			}
			if err.Error() == "" || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateRequired_cloudTokenMissingOrgWorkspace(t *testing.T) {
	v := viper.New()
	v.Set(KeyBearerToken, "token")
	cfg := NewConfig(v)
	err := cfg.ValidateRequired()
	if err == nil {
		t.Fatal("expected error for cloud token without org/workspace")
	}
}

func TestValidateRequired_cloudTokenWithOrgWorkspace(t *testing.T) {
	v := viper.New()
	v.Set(KeyBearerToken, "token")
	v.Set(KeyOrganizationID, "org1")
	v.Set(KeyWorkspaceID, "main")
	cfg := NewConfig(v)
	err := cfg.ValidateRequired()
	if err != nil {
		t.Errorf("cloud token with org/workspace should be valid: %v", err)
	}
}

func TestValidateRequired_cloudClientCredsMissingOrgWorkspace(t *testing.T) {
	v := viper.New()
	v.Set(KeyClientID, "cid")
	v.Set(KeyClientSecret, "secret")
	cfg := NewConfig(v)
	err := cfg.ValidateRequired()
	if err == nil {
		t.Fatal("expected error for client creds without org/workspace")
	}
}

func TestValidateRequired_cloudClientCredsComplete(t *testing.T) {
	v := viper.New()
	v.Set(KeyClientID, "cid")
	v.Set(KeyClientSecret, "secret")
	v.Set(KeyOrganizationID, "org1")
	v.Set(KeyWorkspaceID, "main")
	cfg := NewConfig(v)
	err := cfg.ValidateRequired()
	if err != nil {
		t.Errorf("cloud client creds with org/workspace should be valid: %v", err)
	}
}

func TestValidateRequired_noConfig(t *testing.T) {
	cfg := NewConfig(viper.New())
	err := cfg.ValidateRequired()
	if err == nil {
		t.Fatal("expected error when no config set")
	}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
	// Auth failures must surface readable, user-focused errors: tell user what to set.
	if !strings.Contains(err.Error(), EnvOnpremServerURL) && !strings.Contains(err.Error(), "on-prem") {
		t.Errorf("expected missing-credentials error to mention on-prem or %s, got: %s", EnvOnpremServerURL, err.Error())
	}
}

func TestLoadCredentialsFile_noFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Cleanup(func() { _ = os.Unsetenv("HOME") })

	cfg := NewConfig(viper.New())
	err := cfg.LoadCredentialsFile()
	if err != nil {
		t.Errorf("no credentials file should return nil (not error): %v", err)
	}
}

func TestLoadCredentialsFile_validIni(t *testing.T) {
	dir := t.TempDir()
	credDir := filepath.Join(dir, ".cribl")
	if err := os.MkdirAll(credDir, 0755); err != nil {
		t.Fatal(err)
	}
	credPath := filepath.Join(credDir, "credentials")
	iniContent := `[default]
organization_id = test-org
workspace = test-ws
client_id = test-cid
client_secret = test-secret
`
	if err := os.WriteFile(credPath, []byte(iniContent), 0600); err != nil {
		t.Fatal(err)
	}

	origHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	t.Cleanup(func() { _ = os.Setenv("HOME", origHome) })

	cfg := NewConfig(viper.New())
	err := cfg.LoadCredentialsFile()
	if err != nil {
		t.Fatalf("LoadCredentialsFile: %v", err)
	}
	if got := cfg.Get(KeyOrganizationID); got != "test-org" {
		t.Errorf("organization_id: got %q want test-org", got)
	}
	if got := cfg.Get(KeyWorkspaceID); got != "test-ws" {
		t.Errorf("workspace_id: got %q want test-ws", got)
	}
	if got := cfg.Get(KeyClientID); got != "test-cid" {
		t.Errorf("client_id: got %q want test-cid", got)
	}
}

func TestLoadCredentialsFile_fileOnlyFillsUnset(t *testing.T) {
	dir := t.TempDir()
	credDir := filepath.Join(dir, ".cribl")
	if err := os.MkdirAll(credDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(credDir, "credentials"), []byte("[default]\norganization_id = from-file\nworkspace = from-file-ws\n"), 0600); err != nil {
		t.Fatal(err)
	}
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	t.Cleanup(func() { _ = os.Setenv("HOME", origHome) })

	cfg := NewConfig(viper.New())
	cfg.v.Set(KeyOrganizationID, "from-flag") // value already set (e.g. from flag)
	cfg.BindEnv()
	err := cfg.LoadCredentialsFile()
	if err != nil {
		t.Fatal(err)
	}
	// Viper precedence: already-set value wins over file. So org should still be from-flag.
	if got := cfg.Get(KeyOrganizationID); got != "from-flag" {
		t.Errorf("existing value should not be overwritten by file: got %q", got)
	}
	// Workspace was not set, so file should have filled it.
	if got := cfg.Get(KeyWorkspaceID); got != "from-file-ws" {
		t.Errorf("workspace from file: got %q want from-file-ws", got)
	}
}
