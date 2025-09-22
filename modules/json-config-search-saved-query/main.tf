resource "criblio_search_saved_query" "saved_query_from_json" {
  id          = local.query_id
  name        = try(local.final_config.name, "Terraform Query")
  query       = try(local.final_config.query, "dataset=\"*\" | limit 10")
  description = try(local.final_config.description, null)
  earliest    = try(local.final_config.earliest, null)
  latest      = try(local.final_config.latest, null)
  is_private  = try(local.final_config.is_private, null)

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
