# Configure provider
#
# Authenticates for Cribl.Cloud using organization_id, workspace_id, client_id,
# and client_secret declared in variables.tf. Supply their values in terraform.tfvars.
#
# For other authentication methods, see:
# https://docs.cribl.io/cribl-as-code/terraform-auth/

terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = ">= 1.20.138"
    }
  }
}

provider "criblio" {
  organization_id = var.organization_id
  workspace_id    = var.workspace_id
  client_id       = var.client_id
  client_secret   = var.client_secret
}
