package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// ResourceItem is a single resource to be written into a module (same resource type).
type ResourceItem struct {
	TypeName    string
	Name        string
	Attrs       map[string]hcl.Value
	ImportID    string
	// GroupID is the worker group/fleet/default_search id; used to organize output under <output_dir>/<id>/resources/.
	// Empty or "global" for resources without group scope.
	GroupID string
	// LifecycleIgnoreChanges, when non-nil, adds lifecycle { ignore_changes = [...] } to the resource.
	// Used for criblio_group_system_settings group_id=default (cloud disables API host/port updates).
	LifecycleIgnoreChanges []string
}

// SortResourceItems sorts items by TypeName then Name for deterministic output.
func SortResourceItems(items []ResourceItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].TypeName != items[j].TypeName {
			return items[i].TypeName < items[j].TypeName
		}
		return items[i].Name < items[j].Name
	})
}

// ImportBlocksBytes returns the contents of import.tf for the given items.
// Blocks are ordered deterministically by TypeName then Name. Items with empty
// ImportID are skipped. Use for testing and for writing import.tf.
func ImportBlocksBytes(items []ResourceItem) []byte {
	copied := make([]ResourceItem, len(items))
	copy(copied, items)
	SortResourceItems(copied)
	buf := new(bytes.Buffer)
	for _, it := range copied {
		if it.ImportID == "" {
			continue
		}
		fmt.Fprintf(buf, "import {\n  to = %s.%s\n  id = %q\n}\n\n", it.TypeName, it.Name, it.ImportID)
	}
	return buf.Bytes()
}

