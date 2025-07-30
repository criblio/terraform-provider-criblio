package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchRegex(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_regex.my_search_regex", "id", "test_regex"),
						resource.TestCheckResourceAttr("criblio_regex.my_search_regex", "group_id", "default_search"),
						resource.TestCheckResourceAttr("criblio_regex.my_search_regex", "description", "test"),
						resource.TestCheckResourceAttr("criblio_regex.my_search_regex", "tags", "test"),
					),
				},
			},
		})
	})
}
