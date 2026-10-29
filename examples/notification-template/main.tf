terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

# Configure the provider using your usual Cribl credentials and workspace.
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

data "criblio_notification_template" "default_email" {
  id = "default-email"
}

data "criblio_notification_templates" "all" {}

output "template_id" {
  value = criblio_notification_template.email.id
}
