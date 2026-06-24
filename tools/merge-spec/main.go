// Command merge-spec applies Terraform provider overlays to the upstream Cribl OpenAPI spec.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

const (
	defaultInput                = "upstream-openapi.yml"
	defaultOverlay              = "terraform-overlay.yml"
	defaultMgmtInput            = "openapi-mgmt.yml"
	defaultMgmtOverlay          = "terraform-mgmt-overlay.yml"
	defaultOutput               = "merged-spec.yml"
	defaultCloudOnlyPathsOutput = "internal/auth/cloud_only_paths.go"
	groupPrefix                 = "/m/{groupId}"
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

func main() {
	input := flag.String("input", defaultInput, "upstream OpenAPI YAML file")
	overlay := flag.String("overlay", defaultOverlay, "Terraform overlay YAML file")
	mgmtInput := flag.String("mgmt-input", defaultMgmtInput, "management OpenAPI YAML file")
	mgmtOverlay := flag.String("mgmt-overlay", defaultMgmtOverlay, "management Terraform overlay YAML file")
	output := flag.String("output", defaultOutput, "merged OpenAPI YAML file")
	flag.Parse()

	if err := runWithConfig(mergeConfig{
		InputPath:                *input,
		OverlayPath:              *overlay,
		MgmtInputPath:            *mgmtInput,
		MgmtOverlayPath:          *mgmtOverlay,
		OutputPath:               *output,
		CloudOnlyPathsOutputPath: defaultCloudOnlyPathsOutput,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "merge spec: %v\n", err)
		os.Exit(1)
	}
}

type mergeConfig struct {
	InputPath                string
	OverlayPath              string
	MgmtInputPath            string
	MgmtOverlayPath          string
	OutputPath               string
	CloudOnlyPathsOutputPath string
}

func run(inputPath, overlayPath, outputPath, cloudOnlyPathsOutputPath string) error {
	return runWithConfig(mergeConfig{
		InputPath:                inputPath,
		OverlayPath:              overlayPath,
		OutputPath:               outputPath,
		CloudOnlyPathsOutputPath: cloudOnlyPathsOutputPath,
	})
}

func runWithConfig(config mergeConfig) error {
	spec, err := readYAML(config.InputPath)
	if err != nil {
		return fmt.Errorf("read input spec: %v", err)
	}

	cloudOnlyPaths, err := collectCloudOnlyPaths(spec)
	if err != nil {
		return err
	}
	if err := prefixGroupScopedPaths(spec); err != nil {
		return err
	}
	if err := applyOverlayFile(spec, config.OverlayPath); err != nil {
		return err
	}
	if err := unwrapCountedResponses(spec); err != nil {
		return err
	}
	if err := mergeManagementSpec(spec, config.MgmtInputPath, config.MgmtOverlayPath); err != nil {
		return err
	}

	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(spec); err != nil {
		return fmt.Errorf("encode merged spec: %v", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("close YAML encoder: %v", err)
	}

	if err := os.WriteFile(config.OutputPath, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("write output spec: %v", err)
	}
	if err := writeCloudOnlyPaths(config.CloudOnlyPathsOutputPath, cloudOnlyPaths); err != nil {
		return err
	}
	return nil
}

func mergeManagementSpec(spec *yaml.Node, mgmtInputPath, mgmtOverlayPath string) error {
	if mgmtInputPath == "" {
		return nil
	}
	if _, err := os.Stat(mgmtInputPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat management spec: %v", err)
	}
	if _, err := os.Stat(mgmtOverlayPath); err != nil {
		return fmt.Errorf("stat management overlay: %v", err)
	}

	mgmtSpec, err := readYAML(mgmtInputPath)
	if err != nil {
		return fmt.Errorf("read management spec: %v", err)
	}
	if err := applyOverlayFile(mgmtSpec, mgmtOverlayPath); err != nil {
		return fmt.Errorf("apply management overlay: %v", err)
	}
	if err := unwrapCountedResponses(mgmtSpec); err != nil {
		return fmt.Errorf("unwrap management responses: %v", err)
	}
	if err := mergeAnnotatedPaths(spec, mgmtSpec); err != nil {
		return err
	}
	return mergeComponentSchemas(spec, mgmtSpec)
}

func mergeAnnotatedPaths(target, source *yaml.Node) error {
	targetPaths, ok := mappingValue(documentMapping(target), "paths")
	if !ok || targetPaths.Kind != yaml.MappingNode {
		return fmt.Errorf("target spec paths mapping not found")
	}
	sourcePaths, ok := mappingValue(documentMapping(source), "paths")
	if !ok || sourcePaths.Kind != yaml.MappingNode {
		return fmt.Errorf("management spec paths mapping not found")
	}

	for index := 0; index < len(sourcePaths.Content); index += 2 {
		pathKey := sourcePaths.Content[index]
		pathItem := sourcePaths.Content[index+1]
		if !pathHasTerraformAnnotation(pathItem) {
			continue
		}
		if _, exists := mappingValue(targetPaths, pathKey.Value); exists {
			return fmt.Errorf("management path %q already exists in target spec", pathKey.Value)
		}
		targetPaths.Content = append(targetPaths.Content, cloneNode(pathKey), cloneNode(pathItem))
	}
	return nil
}

func pathHasTerraformAnnotation(pathItem *yaml.Node) bool {
	if pathItem.Kind != yaml.MappingNode {
		return false
	}
	for index := 0; index < len(pathItem.Content); index += 2 {
		method := pathItem.Content[index].Value
		if !httpMethods[method] {
			continue
		}
		if operationHasTerraformAnnotation(pathItem.Content[index+1]) {
			return true
		}
	}
	return false
}

func operationHasTerraformAnnotation(operation *yaml.Node) bool {
	for _, key := range []string{
		"x-terraform-resource",
		"x-terraform-read",
		"x-terraform-update",
		"x-terraform-delete",
		"x-terraform-list",
	} {
		if _, ok := mappingValue(operation, key); ok {
			return true
		}
	}
	return false
}

func mergeComponentSchemas(target, source *yaml.Node) error {
	targetSchemas, ok := lookupComponentsSchemas(documentMapping(target))
	if !ok {
		return fmt.Errorf("target components schemas mapping not found")
	}
	sourceSchemas, ok := lookupComponentsSchemas(documentMapping(source))
	if !ok {
		return fmt.Errorf("management components schemas mapping not found")
	}

	for index := 0; index < len(sourceSchemas.Content); index += 2 {
		name := sourceSchemas.Content[index]
		schema := sourceSchemas.Content[index+1]
		if existing, exists := mappingValue(targetSchemas, name.Value); exists {
			if !nodesEqual(existing, schema) {
				return fmt.Errorf("management schema %q already exists with different content", name.Value)
			}
			continue
		}
		targetSchemas.Content = append(targetSchemas.Content, cloneNode(name), cloneNode(schema))
	}
	return nil
}

func nodesEqual(a, b *yaml.Node) bool {
	var left bytes.Buffer
	leftEncoder := yaml.NewEncoder(&left)
	_ = leftEncoder.Encode(a)
	_ = leftEncoder.Close()

	var right bytes.Buffer
	rightEncoder := yaml.NewEncoder(&right)
	_ = rightEncoder.Encode(b)
	_ = rightEncoder.Close()
	return bytes.Equal(left.Bytes(), right.Bytes())
}

func collectCloudOnlyPaths(spec *yaml.Node) ([]string, error) {
	paths, ok := mappingValue(documentMapping(spec), "paths")
	if !ok || paths.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("spec paths mapping not found")
	}

	var cloudOnlyPaths []string
	for index := 0; index < len(paths.Content); index += 2 {
		path := paths.Content[index].Value
		pathItem := paths.Content[index+1]
		if pathItem.Kind != yaml.MappingNode {
			continue
		}
		if pathIsCloudOnly(pathItem) {
			cloudOnlyPaths = append(cloudOnlyPaths, path)
		}
	}
	sort.Strings(cloudOnlyPaths)
	return cloudOnlyPaths, nil
}

func pathIsCloudOnly(pathItem *yaml.Node) bool {
	foundOperation := false
	for index := 0; index < len(pathItem.Content); index += 2 {
		method := pathItem.Content[index].Value
		if !httpMethods[method] {
			continue
		}
		foundOperation = true
		operation := pathItem.Content[index+1]
		availability, ok := mappingValue(operation, "x-cribl-availability")
		if !ok || availability.Kind != yaml.ScalarNode || availability.Value != "cloud" {
			return false
		}
	}
	return foundOperation
}

func writeCloudOnlyPaths(filename string, paths []string) error {
	var output bytes.Buffer
	output.WriteString("// Code generated by tools/merge-spec. DO NOT EDIT.\n")
	output.WriteString("// Source: x-cribl-availability annotations in upstream-openapi.yml.\n")
	output.WriteString("// Regenerate: make merge\n")
	output.WriteString("package auth\n\n")
	output.WriteString("func init() {\n")
	output.WriteString("\tcloudOnlyPaths = map[string]bool{\n")
	for _, path := range paths {
		fmt.Fprintf(&output, "\t\t%q: true,\n", path)
	}
	output.WriteString("\t}\n")
	output.WriteString("}\n")

	if err := os.WriteFile(filename, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("write cloud-only paths: %v", err)
	}
	return nil
}

func readYAML(filename string) (*yaml.Node, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, err
	}
	return &root, nil
}

