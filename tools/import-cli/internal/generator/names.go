// Package generator produces deterministic Terraform HCL and module directory structure.
package generator

import (
	"regexp"
	"sort"
	"strings"
)

// MaxResourceNameLength is the maximum length for a generated resource name (Terraform limit is 255).
const MaxResourceNameLength = 64

// terraformSafeNameRe matches characters that are not allowed in Terraform resource names.
// Allowed: letters, digits, underscore, hyphen. We normalize to underscore for stability.
var terraformSafeNameRe = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

// StableResourceName returns a Terraform-safe, stable name from the given type and identifier parts.
// Same inputs always produce the same output. Parts are typically id, group_id, pack, etc.
// Invalid characters are replaced with underscore; multiple underscores are collapsed; result is truncated.
func StableResourceName(typeName string, parts []string) string {
	var b strings.Builder
	// Prefix with shortened type (criblio_source -> source) for readability and uniqueness across types
	typeSuffix := typeName
	if idx := strings.Index(typeName, "_"); idx >= 0 {
		typeSuffix = typeName[idx+1:]
	}
	typeSuffix = terraformSafeNameRe.ReplaceAllString(typeSuffix, "_")
	typeSuffix = strings.Trim(strings.TrimSpace(typeSuffix), "_")
	if typeSuffix != "" {
		b.WriteString(typeSuffix)
	}
	for _, p := range parts {
		s := terraformSafeNameRe.ReplaceAllString(p, "_")
		s = strings.Trim(strings.TrimSpace(s), "_")
		if s == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("_")
		}
		b.WriteString(s)
	}
	out := b.String()
	if out == "" {
		out = "resource"
	}
	// Collapse multiple underscores
	for strings.Contains(out, "__") {
		out = strings.ReplaceAll(out, "__", "_")
	}
	out = strings.Trim(out, "_")
	if len(out) > MaxResourceNameLength {
		out = out[:MaxResourceNameLength]
		out = strings.TrimRight(out, "_")
	}
	return out
}

// StableResourceNameFromMap builds a stable name from type and an ordered set of identifier values.
// Keys are sorted so the same logical identifiers produce the same name.
func StableResourceNameFromMap(typeName string, identifiers map[string]string) string {
	keys := make([]string, 0, len(identifiers))
	for k := range identifiers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		if v := identifiers[k]; v != "" {
			parts = append(parts, v)
		}
	}
	return StableResourceName(typeName, parts)
}
