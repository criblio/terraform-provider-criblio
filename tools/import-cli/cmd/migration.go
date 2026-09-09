package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

var groupIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type migrationExclusion struct {
	TypeName string `json:"type"`
	Reason   string `json:"reason"`
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
		if !groupIDPattern.MatchString(source) {
			return nil, fmt.Errorf("invalid source group %q: group IDs may contain only letters, digits, underscores, and hyphens", source)
		}
		if !groupIDPattern.MatchString(target) {
			return nil, fmt.Errorf("invalid target group %q: group IDs may contain only letters, digits, underscores, and hyphens", target)
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

func compactStrings(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if len(out) == 0 || out[len(out)-1] != value {
			out = append(out, value)
		}
	}
	return out
}
