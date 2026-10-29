// Package specschema resolves OpenAPI object compositions for Terraform schemas.
package specschema

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Flatten resolves references and includes properties from every object branch.
// Branch-specific requirements stay optional; unconditional requirements remain.
// It returns private copies and does not modify the upstream document.
func Flatten(schema, schemas *yaml.Node) (*yaml.Node, error) {
	return flatten(schema, schemas, make(map[string]bool))
}

func flatten(schema, schemas *yaml.Node, seen map[string]bool) (*yaml.Node, error) {
	if schema == nil || schema.Kind != yaml.MappingNode {
		return clone(schema), nil
	}
	result := clone(schema)
	if ref, ok := value(result, "$ref"); ok {
		name := strings.TrimPrefix(ref.Value, "#/components/schemas/")
		if name == ref.Value || seen[name] {
			return nil, fmt.Errorf("unsupported or recursive schema reference %q", ref.Value)
		}
		source, ok := value(schemas, name)
		if !ok {
			return nil, fmt.Errorf("schema %q not found", name)
		}
		seen[name] = true
		resolved, err := flatten(source, schemas, seen)
		delete(seen, name)
		if err != nil {
			return nil, err
		}
		remove(result, "$ref")
		merge(resolved, result)
		result = resolved
	}
	if properties, ok := value(result, "properties"); ok {
		for i := 1; i < len(properties.Content); i += 2 {
			resolved, err := flatten(properties.Content[i], schemas, seen)
			if err != nil {
				return nil, err
			}
			properties.Content[i] = resolved
		}
	}
	for _, key := range []string{"items", "additionalProperties"} {
		if child, ok := value(result, key); ok && child.Kind == yaml.MappingNode {
			resolved, err := flatten(child, schemas, seen)
			if err != nil {
				return nil, err
			}
			set(result, key, resolved)
		}
	}
	// Merge allOf alongside explicit properties, and collect all alternatives of
	// inline object unions instead of selecting only their first branch.
	for _, key := range []string{"allOf", "oneOf", "anyOf"} {
		branches, ok := value(result, key)
		if !ok {
			continue
		}
		hasObjects := false
		for _, branch := range branches.Content {
			resolved, err := flatten(branch, schemas, seen)
			if err != nil {
				return nil, err
			}
			if properties, ok := value(resolved, "properties"); ok {
				hasObjects = true
				target, ok := value(result, "properties")
				if !ok {
					target = &yaml.Node{Kind: yaml.MappingNode}
					set(result, "properties", target)
				}
				merge(target, properties)
				if key == "allOf" {
					mergeRequired(result, resolved)
				}
			}
		}
		if hasObjects {
			remove(result, key)
			set(result, "type", scalar("object"))
		}
	}
	// JSON Schema conditionals may also contribute configurable fields.
	for _, key := range []string{"then", "else"} {
		if branch, ok := value(result, key); ok {
			resolved, err := flatten(branch, schemas, seen)
			if err != nil {
				return nil, err
			}
			if properties, ok := value(resolved, "properties"); ok {
				target, ok := value(result, "properties")
				if !ok {
					target = &yaml.Node{Kind: yaml.MappingNode}
					set(result, "properties", target)
				}
				merge(target, properties)
			}
		}
	}
	if dependencies, ok := value(result, "dependencies"); ok {
		for i := 1; i < len(dependencies.Content); i += 2 {
			resolved, err := flatten(dependencies.Content[i], schemas, seen)
			if err != nil {
				return nil, err
			}
			if properties, ok := value(resolved, "properties"); ok {
				target, ok := value(result, "properties")
				if !ok {
					target = &yaml.Node{Kind: yaml.MappingNode}
					set(result, "properties", target)
				}
				merge(target, properties)
			}
		}
	}
	return result, nil
}

func mergeRequired(target, source *yaml.Node) {
	required, ok := value(source, "required")
	if !ok {
		return
	}
	existing, ok := value(target, "required")
	if !ok {
		set(target, "required", clone(required))
		return
	}
	for _, name := range required.Content {
		found := false
		for _, old := range existing.Content {
			if old.Value == name.Value {
				found = true
				break
			}
		}
		if !found {
			existing.Content = append(existing.Content, clone(name))
		}
	}
}

// Preserve explicit definitions while collecting additional nested properties
// and enum values supplied by alternative branches.
func merge(target, source *yaml.Node) {
	for i := 0; i < len(source.Content); i += 2 {
		key, incoming := source.Content[i].Value, source.Content[i+1]
		existing, ok := value(target, key)
		if !ok {
			set(target, key, clone(incoming))
			continue
		}
		if existing.Kind == yaml.MappingNode && incoming.Kind == yaml.MappingNode {
			merge(existing, incoming)
		}
		if key == "enum" && existing.Kind == yaml.SequenceNode && incoming.Kind == yaml.SequenceNode {
			for _, v := range incoming.Content {
				found := false
				for _, old := range existing.Content {
					if old.Value == v.Value {
						found = true
						break
					}
				}
				if !found {
					existing.Content = append(existing.Content, clone(v))
				}
			}
		}
	}
}
func value(node *yaml.Node, key string) (*yaml.Node, bool) {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil, false
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1], true
		}
	}
	return nil, false
}
func set(node *yaml.Node, key string, v *yaml.Node) {
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			node.Content[i+1] = v
			return
		}
	}
	node.Content = append(node.Content, scalar(key), v)
}
func remove(node *yaml.Node, key string) {
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			return
		}
	}
}
func scalar(s string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: s} }
func clone(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	result := *node
	result.Content = make([]*yaml.Node, len(node.Content))
	for i, c := range node.Content {
		result.Content[i] = clone(c)
	}
	return &result
}
