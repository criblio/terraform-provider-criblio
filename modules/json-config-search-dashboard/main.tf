resource "criblio_search_dashboard" "dashboard_from_json" {
  id         = local.dashboard_id
  name       = try(local.final_config.name, "Terraform Dashboard")
  created    = try(local.final_config.created, 1695403200)
  created_by = try(local.final_config.created_by, "terraform")
  modified   = try(local.final_config.modified, 1695403200)
  elements   = try(local.final_config.elements, [])

  # Optional fields
  cache_ttl_seconds    = try(local.final_config.cache_ttl_seconds, null)
  category             = try(local.final_config.category, null)
  description          = try(local.final_config.description, null)
  pack_id              = try(local.final_config.pack_id, null)
  refresh_rate         = try(local.final_config.refresh_rate, null)
  resolved_dataset_ids = try(local.final_config.resolved_dataset_ids, null)

  # Schedule block if provided
  schedule = try(local.final_config.schedule, null) != null ? {
    cron_schedule = local.final_config.schedule.cron_schedule
    enabled       = local.final_config.schedule.enabled
    keep_last_n   = local.final_config.schedule.keep_last_n
    tz            = local.final_config.schedule.tz
    notifications = {
      disabled = local.final_config.schedule.notifications.disabled
    }
  } : null
}
