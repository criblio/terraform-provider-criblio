package main

import (
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestPackListGeneration(t *testing.T) {
	definition := parser.ResourceDef{
		StructName: "PackPipeline", ListStructName: "PackPipelines",
		List: parser.OperationDef{Method: "GET", Path: "/m/{groupId}/p/{pack}/pipelines", Paginated: true, PathParams: []parser.FieldDef{
			{APIName: "groupId", TerraformName: "group_id", GoName: "GroupID", Type: "string", PathParam: true},
			{APIName: "pack", TerraformName: "pack", GoName: "Pack", Type: "string", PathParam: true},
		}},
	}
	content := renderTemplate(t, "list_data_source", definition)
	assertContains(t, content, "restclient.List[PackPipelineModel]")
	assertContains(t, content, "resolvePackIDForRestAPI(ctx, d.client")
	assertContains(t, content, "model.Pack.ValueString())), true)")
	assertNotContains(t, content, "a.client")
}

func TestDynamicListDocsMatchSchema(t *testing.T) {
	definition := parser.ResourceDef{Fields: []parser.FieldDef{{TerraformName: "saved_state", Type: "object", CustomType: "types.DynamicType{}"}}}
	content := dataSourceDocSchema(definition, true)
	assertContains(t, content, "`items` (Dynamic)")
	assertNotContains(t, content, "Nested Schema")
}

func TestListPaginationDetection(t *testing.T) {
	for _, tc := range []struct {
		parameters            string
		totalCount, paginated bool
	}{
		{totalCount: true, paginated: true},
		{parameters: `parameters:
        - name: offset
          in: query
          required: false
          schema:
            type: integer`, paginated: true},
		{},
	} {
		spec := `openapi: 3.0.3
info:
  title: List test
  version: "1"
paths:
  /things:
    get:
      x-terraform-list: Thing
      ` + tc.parameters + `
      responses:
        '200':
          description: Things
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Things'
components:
  schemas:
    Thing:
      type: object
      properties:
        id:
          type: string
    Things:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/Thing'
        totalCount:
          type: integer
`
		if !tc.totalCount {
			spec = strings.Replace(spec, "        totalCount:\n          type: integer\n", "", 1)
		}
		definitions, err := parser.Parse([]byte(spec))
		if err != nil {
			t.Fatal(err)
		}
		if len(definitions) != 1 || definitions[0].List.Paginated != tc.paginated {
			t.Fatalf("pagination detection: got %v, want %v", definitions, tc.paginated)
		}
	}
}
