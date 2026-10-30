resource "criblio_notification_target" "example" {
  id = "monitor-policy-target"
  sns_target = {
    id                        = "monitor-policy-target"
    type                      = "sns"
    region                    = "us-east-1"
    destination_type          = "topic"
    topic_type                = "standard"
    topic_arn                 = "arn:aws:sns:us-east-1:123456789012:monitor-alerts"
    aws_authentication_method = "auto"
    allowlist                 = []
  }
}

resource "criblio_notification_template" "example" {
  id      = "monitor-policy-template"
  type    = "sns"
  content = "[{{status}}] {{monitor_name}}"
}

resource "criblio_notification_policy" "example" {
  id              = "platform-alert-routing"
  order           = 10
  disabled        = false
  final           = true
  wait_to_group   = 0
  group_by_labels = ["host"]
  conditions = jsonencode([
    [{ key = "team", operator = "=", value = "platform" }]
  ])
  template_target_pairs = [{
    target_id   = criblio_notification_target.example.id
    template_id = criblio_notification_template.example.id
  }]
}
