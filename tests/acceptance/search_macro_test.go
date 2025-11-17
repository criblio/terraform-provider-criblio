package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchMacro(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_macro.my_searchmacro", "id", "test_macro_2"),
						resource.TestCheckResourceAttr("criblio_search_macro.my_searchmacro", "description", "test_description"),
						resource.TestCheckResourceAttr("criblio_search_macro.my_searchmacro", "replacement", "test_replacement"),
						resource.TestCheckResourceAttr("criblio_search_macro.my_searchmacro", "tags", "test_tags"),
					),
				},
			},
		})
	})
}
