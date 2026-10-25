package specschema

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestFlattenIncludesEveryConditionalObjectBranch(t *testing.T) {
	var document yaml.Node
	if err := yaml.Unmarshal([]byte(`schemas:
  Common:
    properties:
      common: {type: string}
    allOf:
      - oneOf:
          - $ref: '#/components/schemas/Static'
          - $ref: '#/components/schemas/Backoff'
    dependencies:
      enabled:
        oneOf:
          - properties:
              enabled: {enum: [false]}
          - required: [trackingColumn]
            properties:
              trackingColumn: {type: string}
  Static:
    required: [interval]
    properties:
      interval: {type: number}
      type: {type: string, enum: [static]}
  Backoff:
    required: [multiplier]
    properties:
      multiplier: {type: number}
      type: {type: string, enum: [backoff]}
`), &document); err != nil {
		t.Fatal(err)
	}
	schemas, _ := value(document.Content[0], "schemas")
	source, _ := value(schemas, "Common")
	result, err := Flatten(source, schemas)
	if err != nil {
		t.Fatal(err)
	}
	properties, _ := value(result, "properties")
	for _, name := range []string{"common", "interval", "multiplier", "trackingColumn"} {
		if _, ok := value(properties, name); !ok {
			t.Errorf("missing conditional field %s", name)
		}
	}
	if _, ok := value(result, "required"); ok {
		t.Fatal("branch-specific requirements became unconditional")
	}
	typ, _ := value(properties, "type")
	enums, _ := value(typ, "enum")
	if len(enums.Content) != 2 {
		t.Fatalf("enum values = %v", enums.Content)
	}
	original, _ := value(source, "properties")
	if _, ok := value(original, "interval"); ok {
		t.Fatal("normalization mutated the upstream schema")
	}
}

func TestFlattenRetainsUnconditionalRequirements(t *testing.T) {
	var document yaml.Node
	if err := yaml.Unmarshal([]byte(`properties:
  id: {type: string}
required: [id]
allOf:
  - properties:
      name: {type: string}
    required: [name]
`), &document); err != nil {
		t.Fatal(err)
	}
	result, err := Flatten(document.Content[0], &yaml.Node{Kind: yaml.MappingNode})
	if err != nil {
		t.Fatal(err)
	}
	required, _ := value(result, "required")
	if len(required.Content) != 2 {
		t.Fatalf("required = %#v", required)
	}
}
