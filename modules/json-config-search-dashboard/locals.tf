# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Dashboard ID - use override if provided, otherwise from JSON (required field)
  dashboard_id = var.dashboard_id_override != "" ? var.dashboard_id_override : local.raw_config.id

  # Field mappings specific to search_dashboard
  field_mappings = {
    # Dashboard-level field mappings
    "cacheTtlSeconds"    = "cache_ttl_seconds"
    "cacheTTLSeconds"    = "cache_ttl_seconds" # Handle both camelCase variations
    "createdBy"          = "created_by"
    "displayCreatedBy"   = "display_created_by"
    "displayModifiedBy"  = "display_modified_by"
    "modifiedBy"         = "modified_by"
    "packId"             = "pack_id"
    "refreshRate"        = "refresh_rate"
    "resolvedDatasetIds" = "resolved_dataset_ids"

    # Element-level field mappings
    "hidePanel"       = "hide_panel"
    "horizontalChart" = "horizontal_chart"
    "inputId"         = "input_id"

    # Search-level field mappings
    "parentSearchId" = "parent_search_id"
    "sampleRate"     = "sample_rate"
    "queryId"        = "query_id"
    "runMode"        = "run_mode"

    # Schedule field mappings (same as saved query)
    "cronSchedule" = "cron_schedule"
    "keepLastN"    = "keep_last_n"
    "disabled"     = "disabled"
  }

  # Known direct fields that don't need transformation  
  direct_fields = [
    "id", "name", "description", "category",
    "elements", "schedule", "empty", "index", "title", "type", "value", "variant",
    "layout", "search", "query", "timezone", "earliest", "latest", "values",
    "h", "w", "x", "y", "number", "str", "enabled", "notifications", "tz"
  ]

  # Transform top-level keys (only include known fields)
  base_config = {
    for key, value in local.raw_config :
    (contains(keys(local.field_mappings), key) ? local.field_mappings[key] : key) => value
    if contains(keys(local.field_mappings), key) || contains(local.direct_fields, key)
  }

  # Transform elements to match provider schema more closely
  transformed_elements = [
    for element in try(local.raw_config.elements, []) : {
      element = {
        # Direct fields from JSON
        id          = try(element.id, null)
        type        = try(element.type, null)
        title       = try(element.title, null)
        description = try(element.description, null)
        empty       = try(element.empty, null)
        index       = try(element.index, null)
        variant     = try(element.variant, null)

        # Transform layout block
        layout = try(element.layout, null) != null ? {
          x = element.layout.x
          y = element.layout.y
          w = element.layout.w
          h = element.layout.h
        } : null

        # Transform search block - use only fields that are set
        search = try(element.search, null) != null ? {
          search_query_inline = {
            query    = try(element.search.query, null)
            type     = try(element.search.type, "inline")
            timezone = try(element.search.timezone, null)
            # Keep original time format from UI JSON
            earliest = try(element.search.earliest, null) != null ? {
              str = element.search.earliest
            } : null
            latest = try(element.search.latest, null) != null ? {
              str = element.search.latest
            } : null
          }
        } : null

        # Transform value block if it exists
        value = try(element.value, null) != null ? {
          for k, v in element.value : k => jsonencode(v)
        } : null
      }
    }
  ]

  # Handle nested schedule object transformation if it exists (same as saved query)
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
      id = local.dashboard_id
    },
    # Base config without schedule and elements
    { for k, v in local.base_config : k => v if !contains(["schedule", "elements"], k) },
    # Add transformed elements with proper provider schema structure
    { elements = local.transformed_elements },
    # Add transformed schedule if it exists
    local.transformed_schedule != null ? { schedule = local.transformed_schedule } : {},
    var.additional_config
  )
}
