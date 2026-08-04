package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackRoutesGenerated(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack routes test for on-prem: uses source pack installation")
	}

	packID := "test-pack-routes-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: packRoutesConfig(packID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_pack.routes_pack", "id", packID),
					resource.TestCheckResourceAttr("criblio_pack_routes.my_packroutes", "group_id", "default"),
					resource.TestCheckResourceAttr("criblio_pack_routes.my_packroutes", "pack", packID),
					resource.TestCheckResourceAttr("criblio_pack_routes.my_packroutes", "routes.#", "1"),
				),
			},
			{
				Config: packRoutesConfig(packID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "criblio_pack_routes.my_packroutes",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf(`{"group_id":"default","pack":%q}`, packID),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"items",
					"routes.0.group_id",
				},
			},
		},
	})
}

func packRoutesConfig(packID string) string {
	return fmt.Sprintf(`resource "criblio_pack" "routes_pack" {
  id           = %q
  group_id     = "default"
  description  = "Pack with routes"
  disabled     = true
  display_name = "Pack with routes"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_routes" "my_packroutes" {
  group_id = "default"
  pack     = criblio_pack.routes_pack.id
  routes = [
    {
      name     = "my_name"
      pipeline = "main"
    },
  ]
}
`, packID)
}
