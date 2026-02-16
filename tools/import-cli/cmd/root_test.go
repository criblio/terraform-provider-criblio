package cmd_test

import (
	"bytes"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootCommand(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	require.NotNil(t, root)
	assert.Equal(t, "goatify", root.Use)
	assert.NotEmpty(t, root.Short)
}

func TestRootCommand_SubcommandRegistration(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	subs := root.Commands()
	require.Len(t, subs, 2)
	names := make([]string, len(subs))
	for i, c := range subs {
		names[i] = c.Use
	}
	assert.Contains(t, names, "import")
	assert.Contains(t, names, "version")
}

func TestRootCommand_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string // substrings that must appear in combined stdout+stderr
	}{
		{
			name:     "no args shows help",
			args:     []string{},
			wantErr:  false,
			contains: []string{"Usage", "Available Commands", "import", "version"},
		},
		{
			name:     "root --help",
			args:     []string{"--help"},
			wantErr:  false,
			contains: []string{"Usage", "Available Commands", "import", "version"},
		},
		{
			name:     "import --help",
			args:     []string{"import", "--help"},
			wantErr:  false,
			contains: []string{"import", "Usage"},
		},
		{
			name:     "version prints build info",
			args:     []string{"version"},
			wantErr:  false,
			contains: []string{"version"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := cmd.NewRootCommand()
			out := &bytes.Buffer{}
			errOut := &bytes.Buffer{}
			root.SetOut(out)
			root.SetErr(errOut)
			root.SetArgs(tt.args)
			err := root.Execute()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			combined := out.String() + errOut.String()
			for _, sub := range tt.contains {
				assert.Contains(t, combined, sub, "output should contain %q", sub)
			}
		})
	}
}

// TestRootCommand_HelpIncludesCommandDescriptionsAndExamples verifies goatify --help
// includes command descriptions and example usage (CLI UX stability).
func TestRootCommand_HelpIncludesCommandDescriptionsAndExamples(t *testing.T) {
	t.Parallel()
	root := cmd.NewRootCommand()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"--help"})
	err := root.Execute()
	require.NoError(t, err)
	help := out.String()

	// Command descriptions: root Short and subcommand names
	assert.Contains(t, help, "Export Cribl", "root --help should include command description")
	assert.Contains(t, help, "import", "root --help should list import command")
	assert.Contains(t, help, "version", "root --help should list version command")

	// Example usage must render
	assert.Contains(t, help, "import --dry-run", "root --help should show example usage")
}
