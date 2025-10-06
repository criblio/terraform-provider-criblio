resource "criblio_project" "my_project" {
  consumers = {
    # ...
  }
  description = "Default project configuration"
  destinations = [
    "Output:main",
  ]
  group_id = "Cribl"
  id       = "project-default"
  subscriptions = [
    "sub-default",
  ]
}