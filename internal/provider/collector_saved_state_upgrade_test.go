package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"
)

func TestCollectorSavedStateUpgrade(t *testing.T) {
	ctx := context.Background()
	for _, res := range []resource.ResourceWithUpgradeState{&CollectorResource{}, &PackCollectorResource{}} {
		var response resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &response)
		for version, upgrader := range res.UpgradeState(ctx) {
			for name, attribute := range response.Schema.Attributes {
				variant, ok := attribute.(schema.SingleNestedAttribute)
				if !ok {
					continue
				}
				for _, legacy := range []map[string]any{nil, {}, {"cursor": "configured", "nullable": nil}} {
					state := make(map[string]any)
					for key := range response.Schema.Attributes {
						state[key] = nil
					}
					state["id"] = "events"
					state["group_id"] = "default"
					block := make(map[string]any)
					for key := range variant.Attributes {
						block[key] = nil
					}
					block["saved_state"] = legacy
					state[name] = block
					raw, err := json.Marshal(state)
					require.NoError(t, err)
					var upgraded resource.UpgradeStateResponse
					upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{RawState: &tfprotov6.RawState{JSON: raw}}, &upgraded)
					require.False(t, upgraded.Diagnostics.HasError(), "version %d, %s: %v", version, name, upgraded.Diagnostics)
					require.NotNil(t, upgraded.DynamicValue)
					value, err := upgraded.DynamicValue.Unmarshal(response.Schema.Type().TerraformType(ctx))
					require.NoError(t, err)
					decoded := tfsdk.State{Schema: response.Schema, Raw: value}
					var actual types.Dynamic
					require.False(t, decoded.GetAttribute(ctx, path.Root(name).AtName("saved_state"), &actual).HasError())
					if legacy == nil {
						require.True(t, actual.IsNull())
						continue
					}
					elements := make(map[string]attr.Value)
					for key, value := range legacy {
						if value == nil {
							elements[key] = types.StringNull()
						} else {
							elements[key] = types.StringValue(value.(string))
						}
					}
					want := types.DynamicValue(types.MapValueMust(types.StringType, elements))
					require.True(t, actual.Equal(want), "upgraded saved_state = %v, want %v", actual, want)
				}
			}
		}
	}
}

func TestCollectorSavedStatePreservesLargeDatabaseID(t *testing.T) {
	var model CollectorModel
	err := json.Unmarshal([]byte(`{"collector":{"type":"database"},"savedState":{"databaseCollectorState":{"data":{"id":{"value":9007199254740993}}}}}`), &model)
	require.NoError(t, err)
	state, err := dynamicValueToJSON(model.InputCollectorDatabase.SavedState)
	require.NoError(t, err)
	encoded, err := json.Marshal(state)
	require.NoError(t, err)
	require.JSONEq(t, `{"databaseCollectorState":{"data":{"id":{"value":9007199254740993}}}}`, string(encoded))
	require.Contains(t, string(encoded), "9007199254740993")
}
