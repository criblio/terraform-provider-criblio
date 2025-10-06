resource "criblio_notification_target" "my_notificationtarget" {
  id = "nt-slack-prod"
  pager_duty_target = {
    class       = "incident-class"
    component   = "webserver"
    group       = "oncall-team"
    id          = "myPagerDutyTarget"
    routing_key = "0123456789abcdef0123456789abcdef"
    severity    = "critical"
    system_fields = [
      "cribl_host",
      "cribl_pipe",
    ]
    type = "pager_duty"
  }
  slack_target = {
    id = "mySlackTarget"
    system_fields = [
      "cribl_host",
      "cribl_pipe",
    ]
    type = "slack"
    url  = "https://hooks.slack.com/services/THISisanExample"
  }
  smtp_target = {
    encryption_option = "STARTTLS"
    from              = "alerts@example.com"
    host              = "smtp.example.com"
    id                = "mySmtpTarget"
    password          = "smtpPass123"
    port              = 587
    system_fields = [
      "cribl_host",
      "cribl_pipe",
    ]
    tls = {
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      reject_unauthorized = true
    }
    type     = "smtp"
    username = "smtpUser"
  }
  sns_target = {
    allowlist = [
      "+15551234567",
      "+15559876543",
    ]
    assume_role_arn           = "arn:aws:iam::123456789012:role/MyRole"
    assume_role_external_id   = "external-id-123"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "manual"
    aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    destination_type          = "topic"
    endpoint                  = "https://sns.us-east-1.amazonaws.com"
    id                        = "mySnsTarget"
    message_group_id          = "my-message-group-id"
    phone_number              = "+15557654321"
    region                    = "us-east-1"
    system_fields = [
      "cribl_host",
      "cribl_pipe",
    ]
    topic_arn  = "arn:aws:sns:us-east-1:123456789012:my-topic"
    topic_type = "fifo"
    type       = "sns"
  }
  webhook_target = {
    auth_type = "token"
    format    = "ndjson"
    id        = "myWebhookTarget"
    method    = "POST"
    password  = "webhookPass"
    system_fields = [
      "cribl_host",
      "cribl_pipe",
    ]
    token    = "myWebhookSecretToken"
    type     = "webhook"
    url      = "https://hooks.example.com/webhook"
    username = "webhookUser"
  }
}