package cmd

import (
	"context"
	"slices"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationRestrictedTypesUsesEndpointAvailability(t *testing.T) {
	reg, err := buildRegistry(context.Background())
	require.NoError(t, err)

	restricted := migrationRestrictedTypes(reg)
	assert.True(t, slices.Contains(restricted, "criblio_search_dataset"))
	assert.False(t, slices.Contains(restricted, "criblio_pipeline"))
}

func TestPrepareOnPremToCloudItems(t *testing.T) {
	items := []generator.ResourceItem{
		{
			TypeName: "criblio_group",
			ImportID: "default",
			Attrs: map[string]hcl.Value{
				"id":      {Kind: hcl.KindString, String: "default"},
				"on_prem": {Kind: hcl.KindBool, Bool: true},
			},
		},
		{TypeName: "criblio_pipeline", ImportID: "pipeline-1"},
	}

	_, err := prepareOnPremToCloudItems(items, nil)
	require.NoError(t, err)

	assert.Empty(t, items[0].ImportID)
	assert.NotContains(t, items[0].Attrs, "on_prem")
	assert.Empty(t, items[1].ImportID)
}

func TestPrepareOnPremToCloudItemsMapsGroups(t *testing.T) {
	items := []generator.ResourceItem{{
		TypeName: "criblio_pipeline",
		GroupID:  "worker-a",
		ImportID: "pipeline-1",
		Attrs: map[string]hcl.Value{
			"group_id": {Kind: hcl.KindString, String: "worker-a"},
		},
	}}

	transforms, err := prepareOnPremToCloudItems(items, map[string]string{"worker-a": "cloud-a"})
	require.NoError(t, err)
	assert.Equal(t, "cloud-a", items[0].GroupID)
	assert.Equal(t, "cloud-a", items[0].Attrs["group_id"].String)
	assert.Contains(t, transforms, "mapped group worker-a to cloud-a")
}

func TestPrepareOnPremToCloudItemsRequiresCustomGroupMapping(t *testing.T) {
	items := []generator.ResourceItem{{TypeName: "criblio_pipeline", GroupID: "worker-a"}}
	_, err := prepareOnPremToCloudItems(items, nil)
	require.ErrorContains(t, err, "requires an explicit --group-map")
}

func TestParseGroupMappings(t *testing.T) {
	mappings, err := parseGroupMappings([]string{"worker-a=cloud-a", "worker-b=cloud-b"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"worker-a": "cloud-a", "worker-b": "cloud-b"}, mappings)

	_, err = parseGroupMappings([]string{"invalid"})
	require.Error(t, err)
}

func TestParseGroupMappingsRejectsUnsafeGroupIDs(t *testing.T) {
	tests := []struct {
		name    string
		mapping string
	}{
		{name: "source traversal", mapping: "../worker=cloud-a"},
		{name: "target traversal", mapping: "worker-a=../../other-dir"},
		{name: "source separator", mapping: "workers/east=cloud-a"},
		{name: "target separator", mapping: `worker-a=cloud\east`},
		{name: "dot target", mapping: "worker-a=."},
		{name: "space", mapping: "worker a=cloud-a"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseGroupMappings([]string{test.mapping})
			require.ErrorContains(t, err, "invalid")
		})
	}
}
