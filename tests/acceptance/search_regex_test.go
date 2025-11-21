package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchRegex(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
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
