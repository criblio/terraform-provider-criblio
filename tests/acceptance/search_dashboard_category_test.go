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
				Config: searchDashboardCategoryConfig("phase2 search dashboard category"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default_search"),
					resource.TestCheckResourceAttr(resourceName, "id", "phase2_search_dashboard_category"),
					resource.TestCheckResourceAttr(resourceName, "name", "phase2_search_dashboard_category"),
					resource.TestCheckResourceAttr(resourceName, "description", "phase2 search dashboard category"),
					resource.TestCheckResourceAttr(resourceName, "is_pack", "false"),
				),
			},
			{
				Config: searchDashboardCategoryConfig("phase2 search dashboard category updated"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "phase2 search dashboard category updated"),
			},
			{
				Config:   searchDashboardCategoryConfig("phase2 search dashboard category updated"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default_search","id":"phase2_search_dashboard_category"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func searchDashboardCategoryConfig(description string) string {
	return `resource "criblio_search_dashboard_category" "my_searchdashboardcategory" {
  description = "` + description + `"
  group_id    = "default_search"
  id          = "phase2_search_dashboard_category"
  is_pack     = false
  name        = "phase2_search_dashboard_category"
}
`
}
