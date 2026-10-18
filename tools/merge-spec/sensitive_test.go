package main

import (
	"path/filepath"
	"testing"
)

func TestCredentialOverlayAnnotations(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(filepath.Join("..", ".."))
	output := filepath.Join(dir, "merged-spec.yml")
	if err := run("upstream-openapi.yml", "terraform-overlay.yml", output, filepath.Join(dir, "availability.go")); err != nil {
		t.Fatal(err)
	}
	document, err := readYAML(output)
	if err != nil {
		t.Fatal(err)
	}
	root := documentMapping(document)
	for _, path := range []string{
		"InputHttpRaw.properties.authTokens",
		"InputHttpRaw.properties.authTokensExt.items.properties.token",
		"InputHttp.properties.authTokens",
		"AuthTypeTemplatemanualApiKeyAuthType.properties.manualAPIKey",
		"OutputGoogleChronicle.properties.serviceAccountCredentials",
		"InputCollectorRest.allOf.1.properties.collector.properties.conf.properties.password",
		"InputCollectorRest.allOf.1.properties.collector.properties.conf.properties.token",
		"InputCollectorS3.allOf.1.properties.collector.properties.conf.properties.awsSecretKey",
		"WebhookTarget.properties.password",
		"SlackTarget.properties.url",
	} {
		field, err := lookupTarget(root, "$.components.schemas."+path)
		if err != nil {
			t.Fatal(err)
		}
		annotation, ok := mappingValue(field, "x-terraform-sensitive")
		if !ok || annotation.Value != "true" {
			t.Errorf("%s must be annotated sensitive", path)
		}
		if preference, ok := mappingValue(field, "x-terraform-prefer-state"); ok && preference.Value == "true" {
			t.Errorf("%s sensitivity must not change state merging", path)
		}
	}
	for _, path := range []string{
		"InputHttpRaw.properties.authTokensExt.items.properties.tokenSecret",
		"InputCollectorRest.allOf.1.properties.collector.properties.conf.properties.credentialsSecret",
		"InputCollectorS3.allOf.1.properties.collector.properties.conf.properties.awsSecret",
	} {
		field, err := lookupTarget(root, "$.components.schemas."+path)
		if err != nil {
			t.Fatal(err)
		}
		if annotation, ok := mappingValue(field, "x-terraform-sensitive"); ok && annotation.Value == "true" {
			t.Errorf("stored secret reference %s must remain non-sensitive", path)
		}
	}
}
