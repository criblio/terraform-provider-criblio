package parser

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/pb33f/libopenapi"
	"go.yaml.in/yaml/v3"
)

var httpMethods = map[string]bool{
	"get":     true,
	"put":     true,
	"post":    true,
	"delete":  true,
	"options": true,
	"head":    true,
	"patch":   true,
	"trace":   true,
}

// ParseFile reads an OpenAPI document and returns annotated Terraform resources.
func ParseFile(filename string) ([]ResourceDef, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read spec: %v", err)
	}
	return Parse(content)
}

// Parse reads an OpenAPI document and returns annotated Terraform resources.
func Parse(content []byte) ([]ResourceDef, error) {
	if _, err := libopenapi.NewDocument(content); err != nil {
		return nil, fmt.Errorf("parse OpenAPI document: %v", err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("parse OpenAPI YAML: %v", err)
	}

	doc := documentMapping(&root)
	schemas, ok := lookupSchemas(doc)
	if !ok {
		return nil, fmt.Errorf("components.schemas mapping not found")
	}
	paths, ok := mappingValue(doc, "paths")
	if !ok || paths.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("paths mapping not found")
	}
	examples, _ := lookupExamples(doc)

	resources := map[string]*ResourceDef{}
	for index := 0; index < len(paths.Content); index += 2 {
		path := paths.Content[index].Value
		pathItem := paths.Content[index+1]
		if pathItem.Kind != yaml.MappingNode {
			continue
		}
		if err := collectOperations(resources, schemas, examples, path, pathItem); err != nil {
			return nil, err
		}
	}

	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	sort.Strings(names)

	defs := make([]ResourceDef, 0, len(names))
	for _, name := range names {
		resource := resources[name]
		if resource.SchemaName == "" {
			resource.SchemaName = resource.Create.RequestSchema
		}
		if err := populateFields(resource, schemas); err != nil {
			return nil, err
		}
		defs = append(defs, *resource)
	}
	return defs, nil
}

func collectOperations(resources map[string]*ResourceDef, schemas, examples *yaml.Node, path string, pathItem *yaml.Node) error {
	for index := 0; index < len(pathItem.Content); index += 2 {
		method := pathItem.Content[index].Value
		if !httpMethods[method] {
			continue
		}
		operation := pathItem.Content[index+1]
		if operation.Kind != yaml.MappingNode {
			continue
		}

		if boolAnnotation(operation, "x-terraform-resource") {
			name, ok := stringAnnotation(operation, "x-terraform-resource-name")
			if !ok || name == "" {
				return fmt.Errorf("%s %s missing x-terraform-resource-name", strings.ToUpper(method), path)
			}
			resource := ensureResource(resources, name)
			resource.Create = operationDef(method, path, operation, examples)
			resource.SchemaName = resource.Create.RequestSchema
			resource.Action = boolAnnotation(operation, "x-terraform-action")
			resource.NoRead = boolAnnotation(operation, "x-terraform-no-read")
		}
		if name, ok := stringAnnotation(operation, "x-terraform-list"); ok && name != "" {
			resource := ensureResource(resources, name)
			resource.List = operationDef(method, path, operation, examples)
			listName, ok := stringAnnotation(operation, "x-terraform-list-name")
			if !ok || listName == "" {
				listName = name + "s"
			}
			resource.ListName = listName
			resource.ListFileStem = snake(listName)
			resource.ListStructName = exportName(listName)
			resource.ListTypeName = "criblio_" + snake(listName)
		}

		for _, annotation := range []struct {
			key string
			set func(*ResourceDef, OperationDef)
		}{
			{"x-terraform-read", func(r *ResourceDef, op OperationDef) { r.Read = op }},
			{"x-terraform-update", func(r *ResourceDef, op OperationDef) { r.Update = op }},
			{"x-terraform-delete", func(r *ResourceDef, op OperationDef) { r.Delete = op }},
		} {
			name, ok := stringAnnotation(operation, annotation.key)
			if !ok || name == "" {
				continue
			}
			resource := ensureResource(resources, name)
			annotation.set(resource, operationDef(method, path, operation, examples))
		}
	}
	return nil
}

func ensureResource(resources map[string]*ResourceDef, name string) *ResourceDef {
	resource, ok := resources[name]
	if ok {
		return resource
	}
	resource = &ResourceDef{
		Name:       name,
		FileStem:   snake(name),
		TypeName:   "criblio_" + snake(name),
		StructName: exportName(name),
	}
	resources[name] = resource
	return resource
}

