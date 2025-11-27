package hooks

import (
	"os"
	"testing"
)

func TestTrimPath(t *testing.T) {
	example := "/api/v1/bar"
	output := trimPath(example)
	expected := "bar"
	if output != expected {
		t.Errorf("got wrong output from trimPath, expected '%s' and got '%s'", expected, output)
	}
}

func TestIsGatewayPath(t *testing.T) {
	// Test cases for gateway paths
	gatewayPaths := []struct {
		path     string
		expected bool
		desc     string
	}{
		{"/v1/organizations/my-org/workspaces", true, "workspace creation path"},
		{"/api/v1/organizations/my-org/workspaces", true, "workspace creation path with api prefix"},
		{"v1/organizations/my-org/workspaces/workspace-id", true, "workspace operations path"},
		{"api/v1/organizations/my-org/workspaces/workspace-id", true, "workspace operations path with api prefix"},
		{"/v1/workspaces/workspace-id/sources", false, "regular workspace API path"},
		{"/api/v1/workspaces/workspace-id/destinations", false, "regular workspace API path"},
		{"/v1/system/health", false, "system health path"},
		{"", false, "empty path"},
		{"/", false, "root path"},
	}

	for _, test := range gatewayPaths {
		result := isGatewayPath(test.path)
		if result != test.expected {
			t.Errorf("isGatewayPath(%q) = %v, expected %v (%s)", test.path, result, test.expected, test.desc)
		}
	}
}

func TestConstructGatewayURL(t *testing.T) {
	// Test with default domain
	result := constructGatewayURL("", nil)
	expected := "https://gateway.cribl.cloud"
	if result != expected {
		t.Errorf("constructGatewayURL('', nil) = %q, expected %q", result, expected)
	}

	// Test with provider cloud domain
	result = constructGatewayURL("cribl-playground.cloud", nil)
	expected = "https://gateway.cribl-playground.cloud"
	if result != expected {
		t.Errorf("constructGatewayURL('cribl-playground.cloud', nil) = %q, expected %q", result, expected)
	}

	// Test with config cloud domain
	config := &CriblConfig{
		CloudDomain: "cribl-staging.cloud",
	}
	result = constructGatewayURL("", config)
	expected = "https://gateway.cribl-staging.cloud"
	if result != expected {
		t.Errorf("constructGatewayURL('', config) = %q, expected %q", result, expected)
	}

	// Test provider takes precedence over config
	result = constructGatewayURL("cribl-prod.cloud", config)
	expected = "https://gateway.cribl-prod.cloud"
	if result != expected {
		t.Errorf("constructGatewayURL('cribl-prod.cloud', config) = %q, expected %q", result, expected)
	}
}

func TestConstructBaseURL(t *testing.T) {
	//no env, concrete url
	input := ConstructBaseUrlInput{
		BaseURL: "foo",
	}
	expected := "foo"
	result := constructBaseURL(input, nil)
	if result != expected {
		t.Errorf("constructBaseUrl returned %q, expected %q", result, expected)
	}

	//localhost url
	input = ConstructBaseUrlInput{
		BaseURL: "127.0.0.1",
	}
	expected = "127.0.0.1"
	result = constructBaseURL(input, nil)
	if result != expected {
		t.Errorf("constructBaseUrl returned %q, expected %q", result, expected)
	}

	//input from config
	input = ConstructBaseUrlInput{}
	config := CriblConfig{
		OrganizationID: "ConfigOrgID",
		Workspace:      "ConfigWorkspaceID",
		CloudDomain:    "ConfigCloudDomain",
	}
	expected = "https://ConfigWorkspaceID-ConfigOrgID.ConfigCloudDomain"
	result = constructBaseURL(input, &config)
	if result != expected {
		t.Errorf("constructBaseUrl returned %q, expected %q", result, expected)
	}

	//input from env -- should override config
	os.Setenv("CRIBL_ORGANIZATION_ID", "EnvOrgID")
	os.Setenv("CRIBL_WORKSPACE_ID", "EnvWorkspaceID")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "EnvCloudDomain")
	expected = "https://EnvWorkspaceID-EnvOrgID.EnvCloudDomain"
	result = constructBaseURL(input, &config)
	if result != expected {
		t.Errorf("constructBaseUrl returned %q, expected %q", result, expected)
	}

	//input from provider -- should override env and config
	input = ConstructBaseUrlInput{
		ProviderOrgID:       "ProviderOrgID",
		ProviderWorkspaceID: "ProviderWorkspaceID",
		ProviderCloudDomain: "ProviderCloudDomain",
	}
	expected = "https://ProviderWorkspaceID-ProviderOrgID.ProviderCloudDomain"
	result = constructBaseURL(input, &config)
	if result != expected {
		t.Errorf("constructBaseUrl returned %q, expected %q", result, expected)
	}

	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")

	//check WorkspaceID default
	input = ConstructBaseUrlInput{
		ProviderOrgID:       "",
		ProviderWorkspaceID: "",
		ProviderCloudDomain: "ProviderCloudDomain",
	}
	expected = "https://main-ian.ProviderCloudDomain"
	result = constructBaseURL(input, nil)
	if result != expected {
		t.Errorf("constructBaseUrl returned %q, expected %q", result, expected)
	}

}

func TestOnPremRestrictedEndpoints(t *testing.T) {
	// Test that restricted endpoints return errors for on-prem deployments
	restrictedPaths := map[string]bool{"/products/search/jobs": true,
		"/products/lake/lakes":                         true,
		"/products/lake/lakehouses":                    true,
		"/products/lakehouses":                         true,
		"/v1/organizations/org/workspaces":             true,
		"/api/v1/organizations/org/workspaces":         true,
		"/search/jobs":                                 true,
		"/products/lake/datasets":                      true,
		"/products/search/datasets":                    true,
		"/api/v1/m/default_search/search/saved":        true,
		"/api/v1/m/default_search/search/usage-groups": true,
		"/api/v1/m/products/search":                    true,
		"/api/v1/m/default_search/search/saved-query":  true,
		"/m/default_search/search/dashboards":          true,
		"v1/organizations/foo":                         true,
		"organizations/foo":                            true,
		"/m/default/":                                  false,
		"/foo":                                         false,
	}

	for k, v := range restrictedPaths {
		if isRestrictedOnPremEndpoint(k) != v {
			t.Errorf("Expected %s to return %t, returned %t", k, v, !v)
		}
	}
}
