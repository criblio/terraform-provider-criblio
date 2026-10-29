package main

import (
	"fmt"
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/specschema"
	"go.yaml.in/yaml/v3"
)

// applyCollectorSchemas derives collector wrappers from the upstream saved-job
// and concrete collector schemas. Legacy definitions supply compatibility only.
func applyCollectorSchemas(spec *yaml.Node, overlayPath string) error {
	overlay, err := readYAML(overlayPath)
	if err != nil {
		return err
	}
	config, ok := mappingValue(documentMapping(overlay), "collector_schemas")
	if !ok {
		return nil
	}
	job, ok := mappingValue(config, "job")
	if !ok || job.Kind != yaml.ScalarNode {
		return fmt.Errorf("collector_schemas: job scalar is required")
	}
	base, ok := mappingValue(config, "base")
	if !ok || base.Kind != yaml.ScalarNode {
		return fmt.Errorf("collector_schemas: base scalar is required")
	}
	variants, ok := mappingValue(config, "variants")
	if !ok || variants.Kind != yaml.MappingNode {
		return fmt.Errorf("collector_schemas: variants mapping is required")
	}
	schemas, ok := lookupComponentsSchemas(documentMapping(spec))
	if !ok {
		return fmt.Errorf("collector components schemas not found")
	}
	// Capture compatibility definitions before replacing their shared base.
	prior := cloneNode(schemas)
	nativeJob, ok := mappingValue(schemas, job.Value)
	if !ok {
		return fmt.Errorf("collector job schema %q not found", job.Value)
	}
	legacyBase, ok := mappingValue(prior, base.Value)
	if !ok {
		return fmt.Errorf("collector compatibility base %q not found", base.Value)
	}
	nativeBase, err := compatibleCollectorSchema(nativeJob, legacyBase, schemas, prior)
	if err != nil {
		return fmt.Errorf("derive collector base: %w", err)
	}
	properties, _ := mappingValue(nativeBase, "properties")
	deleteMappingKey(properties, "collector")
	if err := applyCollectorRuntimeSchema(nativeBase, schemas, config); err != nil {
		return err
	}
	setMappingValue(schemas, scalar(base.Value), nativeBase)
	for i := 0; i < len(variants.Content); i += 2 {
		legacyName, nativeName := variants.Content[i].Value, variants.Content[i+1].Value
		legacy, ok := mappingValue(prior, legacyName)
		if !ok {
			return fmt.Errorf("collector compatibility schema %q not found", legacyName)
		}
		native, ok := mappingValue(schemas, nativeName)
		if !ok {
			return fmt.Errorf("upstream collector schema %q not found", nativeName)
		}
		wrapper := cloneNode(legacy)
		collector, err := lookupTarget(wrapper, "$.allOf.1.properties.collector")
		if err != nil {
			return fmt.Errorf("collector wrapper %q: %w", legacyName, err)
		}
		derived, err := compatibleCollectorSchema(native, collector, schemas, prior)
		if err != nil {
			return fmt.Errorf("derive %s: %w", legacyName, err)
		}
		*collector = *derived
		allOf, _ := mappingValue(wrapper, "allOf")
		if len(allOf.Content) > 0 {
			jobWithScheduling, hasScheduling, err := collectorJobWithScheduling(nativeJob, native, schemas)
			if err != nil {
				return fmt.Errorf("derive %s scheduling: %w", legacyName, err)
			}
			_, hasRef := mappingValue(allOf.Content[0], "$ref")
			if hasScheduling || !hasRef {
				// Project collector-specific scheduling into the saved-job API layout
				// and retain compatibility extensions on a private base.
				privateBase, err := compatibleCollectorSchema(jobWithScheduling, allOf.Content[0], schemas, prior)
				if err != nil {
					return fmt.Errorf("derive %s scheduling: %w", legacyName, err)
				}
				properties, _ := mappingValue(privateBase, "properties")
				deleteMappingKey(properties, "collector")
				if err := applyCollectorRuntimeSchema(privateBase, schemas, config); err != nil {
					return err
				}
				allOf.Content[0] = privateBase
			}
		}
		setMappingValue(schemas, scalar(legacyName), wrapper)
	}
	return nil
}

// Keep the upstream runtime shape while retaining the Terraform compatibility annotations.
func applyCollectorRuntimeSchema(base, schemas, config *yaml.Node) error {
	runtime, ok := mappingValue(config, "runtime")
	if !ok {
		return nil
	}
	schema, ok := mappingValue(schemas, runtime.Value)
	if !ok {
		return fmt.Errorf("collector runtime schema %q not found", runtime.Value)
	}
	upstream, err := lookupTarget(schema, "$.properties.savedState")
	if err != nil {
		return fmt.Errorf("collector runtime savedState: %w", err)
	}
	properties, _ := mappingValue(base, "properties")
	state := cloneNode(upstream)
	if legacy, ok := mappingValue(properties, "savedState"); ok {
		for i := 0; i < len(legacy.Content); i += 2 {
			if strings.HasPrefix(legacy.Content[i].Value, "x-terraform-") {
				setMappingValue(state, cloneNode(legacy.Content[i]), cloneNode(legacy.Content[i+1]))
			}
		}
	}
	setMappingValue(properties, scalar("savedState"), state)
	return nil
}

