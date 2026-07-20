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

func TestIsRoutesImportStateTreatsEmptyRoutesAsImport(t *testing.T) {
	routeType := types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}
	state := RoutesModel{
		Routes: types.ListValueMust(routeType, nil),
	}
	if !isRoutesImportState(&state) {
		t.Fatalf("expected empty routes list to be treated as import state")
	}
}

func TestIsPackRoutesImportStateTreatsEmptyRoutesAsImport(t *testing.T) {
	routeType := types.ObjectType{AttrTypes: PackRoutesRoutesAttrTypes()}
	state := PackRoutesModel{
		Routes: types.ListValueMust(routeType, nil),
	}
	if !isPackRoutesImportState(&state) {
		t.Fatalf("expected empty pack routes list to be treated as import state")
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

func TestRoutesModelUpdateBodyNormalizesCommentsGroupsAndRouteGroupID(t *testing.T) {
	commentAttrs := map[string]attr.Value{
		"comment":  types.StringValue("my-comment"),
		"group_id": types.StringNull(),
		"id":       types.StringNull(),
		"index":    types.Int64Null(),
	}
	groupAttrsA := map[string]attr.Value{
		"description": types.StringValue("group A"),
		"index":       types.Int64Null(),
		"name":        types.StringValue("first"),
	}
	groupAttrsB := map[string]attr.Value{
		"description": types.StringValue("group B"),
		"index":       types.Int64Value(9),
		"name":        types.StringValue("second"),
	}
	routeAttrs := map[string]attr.Value{
		"clones":                   types.ListNull(types.MapType{ElemType: types.StringType}),
		"context":                  types.StringNull(),
		"description":              types.StringNull(),
		"disabled":                 types.BoolNull(),
		"enable_output_expression": types.BoolNull(),
		"filter":                   types.StringNull(),
		"group_id":                 types.StringNull(),
		"name":                     types.StringValue("my-route"),
		"output":                   types.StringNull(),
		"output_expression":        types.StringNull(),
		"pipeline":                 types.StringValue("main"),
		"target_context":           types.StringNull(),
		"final":                    types.BoolNull(),
		"id":                       types.StringNull(),
	}
	model := RoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: RoutesCommentsAttrTypes()}, []attr.Value{
			types.ObjectValueMust(RoutesCommentsAttrTypes(), commentAttrs),
		}),
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: RoutesGroupsAttrTypes()}, map[string]attr.Value{
			"mygroup_b": types.ObjectValueMust(RoutesGroupsAttrTypes(), groupAttrsB),
			"mygroup_a": types.ObjectValueMust(RoutesGroupsAttrTypes(), groupAttrsA),
		}),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, []attr.Value{
			types.ObjectValueMust(RoutesRoutesAttrTypes(), routeAttrs),
		}),
	}

	body, err := model.updateBody()
	if err != nil {
		t.Fatalf("updateBody returned error: %v", err)
	}

	comments, ok := body["comments"].([]any)
	if !ok || len(comments) != 1 {
		t.Fatalf("comments = %#v", body["comments"])
	}
	commentMap, ok := comments[0].(map[string]any)
	if !ok {
		t.Fatalf("comment type = %T", comments[0])
	}
	if got := commentMap["id"]; got != "tf-comment-0" {
		t.Fatalf("comment id = %#v", got)
	}
	if got := commentMap["index"]; got != 0 {
		t.Fatalf("comment index = %#v", got)
	}

	groups, ok := body["groups"].(map[string]any)
	if !ok {
		t.Fatalf("groups type = %T", body["groups"])
	}
	groupA, ok := groups["mygroup_a"].(map[string]any)
	if !ok {
		t.Fatalf("group mygroup_a = %#v", groups["mygroup_a"])
	}
	if got := groupA["index"]; got != 0 {
		t.Fatalf("group mygroup_a index = %#v", got)
	}
	groupB, ok := groups["mygroup_b"].(map[string]any)
	if !ok {
		t.Fatalf("group mygroup_b = %#v", groups["mygroup_b"])
	}
	if got := groupB["index"]; got != int64(9) {
		t.Fatalf("group mygroup_b index = %#v", got)
	}

	routes, ok := body["routes"].([]any)
	if !ok || len(routes) != 1 {
		t.Fatalf("routes = %#v", body["routes"])
	}
	route, ok := routes[0].(map[string]any)
	if !ok {
		t.Fatalf("route type = %T", routes[0])
	}
	if got := route["groupId"]; got != "default" {
		t.Fatalf("route groupId = %#v", got)
	}
}

