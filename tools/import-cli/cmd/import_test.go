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
	assert.Contains(t, help, "--output-dir")
	assert.Contains(t, help, "--include")
	assert.Contains(t, help, "--exclude")
	assert.Contains(t, help, "--dry-run")
	assert.Contains(t, help, "--verbose")
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
