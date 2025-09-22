resource "criblio_destination" "syslog_from_json" {
  group_id = var.group_id
  id       = local.destination_id

  output_syslog = local.final_config
}

# Validation resource to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "syslog"
      error_message = "Destination type must be 'syslog', got '${local.destination_type}'"
    }
  }
}
