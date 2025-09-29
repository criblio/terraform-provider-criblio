# Outputs for Cribl Worker Installation Module

output "user_data_script" {
  description = "The complete user data script for installing a Cribl Cloud hybrid worker"
  value       = local.user_data_script
  # The user data script contains the leader auth token
  sensitive = true
}
