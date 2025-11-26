package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEdgeSubfleet(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ResourceName:    "criblio_group.my_edge_subfleet[0]",
					ImportState:     true,
					ImportStateId:   "my-edge-subfleet",
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.my_edge_subfleet[0]", "id", "my-edge-subfleet"),
						resource.TestCheckResourceAttr("criblio_group.my_edge_subfleet[0]", "name", "my-edge-subfleet"),
						resource.TestCheckResourceAttr("criblio_group.my_edge_subfleet[0]", "product", "edge"),
					),
				},
			},
		})
	})
}
