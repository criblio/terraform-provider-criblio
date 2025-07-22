package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackDestination(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "id", "pack-with-dest"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "description", "Pack from source with destination"),
						resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "id", "my_id"),
						resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "author", "my_author"),
						resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "description", "my_description"),
						resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "my_display_name", "my_display_name"),
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
