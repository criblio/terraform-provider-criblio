resource "criblio_search_dashboard" "my_searchdashboard" {
  cache_ttl_seconds = 6.24
  category          = "...my_category..."
  description       = "...my_description..."
  elements = [
    {
      dashboard_element = {
        config = {
          markdown = "...my_markdown..."
        }
        hide_panel       = true
        horizontal_chart = false
        id               = "...my_id..."
        layout = {
          h = 0
          w = 9
          x = 8
          y = 2
        }
        search = {
          alias    = "...my_alias..."
          local_id = "...my_local_id..."
          query    = "...my_query..."
        }
        title_action = {
          label           = "...my_label..."
          open_in_new_tab = true
          url             = "...my_url..."
        }
        type    = "markdown.copilot"
        variant = "markdown"
      }
    }
  ]
  groups = {
    key = {
      action = {
        label = "...my_label..."
        params = {
          key = "value"
        }
        target = "...my_target..."
      }
      collapsed = true
      input_id  = "...my_input_id..."
      title     = "...my_title..."
    }
  }
  id           = "dash-overview"
  name         = "...my_name..."
  refresh_rate = 2.1
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