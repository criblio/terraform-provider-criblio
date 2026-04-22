resource "criblio_search_dataset_ruleset" "example" {
  id = "default"

  rules = [
    {
      id               = "rule_1"
      name             = "security logs"
      description      = "Route vendor Cribl events to main"
      kusto_expression = "vendor == \"cribl\""
      send_data_to     = "destinationDataset"
      dataset          = "my-dataset-id"
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

output "search_dataset_ruleset" {
  value = criblio_search_dataset_ruleset.example
}

output "search_dataset_ruleset_id" {
  value = criblio_search_dataset_ruleset.example.id
}
