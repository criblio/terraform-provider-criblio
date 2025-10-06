resource "criblio_notification" "my_notification" {
  condition = "true"
  conf = {
    message            = "Message for notification"
    saved_query_id     = "savedQueryId"
    trigger_comparator = ">"
    trigger_count      = 10
    trigger_type       = "resultsCount"
  }
  disabled = false
  group    = "myNotificationGroup"
  id       = "myUniqueNotificationId"
  target_configs = [
    {
      conf = {
        attachment_type = "attachment"
        include_results = true
      }
      id = "myTargetConfigId"
    }
  ]
  targets = [
    "target1",
    "target2",
  ]
}