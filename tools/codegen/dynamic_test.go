package main

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestDynamicObjectFieldRendering(t *testing.T) {
	resource := parser.ResourceDef{
		StructName: "Collector", ListStructName: "Collectors",
		OneOfVariants: []parser.OneOfVariantDef{{
			GoName: "InputCollectorDatabase", ModelName: "InputCollectorDatabaseModel",
			TerraformName: "input_collector_database", DiscriminatorValue: "database",
			Fields: []parser.FieldDef{{
				APIName: "savedState", TerraformName: "saved_state", GoName: "SavedState",
				Type: "object", CustomType: "types.DynamicType{}", Optional: true, Computed: true,
			}},
		}},
	}
	assertContains(t, renderTemplate(t, "types", resource), "SavedState types.Dynamic")
	assertContains(t, renderTemplate(t, "types", resource), "return dynamicValueFromAPI(value)")
	assertContains(t, renderTemplate(t, "resource", resource), `"saved_state": schema.DynamicAttribute{`)
	assertContains(t, renderTemplate(t, "data_source", resource), `"saved_state": schema.DynamicAttribute{`)
	assertContains(t, renderTemplate(t, "list_data_source", resource), `"items": schema.DynamicAttribute{Computed: true}`)
	assertContains(t, renderTemplate(t, "doc", resource), "`saved_state` (Dynamic)")
}

func TestCollectorSavedStateModifyPlanRendering(t *testing.T) {
	for _, name := range []string{"Collector", "PackCollector"} {
		t.Run(name, func(t *testing.T) {
			output := renderTemplate(t, "resource", parser.ResourceDef{StructName: name})
			assertContains(t, output, "var _ resource.ResourceWithModifyPlan = &"+name+"Resource{}")
			assertContains(t, output, "modifyCollectorSavedStatePlan(ctx, req, resp)")
		})
	}
}
