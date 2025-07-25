terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "determined-gian-gkh6kzw"
  workspace_id = "main"
}

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
      disabled = true
      items = []
    }
  }
}

resource "criblio_search_saved_query" "my_searchsavedquery_with_notifications" {
  name        = "test_with_notifications"
  description = "test with notifications"
  is_private  = true
  id          = "test_saved_query_with_notifications"
  query       = "cribl dataset=\"cribl_internal_logs\" source=*searches.log message=\"search finished\" | summarize count(), elapsedMS=sum(stats.elapsedMs), eventsFound=sum(stats.eventsFound) by user=stats.user"
  schedule = {
    enabled       = true
    cron_schedule = "0 0 * * *"
    tz            = "UTC"
    keep_last_n   = 2
    notifications = {
      disabled = false
      items = [
        {
          condition = "count() > 100"
          conf = {
            trigger_type        = "resultsCount"
            trigger_comparator  = ">"
            trigger_count       = 10
            saved_query_id      = "test_saved_query_with_notifications"
          }
          disabled = false
          id       = "test_notification"
          metadata = [
            {
              name  = "test_metadata"
              value = "test_value"
            }
          ]
          targets = [
            criblio_notification_target.my_notificationtarget.id
          ]
        }
      ]
    }
  }
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
