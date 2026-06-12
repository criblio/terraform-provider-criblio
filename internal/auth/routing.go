package auth

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

// TrimPath removes leading slashes and the api/v1 prefix from path.
func TrimPath(path string) string {
	path = strings.TrimLeft(path, "/")
	path = strings.TrimPrefix(path, "api/v1/")
	path = strings.TrimPrefix(path, "api/v1")

	return path
}

// IsGatewayPath reports whether path is a gateway management endpoint.
func IsGatewayPath(path string) bool {
	cleanPath := strings.TrimLeft(path, "/")
	cleanPath = strings.TrimPrefix(cleanPath, "api/")
	cleanPath = strings.TrimPrefix(cleanPath, "v1/")

	return strings.HasPrefix(cleanPath, "organizations/")
}

// IsGatewayHost reports whether host is a gateway host.
func IsGatewayHost(host string) bool {
	return strings.Contains(host, "gateway.")
}

// IsLocalHost reports whether host is localhost or a loopback IP.
func IsLocalHost(host string) bool {
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
	if ip == nil {
		return false
	}

	return ip.IsLoopback()
}

// ConstructGatewayURL builds the gateway URL using Provider > Environment > Credentials > Default precedence.
func ConstructGatewayURL(providerCloudDomain string, config *CriblConfig) string {
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

// ConstructBaseURL builds a workspace URL from credentials when needed.
func ConstructBaseURL(input ConstructBaseURLInput, config *CriblConfig) string {
	workspaceEnv := os.Getenv("CRIBL_WORKSPACE_ID")
	orgEnv := os.Getenv("CRIBL_ORGANIZATION_ID")
	domainEnv := os.Getenv("CRIBL_CLOUD_DOMAIN")

	baseURL := input.BaseURL

	if baseURL != "" && IsLocalHost(baseURL) {
		log.Printf("[DEBUG] Localhost URL detected, keeping as-is: %s", baseURL)
		return baseURL
	}

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
		workspace = "main"
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
		organizationID = "ian"
		orgSource = "default"
	}

	if config != nil {
		log.Printf("[DEBUG] Organization selection: env='%s', config='%s', final='%s', source='%s'",
			orgEnv, config.OrganizationID, organizationID, orgSource)
	} else {
		log.Printf("[DEBUG] Organization selection: env='%s', final='%s', source='%s'",
			orgEnv, organizationID, orgSource)
	}

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
		cloudDomain = "cribl.cloud"
		domainSource = "default"
	}

	log.Printf("[DEBUG] Final precedence - Workspace: '%s' (from %s), Org: '%s' (from %s), Domain: '%s' (from %s)",
		workspace, workspaceSource, organizationID, orgSource, cloudDomain, domainSource)

	constructedURL := fmt.Sprintf("https://%s-%s.%s", workspace, organizationID, cloudDomain)
	log.Printf("[DEBUG] Constructed URL: %s from workspace=%s, org=%s, domain=%s",
		constructedURL, workspace, organizationID, cloudDomain)

	return constructedURL
}

// IsRestrictedOnPremEndpoint reports whether path is unavailable on-prem.
// The cloudOnlyPaths map is generated by tools/merge-spec from upstream
// x-cribl-availability annotations.
func IsRestrictedOnPremEndpoint(path string) bool {
	normalized := normalizeRestrictedEndpointPath(path)
	if cloudOnlyPaths[normalized] {
		return true
	}

	for template := range cloudOnlyPaths {
		if pathMatchesTemplate(normalized, template) {
			return true
		}
	}
	return false
}

func normalizeRestrictedEndpointPath(path string) string {
	path, _, _ = strings.Cut(path, "?")
	path = strings.TrimSpace(path)
	path = strings.TrimLeft(path, "/")
	path = strings.TrimPrefix(path, "api/v1/")
	path = strings.TrimPrefix(path, "api/v1")
	path = "/" + strings.TrimLeft(path, "/")

	if after, ok := strings.CutPrefix(path, "/m/"); ok {
		if slash := strings.Index(after, "/"); slash >= 0 {
			path = after[slash:]
		}
	}
	return path
}

func pathMatchesTemplate(path, template string) bool {
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")
	templateSegments := strings.Split(strings.Trim(template, "/"), "/")
	if len(pathSegments) != len(templateSegments) {
		return false
	}

	for index := range templateSegments {
		if isPathTemplateSegment(templateSegments[index]) {
			if pathSegments[index] == "" {
				return false
			}
			continue
		}
		if pathSegments[index] != templateSegments[index] {
			return false
		}
	}
	return true
}

func isPathTemplateSegment(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
}