func collectorJobWithScheduling(job, collector, schemas *yaml.Node) (*yaml.Node, bool, error) {
	expanded, err := specschema.Flatten(collector, schemas)
	if err != nil {
		return nil, false, err
	}
	scheduling, err := lookupTarget(expanded, "$.properties.conf.properties.__scheduling")
	if err != nil {
		return job, false, nil
	}
	projected, err := specschema.Flatten(job, schemas)
	if err != nil {
		return nil, false, err
	}
	run, err := lookupTarget(projected, "$.properties.schedule.properties.run")
	if err != nil {
		return nil, false, fmt.Errorf("upstream job scheduling target: %w", err)
	}
	mergeMapping(run, scheduling)
	return projected, true, nil
}

func compatibleCollectorSchema(native, legacy, nativeSchemas, legacySchemas *yaml.Node) (*yaml.Node, error) {
	derived, err := specschema.Flatten(native, nativeSchemas)
	if err != nil {
		return nil, err
	}
	prior, err := specschema.Flatten(legacy, legacySchemas)
	if err != nil {
		return nil, err
	}
	prepareNativeCollectorSchema(derived)
	applyCollectorCompatibility(derived, prior)
	return derived, nil
}

func prepareNativeCollectorSchema(schema *yaml.Node) {
	deleteMappingKey(schema, "required")
	setMappingValue(schema, scalar("x-terraform-optional-computed"), boolScalar(true))
	if properties, ok := mappingValue(schema, "properties"); ok {
		for i := 0; i < len(properties.Content); {
			name := properties.Content[i].Value
			if strings.HasPrefix(name, "__") || name == "defaultBreakers" {
				properties.Content = append(properties.Content[:i], properties.Content[i+2:]...)
				continue
			}
			prepareNativeCollectorSchema(properties.Content[i+1])
			i += 2
		}
	}
	if items, ok := mappingValue(schema, "items"); ok && items.Kind == yaml.MappingNode {
		prepareNativeCollectorSchema(items)
	}
	// These conditions have already contributed their configurable properties.
	for _, key := range []string{"dependencies", "if", "then", "else"} {
		deleteMappingKey(schema, key)
	}
}

// Keep new upstream properties, but restore types, validation, field names, and
// annotations for every previously accepted attribute. Legacy-only fields remain.
func applyCollectorCompatibility(native, legacy *yaml.Node) {
	for _, key := range []string{"type", "enum", "const", "pattern", "minLength", "maxLength", "minimum", "maximum", "minItems", "maxItems", "uniqueItems", "readOnly", "writeOnly", "required", "additionalProperties"} {
		deleteMappingKey(native, key)
	}
	for i := 0; i < len(native.Content); {
		if strings.HasPrefix(native.Content[i].Value, "x-terraform-") || strings.HasPrefix(native.Content[i].Value, "x-speakeasy-") {
			native.Content = append(native.Content[:i], native.Content[i+2:]...)
			continue
		}
		i += 2
	}
	for i := 0; i < len(legacy.Content); i += 2 {
		key, value := legacy.Content[i], legacy.Content[i+1]
		if key.Value == "properties" {
			properties, ok := mappingValue(native, "properties")
			if !ok {
				setMappingValue(native, cloneNode(key), cloneNode(value))
				continue
			}
			for j := 0; j < len(value.Content); j += 2 {
				name, old := value.Content[j], value.Content[j+1]
				if current, ok := mappingValue(properties, name.Value); ok {
					applyCollectorCompatibility(current, old)
				} else {
					setMappingValue(properties, cloneNode(name), cloneNode(old))
				}
			}
			// Preserve declaration order as well as behavior to keep regeneration
			// diffs focused on the newly available upstream fields.
			ordered := make([]*yaml.Node, 0, len(properties.Content))
			seen := make(map[string]bool)
			for j := 0; j < len(value.Content); j += 2 {
				name := value.Content[j]
				current, _ := mappingValue(properties, name.Value)
				ordered = append(ordered, cloneNode(name), current)
				seen[name.Value] = true
			}
			for j := 0; j < len(properties.Content); j += 2 {
				if !seen[properties.Content[j].Value] {
					ordered = append(ordered, properties.Content[j], properties.Content[j+1])
				}
			}
			properties.Content = ordered
			continue
		}
		if key.Value == "items" {
			if current, ok := mappingValue(native, "items"); ok && current.Kind == yaml.MappingNode && value.Kind == yaml.MappingNode {
				applyCollectorCompatibility(current, value)
				continue
			}
		}
		setMappingValue(native, cloneNode(key), cloneNode(value))
	}
}
