package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSubscription(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription.0", "id", "my_subscription"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription.0", "description", "test subscription"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription.0", "filter", "test"),
						resource.TestCheckResourceAttr("criblio_subscription.my_subscription.0", "group_id", "default"),
					),
				},
			},
		})
	})
}
