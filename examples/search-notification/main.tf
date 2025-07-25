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

resource "criblio_notification" "notification_test" {
  condition = "search"
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
  target_configs = [
    {
      conf = {
        body = "test_body"
        email_recipient = {
          bcc = "test@test.com"
          cc  = "test@test.com"
          to  = "test@test.com"
        }
        subject = "test_subject"
      }
      id = "test_target_config_1"
    }
  ]
  targets = [
    criblio_notification_target.my_notificationtarget.id
  ]
}


resource "criblio_notification_target" "my_notificationtarget" {
  id = "test_slack_target_1"
  slack_target = {
    id = "test_slack_target_1"
    system_fields = [
      "test"
    ]
    type = "slack"
    url  = "https://hooks.slack.com/services/T00000000/B00000000/X00000000"
  }
}