package client

import "github.com/criblio/terraform-provider-criblio/internal/useragent"

// BulkExporterUserAgent returns the User-Agent used by the import CLI.
// For now it matches the User-Agent User-Agent.
func BulkExporterUserAgent() string {
	return useragent.TerraformProvider
}
