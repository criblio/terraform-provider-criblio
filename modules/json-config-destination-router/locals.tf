# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Destination ID - use override if provided, otherwise from JSON (required field)
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to router (minimal transformations needed)
  field_mappings = {
    "systemFields" = "system_fields"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "environment", "pipeline", "streamtags", "system_fields", "type", "rules"]

  # Apply field mappings and filter only known fields
  base_config = {
    for key, value in local.raw_config :
    (contains(keys(local.field_mappings), key) ? local.field_mappings[key] : key) => value
    if contains(keys(local.field_mappings), key) || contains(local.direct_fields, key)
  }

  # Transform rules array if present
  transformed_rules = try(local.raw_config.rules, null) != null ? [
    for rule in local.raw_config.rules : {
      description = try(rule.description, null)
      filter      = rule.filter
      output      = rule.output
      final       = try(rule.final, null)
    }
  ] : null

  # Combine all configuration
  final_config = merge(
    {
      id   = local.destination_id
      type = "router"
    },
    # Base config without rules
    { for k, v in local.base_config : k => v if k != "rules" },
    # Add transformed rules if they exist
    local.transformed_rules != null ? { rules = local.transformed_rules } : {},
    var.additional_config
  )
}
