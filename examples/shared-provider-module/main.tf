# Shared Provider Module Example
# This pattern uses a module to standardize provider configuration

# Use the shared provider module
module "cribl_provider_config" {
  source = "../guides"
  
  # Pass provider variables
  cribl_client_id    = var.cribl_client_id
  cribl_client_secret = var.cribl_client_secret
  cribl_organization = var.cribl_organization
  cribl_workspace    = var.cribl_workspace
}

# Configure providers using module outputs (if the guides module provides them)
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "~> 1.0"
    }
  }
}

provider "criblio" {
  client_id       = var.cribl_client_id
  client_secret   = var.cribl_client_secret
  organization_id = var.cribl_organization
  workspace_id    = var.cribl_workspace
}

# Use your Cribl modules
module "bootstrap_token" {
  source = "../../module/test/cribl-bootstrap-token-module"
  
  client_id     = var.cribl_client_id
  client_secret = var.cribl_client_secret
  organization  = var.cribl_organization
  workspace     = var.cribl_workspace
  group         = var.cribl_group
}

module "cribl_worker" {
  source = "../../module/test/cribl-worker-module"
  
  organization = var.cribl_organization
  workspace    = var.cribl_workspace
  group        = var.cribl_group
  auth_token   = module.bootstrap_token.bootstrap_token
} 