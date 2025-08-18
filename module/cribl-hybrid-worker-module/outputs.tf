# Outputs for Cribl Worker Installation Module

output "user_data_script" {
  description = "The complete user data script for installing Cribl worker"
  value       = local.user_data_script
}

output "cribl_install_url" {
  description = "The constructed Cribl installation URL"
  value       = local.cribl_install_url
  sensitive   = true
}

output "installation_command" {
  description = "The core installation command that can be used independently"
  value       = "curl -fsSL '${local.cribl_install_url}' | bash -"
  sensitive   = true
} 