func TestApplyRoutesAPIToStateImportDefaultsMissingRouteGroupID(t *testing.T) {
	routeTypes := RoutesRoutesAttrTypes()
	apiRoute := types.ObjectValueMust(routeTypes, map[string]attr.Value{
		"clones":                   types.ListValueMust(types.MapType{ElemType: types.StringType}, nil),
		"context":                  types.StringNull(),
		"description":              types.StringNull(),
		"disabled":                 types.BoolValue(false),
		"enable_output_expression": types.BoolValue(false),
		"filter":                   types.StringValue("true"),
		"group_id":                 types.StringNull(),
		"name":                     types.StringValue("route-1"),
		"output":                   types.StringNull(),
		"output_expression":        types.StringNull(),
		"pipeline":                 types.StringValue("main"),
		"target_context":           types.StringNull(),
		"final":                    types.BoolValue(false),
		"id":                       types.StringValue("route-1"),
	})
	api := RoutesModel{
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: routeTypes}, []attr.Value{apiRoute}),
	}
	state := RoutesModel{
		ID:     types.StringValue("default"),
		Routes: types.ListNull(types.ObjectType{AttrTypes: routeTypes}),
	}

	applyRoutesAPIToState(&api, &state, false, false)

	route := state.Routes.Elements()[0].(types.Object)
	groupID := route.Attributes()["group_id"]
	if groupID.IsNull() || groupID.IsUnknown() {
		t.Fatalf("group_id not defaulted during import: %#v", groupID)
	}
	if got := groupID.(types.String).ValueString(); got != "default" {
		t.Fatalf("group_id = %q, want default", got)
	}
}

func TestApplyRoutesAPIToStateFillsUnknownGroupAndCommentAttrsAfterWrite(t *testing.T) {
	commentTypes := RoutesCommentsAttrTypes()
	groupTypes := RoutesGroupsAttrTypes()

	stateComment := types.ObjectValueMust(commentTypes, map[string]attr.Value{
		"comment":  types.StringValue("my-comment"),
		"group_id": types.StringUnknown(),
		"id":       types.StringUnknown(),
		"index":    types.Int64Unknown(),
	})
	apiComment := types.ObjectValueMust(commentTypes, map[string]attr.Value{
		"comment":  types.StringValue("my-comment"),
		"group_id": types.StringValue("default"),
		"id":       types.StringValue("tf-comment-0"),
		"index":    types.Int64Value(0),
	})

	stateGroup := types.ObjectValueMust(groupTypes, map[string]attr.Value{
		"description": types.StringValue("group A"),
		"index":       types.Int64Unknown(),
		"name":        types.StringValue("firstgroup"),
	})
	apiGroup := types.ObjectValueMust(groupTypes, map[string]attr.Value{
		"description": types.StringValue("group A"),
		"index":       types.Int64Value(1),
		"name":        types.StringValue("firstgroup"),
	})

	api := RoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: commentTypes}, []attr.Value{apiComment}),
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"mygroup": apiGroup,
		}),
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}
	state := RoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: commentTypes}, []attr.Value{stateComment}),
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"mygroup": stateGroup,
		}),
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}

	applyRoutesAPIToState(&api, &state, true, false)

	comment, ok := state.Comments.Elements()[0].(types.Object)
	if !ok {
		t.Fatalf("comment element type = %T", state.Comments.Elements()[0])
	}
	for _, name := range []string{"group_id", "id", "index"} {
		if comment.Attributes()[name].IsUnknown() {
			t.Fatalf("comment %s is still unknown", name)
		}
	}

	groupValue, ok := state.Groups.Elements()["mygroup"]
	if !ok {
		t.Fatalf("group mygroup missing: %#v", state.Groups)
	}
	group, ok := groupValue.(types.Object)
	if !ok {
		t.Fatalf("group element type = %T", groupValue)
	}
	if group.Attributes()["index"].IsUnknown() {
		t.Fatalf("group index is still unknown")
	}
}

