package cmd

import (
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/version"
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root command for the CLI.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   version.AppNameOrDefault(),
		Short: "Export Cribl config to Terraform HCL and generate import blocks",
	}
	root.AddCommand(NewImportCommand())
	root.AddCommand(NewVersionCommand())
	return root
}
