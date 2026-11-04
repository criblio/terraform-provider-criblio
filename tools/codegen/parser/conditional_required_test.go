package parser

import (
	"reflect"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestBooleanRequirementsFromReferencedSchema(t *testing.T) {
	var document yaml.Node
	if err := yaml.Unmarshal([]byte(`Metadata:
  type: object
  properties:
    enableAcceleration: {type: boolean, x-terraform-name: accelerated}
    earliest: {type: string, x-terraform-name: start_time}
  if:
    required: [enableAcceleration]
    properties:
      enableAcceleration: {const: true}
  then:
    required: [earliest]
`), &document); err != nil {
		t.Fatal(err)
	}
	schemas := document.Content[0]
	var property yaml.Node
	if err := yaml.Unmarshal([]byte("$ref: '#/components/schemas/Metadata'"), &property); err != nil {
		t.Fatal(err)
	}
	field, err := fieldDef("Dataset", "metadata", property.Content[0], schemas)
	if err != nil {
		t.Fatal(err)
	}
	want := []BooleanRequirement{{Property: "accelerated", Value: true, Required: []string{"start_time"}}}
	if !reflect.DeepEqual(field.BooleanRequirements, want) {
		t.Fatalf("rules = %+v, want %+v", field.BooleanRequirements, want)
	}
}

func TestBooleanRequirementsConditions(t *testing.T) {
	fields := []FieldDef{{APIName: "enabled", TerraformName: "enabled", Type: "boolean"}, {APIName: "start", TerraformName: "start", Type: "string"}}
	for _, test := range []struct {
		name, schema        string
		wantRule, wantError bool
	}{
		{"false", "if: {required: [enabled], properties: {enabled: {const: false}}}\nthen: {required: [start]}", true, false},
		{"missing-trigger-requirement", "if: {properties: {enabled: {const: true}}}\nthen: {required: [start]}", false, false},
		{"non-boolean", "if: {required: [enabled], properties: {enabled: {const: detailed}}}\nthen: {required: [start]}", false, false},
		{"additional-condition", "if: {required: [enabled], properties: {enabled: {const: true}}, minProperties: 1}\nthen: {required: [start]}", false, false},
		{"unknown-required-property", "if: {required: [enabled], properties: {enabled: {const: true}}}\nthen: {required: [missing]}", false, true},
		{"unknown-trigger-property", "if: {required: [missing], properties: {missing: {const: true}}}\nthen: {required: [start]}", false, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var schema yaml.Node
			if err := yaml.Unmarshal([]byte(test.schema), &schema); err != nil {
				t.Fatal(err)
			}
			rules, err := booleanRequirements(schema.Content[0], fields)
			if (err != nil) != test.wantError || (len(rules) > 0) != test.wantRule {
				t.Fatalf("rules = %+v, error = %v", rules, err)
			}
			if test.wantRule && rules[0].Value {
				t.Fatal("false condition became true")
			}
		})
	}
}
