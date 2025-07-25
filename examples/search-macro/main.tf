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

resource "criblio_search_macro" "my_searchmacro" {
  created     = 1753453443
  created_by  = "test_created_by"
  description = "test_description"
  id          = "test_macro"
  modified    = 1753453443
  replacement = "test_replacement"
  tags        = "test_tags"
}
