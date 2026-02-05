package cmd

import (
	"fmt"
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultOutputDir = "."
)

// NewImportCommand returns the import subcommand.
func NewImportCommand() *cobra.Command {
	var (
		outputDir   string
		include     []string
		exclude     []string
		dryRun      bool
		verbose     bool
		serverURL   string
		orgID       string
		workspaceID string
		cloudDomain string
	)
	v := viper.New()
	cfg := config.NewConfig(v)
	cfg.BindEnv()

	imp := &cobra.Command{
		Use:   "import",
		Short: "Generate Terraform HCL and import blocks from Cribl resources",
		Long:  "Reads resources from Cribl and writes Terraform HCL plus import blocks so you can run terraform import.",
		RunE: func(c *cobra.Command, args []string) error {
			if err := ValidateImportFlags(include, exclude); err != nil {
				return err
			}
			if err := cfg.LoadCredentialsFile(); err != nil {
				return err
			}
			if err := cfg.ValidateRequired(); err != nil {
				return err
			}
			if verbose {
				printResolvedConfig(c, cfg)
			}
			return nil
		},
	}

	imp.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir, "Output directory for generated Terraform")
	imp.Flags().StringSliceVar(&include, "include", nil, "Resource types to include")
	imp.Flags().StringSliceVar(&exclude, "exclude", nil, "Resource types to exclude")
	imp.Flags().BoolVar(&dryRun, "dry-run", false, "Preview resources without generating files")
	imp.Flags().BoolVar(&verbose, "verbose", false, "Enable debug logging")

	imp.Flags().StringVar(&serverURL, "server-url", "", "On-prem base URL")
	imp.Flags().StringVar(&orgID, "org-id", "", "Cribl org identifier")
	imp.Flags().StringVar(&workspaceID, "workspace-id", "", "Workspace identifier")
	imp.Flags().StringVar(&cloudDomain, "cloud-domain", "", "Cloud domain override")
	_ = cfg.BindPFlag(config.KeyOnpremServerURL, imp.Flags().Lookup("server-url"))
	_ = cfg.BindPFlag(config.KeyOrganizationID, imp.Flags().Lookup("org-id"))
	_ = cfg.BindPFlag(config.KeyWorkspaceID, imp.Flags().Lookup("workspace-id"))
	_ = cfg.BindPFlag(config.KeyCloudDomain, imp.Flags().Lookup("cloud-domain"))

	return imp
}

// printResolvedConfig prints the resolved config (no secrets) for verbose mode.
func printResolvedConfig(cmd *cobra.Command, cfg *config.Config) {
	out := cmd.OutOrStderr()
	serverURL := cfg.Get(config.KeyOnpremServerURL)
	if serverURL != "" {
		fmt.Fprintf(out, "server_url: %s (on-prem)\n", serverURL)
	} else {
		orgID := cfg.Get(config.KeyOrganizationID)
		workspaceID := cfg.Get(config.KeyWorkspaceID)
		cloudDomain := cfg.Get(config.KeyCloudDomain)
		if cloudDomain == "" {
			cloudDomain = "cribl.cloud"
		}
		fmt.Fprintf(out, "organization_id: %s, workspace_id: %s, cloud_domain: %s (cloud)\n", orgID, workspaceID, cloudDomain)
	}
	if cfg.Get(config.KeyBearerToken) != "" {
		fmt.Fprintln(out, "auth: bearer token")
	} else if cfg.Get(config.KeyClientID) != "" {
		fmt.Fprintln(out, "auth: client credentials")
	} else if cfg.Get(config.KeyOnpremUsername) != "" {
		fmt.Fprintln(out, "auth: username/password")
	}
}

// ValidateImportFlags returns an error if include and exclude overlap.
func ValidateImportFlags(include, exclude []string) error {
	excludeSet := make(map[string]struct{})
	for _, t := range exclude {
		excludeSet[strings.TrimSpace(t)] = struct{}{}
	}
	for _, t := range include {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, ok := excludeSet[t]; ok {
			return fmt.Errorf("resource type %q cannot be in both --include and --exclude", t)
		}
	}
	return nil
}
