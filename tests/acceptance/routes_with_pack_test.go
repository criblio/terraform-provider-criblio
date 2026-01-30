package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestRoutesWithPacks(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}
	// Testdata has no provider block (ConfigDirectory constraint); provider uses env vars.
	if os.Getenv("DEPLOYMENT") != "onprem" && (os.Getenv("CRIBL_CLIENT_ID") == "" || os.Getenv("CRIBL_CLIENT_SECRET") == "") {
		t.Skip("Acceptance test requires CRIBL_CLIENT_ID and CRIBL_CLIENT_SECRET (and CRIBL_ORGANIZATION_ID, CRIBL_WORKSPACE_ID, CRIBL_CLOUD_DOMAIN) when not on-prem")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.billing_pipeline", "id", "billing_pipeline"),
						resource.TestCheckResourceAttr("criblio_pack.billing_pipeline", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.billing_pipeline", "display_name", "Billing Pipeline"),
						resource.TestCheckResourceAttr("criblio_routes.routes", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_routes.routes", "routes.#", "1"),
						resource.TestCheckResourceAttr("criblio_routes.routes", "routes.0.name", "Billing Pipeline"),
						resource.TestCheckResourceAttr("criblio_routes.routes", "routes.0.pipeline", "pack:billing_pipeline"),
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
