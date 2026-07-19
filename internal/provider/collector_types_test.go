package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCollectorMarshalJSONIncludesSavedJobType(t *testing.T) {
	model := CollectorModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("rest-api-demo-collector"),
		InputCollectorRest: &InputCollectorRestModel{
			ID: types.StringValue("rest-api-demo-collector"),
		},
	}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if got["type"] != "collection" {
		t.Fatalf("expected saved job type collection, got %#v in %s", got["type"], data)
	}
	if got["id"] != "rest-api-demo-collector" {
		t.Fatalf("expected collector id in body, got %#v in %s", got["id"], data)
	}
}
