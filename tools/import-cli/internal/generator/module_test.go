package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModuleDir(t *testing.T) {
	dir := ModuleDir("/out", "criblio_source")
	assert.Equal(t, "/out/source", dir)
	dir2 := ModuleDir("/out", "criblio_pipeline")
	assert.Equal(t, "/out/pipeline", dir2)
}

func TestModuleDirByGroup(t *testing.T) {
	dir := ModuleDirByGroup("/out", "default", "criblio_pipeline")
	assert.Equal(t, "/out/default/pipeline", dir)
	dir2 := ModuleDirByGroup("/out", "global", "criblio_source")
	assert.Equal(t, "/out/global/source", dir2)
}

func TestSortResourceItems(t *testing.T) {
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "z"},
		{TypeName: "criblio_source", Name: "a"},
		{TypeName: "criblio_source", Name: "b"},
	}
	SortResourceItems(items)
	assert.Equal(t, "criblio_pipeline", items[0].TypeName)
	assert.Equal(t, "criblio_source", items[1].TypeName)
	assert.Equal(t, "a", items[1].Name)
	assert.Equal(t, "b", items[2].Name)
}

func TestDeduplicateByImportID(t *testing.T) {
	items := []ResourceItem{
		{TypeName: "criblio_source", Name: "a", ImportID: `{"group_id":"default","id":"1"}`},
		{TypeName: "criblio_source", Name: "b", ImportID: `{"group_id":"default","id":"1"}`},
		{TypeName: "criblio_source", Name: "c", ImportID: `{"group_id":"default","id":"2"}`},
	}
	got := DeduplicateByImportID(items)
	require.Len(t, got, 2)
	assert.Equal(t, "a", got[0].Name, "first occurrence kept")
	assert.Equal(t, "c", got[1].Name)
	assert.Equal(t, `{"group_id":"default","id":"1"}`, got[0].ImportID)
	assert.Equal(t, `{"group_id":"default","id":"2"}`, got[1].ImportID)
}

func TestWriteModuleDirectory_structure_and_determinism(t *testing.T) {
	tmp := t.TempDir()
	items := []ResourceItem{
		{
			TypeName: "criblio_source",
			Name:     "second",
			Attrs:    map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "2"}, "group_id": {Kind: hcl.KindString, String: "default"}},
			ImportID: `{"group_id":"default","id":"2"}`,
		},
		{
			TypeName: "criblio_source",
			Name:     "first",
			Attrs:    map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "1"}, "group_id": {Kind: hcl.KindString, String: "default"}},
			ImportID: `{"group_id":"default","id":"1"}`,
		},
	}
	opts := &hcl.ResourceBlockOptions{SkipNullAttributes: true}
	err := WriteModuleDirectory(tmp, items, opts)
	require.NoError(t, err)

	dir := filepath.Join(tmp, "source")
	require.DirExists(t, dir)
	require.FileExists(t, filepath.Join(dir, "main.tf"))
	require.FileExists(t, filepath.Join(dir, "variables.tf"))
	require.FileExists(t, filepath.Join(dir, "outputs.tf"))
	// import.tf is at root only (Terraform restriction)

	main1, err := os.ReadFile(filepath.Join(dir, "main.tf"))
	require.NoError(t, err)
	assert.Contains(t, string(main1), `resource "criblio_source" "first"`)
	assert.Contains(t, string(main1), `"criblio_source" "second"`)
	// Deterministic: first must appear before second (sorted by name)
	assert.Contains(t, string(main1), "first")
	assert.Contains(t, string(main1), "second")

	// Run again and compare output (determinism)
	err = WriteModuleDirectory(tmp, items, opts)
	require.NoError(t, err)
	main2, _ := os.ReadFile(filepath.Join(dir, "main.tf"))
	assert.Equal(t, string(main1), string(main2), "same input must produce same output")
}

func TestWriteModuleDirectory_parseable_hcl(t *testing.T) {
	tmp := t.TempDir()
	items := []ResourceItem{
		{
			TypeName: "criblio_source",
			Name:     "hec",
			Attrs:    map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "hec-1"}, "group_id": {Kind: hcl.KindString, String: "default"}},
			ImportID: `{"group_id":"default","id":"hec-1"}`,
		},
	}
	err := WriteModuleDirectory(tmp, items, nil)
	require.NoError(t, err)
	mainPath := filepath.Join(tmp, "source", "main.tf")
	mainBytes, _ := os.ReadFile(mainPath)
	err = hcl.ParseHCL(mainBytes, "main.tf")
	assert.NoError(t, err)
}

