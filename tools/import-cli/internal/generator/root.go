package generator

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// RootModuleInfo describes a child module for the root main.tf.
type RootModuleInfo struct {
	Name      string   // Terraform module name (e.g. "pipeline", "default_pipeline")
	Path      string   // Source path relative to root (e.g. "./pipeline", "./default/pipeline")
	Variables []string // Secret variable names to pass from root (e.g. "secret_default_test_secret_value")
}

// WriteRootFiles writes the top-level main.tf, providers.tf, variables.tf, and import.tf so users can run
// terraform init, plan, and apply from the output directory.
// Import blocks are at root only (Terraform requires import in root module, not in child modules).
func WriteRootFiles(baseDir string, moduleInfos []RootModuleInfo, items []ResourceItem, byGroup bool) error {
	if err := DefaultFS.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	providersContent := rootProvidersTF()
	if err := DefaultFS.WriteFileAtomic(baseDir, "providers.tf", providersContent, 0644); err != nil {
		return fmt.Errorf("write providers.tf: %w", err)
	}
	mainContent := rootMainTF(moduleInfos)
	if err := DefaultFS.WriteFileAtomic(baseDir, "main.tf", mainContent, 0644); err != nil {
		return fmt.Errorf("write main.tf: %w", err)
	}
	variablesContent := rootVariablesTF(moduleInfos)
	if len(variablesContent) > 0 {
		if err := DefaultFS.WriteFileAtomic(baseDir, "variables.tf", variablesContent, 0644); err != nil {
			return fmt.Errorf("write variables.tf: %w", err)
		}
	}
	importContent := rootImportTF(items, byGroup)
	if len(importContent) > 0 {
		if err := DefaultFS.WriteFileAtomic(baseDir, "import.tf", importContent, 0644); err != nil {
			return fmt.Errorf("write import.tf: %w", err)
		}
	}
	return nil
}

func rootProvidersTF() []byte {
	return []byte(`# Terraform and provider configuration for imported Cribl resources.
# Configure authentication via environment variables or a .cribl/credentials file.
# See: https://registry.terraform.io/providers/criblio/criblio/latest/docs

terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = ">= 1.0"
    }
  }
}

provider "criblio" {
  # Uses CRIBL_* env vars or ~/.cribl/credentials (see README)
  # Cloud: CRIBL_CLIENT_ID, CRIBL_CLIENT_SECRET, CRIBL_ORGANIZATION_ID, CRIBL_WORKSPACE_ID
  # On-prem: CRIBL_ONPREM_SERVER_URL, CRIBL_BEARER_TOKEN or CRIBL_ONPREM_USERNAME/CRIBL_ONPREM_PASSWORD
}
`)
}

func rootImportTF(items []ResourceItem, byGroup bool) []byte {
	copied := make([]ResourceItem, 0, len(items))
	for _, it := range items {
		if it.ImportID == "" {
			continue
		}
		copied = append(copied, it)
	}
	if len(copied) == 0 {
		return nil
	}
	SortResourceItems(copied)
	// Deduplicate by address: Terraform allows only one import block per resource.
	seenAddr := make(map[string]bool)
	var buf bytes.Buffer
	buf.WriteString("# Import blocks - must be in root module (Terraform restriction).\n")
	buf.WriteString("# Run terraform plan to execute imports.\n\n")
	for _, it := range copied {
		modName := moduleNameForItem(it, byGroup)
		addr := fmt.Sprintf("module.%s.%s.%s", modName, it.TypeName, it.Name)
		if seenAddr[addr] {
			continue
		}
		seenAddr[addr] = true
		fmt.Fprintf(&buf, "import {\n  to = %s\n  id = %q\n}\n\n", addr, it.ImportID)
	}
	return buf.Bytes()
}

func moduleNameForItem(it ResourceItem, byGroup bool) string {
	if byGroup {
		gid := it.GroupID
		if gid == "" {
			gid = "global"
		}
		return sanitizeModuleName(gid) + "_" + typeNameToFolderName(it.TypeName)
	}
	return typeNameToFolderName(it.TypeName)
}

func rootMainTF(moduleInfos []RootModuleInfo) []byte {
	// Sort for deterministic output
	sorted := make([]RootModuleInfo, len(moduleInfos))
	copy(sorted, moduleInfos)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })

	var buf bytes.Buffer
	buf.WriteString("# Root module - run terraform init, plan, apply from this directory.\n")
	buf.WriteString("# Each module below contains imported Cribl resources.\n\n")
	for _, m := range sorted {
		buf.WriteString(fmt.Sprintf("module %q {\n  source = %q\n", m.Name, m.Path))
		for _, v := range m.Variables {
			buf.WriteString(fmt.Sprintf("  %s = var.%s\n", v, v))
		}
		buf.WriteString("}\n\n")
	}
	return buf.Bytes()
}

func rootVariablesTF(moduleInfos []RootModuleInfo) []byte {
	seen := make(map[string]bool)
	var names []string
	for _, m := range moduleInfos {
		for _, v := range m.Variables {
			if !seen[v] {
				seen[v] = true
				names = append(names, v)
			}
		}
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	var buf bytes.Buffer
	buf.WriteString("# Root-level variables for secrets passed to child modules.\n")
	buf.WriteString("# Set via TF_VAR_*, -var, -var-file, or .tfvars (e.g. export TF_VAR_secret_xxx=\"...\")\n\n")
	for _, n := range names {
		fmt.Fprintf(&buf, "variable %q {\n  type      = string\n  sensitive = true\n}\n\n", n)
	}
	return buf.Bytes()
}

// RootModuleInfosFromItems returns RootModuleInfo for the default (by-type) layout.
func RootModuleInfosFromItems(items []ResourceItem) []RootModuleInfo {
	byType := make(map[string][]ResourceItem)
	for _, it := range items {
		byType[it.TypeName] = append(byType[it.TypeName], it)
	}
	var infos []RootModuleInfo
	for typeName, moduleItems := range byType {
		folder := typeNameToFolderName(typeName)
		vars := collectSecretVariableNames(moduleItems)
		infos = append(infos, RootModuleInfo{
			Name:      folder,
			Path:      "./" + folder,
			Variables: vars,
		})
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].Name < infos[j].Name })
	return infos
}

// RootModuleInfosFromItemsByGroup returns RootModuleInfo for the group-by layout.
func RootModuleInfosFromItemsByGroup(items []ResourceItem) []RootModuleInfo {
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
	var infos []RootModuleInfo
	for gid, byType := range byGroupAndType {
		for typeName, moduleItems := range byType {
			folder := typeNameToFolderName(typeName)
			moduleName := sanitizeModuleName(gid) + "_" + folder
			vars := collectSecretVariableNames(moduleItems)
			infos = append(infos, RootModuleInfo{
				Name:      moduleName,
				Path:      "./" + gid + "/" + folder,
				Variables: vars,
			})
		}
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].Name < infos[j].Name })
	return infos
}

func collectSecretVariableNames(items []ResourceItem) []string {
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
	return names
}

func sanitizeModuleName(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
