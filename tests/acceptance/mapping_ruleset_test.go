package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMappingRuleset(t *testing.T) {
	resourceName := "criblio_mapping_ruleset.my_mapping_ruleset"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: mappingRulesetConfig("phase2 mapping", "true", "phase2_mapping_ruleset"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "product", "edge"),
					resource.TestCheckResourceAttr(resourceName, "id", "phase2_mapping_ruleset"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.description", "phase2 mapping"),
				),
			},
			{
				Config: mappingRulesetConfig("phase2 mapping updated", "!cribl.group", "phase2_mapping_ruleset"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.description", "phase2 mapping updated"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.filter", "!cribl.group"),
				),
			},
			{Config: mappingRulesetConfig("phase2 mapping updated", "!cribl.group", "phase2_mapping_ruleset"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"product":"edge","id":"phase2_mapping_ruleset"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func mappingRulesetConfig(description, filter, id string) string {
	return `resource "criblio_mapping_ruleset" "my_mapping_ruleset" {
  product = "edge"
  id      = "` + id + `"
  conf = {
    functions = [
      {
        id          = "eval"
        filter      = "` + filter + `"
        disabled    = false
        final       = true
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
