resource "criblio_certificate" "my_certificate" {
  ca               = "...my_ca..."
  ca_path          = "...my_ca_path..."
  cert             = "...my_cert..."
  cert_expiry_date = "...my_cert_expiry_date..."
  cert_path        = "...my_cert_path..."
  description      = "...my_description..."
  group_id         = "default"
  id               = "cert-001"
  passphrase       = "...my_passphrase..."
  passphrase_path  = "...my_passphrase_path..."
  priv_key         = "...my_priv_key..."
  priv_key_path    = "...my_priv_key_path..."
}