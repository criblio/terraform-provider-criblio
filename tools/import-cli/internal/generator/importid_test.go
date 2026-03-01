package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildImportID(t *testing.T) {
	t.Run("id_only", func(t *testing.T) {
		id, err := BuildImportID("id", map[string]string{"id": "my-id"})
		require.NoError(t, err)
		assert.Equal(t, "my-id", id)
	})
	t.Run("json_format", func(t *testing.T) {
		id, err := BuildImportID("json:group_id,id", map[string]string{"group_id": "default", "id": "input-1"})
		require.NoError(t, err)
		assert.Contains(t, id, "group_id")
		assert.Contains(t, id, "id")
		assert.Contains(t, id, "default")
		assert.Contains(t, id, "input-1")
	})
	t.Run("json_deterministic", func(t *testing.T) {
		m := map[string]string{"group_id": "g", "id": "i"}
		a, err := BuildImportID("json:group_id,id", m)
		require.NoError(t, err)
		b, err := BuildImportID("json:id,group_id", m)
		require.NoError(t, err)
		assert.Equal(t, a, b, "JSON keys should be sorted for deterministic output")
	})
	t.Run("empty_format_error", func(t *testing.T) {
		_, err := BuildImportID("", map[string]string{"id": "x"})
		require.Error(t, err)
	})
}
