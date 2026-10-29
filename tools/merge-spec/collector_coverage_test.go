package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
	"go.yaml.in/yaml/v3"
)

// Compare raw upstream branches against the fields that codegen actually emits.
// This walker intentionally does not use the normalization implementation.
func TestCollectorUpstreamFieldCoverage(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(filepath.Join("..", ".."))
	output := filepath.Join(dir, "merged.yml")
	if err := run("upstream-openapi.yml", "terraform-overlay.yml", output, filepath.Join(dir, "availability.go")); err != nil {
		t.Fatal(err)
	}
	resources, err := parser.ParseFile(output)
	if err != nil {
		t.Fatal(err)
	}
	source, err := readYAML("upstream-openapi.yml")
	if err != nil {
		t.Fatal(err)
	}
	schemas, ok := lookupComponentsSchemas(documentMapping(source))
	if !ok {
		t.Fatal("upstream schemas missing")
	}
	overlay, err := readYAML("terraform-overlay.yml")
	if err != nil {
		t.Fatal(err)
	}
	config, _ := mappingValue(documentMapping(overlay), "collector_schemas")
	variants, _ := mappingValue(config, "variants")
	mapped := make(map[string]bool)
	for i := 1; i < len(variants.Content); i += 2 {
		mapped[variants.Content[i].Value] = true
	}
	union, _ := mappingValue(schemas, "Collector")
	branches, _ := mappingValue(union, "oneOf")
	for _, branch := range branches.Content {
		ref, _ := mappingValue(branch, "$ref")
		name := strings.TrimPrefix(ref.Value, "#/components/schemas/")
		if !mapped[name] {
			t.Errorf("upstream collector type %s has no Terraform wrapper", name)
		}
	}
	var collector parser.ResourceDef
	for _, resource := range resources {
		if resource.StructName == "Collector" {
			collector = resource
			break
		}
	}
	if len(collector.OneOfVariants) == 0 {
		t.Fatal("collector variants missing")
	}
	merged, err := readYAML(output)
	if err != nil {
		t.Fatal(err)
	}
	mergedSchemas, _ := lookupComponentsSchemas(documentMapping(merged))
	for _, variant := range collector.OneOfVariants {
		for _, field := range variant.Fields {
			if field.APIName == "savedState" && (field.CustomType != "types.DynamicType{}" || !field.Optional || !field.Computed) {
				t.Errorf("%s savedState lost dynamic compatibility: %+v", variant.SchemaName, field)
			}
		}
	}
	base, _ := mappingValue(mergedSchemas, "InputCollectorBase")
	checkpoint, err := lookupTarget(base, "$.properties.savedState.additionalProperties.$ref")
	if err != nil || checkpoint.Value != "#/components/schemas/AdditionalPropertiesTypeEnrichedFieldsSavedState" {
		t.Fatalf("merged savedState lost its upstream checkpoint shape: %v", err)
	}
	jobName, _ := mappingValue(config, "job")
	job, _ := mappingValue(schemas, jobName.Value)
	jobPaths := make(map[string]bool)
	upstreamPropertyPaths(t, job, schemas, "", jobPaths, map[string]bool{})
	for i := 0; i < len(variants.Content); i += 2 {
		legacyName, nativeName := variants.Content[i].Value, variants.Content[i+1].Value
		t.Run(legacyName, func(t *testing.T) {
			var fields []parser.FieldDef
			for _, variant := range collector.OneOfVariants {
				if variant.SchemaName == legacyName {
					fields = variant.Fields
					break
				}
			}
			if len(fields) == 0 {
				t.Fatal("generated variant missing")
			}
			actual := make(map[string]bool)
			generatedPropertyPaths(fields, "", actual)
			native, _ := mappingValue(schemas, nativeName)
			expected := make(map[string]bool)
			upstreamPropertyPaths(t, native, schemas, "collector", expected, map[string]bool{})
			for path := range jobPaths {
				if path != "type" && path != "collector" && !strings.HasPrefix(path, "collector.") {
					expected[path] = true
				}
			}
			for path := range expected {
				if !actual[path] {
					t.Errorf("upstream field %s is absent from generated %s", path, legacyName)
				}
			}
		})
	}
}

func generatedPropertyPaths(fields []parser.FieldDef, prefix string, paths map[string]bool) {
	for _, field := range fields {
		path := field.APIName
		if prefix != "" {
			path = prefix + "." + path
		}
		paths[path] = true
		generatedPropertyPaths(field.Fields, path, paths)
	}
}

func upstreamPropertyPaths(t *testing.T, node, schemas *yaml.Node, prefix string, paths map[string]bool, seen map[string]bool) {
	t.Helper()
	if node == nil || node.Kind != yaml.MappingNode {
		return
	}
	if ref, ok := mappingValue(node, "$ref"); ok {
		name := strings.TrimPrefix(ref.Value, "#/components/schemas/")
		if !seen[name] {
			source, ok := mappingValue(schemas, name)
			if !ok {
				t.Fatalf("missing upstream schema %s", name)
			}
			seen[name] = true
			upstreamPropertyPaths(t, source, schemas, prefix, paths, seen)
			delete(seen, name)
		}
	}
	if properties, ok := mappingValue(node, "properties"); ok {
		for i := 0; i < len(properties.Content); i += 2 {
			name := properties.Content[i].Value
			if name == "__scheduling" {
				upstreamPropertyPaths(t, properties.Content[i+1], schemas, "schedule.run", paths, seen)
				continue
			}
			if strings.HasPrefix(name, "__") || name == "defaultBreakers" {
				continue
			}
			path := name
			if prefix != "" {
				path = prefix + "." + name
			}
			paths[path] = true
			upstreamPropertyPaths(t, properties.Content[i+1], schemas, path, paths, seen)
		}
	}
	if items, ok := mappingValue(node, "items"); ok {
		upstreamPropertyPaths(t, items, schemas, prefix, paths, seen)
	}
	for _, key := range []string{"allOf", "oneOf", "anyOf"} {
		if branches, ok := mappingValue(node, key); ok {
			for _, branch := range branches.Content {
				upstreamPropertyPaths(t, branch, schemas, prefix, paths, seen)
			}
		}
	}
	for _, key := range []string{"then", "else"} {
		if branch, ok := mappingValue(node, key); ok {
			upstreamPropertyPaths(t, branch, schemas, prefix, paths, seen)
		}
	}
	if dependencies, ok := mappingValue(node, "dependencies"); ok {
		for i := 1; i < len(dependencies.Content); i += 2 {
			upstreamPropertyPaths(t, dependencies.Content[i], schemas, prefix, paths, seen)
		}
	}
}