func operationDef(method, path string, operation, examples *yaml.Node) OperationDef {
	return OperationDef{
		Method:         strings.ToUpper(method),
		Path:           path,
		OperationID:    scalarValue(operation, "operationId"),
		RequestSchema:  schemaRefName(requestSchema(operation)),
		ResponseSchema: schemaRefName(responseSchema(operation)),
		PathParams:     pathParams(operation),
		QueryParams:    queryParams(operation),
		Examples:       requestExamples(operation, examples),
		ReadAfterWrite: boolAnnotation(operation, "x-terraform-read-after-write"),
		PreserveInputsAfterWrite: boolAnnotation(
			operation,
			"x-terraform-preserve-inputs-after-write",
		),
		ResetBody:  resetBody(operation),
		DeleteHook: scalarValue(operation, "x-terraform-delete-hook"),
	}
}

func populateFields(resource *ResourceDef, schemas *yaml.Node) error {
	schema, ok := mappingValue(schemas, resource.SchemaName)
	if !ok {
		return fmt.Errorf("resource %q schema %q not found", resource.Name, resource.SchemaName)
	}

	postFields := schemaPropertySet(schema)
	updateFields := map[string]bool{}
	if resource.Update.RequestSchema != "" {
		if updateSchema, ok := mappingValue(schemas, resource.Update.RequestSchema); ok {
			updateFields = schemaPropertySet(updateSchema)
		}
	}
	getFields := map[string]bool{}
	if resource.Read.ResponseSchema != "" {
		if readSchema, ok := mappingValue(schemas, resource.Read.ResponseSchema); ok {
			getFields = schemaPropertySet(readSchema)
		}
	}

	fields, variants, err := parseSchemaFields(resource.StructName, schema, schemas, postFields, updateFields, getFields)
	if err != nil {
		return err
	}
	if resource.Read.ResponseSchema != "" && resource.Read.ResponseSchema != resource.SchemaName {
		if readSchema, ok := mappingValue(schemas, resource.Read.ResponseSchema); ok {
			readFields, readVariants, err := parseSchemaFields(resource.StructName, readSchema, schemas, postFields, updateFields, getFields)
			if err != nil {
				return err
			}
			existing := map[string]bool{}
			for _, field := range fields {
				existing[field.APIName] = true
			}
			for _, field := range readFields {
				if existing[field.APIName] {
					continue
				}
				field.Computed = true
				field.Required = false
				field.Optional = false
				fields = append(fields, field)
			}
			variants = append(variants, readVariants...)
		}
	}
	fields = appendPathParams(fields, resource.Create.PathParams)
	fields = appendPathParams(fields, resource.Read.PathParams)
	fields = appendPathParams(fields, resource.Update.PathParams)
	fields = appendPathParams(fields, resource.Delete.PathParams)
	fields = appendQueryParams(fields, resource.Create.QueryParams)
	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].TerraformName < fields[j].TerraformName
	})
	resource.Fields = fields
	resource.OneOfVariants = variants
	applyResourceCompatibility(resource)
	return nil
}

func applyResourceCompatibility(resource *ResourceDef) {
	if resource == nil {
		return
	}
	if resource.Action {
		for index := range resource.Fields {
			field := &resource.Fields[index]
			if !field.Computed {
				field.ForceNew = true
			}
		}
	}
	if strings.HasPrefix(resource.TypeName, "criblio_search_") {
		fields := resource.Fields[:0]
		for _, field := range resource.Fields {
			if field.PathParam && field.TerraformName == "group_id" {
				continue
			}
			fields = append(fields, field)
		}
		resource.Fields = fields
	}
	if resource.StructName == "Notification" {
		fields := resource.Fields[:0]
		for _, field := range resource.Fields {
			if field.PathParam && field.TerraformName == "group_id" {
				continue
			}
			fields = append(fields, field)
		}
		resource.Fields = fields
	}
	if resource.StructName == "Collector" {
		makeCollectorVariantsOptionalComputed(resource.OneOfVariants)
	}
	if resource.StructName == "SearchDataset" {
		makeSearchDatasetHoistedFieldsComputed(resource.Fields)
		renameSearchDatasetProviderFields(resource.OneOfVariants)
		makeFieldsOptionalComputedFromValues(resource.OneOfVariants)
	}
	if resource.StructName != "MappingRuleset" {
		return
	}
	for index := range resource.Fields {
		field := &resource.Fields[index]
		if field.TerraformName == "id" {
			field.Required = false
			field.Optional = true
			field.Computed = true
		}
		if field.TerraformName == "conf" {
			makeMappingRulesetFunctionDefaultsOptional(field)
		}
	}
}

func makeSearchDatasetHoistedFieldsComputed(fields []FieldDef) {
	for index := range fields {
		field := &fields[index]
		switch field.TerraformName {
		case "id", "description", "provider_id", "type":
			field.Required = false
			field.Optional = false
			field.Computed = true
		}
	}
}

func renameSearchDatasetProviderFields(variants []OneOfVariantDef) {
	for variantIndex := range variants {
		renameProviderField(variants[variantIndex].Fields)
	}
}

