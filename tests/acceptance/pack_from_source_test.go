package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackFromSource(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: pSourceConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "id", "pack-from-source"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "description", "Pack from source"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "display_name", "Pack from source"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "version", "1.0.0"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "source",
							"file:/opt/cribl_data/failover/groups/default/default/HelloPacks"),
					),
				},
				{
					Config: pSourceConfig,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	})
}

var pSourceConfig = `
resource "criblio_pack" "my_pack" {
  id           = "pack-from-source"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"

}

# Output the pack details to see the read-only attributes
output "pack_details" {
  value = criblio_pack.my_pack
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
}

`
