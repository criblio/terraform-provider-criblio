package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestEdgeFleet(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: edgeFleetConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.my_edge_fleet", "id", "my-edge-fleet"),
						resource.TestCheckResourceAttr("criblio_group.my_edge_fleet", "name", "my-edge-fleet"),
						resource.TestCheckResourceAttr("criblio_group.my_edge_fleet", "product", "edge"),
					),
				},
				{
					Config: edgeFleetConfig,
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

const edgeFleetConfig = `
resource "criblio_group" "my_edge_fleet" {
  config_version        = data.criblio_config_version.my_configversion.id
  estimated_ingest_rate = 1024
  id                    = "my-edge-fleet"
  is_fleet              = true
  name                  = "my-edge-fleet"
  on_prem               = true
  product               = "edge"
  provisioned           = false
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

data "criblio_config_version" "my_configversion" {
  id         = "default"
}

output "edge_fleet" {
  value = criblio_group.my_edge_fleet
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider"
  version = "999.99.9"
}
`
