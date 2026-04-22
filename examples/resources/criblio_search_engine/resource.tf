resource "criblio_search_engine" "my_searchengine" {
  description = "Primary ingestion engine for the default local Dataset."
  id          = "local_ingest_primary"
  tier_size   = "small"
}