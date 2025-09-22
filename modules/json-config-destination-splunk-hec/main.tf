resource "criblio_destination" "splunk_hec_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # Splunk HEC configuration from JSON
  output_splunk_hec = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "splunk_hec"
      error_message = "Destination type must be 'splunk_hec' for this module. Found: '${local.destination_type}'"
    }
  }
}
