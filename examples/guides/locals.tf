# Define common variables for modules
locals {
  cloud_tenant = var.cloud_tenant
  staging_suffix = ""
  group = var.group-hybrid
  cribl_version = var.cribl_version
  workspace = var.workspace
  say_my_name = "azamora"  # Your identifier
  
  # Derived values - can use module output or file content
  clean_auth_token = coalesce(module.bootstrap_token.bootstrap_token, try(trimspace(data.local_file.auth_token.content), ""))
  
  # Template variables for user data script
  cribl_install_url = "https://${local.workspace}-${local.cloud_tenant}.cribl${local.staging_suffix}.cloud/init/install-worker.sh?group=${local.group}&token=${local.clean_auth_token}&user=cribl&user_group=cribl&install_dir=%2Fopt%2Fcribl"

  # Store auth token as a variable
  auth_token = trimspace(data.local_file.auth_token.content)
  
  # Hybrid deployment configuration
  hybrid_config = {
    auth_token = local.auth_token
    master_host = "0.0.0.0"  # From your example
    master_port = 4200        # From your example
    protocol    = "http2"     # From your example
  }
}