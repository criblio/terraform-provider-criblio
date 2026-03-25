resource "criblio_certificate" "my_certificate" {
  ca          = "LS0tLS1CR...FLS0tLS0K"
  cert        = "LS0tLS1CR...FLS0tLS0K"
  description = "Short description of x509 certificate"
  group_id    = "default"
  id          = "myUniqueCertId"
  in_use = [
    "list",
    "of",
    "configurations",
  ]
  passphrase = "SecurityPassphrase"
  priv_key   = "dont-share-this-key"
}