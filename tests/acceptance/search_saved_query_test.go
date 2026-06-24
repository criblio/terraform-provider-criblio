package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchSavedQuery(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	resourceName := "criblio_search_saved_query.my_searchsavedquery"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: searchSavedQueryConfig("phase2 saved query", "dataset=\"cribl_internal_logs\" | limit 10"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group_id", "default_search"),
						resource.TestCheckResourceAttr(resourceName, "id", "phase2_search_saved_query"),
						resource.TestCheckResourceAttr(resourceName, "name", "phase2 saved query"),
						resource.TestCheckResourceAttr(resourceName, "description", "phase2 saved query"),
						resource.TestCheckResourceAttr(resourceName, "is_private", "true"),
					),
				},
				{
					Config: searchSavedQueryConfig("phase2 saved query updated", "dataset=\"cribl_internal_logs\" | limit 20"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "phase2 saved query updated"),
						resource.TestCheckResourceAttr(resourceName, "description", "phase2 saved query updated"),
						resource.TestCheckResourceAttr(resourceName, "query", "dataset=\"cribl_internal_logs\" | limit 20"),
					),
				},
				{
					Config:   searchSavedQueryConfig("phase2 saved query updated", "dataset=\"cribl_internal_logs\" | limit 20"),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default_search","id":"phase2_search_saved_query"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func searchSavedQueryConfig(description, query string) string {
	return `resource "criblio_search_saved_query" "my_searchsavedquery" {
  description = "` + description + `"
  earliest    = "-1h"
  group_id    = "default_search"
  id          = "phase2_search_saved_query"
  is_private  = true
  latest      = "now"
  name        = "` + description + `"
  query       = <<-EOT
` + query + `
EOT
}
`
}
