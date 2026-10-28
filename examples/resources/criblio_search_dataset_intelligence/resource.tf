# The dataset must already exist; reference criblio_search_dataset.<name>.id
# when Terraform manages it. Requires agentic Search and accepted AI consent.
resource "criblio_search_dataset_intelligence" "example" {
  dataset_id = "application_logs"
}
