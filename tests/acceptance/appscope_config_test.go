package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppscopeConfig(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "description", "A sample AppScope configuration"),
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "lib", "cribl"),
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "id", "sample_appscope_config"),
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "tags", "cribl, test"),
					),
				},
			},
		})
	})
}
