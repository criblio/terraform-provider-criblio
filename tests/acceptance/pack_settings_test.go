package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackSettings(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack test for on-prem: uses HelloPacks which causes API 500 errors")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.edge_supervisor_health", "id", "EdgeSupervisorHealth"),
						resource.TestCheckResourceAttr("criblio_pack.edge_supervisor_health", "group_id", "default_fleet"),
						resource.TestCheckResourceAttr("criblio_pack.edge_supervisor_health", "description", "EdgeSupervisorHealth - Monitors Supervisor process health on Search Deux instances"),
						resource.TestCheckResourceAttr("criblio_pack.edge_supervisor_health", "display_name", "EdgeSupervisorHealth"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.edge_supervisor_health_main", "id", "main"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.edge_supervisor_health_main", "pack", "EdgeSupervisorHealth"),
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
