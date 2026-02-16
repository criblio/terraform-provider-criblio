package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportCommand_Help_ShowsAllFlags(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"import", "--help"})
	err := root.Execute()
	require.NoError(t, err)
	help := out.String()

	// All supported flags (goatify import --help)
	assert.Contains(t, help, "--output-dir", "import --help should document --output-dir")
	assert.Contains(t, help, "--include", "import --help should document --include")
	assert.Contains(t, help, "--exclude", "import --help should document --exclude")
	assert.Contains(t, help, "--dry-run", "import --help should document --dry-run")
	assert.Contains(t, help, "--verbose", "import --help should document --verbose")
	assert.Contains(t, help, "--server-url", "import --help should document --server-url")
	assert.Contains(t, help, "--org-id", "import --help should document --org-id")
	assert.Contains(t, help, "--workspace-id", "import --help should document --workspace-id")
	assert.Contains(t, help, "--cloud-domain", "import --help should document --cloud-domain")

	// Description of --dry-run (preview only; no file writes; List* only)
	assert.Contains(t, help, "Preview", "import --help should describe --dry-run (Preview resources)")
	assert.Contains(t, help, "dry-run", "import --help should document --dry-run")
	// --include and --exclude filters work as documented
	assert.Contains(t, help, "include", "import --help should document --include filter")
	assert.Contains(t, help, "exclude", "import --help should document --exclude filter")
}

func TestImportCommand_Help_ShowsExampleUsage(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"import", "--help"})
	err := root.Execute()
	require.NoError(t, err)
	help := out.String()
	// Example commands must render (no regressions in CLI UX)
	assert.Contains(t, help, "import --dry-run", "import --help should show example usage")
	assert.Contains(t, help, "import --server-url", "import --help should show example with --server-url")
}

func TestImportCommand_DefaultBehavior(t *testing.T) {
	// Run without credentials so we hit validation and never call the API.
	origHome := os.Getenv("HOME")
	origCribl := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", origCribl)
	})
	// Use a temp dir as HOME so ~/.cribl/credentials is not used
	dir := t.TempDir()
	_ = os.Setenv("HOME", dir)
	_ = os.Unsetenv("CRIBL_ONPREM_SERVER_URL")
	_ = os.Unsetenv("CRIBL_BEARER_TOKEN")
	_ = os.Unsetenv("CRIBL_CLIENT_ID")
	_ = os.Unsetenv("CRIBL_CLIENT_SECRET")
	_ = os.Unsetenv("CRIBL_ORGANIZATION_ID")
	_ = os.Unsetenv("CRIBL_WORKSPACE_ID")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs([]string{"import"})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no valid configuration")
}

func TestImportCommand_ValidConfigInitializesClient(t *testing.T) {
	// Successful auth: valid config passes validation and SDK client initializes (no duplicate auth logic).
	// Auth is via env (CRIBL_BEARER_TOKEN) or credentials file; import command has no --bearer-token flag.
	origHome := os.Getenv("HOME")
	origURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	origToken := os.Getenv("CRIBL_BEARER_TOKEN")
	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", origURL)
		_ = os.Setenv("CRIBL_BEARER_TOKEN", origToken)
	})
	dir := t.TempDir()
	_ = os.Setenv("HOME", dir)
	_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", "https://cribl.example.com")
	_ = os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs([]string{"import", "--server-url", "https://cribl.example.com", "--dry-run"})
	err := root.Execute()
	require.NoError(t, err)
	// Dry-run prints resource counts and types to stderr (with verbose/config output).
	stderr := errOut.String()
	assert.Contains(t, stderr, "Preview:", "dry-run should print preview header")
	assert.Contains(t, stderr, "Total:", "dry-run should print total line")
	assert.Contains(t, stderr, "criblio_", "dry-run should list at least one resource type")
}

// TestImportCommand_DryRun_IncludeFilter verifies --include limits output to listed resource types.
func TestImportCommand_DryRun_IncludeFilter(t *testing.T) {
	origHome := os.Getenv("HOME")
	origURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	origToken := os.Getenv("CRIBL_BEARER_TOKEN")
	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", origURL)
		_ = os.Setenv("CRIBL_BEARER_TOKEN", origToken)
	})
	dir := t.TempDir()
	_ = os.Setenv("HOME", dir)
	_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", "https://cribl.example.com")
	_ = os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs([]string{"import", "--dry-run", "--include", "criblio_source", "--include", "criblio_pipeline"})
	err := root.Execute()
	require.NoError(t, err)
	stderr := errOut.String()
	assert.Contains(t, stderr, "criblio_source", "--include should include criblio_source")
	assert.Contains(t, stderr, "criblio_pipeline", "--include should include criblio_pipeline")
	// Only two types should be listed (count or error per line)
	assert.Contains(t, stderr, "Total: 2 resource types", "filter should produce exactly 2 types")
}

// TestImportCommand_DryRun_ExcludeFilter verifies --exclude omits listed resource types.
func TestImportCommand_DryRun_ExcludeFilter(t *testing.T) {
	origHome := os.Getenv("HOME")
	origURL := os.Getenv("CRIBL_ONPREM_SERVER_URL")
	origToken := os.Getenv("CRIBL_BEARER_TOKEN")
	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", origURL)
		_ = os.Setenv("CRIBL_BEARER_TOKEN", origToken)
	})
	dir := t.TempDir()
	_ = os.Setenv("HOME", dir)
	_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", "https://cribl.example.com")
	_ = os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs([]string{"import", "--dry-run", "--exclude", "criblio_source"})
	err := root.Execute()
	require.NoError(t, err)
	stderr := errOut.String()
	// Excluded type must not appear as a listed resource line (e.g. "  criblio_source: ...").
	assert.NotContains(t, stderr, "  criblio_source:", "--exclude should omit criblio_source from listing")
	assert.Contains(t, stderr, "Preview:", "dry-run should still print preview")
}

func TestImportCommand_Validation_IncludeExcludeOverlap(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs([]string{"import", "--include", "sources", "--exclude", "sources"})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be in both")
	assert.Contains(t, err.Error(), "--include")
	assert.Contains(t, err.Error(), "--exclude")
}

func TestValidateImportFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		include []string
		exclude []string
		wantErr bool
	}{
		{"no overlap", []string{"sources"}, []string{"destinations"}, false},
		{"empty both", nil, nil, false},
		{"overlap single", []string{"sources"}, []string{"sources"}, true},
		{"overlap multiple", []string{"a", "b"}, []string{"b", "c"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := cmd.ValidateImportFlags(tt.include, tt.exclude)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
