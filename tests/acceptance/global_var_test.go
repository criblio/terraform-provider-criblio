package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestGlobalVar(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: gVarConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "id", "my_globalvar"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "group_id", "default"),
					),
				},
				{
					Config: gVarConfig,
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

var gVarConfig = `

resource "criblio_global_var" "my_globalvar" {
  description = "test"
  group_id    = "default"
  id          = "my_globalvar"
  lib         = "test"
  tags        = "test"
  type        = "number"
  value       = 100
}

output "global_var" {
  value = criblio_global_var.my_globalvar
}

/*
data "criblio_global_var" "my_globalvar" {
  group_id = "...my_group_id..."
  with     = "...my_with..."
}
*/

 provider "criblio" {
    server_url = "https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
 }
`