func renameProviderField(fields []FieldDef) {
	for index := range fields {
		field := &fields[index]
		if field.APIName == "provider" {
			field.TerraformName = "provider_id"
			field.GoName = "ProviderID"
		}
		renameProviderField(field.Fields)
	}
}

func makeCollectorVariantsOptionalComputed(variants []OneOfVariantDef) {
	for variantIndex := range variants {
		makeFieldsOptionalComputed(variants[variantIndex].Fields)
		addCollectorPlanModifierHooks(variants[variantIndex].Fields)
	}
}

func makeFieldsOptionalComputedFromValues(variants []OneOfVariantDef) {
	for variantIndex := range variants {
		makeFieldsOptionalComputed(variants[variantIndex].Fields)
	}
}

func makeFieldsOptionalComputed(fields []FieldDef) {
	for fieldIndex := range fields {
		field := &fields[fieldIndex]
		field.Required = false
		field.Optional = true
		field.Computed = true
		field.OptionalComputed = true
		makeFieldsOptionalComputed(field.Fields)
	}
}

func addCollectorPlanModifierHooks(fields []FieldDef) {
	for fieldIndex := range fields {
		field := &fields[fieldIndex]
		switch field.TerraformName {
		case "conf", "input", "schedule":
			if field.Type == "object" {
				field.PlanModifierHook = "collectorPreferConfigOrStatePlanModifiers"
			}
		}
		addCollectorPlanModifierHooks(field.Fields)
	}
}

func makeMappingRulesetFunctionDefaultsOptional(field *FieldDef) {
	for index := range field.Fields {
		nested := &field.Fields[index]
		if nested.TerraformName == "functions" {
			for nestedIndex := range nested.Fields {
				functionField := &nested.Fields[nestedIndex]
				if functionField.TerraformName == "id" || functionField.TerraformName == "final" {
					functionField.Required = false
					functionField.Optional = true
					functionField.Computed = false
					functionField.FixedValue = ""
				}
				if functionField.TerraformName == "conf" {
					makeMappingRulesetAddNameOptional(functionField)
				}
			}
			return
		}
	}
}

func makeMappingRulesetAddNameOptional(field *FieldDef) {
	for index := range field.Fields {
		nested := &field.Fields[index]
		if nested.TerraformName != "add" {
			continue
		}
		for nestedIndex := range nested.Fields {
			addField := &nested.Fields[nestedIndex]
			if addField.TerraformName != "name" {
				continue
			}
			addField.Required = false
			addField.Optional = true
			addField.Computed = false
			addField.FixedValue = ""
		}
	}
}

func parseSchemaFields(modelName string, schema, schemas *yaml.Node, postFields, updateFields, getFields map[string]bool) ([]FieldDef, []OneOfVariantDef, error) {
	if _, ok := mappingValue(schema, "allOf"); ok {
		resolved, ok, err := objectSchema(schema, schemas)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			schema = resolved
		}
	}
	required := requiredSet(schema)
	properties, ok := mappingValue(schema, "properties")

	var fields []FieldDef
	var variants []OneOfVariantDef
	if ok && properties.Kind == yaml.MappingNode {
		for index := 0; index < len(properties.Content); index += 2 {
			apiName := properties.Content[index].Value
			property := properties.Content[index+1]
			if ignoredAnnotation(property) {
				continue
			}

			if oneOf, ok := mappingValue(property, "oneOf"); ok && oneOf.Kind == yaml.SequenceNode {
				parsed, err := parseOneOfVariants(oneOf, schemas)
				if err != nil {
					return nil, nil, err
				}
				variants = append(variants, parsed...)
				continue
			}

			field, err := fieldDef(modelName, apiName, property, schemas)
			if err != nil {
				return nil, nil, err
			}
			applyFieldAnnotations(&field, property, required[apiName], postFields[apiName], updateFields[apiName], getFields[apiName])
			fields = append(fields, field)
		}
	}
	if oneOf, ok := mappingValue(schema, "oneOf"); ok && oneOf.Kind == yaml.SequenceNode {
		parsed, err := parseOneOfVariants(oneOf, schemas)
		if err != nil {
			return nil, nil, err
		}
		variants = append(variants, parsed...)
	}
	return fields, variants, nil
}

func parseOneOfVariants(oneOf, schemas *yaml.Node) ([]OneOfVariantDef, error) {
	var variants []OneOfVariantDef
	for _, variantRef := range oneOf.Content {
		schemaName := schemaRefName(variantRef)
		if schemaName == "" {
			continue
		}
		variantSchema, ok := mappingValue(schemas, schemaName)
		if !ok {
			return nil, fmt.Errorf("oneOf schema %q not found", schemaName)
		}
		variantFields, _, err := parseSchemaFields(schemaName, variantSchema, schemas, schemaPropertySet(variantSchema), schemaPropertySet(variantSchema), schemaPropertySet(variantSchema))
		if err != nil {
			return nil, err
		}
		tfName := schemaName
		if renamed, ok := stringAnnotation(variantSchema, "x-terraform-name"); ok && renamed != "" {
			tfName = renamed
		}
		variants = append(variants, OneOfVariantDef{
			APIName:            schemaName,
			TerraformName:      snake(tfName),
			GoName:             exportName(schemaName),
			ModelName:          exportName(schemaName) + "Model",
			SchemaName:         schemaName,
			DiscriminatorValue: discriminatorValue(variantFields),
			Fields:             variantFields,
		})
	}
	return variants, nil
}

