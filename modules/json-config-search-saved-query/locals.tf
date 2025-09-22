# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Query ID - use override if provided, otherwise from JSON, fallback to generated ID
  query_id = var.query_id_override != "" ? var.query_id_override : try(local.raw_config.id, "terraform_query_${substr(md5(jsonencode(local.raw_config)), 0, 8)}")

  # Field mappings specific to search_saved_query
  field_mappings = {
    # Search saved query specific mappings
    "isPrivate" = "is_private"

    # Schedule nested fields
    "cronSchedule" = "cron_schedule"
    "keepLastN"    = "keep_last_n"
    "disabled"     = "disabled"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "name", "query", "description", "earliest", "latest", "schedule", "enabled", "notifications", "tz"]

  # Transform top-level keys (only include known fields)
  base_config = {
    for key, value in local.raw_config :
    (contains(keys(local.field_mappings), key) ? local.field_mappings[key] : key) => value
    if contains(keys(local.field_mappings), key) || contains(local.direct_fields, key)
  }

  # Handle nested schedule object transformation if it exists
  transformed_schedule = try(local.raw_config.schedule, null) != null ? {
    cron_schedule = try(local.raw_config.schedule.cronSchedule, local.raw_config.schedule.cron_schedule, null)
    enabled       = try(local.raw_config.schedule.enabled, null)
    keep_last_n   = try(local.raw_config.schedule.keepLastN, local.raw_config.schedule.keep_last_n, null)
    tz            = try(local.raw_config.schedule.tz, null)
    notifications = try(local.raw_config.schedule.notifications, null) != null ? {
      disabled = try(local.raw_config.schedule.notifications.disabled, null)
    } : null
  } : null

  # Final configuration with additional config merged
  final_config = merge(
    {
      id = local.query_id
    },
    # Base config without schedule
    { for k, v in local.base_config : k => v if k != "schedule" },
    # Add transformed schedule if it exists
    local.transformed_schedule != null ? { schedule = local.transformed_schedule } : {},
    var.additional_config
  )
}
