package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLakeHouse(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_cribl_lake_house.my_cribllakehouse", "id", "my_lakehouse"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_house.my_cribllakehouse", "description", "My Lakehouse"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_house.my_cribllakehouse", "tier_size", "medium"),
					),
				},
			},
		})
	})
}
