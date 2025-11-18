package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackLookups(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(2 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.lookups_pack", "id", "pack-with-lookups"),
						resource.TestCheckResourceAttr("criblio_pack.lookups_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.lookups_pack", "description", "Pack with lookups"),
						resource.TestCheckResourceAttr("criblio_pack_lookups.my_packlookups", "id", "my_id"),
						resource.TestCheckResourceAttr("criblio_pack_lookups.my_packlookups", "group_id", "default"),
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
