// HCL parsing for validation of generated config.
package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"
)

// ParseHCL parses the given HCL source. Returns an error if parsing fails.
// Use in tests to validate that generated HCL parses successfully.
func ParseHCL(src []byte, filename string) error {
	parser := hclparse.NewParser()
	_, diags := parser.ParseHCL(src, filename)
	if diags != nil && diags.HasErrors() {
		return fmt.Errorf("parse %s: %s", filename, diags.Error())
	}
	return nil
}
