terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "1.22.2"
    }
  }
}

provider "criblio" {
  cloud_domain    = "..." # Optional - can use CRIBL_CLOUD_DOMAIN environment variable
  organization_id = "..." # Optional - can use CRIBL_ORGANIZATION_ID environment variable
  server_url      = "..." # Optional
  workspace_id    = "..." # Optional - can use CRIBL_WORKSPACE_ID environment variable
}