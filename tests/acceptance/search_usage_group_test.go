package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchUsageGroup(t *testing.T) {
       if os.Getenv("DEPLOYMENT") == "onprem" {
               t.Skip("Skipping resource for On-Prem deployments as it is not supported onprem")
       }
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_usage_group.my_searchusagegroup", "id", "test_usage_group"),
						resource.TestCheckResourceAttr("criblio_search_usage_group.my_searchusagegroup", "description", "test"),
						resource.TestCheckResourceAttr("criblio_search_usage_group.my_searchusagegroup", "enabled", "true"),
						resource.TestCheckResourceAttr("criblio_search_usage_group.my_searchusagegroup", "coordinator_heap_memory_limit", "8"),
						resource.TestCheckResourceAttr("criblio_search_usage_group.my_searchusagegroup", "rules", "{\"test\": \"test\"}"),
						resource.TestCheckResourceAttr("criblio_search_usage_group.my_searchusagegroup", "users_count", "10"),
					),
				},
			},
		})
	})
}
