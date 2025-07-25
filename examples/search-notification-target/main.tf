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

resource "criblio_notification_target" "my_notificationtarget" {
  id = "test_notification_target"
  sns_target = {
    allowlist = [
      "test"
    ]
    assume_role_arn           = "arn:aws:iam::123456789012:role/test"
    assume_role_external_id   = "test"
    aws_authentication_method = "auto"
    destination_type          = "topic"
    endpoint                  = "https://example.com/test"
    id                        = "test_notification_target"
    message_group_id          = "test"
    region                    = "us-east-1"
    system_fields = [
      "test"
    ]
    topic_arn  = "arn:aws:sns:us-east-1:123456789012:test"
    topic_type = "fifo"
    type       = "sns"
  }
}