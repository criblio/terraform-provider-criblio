package generator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// BuildImportID returns the Terraform import ID string for the given format and identifiers.
// Format examples: "id", "group_id", "json:group_id,id", "json:group_id,id,pack".
// For "json:key1,key2,..." the result is a JSON object with those keys (sorted).
// For a single key the result is the value; for multiple keys without "json:" the result is comma-separated (order by format).
func BuildImportID(format string, identifiers map[string]string) (string, error) {
	format = strings.TrimSpace(format)
	if format == "" {
		return "", fmt.Errorf("import ID format is empty")
	}
	if strings.HasPrefix(format, "json:") {
		keys := strings.Split(strings.TrimPrefix(format, "json:"), ",")
		for i := range keys {
			keys[i] = strings.TrimSpace(keys[i])
		}
		obj := make(map[string]string)
		for _, k := range keys {
			if k == "" {
				continue
			}
			if v, ok := identifiers[k]; ok {
				obj[k] = v
			}
		}
		// Sort keys for deterministic JSON
		names := make([]string, 0, len(obj))
		for k := range obj {
			names = append(names, k)
		}
		sort.Strings(names)
		ordered := make(map[string]string, len(names))
		for _, k := range names {
			ordered[k] = obj[k]
		}
		b, err := json.Marshal(ordered)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	// Single key or comma-separated
	keys := strings.Split(format, ",")
	for i := range keys {
		keys[i] = strings.TrimSpace(keys[i])
	}
	if len(keys) == 1 && keys[0] != "" {
		if v, ok := identifiers[keys[0]]; ok {
			return v, nil
		}
		return "", fmt.Errorf("identifier %q not found", keys[0])
	}
	var parts []string
	for _, k := range keys {
		if k == "" {
			continue
		}
		if v, ok := identifiers[k]; ok {
			parts = append(parts, v)
		}
	}
	return strings.Join(parts, ","), nil
}
