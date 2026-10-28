resource "criblio_notification_template" "email" {
  id          = "custom-email"
  type        = "smtp"
  description = "Custom email notification template"
  content = jsonencode({
    to      = "{{metadata.to}}"
    subject = "[{{status}}] {{monitor_name}}"
    body    = "{{description}}"
  })
}
