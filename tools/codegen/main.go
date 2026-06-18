// Command codegen renders Terraform provider scaffolds from merged OpenAPI.
package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func main() {
	specPath := flag.String("spec", "merged-spec.yml", "merged OpenAPI YAML file")
	ignorePath := flag.String("ignore", ".codegen-ignore", "codegen ignore file")
	outputDir := flag.String("output-dir", "", "optional output directory prefix")
	resourceName := flag.String("resource", "", "optional resource name filter")
	flag.Parse()

	if err := run(*specPath, *ignorePath, *outputDir, *resourceName); err != nil {
		fmt.Fprintf(os.Stderr, "codegen: %v\n", err)
		os.Exit(1)
	}
}

func run(specPath, ignorePath, outputDir, resourceName string) error {
	resources, err := parser.ParseFile(specPath)
	if err != nil {
		return err
	}
	if resourceName != "" {
		resources = slices.DeleteFunc(resources, func(resource parser.ResourceDef) bool {
			return !resourceMatches(resource, resourceName)
		})
	}

	ignored, err := readIgnoreFile(ignorePath)
	if err != nil {
		return err
	}
	files, err := newRenderer(outputDir, ignored).render(resources)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Skipped {
			fmt.Printf("skip %s\n", file.Path)
			continue
		}
		fmt.Printf("write %s\n", file.Path)
	}
	return nil
}

func resourceMatches(resource parser.ResourceDef, name string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(name, "_", ""))
	for _, candidate := range []string{
		resource.Name,
		resource.FileStem,
		strings.TrimPrefix(resource.TypeName, "criblio_"),
	} {
		if strings.EqualFold(candidate, name) {
			return true
		}
		if strings.ToLower(strings.ReplaceAll(candidate, "_", "")) == normalized {
			return true
		}
	}
	return false
}
