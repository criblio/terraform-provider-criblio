locals {
  panels = [
    {
      id       = "poll-outcomes"
      title    = "Poll outcomes"
      query    = "sum by (source, outcome) (rate(goatbot_poll_total[15m]))"
      metric   = "goatbot_poll_total"
      group_by = ["source", "outcome"]
      range    = "15m"
    },
    {
      id       = "worker-outcomes"
      title    = "Worker run outcomes"
      query    = var.worker_query
      metric   = "goatbot_run_outcome_total"
      group_by = var.worker_group_by
      range    = "5m"
    }
  ]
}

resource "criblio_search_dashboard" "metrics_test" {
  id          = var.dashboard_id
  name        = var.dashboard_name
  description = "Metric dashboard round-trip test for AGT-596 and PLAT-16872"
  category    = var.category
  groups      = {}

  elements = concat([
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
    }
    ], [for index, panel in local.panels : {
      dashboard_element_visualization = {
        id     = panel.id
        type   = "chart.line"
        title  = panel.title
        layout = { x = 0, y = 2 + index * 4, w = 12, h = 4 }
        config = {}
        search = {
          search_query_metric = {
            type                 = "metric"
            source               = "customerMetricsPromql"
            dataset_id           = var.dataset_id
            expected_output_type = "range"
            earliest             = { str = "-1h" }
            latest               = { str = "now" }
            timezone             = "local"
            step_strategy        = var.step_strategy
            queries = [{
              local_id            = "A"
              alias               = "Query A"
              enabled             = true
              mode                = "code"
              dataset_id          = var.dataset_id
              show_promql_preview = false
              query               = panel.query
              builder = {
                metric   = panel.metric
                filters  = []
                agg      = "sum"
                group_by = panel.group_by
                range    = panel.range
              }
            }]
          }
        }
        title_action = {
          label           = "Open in Metrics"
          open_in_new_tab = true
          url             = "/search/metrics/explorer?dataset=${urlencode(var.dataset_id)}&qA=${urlencode(panel.query)}&dashboardId=${urlencode(var.dashboard_id)}&dashboardElementId=${panel.id}"
        }
      }
  }])
}

data "criblio_search_dashboard" "by_id" {
  id = criblio_search_dashboard.metrics_test.id
}

output "dashboard_id" {
  value = criblio_search_dashboard.metrics_test.id
}

output "read_back_elements" {
  value = data.criblio_search_dashboard.by_id.elements
}
