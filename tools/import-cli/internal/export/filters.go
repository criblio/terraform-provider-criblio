// Package export converts discovery results into generator ResourceItems.
package export

import (
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/exclusions"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// sanitizeConvertError returns a short, safe message for user-facing output.
// Never includes raw JSON (which may contain credentials or sensitive data).
func sanitizeConvertError(err error) string {
	if err == nil {
		return ""
	}
	s := err.Error()
	// SDK unmarshal errors include raw JSON; never expose to users.
	if strings.Contains(s, "could not unmarshal") || strings.Contains(s, "error unmarshaling json") {
		return "unsupported type (SDK unmarshal failed)"
	}
	// Truncate long errors that might contain sensitive data
	if len(s) > 120 {
		s = s[:120] + "..."
	}
	return s
}

func skipResourceByID(typeName string, idMap map[string]string) bool {
	// Search worker group: we don't support riptide yet.
	if typeName == "criblio_source" && idMap["group_id"] == "default_search" {
		return true
	}
	// criblio_pack: pack id is in "id"
	if typeName == "criblio_pack" && custom.SkipPacks[idMap["id"]] {
		return true
	}
	// All pack-scoped resources: skip when pack is in SkipPacks
	if pack := idMap["pack"]; pack != "" && custom.SkipPacks[pack] {
		return true
	}
	// criblio_pack_lookups: skip built-in lookups (id starts with "cribl."); provider id pattern rejects them.
	if typeName == "criblio_pack_lookups" {
		if id := idMap["id"]; id != "" && strings.HasPrefix(id, "cribl.") {
			return true
		}
	}
	// criblio_pack_vars: skip vars whose id contains dots (e.g. cribl.my_globalvar); provider id pattern is ^[a-zA-Z0-9_-]+$
	if typeName == "criblio_pack_vars" {
		if id := idMap["id"]; id != "" && strings.Contains(id, ".") {
			return true
		}
	}
	// Skip when id equals group_id (wrong format): the API returns items with their own ids (e.g. "test.csv", "my_event_breaker_rule"),
	// not group names. Exception: criblio_group_system_settings and criblio_routes intentionally use group_id as id (one resource per group).
	if gid := idMap["group_id"]; gid != "" {
		if id := idMap["id"]; id != "" && id == gid {
			switch typeName {
			case "criblio_group_system_settings", "criblio_routes":
				// Intentional: one resource per group, id=group_id
			default:
				return true
			}
		}
	}
	ids, ok := exclusions.SkipExportIDs[typeName]
	if !ok || len(ids) == 0 {
		return false
	}
	// Match ImportIDFormat "id" and group-only maps (e.g. criblio_group uses group_id without "id").
	if v := idMap["id"]; v != "" && ids[v] {
		return true
	}
	if v := idMap["group_id"]; v != "" && ids[v] {
		return true
	}
	return false
}

func skipResourceWhenLibCribl(attrs map[string]hcl.Value) bool {
	v, ok := attrs["lib"]
	if !ok || v.Kind != hcl.KindString {
		return false
	}
	return v.String == "cribl"
}

// DefaultResource reports whether the resource is a built-in Cribl default.
func DefaultResource(typeName string, idMap map[string]string, attrs map[string]hcl.Value) bool {
	if skipResourceWhenLibCribl(attrs) {
		return true
	}
	if criblDefaultTag(attrs) {
		return true
	}
	id := idMap["id"]
	pack := idMap["pack"]
	switch typeName {
	case "criblio_destination":
		return custom.DefaultDestinationIDs[id]
	case "criblio_pack_destination":
		if custom.DefaultPackIDs[pack] {
			return true
		}
		return custom.DefaultDestinationIDs[id]
	case "criblio_search_dataset":
		return custom.DefaultSearchDatasetIDs[id]
	case "criblio_search_dataset_provider":
		return custom.DefaultSearchDatasetProviderIDs[id]
	case "criblio_cribl_lake_dataset":
		return custom.DefaultCriblLakeDatasetIDs[id]
	case "criblio_group":
		// Groups are filtered in post-processing (filterEmptyDefaultGroups) to keep
		// default groups that have user-created resources within them.
		return false
	case "criblio_pack":
		return custom.DefaultPackIDs[id]
	case "criblio_pipeline":
		if custom.DefaultPipelineIDs[id] {
			return true
		}
		if strings.HasPrefix(id, "pack:") {
			packName := strings.TrimPrefix(id, "pack:")
			return custom.DefaultPackIDs[packName]
		}
		return false
	case "criblio_grok":
		return custom.DefaultGrokIDs[id]
	case "criblio_parquet_schema":
		return custom.DefaultParquetSchemaIDs[id]
	case "criblio_schema":
		return custom.DefaultSchemaIDs[id]
	case "criblio_source":
		return custom.DefaultSourceIDs[id]
	case "criblio_event_breaker_ruleset":
		return custom.DefaultEventBreakerRulesetIDs[id]
	case "criblio_pack_breakers":
		if custom.DefaultPackIDs[pack] {
			return true
		}
		return custom.DefaultEventBreakerRulesetIDs[id]
	case "criblio_pack_pipeline":
		// Only check if the pack itself is a default; pipeline IDs like "main" are valid inside user packs.
		return custom.DefaultPackIDs[pack]
	case "criblio_pack_source":
		// Only check if the pack itself is a default; source IDs are valid inside user packs.
		return custom.DefaultPackIDs[pack]
	case "criblio_pack_vars":
		return custom.DefaultPackIDs[pack]
	case "criblio_routes", "criblio_pack_routes":
		// Routes are a singleton per group/pack; users modify but don't create them from scratch.
		return true
	case "criblio_search_dataset_ruleset":
		return custom.DefaultSearchDatasetRulesetIDs[id]
	case "criblio_search_datatype_ruleset":
		return custom.DefaultSearchDatatypeRulesetIDs[id]
	case "criblio_search_saved_query":
		return custom.DefaultSearchSavedQueryIDs[id]
	}
	return false
}

// criblDefaultTag reports whether attrs contains the cribl:default tag.
func criblDefaultTag(attrs map[string]hcl.Value) bool {
	v, ok := attrs["tags"]
	if !ok {
		return false
	}
	if v.Kind == hcl.KindString {
		return v.String == custom.CriblDefaultTag
	}
	if v.Kind == hcl.KindList {
		for _, item := range v.List {
			if item.Kind == hcl.KindString && item.String == custom.CriblDefaultTag {
				return true
			}
		}
	}
	return false
}