// SecretVariablesBytes returns the contents of variables.tf for secret variables
// referenced in the given items (from ReplaceSecretValuesWithVariableRefs). Each
// variable is declared with sensitive = true.
func SecretVariablesBytes(items []ResourceItem) []byte {
	seen := make(map[string]bool)
	var names []string
	for _, it := range items {
		for _, n := range hcl.CollectSecretVariableNames(it.Attrs) {
			if !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		typeName := ""
		if len(items) > 0 {
			typeName = items[0].TypeName
		}
		return []byte("# Variables for " + typeName + " module\n# Add variable blocks as needed.\n")
	}
	buf := new(bytes.Buffer)
	buf.WriteString("# Secret variables referenced by imported resources.\n# Set these (e.g. via environment or a secret store) before plan/apply.\n\n")
	for _, n := range names {
		fmt.Fprintf(buf, "variable %q {\n  type      = string\n  sensitive = true\n}\n\n", n)
	}
	return buf.Bytes()
}

// typeNameToFolderName converts criblio_pipeline -> pipeline, criblio_source -> source, etc.
func typeNameToFolderName(typeName string) string {
	if strings.HasPrefix(typeName, "criblio_") {
		return typeName[8:]
	}
	return typeName
}

// ModuleDir returns the output directory path for a resource type: <base_dir>/<folder_name>/ (e.g. pipeline, source).
func ModuleDir(baseDir, resourceType string) string {
	return filepath.Join(baseDir, typeNameToFolderName(resourceType))
}

// ModuleDirByGroup returns the output directory when grouping by worker group: <base_dir>/<group_id>/<folder_name>/.
func ModuleDirByGroup(baseDir, groupID, resourceType string) string {
	return filepath.Join(baseDir, groupID, typeNameToFolderName(resourceType))
}

// WriteModuleDirectory creates <base_dir>/<resource_type>/ and writes main.tf, variables.tf, outputs.tf.
// Items must all have the same TypeName. Order is deterministic (sorted by Name).
// Uses DefaultFS (atomic writes). For tests with a mock filesystem, use WriteModuleDirectoryWithFS.
func WriteModuleDirectory(baseDir string, items []ResourceItem, opts *hcl.ResourceBlockOptions) error {
	return WriteModuleDirectoryWithFS(DefaultFS, baseDir, items, opts)
}

// DeduplicateByImportID returns a slice with at most one item per ImportID (same type).
// Keeps the first occurrence so main.tf and root import.tf do not get duplicate resource blocks.
// Items with empty ImportID are kept as-is (e.g. single resources per type).
func DeduplicateByImportID(items []ResourceItem) []ResourceItem {
	if len(items) == 0 {
		return items
	}
	seen := make(map[string]bool)
	out := make([]ResourceItem, 0, len(items))
	for _, it := range items {
		key := it.ImportID
		if key == "" {
			out = append(out, it)
			continue
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, it)
	}
	return out
}

// EnsureUniqueNames deduplicates resource names in place by appending _2, _3, ... when a name
// repeats (e.g. after truncation to MaxResourceNameLength). Call before writing so main.tf and
// import.tf use the same unique names.
func EnsureUniqueNames(items []ResourceItem) {
	seen := make(map[string]int)
	for i := range items {
		name := items[i].Name
		n := seen[name]
		seen[name]++
		if n > 0 {
			// Name collision: append _2, _3, ... (trim if needed to stay under max length)
			suffix := fmt.Sprintf("_%d", n+1)
			maxBase := MaxResourceNameLength - len(suffix)
			if maxBase < 1 {
				maxBase = 1
			}
			base := name
			if len(base) > maxBase {
				base = base[:maxBase]
				base = strings.TrimRight(base, "_")
			}
			items[i].Name = base + suffix
		}
	}
}

// WriteModuleDirectoryWithFS is like WriteModuleDirectory but uses the given FileSystem.
// Use a mock FS in unit tests to avoid real disk I/O.
func WriteModuleDirectoryWithFS(fs FileSystem, baseDir string, items []ResourceItem, opts *hcl.ResourceBlockOptions) error {
	return WriteModuleDirectoryWithFSAndGroup(fs, baseDir, items, opts, "")
}

// WriteModuleDirectoryWithFSAndGroup writes a type directory; when groupID is non-empty, uses <base_dir>/<groupID>/<type>/.
func WriteModuleDirectoryWithFSAndGroup(fs FileSystem, baseDir string, items []ResourceItem, opts *hcl.ResourceBlockOptions, groupID string) error {
	if len(items) == 0 {
		return nil
	}
	typeName := items[0].TypeName
	for _, it := range items {
		if it.TypeName != typeName {
			return fmt.Errorf("mixed resource types in same module: %q vs %q", typeName, it.TypeName)
		}
	}
	items = DeduplicateByImportID(items)
	if len(items) == 0 {
		return nil
	}
	EnsureUniqueNames(items)
	SortResourceItems(items)
	var dir string
	if groupID != "" {
		dir = ModuleDirByGroup(baseDir, groupID, typeName)
	} else {
		dir = ModuleDir(baseDir, typeName)
	}
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create module dir: %w", err)
	}

	// main.tf: resource blocks (order by name)
	resources := make([]hcl.ResourceInput, len(items))
	for i, it := range items {
		resources[i] = hcl.ResourceInput{
			TypeName:               it.TypeName,
			Name:                   it.Name,
			Attrs:                  it.Attrs,
			LifecycleIgnoreChanges: it.LifecycleIgnoreChanges,
		}
	}
	f, err := hcl.FileWithResources(resources, opts)
	if err != nil {
		return fmt.Errorf("main.tf: %w", err)
	}
	if err := fs.WriteFileAtomic(dir, "main.tf", f.Bytes(), 0644); err != nil {
		return err
	}

	// Import blocks are written at root only (Terraform allows import only in root module).
	// Remove stale import.tf from module dir if present (from prior runs).
	_ = os.Remove(filepath.Join(dir, "import.tf"))

	// versions.tf: required_providers so child modules explicitly declare criblio/criblio
	// (avoids Terraform inferring wrong provider address like hashicorp/criblio with dev_overrides)
	versionsContent := []byte(`terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = ">= 1.0"
    }
  }
}
`)
	if err := fs.WriteFileAtomic(dir, "versions.tf", versionsContent, 0644); err != nil {
		return fmt.Errorf("write versions.tf: %w", err)
	}

	// variables.tf: secret variables referenced in resources, plus comment if none
	variablesContent := SecretVariablesBytes(items)
	if err := fs.WriteFileAtomic(dir, "variables.tf", variablesContent, 0644); err != nil {
		return err
	}

	// outputs.tf: minimal placeholder
	outputsContent := []byte("# Outputs for " + typeName + " module\n# Add output blocks as needed.\n")
	if err := fs.WriteFileAtomic(dir, "outputs.tf", outputsContent, 0644); err != nil {
		return err
	}
	return nil
}

// ProgressFunc reports progress to the user; nil means no progress output.
type ProgressFunc func(format string, args ...interface{})

// WriteAllModuleDirectories groups items by TypeName and writes each group to <base_dir>/<type>/.
// Creates baseDir (or reuses it) before writing. Uses DefaultFS (atomic writes).
func WriteAllModuleDirectories(baseDir string, items []ResourceItem, opts *hcl.ResourceBlockOptions, progress ProgressFunc) error {
	return WriteAllModuleDirectoriesWithLayout(baseDir, items, opts, false, progress)
}

// WriteAllModuleDirectoriesByGroup groups items by GroupID and TypeName, writing each to <base_dir>/<group_id>/<type>/.
func WriteAllModuleDirectoriesByGroup(baseDir string, items []ResourceItem, opts *hcl.ResourceBlockOptions, progress ProgressFunc) error {
	return WriteAllModuleDirectoriesWithLayout(baseDir, items, opts, true, progress)
}

