// Package export converts discovery results into generator ResourceItems.
package export

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

func groupIDFromIDMap(idMap map[string]string) string {
	if g := idMap["group_id"]; g != "" {
		return g
	}
	return "global"
}

// groupIDForOutput returns the output group/folder for module layout.
// - criblio_cribl_lake_house: always "global" (Lake API is not group-scoped)
// - default_search -> "search" (rename for clarity)
// - criblio_search_* types (search_dataset, search_dataset_provider, etc.) -> "search" so they live under search/
func groupIDForOutput(typeName string, gid string) string {
	if typeName == "criblio_cribl_lake_house" {
		return "global"
	}
	if gid == "default_search" {
		return "search"
	}
	// Any criblio_search_* type (search_dataset, search_dataset_provider, etc.) goes under search folder
	if strings.HasPrefix(typeName, "criblio_search_") {
		return "search"
	}
	return gid
}

// allowedOutputFoldersFromGroupIDs maps resolved worker/search group IDs to module layout folder names.
// default_search is also mapped to "search" so exports can include Search UI resources when that group is selected.
func allowedOutputFoldersFromGroupIDs(groupIDs []string) map[string]bool {
	m := make(map[string]bool, len(groupIDs)+1)
	for _, g := range groupIDs {
		m[g] = true
		if g == "default_search" {
			m["search"] = true
		}
	}
	return m
}

// skipExportForGroupFilter reports whether a resource should be omitted when the user passed --group.
// Only resources whose output folder (see groupIDForOutput) is in allowedOutputFoldersFromGroupIDs(groupIDs)
// are exported; "global" and "search" are excluded unless the user included the corresponding scope
// (search requires default_search in groupIDs).
func skipExportForGroupFilter(typeName string, idMap map[string]string, groupFilter []string, groupIDs []string) bool {
	if len(groupFilter) == 0 {
		return false
	}
	out := groupIDForOutput(typeName, groupIDFromIDMap(idMap))
	allowed := allowedOutputFoldersFromGroupIDs(groupIDs)
	return !allowed[out]
}

// toRequestParams maps lowercase identifier keys (group_id, id, pack) to
// request param names (GroupID, ID, Pack) expected by the SDK and converter.
func toRequestParams(idMap map[string]string) map[string]string {
	out := make(map[string]string)
	if v := idMap["group_id"]; v != "" {
		out["GroupID"] = v
		// GetGroupsByID and similar use ID in the path; use group_id when id is not set.
		if idMap["id"] == "" {
			out["ID"] = v
		}
	}
	if v := idMap["id"]; v != "" {
		out["ID"] = v
	}
	if v := idMap["pack"]; v != "" {
		out["Pack"] = v
	}
	if v := idMap["lake_id"]; v != "" {
		out["LakeID"] = v
	}
	return out
}

// rawJSONToItemMap converts raw API item JSON to map[string]string (keys preserved, values JSON-encoded)
// for use with hcl.ItemMapToBlock.
func rawJSONToItemMap(itemJSON []byte) map[string]string {
	var m map[string]interface{}
	if err := json.Unmarshal(itemJSON, &m); err != nil || len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			continue
		}
		out[k] = string(b)
	}
	return out
}

// firstItemMapFromModel returns the first item from the model's slice field (e.g. Items) as map[string]string.
// readOnlyAttrTfsdk is the tfsdk attribute name (e.g. "items"); the Go field name is derived for reflection.
func firstItemMapFromModel(model interface{}, readOnlyAttrTfsdk string) map[string]string {
	goFieldName := tfsdkNameToGoFieldName(readOnlyAttrTfsdk)
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	field := val.FieldByName(goFieldName)
	if !field.IsValid() || field.Kind() != reflect.Slice || field.Len() == 0 {
		return nil
	}
	first := field.Index(0)
	if first.Kind() != reflect.Map {
		return nil
	}
	out := make(map[string]string)
	for _, k := range first.MapKeys() {
		key := k.String()
		mv := first.MapIndex(k)
		if !mv.IsValid() {
			continue
		}
		if n, ok := mv.Interface().(jsontypes.Normalized); ok {
			if n.IsNull() || n.IsUnknown() {
				continue
			}
			out[key] = n.ValueString()
		}
	}
	return out
}

func tfsdkNameToGoFieldName(tfsdk string) string {
	if tfsdk == "" {
		return ""
	}
	return strings.ToUpper(tfsdk[:1]) + tfsdk[1:]
}
