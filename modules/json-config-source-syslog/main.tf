resource "criblio_source" "syslog_from_json" {
  id       = local.source_id
  group_id = var.group_id

  # Syslog configuration from JSON
  input_syslog = {
    input_syslog_syslog1 = local.final_config
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Validation to ensure source type matches
resource "terraform_data" "validate_source_type" {
  lifecycle {
    precondition {
      condition     = local.source_type == "syslog"
      error_message = "Source type must be 'syslog' for this module. Found: '${local.source_type}'"
    }
  }
}
