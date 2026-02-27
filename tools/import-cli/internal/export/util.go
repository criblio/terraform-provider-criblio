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
// - criblio_search_* under default_search -> "search"
func groupIDForOutput(typeName string, gid string) string {
	if typeName == "criblio_cribl_lake_house" {
		return "global"
	}
	if gid == "default_search" {
		return "search"
	}
	return gid
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
