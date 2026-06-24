package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchSource(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	resourceName := "criblio_search_source.my_searchsource"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: searchSourceConfig("phase2 search source", 31170),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group_id", "default_search"),
						resource.TestCheckResourceAttr(resourceName, "id", "phase2_search_source"),
						resource.TestCheckResourceAttr(resourceName, "type", "cribl_http"),
						resource.TestCheckResourceAttr(resourceName, "description", "phase2 search source"),
						resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					),
				},
				{
					Config: searchSourceConfig("phase2 search source updated", 31170),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "phase2 search source updated"),
					),
				},
				{
					Config:   searchSourceConfig("phase2 search source updated", 31170),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default_search","id":"phase2_search_source"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func searchSourceConfig(description string, port int) string {
	return `resource "criblio_search_source" "my_searchsource" {
  cribl_api   = "/cribl/_bulk"
  description = "` + description + `"
  disabled    = false
  group_id    = "default_search"
  host        = "0.0.0.0"
  id          = "phase2_search_source"
  port        = ` + fmt.Sprint(port) + `
  type        = "cribl_http"

  tls = {
    disabled = true
  }
}
`
}
