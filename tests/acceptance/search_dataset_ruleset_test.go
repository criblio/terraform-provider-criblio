package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDataSetRuleset(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: searchDatasetRulesetConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset_ruleset.example", "id", "default"),
						resource.TestCheckResourceAttr("criblio_search_dataset_ruleset.example", "rules.#", "2"),
						resource.TestCheckResourceAttr("criblio_search_dataset_ruleset.metrics", "id", "metrics"),
						resource.TestCheckResourceAttr("criblio_search_dataset_ruleset.metrics", "rules.#", "1"),
					),
				},
			},
		})
	})
}

func searchDatasetRulesetConfig() string {
	return `
resource "criblio_search_dataset_ruleset" "example" {
  id = "default"
  rules = [
    {
      id               = "rule_1"
      name             = "security logs"
      description      = "Route vendor Cribl events to the main dataset"
      kusto_expression = "vendor == \"cribl\""
      send_data_to     = "destinationDataset"
      dataset          = "main"
      disabled         = false
    },
    {
      id               = "rule_2"
      name             = "test"
      description      = "test data"
      kusto_expression = "*"
      send_data_to     = "destinationDataset"
      dataset          = "main"
      disabled         = false
    }
  ]
}

resource "criblio_search_dataset_ruleset" "metrics" {
  id         = "metrics"
  depends_on = [criblio_search_dataset_ruleset.example]

  rules = [
    {
      id               = "prometheus_metrics"
      name             = "Prometheus metrics"
      description      = "Route Prometheus remote-write metrics to an existing dataset"
      kusto_expression = "__inputId == \"prometheus_rw:production\""
      send_data_to     = "destinationDataset"
      dataset          = "metrics"
      disabled         = false
    }
  ]
}
`
}
