resource "criblio_search_datatype_ruleset" "example" {
  id = "default"

  rules = [
    {
      id               = "datatype_rule_1"
      name             = "NDJSON default"
      description      = "Route unparsed NDJSON to generic_ndjson"
      kusto_expression = "*"
      datatype         = "generic_ndjson"
      disabled         = false
    }
  ]
}
