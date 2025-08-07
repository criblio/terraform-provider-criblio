resource "criblio_search_saved_query" "my_searchsavedquery" {
  description = "test saved query"
  earliest    = "-1h"
  id          = "test_saved"
  is_private  = true
  latest      = "now"
  name        = "test_saved"
  query       = "dataset=\"cribl_edge_appscope_metrics\" | limit 1000"
  schedule = {
    cron_schedule = "...my_cron_schedule..."
    enabled       = true
    keep_last_n   = 7.78
    notifications = {
      disabled = false
    }
    tz = "...my_tz..."
  }
}