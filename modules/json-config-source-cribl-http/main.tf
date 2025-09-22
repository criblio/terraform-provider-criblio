resource "criblio_source" "cribl_http_from_json" {
  id       = local.source_id
  group_id = var.group_id

  # Cribl HTTP configuration from JSON
  input_cribl_http = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure source type matches
resource "terraform_data" "validate_source_type" {
  lifecycle {
    precondition {
      condition     = local.source_type == "cribl_http"
      error_message = "Source type must be 'cribl_http' for this module. Found: '${local.source_type}'"
    }
  }
}
