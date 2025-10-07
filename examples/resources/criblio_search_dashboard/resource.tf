resource "criblio_search_dashboard" "my_searchdashboard" {
  cache_ttl_seconds   = 300
  category            = "observability"
  created             = 1696166400
  created_by          = "user123"
  description         = "Dashboard for monitoring system metrics"
  display_created_by  = "User 123"
  display_modified_by = "User 456"
  elements = [
    {
      element = {
        description      = "5xx error rate for checkout"
        empty            = false
        hide_panel       = false
        horizontal_chart = false
        id               = "panel-http-5xx"
        index            = 3
        input_id         = "service"
        layout = {
          h = 0
          w = 9
          x = 8
          y = 2
        }
        search = {
          search_query_inline = {
            earliest = {
              number = 1696166400
            }
            latest = {
              number = 1696170000
            }
            parent_search_id = "parentSearch1"
            query            = "error OR warn"
            sample_rate      = 0.5
            timezone         = "UTC"
            type             = "inline"
          }
        }
        title = "HTTP 5xx errors"
        type  = "chart.line"
        value = {
          key = jsonencode("value")
        }
        variant = "markdown"
      }
    }
  ]
  id           = "dashboard123"
  modified     = 1696170000
  modified_by  = "user456"
  name         = "System Metrics Dashboard"
  pack_id      = "New Pack Id"
  refresh_rate = 60
  resolved_dataset_ids = [
    "string",
    "int",
  ]
  schedule = {
    cron_schedule = "0 * * * *"
    enabled       = true
    keep_last_n   = 5
    notifications = {
      disabled = false
    }
    tz = "UTC"
  }
}