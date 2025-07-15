package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEdgeFleet(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					ImportStateId:   "criblio_group.my_edge_fleet",
					ResourceName:    "criblio_group.my_edge_fleet",
					ImportState:     true,
				},
				{
					ResourceName:    "criblio_group.my_edge_fleet",
					ImportState:     true,
					ImportStateId:   "criblio_group.my_edge_fleet",
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.my_edge_fleet", "id", "my-edge-fleet"),
						resource.TestCheckResourceAttr("criblio_group.my_edge_fleet", "name", "my-edge-fleet"),
						resource.TestCheckResourceAttr("criblio_group.my_edge_fleet", "product", "edge"),
					),
				},
			},
		})
	})
}