func applyFieldAnnotations(field *FieldDef, property *yaml.Node, required, requestField, updateField, getField bool) {
	field.RequestField = requestField
	field.UpdateField = updateField
	field.Required = required
	field.Optional = !field.Required
	if boolAnnotation(property, "readOnly") || boolAnnotation(property, "x-terraform-computed") || (getField && !requestField) {
		field.Computed = true
		field.Required = false
		field.Optional = false
	}
	if boolAnnotation(property, "x-terraform-optional-computed") {
		field.Computed = true
		field.Required = false
		field.Optional = true
		field.OptionalComputed = true
	}
	if boolAnnotation(property, "writeOnly") {
		field.Sensitive = true
		field.PreferState = true
		field.ApplyStrategy = "stringFromAPIOrPrior"
	}
	if boolAnnotation(property, "x-terraform-sensitive") {
		field.Sensitive = true
	}
	if boolAnnotation(property, "x-terraform-prefer-state") || scalarValue(property, "x-speakeasy-plan-modifiers") == "PreferState" {
		field.PreferState = true
	}
	if suppressDiffAnnotation(property) {
		field.SuppressDiff = true
	}
	if boolAnnotation(property, "x-terraform-force-new") {
		field.ForceNew = true
	}
	if fixedValue := fixedValueAnnotation(property); fixedValue != "" {
		field.FixedValue = fixedValue
		field.Required = false
		field.Optional = true
		field.Computed = true
	}
	if field.PreferState && field.Sensitive {
		field.ApplyStrategy = "stringFromAPIOrPrior"
	} else if field.PreferState {
		field.ApplyStrategy = "preferState"
	}
}

