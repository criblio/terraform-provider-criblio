resource "criblio_search_saved_query" "my_searchsavedquery" {
  name        = "test"
  description = "test"
  is_private  = true
  id          = "test_saved_query"
  query       = "cribl dataset=\"cribl_internal_logs\" source=*searches.log message=\"search finished\" | summarize count(), elapsedMS=sum(stats.elapsedMs), eventsFound=sum(stats.eventsFound) by user=stats.user"
  schedule = {
    enabled       = false
    cron_schedule = "0 0 * * *"
    tz            = "UTC"
    keep_last_n   = 2
    notifications = {
      disabled = false
      items = []
    }
  }
}

resource "criblio_search_saved_query" "my_searchsavedquery_with_notifications" {
  name        = "test_with_notifications"
  description = "test with notifications"
  is_private  = true
  id          = "test_saved_query_with_notifications_2"
  query       = "cribl dataset=\"cribl_internal_logs\" source=*searches.log message=\"search finished\" | summarize count(), elapsedMS=sum(stats.elapsedMs), eventsFound=sum(stats.eventsFound) by user=stats.user"
  earliest    = "-2h"
  latest      = "now"
  schedule = {
    enabled       = true
    cron_schedule = "0 0 * * *"
    tz            = "UTC"
    keep_last_n   = 2
    notifications = {
      disabled = false
    }
  }
  depends_on = [
    criblio_notification_target.my_notificationtarget
  ]
}

resource "criblio_notification" "my_notification" {
  id       = "test_notification"
  condition = "search"
  conf = {
    trigger_type        = "resultsCount"
    trigger_comparator  = ">"
    trigger_count       = 10
    saved_query_id      = criblio_search_saved_query.my_searchsavedquery_with_notifications.id
    message             = "Date: {{timestamp}}\n\nA notification was triggered for the scheduled search: {{searchId}},\nTenant ID: {{tenantId}}\nSearch ID: {{savedQueryId}}.\nNotification: {{notificationId}}"
  }
  disabled = false
  group    = "default_search"
  targets = [criblio_notification_target.my_notificationtarget.id]
  depends_on = [
    criblio_notification_target.my_notificationtarget
  ]
  target_configs = [
    {
      id = criblio_notification_target.my_notificationtarget.id
      conf = {
        include_results = true
        attachment_type = "inline"
      }
    }
  ]
}

resource "criblio_notification_target" "my_notificationtarget" {
  id = "test_slack_target_2"
  slack_target = {
    id = "test_slack_target_2"
    system_fields = [
      "test"
    ]
    type = "slack"
    url  = "https://hooks.slack.com/services/T00000000/B00000000/X00000000"
  }
}
