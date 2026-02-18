package cmd

import (
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/version"
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root command for the CLI.
func NewRootCommand() *cobra.Command {
	appName := version.AppNameOrDefault()
	root := &cobra.Command{
		Use:   appName,
		Short: "Export Cribl config to Terraform HCL and generate import blocks",
		Long:  "Export Cribl configuration to Terraform HCL and generate import blocks so you can run terraform import. Supports Cribl Cloud and on-prem; authentication via environment variables or credentials file.",
		Example: "  " + appName + " import --dry-run\n  " + appName + " import --output-dir ./tf",
	}
	root.AddCommand(NewImportCommand())
	root.AddCommand(NewVersionCommand())
	return root
}
