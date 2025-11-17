package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSubscription(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is 'prohibited by current license'")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription", "id", "my_subscription"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription", "description", "test subscription"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription", "filter", "test"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription", "group_id", "default"),
					),
				},
			},
		})
	})
}
