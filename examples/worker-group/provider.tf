terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  client_id        = var.client_id
  client_secret    = var.client_secret
  organization_id  = var.organization_id
  workspace_id     = var.workspace_id
  cloud_domain     = var.cloud_domain
}