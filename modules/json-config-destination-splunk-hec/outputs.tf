output "destination_id" {
  description = "The ID of the created destination"
  value       = criblio_destination.splunk_hec_from_json.id
}

output "destination_type" {
  description = "The type of the destination"
  value       = local.destination_type
}
