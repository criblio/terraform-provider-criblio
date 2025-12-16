resource "criblio_secret" "my_secret" {
  api_key     = "api_9f3c2a1b0e5d4c7a"
  description = "API key for ingestion service"
  group_id    = "default"
  id          = "sec-3f7a2c9d"
  password    = "app-password-123"
  secret_key  = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  secret_type = "text"
  tags        = "env:prod,team:security"
  username    = "ingest-bot"
  value       = "token-abc123"
}