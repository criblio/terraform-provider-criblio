package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestSubscription(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: subConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription", "id", "my_subscription"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription", "description", "test subscription"),
					),
				},
				{
					Config: subConfig,
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

var subConfig = `

resource "criblio_subscription" "my_subscription" {
  description = "test subscription"
  disabled    = true
  filter      = "test"
  group_id    = "default"
  id          = "my_subscription"
  pipeline    = "passthru"
}

output "subscription" {
  value = criblio_subscription.my_subscription
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}

`
