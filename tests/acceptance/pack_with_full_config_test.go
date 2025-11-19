package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackWithFullConfig(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.full_config_pack", "id", "pack-with-full-config"),
						resource.TestCheckResourceAttr("criblio_pack.full_config_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.full_config_pack", "description", "Pack with full pipeline configuration"),
						resource.TestCheckResourceAttr("criblio_pack.full_config_pack", "display_name", "Pack with Full Config"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.AuditdLogs_main", "id", "main"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.AuditdLogs_main", "pack", "pack-with-full-config"),
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
