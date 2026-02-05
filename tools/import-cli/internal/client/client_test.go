package client

import (
	"os"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/spf13/viper"
)

func TestNewFromConfig_setsEnvFromConfig(t *testing.T) {
	// Mutates process env; restore in cleanup. Do not run in parallel.
	v := viper.New()
	v.Set(config.KeyOrganizationID, "test-org-1")
	v.Set(config.KeyWorkspaceID, "main")
	cfg := config.NewConfig(v)

	origOrg := os.Getenv(config.EnvOrganizationID)
	origWorkspace := os.Getenv(config.EnvWorkspaceID)
	t.Cleanup(func() {
		_ = os.Setenv(config.EnvOrganizationID, origOrg)
		_ = os.Setenv(config.EnvWorkspaceID, origWorkspace)
	})

	_, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig: %v", err)
	}
	if got := os.Getenv(config.EnvOrganizationID); got != "test-org-1" {
		t.Errorf("CRIBL_ORGANIZATION_ID: got %q want test-org-1", got)
	}
	if got := os.Getenv(config.EnvWorkspaceID); got != "main" {
		t.Errorf("CRIBL_WORKSPACE_ID: got %q want main", got)
	}
}

func TestNewFromConfig_unsetsEnvWhenEmpty(t *testing.T) {
	t.Setenv(config.EnvOrganizationID, "was-set") // restored automatically after test

	v := viper.New()
	v.Set(config.KeyOrganizationID, "") // explicit empty
	cfg := config.NewConfig(v)

	_, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig: %v", err)
	}
	if got := os.Getenv(config.EnvOrganizationID); got != "" {
		t.Errorf("empty config value should unset env: got %q", got)
	}
}

func TestNewFromConfig_returnsNonNilClient(t *testing.T) {
	v := viper.New()
	v.Set(config.KeyOnpremServerURL, "https://local.cribl")
	v.Set(config.KeyBearerToken, "token")
	cfg := config.NewConfig(v)

	sdkClient, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig: %v", err)
	}
	if sdkClient == nil {
		t.Error("NewFromConfig returned nil client")
	}
}

func TestNewFromConfig_onPremEnvSet(t *testing.T) {
	origURL := os.Getenv(config.EnvOnpremServerURL)
	origToken := os.Getenv(config.EnvBearerToken)
	t.Cleanup(func() {
		_ = os.Setenv(config.EnvOnpremServerURL, origURL)
		_ = os.Setenv(config.EnvBearerToken, origToken)
	})

	v := viper.New()
	v.Set(config.KeyOnpremServerURL, "https://onprem.example.com")
	v.Set(config.KeyBearerToken, "secret-token")
	cfg := config.NewConfig(v)

	_, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig: %v", err)
	}
	if got := os.Getenv(config.EnvOnpremServerURL); got != "https://onprem.example.com" {
		t.Errorf("CRIBL_ONPREM_SERVER_URL: got %q", got)
	}
	if got := os.Getenv(config.EnvBearerToken); got != "secret-token" {
		t.Errorf("CRIBL_BEARER_TOKEN: got %q", got)
	}
}
