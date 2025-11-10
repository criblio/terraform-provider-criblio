package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	// Initialize Viper for configuration
	initViper()

	rootCmd := &cobra.Command{
		Use:   "criblio-import",
		Short: "Import Cribl configuration to Terraform",
		Long: `A CLI tool to fetch existing Cribl configuration from /bulk/diag/download
and convert YAML configuration files into Terraform modules/resources.

Examples:
  # Basic usage with environment variables
  export CRIBL_BEARER_TOKEN="your-token"
  export CRIBL_WORKSPACE_ID="main"
  export CRIBL_ORGANIZATION_ID="my-org"
  criblio-import --output ./terraform-configs

  # With explicit authentication via flags
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

// initViper initializes Viper for configuration management
// Priority: Flags > Environment Variables > Config File > Defaults
func initViper() {
	// Set environment variable prefix
	viper.SetEnvPrefix("CRIBL")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Set config file name and paths
	viper.SetConfigName("credentials")
	viper.SetConfigType("ini")
	viper.AddConfigPath("$HOME/.cribl")
	viper.AddConfigPath(".")

	// Read config file (ignore errors if file doesn't exist)
	_ = viper.ReadInConfig()
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
	viper.BindPFlag("output", cmd.Flags().Lookup("output"))

	cmd.Flags().StringVarP(&format, "format", "f", "modules", "Output format: resources or modules")
	viper.BindPFlag("format", cmd.Flags().Lookup("format"))

	cmd.Flags().StringSliceVarP(&include, "include", "i", []string{}, "Only include these resource types (comma-separated)")
	viper.BindPFlag("include", cmd.Flags().Lookup("include"))

	cmd.Flags().StringSliceVarP(&exclude, "exclude", "e", []string{}, "Exclude these resource types (comma-separated)")
	viper.BindPFlag("exclude", cmd.Flags().Lookup("exclude"))

	cmd.Flags().StringVarP(&varFile, "var-file", "v", "", "Generate variables.tf for sensitive fields")
	viper.BindPFlag("var-file", cmd.Flags().Lookup("var-file"))

	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview what would be generated without writing files")
	viper.BindPFlag("dry-run", cmd.Flags().Lookup("dry-run"))

	cmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))

	// Authentication flags (with env var support)
	cmd.Flags().StringVar(&workspaceID, "workspace-id", "", "Cribl workspace ID (env: CRIBL_WORKSPACE_ID)")
	viper.BindPFlag("workspace-id", cmd.Flags().Lookup("workspace-id"))

	cmd.Flags().StringVar(&orgID, "organization-id", "", "Cribl organization ID (env: CRIBL_ORGANIZATION_ID)")
	viper.BindPFlag("organization-id", cmd.Flags().Lookup("organization-id"))

	cmd.Flags().StringVar(&bearerToken, "bearer-token", "", "Bearer token for authentication (env: CRIBL_BEARER_TOKEN)")
	viper.BindPFlag("bearer-token", cmd.Flags().Lookup("bearer-token"))

	cmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client ID (env: CRIBL_CLIENT_ID)")
	viper.BindPFlag("client-id", cmd.Flags().Lookup("client-id"))

	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth client secret (env: CRIBL_CLIENT_SECRET)")
	viper.BindPFlag("client-secret", cmd.Flags().Lookup("client-secret"))

	// On-prem flags
	cmd.Flags().StringVar(&onpremURL, "onprem-server-url", "", "On-prem server URL (env: CRIBL_ONPREM_SERVER_URL)")
	viper.BindPFlag("onprem-server-url", cmd.Flags().Lookup("onprem-server-url"))

	cmd.Flags().StringVar(&onpremUser, "onprem-username", "", "On-prem username (env: CRIBL_ONPREM_USERNAME)")
	viper.BindPFlag("onprem-username", cmd.Flags().Lookup("onprem-username"))

	cmd.Flags().StringVar(&onpremPass, "onprem-password", "", "On-prem password (env: CRIBL_ONPREM_PASSWORD)")
	viper.BindPFlag("onprem-password", cmd.Flags().Lookup("onprem-password"))

	// Terraform automation flags
	cmd.Flags().BoolVar(&autoImport, "auto-import", false, "Automatically import generated resources into Terraform state")
	viper.BindPFlag("auto-import", cmd.Flags().Lookup("auto-import"))

	cmd.Flags().BoolVar(&autoApply, "auto-apply", false, "Automatically apply after import (requires --auto-import)")
	viper.BindPFlag("auto-apply", cmd.Flags().Lookup("auto-apply"))

	// Bind command to Viper and read values after flag parsing
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Read values from Viper (env vars/config file) if flags not set
		if !cmd.Flags().Changed("workspace-id") {
			workspaceID = viper.GetString("workspace-id")
		}
		if !cmd.Flags().Changed("organization-id") {
			orgID = viper.GetString("organization-id")
		}
		if !cmd.Flags().Changed("bearer-token") {
			bearerToken = viper.GetString("bearer-token")
		}
		if !cmd.Flags().Changed("client-id") {
			clientID = viper.GetString("client-id")
		}
		if !cmd.Flags().Changed("client-secret") {
			clientSecret = viper.GetString("client-secret")
		}
		if !cmd.Flags().Changed("onprem-server-url") {
			onpremURL = viper.GetString("onprem-server-url")
		}
		if !cmd.Flags().Changed("onprem-username") {
			onpremUser = viper.GetString("onprem-username")
		}
		if !cmd.Flags().Changed("onprem-password") {
			onpremPass = viper.GetString("onprem-password")
		}
		return nil
	}

	return cmd
}
