// Diagnostic-to-CLI error conversion for the converter package.
package converter

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// CLIError is a normalised error for CLI output, including resource context.
type CLIError struct {
	ResourceType string
	ResourceID   string
	Summary      string
	Detail       string
}

func (e CLIError) Error() string {
	if e.ResourceType != "" || e.ResourceID != "" {
		return fmt.Sprintf("%s %s: %s", e.ResourceType, e.ResourceID, e.Summary)
	}
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Summary, e.Detail)
	}
	return e.Summary
}

// DiagnosticsToCLIErrors converts Terraform diagnostics into CLI errors with optional resource context.
func DiagnosticsToCLIErrors(diags diag.Diagnostics, resourceType, resourceID string) []CLIError {
	if len(diags) == 0 {
		return nil
	}
	out := make([]CLIError, 0, len(diags))
	for _, d := range diags {
		summary := d.Summary()
		detail := d.Detail()
		if summary == "" {
			summary = "conversion error"
		}
		out = append(out, CLIError{
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Summary:      summary,
			Detail:       detail,
		})
	}
	return out
}

// DiagnosticsToError returns a single error that combines all diagnostics, with resource context.
func DiagnosticsToError(diags diag.Diagnostics, resourceType, resourceID string) error {
	if diags == nil || !diags.HasError() {
		return nil
	}
	clis := DiagnosticsToCLIErrors(diags, resourceType, resourceID)
	if len(clis) == 0 {
		return nil
	}
	var b strings.Builder
	for i, e := range clis {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	return fmt.Errorf("%s", b.String())
}
