resource "criblio_destination" "elasticsearch_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # Elasticsearch configuration from JSON
  output_elastic = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "elastic" || local.destination_type == "elasticsearch"
      error_message = "Destination type must be 'elastic' or 'elasticsearch' for this module. Found: '${local.destination_type}'"
    }
  }
}
