package main

import (
	"path/filepath"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestNativeCollectorSchemaPreservesLegacyContract(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.yml")
	writeFile(t, specPath, `components:
  schemas:
    Legacy:
      type: object
      properties:
        timeout:
          type: integer
          x-terraform-sensitive: true
        legacyOnly:
          type: string
        nested:
          type: object
          properties:
            old:
              type: string
    Upstream:
      type: object
      required: [timeout, added]
      properties:
        timeout:
          type: number
          minimum: 1
        added:
          $ref: '#/components/schemas/Added'
        nested:
          type: object
          required: [new]
          properties:
            new:
              type: boolean
        __scheduling:
          type: object
        defaultBreakers:
          type: string
    Added:
      type: object
      required: [value]
      properties:
        value:
          type: string
`)
	spec, err := readYAML(specPath)
	if err != nil {
		t.Fatal(err)
	}
	root := documentMapping(spec)
	schemas, ok := lookupComponentsSchemas(root)
	if !ok {
		t.Fatal("missing schemas")
	}
	native, _ := mappingValue(schemas, "Upstream")
	legacy, _ := mappingValue(schemas, "Legacy")
	derived, err := compatibleCollectorSchema(native, legacy, schemas, schemas)
	if err != nil {
		t.Fatal(err)
	}
	setMappingValue(schemas, scalar("Legacy"), derived)
	for path, want := range map[string]string{
		"Legacy.properties.timeout.type":                                         "integer",
		"Legacy.properties.timeout.x-terraform-sensitive":                        "true",
		"Legacy.properties.legacyOnly.type":                                      "string",
		"Legacy.properties.nested.properties.old.type":                           "string",
		"Legacy.properties.nested.properties.new.x-terraform-optional-computed":  "true",
		"Legacy.properties.added.properties.value.x-terraform-optional-computed": "true",
	} {
		got, err := lookupTarget(root, "$.components.schemas."+path)
		if err != nil {
			t.Fatal(err)
		}
		if got.Value != want {
			t.Errorf("%s = %q, want %q", path, got.Value, want)
		}
	}
	for _, path := range []string{"Legacy.required", "Legacy.properties.__scheduling", "Legacy.properties.defaultBreakers", "Legacy.properties.added.required", "Legacy.properties.timeout.minimum"} {
		if _, err := lookupTarget(root, "$.components.schemas."+path); err == nil {
			t.Errorf("unexpected %s", path)
		}
	}
	// Expanding a new field must not mutate the shared upstream reference.
	if _, err := lookupTarget(root, "$.components.schemas.Added.required"); err != nil {
		t.Fatal(err)
	}
}

func TestOverlayCopyDoesNotMutateSource(t *testing.T) {
	dir := t.TempDir()
	specPath, overlayPath := filepath.Join(dir, "spec.yml"), filepath.Join(dir, "overlay.yml")
	writeFile(t, specPath, `components:
  schemas:
    Base:
      properties:
        enabled:
          type: boolean
    Database:
      allOf:
        - $ref: '#/components/schemas/Base'
`)
	writeFile(t, overlayPath, `actions:
  - target: $.components.schemas.Database.allOf.0
    copy_from: $.components.schemas.Base
    update:
      properties:
        trackingColumn:
          type: string
`)
	spec, err := readYAML(specPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := applyOverlayFile(spec, overlayPath); err != nil {
		t.Fatal(err)
	}
	root := documentMapping(spec)
	for _, path := range []string{"Database.allOf.0.properties.enabled", "Database.allOf.0.properties.trackingColumn"} {
		if _, err := lookupTarget(root, "$.components.schemas."+path); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := lookupTarget(root, "$.components.schemas.Base.properties.trackingColumn"); err == nil {
		t.Fatal("copy mutated base schema")
	}
}

func TestCollectorSchedulingProjectionIncludesConditionalFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.yml")
	writeFile(t, path, `job:
  properties:
    schedule:
      properties:
        run:
          properties:
            mode: {type: string}
collector:
  properties:
    conf:
      properties:
        __scheduling:
          properties:
            stateTracking:
              properties:
                enabled: {type: boolean}
              dependencies:
                enabled:
                  oneOf:
                    - properties:
                        enabled: {enum: [false]}
                    - required: [trackingColumn]
                      properties:
                        trackingColumn: {type: string}
`)
	spec, err := readYAML(path)
	if err != nil {
		t.Fatal(err)
	}
	root := documentMapping(spec)
	job, _ := mappingValue(root, "job")
	collector, _ := mappingValue(root, "collector")
	projected, ok, err := collectorJobWithScheduling(job, collector, &yaml.Node{Kind: yaml.MappingNode})
	if err != nil || !ok {
		t.Fatalf("projection: ok=%v err=%v", ok, err)
	}
	if _, err := lookupTarget(projected, "$.properties.schedule.properties.run.properties.stateTracking.properties.trackingColumn"); err != nil {
		t.Fatal(err)
	}
	if _, err := lookupTarget(job, "$.properties.schedule.properties.run.properties.stateTracking"); err == nil {
		t.Fatal("projection changed the upstream shared job")
	}
}
