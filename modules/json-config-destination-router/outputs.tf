output "destination_id" {
  description = "The ID of the created router destination"
  value       = criblio_destination.router_from_json.output_router.id
}

output "destination_type" {
  description = "The type of the created destination"
  value       = criblio_destination.router_from_json.output_router.type
}
