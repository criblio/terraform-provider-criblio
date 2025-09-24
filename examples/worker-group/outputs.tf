output "group" {
  value = criblio_group.azure_eastus_stream_group
  description = "Azure East US Stream Group object"
}

output "hybrid_worker_group" {
  value = criblio_group.hybrid_worker_group
  description = "Hybrid Worker Group object"
}

output "hybrid_worker_group_bootstrap" {
  value = module.hybrid_worker_group_bootstrap.user_data_script
  description = "Hybrid Worker Group User Data for bootstrap"
  sensitive   = true
  #Marked as sensitive because it contains the leader auth token
}