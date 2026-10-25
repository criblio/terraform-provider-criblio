package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSearchEngineTierBoundsAreReadOnly(t *testing.T) {
	var response resource.SchemaResponse
	(&SearchEngineResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)
	for _, name := range []string{"minimum_tier_size", "maximum_tier_size"} {
		attribute := response.Schema.Attributes[name]
		if !attribute.IsComputed() || attribute.IsOptional() || attribute.IsRequired() {
			t.Errorf("%s must be computed only", name)
		}
	}
}

func TestSearchEngineTierBoundsUseAPIValues(t *testing.T) {
	var api SearchEngineModel
	if err := json.Unmarshal([]byte(`{"minimumTierSize":"small","maximumTierSize":"large"}`), &api); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name  string
		state SearchEngineModel
	}{
		{
			name: "omitted",
			state: SearchEngineModel{
				MinimumTierSize: types.StringNull(),
				MaximumTierSize: types.StringNull(),
			},
		},
		{
			name: "previous bounds",
			state: SearchEngineModel{
				MinimumTierSize: types.StringValue("xsmall"),
				MaximumTierSize: types.StringValue("xlarge"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.state.TierSize = types.StringValue("medium")
			// Create, ordinary Read, and Update preserve configured inputs.
			applySearchEngineAPIToState(&api, &tt.state, true, false)
			if !tt.state.MinimumTierSize.Equal(api.MinimumTierSize) || !tt.state.MaximumTierSize.Equal(api.MaximumTierSize) {
				t.Fatalf("API bounds not reflected: minimum=%s maximum=%s", tt.state.MinimumTierSize, tt.state.MaximumTierSize)
			}
			data, err := json.Marshal(tt.state)
			if err != nil {
				t.Fatal(err)
			}
			var body map[string]any
			if err := json.Unmarshal(data, &body); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"minimumTierSize", "maximumTierSize"} {
				if _, present := body[name]; present {
					t.Errorf("read-only %s leaked into request: %s", name, data)
				}
			}
			if body["tierSize"] != "medium" {
				t.Errorf("configured tierSize lost from request: %s", data)
			}
		})
	}
}
