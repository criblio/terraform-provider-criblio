package converter

// This file shows an example implementation of a resource converter
// for converting Cribl source YAML to Terraform HCL

import (
	"fmt"
	"strings"
	
	"github.com/hashicorp/hcl/v2/hclwrite"
	"gopkg.in/yaml.v3"
)

// Example: Convert a source YAML configuration to Terraform HCL
// This demonstrates the pattern for all resource converters

type SourceConverter struct {
	fieldMapper *FieldMapper
}

func NewSourceConverter() *SourceConverter {
	return &SourceConverter{
		fieldMapper: NewFieldMapper(),
	}
}

// ConvertSource converts a source YAML config to Terraform HCL
func (c *SourceConverter) ConvertSource(sourceID string, yamlConfig map[string]interface{}) (*hclwrite.File, error) {
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	// Create resource block
	resourceBlock := body.AppendNewBlock("resource", []string{
		"criblio_source",
		normalizeResourceName(sourceID),
	})
	resourceBody := resourceBlock.Body()

	// Extract source type - this determines which input_* attribute to use
	sourceType, ok := yamlConfig["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'type' field in source config")
	}

	// Set required fields
	resourceBody.SetAttributeValue("id", hclwrite.TokensForValue(
		hclwrite.TokensForString(sourceID),
	))

	// Map group_id if present
	if groupID, ok := yamlConfig["group"].(string); ok {
		resourceBody.SetAttributeValue("group_id", hclwrite.TokensForValue(
			hclwrite.TokensForString(groupID),
		))
	} else {
		// Default group
		resourceBody.SetAttributeValue("group_id", hclwrite.TokensForValue(
			hclwrite.TokensForString("default"),
		))
	}

	// Handle type-specific input attributes
	// Remove 'type' from config before processing
	inputConfig := make(map[string]interface{})
	for k, v := range yamlConfig {
		if k != "type" && k != "id" && k != "group" {
			inputConfig[k] = v
		}
	}

	// Map source type to Terraform input attribute name
	inputAttrName := mapSourceTypeToInputAttr(sourceType)
	
	// Convert nested config to Terraform attribute
	if len(inputConfig) > 0 {
		inputBlock := resourceBody.AppendNewBlock(inputAttrName, nil)
		inputBody := inputBlock.Body()

		// Map each field using field mapper (camelCase -> snake_case)
		for key, value := range inputConfig {
			terraformKey := c.fieldMapper.MapField(key)
			
			// Handle nested objects (e.g., tls, pq)
			// Note: For SingleNestedAttribute, use = syntax: attr = { ... }
			// For blocks, use block syntax: block { ... }
			if nestedMap, ok := value.(map[string]interface{}); ok {
				// Check if this is a SingleNestedAttribute (like cloud) or a block
				// For now, treat as block - adjust based on schema
				nestedBlock := inputBody.AppendNewBlock(terraformKey, nil)
				nestedBody := nestedBlock.Body()
				
				for nk, nv := range nestedMap {
					terraformNestedKey := c.fieldMapper.MapField(nk)
					setAttributeValue(nestedBody, terraformNestedKey, nv)
				}
			} else {
				setAttributeValue(inputBody, terraformKey, value)
			}
		}
	}

	return file, nil
}

// mapSourceTypeToInputAttr maps Cribl source type to Terraform input attribute
func mapSourceTypeToInputAttr(sourceType string) string {
	typeMap := map[string]string{
		"http":              "input_http",
		"tcp":               "input_tcp",
		"syslog":            "input_syslog",
		"cribl_http":        "input_cribl_http",
		"cribl_tcp":         "input_cribl_tcp",
		"open_telemetry":    "input_open_telemetry",
		"kafka":             "input_kafka",
		"s3":                "input_s3",
		// ... add more mappings
	}

	if attr, ok := typeMap[sourceType]; ok {
		return attr
	}

	// Default: convert type name to snake_case
	return "input_" + toSnakeCase(sourceType)
}

// normalizeResourceName converts resource ID to valid Terraform identifier
func normalizeResourceName(name string) string {
	// Replace invalid characters with underscores
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")
	
	// Ensure it starts with a letter
	if len(name) > 0 && !isLetter(rune(name[0])) {
		name = "r_" + name
	}
	
	return name
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// setAttributeValue sets an attribute value in HCL, handling different types
func setAttributeValue(body *hclwrite.Body, key string, value interface{}) {
	switch v := value.(type) {
	case string:
		body.SetAttributeValue(key, hclwrite.TokensForValue(
			hclwrite.TokensForString(v),
		))
	case bool:
		body.SetAttributeValue(key, hclwrite.TokensForBool(v))
	case int, int32, int64:
		body.SetAttributeValue(key, hclwrite.TokensForValue(
			hclwrite.TokensForNumber(fmt.Sprintf("%d", v)),
		))
	case float64, float32:
		body.SetAttributeValue(key, hclwrite.TokensForValue(
			hclwrite.TokensForNumber(fmt.Sprintf("%g", v)),
		))
	case []interface{}:
		tokens := hclwrite.TokensForValue(hclwrite.TokensForList([]interface{}{}))
		for _, item := range v {
			// Add item to list
		}
		body.SetAttributeRaw(key, tokens)
	default:
		// For complex types, convert to JSON string or handle specially
		body.SetAttributeValue(key, hclwrite.TokensForValue(
			hclwrite.TokensForString(fmt.Sprintf("%v", v)),
		))
	}
}

