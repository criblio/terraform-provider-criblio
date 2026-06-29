package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestTagsAPIMapUsesEmptyArrays(t *testing.T) {
	tags, err := tagsAPIMap(context.Background(), packRequestBodyTagsObjectFromAPI(&packTagsAPI{}))
	if err != nil {
		t.Fatalf("build tags API map: %v", err)
	}

	body, err := json.Marshal(map[string]any{
		"package": map[string]any{
			"tags": tags,
		},
	})
	if err != nil {
		t.Fatalf("marshal tags body: %v", err)
	}

	want := `{"package":{"tags":{"dataType":[],"domain":[],"streamtags":[],"technology":[]}}}`
	if string(body) != want {
		t.Fatalf("unexpected tags body\nwant: %s\n got: %s", want, string(body))
	}
}

func TestTagsAPIMapPreservesConfiguredTags(t *testing.T) {
	tags, err := tagsAPIMap(context.Background(), packRequestBodyTagsObjectFromAPI(&packTagsAPI{
		DataType:   []string{"logs"},
		Domain:     []string{"security"},
		Streamtags: []string{"prod"},
		Technology: []string{"cribl"},
	}))
	if err != nil {
		t.Fatalf("build tags API map: %v", err)
	}

	body, err := json.Marshal(map[string]any{
		"package": map[string]any{
			"tags": tags,
		},
	})
	if err != nil {
		t.Fatalf("marshal tags body: %v", err)
	}

	want := `{"package":{"tags":{"dataType":["logs"],"domain":["security"],"streamtags":["prod"],"technology":["cribl"]}}}`
	if string(body) != want {
		t.Fatalf("unexpected tags body\nwant: %s\n got: %s", want, string(body))
	}
}

func TestPackPlanDecodeAllowsUnknownItems(t *testing.T) {
	stringListType := types.ListType{ElemType: types.StringType}
	tagsType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"data_type":  stringListType,
		"domain":     stringListType,
		"streamtags": stringListType,
		"technology": stringListType,
	}}
	itemType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"author":                 types.StringType,
		"description":            types.StringType,
		"display_name":           types.StringType,
		"exports":                stringListType,
		"id":                     types.StringType,
		"inputs":                 types.Float64Type,
		"min_log_stream_version": types.StringType,
		"outputs":                types.Float64Type,
		"settings":               types.MapType{ElemType: jsontypes.NormalizedType{}},
		"source":                 types.StringType,
		"spec":                   types.StringType,
		"tags":                   tagsType,
		"version":                types.StringType,
		"warnings":               jsontypes.NormalizedType{},
	}}
	packTypes := map[string]attr.Type{
		"allow_custom_functions": types.BoolType,
		"author":                 types.StringType,
		"description":            types.StringType,
		"disabled":               types.BoolType,
		"display_name":           types.StringType,
		"exports":                stringListType,
		"filename":               types.StringType,
		"force":                  types.BoolType,
		"group_id":               types.StringType,
		"id":                     types.StringType,
		"inputs":                 types.Float64Type,
		"items":                  types.ListType{ElemType: itemType},
		"min_log_stream_version": types.StringType,
		"outputs":                types.Float64Type,
		"source":                 types.StringType,
		"spec":                   types.StringType,
		"tags":                   tagsType,
		"version":                types.StringType,
	}
	plan := types.ObjectValueMust(packTypes, map[string]attr.Value{
		"allow_custom_functions": types.BoolNull(),
		"author":                 types.StringNull(),
		"description":            types.StringValue("Pack from source"),
		"disabled":               types.BoolNull(),
		"display_name":           types.StringValue("Search Pack"),
		"exports":                types.ListNull(types.StringType),
		"filename":               types.StringNull(),
		"force":                  types.BoolNull(),
		"group_id":               types.StringValue("default_search"),
		"id":                     types.StringValue("my_search_pack"),
		"inputs":                 types.Float64Null(),
		"items":                  types.ListUnknown(itemType),
		"min_log_stream_version": types.StringNull(),
		"outputs":                types.Float64Null(),
		"source":                 types.StringValue("https://example.com/my-pack.crbl"),
		"spec":                   types.StringNull(),
		"tags":                   types.ObjectNull(tagsType.AttrTypes),
		"version":                types.StringNull(),
	})

	var model *PackResourceModel
	diags := plan.As(context.Background(), &model, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})

	if diags.HasError() {
		t.Fatalf("expected unknown items to decode without diagnostics, got: %v", diags)
	}
	if model == nil {
		t.Fatal("expected decoded pack model")
	}
	if model.Items != nil {
		t.Fatalf("expected unknown items to decode as nil, got %#v", model.Items)
	}
	if got := model.ID.ValueString(); got != "my_search_pack" {
		t.Fatalf("expected id to decode, got %q", got)
	}
}

