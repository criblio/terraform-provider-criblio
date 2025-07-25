resource "criblio_search_dashboard" "my_searchdashboard" {
  cache_ttl_seconds   = 6.24
  category            = "...my_category..."
  created             = 6.32
  created_by          = "...my_created_by..."
  description         = "...my_description..."
  display_created_by  = "...my_display_created_by..."
  display_modified_by = "...my_display_modified_by..."
  elements = [
    {
      element = {
        description      = "...my_description..."
        empty            = false
        hide_panel       = false
        horizontal_chart = true
        id               = "...my_id..."
        index            = 5.35
        input_id         = "...my_input_id..."
        layout = {
          h = 3.81
          w = 5.73
          x = 5.53
          y = 6.72
        }
        search = {
          search_query_inline = {
            earliest = {
              number = 2.23
            }
            latest = {
              number = 2.43
            }
            parent_search_id = "...my_parent_search_id..."
            query            = "...my_query..."
            sample_rate      = 5.02
            timezone         = "...my_timezone..."
            type             = "inline"
          }
        }
        title = "...my_title..."
        type  = "chart.line"
        value = {
          key = jsonencode("value")
        }
        variant = "input"
      }
    }
  ]
  id           = "...my_id..."
  modified     = 3.77
  modified_by  = "...my_modified_by..."
  name         = "...my_name..."
  pack_id      = "...my_pack_id..."
  refresh_rate = 2.1
  resolved_dataset_ids = [
    "..."
  ]
  schedule = {
    cron_schedule = "...my_cron_schedule..."
    enabled       = true
    keep_last_n   = 5.1
    notifications = {
      disabled = false
      items = [
        {
          condition = "...my_condition..."
          conf = {
            # ...
          }
          disabled = false
          group    = "...my_group..."
          id       = "...my_id..."
          metadata = [
            {
              name  = "...my_name..."
              value = "...my_value..."
            }
          ]
          target_configs = [
            {
              conf = {
                body = "...my_body..."
                email_recipient = {
                  bcc = "Lillian95@gmail.com"
                  cc  = "Maud.Dickens@hotmail.com"
                  to  = "Hillard.Feeney@yahoo.com"
                }
                subject = "...my_subject..."
              }
              id = "...my_id..."
            }
          ]
          targets = [
            "..."
          ]
        }
      ]
    }
    resume_missed  = true
    resume_on_boot = false
    tz             = "...my_tz..."
  }
}