func TestApplyRoutesAPIToStateAssignsDeterministicGroupIndexWithoutAPIMatch(t *testing.T) {
	groupTypes := RoutesGroupsAttrTypes()

	state := RoutesModel{
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"mygroup": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("group A"),
				"index":       types.Int64Unknown(),
				"name":        types.StringValue("firstgroup"),
			}),
			"mygroup2": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("group B"),
				"index":       types.Int64Unknown(),
				"name":        types.StringValue("secondgroup"),
			}),
		}),
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}
	api := RoutesModel{
		// API keys can differ or omit index; deterministic fallback should still resolve unknown indexes.
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"firstgroup": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("group A"),
				"index":       types.Int64Null(),
				"name":        types.StringValue("firstgroup"),
			}),
		}),
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}

	applyRoutesAPIToState(&api, &state, true, false)

	for _, key := range []string{"mygroup", "mygroup2"} {
		groupValue, ok := state.Groups.Elements()[key]
		if !ok {
			t.Fatalf("group %s missing from state", key)
		}
		group, ok := groupValue.(types.Object)
		if !ok {
			t.Fatalf("group %s type = %T", key, groupValue)
		}
		if group.Attributes()["index"].IsUnknown() {
			t.Fatalf("group %s index is still unknown", key)
		}
	}
}

func TestRoutesModelUpdateBodySkipsClaimedCommentIndexes(t *testing.T) {
	commentTypes := RoutesCommentsAttrTypes()
	model := RoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: commentTypes}, []attr.Value{
			types.ObjectValueMust(commentTypes, map[string]attr.Value{
				"comment":  types.StringValue("c1"),
				"group_id": types.StringNull(),
				"id":       types.StringNull(),
				"index":    types.Int64Value(1),
			}),
			types.ObjectValueMust(commentTypes, map[string]attr.Value{
				"comment":  types.StringValue("c2"),
				"group_id": types.StringNull(),
				"id":       types.StringNull(),
				"index":    types.Int64Null(),
			}),
		}),
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: RoutesGroupsAttrTypes()}, nil),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}

	body, err := model.updateBody()
	if err != nil {
		t.Fatalf("updateBody returned error: %v", err)
	}

	comments := body["comments"].([]any)
	first := comments[0].(map[string]any)
	second := comments[1].(map[string]any)
	if got := first["index"]; got != int64(1) {
		t.Fatalf("first comment index = %#v, want 1", got)
	}
	if got := second["index"]; got == int64(1) || got == 1 || got == float64(1) {
		t.Fatalf("second comment index collided at %#v", got)
	}
}

func TestRoutesModelUpdateBodySkipsClaimedGroupIndexes(t *testing.T) {
	groupTypes := RoutesGroupsAttrTypes()
	model := RoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: RoutesCommentsAttrTypes()}, nil),
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"a_group": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("A"),
				"index":       types.Int64Value(1),
				"name":        types.StringValue("A"),
			}),
			"b_group": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("B"),
				"index":       types.Int64Null(),
				"name":        types.StringValue("B"),
			}),
		}),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}

	body, err := model.updateBody()
	if err != nil {
		t.Fatalf("updateBody returned error: %v", err)
	}

	groups := body["groups"].(map[string]any)
	groupA := groups["a_group"].(map[string]any)
	groupB := groups["b_group"].(map[string]any)
	if got := groupA["index"]; got != int64(1) {
		t.Fatalf("a_group index = %#v, want 1", got)
	}
	if got := groupB["index"]; got != 0 {
		t.Fatalf("b_group index = %#v, want 0", got)
	}
}

