package main

import (
	"path/filepath"
	"testing"
)

func TestWorkspaceSettingsPreservesBothScopes(t *testing.T) {
	file := filepath.Join(t.TempDir(), "spec.yml")
	writeFile(t, file, `openapi: 3.1.0
paths:
  /system/settings/conf:
    get:
      operationId: getSystemSettingsConf
    patch:
      operationId: updateSystemSettingsConf
  /system/limits:
    get:
      operationId: getLimits
`)
	spec, err := readYAML(file)
	if err != nil {
		t.Fatal(err)
	}
	// Prefixing twice must neither remove the workspace route nor duplicate it.
	for range 2 {
		if err := prefixGroupScopedPaths(spec); err != nil {
			t.Fatal(err)
		}
	}
	for _, method := range []string{"get", "patch"} {
		workspace, err := lookupTarget(documentMapping(spec), "$.paths./system/settings/conf."+method)
		if err != nil {
			t.Fatal(err)
		}
		if params, ok := mappingValue(workspace, "parameters"); ok && hasNamedParameter(params, "groupId") {
			t.Fatal("workspace settings gained a group selector")
		}
		group, err := lookupTarget(documentMapping(spec), "$.paths./m/{groupId}/system/settings/conf."+method)
		if err != nil {
			t.Fatal(err)
		}
		params, ok := mappingValue(group, "parameters")
		if !ok || !hasNamedParameter(params, "groupId") {
			t.Fatal("group settings lost their group selector")
		}
	}
	if _, err := lookupTarget(documentMapping(spec), "$.paths./system/limits"); err == nil {
		t.Fatal("unrelated group-only paths must retain their existing behavior")
	}
}
