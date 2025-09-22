resource "criblio_destination" "cribl_tcp_from_json" {
  group_id = var.group_id
  id       = local.destination_id

  output_cribl_tcp = local.final_config
}

# Validation resource to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "cribl_tcp"
      error_message = "Destination type must be 'cribl_tcp', got '${local.destination_type}'"
    }
  }
}
