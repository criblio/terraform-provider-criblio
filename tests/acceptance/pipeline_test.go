package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPipeline(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ExpectNonEmptyPlan: true, 
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pipeline.my_pipeline", "id", "pipeline-1"),
						resource.TestCheckResourceAttr("criblio_pipeline.my_pipeline", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pipeline.my_pipeline", "conf.description", "my_description"),
						resource.TestCheckResourceAttr("criblio_pipeline.my_pipeline", "conf.async_func_timeout", "60"),
						resource.TestCheckResourceAttr("criblio_pipeline.my_pipeline", "conf.output", "my_output"),
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
