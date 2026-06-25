package tests

import (
	"os"
	"strconv"
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
					Config: searchSavedQueryConfig("test saved query", "dataset=\"cribl_internal_logs\" | limit 10"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "test_search_saved_query"),
						resource.TestCheckResourceAttr(resourceName, "name", "test saved query"),
						resource.TestCheckResourceAttr(resourceName, "description", "test saved query"),
						resource.TestCheckResourceAttr(resourceName, "is_private", "true"),
					),
				},
				{
					Config: searchSavedQueryConfig("test saved query updated", "dataset=\"cribl_internal_logs\" | limit 20"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "test saved query updated"),
						resource.TestCheckResourceAttr(resourceName, "description", "test saved query updated"),
						resource.TestCheckResourceAttr(resourceName, "query", "dataset=\"cribl_internal_logs\" | limit 20"),
					),
				},
				{
					Config:   searchSavedQueryConfig("test saved query updated", "dataset=\"cribl_internal_logs\" | limit 20"),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     "test_search_saved_query",
					ImportStateVerify: true,
				},
			},
		})
	})
}

func searchSavedQueryConfig(description, query string) string {
	return `resource "criblio_search_saved_query" "my_searchsavedquery" {
  description = ` + strconv.Quote(description) + `
  earliest    = "-1h"
  id          = "test_search_saved_query"
  is_private  = true
  latest      = "now"
  name        = ` + strconv.Quote(description) + `
  query       = ` + strconv.Quote(query) + `
}
`
}
