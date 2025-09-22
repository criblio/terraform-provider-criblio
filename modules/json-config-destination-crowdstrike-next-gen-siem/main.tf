resource "criblio_destination" "crowdstrike_from_json" {
  group_id = var.group_id
  id       = local.destination_id

  output_crowdstrike_next_gen_siem = local.final_config
}

# Validation resource to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "crowdstrike_next_gen_siem"
      error_message = "Destination type must be 'crowdstrike_next_gen_siem', got '${local.destination_type}'"
    }
  }
}
