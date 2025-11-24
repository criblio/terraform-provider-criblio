package hooks

import (
	"fmt"
	"testing"
)

func TestTrimPath(t *testing.T) {
	example := "/api/v1/bar"
	output := trimPath(example)
	expected := "bar"
	if output != expected {
		t.Fatal(fmt.Sprintf("got wrong output from trimPath, expected '%s' and got '%s'", expected, output))
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
