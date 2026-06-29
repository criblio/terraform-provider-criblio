package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

func TestGroupAPIFromModelNormalizesInheritedGroupAsFleetForAPI(t *testing.T) {
	model := &GroupResourceModel{
		ID:       types.StringValue("my-edge-subfleet"),
		Name:     types.StringValue("my-edge-subfleet"),
		Inherits: types.StringValue("default_fleet"),
		IsFleet:  types.BoolValue(false),
		OnPrem:   types.BoolValue(true),
		Product:  types.StringValue("edge"),
		Type:     types.StringValue("edge"),
	}

	apiModel := groupAPIFromModel(model)

	if apiModel.IsFleet == nil {
		t.Fatal("expected isFleet to be set")
	}
	if !*apiModel.IsFleet {
		t.Fatal("expected inherited group API payload to force isFleet=true")
	}
	if model.IsFleet.ValueBool() {
		t.Fatal("expected Terraform model to preserve configured is_fleet=false")
	}
}

func TestGroupPlanDecodeAllowsUnknownCloud(t *testing.T) {
	cloudTypes := map[string]attr.Type{
		"provider": types.StringType,
		"region":   types.StringType,
	}
	groupTypes := map[string]attr.Type{
		"cloud":                 types.ObjectType{AttrTypes: cloudTypes},
		"description":           types.StringType,
		"estimated_ingest_rate": types.Float64Type,
		"id":                    types.StringType,
		"inherits":              types.StringType,
		"is_fleet":              types.BoolType,
		"max_worker_age":        types.StringType,
		"name":                  types.StringType,
		"on_prem":               types.BoolType,
		"product":               types.StringType,
		"provisioned":           types.BoolType,
		"streamtags":            types.ListType{ElemType: types.StringType},
		"tags":                  types.StringType,
		"type":                  types.StringType,
		"worker_remote_access":  types.BoolType,
	}
	plan := types.ObjectValueMust(groupTypes, map[string]attr.Value{
		"cloud":                 types.ObjectUnknown(cloudTypes),
		"description":           types.StringNull(),
		"estimated_ingest_rate": types.Float64Null(),
		"id":                    types.StringValue("my-hybrid-group"),
		"inherits":              types.StringNull(),
		"is_fleet":              types.BoolValue(false),
		"max_worker_age":        types.StringValue("2h"),
		"name":                  types.StringValue("my-hybrid-group"),
		"on_prem":               types.BoolValue(true),
		"product":               types.StringValue("stream"),
		"provisioned":           types.BoolValue(false),
		"streamtags":            types.ListValueMust(types.StringType, []attr.Value{types.StringValue("datacenter1")}),
		"tags":                  types.StringNull(),
		"type":                  types.StringNull(),
		"worker_remote_access":  types.BoolValue(false),
	})

	var model *GroupResourceModel
	diags := plan.As(context.Background(), &model, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})

	if diags.HasError() {
		t.Fatalf("expected unknown cloud to decode without diagnostics, got: %v", diags)
	}
	if model == nil {
		t.Fatal("expected decoded group model")
	}
	if model.Cloud != nil {
		t.Fatalf("expected unknown cloud to decode as nil, got %#v", model.Cloud)
	}
	if got := model.ID.ValueString(); got != "my-hybrid-group" {
		t.Fatalf("expected id to decode, got %q", got)
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
