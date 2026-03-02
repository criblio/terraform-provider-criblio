package generator

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockFS implements FileSystem with in-memory storage for unit tests.
// No real disk I/O; use to mock filesystem interactions.
type mockFS struct {
	mu           sync.Mutex
	MkdirAllCalls []string            // paths passed to MkdirAll
	Files        map[string][]byte    // full path -> content
}

func (m *mockFS) MkdirAll(path string, perm os.FileMode) error {
	m.mu.Lock()
	m.MkdirAllCalls = append(m.MkdirAllCalls, path)
	m.mu.Unlock()
	return nil
}

func (m *mockFS) WriteFileAtomic(dir, name string, data []byte, perm os.FileMode) error {
	path := filepath.Join(dir, name)
	m.mu.Lock()
	if m.Files == nil {
		m.Files = make(map[string][]byte)
	}
	m.Files[path] = append([]byte(nil), data...)
	m.mu.Unlock()
	return nil
}

func TestWriteModuleDirectoryWithFS_mockFS(t *testing.T) {
	mock := &mockFS{}
	baseDir := "/out"
	items := []ResourceItem{
		{
			TypeName: "criblio_source",
			Name:     "first",
			Attrs:    map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "id1"}, "group_id": {Kind: hcl.KindString, String: "default"}},
			ImportID: `{"group_id":"default","id":"id1"}`,
		},
	}
	opts := &hcl.ResourceBlockOptions{SkipNullAttributes: true}

	err := WriteModuleDirectoryWithFS(mock, baseDir, items, opts)
	require.NoError(t, err)

	// MkdirAll was called with the module dir
	mock.mu.Lock()
	mkdirCalls := append([]string(nil), mock.MkdirAllCalls...)
	files := make(map[string][]byte, len(mock.Files))
	for k, v := range mock.Files {
		files[k] = v
	}
	mock.mu.Unlock()

	expectedDir := filepath.Join(baseDir, "source")
	assert.Contains(t, mkdirCalls, expectedDir, "MkdirAll should be called with module path")

	// main.tf, versions.tf, variables.tf, outputs.tf written (import.tf is at root only)
	mainPath := filepath.Join(expectedDir, "main.tf")
	versionsPath := filepath.Join(expectedDir, "versions.tf")
	variablesPath := filepath.Join(expectedDir, "variables.tf")
	outputsPath := filepath.Join(expectedDir, "outputs.tf")

	assert.Contains(t, files, mainPath, "main.tf should be written")
	assert.Contains(t, files, versionsPath, "versions.tf should be written")
	assert.Contains(t, files, variablesPath, "variables.tf should be written")
	assert.Contains(t, files, outputsPath, "outputs.tf should be written")

	versionsContent := string(files[versionsPath])
	assert.Contains(t, versionsContent, "criblio/criblio")
	assert.Contains(t, versionsContent, "required_providers")

	mainContent := string(files[mainPath])
	assert.Contains(t, mainContent, `resource "criblio_source" "first"`)
	assert.Contains(t, mainContent, "id1")

	variablesContent := string(files[variablesPath])
	assert.Contains(t, variablesContent, "Variables for criblio_source")

	outputsContent := string(files[outputsPath])
	assert.Contains(t, outputsContent, "Outputs for criblio_source")
}

func TestOsFS_WriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	fs := &osFS{}
	err := fs.WriteFileAtomic(dir, "test.txt", []byte("hello"), 0644)
	require.NoError(t, err)
	content, err := os.ReadFile(filepath.Join(dir, "test.txt"))
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), content)
}