func TestApplyRoutesAPIToStateSkipsClaimedGroupIndexes(t *testing.T) {
	groupTypes := RoutesGroupsAttrTypes()
	api := RoutesModel{
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, nil),
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}
	state := RoutesModel{
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"a_group": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("A"),
				"index":       types.Int64Value(1),
				"name":        types.StringValue("A"),
			}),
			"b_group": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("B"),
				"index":       types.Int64Unknown(),
				"name":        types.StringValue("B"),
			}),
		}),
		ID:     types.StringValue("default"),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: RoutesRoutesAttrTypes()}, nil),
	}

	applyRoutesAPIToState(&api, &state, true, false)
	groupA := state.Groups.Elements()["a_group"].(types.Object)
	groupB := state.Groups.Elements()["b_group"].(types.Object)
	aIndex := groupA.Attributes()["index"].(types.Int64).ValueInt64()
	bIndex := groupB.Attributes()["index"].(types.Int64).ValueInt64()
	if aIndex == bIndex {
		t.Fatalf("group indexes collided at %d", aIndex)
	}
}

func TestNormalizeRouteCommentsAssignsWhenIndexIsNil(t *testing.T) {
	output := map[string]any{
		"comments": []any{
			map[string]any{"comment": "c1", "index": nil},
			map[string]any{"comment": "c2", "index": 1},
		},
	}

	normalizeRouteComments(output)

	comments := output["comments"].([]any)
	first := comments[0].(map[string]any)
	second := comments[1].(map[string]any)
	if first["index"] == nil {
		t.Fatalf("first comment index should be assigned")
	}
	if first["index"] == second["index"] {
		t.Fatalf("first comment index collided with claimed value %#v", second["index"])
	}
}

func TestNormalizeRouteGroupsAssignsWhenIndexIsNil(t *testing.T) {
	output := map[string]any{
		"groups": map[string]any{
			"a_group": map[string]any{"name": "A", "index": nil},
			"b_group": map[string]any{"name": "B", "index": 1},
		},
	}

	normalizeRouteGroups(output)

	groups := output["groups"].(map[string]any)
	groupA := groups["a_group"].(map[string]any)
	groupB := groups["b_group"].(map[string]any)
	if groupA["index"] == nil {
		t.Fatalf("a_group index should be assigned")
	}
	if groupA["index"] == groupB["index"] {
		t.Fatalf("a_group index collided with claimed value %#v", groupB["index"])
	}
}

