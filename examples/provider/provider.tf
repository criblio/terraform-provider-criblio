terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "1.0.45"
    }
  }
}

provider "criblio" {
  # Configuration options
  client_id = var.cribl_client_id
  client_secret = var.cribl_client_secret
  organization_id = var.cloud_tenant
  workspace_id = var.workspace
}