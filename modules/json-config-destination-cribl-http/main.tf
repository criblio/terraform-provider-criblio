resource "criblio_destination" "cribl_http_from_json" {
  group_id = var.group_id
  id       = local.destination_id

  output_cribl_http = local.final_config
}

# Validation resource to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "cribl_http"
      error_message = "Destination type must be 'cribl_http', got '${local.destination_type}'"
    }
  }
}
