package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestHybridGroupSettings(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_hybrid_group_system_settings.hybrid_group_settings", "group_id", "my-hybrid-group"),
						resource.TestCheckResourceAttr("criblio_hybrid_group_system_settings.hybrid_group_settings", "api.base_url", "https://leader.example.com:9000"),
						resource.TestCheckResourceAttr("criblio_hybrid_group_system_settings.hybrid_group_settings", "api.protocol", "https"),
						resource.TestCheckResourceAttr("criblio_hybrid_group_system_settings.hybrid_group_settings", "upgrade_settings.upgrade_source", "cdn"),
					),
				},
			},
		})
	})
}
