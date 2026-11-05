package main

import (
	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
	"testing"
)

func TestResourceFactoryHook(t *testing.T) {
	defs, err := parser.Parse([]byte(`openapi: 3.1.0
info: {title: Test, version: '1'}
paths:
  /limits:
    patch:
      x-terraform-resource: true
      x-terraform-resource-name: Limits
      x-terraform-resource-factory-hook: newManagedLimitsResource
      requestBody:
        content:
          application/json:
            schema: {$ref: '#/components/schemas/Limits'}
components:
  schemas:
    Limits:
      type: object
      properties:
        lookupAllowedPaths: {type: array, items: {type: string}}
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) != 1 || defs[0].ResourceFactoryHook != "newManagedLimitsResource" {
		t.Fatalf("factory annotation not parsed: %+v", defs)
	}
	content := renderTemplate(t, "resource", defs[0])
	assertContains(t, content, "return newManagedLimitsResource()")
}
