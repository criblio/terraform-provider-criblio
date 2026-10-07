package main

import (
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestPackCollectorSharesVariantImplementations(t *testing.T) {
	resource := parser.ResourceDef{
		Name: "pack_collector", FileStem: "pack_collector",
		TypeName: "criblio_pack_collector", StructName: "PackCollector",
		OneOfVariants: []parser.OneOfVariantDef{{
			GoName: "InputCollectorGCS", TerraformName: "input_collector_gcs",
			ModelName: "InputCollectorGCSModel", DiscriminatorValue: "gcs",
		}},
	}
	got := renderTemplate(t, "types", resource)
	assertContains(t, got, `output["type"] = "collection"`)
	assertContains(t, got, `normalizeCollectorUnionValues(raw)`)
	if strings.Contains(got, "type InputCollectorGCSModel struct") || strings.Contains(got, "func normalizeCollectorUnionValues(") {
		t.Fatal("pack collector must reuse collector variant types and normalization")
	}
	aliases := discriminatorCaseValues(resource, resource.OneOfVariants[0])
	if len(aliases) != 2 || aliases[1] != "google_cloud_storage" {
		t.Fatalf("collector discriminator aliases lost: %v", aliases)
	}
	got = renderTemplate(t, "resource", resource)
	assertContains(t, got, "oneOfRequestModelWithHoistedIdentity(model)")
	if strings.Contains(got, "ResourceWithUpgradeState") {
		t.Fatal("new pack collector resource must not inherit legacy collector state upgrades")
	}
}
