package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDashboardCategory(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	resourceName := "criblio_search_dashboard_category.my_searchdashboardcategory"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: searchDashboardCategoryConfig("test search dashboard category"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_search_dashboard_category"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_search_dashboard_category"),
					resource.TestCheckResourceAttr(resourceName, "description", "test search dashboard category"),
					resource.TestCheckResourceAttr(resourceName, "is_pack", "false"),
				),
			},
			{
				Config: searchDashboardCategoryConfig("test search dashboard category updated"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "test search dashboard category updated"),
			},
			{
				Config:   searchDashboardCategoryConfig("test search dashboard category updated"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "test_search_dashboard_category",
				ImportStateVerify: true,
			},
		},
	})
}

func searchDashboardCategoryConfig(description string) string {
	return `resource "criblio_search_dashboard_category" "my_searchdashboardcategory" {
  description = "` + description + `"
  id          = "test_search_dashboard_category"
  is_pack     = false
  name        = "test_search_dashboard_category"
}
`
}
