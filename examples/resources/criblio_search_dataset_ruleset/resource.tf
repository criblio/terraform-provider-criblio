resource "criblio_search_dataset_ruleset" "my_searchdatasetruleset" {
  id = "default"
  rules = [
    {
      dataset                   = "main"
      description               = "data catch all"
      disabled                  = false
      extend_expression         = "vendor = toupper(vendor)"
      extend_expression_enabled = true
      id                        = "default"
      kusto_expression          = "*"
      name                      = "main"
      send_data_to              = "destinationDataset"
    }
  ]
}

resource "criblio_search_dataset_ruleset" "metrics" {
  id = "metrics"

  # Both rulesets deploy the default_search group. Serialize their updates to
  # avoid concurrent group deployments.
  depends_on = [criblio_search_dataset_ruleset.my_searchdatasetruleset]

  rules = [
    {
      dataset          = "main"
      description      = "Route Prometheus remote-write metrics to an existing dataset"
      disabled         = false
      id               = "prometheus_metrics"
      kusto_expression = "__inputId == \"prometheus_rw:production\""
      name             = "Prometheus metrics"
      send_data_to     = "destinationDataset"
    }
  ]
}
