package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestPackTagsStateUpgrade(t *testing.T) {
	ctx := context.Background()
	pack := &PackResource{}
	var schemaResponse resource.SchemaResponse
	pack.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	for _, tags := range []any{
		nil,
		map[string]any{"data_type": []string{"logs"}, "domain": []string{"observability"}, "streamtags": []string{"processing"}, "technology": []string{"aws"}},
	} {
		state := make(map[string]any)
		for key := range schemaResponse.Schema.Attributes {
			state[key] = nil
		}
		state["id"] = "test"
		state["group_id"] = "default"
		state["tags"] = tags
		raw, err := json.Marshal(state)
		require.NoError(t, err)
		var response resource.UpgradeStateResponse
		pack.UpgradeState(ctx)[0].StateUpgrader(ctx, resource.UpgradeStateRequest{RawState: &tfprotov6.RawState{JSON: raw}}, &response)
		require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
		require.NotNil(t, response.DynamicValue)
		value, err := response.DynamicValue.Unmarshal(schemaResponse.Schema.Type().TerraformType(ctx))
		require.NoError(t, err)
		var attributes map[string]tftypes.Value
		require.NoError(t, value.As(&attributes))
		require.Equal(t, tags == nil, attributes["tags"].IsNull())
		if tags != nil {
			var tagValues map[string]tftypes.Value
			require.NoError(t, attributes["tags"].As(&tagValues))
			var dataType []tftypes.Value
			require.NoError(t, tagValues["data_type"].As(&dataType))
			var actual string
			require.NoError(t, dataType[0].As(&actual))
			require.Equal(t, "logs", actual)
		}
	}
	_, err := upgradePackTagsStateJSON([]byte(`invalid`))
	require.Error(t, err)
	var missing resource.UpgradeStateResponse
	pack.UpgradeState(ctx)[0].StateUpgrader(ctx, resource.UpgradeStateRequest{}, &missing)
	require.True(t, missing.Diagnostics.HasError())
}
