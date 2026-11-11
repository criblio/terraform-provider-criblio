package main

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestWorkspaceSettingsGeneration(t *testing.T) {
	definition := parser.ResourceDef{
		Name: "WorkspaceSystemSettings", StructName: "WorkspaceSystemSettings", TypeName: "criblio_workspace_system_settings",
		ResourceFactoryHook: "newWorkspaceSystemSettingsManagedResource", RefreshFromAPI: true,
		Create: parser.OperationDef{Method: "PATCH", Path: "/system/settings/conf"},
		Read:   parser.OperationDef{Method: "GET", Path: "/system/settings/conf"},
		Update: parser.OperationDef{Method: "PATCH", Path: "/system/settings/conf"},
		Fields: []parser.FieldDef{{APIName: "id", TerraformName: "id", GoName: "ID", Type: "string", FixedValue: "conf", Optional: true, Computed: true}},
	}
	client := renderTemplate(t, "client", definition)
	assertContains(t, client, `a.client, "/system/settings/conf"`)
	assertNotContains(t, client, "/m/")
	assertNotContains(t, client, "restclient.Delete")
	assertContains(t, renderTemplate(t, "resource", definition), "return newWorkspaceSystemSettingsManagedResource()")
	assertContains(t, renderTemplate(t, "data_source", definition), `model.ID = types.StringValue("conf")`)
	assertContains(t, generatedImportBlock(definition), `id = "conf"`)
	assertContains(t, generatedImportCommand(definition), `"conf"`)
	assertNotContains(t, generatedImportBlock(definition), "example\"")
}
