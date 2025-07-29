# Output the auth token for reference
output "auth_token" {
  description = "Auth token from distributed instance"
  value       = local.auth_token
  sensitive   = true
}

output "hybrid_config" {
  description = "Configuration for hybrid node deployment"
  value = local.hybrid_config
  sensitive = true
}

# Outputs from AWS modules
output "vpc_id" {
  description = "ID of the created VPC"
  value       = module.awsnetworks.vpcid
}

output "subnet_id" {
  description = "ID of the created subnet"
  value       = module.awsnetworks.subnet_id
}

output "security_group_id" {
  description = "ID of the network security group"
  value       = module.awsnetworks.security_group_id
}

output "linux_instance_ips" {
  description = "Public IP addresses of Linux instances"
  value       = module.linux_instances.ec2-ip
}

output "linux_instance_ids" {
  description = "Instance IDs of Linux instances"
  value       = module.linux_instances.instances
}

output "linux_security_group_id" {
  description = "Security group ID for Linux instances"
  value       = module.linux_instances.security_group
}

output "tokenvalue" {
  description = "Security group ID for Linux instances"
  value       = data.local_file.auth_token.content
}

output "debug_full_url" {
  value = "https://${local.workspace}-${local.cloud_tenant}.cribl${local.staging_suffix}.cloud/init/install-worker.sh?group=${local.group}&token=${data.local_file.auth_token.content}&user=cribl&user_group=cribl&install_dir=%2Fopt%2Fcribl"
  sensitive = false
} 