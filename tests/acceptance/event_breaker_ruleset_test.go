package tests

import (
	"os"
	"time"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestEventBreakerRuleset(t *testing.T) {
        if os.Getenv("DEPLOYMENT") == "onprem" {
                time.Sleep(1 * time.Second)
        }

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "id", "test_eventbreakerruleset"),
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "description", "test"),
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "lib", "custom"),
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "tags", "test"),
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
