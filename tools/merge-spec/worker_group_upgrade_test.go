package main

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestWorkerGroupUpgradeLeaderRoute(t *testing.T) {
	var spec yaml.Node
	err := yaml.Unmarshal([]byte(`paths:
  /system/distributed/upgrade/{group}:
    post:
      parameters:
        - name: group
          in: path
          required: true
  /system/certificates:
    get: {}
`), &spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := prefixGroupScopedPaths(&spec); err != nil {
		t.Fatal(err)
	}
	paths, _ := mappingValue(documentMapping(&spec), "paths")
	upgrade, ok := mappingValue(paths, "/system/distributed/upgrade/{group}")
	if !ok {
		t.Fatal("leader upgrade route must not be prefixed")
	}
	post, _ := mappingValue(upgrade, "post")
	params, _ := mappingValue(post, "parameters")
	if len(params.Content) != 1 || hasNamedParameter(params, "groupId") {
		t.Fatal("leader upgrade must not acquire a second group parameter")
	}
	if _, ok := mappingValue(paths, "/m/{groupId}/system/certificates"); !ok {
		t.Fatal("other system routes must retain Worker Group scoping")
	}
}
