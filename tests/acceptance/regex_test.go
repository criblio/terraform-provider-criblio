package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRegex(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "description", "test"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "lib", "custom"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "tags", "test"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "id", "test_regex"),
					),
				},
			},
		})
	})
}
