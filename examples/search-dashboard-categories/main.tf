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

resource "criblio_search_dashboard_category" "my_searchdashboardcategory" {
  description = "test"
  id          = "test_dashboard_category"
  is_pack     = true
  name        = "test_dashboard_category"
}

resource "criblio_search_dashboard_category" "my_searchdashboardcategory_not_pack" {
  description = "test"
  id          = "test_dashboard_category_not_pack"
  is_pack     = false
  name        = "test_dashboard_category_not_pack"
}