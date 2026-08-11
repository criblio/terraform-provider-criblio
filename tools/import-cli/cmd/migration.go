package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/export"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

const migrationReportFilename = "migration-report.json"

type migrationExclusion struct {
	TypeName string `json:"type"`
	Reason   string `json:"reason"`
}

type migrationReport struct {
	Mode              string               `json:"mode"`
	ExportedResources int                  `json:"exported_resources"`
	ExportedByType    map[string]int       `json:"exported_by_type"`
	Excluded          []migrationExclusion `json:"excluded"`
	Transforms        []string             `json:"transforms,omitempty"`
	GroupMappings     map[string]string    `json:"group_mappings,omitempty"`
	UnresolvedSecrets []string             `json:"unresolved_secrets,omitempty"`
	DefaultsExcluded  int                  `json:"defaults_excluded,omitempty"`
	Warnings          []string             `json:"warnings,omitempty"`
}

func migrationPolicy(reg *registry.Registry) []migrationExclusion {
	reasons := map[string]string{
		"criblio_commit":                       "action-only resource",
		"criblio_deploy":                       "action-only resource",
		"criblio_group_system_settings":        "Cloud-managed settings are not portable",
		"criblio_group":                        "target Cloud groups must be provisioned separately and selected with --group-map",
		"criblio_key":                          "sensitive key material cannot be exported safely",
		"criblio_lakehouse_dataset_connection": "Cloud-managed relationship without a portable read API",
		"criblio_workspace":                    "target Cloud workspace is selected through provider configuration",
	}
	for _, entry := range reg.Entries() {
		switch {
		case auth.IsRestrictedOnPremEndpoint(entry.RESTListPath), auth.IsRestrictedOnPremEndpoint(entry.RESTGetPath):
			reasons[entry.TypeName] = "endpoint is unavailable on the on-prem source"
		case auth.IsRestrictedCloudEndpoint(entry.RESTListPath), auth.IsRestrictedCloudEndpoint(entry.RESTGetPath):
			reasons[entry.TypeName] = "endpoint is on-prem-only and unavailable in Cloud"
		}
	}
	excluded := make([]migrationExclusion, 0, len(reasons))
	for typeName, reason := range reasons {
		excluded = append(excluded, migrationExclusion{TypeName: typeName, Reason: reason})
	}
	sort.Slice(excluded, func(i, j int) bool { return excluded[i].TypeName < excluded[j].TypeName })
	return excluded
}

func migrationRestrictedTypes(reg *registry.Registry) []string {
	excluded := migrationPolicy(reg)
	types := make([]string, len(excluded))
	for i, item := range excluded {
		types[i] = item.TypeName
	}
	return types
}

func parseGroupMappings(values []string) (map[string]string, error) {
	mappings := make(map[string]string, len(values))
	for _, value := range values {
		source, target, ok := strings.Cut(value, "=")
		source = strings.TrimSpace(source)
		target = strings.TrimSpace(target)
		if !ok || source == "" || target == "" {
			return nil, fmt.Errorf("invalid --group-map %q: expected SOURCE=TARGET", value)
		}
		if existing, ok := mappings[source]; ok && existing != target {
			return nil, fmt.Errorf("source group %q maps to both %q and %q", source, existing, target)
		}
		mappings[source] = target
	}
	return mappings, nil
}

func prepareOnPremToCloudItems(items []generator.ResourceItem, mappings map[string]string) ([]string, error) {
	transforms := make([]string, 0)
	for i := range items {
		item := &items[i]
		item.ImportID = ""
		item.Migration = true
		sourceGroup := item.GroupID
		if sourceGroup != "" && sourceGroup != "global" {
			targetGroup, ok := mappings[sourceGroup]
			if !ok {
				if isIdentitySafeGroup(sourceGroup) {
					targetGroup = sourceGroup
				} else {
					return nil, fmt.Errorf("group %q requires an explicit --group-map %s=TARGET", sourceGroup, sourceGroup)
				}
			}
			if targetGroup != sourceGroup {
				transforms = append(transforms, fmt.Sprintf("mapped group %s to %s", sourceGroup, targetGroup))
			}
			item.GroupID = targetGroup
			rewriteGroupReferences(item.Attrs, mappings, sourceGroup, targetGroup)
		}
		if item.TypeName == "criblio_group" {
			if _, ok := item.Attrs["on_prem"]; ok {
				delete(item.Attrs, "on_prem")
				transforms = append(transforms, fmt.Sprintf("removed on_prem from %s", item.Name))
			}
			delete(item.Attrs, "provisioned")
			delete(item.Attrs, "worker_remote_access")
		}
	}
	sort.Strings(transforms)
	return compactStrings(transforms), nil
}

