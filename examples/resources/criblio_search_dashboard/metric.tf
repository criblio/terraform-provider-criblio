resource "criblio_search_dashboard" "metrics" {
  id          = "metrics_starter"
  name        = "Metrics Starter"
  description = "Customer metrics dashboard"
  groups      = {}

  elements = [
    {
      dashboard_element_input = {
        id       = "time-range"
        input_id = "time"
        type     = "input.timerange"
        variant  = "input"
        title    = "Time Range"
        layout   = { x = 0, y = 0, w = 4, h = 2 }
        config = {
          defaultValue = jsonencode({ earliest = "-1h", latest = "now" })
        }
      }
    },
    {
      dashboard_element_visualization = {
        id     = "worker-outcomes"
        type   = "chart.line"
        title  = "Worker run outcomes"
        layout = { x = 0, y = 2, w = 12, h = 4 }
        config = {}
        search = {
          search_query_metric = {
            type                 = "metric"
            source               = "customerMetricsPromql"
            dataset_id           = "metrics"
            expected_output_type = "range"
            earliest             = { str = "-1h" }
            latest               = { str = "now" }
            timezone             = "local"
            step_strategy        = { type = "auto", target_points = 240 }
            queries = [
              {
                local_id            = "A"
                alias               = "Query A"
                enabled             = true
                mode                = "code"
                dataset_id          = "metrics"
                show_promql_preview = false
                query               = "sum by (issue, outcome) (goatbot_run_outcome_total)"
                builder = {
                  metric   = "goatbot_run_outcome_total"
                  filters  = []
                  agg      = "sum"
                  group_by = ["issue", "outcome"]
                  range    = "5m"
                }
              }
            ]
          }
        }
        title_action = {
          label           = "Open in Metrics"
          open_in_new_tab = true
          url             = "/search/metrics/explorer?dataset=metrics&dashboardId=metrics_starter&dashboardElementId=worker-outcomes"
        }
      }
    }
  ]
}
