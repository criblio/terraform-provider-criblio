output "destination_id" {
  description = "The ID of the created destination"
  value       = criblio_destination.this.id
}

output "group_id" {
  description = "The group ID of the destination"
  value       = criblio_destination.this.group_id
}

output "destination_type" {
  description = "The type of destination created"
  value       = var.destination_type
}

output "destination_resource" {
  description = "The full destination resource (for advanced use cases)"
  value       = criblio_destination.this
  sensitive   = true
}