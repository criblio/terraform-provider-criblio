resource "criblio_search_dataset_acl" "example" {
  dataset_id = "application_logs"

  users = {
    "google-oauth2|example-member-id" = "DatasetRead"
    "example-client-id@clients"       = "DatasetRead"
  }

  teams = {
    "engineering-team" = "DatasetRead"
  }
}
