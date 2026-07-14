package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMappingRuleset(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping mapping ruleset test for on-prem: mappings API returns 403")
	}

	testMappingRulesetForProduct(t, "edge")
}

func TestMappingRulesetStream(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping mapping ruleset test for on-prem: mappings API returns 403")
	}

	testMappingRulesetForProduct(t, "stream")
}

func testMappingRulesetForProduct(t *testing.T, product string) {
	t.Helper()

	resourceName := "criblio_mapping_ruleset.my_mapping_ruleset"
	id := "default"
	final := "true"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: mappingRulesetConfig(product, final, "test mapping", "true", id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "product", product),
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.description", "test mapping"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.final", final),
				),
			},
			{
				Config: mappingRulesetConfig(product, final, "test mapping updated", "!cribl.group", id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.description", "test mapping updated"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.filter", "!cribl.group"),
				),
			},
			{Config: mappingRulesetConfig(product, final, "test mapping updated", "!cribl.group", id), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"product":"` + product + `","id":"` + id + `"}`,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"conf",
				},
			},
		},
	})
}

func mappingRulesetConfig(product, final, description, filter, id string) string {
	return `resource "criblio_mapping_ruleset" "my_mapping_ruleset" {
  active  = true
  product = "` + product + `"
  id      = "` + id + `"
  conf = {
    functions = [
      {
        id          = "eval"
        filter      = "` + filter + `"
        disabled    = false
        final       = ` + final + `
        description = "` + description + `"
        group_id    = "default"
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'default'"
            }
          ]
        }
      }
    ]
  }
}
`
}
