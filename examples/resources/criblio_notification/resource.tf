resource "criblio_notification" "my_notification" {
  condition = "true"
  conf = {
    data_volume          = "1TB"
    message              = "Message for notification"
    name                 = "splunk:in_splunk_tcp"
    notify_on_resolution = true
    saved_query_id       = "savedQueryId"
    time_window          = "60s"
    trigger_comparator   = ">"
    trigger_count        = 10
    trigger_type         = "resultsCount"
    usage_threshold      = 90
    worker_group         = "...my_worker_group..."
  }
  disabled = false
  group    = "myNotificationGroup"
  id       = "myUniqueNotificationId"
  metadata = [
    {
      name  = "env"
      value = "production"
    }
  ]
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