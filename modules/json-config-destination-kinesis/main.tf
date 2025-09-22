resource "criblio_destination" "kinesis_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # Kinesis configuration from JSON
  output_kinesis = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "kinesis"
      error_message = "Destination type must be 'kinesis' for this module. Found: '${local.destination_type}'"
    }
  }
}
