resource "criblio_destination" "router_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  output_router = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.final_config.type == "router"
      error_message = "Destination type must be 'router' for this module. Found: '${local.final_config.type}'"
    }
  }
}
