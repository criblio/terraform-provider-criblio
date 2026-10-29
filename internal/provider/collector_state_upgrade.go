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

// UpgradeState migrates legacy schedule ranges and string-map collector state.
func (r *CollectorResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: collectorStateUpgrader(r, true),
		1: collectorStateUpgrader(r, false),
	}
}

// UpgradeState migrates pack collectors written with string-map saved_state.
func (r *PackCollectorResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{0: collectorStateUpgrader(r, false)}
}

func collectorStateUpgrader(res resource.Resource, upgradeSchedule bool) resource.StateUpgrader {
	return resource.StateUpgrader{StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
		if req.RawState == nil || len(req.RawState.JSON) == 0 {
			resp.Diagnostics.AddError("Unable to upgrade collector state", "Prior collector state is not available as JSON.")
			return
		}
		raw := req.RawState.JSON
		var err error
		if upgradeSchedule {
			raw, err = upgradeCollectorStateJSON(raw)
		} else {
			raw, err = upgradeCollectorSavedStateJSON(raw)
		}
		if err != nil {
			resp.Diagnostics.AddError("Unable to upgrade collector state", err.Error())
			return
		}
		var schema resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schema)
		typ := schema.Schema.Type().TerraformType(ctx)
		value, err := tftypes.ValueFromJSON(raw, typ)
		if err != nil {
			resp.Diagnostics.AddError("Unable to decode upgraded collector state", err.Error())
			return
		}
		dynamicValue, err := tfprotov6.NewDynamicValue(typ, value)
		if err != nil {
			resp.Diagnostics.AddError("Unable to encode upgraded collector state", err.Error())
			return
		}
		resp.DynamicValue = &dynamicValue
	}}
}

func upgradeCollectorSavedStateJSON(raw []byte) ([]byte, error) {
	var state map[string]json.RawMessage
	if err := json.Unmarshal(raw, &state); err != nil {
		return nil, fmt.Errorf("decode prior collector state: %w", err)
	}
	for name, data := range state {
		if !strings.HasPrefix(name, "input_collector_") || string(data) == "null" {
			continue
		}
		var block map[string]json.RawMessage
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, fmt.Errorf("decode prior %s: %w", name, err)
		}
		savedState, ok := block["saved_state"]
		if !ok || string(savedState) == "null" {
			continue
		}
		wrapped, err := json.Marshal(struct {
			Type  tftypes.Type    `json:"type"`
			Value json.RawMessage `json:"value"`
		}{Type: tftypes.Map{ElementType: tftypes.String}, Value: savedState})
		if err != nil {
			return nil, fmt.Errorf("encode dynamic saved_state: %w", err)
		}
		block["saved_state"] = wrapped
		encoded, err := json.Marshal(block)
		if err != nil {
			return nil, err
		}
		state[name] = encoded
	}
	return json.Marshal(state)
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
	encoded, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	return upgradeCollectorSavedStateJSON(encoded)
}
