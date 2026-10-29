package parser

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestSchemaAnnotationsDriveSensitivity(t *testing.T) {
	var document yaml.Node
	if err := yaml.Unmarshal([]byte(`
type: object
properties:
  authTokens:
    type: array
    items:
      type: string
    x-terraform-sensitive: true
  conf:
    type: object
    properties:
      opaque:
        type: string
        x-terraform-sensitive: true
      tokenSecret:
        type: string
  password:
    type: string
`), &document); err != nil {
		t.Fatal(err)
	}
	fields, _, err := parseSchemaFields("Example", document.Content[0], &yaml.Node{Kind: yaml.MappingNode}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, resourceName := range []string{"Source", "PackSource", "Collector", "PackCollector", "Destination", "PackDestination", "NotificationTarget"} {
		t.Run(resourceName, func(t *testing.T) {
			resource := ResourceDef{StructName: resourceName, OneOfVariants: []OneOfVariantDef{{Fields: fields}}}
			applyResourceCompatibility(&resource)
			got := resource.OneOfVariants[0].Fields
			if !got[0].Sensitive || !got[1].Fields[0].Sensitive {
				t.Fatal("sensitive annotations must survive resource compatibility handling")
			}
			if got[1].Fields[1].Sensitive || got[2].Sensitive {
				t.Fatal("names alone must not make an unannotated field sensitive")
			}
		})
	}
}

func TestSensitiveAnnotationDoesNotChangeSourceStateMerging(t *testing.T) {
	var property yaml.Node
	if err := yaml.Unmarshal([]byte("type: string\nx-terraform-sensitive: true\n"), &property); err != nil {
		t.Fatal(err)
	}
	field := FieldDef{TerraformName: "token", Type: "string"}
	applyFieldAnnotations(&field, property.Content[0], false, true, true, true)
	resource := ResourceDef{StructName: "Source", OneOfVariants: []OneOfVariantDef{{Fields: []FieldDef{field}}}}
	applyResourceCompatibility(&resource)
	got := resource.OneOfVariants[0].Fields[0]
	if !got.Sensitive || got.PreferState {
		t.Fatal("sensitivity alone must not change source state merging")
	}
}
