package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// UpgradeState migrates collector state written before schedule time ranges supported strings.
func (r *CollectorResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				if req.RawState == nil || len(req.RawState.JSON) == 0 {
					resp.Diagnostics.AddError("Unable to upgrade collector state", "Prior collector state is not available as JSON.")
					return
				}
				upgradedJSON, err := upgradeCollectorStateJSON(req.RawState.JSON)
				if err != nil {
					resp.Diagnostics.AddError("Unable to upgrade collector state", err.Error())
					return
				}
				var schemaResp resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
				terraformType := schemaResp.Schema.Type().TerraformType(ctx)
				value, err := tftypes.ValueFromJSON(upgradedJSON, terraformType)
				if err != nil {
					resp.Diagnostics.AddError("Unable to decode upgraded collector state", err.Error())
					return
				}
				dynamicValue, err := tfprotov6.NewDynamicValue(terraformType, value)
				if err != nil {
					resp.Diagnostics.AddError("Unable to encode upgraded collector state", err.Error())
					return
				}
				resp.DynamicValue = &dynamicValue
			},
		},
	}
}

func upgradeCollectorStateJSON(raw []byte) ([]byte, error) {
	var state map[string]any
	if err := json.Unmarshal(raw, &state); err != nil {
		return nil, fmt.Errorf("decode prior state: %w", err)
	}
	for key, value := range state {
		if !strings.HasPrefix(key, "input_collector_") {
			continue
		}
		block, ok := value.(map[string]any)
		if !ok {
			continue
		}
		schedule, ok := block["schedule"].(map[string]any)
		if !ok {
			continue
		}
		run, ok := schedule["run"].(map[string]any)
		if !ok {
			continue
		}
		for _, field := range []string{"earliest", "latest"} {
			if number, ok := run[field].(float64); ok {
				run[field] = strconv.FormatFloat(number, 'f', -1, 64)
			}
		}
	}
	return json.Marshal(state)
}
