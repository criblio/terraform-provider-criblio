package main

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestRefreshFromAPIPreservesPathIdentity(t *testing.T) {
	definition := parser.ResourceDef{
		StructName: "Limits", TypeName: "criblio_limits", RefreshFromAPI: true,
		Read: parser.OperationDef{Path: "/m/{groupId}/system/limits", Method: "GET"},
		Fields: []parser.FieldDef{
			{TerraformName: "group_id", GoName: "GroupID", Type: "string", Required: true, PathParam: true},
			{TerraformName: "min_free_space", GoName: "MinFreeSpace", Type: "string", Required: true},
		},
	}
	content := renderTemplate(t, "resource", definition)
	assertContains(t, content, "apiModel.GroupID = model.GroupID")
	assertContains(t, content, "model = *apiModel")
}
