package cmd

import (
	"fmt"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand returns the version subcommand.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version and build information",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Fprintf(c.OutOrStdout(), "%s version %s (commit %s, date %s)\n",
				version.AppNameOrDefault(),
				orUnknown(version.Version),
				orUnknown(version.Commit),
				orUnknown(version.Date))
			return nil
		},
	}
}

func orUnknown(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
