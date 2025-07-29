# Basic Example using Shared Provider Configuration
# This example shows how to reference the shared provider setup

# Import the shared provider configuration
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "~> 1.0"
    }
  }
}

# Use the bootstrap token module
module "bootstrap_token" {
  source = "../examples/cribl-bootstrap-token-module"
  
  client_id     = var.cribl_client_id
  client_secret = var.cribl_client_secret
  organization  = var.cribl_organization
  workspace     = var.cribl_workspace
  group         = var.cribl_group
}

# Use the worker installation module
module "cribl_worker" {
  source = "../examples/cribl-worker-module"
  
  organization = var.cribl_organization
  workspace    = var.cribl_workspace
  group        = var.cribl_group
  auth_token   = module.bootstrap_token.bootstrap_token
}

# Provider configuration (inherits from shared variables)
provider "criblio" {
  client_id       = var.cribl_client_id
  client_secret   = var.cribl_client_secret
  organization_id = var.cribl_organization
  workspace_id    = var.cribl_workspace
}

# Output the results
output "bootstrap_token" {
  description = "Retrieved bootstrap token"
  value       = module.bootstrap_token.bootstrap_token
  sensitive   = true
}

output "installation_script" {
  description = "Complete Cribl worker installation script"
  value       = module.cribl_worker.user_data_script
  sensitive   = true
} 