func isIdentitySafeGroup(groupID string) bool {
	switch groupID {
	case "default", "default_fleet", "defaultHybrid":
		return true
	default:
		return false
	}
}

func rewriteGroupReferences(attrs map[string]hcl.Value, mappings map[string]string, sourceGroup, targetGroup string) {
	var rewrite func(hcl.Value) hcl.Value
	rewrite = func(value hcl.Value) hcl.Value {
		switch value.Kind {
		case hcl.KindMap:
			for key, child := range value.Map {
				if (key == "group_id" || key == "inherits") && child.Kind == hcl.KindString {
					if mapped, ok := mappings[child.String]; ok {
						child.String = mapped
					} else if child.String == sourceGroup {
						child.String = targetGroup
					}
				}
				value.Map[key] = rewrite(child)
			}
		case hcl.KindList:
			for i := range value.List {
				value.List[i] = rewrite(value.List[i])
			}
		}
		return value
	}
	for key, value := range attrs {
		if (key == "group_id" || key == "inherits") && value.Kind == hcl.KindString {
			if mapped, ok := mappings[value.String]; ok {
				value.String = mapped
			} else if value.String == sourceGroup {
				value.String = targetGroup
			}
		}
		attrs[key] = rewrite(value)
	}
	if id, ok := attrs["id"]; ok && id.Kind == hcl.KindString && id.String == sourceGroup {
		id.String = targetGroup
		attrs["id"] = id
	}
}

func buildMigrationReport(items []generator.ResourceItem, excluded []migrationExclusion, transforms []string, mappings map[string]string, result *export.ExportResult) migrationReport {
	report := migrationReport{
		Mode:              "onprem-to-cloud",
		ExportedResources: len(items),
		ExportedByType:    make(map[string]int),
		Excluded:          excluded,
		Transforms:        transforms,
		GroupMappings:     mappings,
	}
	secretSet := make(map[string]bool)
	for _, item := range items {
		report.ExportedByType[item.TypeName]++
		for _, name := range hcl.CollectSecretVariableNames(item.Attrs) {
			secretSet[name] = true
		}
	}
	for name := range secretSet {
		report.UnresolvedSecrets = append(report.UnresolvedSecrets, name)
	}
	sort.Strings(report.UnresolvedSecrets)
	if result != nil {
		report.DefaultsExcluded = result.DefaultsSkipped
		for _, skipped := range result.ListSkipped {
			if skipped.Count > 0 || skipped.Reason != "list returned 0 identifiers" {
				report.Warnings = append(report.Warnings, fmt.Sprintf("%s: %s", skipped.TypeName, skipped.Reason))
			}
		}
		report.Warnings = append(report.Warnings, result.ConvertSkipped...)
	}
	sort.Strings(report.Warnings)
	return report
}

func writeMigrationReport(outputDir string, report migrationReport) error {
	content, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode migration report: %w", err)
	}
	content = append(content, '\n')
	if err := generator.DefaultFS.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create migration output directory: %w", err)
	}
	if err := generator.DefaultFS.WriteFileAtomic(outputDir, migrationReportFilename, content, 0644); err != nil {
		return fmt.Errorf("write migration report: %w", err)
	}
	return nil
}

func compactStrings(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if len(out) == 0 || out[len(out)-1] != value {
			out = append(out, value)
		}
	}
	return out
}
