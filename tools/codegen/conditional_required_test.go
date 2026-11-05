package main

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestConditionalRequiredObjectGeneration(t *testing.T) {
	definition := parser.ResourceDef{
		StructName: "Dataset", TypeName: "criblio_dataset",
		Fields: []parser.FieldDef{{
			TerraformName: "metadata", GoName: "Metadata", Type: "object", Optional: true,
			Fields:              []parser.FieldDef{{TerraformName: "enabled", Type: "boolean", Optional: true}, {TerraformName: "start", Type: "string", Optional: true}},
			BooleanRequirements: []parser.BooleanRequirement{{Property: "enabled", Value: true, Required: []string{"start"}}},
		}},
	}
	content := renderTemplate(t, "resource", definition)
	assertContains(t, content, `custom_objectvalidators "github.com/criblio/terraform-provider-criblio/internal/validators/objectvalidators"`)
	assertContains(t, content, `"github.com/hashicorp/terraform-plugin-framework/schema/validator"`)
	assertContains(t, content, `Validators: []validator.Object{`)
	assertContains(t, content, `custom_objectvalidators.RequiredWhenBoolean("enabled", true, "start")`)
	// Data source results are API-owned and must not acquire config validators.
	assertNotContains(t, renderTemplate(t, "data_source", definition), "RequiredWhenBoolean")
	definition.Fields[0].BooleanRequirements = nil
	assertNotContains(t, renderTemplate(t, "resource", definition), "custom_objectvalidators")
}
