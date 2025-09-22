output "dashboard_id" {
  description = "The ID of the created search dashboard"
  value       = criblio_search_dashboard.dashboard_from_json.id
}

output "dashboard_name" {
  description = "The name of the created search dashboard"
  value       = criblio_search_dashboard.dashboard_from_json.name
}

output "final_config" {
  description = "The final merged configuration used for the dashboard"
  value       = local.final_config
  sensitive   = true
}
