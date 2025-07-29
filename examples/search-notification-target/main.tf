resource "criblio_notification_target" "my_notificationtarget" {
  id = "test_notification_target_1"
  sns_target = {
    allowlist = [
      "test"
    ]
    assume_role_arn           = "arn:aws:iam::123456789012:role/test"
    assume_role_external_id   = "test"
    aws_authentication_method = "auto"
    destination_type          = "topic"
    endpoint                  = "https://example.com/test"
    id                        = "test_notification_target_1"
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