func applyOverlayFile(spec *yaml.Node, overlayPath string) error {
	overlay, err := readYAML(overlayPath)
	if err != nil {
		return fmt.Errorf("read overlay: %v", err)
	}

	actions, ok := mappingValue(documentMapping(overlay), "actions")
	if !ok || actions.Kind != yaml.SequenceNode {
		return fmt.Errorf("overlay must contain an actions sequence")
	}

	for index, action := range actions.Content {
		target, update, err := parseAction(action)
		if err != nil {
			return fmt.Errorf("overlay action %d: %v", index, err)
		}
		targetNode, err := lookupTarget(documentMapping(spec), target)
		if err != nil {
			return fmt.Errorf("overlay action %d target %q: %v", index, target, err)
		}
		mergeMapping(targetNode, update)
	}

	return nil
}

func parseAction(action *yaml.Node) (string, *yaml.Node, error) {
	if action.Kind != yaml.MappingNode {
		return "", nil, fmt.Errorf("action must be a mapping")
	}

	targetNode, ok := mappingValue(action, "target")
	if !ok || targetNode.Kind != yaml.ScalarNode {
		return "", nil, fmt.Errorf("target scalar is required")
	}
	updateNode, ok := mappingValue(action, "update")
	if !ok || updateNode.Kind != yaml.MappingNode {
		return "", nil, fmt.Errorf("update mapping is required")
	}
	return targetNode.Value, updateNode, nil
}

