package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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

func TestWriteMigrationReportIncludesSecretsAndExclusions(t *testing.T) {
	items := []generator.ResourceItem{{
		TypeName: "criblio_secret",
		Attrs: map[string]hcl.Value{
			"value": {Kind: hcl.KindVariableRef, VarName: "secret_value"},
		},
	}}
	report := buildMigrationReport(items, []migrationExclusion{{TypeName: "criblio_key", Reason: "sensitive"}}, nil, map[string]string{"a": "b"}, nil)
	outputDir := t.TempDir()
	require.NoError(t, writeMigrationReport(outputDir, report))

	content, err := os.ReadFile(filepath.Join(outputDir, migrationReportFilename))
	require.NoError(t, err)
	var decoded migrationReport
	require.NoError(t, json.Unmarshal(content, &decoded))
	assert.Equal(t, 1, decoded.ExportedResources)
	assert.Equal(t, []string{"secret_value"}, decoded.UnresolvedSecrets)
	assert.Equal(t, "criblio_key", decoded.Excluded[0].TypeName)
}
