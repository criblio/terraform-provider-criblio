package main

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestDatasetACLGenerationHooksAndIdentity(t *testing.T) {
	def := parser.ResourceDef{
		Name: "SearchDatasetAcl", StructName: "SearchDatasetACL", TypeName: "criblio_search_dataset_acl",
		ResourceFactoryHook: "newSearchDatasetAclManagedResource", DataSourceFactoryHook: "newSearchDatasetAclManagedDataSource", ClientStub: true,
		Fields: []parser.FieldDef{{APIName: "datasetId", TerraformName: "dataset_id", GoName: "DatasetID", Type: "string", Required: true, ForceNew: true, PathParam: true}},
		Read: parser.OperationDef{Method: "GET", Path: "/m/{groupId}/search/datasets/{id}/acl", PathParams: []parser.FieldDef{
			{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID", Type: "string", PathParam: true, Required: true},
			{APIName: "id", TerraformName: "dataset_id", GoName: "DatasetID", Type: "string", PathParam: true, Required: true},
		}},
	}
	assertContains(t, renderTemplate(t, "resource", def), "return newSearchDatasetAclManagedResource()")
	assertContains(t, renderTemplate(t, "data_source", def), "return newSearchDatasetAclManagedDataSource()")
	client := renderTemplate(t, "client", def)
	assertContains(t, client, "the generated client is disabled")
	assertNotContains(t, client, "restclient.Get")
	assertNotContains(t, client, "restclient.Post")
	assertNotContains(t, client, "model.ID.ValueString()")
	assertNotContains(t, generatedImportBlock(def), "jsonencode")
	assertContains(t, renderTemplate(t, "doc", def), "Each configured map is authoritative")
}