func lookupTarget(root *yaml.Node, target string) (*yaml.Node, error) {
	const prefix = "$."
	if !strings.HasPrefix(target, prefix) {
		return nil, fmt.Errorf("must start with %q", prefix)
	}

	node := root
	for _, segment := range strings.Split(strings.TrimPrefix(target, prefix), ".") {
		if segment == "" {
			return nil, fmt.Errorf("empty target segment")
		}
		next, ok := mappingValue(node, segment)
		if !ok {
			return nil, fmt.Errorf("segment %q not found", segment)
		}
		node = next
	}
	return node, nil
}

func prefixGroupScopedPaths(spec *yaml.Node) error {
	paths, ok := mappingValue(documentMapping(spec), "paths")
	if !ok || paths.Kind != yaml.MappingNode {
		return fmt.Errorf("spec paths mapping not found")
	}

	var additions []*yaml.Node
	var removeKeys []string
	for index := 0; index < len(paths.Content); index += 2 {
		key := paths.Content[index]
		value := paths.Content[index+1]
		if !shouldPrefixPath(key.Value) {
			continue
		}

		prefixedPath := groupPrefix + key.Value
		if _, exists := mappingValue(paths, prefixedPath); exists {
			removeKeys = append(removeKeys, key.Value)
			continue
		}

		keyCopy := cloneNode(key)
		keyCopy.Value = prefixedPath
		valueCopy := cloneNode(value)
		ensureGroupIDParameter(valueCopy)
		additions = append(additions, keyCopy, valueCopy)
		removeKeys = append(removeKeys, key.Value)
	}

	for _, key := range removeKeys {
		deleteMappingKey(paths, key)
	}
	paths.Content = append(paths.Content, additions...)
	return nil
}

