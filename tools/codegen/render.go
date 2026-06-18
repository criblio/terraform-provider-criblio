package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

type renderer struct {
	outputDir string
	ignored   ignoreSet
}

type renderedFile struct {
	Path    string
	Skipped bool
}

func newRenderer(outputDir string, ignored ignoreSet) renderer {
	return renderer{
		outputDir: outputDir,
		ignored:   ignored,
	}
}

func (r renderer) render(resources []parser.ResourceDef) ([]renderedFile, error) {
	var files []renderedFile
	for _, resource := range resources {
		for _, output := range outputFiles(resource) {
			path := output.Path
			if r.outputDir != "" {
				path = filepath.Join(r.outputDir, path)
			}
			if r.ignored.ignored(output.Path) || r.ignored.ignored(path) {
				files = append(files, renderedFile{Path: path, Skipped: true})
				continue
			}
			content, err := executeTemplate(output.Kind, resource)
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return nil, fmt.Errorf("create output directory: %v", err)
			}
			if err := os.WriteFile(path, content, 0644); err != nil {
				return nil, fmt.Errorf("write %s: %v", path, err)
			}
			files = append(files, renderedFile{Path: path})
		}
	}
	return files, nil
}

func outputFiles(resource parser.ResourceDef) []parser.OutputFile {
	prefix := "internal/provider/" + resource.Name
	return []parser.OutputFile{
		{Path: prefix + "_types.go", Kind: "types"},
		{Path: prefix + "_client.go", Kind: "client"},
		{Path: prefix + "_resource.go", Kind: "resource"},
		{Path: prefix + "_data_source.go", Kind: "data_source"},
		{Path: "docs/resources/" + resource.Name + ".md", Kind: "doc"},
		{Path: "tests/acceptance/" + resource.Name + "_test.go", Kind: "test"},
	}
}

func executeTemplate(kind string, resource parser.ResourceDef) ([]byte, error) {
	body, ok := templateBodies[kind]
	if !ok {
		return nil, fmt.Errorf("template %q not found", kind)
	}
	tmpl, err := template.New(kind).Funcs(template.FuncMap{
		"goType":          goType,
		"schemaAttribute": schemaAttribute,
		"schemaTypeName":  schemaTypeName,
		"schemaSections":  schemaSections,
		"zeroValue":       zeroValue,
		"pathExpr":        pathExpr,
		"jsonName":        jsonName,
		"apiType":         apiType,
		"legacyGoType":    legacyGoType,
		"resourceType":    resourceType,
		"exampleUsage":    exampleUsage,
		"importBlock":     importBlock,
		"importCommand":   importCommand,
		"joinDocFields":   joinDocFields,
		"pathParamFields": pathParamFields,
		"jsonImport":      jsonImport,
	}).Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %v", kind, err)
	}
	var output bytes.Buffer
	if err := tmpl.Execute(&output, resource); err != nil {
		return nil, fmt.Errorf("execute template %q: %v", kind, err)
	}
	return output.Bytes(), nil
}

