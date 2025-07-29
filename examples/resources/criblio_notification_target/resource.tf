resource "criblio_notification_target" "my_notificationtarget" {
  id = "...my_id..."
  pager_duty_target = {
    class       = "...my_class..."
    component   = "...my_component..."
    group       = "...my_group..."
    id          = "...my_id..."
    routing_key = "...my_routing_key..."
    severity    = "critical"
    system_fields = [
      "..."
    ]
    type = "pager_duty"
  }
  slack_target = {
    id = "...my_id..."
    system_fields = [
      "..."
    ]
    type = "slack"
    url  = "...my_url..."
  }
  smtp_target = {
    encryption_option = "SSL"
    from              = "Jasmin4@yahoo.com"
    host              = "...my_host..."
    id                = "...my_id..."
    password          = "...my_password..."
    port              = 5508
    system_fields = [
      "..."
    ]
    tls = {
      max_version         = "TLSv1"
      min_version         = "TLSv1.2"
      reject_unauthorized = false
    }
    type     = "smtp"
    username = "...my_username..."
  }
  sns_target = {
    allowlist = [
      "..."
    ]
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret_key            = "...my_aws_secret_key..."
    destination_type          = "topic"
    endpoint                  = "...my_endpoint..."
    id                        = "...my_id..."
    message_group_id          = "...my_message_group_id..."
    phone_number              = "...my_phone_number..."
    region                    = "...my_region..."
    system_fields = [
      "..."
    ]
    topic_arn  = "...my_topic_arn..."
    topic_type = "fifo"
    type       = "sns"
  }
  webhook_target = {
    auth_type = "token"
    format    = "custom"
    id        = "...my_id..."
    method    = "PUT"
    password  = "...my_password..."
    system_fields = [
      "..."
    ]
    token    = "...my_token..."
    type     = "webhook"
    url      = "...my_url..."
    username = "...my_username..."
  }
}