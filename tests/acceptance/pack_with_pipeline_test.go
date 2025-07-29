package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackPipeline(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.pipeline_pack", "id", "pack-with-pipeline"),
						resource.TestCheckResourceAttr("criblio_pack.pipeline_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.pipeline_pack", "description", "Pack with pipeline"),
						resource.TestCheckResourceAttr("criblio_pack.pipeline_pack", "display_name", "Pack with pipeline"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "id", "my_id"),
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