func joinDocFields(fields []parser.FieldDef) string {
	var lines []string
	for _, field := range fields {
		line := fmt.Sprintf("- `%s` (%s", field.TerraformName, schemaTypeName(field))
		if field.Sensitive {
			line += ", Sensitive"
		}
		line += ")"
		if field.Description != "" {
			line += " " + field.Description
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

type schemaDocSections struct {
	Required []parser.FieldDef
	Optional []parser.FieldDef
	Computed []parser.FieldDef
}

func schemaSections(resource parser.ResourceDef) schemaDocSections {
	var sections schemaDocSections
	for _, field := range resource.Fields {
		switch {
		case field.Required:
			sections.Required = append(sections.Required, field)
		case field.Computed && !field.Optional && !field.PathParam && !field.ForceNew:
			sections.Computed = append(sections.Computed, field)
		default:
			sections.Optional = append(sections.Optional, field)
		}
	}
	return sections
}

func goType(field parser.FieldDef) string {
	if field.CustomType == "jsontypes.NormalizedType{}" {
		return "jsontypes.Normalized"
	}
	switch field.Type {
	case "boolean":
		return "types.Bool"
	case "integer":
		return "types.Int64"
	case "number":
		return "types.Float64"
	case "array":
		return "types.List"
	case "object":
		return "types.Map"
	default:
		return "types.String"
	}
}

func legacyGoType(field parser.FieldDef) string {
	if field.CustomType == "jsontypes.NormalizedType{}" {
		return "jsontypes.Normalized"
	}
	switch field.Type {
	case "boolean":
		return "types.Bool"
	case "integer":
		return "types.Int64"
	case "number":
		return "types.Float64"
	case "array":
		return "[]types.String"
	case "object":
		return "types.Map"
	default:
		return "types.String"
	}
}

func apiType(field parser.FieldDef) string {
	switch field.Type {
	case "boolean":
		return "*bool"
	case "integer":
		return "*int64"
	case "number":
		return "*float64"
	case "array":
		return "[]string"
	case "object":
		return "map[string]string"
	default:
		return "*string"
	}
}

func schemaAttribute(field parser.FieldDef) string {
	switch field.Type {
	case "boolean":
		return "schema.BoolAttribute"
	case "integer":
		return "schema.Int64Attribute"
	case "number":
		return "schema.Float64Attribute"
	case "array":
		return "schema.ListAttribute"
	case "object":
		return "schema.MapAttribute"
	default:
		return "schema.StringAttribute"
	}
}

func schemaTypeName(field parser.FieldDef) string {
	switch field.Type {
	case "boolean":
		return "Boolean"
	case "integer":
		return "Integer"
	case "number":
		return "Number"
	case "array":
		return "List of String"
	case "object":
		return "Map of String"
	default:
		return "String"
	}
}

func zeroValue(field parser.FieldDef) string {
	switch goType(field) {
	case "types.Bool":
		return "types.BoolValue(false)"
	case "types.Int64":
		return "types.Int64Value(0)"
	case "types.Float64":
		return "types.Float64Value(0)"
	case "types.List":
		return "types.ListValueMust(types.StringType, nil)"
	case "types.Map":
		return "types.MapValueMust(types.StringType, nil)"
	case "jsontypes.Normalized":
		return `jsontypes.NewNormalizedValue("")`
	default:
		return `types.StringValue("")`
	}
}

func jsonName(field parser.FieldDef) string {
	return field.APIName + ",omitempty"
}

func resourceType(resource parser.ResourceDef) string {
	return "criblio_" + resource.Name
}

func exampleUsage(resource parser.ResourceDef) string {
	path := filepath.Join("examples", "resources", resourceType(resource), "resource.tf")
	content, err := os.ReadFile(path)
	if err != nil {
		return generatedExample(resource)
	}
	return strings.TrimSpace(string(content))
}

func importBlock(resource parser.ResourceDef) string {
	path := filepath.Join("examples", "resources", resourceType(resource), "import-by-string-id.tf")
	content, err := os.ReadFile(path)
	if err != nil {
		return generatedImportBlock(resource)
	}
	return strings.TrimSpace(string(content))
}

func importCommand(resource parser.ResourceDef) string {
	path := filepath.Join("examples", "resources", resourceType(resource), "import.sh")
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("terraform import %s.my_%s '{\"group_id\": \"default\", \"id\": \"cert-001\"}'", resourceType(resource), resourceType(resource))
	}
	return strings.TrimSpace(string(content))
}

func pathParamFields(resource parser.ResourceDef) []parser.FieldDef {
	var fields []parser.FieldDef
	for _, field := range resource.Fields {
		if field.PathParam {
			fields = append(fields, field)
		}
	}
	return fields
}

func jsonImport(resource parser.ResourceDef) bool {
	return len(pathParamFields(resource)) > 1
}

func generatedExample(resource parser.ResourceDef) string {
	var output strings.Builder
	fmt.Fprintf(&output, "resource %q %q {\n", resourceType(resource), "example")
	for _, field := range resource.Fields {
		if field.Computed && !field.Optional {
			continue
		}
		fmt.Fprintf(&output, "  %s = %q\n", field.TerraformName, exampleValue(field))
	}
	output.WriteString("}")
	return output.String()
}

func generatedImportBlock(resource parser.ResourceDef) string {
	return fmt.Sprintf(`import {
  to = %s.my_%s
  id = jsonencode({
    group_id = "default"
    id       = "cert-001"
  })
}`, resourceType(resource), resourceType(resource))
}

func exampleValue(field parser.FieldDef) string {
	if field.PathParam && field.TerraformName == "group_id" {
		return "default"
	}
	if field.TerraformName == "id" {
		return "example"
	}
	return "example"
}

func pathExpr(op parser.OperationDef) string {
	path := op.Path
	for _, param := range op.PathParams {
		path = strings.ReplaceAll(path, "{"+param.APIName+"}", `%s`)
	}
	if len(op.PathParams) == 0 {
		return fmt.Sprintf("%q", path)
	}
	args := []string{fmt.Sprintf("%q", path)}
	for _, param := range op.PathParams {
		args = append(args, "model."+param.GoName+".ValueString()")
	}
	return "fmt.Sprintf(" + strings.Join(args, ", ") + ")"
}
