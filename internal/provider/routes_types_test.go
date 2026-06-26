package provider

import (
	"encoding/json"
	"testing"

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
