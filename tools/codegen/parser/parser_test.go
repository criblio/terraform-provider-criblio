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

	groupID := fieldByTFName(t, certificate.Fields, "group_id")
	if !groupID.Required || !groupID.ForceNew || !groupID.PathParam {
		t.Fatalf("group_id flags = required:%v forceNew:%v pathParam:%v", groupID.Required, groupID.ForceNew, groupID.PathParam)
	}

	id := fieldByTFName(t, certificate.Fields, "id")
	if !id.Required || !id.ForceNew {
		t.Fatalf("id flags = required:%v forceNew:%v", id.Required, id.ForceNew)
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
