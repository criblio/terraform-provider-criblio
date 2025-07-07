package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackFromFile(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: pFileConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "id", "pack-from-file"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "description", "Pack from file"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "display_name", "Pack from file"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "version", "1.0.0"),
					),
				},
				{
					Config: pFileConfig,
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

var pFileConfig = `

resource "criblio_pack" "my_pack" {
  id           = "pack-from-file"
  group_id     = "default"
  description  = "Pack from file"
  disabled     = true
  display_name = "Pack from file"
  filename     = "cribl-palo-alto-networks-source-1.0.0.crbl"
  version      = "1.0.0"

}

# Output the pack details to see the read-only attributes
output "pack_details" {
  value = criblio_pack.my_pack
}

data "criblio_pack" "my_pack" {
  group_id = "default"
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}
`
