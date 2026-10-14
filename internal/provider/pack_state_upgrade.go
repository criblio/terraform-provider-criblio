package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// UpgradeState migrates the original fixed tag object to dynamic pack tags.
func (pack *PackResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				if req.RawState == nil {
					resp.Diagnostics.AddError("Unable to upgrade pack state", "Prior pack state is not available as JSON.")
					return
				}
				upgraded, err := upgradePackTagsStateJSON(req.RawState.JSON)
				if err != nil {
					resp.Diagnostics.AddError("Unable to upgrade pack state", err.Error())
					return
				}
				var schemaResponse resource.SchemaResponse
				pack.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
				typ := schemaResponse.Schema.Type().TerraformType(ctx)
				value, err := tftypes.ValueFromJSON(upgraded, typ)
				if err != nil {
					resp.Diagnostics.AddError("Unable to decode upgraded pack state", err.Error())
					return
				}
				dynamicValue, err := tfprotov6.NewDynamicValue(typ, value)
				if err != nil {
					resp.Diagnostics.AddError("Unable to encode upgraded pack state", err.Error())
					return
				}
				resp.DynamicValue = &dynamicValue
			},
		},
	}
}

func upgradePackTagsStateJSON(raw []byte) ([]byte, error) {
	var state map[string]json.RawMessage
	if err := json.Unmarshal(raw, &state); err != nil {
		return nil, fmt.Errorf("decode prior pack state: %w", err)
	}
	if tags, exists := state["tags"]; exists && string(tags) != "null" {
		attributeTypes := make(map[string]tftypes.Type)
		for _, key := range []string{"data_type", "domain", "streamtags", "technology"} {
			attributeTypes[key] = tftypes.List{ElementType: tftypes.String}
		}
		wrapped, err := json.Marshal(struct {
			Type  tftypes.Type    `json:"type"`
			Value json.RawMessage `json:"value"`
		}{Type: tftypes.Object{AttributeTypes: attributeTypes}, Value: tags})
		if err != nil {
			return nil, fmt.Errorf("encode dynamic pack tags: %w", err)
		}
		state["tags"] = wrapped
	}
	return json.Marshal(state)
}
