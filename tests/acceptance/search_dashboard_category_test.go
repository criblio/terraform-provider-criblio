package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDashboardCategory(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory", "id", "test_dashboard_category"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory", "name", "test_dashboard_category"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory", "description", "test"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory", "is_pack", "true"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory_not_pack", "id", "test_dashboard_category_not_pack"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory_not_pack", "name", "test_dashboard_category_not_pack"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory_not_pack", "description", "test"),
						resource.TestCheckResourceAttr("criblio_search_dashboard_category.my_searchdashboardcategory_not_pack", "is_pack", "false"),
					),
				},
			},
		})
	})
}
