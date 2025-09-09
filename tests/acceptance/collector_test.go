package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCollector(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector", "id", "splunk-demo-collector"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector", "environment", "demo"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector", "collector.type", "splunk"),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "id", "rest-api-demo-collector"),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "environment", "demo"),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "collector.type", "rest"),
					),
				},
			},
		})
	})
}
