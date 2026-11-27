package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// matchOrderedConfig associates configured rows with prior rows without treating
// a list index as identity. Within a repeated key, reserve unambiguous matches
// before considering rows that omit one or more distinguishing settings.
func matchOrderedConfig(config, prior []attr.Value, key string, compatible func(types.Object, types.Object) bool) []int {
	matches := make([]int, len(config))
	newGroups := make(map[string][]int)
	oldGroups := make(map[string][]int)
	for i, row := range config {
		matches[i] = -1
		if name := objectString(row, key); name != "" {
			newGroups[name] = append(newGroups[name], i)
		}
	}
	for i, row := range prior {
		if name := objectString(row, key); name != "" {
			oldGroups[name] = append(oldGroups[name], i)
		}
	}
	used := make([]bool, len(prior))
	for name, group := range newGroups {
		old := oldGroups[name]
		if len(old) == 0 {
			continue
		}
		if len(group) == 1 && len(old) == 1 {
			matches[group[0]] = old[0]
			continue
		}
		for {
			// A candidate with two sole claimants is still ambiguous.
			claims := make(map[int][]int)
			for _, i := range group {
				if matches[i] >= 0 {
					continue
				}
				candidate := -1
				for _, j := range old {
					if used[j] || !compatible(config[i].(types.Object), prior[j].(types.Object)) {
						continue
					}
					if candidate >= 0 {
						candidate = -1
						break
					}
					candidate = j
				}
				if candidate >= 0 {
					claims[candidate] = append(claims[candidate], i)
				}
			}
			progress := false
			for j, rows := range claims {
				if len(rows) == 1 {
					matches[rows[0]] = j
					used[j] = true
					progress = true
				}
			}
			if !progress {
				break
			}
		}
	}
	return matches
}

func objectString(row attr.Value, key string) string {
	if row.IsNull() || row.IsUnknown() {
		return ""
	}
	return row.(types.Object).Attributes()[key].(types.String).ValueString()
}
