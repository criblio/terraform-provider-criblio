package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackBreakers(t *testing.T) {
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
						resource.TestCheckResourceAttr("criblio_pack.breakers_pack", "id", "pack-breakers"),
						resource.TestCheckResourceAttr("criblio_pack.breakers_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.breakers_pack", "description", "Pack breakers"),
						resource.TestCheckResourceAttr("criblio_pack_breakers.my_packbreakers", "id", "test_packbreakers"),
						resource.TestCheckResourceAttr("criblio_pack_breakers.my_packbreakers", "group_id", "default"),
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
