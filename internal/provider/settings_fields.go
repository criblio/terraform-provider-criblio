package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type settingsPrivate interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
	SetKey(context.Context, string, []byte) diag.Diagnostics
}

// settingsFieldMask records configured leaves and nested object fields.
// Persisting the mask retains ownership when the API removes a managed value.
type settingsFieldMask map[string]settingsFieldMask

func settingsManagedFields(ctx context.Context, mask settingsFieldMask) (types.List, diag.Diagnostics) {
	fields := make([]string, 0)
	var collect func(settingsFieldMask, string)
	collect = func(mask settingsFieldMask, prefix string) {
		for name, nested := range mask {
			field := prefix + name
			if len(nested) == 0 {
				fields = append(fields, field)
				continue
			}
			collect(nested, field+".")
		}
	}
	collect(mask, "")
	sort.Strings(fields)
	return types.ListValueFrom(ctx, types.StringType, fields)
}

func settingsMask(payload map[string]json.RawMessage) settingsFieldMask {
	mask := make(settingsFieldMask, len(payload))
	for name, raw := range payload {
		var nested map[string]json.RawMessage
		if json.Unmarshal(raw, &nested) == nil && nested != nil {
			mask[name] = settingsMask(nested)
		} else {
			mask[name] = nil
		}
	}
	return mask
}

func mergeSettingsFields(current, configured map[string]json.RawMessage) (map[string]json.RawMessage, error) {
	merged := make(map[string]json.RawMessage, len(current)+len(configured))
	for name, value := range current {
		merged[name] = value
	}
	for name, value := range configured {
		var requested, existing map[string]json.RawMessage
		if json.Unmarshal(value, &requested) == nil && requested != nil {
			if raw, ok := current[name]; ok && string(raw) != "null" {
				if err := json.Unmarshal(raw, &existing); err != nil {
					return nil, fmt.Errorf("existing %s is not an object: %w", name, err)
				}
			}
			nested, err := mergeSettingsFields(existing, requested)
			if err != nil {
				return nil, err
			}
			value, err = json.Marshal(nested)
			if err != nil {
				return nil, err
			}
		}
		merged[name] = value
	}
	return merged, nil
}

func selectSettingsFields(current map[string]json.RawMessage, mask settingsFieldMask) (map[string]json.RawMessage, error) {
	selected := make(map[string]json.RawMessage, len(mask))
	for name, nestedMask := range mask {
		value, ok := current[name]
		if !ok {
			continue
		}
		if nestedMask != nil && string(value) != "null" {
			var nested map[string]json.RawMessage
			if err := json.Unmarshal(value, &nested); err != nil {
				return nil, fmt.Errorf("managed %s is not an object: %w", name, err)
			}
			filtered, err := selectSettingsFields(nested, nestedMask)
			if err != nil {
				return nil, err
			}
			value, err = json.Marshal(filtered)
			if err != nil {
				return nil, err
			}
		}
		selected[name] = value
	}
	return selected, nil
}
