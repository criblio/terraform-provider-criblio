/*
This is commented out because creating a workspace is expensive operation and we are just using the data source to get the workspaces for acceptance tests.
resource "criblio_workspace" "my_workspace" {
  alias           = "Production Environment Test"
  description     = "My new workspace for production customer data processing test"
  organization_id = "beautiful-nguyen-y8y4azd"
  region          = "us-west-2"
  tags = [
    "production",
    "customer-data",
    "test"
  ]
  workspace_id = "testworkspace"
}
*/
data "criblio_workspaces" "my_workspaces" {
  organization_id = "beautiful-nguyen-y8y4azd"
}

output "my_workspaces" {
  value = data.criblio_workspaces.my_workspaces
}
