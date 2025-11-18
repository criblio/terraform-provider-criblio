package tests

import (
	"os"
	"time"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestRoutes(t *testing.T) {
        if os.Getenv("DEPLOYMENT") == "onprem" {
                time.Sleep(1 * time.Second)
        }
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_routes.my_routes", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_routes.my_routes", "routes.#", "2"),
						resource.TestCheckResourceAttr("criblio_routes.my_routes", "routes.0.name", "my_route_1"),
						resource.TestCheckResourceAttr("criblio_routes.my_routes", "routes.0.pipeline", "main"),
						resource.TestCheckResourceAttr("criblio_routes.my_routes", "routes.1.name", "my_route_2"),
						resource.TestCheckResourceAttr("criblio_routes.my_routes", "routes.1.pipeline", "main"),
					),
				},
				{
					ConfigDirectory: config.TestNameDirectory(),
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
