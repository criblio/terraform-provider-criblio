package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchMacro(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	resourceName := "criblio_search_macro.my_searchmacro"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: searchMacroConfig("test search macro", "source=*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_search_macro"),
					resource.TestCheckResourceAttr(resourceName, "description", "test search macro"),
					resource.TestCheckResourceAttr(resourceName, "replacement", "source=*"),
					resource.TestCheckResourceAttr(resourceName, "tags", "test"),
				),
			},
			{
				Config: searchMacroConfig("test search macro updated", "source=* | limit 10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test search macro updated"),
					resource.TestCheckResourceAttr(resourceName, "replacement", "source=* | limit 10"),
				),
			},
			{
				Config:   searchMacroConfig("test search macro updated", "source=* | limit 10"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "test_search_macro",
				ImportStateVerify: true,
			},
		},
	})
}

func searchMacroConfig(description, replacement string) string {
	return `resource "criblio_search_macro" "my_searchmacro" {
  description = "` + description + `"
  id          = "test_search_macro"
  replacement = "` + replacement + `"
  tags        = "test"
}
`
}
