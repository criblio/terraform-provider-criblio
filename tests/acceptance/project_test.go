package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestProject(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: projConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_project.my_project", "id", "my_project"),
						resource.TestCheckResourceAttr("criblio_project.my_project", "group_id", "default"),
					),
				},
				{
					Config: projConfig,
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

var projConfig = `

resource "criblio_project" "my_project" {
  consumers = {
  }
  description = "test project"
  destinations = [
  ]
  group_id = "default"
  id       = "my_project"
  subscriptions = [
  ]
}

output "project" {
  value = criblio_project.my_project
}

data "criblio_project" "my_project" {
  group_id = "...my_group_id..."
  id       = "...my_id..."
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
}
`
