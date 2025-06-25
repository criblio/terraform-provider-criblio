terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "999.99.9" # This ensures that we only test with the locally built provider
    }
  }
}

provider "criblio" {
  server_url ="https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
}