func fieldDef(modelName, apiName string, property, schemas *yaml.Node) (FieldDef, error) {
	tfName := apiName
	if renamed, ok := stringAnnotation(property, "x-terraform-name"); ok && renamed != "" {
		tfName = renamed
	}
	schemaForType := property
	if schemaName := directSchemaRefName(property); schemaName != "" {
		resolved, found := mappingValue(schemas, schemaName)
		if !found {
			return FieldDef{}, fmt.Errorf("field schema %q not found", schemaName)
		}
		schemaForType = resolved
	}
	fieldType := schemaType(schemaForType)
	if objectSchema, ok, err := objectSchemaForProperty(property, schemas); err != nil {
		return FieldDef{}, err
	} else if ok && objectSchema != nil {
		fieldType = "object"
	}
	fieldElementType := elementType(schemaForType)
	if fieldType == "array" {
		if items, ok := mappingValue(schemaForType, "items"); ok {
			if objectSchema, ok, err := objectSchemaForProperty(items, schemas); err != nil {
				return FieldDef{}, err
			} else if ok && objectSchema != nil {
				fieldElementType = "object"
			}
		}
	}
	field := FieldDef{
		APIName:       apiName,
		TerraformName: snake(tfName),
		GoName:        exportName(tfName),
		Type:          fieldType,
		ElementType:   fieldElementType,
		Description:   scalarValue(property, "description"),
		CustomType:    scalarValue(property, "x-terraform-custom-type"),
		ElementCustomType: scalarValue(
			property,
			"x-terraform-element-custom-type",
		),
		PlanModifierHook: scalarValue(property, "x-terraform-plan-modifier-hook"),
		UseStateForUnknown: boolAnnotation(
			property,
			"x-terraform-use-state-for-unknown",
		),
		EmitEmpty:    boolAnnotation(property, "x-terraform-emit-empty"),
		FixedValue:   fixedValueAnnotation(property),
		ObjectAsJSON: boolAnnotation(property, "x-terraform-object-as-json"),
		NotNull:      boolAnnotation(property, "x-terraform-not-null"),
		ValidJSON:    boolAnnotation(property, "x-terraform-valid-json"),
		PipelineFunctionID: boolAnnotation(
			property,
			"x-terraform-pipeline-function-id-validator",
		),
		ReadOnly:  boolAnnotation(property, "readOnly"),
		WriteOnly: boolAnnotation(property, "writeOnly"),
		Enum:      enumValues(property),
	}
	if field.Type == "array" && field.ElementType == "object" {
		items, ok := mappingValue(schemaForType, "items")
		if !ok {
			return FieldDef{}, fmt.Errorf("%s.%s array field missing items schema", modelName, apiName)
		}
		if schemaName := schemaRefName(items); schemaName != "" {
			var found bool
			items, found = mappingValue(schemas, schemaName)
			if !found {
				return FieldDef{}, fmt.Errorf("array item schema %q not found", schemaName)
			}
		}
		if resolved, ok, err := objectSchema(items, schemas); err != nil {
			return FieldDef{}, err
		} else if ok {
			items = resolved
		}
		fields, _, err := parseSchemaFields(modelName+exportName(tfName), items, schemas, schemaPropertySet(items), schemaPropertySet(items), schemaPropertySet(items))
		if err != nil {
			return FieldDef{}, err
		}
		field.Fields = fields
		nestedName := nestedModelPrefix(modelName, tfName)
		field.NestedModelName = nestedName + "Model"
		field.NestedAPIModelName = nestedName + "APIModel"
		field.NestedAttrTypes = nestedName + "AttrTypes"
	}
	if field.Type == "object" && !field.ObjectAsJSON {
		objectSchema := schemaForType
		if resolved, ok, err := objectSchemaForProperty(property, schemas); err != nil {
			return FieldDef{}, err
		} else if ok {
			objectSchema = resolved
		}
		if additional, ok := mappingValue(objectSchema, "additionalProperties"); ok && additional.Kind == yaml.MappingNode {
			mapSchema := additional
			if resolved, ok, err := objectSchemaForProperty(additional, schemas); err != nil {
				return FieldDef{}, err
			} else if ok {
				mapSchema = resolved
			}
			if properties, ok := mappingValue(mapSchema, "properties"); ok && properties.Kind == yaml.MappingNode {
				fields, _, err := parseSchemaFields(modelName+exportName(tfName), mapSchema, schemas, schemaPropertySet(mapSchema), schemaPropertySet(mapSchema), schemaPropertySet(mapSchema))
				if err != nil {
					return FieldDef{}, err
				}
				field.ElementType = "object"
				field.Fields = fields
				nestedName := nestedModelPrefix(modelName, tfName)
				field.NestedModelName = nestedName + "Model"
				field.NestedAPIModelName = nestedName + "APIModel"
				field.NestedAttrTypes = nestedName + "AttrTypes"
				return field, nil
			}
		}
		if properties, ok := mappingValue(objectSchema, "properties"); ok && properties.Kind == yaml.MappingNode {
			fields, _, err := parseSchemaFields(modelName+exportName(tfName), objectSchema, schemas, schemaPropertySet(objectSchema), schemaPropertySet(objectSchema), schemaPropertySet(objectSchema))
			if err != nil {
				return FieldDef{}, err
			}
			field.Fields = fields
			nestedName := nestedModelPrefix(modelName, tfName)
			field.NestedModelName = nestedName + "Model"
			field.NestedAPIModelName = nestedName + "APIModel"
			field.NestedAttrTypes = nestedName + "AttrTypes"
		}
	}
	return field, nil
}

func ignoredAnnotation(property *yaml.Node) bool {
	return boolAnnotation(property, "x-terraform-ignore") || boolAnnotation(property, "x-speakeasy-terraform-ignore")
}

func enumValues(property *yaml.Node) []string {
	if value := scalarValue(property, "const"); value != "" {
		return []string{value}
	}
	enum, ok := mappingValue(property, "enum")
	if !ok || enum.Kind != yaml.SequenceNode {
		return nil
	}
	values := make([]string, 0, len(enum.Content))
	for _, item := range enum.Content {
		if item.Kind == yaml.ScalarNode {
			values = append(values, item.Value)
		}
	}
	return values
}

func discriminatorValue(fields []FieldDef) string {
	for _, field := range fields {
		if field.APIName == "type" && len(field.Enum) == 1 {
			return field.Enum[0]
		}
		if field.APIName != "collector" {
			continue
		}
		for _, nested := range field.Fields {
			if nested.APIName == "type" && len(nested.Enum) == 1 {
				return nested.Enum[0]
			}
		}
	}
	return ""
}

func nestedModelPrefix(modelName, fieldName string) string {
	prefix := modelName + exportName(fieldName)
	if prefix == modelName+"API" {
		return modelName + "APIField"
	}
	return prefix
}

func objectSchemaForProperty(property, schemas *yaml.Node) (*yaml.Node, bool, error) {
	if scalarValue(property, "type") == "array" {
		return nil, false, nil
	}
	schema := property
	if schemaName := schemaRefName(property); schemaName != "" {
		var found bool
		schema, found = mappingValue(schemas, schemaName)
		if !found {
			return nil, false, fmt.Errorf("object schema %q not found", schemaName)
		}
	}
	return objectSchema(schema, schemas)
}