func TestPackRoutesModelUpdateBodyUsesRoutesNormalization(t *testing.T) {
	commentTypes := PackRoutesCommentsAttrTypes()
	groupTypes := PackRoutesGroupsAttrTypes()
	routeTypes := PackRoutesRoutesAttrTypes()
	model := PackRoutesModel{
		Comments: types.ListValueMust(types.ObjectType{AttrTypes: commentTypes}, []attr.Value{
			types.ObjectValueMust(commentTypes, map[string]attr.Value{
				"comment":  types.StringValue("c1"),
				"group_id": types.StringNull(),
				"id":       types.StringNull(),
				"index":    types.Int64Null(),
			}),
		}),
		Groups: types.MapValueMust(types.ObjectType{AttrTypes: groupTypes}, map[string]attr.Value{
			"z_group": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("Z"),
				"index":       types.Int64Value(1),
				"name":        types.StringValue("Z"),
			}),
			"a_group": types.ObjectValueMust(groupTypes, map[string]attr.Value{
				"description": types.StringValue("A"),
				"index":       types.Int64Null(),
				"name":        types.StringValue("A"),
			}),
		}),
		Routes: types.ListValueMust(types.ObjectType{AttrTypes: routeTypes}, []attr.Value{
			types.ObjectValueMust(routeTypes, map[string]attr.Value{
				"clones":                   types.ListNull(types.MapType{ElemType: types.StringType}),
				"context":                  types.StringNull(),
				"description":              types.StringNull(),
				"disabled":                 types.BoolNull(),
				"enable_output_expression": types.BoolNull(),
				"filter":                   types.StringNull(),
				"final":                    types.BoolNull(),
				"group_id":                 types.StringNull(),
				"id":                       types.StringNull(),
				"name":                     types.StringValue("r1"),
				"output":                   types.StringNull(),
				"output_expression":        types.StringNull(),
				"pipeline":                 types.StringValue("main"),
				"target_context":           types.StringNull(),
			}),
		}),
	}

	body, err := model.updateBody()
	if err != nil {
		t.Fatalf("updateBody returned error: %v", err)
	}

	comments := body["comments"].([]any)
	comment := comments[0].(map[string]any)
	if comment["id"] == nil {
		t.Fatalf("comment id should be assigned")
	}
	if comment["index"] == nil {
		t.Fatalf("comment index should be assigned")
	}

	groups := body["groups"].(map[string]any)
	if groups["a_group"].(map[string]any)["index"] == groups["z_group"].(map[string]any)["index"] {
		t.Fatalf("group indexes should not collide")
	}

	routes := body["routes"].([]any)
	route := routes[0].(map[string]any)
	if route["groupId"] != "default" {
		t.Fatalf("route groupId = %#v, want default", route["groupId"])
	}
}

func TestRoutesListWithDefaultGroupIDSetsDefaultForNullAndUnknown(t *testing.T) {
	routeTypes := RoutesRoutesAttrTypes()
	routes := types.ListValueMust(types.ObjectType{AttrTypes: routeTypes}, []attr.Value{
		types.ObjectValueMust(routeTypes, map[string]attr.Value{
			"clones":                   types.ListNull(types.MapType{ElemType: types.StringType}),
			"context":                  types.StringNull(),
			"description":              types.StringNull(),
			"disabled":                 types.BoolNull(),
			"enable_output_expression": types.BoolNull(),
			"filter":                   types.StringNull(),
			"group_id":                 types.StringNull(),
			"name":                     types.StringValue("r1"),
			"output":                   types.StringNull(),
			"output_expression":        types.StringNull(),
			"pipeline":                 types.StringValue("main"),
			"target_context":           types.StringNull(),
			"final":                    types.BoolNull(),
			"id":                       types.StringNull(),
		}),
		types.ObjectValueMust(routeTypes, map[string]attr.Value{
			"clones":                   types.ListNull(types.MapType{ElemType: types.StringType}),
			"context":                  types.StringNull(),
			"description":              types.StringNull(),
			"disabled":                 types.BoolNull(),
			"enable_output_expression": types.BoolNull(),
			"filter":                   types.StringNull(),
			"group_id":                 types.StringUnknown(),
			"name":                     types.StringValue("r2"),
			"output":                   types.StringNull(),
			"output_expression":        types.StringNull(),
			"pipeline":                 types.StringValue("main"),
			"target_context":           types.StringNull(),
			"final":                    types.BoolNull(),
			"id":                       types.StringNull(),
		}),
	})

	normalized := routesListWithDefaultGroupID(routes)
	for idx, element := range normalized.Elements() {
		route := element.(types.Object)
		groupID := route.Attributes()["group_id"]
		if groupID.IsNull() || groupID.IsUnknown() {
			t.Fatalf("route %d group_id not normalized: %#v", idx, groupID)
		}
		if groupID.(types.String).ValueString() != "default" {
			t.Fatalf("route %d group_id = %q, want default", idx, groupID.(types.String).ValueString())
		}
	}
}
