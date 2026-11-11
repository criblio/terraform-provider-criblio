resource "criblio_search_dataset_acl" "logs" {
  dataset_id = "test_1"

  users = {
    "google-oauth2|115928590796700585084" = "DatasetRead"
  }

  teams = {
    "test-team" = "DatasetRead"
  }
}

data "criblio_search_dataset_acl" "main" {
  dataset_id = criblio_search_dataset_acl.logs.dataset_id
  depends_on = [criblio_search_dataset_acl.logs]
}