func objectSchema(schema, schemas *yaml.Node) (*yaml.Node, bool, error) {
	if properties, ok := mappingValue(schema, "properties"); ok && properties.Kind == yaml.MappingNode {
		return schema, true, nil
	}

	if oneOf, ok := mappingValue(schema, "oneOf"); ok && oneOf.Kind == yaml.SequenceNode {
		for _, item := range oneOf.Content {
			itemSchema := item
			if schemaName := schemaRefName(item); schemaName != "" {
				var found bool
				itemSchema, found = mappingValue(schemas, schemaName)
				if !found {
					return nil, false, fmt.Errorf("oneOf schema %q not found", schemaName)
				}
			}
			resolved, ok, err := objectSchema(itemSchema, schemas)
			if err != nil {
				return nil, false, err
			}
			if ok {
				return resolved, true, nil
			}
		}
		return nil, false, nil
	}

	allOf, ok := mappingValue(schema, "allOf")
	if !ok || allOf.Kind != yaml.SequenceNode {
		return nil, false, nil
	}

	properties := &yaml.Node{Kind: yaml.MappingNode}
	required := &yaml.Node{Kind: yaml.SequenceNode}
	seenRequired := map[string]bool{}
	for _, item := range allOf.Content {
		itemSchema := item
		if schemaName := schemaRefName(item); schemaName != "" {
			var found bool
			itemSchema, found = mappingValue(schemas, schemaName)
			if !found {
				return nil, false, fmt.Errorf("allOf schema %q not found", schemaName)
			}
		}
		resolved, ok, err := objectSchema(itemSchema, schemas)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			continue
		}
		if itemProperties, ok := mappingValue(resolved, "properties"); ok && itemProperties.Kind == yaml.MappingNode {
			for index := 0; index < len(itemProperties.Content); index += 2 {
				setOrAppendMappingValue(properties, itemProperties.Content[index], itemProperties.Content[index+1])
			}
		}
		if itemRequired, ok := mappingValue(resolved, "required"); ok && itemRequired.Kind == yaml.SequenceNode {
			for _, name := range itemRequired.Content {
				if seenRequired[name.Value] {
					continue
				}
				seenRequired[name.Value] = true
				required.Content = append(required.Content, name)
			}
		}
	}
	if len(properties.Content) == 0 {
		return nil, false, nil
	}

	merged := &yaml.Node{Kind: yaml.MappingNode}
	merged.Content = append(merged.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "type"}, &yaml.Node{Kind: yaml.ScalarNode, Value: "object"})
	merged.Content = append(merged.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "properties"}, properties)
	if len(required.Content) > 0 {
		merged.Content = append(merged.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "required"}, required)
	}
	return merged, true, nil
}

func setOrAppendMappingValue(node, key, value *yaml.Node) {
	for index := 0; index < len(node.Content); index += 2 {
		if node.Content[index].Value == key.Value {
			node.Content[index+1] = value
			return
		}
	}
	node.Content = append(node.Content, key, value)
}

func appendPathParams(fields []FieldDef, params []FieldDef) []FieldDef {
	existing := map[string]bool{}
	for _, field := range fields {
		existing[field.TerraformName] = true
	}
	for _, param := range params {
		if existing[param.TerraformName] {
			for index := range fields {
				if fields[index].TerraformName == param.TerraformName {
					if fields[index].FixedValue == "" {
						fields[index].Required = true
						fields[index].Optional = false
						fields[index].Computed = false
					} else {
						fields[index].Required = false
						fields[index].Optional = true
						fields[index].Computed = true
					}
					fields[index].ForceNew = true
					fields[index].PathParam = true
				}
			}
			continue
		}
		fields = append(fields, param)
		existing[param.TerraformName] = true
	}
	return fields
}

func pathParams(operation *yaml.Node) []FieldDef {
	return parametersByLocation(operation, "path")
}

func queryParams(operation *yaml.Node) []FieldDef {
	return parametersByLocation(operation, "query")
}

func parametersByLocation(operation *yaml.Node, location string) []FieldDef {
	parameters, ok := mappingValue(operation, "parameters")
	if !ok || parameters.Kind != yaml.SequenceNode {
		return nil
	}
	var params []FieldDef
	for _, parameter := range parameters.Content {
		if scalarValue(parameter, "in") != location {
			continue
		}
		if location == "query" && !boolAnnotation(parameter, "required") {
			continue
		}
		apiName := scalarValue(parameter, "name")
		field := FieldDef{
			APIName:       apiName,
			TerraformName: snake(apiName),
			GoName:        exportName(apiName),
			Type:          "string",
			Description:   scalarValue(parameter, "description"),
			Required:      true,
			ForceNew:      true,
		}
		if schema, ok := mappingValue(parameter, "schema"); ok {
			if fixedValue := fixedValueAnnotation(schema); fixedValue != "" {
				field.FixedValue = fixedValue
				field.Required = false
				field.Optional = true
				field.Computed = true
			}
		}
		if location == "path" {
			field.PathParam = true
		} else {
			field.QueryParam = true
		}
		params = append(params, field)
	}
	return params
}

