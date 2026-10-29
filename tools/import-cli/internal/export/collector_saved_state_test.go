package export

import (
	"encoding/json"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestCollectorExportOmitsRuntimeSavedState(t *testing.T) {
	for _, typeName := range []string{"criblio_collector", "criblio_pack_collector"} {
		var model any = &provider.CollectorModel{}
		if typeName == "criblio_pack_collector" {
			model = &provider.PackCollectorModel{}
		}
		require.NoError(t, json.Unmarshal([]byte(`{"id":"events","collector":{"type":"database"},"savedState":{"databaseCollectorState":{"data":{"id":{"value":3}}}}}`), model))
		attrs, err := hcl.ModelToValue(model, hclOptionsForType(typeName, registry.Entry{}))
		require.NoError(t, err)
		require.Contains(t, attrs, "input_collector_database")
		require.NotContains(t, attrs["input_collector_database"].Map, "saved_state")
		oneOf := &registry.OneOfConfig{
			DiscriminatorField: "type", NestedDiscriminatorField: "collector.type",
			BlockNamePrefix: "input_collector_", SupportedBlockNames: []string{"input_collector_database"},
		}
		require.NoError(t, addOneOfBlockFromFirstItem(model, attrs, oneOf))
		require.NotContains(t, attrs["input_collector_database"].Map, "saved_state")
		require.Contains(t, attrs["input_collector_database"].Map, "collector")
	}
}
