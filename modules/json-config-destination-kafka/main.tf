resource "criblio_destination" "kafka_from_json" {
  id       = local.destination_id
  group_id = var.group_id

  # Kafka configuration from JSON
  output_kafka = local.final_config

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure destination type matches
resource "terraform_data" "validate_destination_type" {
  lifecycle {
    precondition {
      condition     = local.destination_type == "kafka"
      error_message = "Destination type must be 'kafka' for this module. Found: '${local.destination_type}'"
    }
  }
}
