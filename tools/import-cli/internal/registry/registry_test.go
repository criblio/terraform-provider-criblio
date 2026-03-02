package registry

import (
	"context"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromResources_discoversFromProvider(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, reg)
	assert.Greater(t, reg.Len(), 0, "registry should contain at least one resource type")
	// Registry must contain all resources returned by the provider (no hardcoded list).
	assert.Equal(t, len(constructors), reg.Len(), "registry length must match provider Resources() count")
}

func TestNewFromResources_containsMVPResourceTypes(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	// MVP Terraform resource types that must be present (per provider Resources())
	mvpTypes := []string{
		"criblio_source",
		"criblio_destination",
		"criblio_routes",
		"criblio_pipeline",
		"criblio_global_var",
		"criblio_lookup_file",
		"criblio_key",
		"criblio_certificate",
		"criblio_notification",
		"criblio_notification_target",
		"criblio_pack",
		"criblio_commit",
		"criblio_deploy",
	}

	for _, typeName := range mvpTypes {
		_, ok := reg.ByTypeName(typeName)
		assert.True(t, ok, "registry should contain MVP resource type %q", typeName)
	}
}

func TestNewFromResources_terraformTypeNamesMatchProvider(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	// All type names must be prefixed with provider type
	for _, e := range reg.Entries() {
		assert.Contains(t, e.TypeName, "criblio_", "TypeName should match provider definition (criblio_*)")
		assert.NotEmpty(t, e.ModelTypeName, "each entry should include model type information")
	}
}

func TestNewFromResources_eachEntryHasModelType(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	for _, e := range reg.Entries() {
		assert.NotEmpty(t, e.TypeName, "TypeName must be set")
		assert.NotEmpty(t, e.ModelTypeName, "ModelTypeName must be set for %q", e.TypeName)
		assert.Contains(t, e.ModelTypeName, "Model", "model type should be *ResourceModel convention")
	}
}

func TestRegistry_ByTypeName(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok)
	assert.Equal(t, "criblio_source", e.TypeName)
	assert.Equal(t, "SourceResourceModel", e.ModelTypeName)

	_, ok = reg.ByTypeName("nonexistent_type")
	assert.False(t, ok)
}

func TestRegistry_TypeNames_and_Entries_consistent(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	names := reg.TypeNames()
	entries := reg.Entries()
	require.Len(t, names, len(entries))

	for i := range entries {
		assert.Equal(t, entries[i].TypeName, names[i])
	}
}

// TestNewFromResources_derivedMetadata validates that registry entries include
// List* and Get*ByID SDK method names and import ID format from ImportState logic.
func TestNewFromResources_derivedMetadata(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	// Source: Inputs.ListInput, GetInputByID, JSON group_id+id
	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok)
	assert.Equal(t, "Inputs", e.SDKService, "criblio_source should have SDKService for discovery")
	assert.Equal(t, "ListInput", e.ListMethod, "criblio_source should have ListMethod from SDK")
	assert.Equal(t, "GetInputByID", e.GetMethod, "criblio_source should have GetMethod from SDK")
	assert.Equal(t, "json:group_id,id", e.ImportIDFormat, "criblio_source import ID matches ImportState")

	// Pipeline: Pipelines.ListPipeline, GetPipelineByID, JSON group_id+id
	e, ok = reg.ByTypeName("criblio_pipeline")
	require.True(t, ok)
	assert.Equal(t, "Pipelines", e.SDKService)
	assert.Equal(t, "ListPipeline", e.ListMethod)
	assert.Equal(t, "GetPipelineByID", e.GetMethod)
	assert.Equal(t, "json:group_id,id", e.ImportIDFormat)

	// Notification: id-only import
	e, ok = reg.ByTypeName("criblio_notification")
	require.True(t, ok)
	assert.Equal(t, "ListNotification", e.ListMethod)
	assert.Equal(t, "GetNotificationByID", e.GetMethod)
	assert.Equal(t, "id", e.ImportIDFormat)

	// Lookup file: group_id-only import
	e, ok = reg.ByTypeName("criblio_lookup_file")
	require.True(t, ok)
	assert.Equal(t, "ListLookupFile", e.ListMethod)
	assert.Equal(t, "GetLookupFileByID", e.GetMethod)
	assert.Equal(t, "json:group_id,id", e.ImportIDFormat)
}

