package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestWorkerGroup(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: wgConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.my_group", "id", "my-group"),
						resource.TestCheckResourceAttr("criblio_group.my_group", "name", "my-group"),
					),
				},
				{
					Config: wgConfig,
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

var wgConfig = `
resource "criblio_group" "my_group" {
  cloud = {
    provider = "azure"
    region   = "eastus"
  }
  estimated_ingest_rate = 1024
  id                    = "my-group"
  is_fleet              = false
  name                  = "my-group"
  on_prem               = false
  product               = "stream"
  provisioned           = false
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

output "group" {
  value = criblio_group.my_group
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}
`
