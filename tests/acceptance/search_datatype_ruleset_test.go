package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDataTypeRuleset(t *testing.T) {
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
						resource.TestCheckResourceAttr("criblio_search_datatype_ruleset.example", "id", "default"),
						resource.TestCheckResourceAttr("criblio_search_datatype_ruleset.example", "rules.#", "1"),
						resource.TestCheckResourceAttr("criblio_search_datatype_ruleset.example", "rules.0.datatype", "generic_ndjson"),
					),
				},
			},
		})
	})
}
