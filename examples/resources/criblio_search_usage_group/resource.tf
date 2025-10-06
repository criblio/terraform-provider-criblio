resource "criblio_search_usage_group" "my_searchusagegroup" {
  coordinator_heap_memory_limit = 2048
  description                   = "My Usage Group Description"
  enabled                       = false
  id                            = "myUsageGroupId"
  rules                         = "{ \"see\": \"documentation\" }"
  users_count                   = 1
}