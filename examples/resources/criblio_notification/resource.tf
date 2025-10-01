resource "criblio_notification" "my_notification" {
  condition = "true"
  conf = {
    message            = "Message for notification"
    saved_query_id     = "savedQueryId"
    trigger_comparator = ">"
    trigger_count      = 0
    trigger_type       = "resultsCount"
  }
  disabled = false
  group    = "myNotificationGroup"
  id       = "myUniqueNotificationId"
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
    "target1",
    "target2",
  ]
}