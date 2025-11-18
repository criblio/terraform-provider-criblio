resource "criblio_project" "my_project" {
  //count required for cribl internal testing
  //count is not required for most customer implementations
  count = var.onprem == false ? 1 : 0

  consumers = {
  }
  description = "test project"
  destinations = [
  ]
  group_id = "default"
  id       = "my_project"
  subscriptions = [
  ]
}

output "project" {
  //fancy logic required for cribl internal testing
  //fancy logic is not required for most customer implementations
  value = length(criblio_project.my_project) > 0 ? criblio_project.my_project[0] : null

  //value = criblio_project.my_project[0]
}
