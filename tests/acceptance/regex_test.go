package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestRegex(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: regexConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "description", "test"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "lib", "custom"),
					),
				},
				{
					Config: regexConfig,
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

var regexConfig = `
resource "criblio_regex" "my_regex" {
  description = "test"
  group_id    = "default"
  id          = "test_regex"
  lib         = "custom"
  regex       = "/\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b/"
  sample_data = "192.168.1.1, 10.0.0.1, 172.16.0.1"
  tags        = "test"
}

/*
output "regex" {
  value = criblio_regex.my_regex
}

data "criblio_regex" "my_regex" {
  group_id = "...my_group_id..."
}
*/

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
}
`
