package provider

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSavedQueryVisualizationDisplayConfigRoundTrip(t *testing.T) {
	for _, test := range []struct {
		name    string
		display string
	}{
		{"nested values", `{"colors":{"name":"categorical.observable10"},"opacity":0.5,"enabled":true,"palette":["red","blue"],"missing":null}`},
		{"string values", `{"orientation":"vertical"}`},
		{"empty", `{}`},
		{"null", `null`},
		{"omitted", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			visualization := map[string]any{
				"family":       "chart",
				"fieldsConfig": map[string]any{"fields": []any{map[string]any{"name": "_time", "type": "timestamp"}}},
			}
			var want any
			if test.display != "" {
				if err := json.Unmarshal([]byte(test.display), &want); err != nil {
					t.Fatal(err)
				}
				visualization["displayConfig"] = want
			}
			payload, err := json.Marshal(map[string]any{"id": "saved-query", "chartConfig": map[string]any{"visualization": visualization}})
			if err != nil {
				t.Fatal(err)
			}
			var api SearchSavedQueryModel
			if err := json.Unmarshal(payload, &api); err != nil {
				t.Fatalf("decode API response: %v", err)
			}
			var state SearchSavedQueryModel
			applySearchSavedQueryAPIToState(&api, &state, false, true)
			chart := state.ChartConfig.Attributes()
			viz := chart["visualization"].(types.Object)
			display, ok := viz.Attributes()["display_config"].(jsontypes.Normalized)
			if !ok {
				t.Fatal("display_config must use normalized JSON")
			}
			if display.IsNull() != (want == nil) {
				t.Fatalf("display_config null = %t, want %t", display.IsNull(), want == nil)
			}
			encoded, err := json.Marshal(state)
			if err != nil {
				t.Fatalf("encode update payload: %v", err)
			}
			var output struct {
				ChartConfig struct {
					Visualization map[string]any `json:"visualization"`
				} `json:"chartConfig"`
			}
			if err := json.Unmarshal(encoded, &output); err != nil {
				t.Fatal(err)
			}
			if got := output.ChartConfig.Visualization["displayConfig"]; !reflect.DeepEqual(got, want) {
				t.Fatalf("displayConfig = %#v, want %#v", got, want)
			}
		})
	}
}
