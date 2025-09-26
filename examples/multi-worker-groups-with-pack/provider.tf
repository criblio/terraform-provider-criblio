terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id    = "main"
  cloud_domain    = "cribl-playground.cloud"
}
