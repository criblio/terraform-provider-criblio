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
  id = "metrics"

  # Both rulesets deploy the default_search group. Serialize their updates to
  # avoid concurrent group deployments.
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

output "search_dataset_ruleset" {
  value = criblio_search_dataset_ruleset.example
}

output "search_dataset_ruleset_id" {
  value = criblio_search_dataset_ruleset.example.id
}

data "criblio_search_dataset_ruleset" "metrics" {
  id = "metrics"
}

output "search_dataset_ruleset_metrics" {
  value = data.criblio_search_dataset_ruleset.metrics
}
