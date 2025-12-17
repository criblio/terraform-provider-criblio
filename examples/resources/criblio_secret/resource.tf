resource "criblio_secret" "my_secret" {
  description = "API key for ingestion service"
  group_id    = "default"
  id          = "sec-3f7a2c9d"
  secret_type = "text"
  tags        = "env:prod,team:security"
  value       = "token-abc123"
}