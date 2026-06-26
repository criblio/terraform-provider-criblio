package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGroupAPIFromModelNormalizesEdgeOnPremForAPI(t *testing.T) {
	model := &GroupResourceModel{
		ID:      types.StringValue("my-edge-fleet"),
		Name:    types.StringValue("my-edge-fleet"),
		OnPrem:  types.BoolValue(false),
		Product: types.StringValue("edge"),
		Type:    types.StringValue("edge"),
	}

	apiModel := groupAPIFromModel(model)

	if apiModel.OnPrem == nil {
		t.Fatal("expected onPrem to be set")
	}
	if !*apiModel.OnPrem {
		t.Fatal("expected edge group API payload to force onPrem=true")
	}
	if model.OnPrem.ValueBool() {
		t.Fatal("expected Terraform model to preserve configured on_prem=false")
	}
}

func TestGroupAPIFromModelKeepsStreamOnPremValue(t *testing.T) {
	model := &GroupResourceModel{
		ID:      types.StringValue("my-stream-group"),
		Name:    types.StringValue("my-stream-group"),
		OnPrem:  types.BoolValue(false),
		Product: types.StringValue("stream"),
		Type:    types.StringValue("stream"),
	}

	apiModel := groupAPIFromModel(model)

	if apiModel.OnPrem == nil {
		t.Fatal("expected onPrem to be set")
	}
	if *apiModel.OnPrem {
		t.Fatal("expected stream group API payload to preserve onPrem=false")
	}
}

func TestPreserveLegacyEdgeOnPremKeepsPriorFalse(t *testing.T) {
	model := &GroupResourceModel{
		OnPrem:  types.BoolValue(true),
		Product: types.StringValue("edge"),
		Type:    types.StringValue("edge"),
	}

	preserveLegacyEdgeOnPrem(types.BoolValue(false), model)

	if model.OnPrem.ValueBool() {
		t.Fatal("expected refresh to preserve legacy on_prem=false for edge group")
	}
}

func TestPreserveLegacyEdgeOnPremLeavesStreamAPIValue(t *testing.T) {
	model := &GroupResourceModel{
		OnPrem:  types.BoolValue(true),
		Product: types.StringValue("stream"),
		Type:    types.StringValue("stream"),
	}

	preserveLegacyEdgeOnPrem(types.BoolValue(false), model)

	if !model.OnPrem.ValueBool() {
		t.Fatal("expected stream refresh to keep API onPrem=true")
	}
}

func TestPreserveLegacyEdgeIsFleetKeepsPriorFalse(t *testing.T) {
	model := &GroupResourceModel{
		IsFleet: types.BoolValue(true),
		Product: types.StringValue("edge"),
		Type:    types.StringValue("edge"),
	}

	preserveLegacyEdgeIsFleet(types.BoolValue(false), model)

	if model.IsFleet.ValueBool() {
		t.Fatal("expected refresh to preserve legacy is_fleet=false for edge subfleet")
	}
}

func TestPreserveLegacyEdgeIsFleetLeavesStreamAPIValue(t *testing.T) {
	model := &GroupResourceModel{
		IsFleet: types.BoolValue(true),
		Product: types.StringValue("stream"),
		Type:    types.StringValue("stream"),
	}

	preserveLegacyEdgeIsFleet(types.BoolValue(false), model)

	if !model.IsFleet.ValueBool() {
		t.Fatal("expected stream refresh to keep API isFleet=true")
	}
}
