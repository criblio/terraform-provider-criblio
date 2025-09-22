output "destination_id" {
  description = "The ID of the created destination"
  value       = criblio_destination.kinesis_from_json.id
}

output "destination_type" {
  description = "The type of the destination"
  value       = local.destination_type
}