// processModuleItems applies DeduplicateByImportID and EnsureUniqueNames.
// Returns the processed slice to use for both module writes and root import blocks.
// Must match the processing done in WriteModuleDirectoryWithFSAndGroup.
func processModuleItems(items []ResourceItem) []ResourceItem {
	items = DeduplicateByImportID(items)
	if len(items) == 0 {
		return nil
	}
	EnsureUniqueNames(items)
	return items
}

// WriteAllModuleDirectoriesWithLayout writes output; byGroup true uses <group_id>/<type>/, false uses <type>/.
func WriteAllModuleDirectoriesWithLayout(baseDir string, items []ResourceItem, opts *hcl.ResourceBlockOptions, byGroup bool, progress ProgressFunc) error {
	if err := DefaultFS.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	SortResourceItems(items)
	if byGroup {
		// Remove stale flat directories from previous --flat runs
		removeStaleFlatDirs(baseDir, items)
		byGroupAndType := make(map[string]map[string][]ResourceItem)
		for _, it := range items {
			gid := it.GroupID
			if gid == "" {
				gid = "global"
			}
			if byGroupAndType[gid] == nil {
				byGroupAndType[gid] = make(map[string][]ResourceItem)
			}
			byGroupAndType[gid][it.TypeName] = append(byGroupAndType[gid][it.TypeName], it)
		}
		groupIDs := make([]string, 0, len(byGroupAndType))
		for k := range byGroupAndType {
			groupIDs = append(groupIDs, k)
		}
		sort.Strings(groupIDs)
		var rootItems []ResourceItem
		for _, gid := range groupIDs {
			typeNames := make([]string, 0, len(byGroupAndType[gid]))
			for tn := range byGroupAndType[gid] {
				typeNames = append(typeNames, tn)
			}
			sort.Strings(typeNames)
			for _, typeName := range typeNames {
				moduleItems := byGroupAndType[gid][typeName]
				processed := processModuleItems(moduleItems)
				if len(processed) == 0 {
					continue
				}
				if progress != nil {
					progress("%s/%s: %d resources", gid, typeNameToFolderName(typeName), len(processed))
				}
				if err := WriteGroupTypeDirectory(baseDir, gid, processed, opts); err != nil {
					return fmt.Errorf("group %s/%s: %w", gid, typeNameToFolderName(typeName), err)
				}
				rootItems = append(rootItems, processed...)
			}
		}
		if err := WriteRootFiles(baseDir, RootModuleInfosFromItemsByGroup(rootItems), rootItems, true); err != nil {
			return fmt.Errorf("write root files: %w", err)
		}
		return nil
	}
	byType := make(map[string][]ResourceItem)
	for _, it := range items {
		byType[it.TypeName] = append(byType[it.TypeName], it)
	}
	typeNames := make([]string, 0, len(byType))
	for k := range byType {
		typeNames = append(typeNames, k)
	}
	sort.Strings(typeNames)
	var rootItems []ResourceItem
	for _, typeName := range typeNames {
		moduleItems := byType[typeName]
		processed := processModuleItems(moduleItems)
		if len(processed) == 0 {
			continue
		}
		if progress != nil {
			progress("%s: %d resources", typeNameToFolderName(typeName), len(processed))
		}
		if err := WriteModuleDirectory(baseDir, processed, opts); err != nil {
			return fmt.Errorf("module %s: %w", typeName, err)
		}
		rootItems = append(rootItems, processed...)
	}
	if err := WriteRootFiles(baseDir, RootModuleInfosFromItems(rootItems), rootItems, false); err != nil {
		return fmt.Errorf("write root files: %w", err)
	}
	return nil
}

// WriteGroupTypeDirectory writes a single type directory under a group: <base_dir>/<groupID>/<type>/.
// Items must all have the same TypeName.
func WriteGroupTypeDirectory(baseDir, groupID string, items []ResourceItem, opts *hcl.ResourceBlockOptions) error {
	return WriteModuleDirectoryWithFSAndGroup(DefaultFS, baseDir, items, opts, groupID)
}

// removeStaleFlatDirs removes flat type directories (e.g. baseDir/pipeline/) when using group-by layout,
// so a previous --flat run does not leave stale dirs at root.
func removeStaleFlatDirs(baseDir string, items []ResourceItem) {
	seen := make(map[string]bool)
	for _, it := range items {
		folder := typeNameToFolderName(it.TypeName)
		if !seen[folder] {
			seen[folder] = true
			flatPath := filepath.Join(baseDir, folder)
			if fi, err := os.Stat(flatPath); err == nil && fi.IsDir() {
				_ = os.RemoveAll(flatPath)
			}
		}
	}
}
