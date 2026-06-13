package auth

import "testing"

func TestTrimPath(t *testing.T) {
	example := "/api/v1/bar"
	output := TrimPath(example)
	expected := "bar"
	if output != expected {
		t.Errorf("got wrong output from TrimPath, expected %q and got %q", expected, output)
	}
}

func TestIsGatewayPath(t *testing.T) {
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
		result := IsGatewayPath(test.path)
		if result != test.expected {
			t.Errorf("IsGatewayPath(%q) = %v, expected %v (%s)", test.path, result, test.expected, test.desc)
		}
	}
}

func TestIsGatewayHost(t *testing.T) {
	gatewayHosts := []struct {
		host     string
		expected bool
		desc     string
	}{
		{"foo.bar.com/api/v1/gateway", false, "plain URL path"},
		{"foo.gateway.bar.com", true, "gateway subdomain"},
		{"foo.gateway.com", true, "gateway domain"},
	}

	for _, test := range gatewayHosts {
		result := IsGatewayHost(test.host)
		if result != test.expected {
			t.Errorf("IsGatewayHost(%q) = %v, expected %v (%s)", test.host, result, test.expected, test.desc)
		}
	}
}

func TestIsLocalHost(t *testing.T) {
	gatewayHosts := []struct {
		host     string
		expected bool
		desc     string
	}{
		{"foo.bar.com/api/v1/gateway", false, "plain URL path"},
		{"foo.gateway.bar.com", false, "plain URL path"},
		{"localhost", true, "localhost in words"},
		{"127.0.0.1", true, "localhost in numbers"},
		{"http://127.0.0.1:4344", true, "localhost in numbers with extra"},
		{"http://127.0.0.1", true, "localhost in numbers with extra before"},
		{"127.0.0.1:4344", true, "localhost in numbers with extra after"},
		{"http://localhost:4344", true, "localhost with scheme and port"},
		{"http://localhost", true, "localhost with scheme"},
		{"localhost:4344", true, "localhost with port"},
		{"locahost:4344", false, "misspelled localhost"},
	}

	for _, test := range gatewayHosts {
		result := IsLocalHost(test.host)
		if result != test.expected {
			t.Errorf("IsLocalHost(%q) = %v, expected %v (%s)", test.host, result, test.expected, test.desc)
		}
	}
}

func TestConstructGatewayURL(t *testing.T) {
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	result := ConstructGatewayURL("", nil)
	expected := "https://gateway.cribl.cloud"
	if result != expected {
		t.Errorf("ConstructGatewayURL('', nil) = %q, expected %q", result, expected)
	}

	result = ConstructGatewayURL("cribl-playground.cloud", nil)
	expected = "https://gateway.cribl-playground.cloud"
	if result != expected {
		t.Errorf("ConstructGatewayURL('cribl-playground.cloud', nil) = %q, expected %q", result, expected)
	}

	config := &CriblConfig{
		CloudDomain: "cribl-staging.cloud",
	}
	result = ConstructGatewayURL("", config)
	expected = "https://gateway.cribl-staging.cloud"
	if result != expected {
		t.Errorf("ConstructGatewayURL('', config) = %q, expected %q", result, expected)
	}

	result = ConstructGatewayURL("cribl-prod.cloud", config)
	expected = "https://gateway.cribl-prod.cloud"
	if result != expected {
		t.Errorf("ConstructGatewayURL('cribl-prod.cloud', config) = %q, expected %q", result, expected)
	}

	t.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-env.cloud")
	result = ConstructGatewayURL("", config)
	expected = "https://gateway.cribl-env.cloud"
	if result != expected {
		t.Errorf("ConstructGatewayURL('', config) with env = %q, expected %q", result, expected)
	}
}

func TestConstructBaseURL(t *testing.T) {
	t.Setenv("CRIBL_ORGANIZATION_ID", "")
	t.Setenv("CRIBL_WORKSPACE_ID", "")
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	input := ConstructBaseURLInput{
		BaseURL: "foo",
	}
	expected := "foo"
	result := ConstructBaseURL(input, nil)
	if result != expected {
		t.Errorf("ConstructBaseURL returned %q, expected %q", result, expected)
	}

	input = ConstructBaseURLInput{
		BaseURL: "127.0.0.1",
	}
	expected = "127.0.0.1"
	result = ConstructBaseURL(input, nil)
	if result != expected {
		t.Errorf("ConstructBaseURL returned %q, expected %q", result, expected)
	}

	input = ConstructBaseURLInput{}
	config := CriblConfig{
		OrganizationID: "ConfigOrgID",
		Workspace:      "ConfigWorkspaceID",
		CloudDomain:    "ConfigCloudDomain",
	}
	expected = "https://ConfigWorkspaceID-ConfigOrgID.ConfigCloudDomain"
	result = ConstructBaseURL(input, &config)
	if result != expected {
		t.Errorf("ConstructBaseURL returned %q, expected %q", result, expected)
	}

	t.Setenv("CRIBL_ORGANIZATION_ID", "EnvOrgID")
	t.Setenv("CRIBL_WORKSPACE_ID", "EnvWorkspaceID")
	t.Setenv("CRIBL_CLOUD_DOMAIN", "EnvCloudDomain")
	expected = "https://EnvWorkspaceID-EnvOrgID.EnvCloudDomain"
	result = ConstructBaseURL(input, &config)
	if result != expected {
		t.Errorf("ConstructBaseURL returned %q, expected %q", result, expected)
	}

	input = ConstructBaseURLInput{
		ProviderOrgID:       "ProviderOrgID",
		ProviderWorkspaceID: "ProviderWorkspaceID",
		ProviderCloudDomain: "ProviderCloudDomain",
	}
	expected = "https://ProviderWorkspaceID-ProviderOrgID.ProviderCloudDomain"
	result = ConstructBaseURL(input, &config)
	if result != expected {
		t.Errorf("ConstructBaseURL returned %q, expected %q", result, expected)
	}

	t.Setenv("CRIBL_ORGANIZATION_ID", "")
	t.Setenv("CRIBL_WORKSPACE_ID", "")
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	input = ConstructBaseURLInput{
		ProviderCloudDomain: "ProviderCloudDomain",
	}
	expected = "https://main-ian.ProviderCloudDomain"
	result = ConstructBaseURL(input, nil)
	if result != expected {
		t.Errorf("ConstructBaseURL returned %q, expected %q", result, expected)
	}
}

func TestIsRestrictedOnPremEndpoint(t *testing.T) {
	originalCloudOnlyPaths := cloudOnlyPaths
	cloudOnlyPaths = map[string]bool{
		"/search/datasets":                  true,
		"/products/lake/lakes/{lakeId}":     true,
		"/p/{pack}/search/macros/{id}/move": true,
	}
	t.Cleanup(func() {
		cloudOnlyPaths = originalCloudOnlyPaths
	})

	tests := []struct {
		path     string
		expected bool
	}{
		{"/m/default/search/datasets", true},
		{"/api/v1/m/default/products/lake/lakes/main", true},
		{"/m/default/p/search/search/macros/macro-1/move", true},
		{"/m/default/system/certificates/my-cert", false},
		{"/m/default/system/config-search", false},
	}

	for _, test := range tests {
		result := IsRestrictedOnPremEndpoint(test.path)
		if result != test.expected {
			t.Errorf("IsRestrictedOnPremEndpoint(%q) = %v, expected %v", test.path, result, test.expected)
		}
	}
}
