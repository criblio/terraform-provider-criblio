# modules/cribl-source/outputs.tf (or source-module/outputs.tf in your case)
output "source_id" {
  description = "The ID of the created source"
  value       = criblio_source.this.id
}

output "group_id" {
  description = "The group ID of the source"
  value       = criblio_source.this.group_id
}

output "source_type" {
  description = "The type of source created"
  value       = var.source_type
}

output "port" {
  description = "The port configured for the source"
  value       = var.port
}