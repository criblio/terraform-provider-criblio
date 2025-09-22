output "source_id" {
  description = "The ID of the created source"
  value       = criblio_source.http_from_json.id
}

output "source_type" {
  description = "The type of source that was created"
  value       = local.source_type
}

output "group_id" {
  description = "The group ID where the source was created"
  value       = criblio_source.http_from_json.group_id
}

output "parsed_config" {
  description = "The parsed and transformed JSON configuration"
  value       = local.final_config
  sensitive   = true
}