func appendQueryParams(fields []FieldDef, params []FieldDef) []FieldDef {
	existing := map[string]bool{}
	for _, field := range fields {
		existing[field.TerraformName] = true
	}
	for _, param := range params {
		if existing[param.TerraformName] {
			for index := range fields {
				if fields[index].TerraformName == param.TerraformName {
					if fields[index].FixedValue == "" {
						fields[index].Required = true
						fields[index].Optional = false
						fields[index].Computed = false
					} else {
						fields[index].Required = false
						fields[index].Optional = true
						fields[index].Computed = true
					}
					fields[index].ForceNew = true
					fields[index].QueryParam = true
				}
			}
			continue
		}
		fields = append(fields, param)
		existing[param.TerraformName] = true
	}
	return fields
}

func requestSchema(operation *yaml.Node) *yaml.Node {
	mediaType := requestJSONMediaType(operation)
	if mediaType == nil {
		return nil
	}
	schema, _ := mappingValue(mediaType, "schema")
	return schema
}

func resetBody(operation *yaml.Node) any {
	node, ok := mappingValue(operation, "x-terraform-delete-reset-body")
	if !ok {
		return nil
	}
	var output any
	if err := node.Decode(&output); err != nil {
		return nil
	}
	return output
}

func requestExamples(operation, examples *yaml.Node) []ExampleDef {
	mediaType := requestJSONMediaType(operation)
	if mediaType == nil {
		return nil
	}

	var defs []ExampleDef
	exampleItems, ok := mappingValue(mediaType, "examples")
	if ok && exampleItems.Kind == yaml.MappingNode {
		for index := 0; index < len(exampleItems.Content); index += 2 {
			name := exampleItems.Content[index].Value
			example := resolveExample(exampleItems.Content[index+1], examples)
			value, ok := mappingValue(example, "value")
			if !ok {
				continue
			}
			decoded, err := decodeExampleValue(value)
			if err != nil {
				continue
			}
			defs = append(defs, ExampleDef{
				Name:    name,
				Summary: scalarValue(example, "summary"),
				Value:   decoded,
			})
		}
	}
	if len(defs) > 0 {
		return defs
	}

	example, ok := mappingValue(mediaType, "example")
	if !ok {
		return nil
	}
	decoded, err := decodeExampleValue(example)
	if err != nil {
		return nil
	}
	return []ExampleDef{{Name: "example", Value: decoded}}
}

func requestJSONMediaType(operation *yaml.Node) *yaml.Node {
	requestBody, ok := mappingValue(operation, "requestBody")
	if !ok {
		return nil
	}
	content, ok := mappingValue(requestBody, "content")
	if !ok || content.Kind != yaml.MappingNode {
		return nil
	}
	mediaType, ok := mappingValue(content, "application/json")
	if !ok {
		return nil
	}
	return mediaType
}

func resolveExample(example, examples *yaml.Node) *yaml.Node {
	ref, ok := mappingValue(example, "$ref")
	if !ok || ref.Kind != yaml.ScalarNode {
		return example
	}
	name := strings.TrimPrefix(ref.Value, "#/components/examples/")
	if name == ref.Value || examples == nil {
		return example
	}
	resolved, ok := mappingValue(examples, name)
	if !ok {
		return example
	}
	return resolved
}

