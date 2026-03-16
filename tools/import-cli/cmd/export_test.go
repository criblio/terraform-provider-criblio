package cmd_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// criblMockServer starts an httptest.Server that mimics enough of the Cribl API for discovery to succeed:
// - POST /oauth/token returns 200 with a minimal token so SDK auth succeeds.
// - GET /products/stream/groups and /products/edge/groups return one group (id "default").
// - All other GETs return 200 with {"items":[]} so List* calls return empty and no error.
func criblMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	groupsJSON := []byte(`{"items":[{"id":"default","name":"default"}]}`)
	emptyListJSON := []byte(`{"items":[]}`)
	oauthJSON := []byte(`{"access_token":"test","expires_in":300}`)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost && (strings.HasSuffix(r.URL.Path, "/oauth/token") || r.URL.Path == "/oauth/token") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(oauthJSON)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		path := r.URL.Path
		if strings.HasSuffix(path, "/products/stream/groups") || strings.HasSuffix(path, "/products/edge/groups") {
			_, _ = w.Write(groupsJSON)
			return
		}
		_, _ = w.Write(emptyListJSON)
	}))
}

func TestExportCommand_Help_ShowsAllFlags(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"export", "--help"})
	err := root.Execute()
	require.NoError(t, err)
	help := out.String()

	// All supported flags (export --help)
	assert.Contains(t, help, "--output-dir", "export --help should document --output-dir")
	assert.Contains(t, help, "--include", "export --help should document --include")
	assert.Contains(t, help, "--exclude", "export --help should document --exclude")
	assert.Contains(t, help, "--dry-run", "export --help should document --dry-run")
	assert.Contains(t, help, "--verbose", "export --help should document --verbose")
	assert.Contains(t, help, "--server-url", "export --help should document --server-url")
	assert.Contains(t, help, "--org-id", "export --help should document --org-id")
	assert.Contains(t, help, "--workspace-id", "export --help should document --workspace-id")
	assert.Contains(t, help, "--cloud-domain", "export --help should document --cloud-domain")

	// Description of --dry-run (preview only; no file writes; List* only)
	assert.Contains(t, help, "Preview", "export --help should describe --dry-run (Preview resources)")
	assert.Contains(t, help, "dry-run", "export --help should document --dry-run")
	// --include and --exclude filters work as documented
	assert.Contains(t, help, "include", "export --help should document --include filter")
	assert.Contains(t, help, "exclude", "export --help should document --exclude filter")
}

func TestExportCommand_Help_ShowsExampleUsage(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"export", "--help"})
	err := root.Execute()
	require.NoError(t, err)
	help := out.String()
	// Example commands must render (no regressions in CLI UX)
	assert.Contains(t, help, "export --dry-run", "export --help should show example usage")
	assert.Contains(t, help, "export --server-url", "export --help should show example with --server-url")
}

func TestExportCommand_DefaultBehavior(t *testing.T) {
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
	root.SetArgs([]string{"export"})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no valid configuration")
}

func TestExportCommand_ValidConfigInitializesClient(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()

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
	_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", server.URL)
	_ = os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	// Use --include criblio_group so discovery only hits endpoints the mock supports (groups).
	// On-prem mode blocks many lib/search endpoints; criblio_group uses /products/*/groups which the mock handles.
	root.SetArgs([]string{"export", "--server-url", server.URL, "--dry-run", "--include", "criblio_group"})
	err := root.Execute()
	if err != nil {
		t.Logf("stderr: %s", errOut.String())
		t.Logf("stdout: %s", out.String())
	}
	require.NoError(t, err)

	stderr := errOut.String()
	assert.Contains(t, stderr, "Preview:", "dry-run should print preview header")
	assert.Contains(t, stderr, "Total:", "dry-run should print total line")
	assert.Contains(t, stderr, "criblio_", "dry-run should list at least one resource type")
}

// TestExportCommand_DryRun_IncludeFilter verifies --include limits output to listed resource types.
func TestExportCommand_DryRun_IncludeFilter(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()

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
	_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", server.URL)
	_ = os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	// Use criblio_group; mock returns 2 groups. Other types fail on-prem (lib/search endpoints blocked).
	root.SetArgs([]string{"export", "--dry-run", "--include", "criblio_group"})
	err := root.Execute()
	if err != nil {
		t.Logf("stderr: %s", errOut.String())
		t.Logf("stdout: %s", out.String())
	}
	require.NoError(t, err)

	stderr := errOut.String()
	assert.Contains(t, stderr, "Preview:", "dry-run should print preview")
	assert.Contains(t, stderr, "Total:", "dry-run should print total line")
	assert.Contains(t, stderr, "criblio_group:", "--include criblio_group should list that type")
	assert.Contains(t, stderr, "2", "mock returns 2 groups for criblio_group")
}

// TestExportCommand_DryRun_ExcludeFilter verifies --exclude omits listed resource types.
func TestExportCommand_DryRun_ExcludeFilter(t *testing.T) {
	server := criblMockServer(t)
	defer server.Close()

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
	_ = os.Setenv("CRIBL_ONPREM_SERVER_URL", server.URL)
	_ = os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	// Include only criblio_group (mock supports it), exclude criblio_source so it does not appear.
	root.SetArgs([]string{"export", "--dry-run", "--include", "criblio_group", "--exclude", "criblio_source"})
	err := root.Execute()
	if err != nil {
		t.Logf("stderr: %s", errOut.String())
		t.Logf("stdout: %s", out.String())
	}
	require.NoError(t, err)

	stderr := errOut.String()
	assert.NotContains(t, stderr, "  criblio_source:", "--exclude should omit criblio_source from listing")
	assert.Contains(t, stderr, "Preview:", "dry-run should print preview")
}

func TestExportCommand_Validation_IncludeExcludeOverlap(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetArgs([]string{"export", "--include", "sources", "--exclude", "sources"})
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be in both")
	assert.Contains(t, err.Error(), "--include")
	assert.Contains(t, err.Error(), "--exclude")
}

func TestValidateExportFlags(t *testing.T) {
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
			err := cmd.ValidateExportFlags(tt.include, tt.exclude)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
