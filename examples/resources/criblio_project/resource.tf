resource "criblio_project" "my_project" {
  consumers = {
    # ...
  }
  description = "...my_description..."
  destinations = [
    "..."
  ]
  group_id = "myExistingGroupId"
  id       = "myUniqueProjectIdToCRUD"
  subscriptions = [
    "..."
  ]
}