package export

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHclOptionsForType_criblioDestination_skipsHoistedRootAttrs(t *testing.T) {
	e := registry.Entry{
		OneOf: &registry.OneOfConfig{ReadOnlyAttr: "items"},
	}
	opts := hclOptionsForType("criblio_destination", e)
	require.NotNil(t, opts)
	require.NotNil(t, opts.SkipAttributes)
	assert.True(t, opts.SkipAttributes["items"])
	assert.True(t, opts.SkipAttributes["environment"])
	assert.True(t, opts.SkipAttributes["pipeline"])
	assert.True(t, opts.SkipAttributes["type"])
}
