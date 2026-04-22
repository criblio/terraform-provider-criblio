resource "criblio_notification_target" "my_notificationtarget" {
  id = "test_notification_target_1"
  sns_target = {
    id                        = "test_notification_target_1"
    type                      = "sns"
    region                    = "us-east-1"
    destination_type          = "topic"
    topic_type                = "fifo"
    topic_arn                 = "arn:aws:sns:us-east-1:123456789012:example-topic.fifo"
    message_group_id          = "cribl-notifications"
    endpoint                  = "https://sns.us-east-1.amazonaws.com"
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-sns-notify"
    assume_role_external_id   = "cribl-example"
    aws_authentication_method = "auto"
    allowlist                 = []
    system_fields = [
      "cribl_host",
    ]
  }
}