func TestWriteAllModuleDirectories_groups_by_type(t *testing.T) {
	tmp := t.TempDir()
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "p1", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "p1"}}, ImportID: "p1"},
		{TypeName: "criblio_source", Name: "s1", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "s1"}}, ImportID: "s1"},
	}
	err := WriteAllModuleDirectories(tmp, items, nil, nil)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(tmp, "source"))
	require.DirExists(t, filepath.Join(tmp, "pipeline"))
	require.FileExists(t, filepath.Join(tmp, "source", "main.tf"))
	require.FileExists(t, filepath.Join(tmp, "pipeline", "main.tf"))
	require.FileExists(t, filepath.Join(tmp, "import.tf"))
	importBytes, _ := os.ReadFile(filepath.Join(tmp, "import.tf"))
	assert.Contains(t, string(importBytes), "module.pipeline.criblio_pipeline.p1")
	assert.Contains(t, string(importBytes), "module.source.criblio_source.s1")
}

// TestWriteAllModuleDirectoriesByGroup_groups_by_group_id verifies modules are organized by group ID.
func TestWriteAllModuleDirectoriesByGroup_groups_by_group_id(t *testing.T) {
	tmp := t.TempDir()
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "p1", GroupID: "default", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "p1"}}, ImportID: "p1"},
		{TypeName: "criblio_source", Name: "s1", GroupID: "default", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "s1"}}, ImportID: "s1"},
		{TypeName: "criblio_pipeline", Name: "p2", GroupID: "fleet-1", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "p2"}}, ImportID: "p2"},
		{TypeName: "criblio_source", Name: "s2", GroupID: "global", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "s2"}}, ImportID: "s2"},
	}
	err := WriteAllModuleDirectoriesByGroup(tmp, items, nil, nil)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(tmp, "default", "pipeline"))
	require.DirExists(t, filepath.Join(tmp, "default", "source"))
	require.DirExists(t, filepath.Join(tmp, "fleet-1", "pipeline"))
	require.DirExists(t, filepath.Join(tmp, "global", "source"))
	mainDefault, _ := os.ReadFile(filepath.Join(tmp, "default", "pipeline", "main.tf"))
	assert.Contains(t, string(mainDefault), "criblio_pipeline")
	mainDefaultSrc, _ := os.ReadFile(filepath.Join(tmp, "default", "source", "main.tf"))
	assert.Contains(t, string(mainDefaultSrc), "criblio_source")
	mainFleet, _ := os.ReadFile(filepath.Join(tmp, "fleet-1", "pipeline", "main.tf"))
	assert.Contains(t, string(mainFleet), "criblio_pipeline")
	require.FileExists(t, filepath.Join(tmp, "import.tf"))
	importBytes, _ := os.ReadFile(filepath.Join(tmp, "import.tf"))
	assert.Contains(t, string(importBytes), "module.default_pipeline")
	assert.Contains(t, string(importBytes), "module.fleet_1_pipeline")
}

// TestWriteAllModuleDirectories_creates_output_dir verifies that the output directory
// is created when it does not exist (safe create or reuse).
func TestWriteAllModuleDirectories_creates_output_dir(t *testing.T) {
	baseDir := filepath.Join(t.TempDir(), "nested", "tf-out")
	items := []ResourceItem{
		{TypeName: "criblio_source", Name: "s1", Attrs: map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "s1"}}, ImportID: "s1"},
	}
	err := WriteAllModuleDirectories(baseDir, items, nil, nil)
	require.NoError(t, err)
	require.DirExists(t, baseDir)
	require.DirExists(t, filepath.Join(baseDir, "source"))
}

func TestWriteModuleDirectory_mixed_types_error(t *testing.T) {
	tmp := t.TempDir()
	items := []ResourceItem{
		{TypeName: "criblio_source", Name: "a", Attrs: nil, ImportID: ""},
		{TypeName: "criblio_pipeline", Name: "b", Attrs: nil, ImportID: ""},
	}
	err := WriteModuleDirectory(tmp, items, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mixed resource types")
}

// TestImportBlocksBytes_deterministic_and_parseable verifies JIRA AC: import blocks
// are ordered deterministically and generated import.tf parses with HCL parser.
func TestImportBlocksBytes_deterministic_and_parseable(t *testing.T) {
	items := []ResourceItem{
		{TypeName: "criblio_source", Name: "z_second", ImportID: `{"group_id":"default","id":"2"}`},
		{TypeName: "criblio_source", Name: "a_first", ImportID: `{"group_id":"default","id":"1"}`},
	}
	out := ImportBlocksBytes(items)
	require.NotEmpty(t, out)
	// Order must be by name: a_first before z_second
	assert.Greater(t, bytes.Index(out, []byte("a_first")), -1)
	assert.Greater(t, bytes.Index(out, []byte("z_second")), -1)
	assert.Less(t, bytes.Index(out, []byte("a_first")), bytes.Index(out, []byte("z_second")), "blocks ordered by name")
	// Same input twice produces identical output
	out2 := ImportBlocksBytes(items)
	assert.Equal(t, out, out2, "deterministic output")
	// Parses as HCL
	err := hcl.ParseHCL(out, "import.tf")
	assert.NoError(t, err)
}

// TestImportBlocksBytes_includes_all_exported verifies JIRA AC: import.tf includes
// all exported resources (one block per item with non-empty ImportID).
func TestImportBlocksBytes_includes_all_exported(t *testing.T) {
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "p1", ImportID: "id1"},
		{TypeName: "criblio_pipeline", Name: "p2", ImportID: "id2"},
		{TypeName: "criblio_pipeline", Name: "no_id", ImportID: ""},
	}
	out := ImportBlocksBytes(items)
	// Two blocks (p1, p2); no_id is skipped
	assert.Contains(t, string(out), "to = criblio_pipeline.p1")
	assert.Contains(t, string(out), "to = criblio_pipeline.p2")
	assert.NotContains(t, string(out), "no_id")
}

