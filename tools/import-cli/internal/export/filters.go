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
	// Use "id" key for types with ImportIDFormat "id"; group-scoped types use idMap["id"] as well.
	if v := idMap["id"]; v != "" && ids[v] {
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