func shouldPrefixPath(apiPath string) bool {
	if apiPath == "" || apiPath == "/" || strings.HasPrefix(apiPath, groupPrefix+"/") {
		return false
	}

	excludedPrefixes := []string{
		"/admin/",
		"/master/",
		"/products/",
		"/health",
	}
	if slices.ContainsFunc(excludedPrefixes, func(prefix string) bool {
		return strings.HasPrefix(apiPath, prefix)
	}) {
		return false
	}

	return strings.HasPrefix(apiPath, "/lib/") ||
		strings.HasPrefix(apiPath, "/p/") ||
		strings.HasPrefix(apiPath, "/pipelines") ||
		strings.HasPrefix(apiPath, "/routes") ||
		strings.HasPrefix(apiPath, "/search/") ||
		strings.HasPrefix(apiPath, "/settings/") ||
		strings.HasPrefix(apiPath, "/notification") ||
		strings.HasPrefix(apiPath, "/collectors") ||
		strings.HasPrefix(apiPath, "/executors") ||
		strings.HasPrefix(apiPath, "/functions") ||
		strings.HasPrefix(apiPath, "/system/")
}

func ensureGroupIDParameter(pathItem *yaml.Node) {
	if pathItem.Kind != yaml.MappingNode {
		return
	}

	for index := 0; index < len(pathItem.Content); index += 2 {
		key := pathItem.Content[index]
		if !httpMethods[key.Value] {
			continue
		}
		operation := pathItem.Content[index+1]
		if operation.Kind == yaml.MappingNode {
			ensureOperationGroupIDParameter(operation)
		}
	}
}

func ensureOperationGroupIDParameter(operation *yaml.Node) {
	parameters, ok := mappingValue(operation, "parameters")
	if !ok {
		parameters = &yaml.Node{Kind: yaml.SequenceNode}
		setMappingValue(operation, scalar("parameters"), parameters)
	}
	if parameters.Kind != yaml.SequenceNode || hasNamedParameter(parameters, "groupId") {
		return
	}
	parameters.Content = append([]*yaml.Node{groupIDParameter()}, parameters.Content...)
}

func hasNamedParameter(parameters *yaml.Node, name string) bool {
	for _, parameter := range parameters.Content {
		if parameter.Kind != yaml.MappingNode {
			continue
		}
		nameNode, ok := mappingValue(parameter, "name")
		if ok && nameNode.Value == name {
			return true
		}
	}
	return false
}

func groupIDParameter() *yaml.Node {
	return mapping(
		"name", scalar("groupId"),
		"in", scalar("path"),
		"required", boolScalar(true),
		"schema", mapping("type", scalar("string")),
		"description", scalar("Worker group ID."),
	)
}

func unwrapCountedResponses(spec *yaml.Node) error {
	root := documentMapping(spec)
	paths, ok := mappingValue(root, "paths")
	if !ok || paths.Kind != yaml.MappingNode {
		return fmt.Errorf("spec paths mapping not found")
	}
	schemas, ok := lookupComponentsSchemas(root)
	if !ok {
		return fmt.Errorf("components schemas mapping not found")
	}

	for index := 0; index < len(paths.Content); index += 2 {
		path := paths.Content[index].Value
		pathItem := paths.Content[index+1]
		if pathItem.Kind != yaml.MappingNode {
			continue
		}
		for operationIndex := 0; operationIndex < len(pathItem.Content); operationIndex += 2 {
			method := pathItem.Content[operationIndex].Value
			if !httpMethods[method] {
				continue
			}
			unwrapOperationCountedResponses(path, method, pathItem.Content[operationIndex+1], schemas)
		}
	}
	return nil
}

func lookupComponentsSchemas(root *yaml.Node) (*yaml.Node, bool) {
	components, ok := mappingValue(root, "components")
	if !ok || components.Kind != yaml.MappingNode {
		return nil, false
	}
	schemas, ok := mappingValue(components, "schemas")
	if !ok || schemas.Kind != yaml.MappingNode {
		return nil, false
	}
	return schemas, true
}

func unwrapOperationCountedResponses(path, method string, operation, schemas *yaml.Node) {
	responses, ok := mappingValue(operation, "responses")
	if !ok || responses.Kind != yaml.MappingNode {
		return
	}

	for index := 0; index < len(responses.Content); index += 2 {
		response := responses.Content[index+1]
		content, ok := mappingValue(response, "content")
		if !ok || content.Kind != yaml.MappingNode {
			continue
		}
		for contentIndex := 0; contentIndex < len(content.Content); contentIndex += 2 {
			mediaType := content.Content[contentIndex+1]
			schema, ok := mappingValue(mediaType, "schema")
			if ok {
				unwrapCountedSchemaRef(mediaType, schema, schemas, countedResponseIsSingle(path, method, operation))
			}
		}
	}
}

