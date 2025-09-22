resource "criblio_source" "tcp_from_json" {
  id       = local.source_id
  group_id = var.group_id

  # TCP configuration from JSON
  input_tcp = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure source type matches
resource "terraform_data" "validate_source_type" {
  lifecycle {
    precondition {
      condition     = local.source_type == "tcp"
      error_message = "Source type must be 'tcp' for this module. Found: '${local.source_type}'"
    }
  }
}
