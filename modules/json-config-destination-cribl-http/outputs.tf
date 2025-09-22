output "destination_id" {
  description = "The ID of the created Cribl HTTP destination"
  value       = criblio_destination.cribl_http_from_json.id
}

output "destination_type" {
  description = "The type of the created destination"
  value       = "cribl_http"
}

output "group_id" {
  description = "The group ID where the destination was created"
  value       = var.group_id
}

output "final_config" {
  description = "The final merged configuration used for the destination"
  value       = local.final_config
  sensitive   = true
}
