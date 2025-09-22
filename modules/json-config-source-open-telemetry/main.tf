resource "criblio_source" "open_telemetry_from_json" {
  id       = local.source_id
  group_id = var.group_id

  # OpenTelemetry configuration from JSON
  input_open_telemetry = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure source type matches
resource "terraform_data" "validate_source_type" {
  lifecycle {
    precondition {
      condition     = local.source_type == "open_telemetry"
      error_message = "Source type must be 'open_telemetry' for this module. Found: '${local.source_type}'"
    }
  }
}
