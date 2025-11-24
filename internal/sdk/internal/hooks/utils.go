package hooks

import (
	"strings"
)

func trimPath(path string) string {
	path = strings.TrimLeft(path, "/")
	path = strings.TrimPrefix(path, "api/v1/")
	path = strings.TrimPrefix(path, "api/v1")

	return path
}

func isGatewayPath(path string) bool {
	// Remove leading slash and api prefix for consistent checking
	cleanPath := strings.TrimLeft(path, "/")
	cleanPath = strings.TrimPrefix(cleanPath, "api/")
	cleanPath = strings.TrimPrefix(cleanPath, "v1/")

	// Gateway paths are for organization and workspace management
	return strings.HasPrefix(cleanPath, "organizations/")
}