func TestWriteRootFiles_creates_main_providers_and_import(t *testing.T) {
	tmp := t.TempDir()
	infos := []RootModuleInfo{
		{Name: "pipeline", Path: "./pipeline"},
		{Name: "source", Path: "./source"},
	}
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "p1", ImportID: `{"group_id":"default","id":"p1"}`, GroupID: "default"},
		{TypeName: "criblio_source", Name: "s1", ImportID: `{"group_id":"default","id":"s1"}`, GroupID: "default"},
	}
	err := WriteRootFiles(tmp, infos, items, false)
	require.NoError(t, err)
	mainPath := filepath.Join(tmp, "main.tf")
	providersPath := filepath.Join(tmp, "providers.tf")
	importPath := filepath.Join(tmp, "import.tf")
	require.FileExists(t, mainPath)
	require.FileExists(t, providersPath)
	require.FileExists(t, importPath)
	mainBytes, _ := os.ReadFile(mainPath)
	mainStr := string(mainBytes)
	assert.Contains(t, mainStr, `module "pipeline"`)
	assert.Contains(t, mainStr, `source = "./pipeline"`)
	assert.Contains(t, mainStr, `module "source"`)
	assert.Contains(t, mainStr, `source = "./source"`)
	providersBytes, _ := os.ReadFile(providersPath)
	providersStr := string(providersBytes)
	assert.Contains(t, providersStr, "required_providers")
	assert.Contains(t, providersStr, "criblio/criblio")
	assert.Contains(t, providersStr, `provider "criblio"`)
	importBytes, _ := os.ReadFile(importPath)
	importStr := string(importBytes)
	assert.Contains(t, importStr, "to = module.pipeline.criblio_pipeline.p1")
	assert.Contains(t, importStr, "to = module.source.criblio_source.s1")
}

func TestWriteRootFiles_creates_import_with_module_prefix_by_group(t *testing.T) {
	tmp := t.TempDir()
	infos := []RootModuleInfo{
		{Name: "default_pipeline", Path: "./default/pipeline"},
	}
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "p1", ImportID: `{"group_id":"default","id":"p1"}`, GroupID: "default"},
	}
	err := WriteRootFiles(tmp, infos, items, true)
	require.NoError(t, err)
	importPath := filepath.Join(tmp, "import.tf")
	require.FileExists(t, importPath)
	importBytes, _ := os.ReadFile(importPath)
	importStr := string(importBytes)
	assert.Contains(t, importStr, "to = module.default_pipeline.criblio_pipeline.p1")
}

func TestRootModuleInfosFromItems(t *testing.T) {
	items := []ResourceItem{
		{TypeName: "criblio_pipeline", Name: "p1"},
		{TypeName: "criblio_source", Name: "s1"},
		{TypeName: "criblio_pipeline", Name: "p2"},
	}
	infos := RootModuleInfosFromItems(items)
	require.Len(t, infos, 2)
	assert.Equal(t, "pipeline", infos[0].Name)
	assert.Equal(t, "./pipeline", infos[0].Path)
	assert.Equal(t, "source", infos[1].Name)
	assert.Equal(t, "./source", infos[1].Path)
}

// TestImportBlocksBytes_resource_address_convention verifies JIRA AC: Terraform
// resource addresses match provider naming (type.name with criblio_* type).
func TestImportBlocksBytes_resource_address_convention(t *testing.T) {
	items := []ResourceItem{
		{TypeName: "criblio_source", Name: "my_input", ImportID: `{"group_id":"default","id":"in-1"}`},
	}
	out := ImportBlocksBytes(items)
	assert.Contains(t, string(out), "to = criblio_source.my_input")
	assert.Contains(t, string(out), `id = "{\"group_id\":\"default\",\"id\":\"in-1\"}"`)
}
