package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSubscription(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is 'prohibited by current license'")
	}

	resourceName := "criblio_subscription.my_subscription"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: subscriptionConfig("test subscription", true, "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "phase2_lifecycle_subscription"),
					resource.TestCheckResourceAttr(resourceName, "description", "test subscription"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "filter", "test"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "pipeline", "passthru"),
				),
			},
			{
				Config: subscriptionConfig("updated subscription", false, "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated subscription"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "filter", "updated"),
				),
			},
			{
				Config:   subscriptionConfig("updated subscription", false, "updated"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default","id":"phase2_lifecycle_subscription"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func subscriptionConfig(description string, disabled bool, filter string) string {
	disabledValue := "false"
	if disabled {
		disabledValue = "true"
	}
	return `resource "criblio_subscription" "my_subscription" {
  description = "` + description + `"
  disabled    = ` + disabledValue + `
  filter      = "` + filter + `"
  group_id    = "default"
  id          = "phase2_lifecycle_subscription"
  pipeline    = "passthru"
}
`
}
