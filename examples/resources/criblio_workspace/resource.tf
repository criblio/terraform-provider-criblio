resource "criblio_workspace" "my_workspace" {
  alias           = "Production Environment"
  description     = "Main production workspace for customer data processing"
  organization_id = "...my_organization_id..."
  region          = "us-west-2"
  tags = [
    "production",
    "customer-data",
  ]
  workspace_id = "my-workspace"
}