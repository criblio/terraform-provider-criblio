resource "criblio_search_engine" "my_searchengine" {
  description = "My Search Engine TF"
  tier_size   = "2xlarge"
  id          = "my_search_engine_tf"
}

output "search_engine" {
  value = criblio_search_engine.my_searchengine
}

output "search_engine_id" {
  value = criblio_search_engine.my_searchengine.id
}
