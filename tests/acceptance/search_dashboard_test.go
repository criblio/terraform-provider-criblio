package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDashboard(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dashboard.my_searchdashboard", "id", "sample_test_dashboard"),
						resource.TestCheckResourceAttr("criblio_search_dashboard.my_searchdashboard", "name", "Sample Test Dashboard"),
						resource.TestCheckResourceAttr("criblio_search_dashboard.my_searchdashboard", "description", "A sample dashboard with several panels"),
					),
				},
			},
		})
	})
}
