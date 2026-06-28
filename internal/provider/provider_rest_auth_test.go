package provider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
)

func TestProviderRestCredentialsUsesDefaultProfile(t *testing.T) {
	writeCredentialsFile(t, `[default]
client_id = profile-client
client_secret = profile-secret
organization_id = profile-org
workspace = profile-workspace
cloud_domain = cribl-playground.cloud
`)

	credentials := providerRestCredentials(&providerOAuthConfig{}, map[string]string{
		"organizationId": "ian",
		"workspaceId":    "main",
		"cloudDomain":    "cribl.cloud",
	}, map[string]bool{})

	if credentials.ClientID != "profile-client" {
		t.Fatalf("ClientID = %q, want profile-client", credentials.ClientID)
	}
	if credentials.ClientSecret != "profile-secret" {
		t.Fatalf("ClientSecret = %q, want profile-secret", credentials.ClientSecret)
	}
	if credentials.OrganizationID != "profile-org" {
		t.Fatalf("OrganizationID = %q, want profile-org", credentials.OrganizationID)
	}
	if credentials.Workspace != "profile-workspace" {
		t.Fatalf("Workspace = %q, want profile-workspace", credentials.Workspace)
	}
	if credentials.CloudDomain != "cribl-playground.cloud" {
		t.Fatalf("CloudDomain = %q, want cribl-playground.cloud", credentials.CloudDomain)
	}
}

func TestProviderRestCredentialsProviderValuesOverrideProfile(t *testing.T) {
	writeCredentialsFile(t, `[default]
client_id = profile-client
client_secret = profile-secret
organization_id = profile-org
workspace = profile-workspace
cloud_domain = cribl-playground.cloud
`)

	credentials := providerRestCredentials(&providerOAuthConfig{
		ClientID:     "provider-client",
		ClientSecret: "provider-secret",
	}, map[string]string{
		"organizationId": "provider-org",
		"workspaceId":    "provider-workspace",
		"cloudDomain":    "provider.cloud",
	}, map[string]bool{
		"organizationId": true,
		"workspaceId":    true,
		"cloudDomain":    true,
	})

	if credentials.ClientID != "provider-client" {
		t.Fatalf("ClientID = %q, want provider-client", credentials.ClientID)
	}
	if credentials.ClientSecret != "provider-secret" {
		t.Fatalf("ClientSecret = %q, want provider-secret", credentials.ClientSecret)
	}
	if credentials.OrganizationID != "provider-org" {
		t.Fatalf("OrganizationID = %q, want provider-org", credentials.OrganizationID)
	}
	if credentials.Workspace != "provider-workspace" {
		t.Fatalf("Workspace = %q, want provider-workspace", credentials.Workspace)
	}
	if credentials.CloudDomain != "provider.cloud" {
		t.Fatalf("CloudDomain = %q, want provider.cloud", credentials.CloudDomain)
	}
}

func TestProviderRestCredentialsExplicitDefaultDomainOverridesProfile(t *testing.T) {
	writeCredentialsFile(t, `[default]
client_id = profile-client
client_secret = profile-secret
organization_id = profile-org
workspace = profile-workspace
cloud_domain = cribl-playground.cloud
`)

	credentials := providerRestCredentials(&providerOAuthConfig{}, map[string]string{
		"organizationId": "ian",
		"workspaceId":    "main",
		"cloudDomain":    "cribl.cloud",
	}, map[string]bool{
		"cloudDomain": true,
	})

	if credentials.OrganizationID != "profile-org" {
		t.Fatalf("OrganizationID = %q, want profile-org", credentials.OrganizationID)
	}
	if credentials.Workspace != "profile-workspace" {
		t.Fatalf("Workspace = %q, want profile-workspace", credentials.Workspace)
	}
	if credentials.CloudDomain != "cribl.cloud" {
		t.Fatalf("CloudDomain = %q, want cribl.cloud", credentials.CloudDomain)
	}
}

func TestProviderRestBaseURLUsesResolvedProfile(t *testing.T) {
	credentials := &auth.CriblConfig{
		OrganizationID: "beautiful-nguyen-y8y4azd",
		Workspace:      "main",
		CloudDomain:    "cribl-playground.cloud",
	}

	got := providerRestBaseURL("https://{workspaceId}-{organizationId}.{cloudDomain}", credentials)
	want := "https://main-beautiful-nguyen-y8y4azd.cribl-playground.cloud"
	if got != want {
		t.Fatalf("providerRestBaseURL() = %q, want %q", got, want)
	}
}

func writeCredentialsFile(t *testing.T, content string) {
	t.Helper()
	t.Setenv("CRIBL_CLIENT_ID", "")
	t.Setenv("CRIBL_CLIENT_SECRET", "")
	t.Setenv("CRIBL_ORGANIZATION_ID", "")
	t.Setenv("CRIBL_WORKSPACE_ID", "")
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")
	t.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	t.Setenv("CRIBL_ONPREM_USERNAME", "")
	t.Setenv("CRIBL_ONPREM_PASSWORD", "")
	t.Setenv("CRIBL_PROFILE", "default")

	home := t.TempDir()
	t.Setenv("HOME", home)

	criblDir := filepath.Join(home, ".cribl")
	if err := os.MkdirAll(criblDir, 0755); err != nil {
		t.Fatalf("create credentials dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(criblDir, "credentials"), []byte(content), 0600); err != nil {
		t.Fatalf("write credentials file: %v", err)
	}
}
