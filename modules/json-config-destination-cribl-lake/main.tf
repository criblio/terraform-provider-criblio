resource "criblio_destination" "cribl_lake_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # Cribl Lake configuration from JSON
  output_cribl_lake = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "cribl_lake"
      error_message = "Destination type must be 'cribl_lake' for this module. Found: '${local.destination_type}'"
    }
  }
}