func unwrapCountedSchemaRef(mediaType, schema, schemas *yaml.Node, single bool) {
	if schema.Kind != yaml.MappingNode {
		return
	}

	ref, ok := mappingValue(schema, "$ref")
	if !ok || ref.Kind != yaml.ScalarNode {
		return
	}

	const countedPrefix = "#/components/schemas/Counted"
	if !strings.HasPrefix(ref.Value, countedPrefix) {
		return
	}

	schemaName := strings.TrimPrefix(ref.Value, "#/components/schemas/")
	countedSchema, ok := mappingValue(schemas, schemaName)
	if !ok {
		return
	}
	itemsSchema, ok := countedItemsSchema(countedSchema)
	if !ok {
		return
	}
	if single {
		setMappingValue(mediaType, scalar("schema"), cloneNode(itemsSchema))
		return
	}
	setMappingValue(mediaType, scalar("schema"), arraySchema(itemsSchema))
}

func countedResponseIsSingle(path, method string, operation *yaml.Node) bool {
	for _, key := range []string{
		"x-terraform-resource",
		"x-terraform-create",
		"x-terraform-read",
		"x-terraform-update",
		"x-terraform-delete",
	} {
		if _, ok := mappingValue(operation, key); ok {
			return true
		}
	}

	if method != "get" {
		return true
	}

	lastSegment := path[strings.LastIndex(path, "/")+1:]
	return strings.HasPrefix(lastSegment, "{") && strings.HasSuffix(lastSegment, "}")
}

func countedItemsSchema(countedSchema *yaml.Node) (*yaml.Node, bool) {
	properties, ok := mappingValue(countedSchema, "properties")
	if !ok || properties.Kind != yaml.MappingNode {
		return nil, false
	}
	items, ok := mappingValue(properties, "items")
	if !ok || items.Kind != yaml.MappingNode {
		return nil, false
	}
	itemSchema, ok := mappingValue(items, "items")
	if !ok || itemSchema.Kind != yaml.MappingNode {
		return nil, false
	}
	return itemSchema, true
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

func setMappingValue(node, key, value *yaml.Node) {
	for index := 0; index < len(node.Content); index += 2 {
		if node.Content[index].Value == key.Value {
			node.Content[index+1] = value
			return
		}
	}
	node.Content = append(node.Content, key, value)
}

func deleteMappingKey(node *yaml.Node, key string) {
	if node == nil || node.Kind != yaml.MappingNode {
		return
	}
	for index := 0; index < len(node.Content); index += 2 {
		if node.Content[index].Value == key {
			node.Content = append(node.Content[:index], node.Content[index+2:]...)
			return
		}
	}
}

func mergeMapping(target, update *yaml.Node) {
	if target.Kind != yaml.MappingNode || update.Kind != yaml.MappingNode {
		*target = *cloneNode(update)
		return
	}

	for index := 0; index < len(update.Content); index += 2 {
		updateKey := update.Content[index]
		updateValue := update.Content[index+1]
		targetValue, ok := mappingValue(target, updateKey.Value)
		if ok && targetValue.Kind == yaml.MappingNode && updateValue.Kind == yaml.MappingNode {
			mergeMapping(targetValue, updateValue)
			continue
		}
		setMappingValue(target, cloneNode(updateKey), cloneNode(updateValue))
	}
}

func cloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	copy := *node
	copy.Content = make([]*yaml.Node, len(node.Content))
	for index, child := range node.Content {
		copy.Content[index] = cloneNode(child)
	}
	return &copy
}

func scalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func boolScalar(value bool) *yaml.Node {
	if value {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "true"}
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "false"}
}

func arraySchema(items *yaml.Node) *yaml.Node {
	return mapping(
		"type", scalar("array"),
		"items", cloneNode(items),
	)
}

func mapping(values ...any) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode}
	for index := 0; index < len(values); index += 2 {
		key, ok := values[index].(string)
		if !ok {
			panic("mapping key must be string")
		}
		value, ok := values[index+1].(*yaml.Node)
		if !ok {
			panic("mapping value must be *yaml.Node")
		}
		node.Content = append(node.Content, scalar(key), value)
	}
	return node
}
