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
  value = criblio_project.my_project[0]
}
