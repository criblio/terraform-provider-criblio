package client

import (
	"os"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/spf13/viper"
)

func TestNewFromViper_setsEnvFromViper(t *testing.T) {
	// Mutates process env; restore in cleanup. Do not run in parallel.
	v := viper.New()
	v.Set(config.KeyOrganizationID, "test-org-1")
	v.Set(config.KeyWorkspaceID, "main")

	origOrg := os.Getenv(config.EnvOrganizationID)
	origWorkspace := os.Getenv(config.EnvWorkspaceID)
	t.Cleanup(func() {
		_ = os.Setenv(config.EnvOrganizationID, origOrg)
		_ = os.Setenv(config.EnvWorkspaceID, origWorkspace)
	})

	_, err := NewFromViper(v)
	if err != nil {
		t.Fatalf("NewFromViper: %v", err)
	}
	if got := os.Getenv(config.EnvOrganizationID); got != "test-org-1" {
		t.Errorf("CRIBL_ORGANIZATION_ID: got %q want test-org-1", got)
	}
	if got := os.Getenv(config.EnvWorkspaceID); got != "main" {
		t.Errorf("CRIBL_WORKSPACE_ID: got %q want main", got)
	}
}

func TestNewFromViper_unsetsEnvWhenEmpty(t *testing.T) {
	t.Setenv(config.EnvOrganizationID, "was-set") // restored automatically after test

	v := viper.New()
	v.Set(config.KeyOrganizationID, "") // explicit empty

	_, err := NewFromViper(v)
	if err != nil {
		t.Fatalf("NewFromViper: %v", err)
	}
	if got := os.Getenv(config.EnvOrganizationID); got != "" {
		t.Errorf("empty Viper value should unset env: got %q", got)
	}
}

func TestNewFromViper_returnsNonNilClient(t *testing.T) {
	v := viper.New()
	v.Set(config.KeyOnpremServerURL, "https://local.cribl")
	v.Set(config.KeyBearerToken, "token")

	client, err := NewFromViper(v)
	if err != nil {
		t.Fatalf("NewFromViper: %v", err)
	}
	if client == nil {
		t.Error("NewFromViper returned nil client")
	}
}

func TestNewFromViper_onPremEnvSet(t *testing.T) {
	origURL := os.Getenv(config.EnvOnpremServerURL)
	origToken := os.Getenv(config.EnvBearerToken)
	t.Cleanup(func() {
		_ = os.Setenv(config.EnvOnpremServerURL, origURL)
		_ = os.Setenv(config.EnvBearerToken, origToken)
	})

	v := viper.New()
	v.Set(config.KeyOnpremServerURL, "https://onprem.example.com")
	v.Set(config.KeyBearerToken, "secret-token")

	_, err := NewFromViper(v)
	if err != nil {
		t.Fatalf("NewFromViper: %v", err)
	}
	if got := os.Getenv(config.EnvOnpremServerURL); got != "https://onprem.example.com" {
		t.Errorf("CRIBL_ONPREM_SERVER_URL: got %q", got)
	}
	if got := os.Getenv(config.EnvBearerToken); got != "secret-token" {
		t.Errorf("CRIBL_BEARER_TOKEN: got %q", got)
	}
}
