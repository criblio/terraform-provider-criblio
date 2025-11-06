package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "criblio-import",
		Short: "Import Cribl configuration to Terraform",
		Long: `A CLI tool to fetch existing Cribl configuration from /bulk/diag/download
and convert YAML configuration files into Terraform modules/resources.

Examples:
  # Basic usage with environment variables
  criblio-import --output ./terraform-configs

  # With explicit authentication
  criblio-import --output ./configs \
    --bearer-token $TOKEN \
    --workspace-id main \
    --organization-id my-org

  # Import only specific resource types
  criblio-import --output ./configs --include sources,destinations

  # Preview what would be generated
  criblio-import --output ./preview --dry-run`,
		Version: fmt.Sprintf("%s (commit: %s)", version, commit),
	}

	// Add subcommands and flags
	rootCmd.AddCommand(
		newImportCommand(),
		newVersionCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("criblio-import version %s (commit: %s)\n", version, commit)
		},
	}
}

func newImportCommand() *cobra.Command {
	var (
		outputDir    string
		format       string
		include      []string
		exclude      []string
		varFile      string
		dryRun       bool
		verbose      bool
		workspaceID  string
		orgID        string
		bearerToken  string
		clientID     string
		clientSecret string
		onpremURL    string
		onpremUser   string
		onpremPass   string
		autoImport   bool
		autoApply    bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import Cribl configuration to Terraform",
		Long: `Fetch configuration from Cribl API and convert to Terraform modules.

This command will:
1. Authenticate with Cribl API
2. Download configuration bundle from /bulk/diag/download
3. Parse YAML configuration files
4. Convert to Terraform resource definitions
5. Generate organized Terraform modules

The output directory will contain a ready-to-use Terraform configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement import logic
			// This is a skeleton - actual implementation goes here

			if outputDir == "" {
				return fmt.Errorf("output directory is required (use --output)")
			}

			if dryRun {
				fmt.Println("DRY RUN MODE - No files will be written")
				fmt.Printf("Would output to: %s\n", outputDir)
				fmt.Printf("Format: %s\n", format)
				if len(include) > 0 {
					fmt.Printf("Include: %v\n", include)
				}
				if len(exclude) > 0 {
					fmt.Printf("Exclude: %v\n", exclude)
				}
				return nil
			}

			fmt.Printf("Importing configuration to %s...\n", outputDir)
			fmt.Println("⚠️  This is a placeholder implementation")
			fmt.Println("See CLI_TOOL_DESIGN.md for implementation details")

			// Implementation would go here:
			// 1. Setup authentication
			// 2. Call API /bulk/diag/download
			// 3. Extract archive
			// 4. Parse YAML files
			// 5. Convert to Terraform
			// 6. Write output files

			// If auto-import flag is set, import resources after generation
			if autoImport {
				fmt.Println("\n🔄 Auto-import mode enabled")
				fmt.Println("This feature will be implemented to automatically import generated resources")
				fmt.Println("See internal/terraform/importer.go for the implementation")
			}

			return nil
		},
	}

	// Flags
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for Terraform files (required)")
	cmd.MarkFlagRequired("output")

	cmd.Flags().StringVarP(&format, "format", "f", "modules", "Output format: resources or modules")
	cmd.Flags().StringSliceVarP(&include, "include", "i", []string{}, "Only include these resource types (comma-separated)")
	cmd.Flags().StringSliceVarP(&exclude, "exclude", "e", []string{}, "Exclude these resource types (comma-separated)")
	cmd.Flags().StringVarP(&varFile, "var-file", "v", "", "Generate variables.tf for sensitive fields")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview what would be generated without writing files")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	// Authentication flags
	cmd.Flags().StringVar(&workspaceID, "workspace-id", "", "Cribl workspace ID")
	cmd.Flags().StringVar(&orgID, "organization-id", "", "Cribl organization ID")
	cmd.Flags().StringVar(&bearerToken, "bearer-token", "", "Bearer token for authentication")
	cmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth client secret")

	// On-prem flags
	cmd.Flags().StringVar(&onpremURL, "onprem-server-url", "", "On-prem server URL")
	cmd.Flags().StringVar(&onpremUser, "onprem-username", "", "On-prem username")
	cmd.Flags().StringVar(&onpremPass, "onprem-password", "", "On-prem password")

	// Terraform automation flags
	cmd.Flags().BoolVar(&autoImport, "auto-import", false, "Automatically import generated resources into Terraform state")
	cmd.Flags().BoolVar(&autoApply, "auto-apply", false, "Automatically apply after import (requires --auto-import)")

	return cmd
}
