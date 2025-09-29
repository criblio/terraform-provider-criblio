resource "criblio_search_usage_group" "my_searchusagegroup" {
  coordinator_heap_memory_limit = 8.02
  description                   = "...my_description..."
  enabled                       = true
  id                            = "myUniqueUsageGroupToCRUD"
  rules                         = "{ \"see\": \"documentation\" }"
  users_count                   = 7.77
}