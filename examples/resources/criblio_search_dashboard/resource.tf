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
        color_palette    = "blue"
        description      = "CPU usage chart"
        empty            = false
        h                = 4
        hide_panel       = false
        horizontal_chart = true
        id               = "element1"
        layout = {
          h = 0
          w = 9
          x = 8
          y = 2
        }
        query = {
          search_query_saved = {
            query    = "dataset=my_dataset | stats count"
            query_id = "query123"
            run_mode = "lastRun"
            type     = "saved"
          }
        }
        title   = "CPU Usage"
        type    = "chart.line"
        variant = "default"
        w       = 6
        x       = 0
        x_axis = {
          data_field        = "time"
          inverse           = false
          label_interval    = "1m"
          label_orientation = 0
          name              = "Time"
          offset            = 0
          position          = "bottom"
          type              = "time"
        }
        y = 0
        y_axis = {
          data_field = [
            "cpu",
          ]
          interval   = 10
          max        = 100
          min        = 0
          position   = "left"
          scale      = "linear"
          split_line = true
          type       = "value"
        }
      }
    }
  ]
  id           = "dashboard123"
  modified     = 1696170000
  modified_by  = "user456"
  name         = "System Metrics Dashboard"
  owner        = "teamA"
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
  tags = [
    "monitoring",
    "system",
    "cpu",
  ]
}