resource "criblio_search_datatype_ruleset" "example" {
  id = "default"

  rules = [
    {
      id               = "rule_1"
      name             = "json events"
      description      = "Route events where vendor is cribl to generic_ndjson"
      disabled         = false
      kusto_expression = "vendor == \"cribl\""
      datatype         = "generic_ndjson"
    }
  ]
}

output "search_datatype_ruleset" {
  value = criblio_search_datatype_ruleset.example
}

output "search_datatype_ruleset_id" {
  value = criblio_search_datatype_ruleset.example.id
}
