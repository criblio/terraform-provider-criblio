output "query_id" {
  description = "The ID of the created search saved query"
  value       = criblio_search_saved_query.saved_query_from_json.id
}

output "query_name" {
  description = "The name of the created search saved query"
  value       = criblio_search_saved_query.saved_query_from_json.name
}

output "final_config" {
  description = "The final merged configuration used for the saved query"
  value       = local.final_config
  sensitive   = true
}
