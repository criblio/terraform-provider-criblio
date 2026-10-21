package registry

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func sensitiveSchemaPaths(resourceSchema schema.Schema) map[string]bool {
	paths := make(map[string]bool)
	collectSensitiveAttributes(resourceSchema.Attributes, "", paths)
	collectSensitiveBlocks(resourceSchema.Blocks, "", paths)
	return paths
}

func collectSensitiveAttributes(attributes map[string]schema.Attribute, prefix string, paths map[string]bool) {
	for name, attribute := range attributes {
		path := prefix + name
		if attribute.IsSensitive() {
			paths[path] = true
		}
		var nested map[string]schema.Attribute
		nestedPrefix := path + "."
		switch attribute := attribute.(type) {
		case schema.SingleNestedAttribute:
			nested = attribute.Attributes
		case schema.ListNestedAttribute:
			nested = attribute.NestedObject.Attributes
		case schema.SetNestedAttribute:
			nested = attribute.NestedObject.Attributes
		case schema.MapNestedAttribute:
			nested = attribute.NestedObject.Attributes
			nestedPrefix += "*."
		}
		collectSensitiveAttributes(nested, nestedPrefix, paths)
	}
}

func collectSensitiveBlocks(blocks map[string]schema.Block, prefix string, paths map[string]bool) {
	for name, block := range blocks {
		path := prefix + name + "."
		switch block := block.(type) {
		case schema.SingleNestedBlock:
			collectSensitiveAttributes(block.Attributes, path, paths)
			collectSensitiveBlocks(block.Blocks, path, paths)
		case schema.ListNestedBlock:
			collectSensitiveAttributes(block.NestedObject.Attributes, path, paths)
			collectSensitiveBlocks(block.NestedObject.Blocks, path, paths)
		case schema.SetNestedBlock:
			collectSensitiveAttributes(block.NestedObject.Attributes, path, paths)
			collectSensitiveBlocks(block.NestedObject.Blocks, path, paths)
		}
	}
}
