resource "criblio_secret" "my_secret" {
  description = "API key for ingestion service"
  group_id    = "default"
  id          = "sec-test-001"
  secret_type = "text"
  tags        = "env:prod,team:security"
  value       = "token-abc123"
}

