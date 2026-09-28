package provider

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderDescriptionHasNoPreviewDisclaimer(t *testing.T) {
	t.Parallel()

	var response provider.SchemaResponse
	(&CriblioProvider{}).Schema(context.Background(), provider.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("provider schema diagnostics: %v", response.Diagnostics)
	}
	if response.Schema.MarkdownDescription == "" {
		t.Fatal("provider schema markdown description is empty")
	}

	descriptions := map[string]string{
		"schema_description":          response.Schema.Description,
		"schema_markdown_description": response.Schema.MarkdownDescription,
	}
	for _, filename := range []string{"docs/index.md", "openapi.yml"} {
		content, err := os.ReadFile(filepath.Join("..", "..", filename))
		if err != nil {
			t.Fatalf("read %s: %v", filename, err)
		}
		descriptions[filename] = string(content)
	}

	for name, description := range descriptions {
		t.Run(name, func(t *testing.T) {
			for _, disclaimer := range []string{
				"preview feature",
				"we do not recommend using it in a production environment",
			} {
				if strings.Contains(strings.ToLower(description), disclaimer) {
					t.Errorf("provider description contains %q", disclaimer)
				}
			}
		})
	}
}
