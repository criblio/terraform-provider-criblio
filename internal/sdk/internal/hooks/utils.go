package hooks

import (
	"fmt"
  "os"
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

// isRestrictedOnPremEndpoint determines if a path is for a restricted endpoint that is not supported on on-prem deployments
func isRestrictedOnPremEndpoint(path string) bool {
	// Search endpoints contain "search/" somewhere in the path
	if strings.Contains(path, "search/") {
		// Exclude /system/config-search from restrictions as it's an admin endpoint
		if !strings.Contains(path, "/system/config-search") {
			return true
		}
	}

	// Check for lake endpoints
	if strings.Contains(path, "/lake/") || strings.Contains(path, "products/lake/") {
		return true
	}

	// Check for lakehouse endpoints
	if strings.Contains(path, "lakehouse") {
		return true
	}

	// Check for products/search endpoints
	if strings.Contains(path, "products/search/") {
		return true
	}

	// Check for gateway/management endpoints
	restrictedPrefixes := []string{
		"v1/organizations/", // Gateway management endpoints
		"organizations/",    // Workspace management endpoints
	}

	// Check for restricted prefixes
	for _, prefix := range restrictedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Check for workspace management (these are gateway paths)
	if strings.Contains(path, "workspace") && !strings.Contains(path, "/system/") {
		return true
	}

	return false
}

// constructGatewayURL builds the gateway URL using the appropriate cloud domain
func constructGatewayURL(providerCloudDomain string, config *CriblConfig) string {
	// Get cloud domain with proper precedence: Provider > Environment > Credentials > Default
	cloudDomain := providerCloudDomain
	if cloudDomain == "" {
		cloudDomain = os.Getenv("CRIBL_CLOUD_DOMAIN")
	}
	if cloudDomain == "" && config != nil {
		cloudDomain = config.CloudDomain
	}
	if cloudDomain == "" {
		cloudDomain = "cribl.cloud" // Default domain
	}

	return fmt.Sprintf("https://gateway.%s", cloudDomain)
}