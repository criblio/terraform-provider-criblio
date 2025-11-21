package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchLookup(t *testing.T) {
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
						resource.TestCheckResourceAttr("criblio_lookup_file.my_lookupfile_1", "id", "test_id_lookupfile_1.csv"),
						resource.TestCheckResourceAttr("criblio_lookup_file.my_lookupfile_1", "group_id", "default_search"),
						resource.TestCheckResourceAttr("criblio_lookup_file.my_lookupfile_1", "content", "column1, column2, column3, column4"),
						resource.TestCheckResourceAttr("criblio_lookup_file.my_lookupfile_1", "description", "my_description"),
						resource.TestCheckResourceAttr("criblio_lookup_file.my_lookupfile_1", "mode", "memory"),
						resource.TestCheckResourceAttr("criblio_lookup_file.my_lookupfile_1", "tags", "my_tags"),
					),
				},
			},
		})
	})
}
