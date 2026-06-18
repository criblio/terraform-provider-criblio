package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestRendererHonorsCodegenIgnore(t *testing.T) {
	resources := parseFixture(t)
	ignored, err := readIgnoreFile(filepath.Join("testdata", ".codegen-ignore"))
	if err != nil {
		t.Fatalf("readIgnoreFile returned error: %v", err)
	}
	dir := t.TempDir()
	files, err := newRenderer(dir, ignored).render([]parser.ResourceDef{resourceByName(t, resources, "certificate")})
	if err != nil {
		t.Fatalf("render returned error: %v", err)
	}

	if !hasSkipped(files, filepath.Join(dir, "internal/provider/certificate_resource.go")) {
		t.Fatalf("resource output was not skipped")
	}
	if _, err := os.Stat(filepath.Join(dir, "internal/provider/certificate_resource.go")); !os.IsNotExist(err) {
		t.Fatalf("ignored resource file exists or stat returned unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "internal/provider/certificate_types.go")); err != nil {
		t.Fatalf("types output missing: %v", err)
	}
}

func TestRenderedSnippets(t *testing.T) {
	resources := parseFixture(t)
	certificate := resourceByName(t, resources, "certificate")
	resourceContent := renderTemplate(t, "resource", certificate)
	assertContains(t, resourceContent, "applyCertificateAPIToState(apiModel, &model, true, false)")
	assertContains(t, resourceContent, "applyCertificateAPIToState(apiModel, &model, true, isCertificateImportState(&model))")
	assertContains(t, resourceContent, "if !preserveInputs || (fillMissingInputs && (state.Cert.IsNull() || state.Cert.IsUnknown()))")
	assertContains(t, resourceContent, "api.DisplayName.IsNull()")
	assertContains(t, resourceContent, "if !api.InUse.IsNull() && !api.InUse.IsUnknown()")
	assertContains(t, resourceContent, "stringFromAPIOrPrior(api.Passphrase.ValueString(), state.Passphrase)")
	assertContains(t, resourceContent, "stringplanmodifier.RequiresReplaceIfConfigured()")
	assertContains(t, resourceContent, "custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress)")
	assertContains(t, resourceContent, "custom_listplanmodifier.SuppressDiff(custom_listplanmodifier.ExplicitSuppress)")
	assertContains(t, resourceContent, "custom_objectplanmodifier.SuppressDiff(custom_objectplanmodifier.ExplicitSuppress)")
	assertContains(t, resourceContent, "state.InUse = types.ListValueMust(types.StringType, nil)")
	assertContains(t, resourceContent, "clients, ok := req.ProviderData.(*ProviderClients)")
	assertContains(t, resourceContent, "r.client = clients.RC")
	assertContains(t, resourceContent, `json:"group_id"`)
	assertContains(t, resourceContent, `path.Root("group_id")`)
	assertNotContains(t, resourceContent, "speakeasy_")
	assertNotContains(t, resourceContent, "internal/sdk")

	typesContent := renderTemplate(t, "types", certificate)
	assertContains(t, typesContent, "Conf jsontypes.Normalized")
	assertContains(t, typesContent, "InUse types.List")
	assertContains(t, typesContent, "InUse []types.String")
	assertContains(t, typesContent, "func (m CertificateModel) MarshalJSON()")
	assertContains(t, typesContent, "types.ListValueFrom(context.Background(), types.StringType, input.InUse)")

	dataSourceContent := renderTemplate(t, "data_source", certificate)
	assertContains(t, dataSourceContent, "applyCertificateAPIToState(apiModel, &model, false, false)")

	destination := resourceByName(t, resources, "destination")
	destinationTypes := renderTemplate(t, "types", destination)
	assertContains(t, destinationTypes, "OutputAzureBlob *OutputAzureBlobModel")
	assertContains(t, destinationTypes, "OutputElastic *OutputElasticModel")
	assertContains(t, destinationTypes, "OutputS3 *OutputS3Model")

	destinationResource := renderTemplate(t, "resource", destination)
	assertContains(t, destinationResource, "if api.OutputAzureBlob != nil && (!preserveInputs || (fillMissingInputs && state.OutputAzureBlob == nil))")
	assertContains(t, destinationResource, "state.OutputAzureBlob = &OutputAzureBlobModel{}")
	assertContains(t, destinationResource, "stringFromAPIOrPrior(api.OutputAzureBlob.AccountKey.ValueString(), state.OutputAzureBlob.AccountKey)")
}

func parseFixture(t *testing.T) []parser.ResourceDef {
	t.Helper()
	resources, err := parser.ParseFile(filepath.Join("testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}
	return resources
}

func resourceByName(t *testing.T, resources []parser.ResourceDef, name string) parser.ResourceDef {
	t.Helper()
	for _, resource := range resources {
		if resource.Name == name {
			return resource
		}
	}
	t.Fatalf("resource %q not found", name)
	return parser.ResourceDef{}
}

func renderTemplate(t *testing.T, kind string, resource parser.ResourceDef) string {
	t.Helper()
	content, err := executeTemplate(kind, resource)
	if err != nil {
		t.Fatalf("executeTemplate returned error: %v", err)
	}
	return string(content)
}

func hasSkipped(files []renderedFile, path string) bool {
	for _, file := range files {
		if file.Path == path && file.Skipped {
			return true
		}
	}
	return false
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()
	if !strings.Contains(content, want) {
		t.Fatalf("expected content to contain %q", want)
	}
}

func assertNotContains(t *testing.T, content, want string) {
	t.Helper()
	if strings.Contains(content, want) {
		t.Fatalf("expected content not to contain %q", want)
	}
}
