resource "criblio_search_engine" "my_searchengine" {
  description = "My Search Engine TF"
  group_id    = "default_search"
  id          = "my_search_tf_engine"
  tier_size   = "small"
}

output "search_engine" {
  value = criblio_search_engine.my_searchengine
}

output "search_engine_id" {
  value = criblio_search_engine.my_searchengine.id
}
