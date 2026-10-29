package hcl

import (
	"fmt"
	"strings"
)

// ReplaceSensitiveValuesWithVariableRefs replaces schema-marked credential
// strings with plain variable references. Collection indexes identify variables,
// but do not participate in matching schema paths.
func ReplaceSensitiveValuesWithVariableRefs(attrs map[string]Value, resourceName string, sensitivePaths map[string]bool) {
	var replace func(*Value, string, string, bool)
	replace = func(value *Value, path, schemaPath string, sensitive bool) {
		sensitive = sensitive || sensitiveSchemaPath(schemaPath, sensitivePaths)
		switch value.Kind {
		case KindString:
			if sensitive && value.String != "" {
				*value = Value{Kind: KindVariableRef, VarName: sanitizeVarName(resourceName, path)}
			}
		case KindList, KindTuple:
			for index := range value.List {
				replace(&value.List[index], fmt.Sprintf("%s[%d]", path, index), schemaPath, sensitive)
			}
		case KindMap:
			for name, child := range value.Map {
				replace(&child, path+"."+name, schemaPath+"."+name, sensitive)
				value.Map[name] = child
			}
		}
	}
	for name, value := range attrs {
		replace(&value, name, name, false)
		attrs[name] = value
	}
}

func sensitiveSchemaPath(path string, sensitivePaths map[string]bool) bool {
	if sensitivePaths[path] {
		return true
	}
	parts := strings.Split(path, ".")
	for pattern, sensitive := range sensitivePaths {
		if !sensitive || !strings.Contains(pattern, "*") {
			continue
		}
		patternParts := strings.Split(pattern, ".")
		if len(parts) != len(patternParts) {
			continue
		}
		matches := true
		for index, part := range parts {
			if patternParts[index] != "*" && patternParts[index] != part {
				matches = false
				break
			}
		}
		if matches {
			return true
		}
	}
	return false
}
