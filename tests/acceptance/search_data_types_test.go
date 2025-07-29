package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDataTypes(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					//ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_search_eventbreakerruleset", "id", "test_eventbreakerruleset"),
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_search_eventbreakerruleset", "group_id", "default_search"),
					),
				},
			},
		})
	})
}
