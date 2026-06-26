package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRoutesModelUnmarshalClonesObject(t *testing.T) {
	var model RoutesModel
	err := json.Unmarshal([]byte(`{
		"id": "default",
		"routes": [
			{
				"name": "with clones",
				"pipeline": "main",
				"clones": [
					{"__cloneId": "audit"}
				]
			}
		]
	}`), &model)
	if err != nil {
		t.Fatalf("UnmarshalJSON returned error: %v", err)
	}

	if model.Routes.IsNull() || model.Routes.IsUnknown() || len(model.Routes.Elements()) != 1 {
		t.Fatalf("routes = %#v", model.Routes)
	}
	route, ok := model.Routes.Elements()[0].(types.Object)
	if !ok {
		t.Fatalf("route element type = %T", model.Routes.Elements()[0])
	}
	clones, ok := route.Attributes()["clones"].(types.List)
	if !ok {
		t.Fatalf("clones attribute type = %T", route.Attributes()["clones"])
	}
	if clones.IsNull() || clones.IsUnknown() || len(clones.Elements()) != 1 {
		t.Fatalf("clones = %#v", clones)
	}
	if _, ok := clones.Elements()[0].(types.Map); !ok {
		t.Fatalf("clone element type = %T", clones.Elements()[0])
	}
}

func TestRoutesModelUpdateBodyEmitsEmptyGroupsAndComments(t *testing.T) {
	model := RoutesModel{
		Comments: types.ListNull(types.ObjectType{AttrTypes: RoutesCommentsAttrTypes()}),
		Groups:   types.MapNull(types.ObjectType{AttrTypes: RoutesGroupsAttrTypes()}),
		Routes:   types.ListNull(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}),
	}

	body, err := model.updateBody()
	if err != nil {
		t.Fatalf("updateBody returned error: %v", err)
	}

	comments, ok := body["comments"].([]any)
	if !ok {
		t.Fatalf("comments type = %T", body["comments"])
	}
	if len(comments) != 0 {
		t.Fatalf("comments = %#v", comments)
	}
	groups, ok := body["groups"].(map[string]any)
	if !ok {
		t.Fatalf("groups type = %T", body["groups"])
	}
	if len(groups) != 0 {
		t.Fatalf("groups = %#v", groups)
	}
}

func TestApplyRoutesAPIToStatePreservesNullInputsAfterWrite(t *testing.T) {
	api := RoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: RoutesCommentsAttrTypes()}, nil),
		Groups:   types.MapValueMust(types.ObjectType{AttrTypes: RoutesGroupsAttrTypes()}, nil),
		ID:       types.StringValue("default"),
		Routes:   types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}
	state := RoutesModel{
		Comments: types.ListNull(types.ObjectType{AttrTypes: RoutesCommentsAttrTypes()}),
		Groups:   types.MapNull(types.ObjectType{AttrTypes: RoutesGroupsAttrTypes()}),
		ID:       types.StringValue("default"),
		Routes:   types.ListNull(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}),
	}

	applyRoutesAPIToState(&api, &state, true, false)

	if !state.Comments.IsNull() {
		t.Fatalf("comments = %#v, want null", state.Comments)
	}
	if !state.Groups.IsNull() {
		t.Fatalf("groups = %#v, want null", state.Groups)
	}
	if !state.Routes.IsNull() {
		t.Fatalf("routes = %#v, want null", state.Routes)
	}
}

func TestApplyRoutesAPIToStateFillsUnknownRouteAttributesAfterWrite(t *testing.T) {
	routeTypes := RoutesRoutesAttrTypes()
	stateRoute := types.ObjectValueMust(routeTypes, map[string]attr.Value{
		"clones":                   types.ListUnknown(types.MapType{ElemType: types.StringType}),
		"context":                  types.StringNull(),
		"description":              types.StringUnknown(),
		"disabled":                 types.BoolUnknown(),
		"enable_output_expression": types.BoolUnknown(),
		"filter":                   types.StringValue("true"),
		"group_id":                 types.StringValue("default"),
		"name":                     types.StringValue("with pack"),
		"output":                   types.StringNull(),
		"output_expression":        types.StringNull(),
		"pipeline":                 types.StringValue("main"),
		"target_context":           types.StringNull(),
		"final":                    types.BoolValue(true),
		"id":                       types.StringUnknown(),
	})
	apiRoute := types.ObjectValueMust(routeTypes, map[string]attr.Value{
		"clones":                   types.ListValueMust(types.MapType{ElemType: types.StringType}, nil),
		"context":                  types.StringNull(),
		"description":              types.StringValue(""),
		"disabled":                 types.BoolValue(false),
		"enable_output_expression": types.BoolValue(false),
		"filter":                   types.StringValue("true"),
		"group_id":                 types.StringValue("default"),
		"name":                     types.StringValue("with pack"),
		"output":                   types.StringNull(),
		"output_expression":        types.StringNull(),
		"pipeline":                 types.StringValue("main"),
		"target_context":           types.StringNull(),
		"final":                    types.BoolValue(true),
		"id":                       types.StringValue("route-generated"),
	})
	api := RoutesModel{
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: routeTypes}, []attr.Value{apiRoute}),
	}
	state := RoutesModel{
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: routeTypes}, []attr.Value{stateRoute}),
	}

	applyRoutesAPIToState(&api, &state, true, false)

	route, ok := state.Routes.Elements()[0].(types.Object)
	if !ok {
		t.Fatalf("route element type = %T", state.Routes.Elements()[0])
	}
	for _, name := range []string{"clones", "description", "disabled", "enable_output_expression", "id"} {
		if route.Attributes()[name].IsUnknown() {
			t.Fatalf("%s is still unknown", name)
		}
	}
	if got := route.Attributes()["output"]; !got.IsNull() {
		t.Fatalf("output = %#v, want null", got)
	}
}
