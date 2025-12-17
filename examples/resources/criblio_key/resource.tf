resource "criblio_key" "my_key" {
  algorithm   = "aes-256-cbc"
  description = "My Key Metadata"
  expires     = 1759325416
  group_id    = "default"
  id          = "key-001"
  keyclass    = 0
  kms         = "local"
  use_iv      = true
}