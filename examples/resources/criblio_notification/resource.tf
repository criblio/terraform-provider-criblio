resource "criblio_notification" "my_notification" {
  condition = "...my_condition..."
  conf = {
    message            = "...my_message..."
    saved_query_id     = "...my_saved_query_id..."
    trigger_comparator = "...my_trigger_comparator..."
    trigger_count      = 0.24
    trigger_type       = "...my_trigger_type..."
  }
  disabled = false
  group    = "...my_group..."
  id       = "...my_id..."
  target_configs = [
    {
      conf = {
        attachment_type = "inline"
        include_results = false
      }
      id = "...my_id..."
    }
  ]
  targets = [
    "..."
  ]
}