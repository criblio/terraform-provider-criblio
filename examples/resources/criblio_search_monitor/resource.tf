# Requires an Insights-enabled Enterprise Cribl.Cloud workspace.
resource "criblio_search_monitor" "example" {
  id                        = "worker-events-in"
  name                      = "Worker events in"
  product                   = "stream"
  enabled                   = true
  query                     = jsonencode("sum(rate(total_in_events{namespace=\"data_insights\"}[1m]))")
  params                    = jsonencode({})
  firing_after              = 300
  ok_after                  = 300
  schedule_interval_seconds = 60

  rules = [{
    name          = "Low event rate"
    show_on_chart = true
    included_tags = jsonencode({})
    excluded_tags = jsonencode({})
    conditions = [{
      enabled   = true
      labels    = { severity = "warning" }
      condition = jsonencode({ type = "less_than", threshold = 1 })
    }]
  }]
}
