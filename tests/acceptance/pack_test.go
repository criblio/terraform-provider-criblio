package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPack(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack source test for on-prem: uses HelloPacks test fixture")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	t.Run("source", func(t *testing.T) {
		packID := "test-pack-source-" + suffix
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: packSourceConfig(packID, "created"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.source", "id", packID),
						resource.TestCheckResourceAttr("criblio_pack.source", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.source", "description", "Pack created"),
						resource.TestCheckResourceAttr("criblio_pack.source", "display_name", "Pack created"),
						resource.TestCheckResourceAttr("criblio_pack.source", "version", "1.0.0"),
						resource.TestCheckResourceAttr("criblio_pack.source", "source", "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"),
					),
				},
				{
					Config: packSourceConfig(packID, "updated"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.source", "description", "Pack updated"),
						resource.TestCheckResourceAttr("criblio_pack.source", "display_name", "Pack updated"),
						resource.TestCheckResourceAttr("criblio_pack.source", "version", "1.0.1"),
					),
				},
				{
					Config: packSourceConfig(packID, "updated"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
				{
					ResourceName:      "criblio_pack.source",
					ImportState:       true,
					ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q}`, packID),
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"description",
						"disabled",
						"display_name",
						"filename",
						"version",
					},
				},
			},
		})
	})

	t.Run("file", func(t *testing.T) {
		packFile, err := filepath.Abs(filepath.Join("testdata", "TestPackFromFile", "plan-diff", "cribl-palo-alto-networks-source-1.0.0.crbl"))
		if err != nil {
			t.Fatalf("resolve pack file path: %v", err)
		}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: packFileConfig(packFile),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.file", "id", "pack-from-file"),
						resource.TestCheckResourceAttr("criblio_pack.file", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.file", "description", "Pack from file"),
						resource.TestCheckResourceAttr("criblio_pack.file", "display_name", "Pack from file"),
						resource.TestCheckResourceAttr("criblio_pack.file", "version", "1.0.0"),
					),
				},
				{
					Config: packFileConfig(packFile),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	})

	t.Run("metadata", func(t *testing.T) {
		packID := "test-pack-metadata-" + suffix
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: packMetadataConfig(packID, "created", "1.0.0"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.metadata", "id", packID),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "description", "Pack metadata created"),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "display_name", "Pack metadata created"),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "version", "1.0.0"),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "author", "Observability Team"),
					),
				},
				{
					Config: packMetadataConfig(packID, "updated", "1.0.1"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.metadata", "description", "Pack metadata updated"),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "display_name", "Pack metadata updated"),
						resource.TestCheckResourceAttr("criblio_pack.metadata", "version", "1.0.1"),
					),
				},
				{
					Config: packMetadataConfig(packID, "updated", "1.0.1"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
				{
					ResourceName:      "criblio_pack.metadata",
					ImportState:       true,
					ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q}`, packID),
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"disabled",
						"filename",
					},
				},
			},
		})
	})
}

func packSourceConfig(id, suffix string) string {
	version := "1.0.0"
	if suffix == "updated" {
		version = "1.0.1"
	}
	return fmt.Sprintf(`resource "criblio_pack" "source" {
  id           = %q
  group_id     = "default"
  description  = "Pack %s"
  display_name = "Pack %s"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = %q
}
`, id, suffix, suffix, version)
}

func packFileConfig(filename string) string {
	return fmt.Sprintf(`resource "criblio_pack" "file" {
  id           = "pack-from-file"
  group_id     = "default"
  description  = "Pack from file"
  disabled     = true
  display_name = "Pack from file"
  filename     = %q
  version      = "1.0.0"
}
`, filename)
}

func packMetadataConfig(id, suffix, version string) string {
	return fmt.Sprintf(`resource "criblio_pack" "metadata" {
  id           = %q
  group_id     = "default"
  description  = "Pack metadata %s"
  disabled     = false
  display_name = "Pack metadata %s"
  version      = %q
  author       = "Observability Team"
}
`, id, suffix, suffix, version)
}
