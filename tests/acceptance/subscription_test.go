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
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription_with_enabled", "id", "my_subscription_with_enabled"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription_with_enabled", "description", "test subscription with enabled"),
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

resource "criblio_subscription" "my_subscription_with_enabled" {
  description = "test subscription with enabled"
  disabled    = false
  filter      = "test"
  group_id    = "default"
  id          = "my_subscription_with_enabled"
  pipeline    = "passthru"
}

output "subscription" {
  value = criblio_subscription.my_subscription
}

`
