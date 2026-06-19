package parser

import (
	"path/filepath"
	"testing"
)

func TestParseCertificateResource(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	certificate := resourceByName(t, resources, "certificate")
	if certificate.Create.Path != "/m/{groupId}/system/certificates" {
		t.Fatalf("create path = %q", certificate.Create.Path)
	}
	if certificate.Read.Path != "/m/{groupId}/system/certificates/{id}" {
		t.Fatalf("read path = %q", certificate.Read.Path)
	}
	if len(certificate.Create.Examples) != 1 {
		t.Fatalf("create example count = %d", len(certificate.Create.Examples))
	}
	if certificate.Create.Examples[0].Name != "CertificateCreateExample" {
		t.Fatalf("create example name = %q", certificate.Create.Examples[0].Name)
	}

	groupID := fieldByTFName(t, certificate.Fields, "group_id")
	if !groupID.Required || !groupID.ForceNew || !groupID.PathParam {
		t.Fatalf("group_id flags = required:%v forceNew:%v pathParam:%v", groupID.Required, groupID.ForceNew, groupID.PathParam)
	}

	id := fieldByTFName(t, certificate.Fields, "id")
	if !id.Required || !id.ForceNew || !id.SuppressDiff {
		t.Fatalf("id flags = required:%v forceNew:%v suppressDiff:%v", id.Required, id.ForceNew, id.SuppressDiff)
	}

	privKey := fieldByTFName(t, certificate.Fields, "priv_key")
	if !privKey.Sensitive || !privKey.PreferState || privKey.ApplyStrategy != "stringFromAPIOrPrior" {
		t.Fatalf("priv_key flags = sensitive:%v prefer:%v strategy:%q", privKey.Sensitive, privKey.PreferState, privKey.ApplyStrategy)
	}

	passphrase := fieldByTFName(t, certificate.Fields, "passphrase")
	if !passphrase.Sensitive || !passphrase.PreferState || passphrase.ApplyStrategy != "stringFromAPIOrPrior" {
		t.Fatalf("passphrase flags = sensitive:%v prefer:%v strategy:%q", passphrase.Sensitive, passphrase.PreferState, passphrase.ApplyStrategy)
	}

	expiry := fieldByTFName(t, certificate.Fields, "cert_expiry_date")
	if !expiry.Computed || expiry.Optional || expiry.Required {
		t.Fatalf("cert_expiry_date flags = computed:%v optional:%v required:%v", expiry.Computed, expiry.Optional, expiry.Required)
	}

	inUse := fieldByTFName(t, certificate.Fields, "in_use")
	if !inUse.Computed || inUse.Type != "array" {
		t.Fatalf("in_use computed/type = %v/%q", inUse.Computed, inUse.Type)
	}

	args := fieldByTFName(t, certificate.Fields, "args")
	if args.Type != "array" || args.ElementType != "object" || args.ApplyStrategy != "preferState" || !args.SuppressDiff {
		t.Fatalf("args type/element/strategy/suppressDiff = %q/%q/%q/%v", args.Type, args.ElementType, args.ApplyStrategy, args.SuppressDiff)
	}
	if args.NestedModelName != "CertificateArgsModel" || args.NestedAPIModelName != "CertificateArgsAPIModel" {
		t.Fatalf("args nested model names = %q/%q", args.NestedModelName, args.NestedAPIModelName)
	}
	if len(args.Fields) != 2 {
		t.Fatalf("args nested field count = %d", len(args.Fields))
	}
	argName := fieldByTFName(t, args.Fields, "name")
	if !argName.Required || argName.Type != "string" || !argName.SuppressDiff {
		t.Fatalf("args.name required/type/suppressDiff = %v/%q/%v", argName.Required, argName.Type, argName.SuppressDiff)
	}

	renamed := fieldByTFName(t, certificate.Fields, "display_name")
	if renamed.APIName != "displayName" {
		t.Fatalf("renamed APIName = %q", renamed.APIName)
	}

	custom := fieldByTFName(t, certificate.Fields, "conf")
	if custom.CustomType != "jsontypes.NormalizedType{}" {
		t.Fatalf("custom type = %q", custom.CustomType)
	}

	if hasField(certificate.Fields, "ignored") {
		t.Fatalf("ignored field should not be present")
	}
}

func TestParseOneOfVariants(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	destination := resourceByName(t, resources, "destination")
	if len(destination.OneOfVariants) != 3 {
		t.Fatalf("variant count = %d", len(destination.OneOfVariants))
	}
	azure := variantByName(t, destination.OneOfVariants, "OutputAzureBlob")
	accountKey := fieldByTFName(t, azure.Fields, "account_key")
	if accountKey.ApplyStrategy != "stringFromAPIOrPrior" {
		t.Fatalf("account_key strategy = %q", accountKey.ApplyStrategy)
	}
}

func TestParseMappingRulesetBackwardCompatibleDefaults(t *testing.T) {
	resources, err := ParseFile(filepath.Join("..", "testdata", "fixture.yml"))
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	mappingRuleset := resourceByName(t, resources, "mapping_ruleset")
	id := fieldByTFName(t, mappingRuleset.Fields, "id")
	if id.Required || !id.Optional || !id.Computed {
		t.Fatalf("mapping ruleset id flags = required:%v optional:%v computed:%v", id.Required, id.Optional, id.Computed)
	}

	conf := fieldByTFName(t, mappingRuleset.Fields, "conf")
	functions := fieldByTFName(t, conf.Fields, "functions")
	functionID := fieldByTFName(t, functions.Fields, "id")
	if functionID.Required || !functionID.Optional || functionID.Computed {
		t.Fatalf("mapping function id flags = required:%v optional:%v computed:%v", functionID.Required, functionID.Optional, functionID.Computed)
	}
	final := fieldByTFName(t, functions.Fields, "final")
	if final.Required || !final.Optional || final.Computed {
		t.Fatalf("mapping function final flags = required:%v optional:%v computed:%v", final.Required, final.Optional, final.Computed)
	}
}

func resourceByName(t *testing.T, resources []ResourceDef, name string) ResourceDef {
	t.Helper()
	for _, resource := range resources {
		if resource.Name == name {
			return resource
		}
	}
	t.Fatalf("resource %q not found", name)
	return ResourceDef{}
}

func fieldByTFName(t *testing.T, fields []FieldDef, name string) FieldDef {
	t.Helper()
	for _, field := range fields {
		if field.TerraformName == name {
			return field
		}
	}
	t.Fatalf("field %q not found", name)
	return FieldDef{}
}

func variantByName(t *testing.T, variants []OneOfVariantDef, name string) OneOfVariantDef {
	t.Helper()
	for _, variant := range variants {
		if variant.APIName == name {
			return variant
		}
	}
	t.Fatalf("variant %q not found", name)
	return OneOfVariantDef{}
}

func hasField(fields []FieldDef, name string) bool {
	for _, field := range fields {
		if field.TerraformName == name {
			return true
		}
	}
	return false
}
