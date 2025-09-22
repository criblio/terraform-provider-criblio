resource "criblio_destination" "s3_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # S3 configuration from JSON
  output_s3 = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "s3"
      error_message = "Destination type must be 's3' for this module. Found: '${local.destination_type}'"
    }
  }
}
