// HCL encoding of Value to expression strings, preserving null and structure.
package hcl

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ToHCLExpr returns an HCL expression string for a single Value
// (e.g. "null", `"foo"`, "true", "[1,2]", "{ a = \"b\" }").
func (v Value) ToHCLExpr() string {
	switch v.Kind {
	case KindSensitive:
		return strconv.Quote(v.Sensitive)
	case KindNull:
		return "null"
	case KindList:
		parts := make([]string, 0, len(v.List))
		for _, el := range v.List {
			parts = append(parts, el.ToHCLExpr())
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case KindMap:
		keys := make([]string, 0, len(v.Map))
		for k := range v.Map {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			ev := v.Map[k]
			parts = append(parts, fmt.Sprintf("%s = %s", hclQuoteKey(k), ev.ToHCLExpr()))
		}
		return "{ " + strings.Join(parts, ", ") + " }"
	case KindBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case KindNumber:
		return strconv.FormatFloat(v.Number, 'f', -1, 64)
	case KindString:
		return strconv.Quote(v.String)
	case KindVariableRef:
		if v.VarName != "" && isVariableName(v.VarName) {
			return "var." + v.VarName
		}
		return "null"
	default:
		return "null"
	}
}

func hclQuoteKey(k string) string {
	if needsQuote(k) {
		return strconv.Quote(k)
	}
	return k
}

func needsQuote(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return true
		}
	}
	return false
}
