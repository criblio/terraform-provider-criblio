package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestGlobalVar(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: gVarConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "id", "default"),
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
  description = "Current epoch time"
  group_id    = "default"
  id          = "default"
  lib         = "cribl"
  tags        = "cribl,sample"
  type        = "expression"
  value       = "Math.floor(Date.now()/1000)"
}

output "global_var" {
  value = criblio_global_var.my_globalvar
}

data "criblio_global_var" "my_globalvar" {
  group_id = "default"
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}

`
