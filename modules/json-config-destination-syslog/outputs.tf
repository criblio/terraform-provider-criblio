output "destination_id" {
  description = "The ID of the created Syslog destination"
  value       = criblio_destination.syslog_from_json.id
}

output "destination_type" {
  description = "The type of the created destination"
  value       = "syslog"
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
