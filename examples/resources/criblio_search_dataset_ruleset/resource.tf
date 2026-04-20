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