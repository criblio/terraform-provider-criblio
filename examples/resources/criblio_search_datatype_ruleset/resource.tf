resource "criblio_search_datatype_ruleset" "my_searchdatatyperuleset" {
  id = "default"
  rules = [
    {
      datatype         = "generic_ndjson"
      description      = "Route events where vendor is cribl to generic_ndjson"
      disabled         = false
      id               = "rule_1"
      kusto_expression = "vendor == \"cribl\""
      name             = "json events"
    }
  ]
}