func TestPreservePackPlanKeepsPlannedItems(t *testing.T) {
	plan := packPlanObjectWithItems(t, packInstallInfo{
		Author:      types.StringNull(),
		Description: types.StringNull(),
		DisplayName: types.StringValue("test-pack-metadata"),
		ID:          types.StringValue("test-pack-metadata"),
		Version:     types.StringValue("1.0.0"),
		Warnings:    jsontypes.NewNormalizedValue("null"),
	})
	data := &PackResourceModel{
		Author:      types.StringValue("Observability Team"),
		Description: types.StringValue("Pack metadata updated"),
		DisplayName: types.StringValue("Pack metadata updated"),
		Items: []packInstallInfo{{
			Author:      types.StringValue("Observability Team"),
			Description: types.StringValue("Pack metadata updated"),
			DisplayName: types.StringValue("Pack metadata updated"),
			ID:          types.StringValue("test-pack-metadata"),
			Version:     types.StringValue("1.0.1"),
			Warnings:    jsontypes.NewNormalizedValue("null"),
		}},
		Version: types.StringValue("1.0.1"),
	}

	preservePackPlan(context.Background(), data, plan)

	if len(data.Items) != 1 {
		t.Fatalf("expected one planned item, got %d", len(data.Items))
	}
	if got := data.Items[0].Version.ValueString(); got != "1.0.0" {
		t.Fatalf("expected planned item version to be preserved, got %q", got)
	}
	if got := data.Items[0].DisplayName.ValueString(); got != "test-pack-metadata" {
		t.Fatalf("expected planned item display name to be preserved, got %q", got)
	}
	if got := data.Version.ValueString(); got != "1.0.1" {
		t.Fatalf("expected configured top-level version to be preserved, got %q", got)
	}
}

func packPlanObjectWithItems(t *testing.T, item packInstallInfo) types.Object {
	t.Helper()

	stringListType := types.ListType{ElemType: types.StringType}
	tagsType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"data_type":  stringListType,
		"domain":     stringListType,
		"streamtags": stringListType,
		"technology": stringListType,
	}}
	itemType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"author":                 types.StringType,
		"description":            types.StringType,
		"display_name":           types.StringType,
		"exports":                stringListType,
		"id":                     types.StringType,
		"inputs":                 types.Float64Type,
		"min_log_stream_version": types.StringType,
		"outputs":                types.Float64Type,
		"settings":               types.MapType{ElemType: jsontypes.NormalizedType{}},
		"source":                 types.StringType,
		"spec":                   types.StringType,
		"tags":                   tagsType,
		"version":                types.StringType,
		"warnings":               jsontypes.NormalizedType{},
	}}
	packTypes := map[string]attr.Type{
		"allow_custom_functions": types.BoolType,
		"author":                 types.StringType,
		"description":            types.StringType,
		"disabled":               types.BoolType,
		"display_name":           types.StringType,
		"exports":                stringListType,
		"filename":               types.StringType,
		"force":                  types.BoolType,
		"group_id":               types.StringType,
		"id":                     types.StringType,
		"inputs":                 types.Float64Type,
		"items":                  types.ListType{ElemType: itemType},
		"min_log_stream_version": types.StringType,
		"outputs":                types.Float64Type,
		"source":                 types.StringType,
		"spec":                   types.StringType,
		"tags":                   tagsType,
		"version":                types.StringType,
	}
	itemValue := types.ObjectValueMust(itemType.AttrTypes, map[string]attr.Value{
		"author":                 item.Author,
		"description":            item.Description,
		"display_name":           item.DisplayName,
		"exports":                types.ListNull(types.StringType),
		"id":                     item.ID,
		"inputs":                 types.Float64Null(),
		"min_log_stream_version": types.StringNull(),
		"outputs":                types.Float64Null(),
		"settings":               types.MapNull(jsontypes.NormalizedType{}),
		"source":                 types.StringNull(),
		"spec":                   types.StringNull(),
		"tags":                   types.ObjectNull(tagsType.AttrTypes),
		"version":                item.Version,
		"warnings":               item.Warnings,
	})

	return types.ObjectValueMust(packTypes, map[string]attr.Value{
		"allow_custom_functions": types.BoolNull(),
		"author":                 types.StringValue("Observability Team"),
		"description":            types.StringValue("Pack metadata updated"),
		"disabled":               types.BoolValue(false),
		"display_name":           types.StringValue("Pack metadata updated"),
		"exports":                types.ListNull(types.StringType),
		"filename":               types.StringNull(),
		"force":                  types.BoolNull(),
		"group_id":               types.StringValue("default"),
		"id":                     types.StringValue("test-pack-metadata"),
		"inputs":                 types.Float64Null(),
		"items":                  types.ListValueMust(itemType, []attr.Value{itemValue}),
		"min_log_stream_version": types.StringNull(),
		"outputs":                types.Float64Null(),
		"source":                 types.StringNull(),
		"spec":                   types.StringNull(),
		"tags":                   types.ObjectNull(tagsType.AttrTypes),
		"version":                types.StringValue("1.0.1"),
	})
}
