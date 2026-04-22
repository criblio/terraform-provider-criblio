resource "criblio_notification_target" "alerts" {
  id = "test_notification_target_2"
  sns_target = {
    id                        = "test_notification_target_2"
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
    # Not used when destination_type is topic; empty is typical
    allowlist = []
    system_fields = [
      "cribl_host",
    ]
  }
}

# Stream / source condition: high data volume (matches UI advanced settings)
resource "criblio_notification" "source_high_volume" {
  id        = "test_source_volume_notification"
  group     = "default"
  condition = "high-volume"
  conf = {
    name                 = "splunk:in_splunk_tcp"
    time_window          = "60s"
    data_volume          = "1TB"
    notify_on_resolution = true
  }
  # Each value is a JavaScript expression (UI shows "JS" on the value field), as for source metadata:
  # use a string literal in single quotes, or a backtick-wrapped expression.
  metadata = [
    {
      name  = "env"
      value = "'production'"
    },
    {
      name  = "label"
      value = "`C.vars.example || 'n/a'`"
    },
  ]
  targets = [criblio_notification_target.alerts.id]
}
