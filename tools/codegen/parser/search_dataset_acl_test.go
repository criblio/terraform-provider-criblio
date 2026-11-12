package parser

import "testing"

func TestDatasetACLProjectionAndFactoryHooks(t *testing.T) {
	defs, err := Parse([]byte(`openapi: 3.1.0
info: {title: test, version: '1'}
paths:
  /m/{groupId}/search/datasets/{id}/acl/apply:
    post:
      x-terraform-resource: true
      x-terraform-resource-name: SearchDatasetAcl
      x-terraform-resource-factory-hook: newSearchDatasetAclManagedResource
      x-terraform-data-source-factory-hook: newSearchDatasetAclManagedDataSource
      parameters:
        - {name: groupId, in: path, required: true, schema: {type: string}}
        - {name: id, in: path, required: true, x-terraform-name: dataset_id, schema: {type: string}}
      requestBody:
        content:
          application/json:
            schema: {$ref: '#/components/schemas/ACL'}
components:
  schemas:
    ACL:
      type: object
      required: [datasetId]
      properties:
        datasetId: {type: string, x-terraform-force-new: true, x-terraform-local-only: true}
        users: {type: object, additionalProperties: {type: string}}
        teams: {type: object, additionalProperties: {type: string}}
`))
	if err != nil {
		t.Fatal(err)
	}
	acl := resourceByName(t, defs, "SearchDatasetAcl")
	if acl.ResourceFactoryHook != "newSearchDatasetAclManagedResource" || acl.DataSourceFactoryHook != "newSearchDatasetAclManagedDataSource" {
		t.Fatal("lost ACL projection factory hooks")
	}
	if len(acl.Fields) != 3 {
		t.Fatalf("unexpected fields (group/id must not leak into schema): %+v", acl.Fields)
	}
	id := fieldByTFName(t, acl.Fields, "dataset_id")
	if !id.Required || !id.ForceNew || !id.PathParam || id.RequestField || id.GoName != "DatasetID" {
		t.Fatalf("incorrect dataset identity: %+v", id)
	}
	param := acl.Create.PathParams[1]
	if param.APIName != "id" || param.GoName != "DatasetID" || param.TerraformName != "dataset_id" {
		t.Fatalf("renamed path parameter lost API placeholder: %+v", param)
	}
	for _, name := range []string{"users", "teams"} {
		field := fieldByTFName(t, acl.Fields, name)
		if field.Type != "object" || len(field.Fields) != 0 || !field.Optional || field.Computed {
			t.Fatalf("invalid declarative map: %+v", field)
		}
	}
}
