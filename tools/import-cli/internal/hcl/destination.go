// OneOf conversion from API items (map of JSON values) to type-specific HCL blocks.
// Handles discriminators (e.g. destination type, collector type) and camelCase → snake_case.
package hcl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// camelToSnake converts camelCase or PascalCase to snake_case.
func camelToSnake(s string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// normalizeDiscriminator normalizes the API discriminator value for use in a block name.
// Handles "open_telemetry", "cribl_http", "WebhookTarget" -> "webhook_target", etc.
func normalizeDiscriminator(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	return camelToSnake(s)
}

// ResolveOneOfBlockName maps a normalized API discriminator to a provider block suffix using supportedBlockNames. blockNamePrefix is used to match supported names like "output_prometheus" to normalized "prometheus". Returns suffix and true when a match is found (suffix is the part after prefix, or full name when prefix is empty).
func ResolveOneOfBlockName(normalizedDisc string, supportedBlockNames []string, blockNamePrefix string) (suffix string, ok bool) {
	if len(supportedBlockNames) == 0 || normalizedDisc == "" {
		return "", false
	}
	if blockNamePrefix != "" {
		for _, s := range supportedBlockNames {
			if strings.TrimPrefix(s, blockNamePrefix) == normalizedDisc {
				return strings.TrimPrefix(s, blockNamePrefix), true
			}
		}
	}
	for _, s := range supportedBlockNames {
		if s == normalizedDisc {
			return s, true
		}
	}
	withTarget := normalizedDisc + "_target"
	for _, s := range supportedBlockNames {
		if s == withTarget {
			return s, true
		}
	}
	for _, s := range supportedBlockNames {
		if strings.TrimSuffix(s, "_target") == normalizedDisc {
			return s, true
		}
	}
	return "", false
}

// ResolveOneOfBlockNameRaw parses rawDisc, normalizes it, and resolves to a supported block suffix.
func ResolveOneOfBlockNameRaw(rawDisc string, supportedBlockNames []string, blockNamePrefix string) (suffix string, ok bool) {
	var s string
	if err := json.Unmarshal([]byte(rawDisc), &s); err != nil {
		s = strings.Trim(rawDisc, `"`)
	}
	return ResolveOneOfBlockName(normalizeDiscriminator(s), supportedBlockNames, blockNamePrefix)
}

// ItemMapToBlock is the generic oneOf handler: given an API item (keys = camelCase, values = JSON strings),
// returns the block name (prefix + normalized discriminator + suffix) and the block Value. Use for destination, collector,
// pack_destination, and any resource whose schema uses type-specific blocks keyed by a discriminator.
// If discriminatorAlias is non-nil, API discriminator values are mapped to provider block suffix (e.g. "collection" -> "rest").
// blockNameSuffix is appended when prefix is empty (e.g. "_target" for notification_target -> smtp_target, slack_target).
func ItemMapToBlock(item map[string]string, discriminatorField, blockNamePrefix, blockNameSuffix string, keysToSkip []string, discriminatorAlias map[string]string) (blockName string, value Value, err error) {
	if len(item) == 0 {
		return "", Value{Kind: KindNull}, nil
	}
	raw, ok := item[discriminatorField]
	if !ok || raw == "" {
		return "", Value{}, fmt.Errorf("item missing discriminator field %q", discriminatorField)
	}
	var discStr string
	if err := json.Unmarshal([]byte(raw), &discStr); err != nil {
		discStr = strings.Trim(raw, `"`)
	}
	if discriminatorAlias != nil {
		if alias := discriminatorAlias[discStr]; alias != "" {
			discStr = alias
		}
	}
	normalized := normalizeDiscriminator(discStr)
	blockName = blockNamePrefix + normalized
	if blockNameSuffix != "" && !strings.HasSuffix(blockName, blockNameSuffix) {
		blockName += blockNameSuffix
	}

	skipSet := make(map[string]bool, len(keysToSkip))
	for _, k := range keysToSkip {
		skipSet[k] = true
	}
	m, err := jsonMapToValueMap(item, skipSet)
	if err != nil {
		return "", Value{}, err
	}
	return blockName, Value{Kind: KindMap, Map: m}, nil
}

// DestinationItemToOutputBlock converts a single destination item into output_<type> block (convenience wrapper).
func DestinationItemToOutputBlock(item map[string]string) (blockName string, value Value, err error) {
	return ItemMapToBlock(item, "type", "output_", "", []string{"status"}, nil)
}

// ItemMapToFlatValues converts an API item (keys = camelCase, values = JSON strings) into a flat
// map of HCL values (snake_case keys). Use for resources like global_var where the API returns
// a single payload in a list and the schema has top-level attributes (description, type, value, etc.).
func ItemMapToFlatValues(item map[string]string, keysToSkip []string) (map[string]Value, error) {
	if len(item) == 0 {
		return nil, nil
	}
	skipSet := make(map[string]bool, len(keysToSkip))
	for _, k := range keysToSkip {
		skipSet[k] = true
	}
	return jsonMapToValueMap(item, skipSet)
}

// jsonMapToValueMap converts a map of JSON strings (API keys camelCase) to map[string]Value with snake_case keys.
// Empty lists are omitted so we never emit e.g. urls = [] which would fail schema validators (SizeAtLeast(1)).
func jsonMapToValueMap(item map[string]string, keysToSkip map[string]bool) (map[string]Value, error) {
	out := make(map[string]Value)
	keys := make([]string, 0, len(item))
	for k := range item {
		if keysToSkip[k] {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		jsonStr := item[k]
		sk := camelToSnake(k)
		v, err := jsonStringToValue(jsonStr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", k, err)
		}
		v = omitEmptyListsFromValue(v)
		// Don't emit empty lists; provider often has listvalidator.SizeAtLeast(1).
		if v.Kind == KindList && len(v.List) == 0 {
			continue
		}
		out[sk] = v
	}
	return out, nil
}

// omitEmptyListsFromValue returns a copy of v with empty list values omitted from nested maps.
// Used so nested attributes like urls = [] are not written (avoids SizeAtLeast(1) errors).
func omitEmptyListsFromValue(v Value) Value {
	switch v.Kind {
	case KindMap:
		m := make(map[string]Value, len(v.Map))
		for k, ev := range v.Map {
			cleaned := omitEmptyListsFromValue(ev)
			if cleaned.Kind == KindList && len(cleaned.List) == 0 {
				continue
			}
			m[k] = cleaned
		}
		return Value{Kind: KindMap, Map: m}
	case KindList:
		list := make([]Value, 0, len(v.List))
		for _, el := range v.List {
			list = append(list, omitEmptyListsFromValue(el))
		}
		return Value{Kind: KindList, List: list}
	default:
		return v
	}
}

// jsonStringToValue parses a JSON value string and converts it to Value.
func jsonStringToValue(jsonStr string) (Value, error) {
	jsonStr = strings.TrimSpace(jsonStr)
	if jsonStr == "" || jsonStr == "null" {
		return Value{Kind: KindNull}, nil
	}
	var raw interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		// treat as literal string (e.g. unquoted number or identifier)
		return Value{Kind: KindString, String: jsonStr}, nil
	}
	return anyToValue(raw)
}

func anyToValue(raw interface{}) (Value, error) {
	if raw == nil {
		return Value{Kind: KindNull}, nil
	}
	switch v := raw.(type) {
	case string:
		return Value{Kind: KindString, String: v}, nil
	case float64:
		return Value{Kind: KindNumber, Number: v}, nil
	case bool:
		return Value{Kind: KindBool, Bool: v}, nil
	case []interface{}:
		list := make([]Value, 0, len(v))
		for i, el := range v {
			ev, err := anyToValue(el)
			if err != nil {
				return Value{}, fmt.Errorf("list[%d]: %w", i, err)
			}
			list = append(list, ev)
		}
		return Value{Kind: KindList, List: list}, nil
	case map[string]interface{}:
		m := make(map[string]Value, len(v))
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sk := camelToSnake(k)
			ev, err := anyToValue(v[k])
			if err != nil {
				return Value{}, fmt.Errorf("%s: %w", k, err)
			}
			m[sk] = ev
		}
		return Value{Kind: KindMap, Map: m}, nil
	default:
		return Value{Kind: KindString, String: fmt.Sprint(raw)}, nil
	}
}
