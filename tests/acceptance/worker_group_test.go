package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkerGroup(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ImportState:     true,
					ImportStateId:   "criblio_group.my_group",
					ResourceName:    "criblio_group.my_group",
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.my_group", "id", "newgroup"),
						resource.TestCheckResourceAttr("criblio_group.my_group", "name", "newgroup"),
						resource.TestCheckResourceAttr("criblio_group.my_group", "product", "stream"),
					),
				},
			},
		})
	})
}