func decodeExampleValue(node *yaml.Node) (any, error) {
	var value any
	if err := node.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

func responseSchema(operation *yaml.Node) *yaml.Node {
	responses, ok := mappingValue(operation, "responses")
	if !ok || responses.Kind != yaml.MappingNode {
		return nil
	}
	for _, status := range []string{"200", "201", "202"} {
		if schema := responseJSONSchema(responses, status); schema != nil {
			return schema
		}
	}
	for index := 0; index < len(responses.Content); index += 2 {
		if schema := responseJSONSchema(responses, responses.Content[index].Value); schema != nil {
			return schema
		}
	}
	return nil
}

func responseJSONSchema(responses *yaml.Node, status string) *yaml.Node {
	response, ok := mappingValue(responses, status)
	if !ok {
		return nil
	}
	content, ok := mappingValue(response, "content")
	if !ok || content.Kind != yaml.MappingNode {
		return nil
	}
	mediaType, ok := mappingValue(content, "application/json")
	if !ok {
		return nil
	}
	schema, _ := mappingValue(mediaType, "schema")
	return schema
}

func schemaRefName(node *yaml.Node) string {
	if node == nil {
		return ""
	}
	if name := directSchemaRefName(node); name != "" {
		return name
	}
	if items, ok := mappingValue(node, "items"); ok {
		return schemaRefName(items)
	}
	return ""
}

func directSchemaRefName(node *yaml.Node) string {
	if node == nil {
		return ""
	}
	if ref, ok := mappingValue(node, "$ref"); ok {
		return strings.TrimPrefix(ref.Value, "#/components/schemas/")
	}
	return ""
}

func schemaPropertySet(schema *yaml.Node) map[string]bool {
	set := map[string]bool{}
	properties, ok := mappingValue(schema, "properties")
	if !ok || properties.Kind != yaml.MappingNode {
		return set
	}
	for index := 0; index < len(properties.Content); index += 2 {
		set[properties.Content[index].Value] = true
	}
	return set
}

func requiredSet(schema *yaml.Node) map[string]bool {
	set := map[string]bool{}
	required, ok := mappingValue(schema, "required")
	if !ok || required.Kind != yaml.SequenceNode {
		return set
	}
	for _, item := range required.Content {
		set[item.Value] = true
	}
	return set
}

func schemaType(node *yaml.Node) string {
	typ := scalarValue(node, "type")
	switch typ {
	case "string", "boolean", "integer", "number", "array", "object":
		return typ
	default:
		if ref, ok := mappingValue(node, "$ref"); ok && ref.Kind == yaml.ScalarNode {
			return exportName(strings.TrimPrefix(ref.Value, "#/components/schemas/"))
		}
		return "string"
	}
}

func elementType(node *yaml.Node) string {
	items, ok := mappingValue(node, "items")
	if !ok {
		return ""
	}
	return schemaType(items)
}

func lookupSchemas(root *yaml.Node) (*yaml.Node, bool) {
	components, ok := mappingValue(root, "components")
	if !ok || components.Kind != yaml.MappingNode {
		return nil, false
	}
	return mappingValue(components, "schemas")
}

func lookupExamples(root *yaml.Node) (*yaml.Node, bool) {
	components, ok := mappingValue(root, "components")
	if !ok || components.Kind != yaml.MappingNode {
		return nil, false
	}
	return mappingValue(components, "examples")
}

func documentMapping(root *yaml.Node) *yaml.Node {
	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		return root.Content[0]
	}
	return root
}

func mappingValue(node *yaml.Node, key string) (*yaml.Node, bool) {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil, false
	}
	for index := 0; index < len(node.Content); index += 2 {
		if node.Content[index].Value == key {
			return node.Content[index+1], true
		}
	}
	return nil, false
}

func scalarValue(node *yaml.Node, key string) string {
	value, ok := mappingValue(node, key)
	if !ok || value.Kind != yaml.ScalarNode {
		return ""
	}
	return value.Value
}

func boolAnnotation(node *yaml.Node, key string) bool {
	value, ok := mappingValue(node, key)
	return ok && value.Kind == yaml.ScalarNode && value.Value == "true"
}

func suppressDiffAnnotation(node *yaml.Node) bool {
	value, ok := mappingValue(node, "x-terraform-suppress-diff")
	if !ok || value.Kind != yaml.ScalarNode {
		return false
	}
	switch strings.ToLower(value.Value) {
	case "true", "explicit", "explicit-suppress", "explicit_suppress":
		return true
	default:
		return false
	}
}

func stringAnnotation(node *yaml.Node, key string) (string, bool) {
	value, ok := mappingValue(node, key)
	if !ok || value.Kind != yaml.ScalarNode {
		return "", false
	}
	return value.Value, true
}

func fixedValueAnnotation(node *yaml.Node) string {
	if value := scalarValue(node, "x-terraform-fixed-value"); value != "" {
		return value
	}
	return scalarValue(node, "const")
}

func snake(value string) string {
	var output bytes.Buffer
	var previousLower bool
	for index, r := range value {
		if r == '-' || r == ' ' {
			output.WriteByte('_')
			previousLower = false
			continue
		}
		if unicode.IsUpper(r) && index > 0 && previousLower {
			output.WriteByte('_')
		}
		if r == '_' {
			output.WriteByte('_')
			previousLower = false
			continue
		}
		output.WriteRune(unicode.ToLower(r))
		previousLower = unicode.IsLower(r) || unicode.IsDigit(r)
	}
	return output.String()
}

func exportName(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '/'
	})
	var output strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		lower := strings.ToLower(part)
		if initialism, ok := initialisms[lower]; ok {
			output.WriteString(initialism)
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		output.WriteString(string(runes))
	}
	name := output.String()
	for lower, initialism := range initialisms {
		name = strings.ReplaceAll(name, exportPlain(lower), initialism)
	}
	return name
}

var initialisms = map[string]string{
	"acl":  "ACL",
	"api":  "API",
	"id":   "ID",
	"url":  "URL",
	"uri":  "URI",
	"json": "JSON",
	"tls":  "TLS",
}

func exportPlain(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
