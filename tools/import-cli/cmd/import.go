package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/config"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultOutputDir = "."
)

// logFilterWriter wraps an io.Writer and drops lines containing "[DEBUG]" when suppressDebug is true.
// Used so SDK debug logs only appear when the user passes --verbose.
type logFilterWriter struct {
	w             io.Writer
	suppressDebug bool
}

func (f *logFilterWriter) Write(p []byte) (n int, err error) {
	if f.suppressDebug && strings.Contains(string(p), "[DEBUG]") {
		return len(p), nil
	}
	return f.w.Write(p)
}

// NewImportCommand returns the import subcommand.
func NewImportCommand() *cobra.Command {
	var (
		outputDir   string
		include     []string
		exclude     []string
		group       []string
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

	appName := version.AppNameOrDefault()
	imp := &cobra.Command{
		Use:   "import",
		Short: "Generate Terraform HCL and import blocks from Cribl resources",
		Long:  "Reads resources from Cribl and writes Terraform HCL plus import blocks so you can run terraform import.",
		Example: "  " + appName + " import --dry-run\n  " + appName + " import --server-url https://cribl.example.com --output-dir ./tf",
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
			// Suppress SDK DEBUG logs unless --verbose is set.
			if verbose {
				log.SetOutput(os.Stderr)
			} else {
				log.SetOutput(&logFilterWriter{w: os.Stderr, suppressDebug: true})
			}
			sdkClient, err := client.NewFromConfig(cfg)
			if err != nil {
				return fmt.Errorf("initialize SDK client: %w", err)
			}
			if verbose {
				printResolvedConfig(c, cfg)
			}
			ctx := context.Background()
			reg, err := buildRegistry(ctx)
			if err != nil {
				return fmt.Errorf("build registry: %w", err)
			}
			results, err := discovery.Discover(ctx, sdkClient, reg, include, exclude, group)
			if err != nil {
				return fmt.Errorf("discovery: %w", err)
			}
			if dryRun {
				// Preview only: print resource counts and types. No file writes, no Get*ByID (discovery uses List* only).
				printDryRunPreview(c, results, group)
				return nil
			}
			// Surface SDK errors with resource context; fail if any discovery failed
			var firstErr error
			for _, r := range results {
				if r.Err != nil {
					fmt.Fprintln(c.ErrOrStderr(), r.Err.Error())
					if firstErr == nil {
						firstErr = r.Err
					}
				}
			}
			if firstErr != nil {
				return fmt.Errorf("discovery failed for one or more resource types: %w", firstErr)
			}
			// TODO: generate HCL + import blocks from results, write to outputDir
			_ = outputDir
			return nil
		},
	}

	imp.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir, "Output directory for generated Terraform")
	imp.Flags().StringSliceVar(&include, "include", nil, "Resource types to include (e.g. criblio_source, criblio_pipeline). If set, only these types are discovered; otherwise all discoverable types are used.")
	imp.Flags().StringSliceVar(&exclude, "exclude", nil, "Resource types to exclude (e.g. criblio_notification). Excluded types are omitted from discovery and export.")
	imp.Flags().StringSliceVar(&group, "group", nil, "Restrict discovery and export to these groups only. Use group ID (e.g. default) or label (e.g. 'default (stream)'). Can be repeated. Empty = all groups.")
	imp.Flags().BoolVar(&dryRun, "dry-run", false, "Preview resource counts and types only; no conversion or file writes. Uses List* API only (no Get*ByID).")
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

