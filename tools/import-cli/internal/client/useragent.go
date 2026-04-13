package client

import (
	"fmt"
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/version"
)

// speakeasyGeneratorVersion is the Speakeasy codegen revision; update when
// internal/sdk/criblio.go is regenerated (see its "generator version" comment).
// The CLI release version in the User-Agent is version.Version (set by GoReleaser
// ldflags from the release tag in .goreleaser-goatify.yml).
const speakeasyGeneratorVersion = "2.879.6"

// BulkExporterUserAgent returns the Speakeasy-style User-Agent for the bulk exporter CLI,
// parallel to the default in sdk.New() but with product "bulk-exporter" and the CLI build version.
func BulkExporterUserAgent() string {
	v := strings.TrimSpace(version.Version)
	if v == "" {
		v = "dev"
	}
	return fmt.Sprintf(
		"speakeasy-sdk/bulk-exporter %s %s github.com/criblio/terraform-provider-criblio/tools/import-cli",
		v,
		speakeasyGeneratorVersion,
	)
}
