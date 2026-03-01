package converter

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInjectRequiredIdentifiers(t *testing.T) {
	t.Run("sets_id_and_group_id_on_source_model", func(t *testing.T) {
		model := &provider.SourceResourceModel{}
		identifiers := map[string]string{"GroupID": "default", "ID": "my-input-1"}
		err := InjectRequiredIdentifiers(model, identifiers)
		require.NoError(t, err)
		assert.Equal(t, "default", model.GroupID.ValueString())
		assert.Equal(t, "my-input-1", model.ID.ValueString())
	})

	t.Run("sets_id_and_group_id_on_pipeline_model", func(t *testing.T) {
		model := &provider.PipelineResourceModel{}
		identifiers := map[string]string{"GroupID": "stream-1", "ID": "pipeline-abc"}
		err := InjectRequiredIdentifiers(model, identifiers)
		require.NoError(t, err)
		assert.Equal(t, "stream-1", model.GroupID.ValueString())
		assert.Equal(t, "pipeline-abc", model.ID.ValueString())
	})

	t.Run("ignores_unknown_params", func(t *testing.T) {
		model := &provider.SourceResourceModel{}
		identifiers := map[string]string{"Unknown": "x", "ID": "id1"}
		err := InjectRequiredIdentifiers(model, identifiers)
		require.NoError(t, err)
		assert.Equal(t, "id1", model.ID.ValueString())
		assert.True(t, model.GroupID.IsNull())
	})

	t.Run("nil_or_empty_identifiers_no_error", func(t *testing.T) {
		model := &provider.SourceResourceModel{}
		require.NoError(t, InjectRequiredIdentifiers(model, nil))
		require.NoError(t, InjectRequiredIdentifiers(model, map[string]string{}))
		require.NoError(t, InjectRequiredIdentifiers(nil, map[string]string{"ID": "x"}))
	})

	t.Run("nil_model_no_panic", func(t *testing.T) {
		err := InjectRequiredIdentifiers(nil, map[string]string{"ID": "x"})
		require.NoError(t, err)
	})

	t.Run("non_struct_returns_error", func(t *testing.T) {
		err := InjectRequiredIdentifiers("not a struct", map[string]string{"ID": "x"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be a struct")
	})
}

func TestInjectRequiredIdentifiers_only_sets_types_string(t *testing.T) {
	// Model has ID/GroupID as types.String; we only set those.
	model := &provider.SourceResourceModel{}
	identifiers := map[string]string{"GroupID": "g", "ID": "i"}
	err := InjectRequiredIdentifiers(model, identifiers)
	require.NoError(t, err)
	assert.True(t, model.GroupID.Equal(types.StringValue("g")))
	assert.True(t, model.ID.Equal(types.StringValue("i")))
}
