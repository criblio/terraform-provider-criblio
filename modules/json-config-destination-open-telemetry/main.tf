resource "criblio_destination" "open_telemetry_from_json" {
  group_id = var.group_id
  id       = local.destination_id

  output_open_telemetry = local.final_config
}

# Validation resource to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "open_telemetry"
      error_message = "Destination type must be 'open_telemetry', got '${local.destination_type}'"
    }
  }
}
