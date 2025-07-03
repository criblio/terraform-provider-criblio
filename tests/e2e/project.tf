resource "criblio_project" "my_project" {
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
  value = criblio_project.my_project
}
