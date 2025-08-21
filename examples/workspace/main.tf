terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Configuration options
  #server_url ="https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/main/app/api/v1"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "main"
  server_url = "https://app.cribl-playground.cloud"
}

resource "criblio_workspace" "my_workspace" {
  alias           = "Production Environment"
  description     = "Main production workspace for customer data processing"
  organization_id = "beautiful-nguyen-y8y4azd"
  region          = "us-west-2"
  tags = [
    "production",
    "customer-data",
  ]
  workspace_id = "main"
}
