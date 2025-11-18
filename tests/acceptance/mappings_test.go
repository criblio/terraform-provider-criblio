package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMappings(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(2 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Stream mappings resource checks
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "id", "stream_mappings"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "product", "stream"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.#", "6"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.id", "eval"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.description", "Production leaders"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.filter", "env == \"prod\" && role == \"leader\""),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.disabled", "false"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.final", "false"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.group_id", "default"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.conf.add.#", "1"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.conf.add.0.name", "groupId"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.0.conf.add.0.value", "prod-leaders"),
						// Check the last function has final=true
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.5.final", "true"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.5.description", "Default fallback"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.stream_mappings", "conf.functions.5.filter", "true"),

						// Edge mappings resource checks
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "id", "edge_mappings"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "product", "edge"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.#", "6"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.id", "eval"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.description", "North America network devices"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.filter", "region == \"na\" && device_type == \"network\""),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.disabled", "false"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.final", "false"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.group_id", "default"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.conf.add.#", "1"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.conf.add.0.name", "groupId"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.0.conf.add.0.value", "na-network-devices"),
						// Check the last function has final=true
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.5.final", "true"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.5.description", "Default edge fleet"),
						resource.TestCheckResourceAttr("criblio_mapping_ruleset.edge_mappings", "conf.functions.5.filter", "true"),
					),
				},
			},
		})
	})
}
