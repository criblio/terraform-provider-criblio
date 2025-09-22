resource "criblio_destination" "webhook_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # Webhook configuration from JSON
  output_webhook = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "webhook"
      error_message = "Destination type must be 'webhook' for this module. Found: '${local.destination_type}'"
    }
  }
}