// TestImportMetadata_inSyncWithProvider ensures every resource type in the
// provider has an entry in ImportMetadata(). When you add a new resource to
// the provider, add it to importMetadataBase (and overrides/clearList if needed)
// in import_metadata.go; this test will fail until you do.
func TestImportMetadata_inSyncWithProvider(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	metadata := ImportMetadata()
	var missing []string
	for _, e := range reg.Entries() {
		if _, ok := metadata[e.TypeName]; !ok {
			missing = append(missing, e.TypeName)
		}
	}
	assert.Empty(t, missing, "every provider resource type must have an entry in import_metadata.go (importMetadataBase or overrides); add: %v", missing)
}

// TestNewFromResources_staticOverridesReplaceDerived validates that when
// overrides are provided, they replace the derived List/Get/ImportID values.
func TestNewFromResources_staticOverridesReplaceDerived(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)

	overrides := map[string]EntryOverride{
		"criblio_source": {
			ListMethod:     "ListInputCustom",
			GetMethod:      "", // leave derived
			ImportIDFormat: "id",
		},
		"criblio_pipeline": {
			GetMethod: "GetPipelineByIDCustom",
		},
	}

	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), overrides, nil)
	require.NoError(t, err)

	// Override ListMethod and ImportIDFormat; GetMethod stays derived
	e, ok := reg.ByTypeName("criblio_source")
	require.True(t, ok)
	assert.Equal(t, "ListInputCustom", e.ListMethod, "override should replace derived ListMethod")
	assert.Equal(t, "GetInputByID", e.GetMethod, "empty override should not replace derived GetMethod")
	assert.Equal(t, "id", e.ImportIDFormat, "override should replace derived ImportIDFormat")

	// Override only GetMethod
	e, ok = reg.ByTypeName("criblio_pipeline")
	require.True(t, ok)
	assert.Equal(t, "ListPipeline", e.ListMethod, "unchanged")
	assert.Equal(t, "GetPipelineByIDCustom", e.GetMethod, "override should replace derived GetMethod")
	assert.Equal(t, "json:group_id,id", e.ImportIDFormat, "unchanged")
}

// TestImportIDFormat_buildsValidImportID verifies JIRA AC: Import IDs are formatted
// correctly per resource type. For types with ImportIDFormat set, BuildImportID
// produces a valid import ID string.
func TestImportIDFormat_buildsValidImportID(t *testing.T) {
	ctx := context.Background()
	p := provider.New("test")()
	constructors := p.Resources(ctx)
	reg, err := NewFromResources(ctx, constructors, MetadataFromProvider(), nil, nil)
	require.NoError(t, err)

	tests := []struct {
		typeName    string
		identifiers map[string]string
		wantContain []string // substrings that must appear in the built ID
	}{
		{"criblio_source", map[string]string{"group_id": "default", "id": "input-1"}, []string{"group_id", "id", "default", "input-1"}},
		{"criblio_pipeline", map[string]string{"group_id": "default", "id": "pipeline-1"}, []string{"group_id", "id", "default", "pipeline-1"}},
		{"criblio_notification", map[string]string{"id": "notif-1"}, []string{"notif-1"}},
		{"criblio_lookup_file", map[string]string{"group_id": "default", "id": "model_relative_entropy_top_domains.csv"}, []string{"group_id", "id", "default", "model_relative_entropy_top_domains.csv"}},
	}
	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			e, ok := reg.ByTypeName(tt.typeName)
			require.True(t, ok)
			require.NotEmpty(t, e.ImportIDFormat, "%s must have ImportIDFormat for import", tt.typeName)
			id, err := generator.BuildImportID(e.ImportIDFormat, tt.identifiers)
			require.NoError(t, err)
			assert.NotEmpty(t, id)
			for _, sub := range tt.wantContain {
				assert.Contains(t, id, sub, "import ID should contain %q", sub)
			}
		})
	}
}
