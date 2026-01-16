package hooks

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
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

func isGatewayHost(host string) bool {
	return strings.Contains(host, "gateway.")
}

func isLocalHost(host string) bool {
	if strings.Contains(host, "localhost") {
		return true
	}

	if strings.Contains(host, "http") {
		re := regexp.MustCompile(`https?:\/\/?`)
		host = re.ReplaceAllString(host, "")
	}

	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	ip := net.ParseIP(host)
	if ip != nil {
		return ip.IsLoopback()
	} else {
		return false
	}
}

// isRestrictedOnPremEndpoint determines if a path is for a restricted endpoint that is not supported on on-prem deployments
func isRestrictedOnPremEndpoint(path string) bool {
	switch {
	// Search endpoints contain "search/" somewhere in the path
	case strings.Contains(path, "search/"):
		// Exclude /system/config-search from restrictions as it's an admin endpoint
		if !strings.Contains(path, "/system/config-search") {
			return true
		}
	// Check for lake endpoints
	case strings.Contains(path, "/lake/") || strings.Contains(path, "products/lake/"):
		return true
	// Check for lakehouse endpoints
	case strings.Contains(path, "lakehouse"):
		return true
	// Check for products/search endpoints
	case strings.Contains(path, "products/search"):
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
	var output string

	switch {
	case providerCloudDomain != "":
		output = providerCloudDomain
	case os.Getenv("CRIBL_CLOUD_DOMAIN") != "":
		output = os.Getenv("CRIBL_CLOUD_DOMAIN")
	case config != nil && config.CloudDomain != "":
		output = config.CloudDomain
	default:
		output = "cribl.cloud"
	}

	return fmt.Sprintf("https://gateway.%s", output)
}

// constructBaseURL builds the workspace URL from credentials when needed
func constructBaseURL(input ConstructBaseUrlInput, config *CriblConfig) string {
	// Always check if environment variables are set - they take highest priority
	// If any environment variables are set, reconstruct the URL from components
	workspaceEnv := os.Getenv("CRIBL_WORKSPACE_ID")
	orgEnv := os.Getenv("CRIBL_ORGANIZATION_ID")
	domainEnv := os.Getenv("CRIBL_CLOUD_DOMAIN")

	baseURL := input.BaseURL

	// Special case: if we have a localhost/test URL, keep it even with environment variables
	if baseURL != "" && isLocalHost(baseURL) {
		log.Printf("[DEBUG] Localhost URL detected, keeping as-is: %s", baseURL)
		return baseURL
	}

	// If no environment variables or provider variables are set and we have a concrete URL, use it as-is
	// Provider variables take precedence over concrete baseURL, so check for them first
	if workspaceEnv == "" && orgEnv == "" && domainEnv == "" &&
		input.ProviderWorkspaceID == "" && input.ProviderOrgID == "" && input.ProviderCloudDomain == "" &&
		baseURL != "" && !strings.Contains(baseURL, "{workspaceName}") && !strings.Contains(baseURL, "{organizationId}") {
		log.Printf("[DEBUG] No environment or provider variables set, using provided concrete URL: %s", baseURL)
		return baseURL
	}

	var workspace, workspaceSource string
	switch {
	case input.ProviderWorkspaceID != "":
		workspace = input.ProviderWorkspaceID
		workspaceSource = "provider"
	case os.Getenv("CRIBL_WORKSPACE_ID") != "":
		workspace = os.Getenv("CRIBL_WORKSPACE_ID")
		workspaceSource = "environment"
	case config != nil && config.Workspace != "":
		workspace = config.Workspace
		workspaceSource = "credentials"
	default:
		workspace = "main" // Default workspace name
		workspaceSource = "default"
	}

	if config != nil {
		log.Printf("[DEBUG] Workspace selection: env='%s', config='%s', final='%s', source='%s'",
			workspaceEnv, config.Workspace, workspace, workspaceSource)
	} else {
		log.Printf("[DEBUG] Workspace selection: env='%s', final='%s', source='%s'",
			workspaceEnv, workspace, workspaceSource)
	}

	var orgSource, organizationID string
	switch {
	case input.ProviderOrgID != "":
		organizationID = input.ProviderOrgID
		orgSource = "provider"
	case os.Getenv("CRIBL_ORGANIZATION_ID") != "":
		organizationID = os.Getenv("CRIBL_ORGANIZATION_ID")
		orgSource = "environment"
	case config != nil && config.OrganizationID != "":
		organizationID = config.OrganizationID
		orgSource = "credentials"
	default:
		organizationID = "ian" // Default organization ID
		orgSource = "default"
	}

	if config != nil {
		log.Printf("[DEBUG] Organization selection: env='%s', config='%s', final='%s', source='%s'",
			orgEnv, config.OrganizationID, organizationID, orgSource)
	} else {
		log.Printf("[DEBUG] Organization selection: env='%s', final='%s', source='%s'",
			orgEnv, organizationID, orgSource)
	}

	// Get cloud domain with proper precedence: Environment > Config > Default
	var cloudDomain, domainSource string
	switch {
	case input.ProviderCloudDomain != "":
		cloudDomain = input.ProviderCloudDomain
		domainSource = "provider"
	case os.Getenv("CRIBL_CLOUD_DOMAIN") != "":
		cloudDomain = os.Getenv("CRIBL_CLOUD_DOMAIN")
		domainSource = "environment"
	case config != nil && config.CloudDomain != "":
		cloudDomain = config.CloudDomain
		domainSource = "config"
	default:
		cloudDomain = "cribl.cloud" // Default domain
		domainSource = "default"
	}

	log.Printf("[DEBUG] Final precedence - Workspace: '%s' (from %s), Org: '%s' (from %s), Domain: '%s' (from %s)",
		workspace, workspaceSource, organizationID, orgSource, cloudDomain, domainSource)

	// Construct the workspace URL: https://{workspace}-{organizationId}.{cloudDomain}
	constructedURL := fmt.Sprintf("https://%s-%s.%s", workspace, organizationID, cloudDomain)
	log.Printf("[DEBUG] Constructed URL: %s from workspace=%s, org=%s, domain=%s",
		constructedURL, workspace, organizationID, cloudDomain)

	return constructedURL
}