// printDryRunPreview prints resource types and counts only. Used for --dry-run.
// Discovery uses List* endpoints only; no Get*ByID calls or file writes.
// Writes to stderr so the preview appears in the terminal with verbose/config output.
// When groupFilter is non-empty, only group-scoped types (and criblio_group) are shown.
func printDryRunPreview(cmd *cobra.Command, results []discovery.Result, groupFilter []string) {
	out := cmd.ErrOrStderr()
	// Aggregate by type name (registry may have multiple entries per type)
	type agg struct {
		count          int
		errs           int
		firstErr       string
		details        []string         // e.g. group names for criblio_group
		perGroupCounts map[string]int   // group-scoped resource counts per group
	}
	byType := make(map[string]*agg)
	for _, r := range results {
		// When filtering by group, only include group-scoped types and criblio_group.
		if len(groupFilter) > 0 && r.TypeName != "criblio_group" && len(r.PerGroupCounts) == 0 {
			continue
		}
		a, ok := byType[r.TypeName]
		if !ok {
			a = &agg{firstErr: shortenError(r.Err, 80)}
			byType[r.TypeName] = a
		}
		if r.Err != nil {
			a.errs++
			if a.firstErr == "" {
				a.firstErr = shortenError(r.Err, 80)
			}
		} else {
			a.count += r.Count
		}
		if len(r.Details) > 0 {
			a.details = r.Details
		}
		if len(r.PerGroupCounts) > 0 {
			if a.perGroupCounts == nil {
				a.perGroupCounts = make(map[string]int)
			}
			for k, v := range r.PerGroupCounts {
				a.perGroupCounts[k] += v
			}
		}
	}
	names := make([]string, 0, len(byType))
	for n := range byType {
		names = append(names, n)
	}
	sort.Strings(names)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Preview: resource types and counts that would be exported (--dry-run)")
	fmt.Fprintln(out, "Resource types are filtered by --include and --exclude when set. Groups filtered by --group when set.")
	var totalResources int
	var typesShown int
	var typesWithErr int
	for _, name := range names {
		a := byType[name]
		if a.errs > 0 {
			msg := a.firstErr
			if a.errs > 1 {
				msg = fmt.Sprintf("%s (%d errors)", msg, a.errs)
			}
			fmt.Fprintf(out, "  %s: error: %s\n", name, msg)
			typesShown++
			typesWithErr++
		} else if a.count > 0 {
			fmt.Fprintf(out, "  %s: %d\n", name, a.count)
			totalResources += a.count
			typesShown++
			for _, d := range a.details {
				fmt.Fprintf(out, "    - %s\n", d)
			}
			// Per-group breakdown for group-scoped resources (sorted for stable output)
			if len(a.perGroupCounts) > 0 {
				groupLabels := make([]string, 0, len(a.perGroupCounts))
				for k := range a.perGroupCounts {
					groupLabels = append(groupLabels, k)
				}
				sort.Strings(groupLabels)
				for _, label := range groupLabels {
					fmt.Fprintf(out, "    - %s: %d\n", label, a.perGroupCounts[label])
				}
			}
		}
	}
	fmt.Fprintf(out, "Total: %d resource types", typesShown)
	if typesWithErr > 0 {
		fmt.Fprintf(out, " (%d with errors)", typesWithErr)
	}
	fmt.Fprintf(out, ", %d resources\n", totalResources)
}

// shortenError returns a single-line, truncated error message for user-facing output.
func shortenError(err error, maxLen int) string {
	if err == nil {
		return ""
	}
	s := err.Error()
	if i := strings.IndexAny(s, "\n\r"); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	// Drop common noise for brevity
	if strings.HasPrefix(s, "criblio_") && strings.Contains(s, ": ") {
		if j := strings.Index(s, ": "); j > 0 {
			s = strings.TrimSpace(s[j+2:])
		}
	}
	if len(s) > maxLen {
		s = s[:maxLen-3] + "..."
	}
	return s
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

// buildRegistry builds the resource registry from the provider's Resources().
func buildRegistry(ctx context.Context) (*registry.Registry, error) {
	ver := version.Version
	if ver == "" {
		ver = "dev"
	}
	p := provider.New(ver)()
	constructors := p.Resources(ctx)
	return registry.NewFromResources(ctx, constructors, registry.MetadataFromProvider(), nil)
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
