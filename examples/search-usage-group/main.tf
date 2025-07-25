terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "determined-gian-gkh6kzw"
  workspace_id = "main"
}

resource "criblio_search_usage_group" "my_searchusagegroup" {
  coordinator_heap_memory_limit = 8
  description                   = "test"
  enabled                       = true
  id                            = "test_usage_group"
  rules                         = "{\"test\": \"test\"}"
  users_count                   = 10
}
