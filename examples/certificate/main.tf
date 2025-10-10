resource "criblio_certificate" "my_certificate" {


  ca       = file("${path.module}/server.crt")
  cert     = file("${path.module}/server.crt")
  priv_key = file("${path.module}/server.key")

  description = "Demo x509 certificate for Cribl configuration"
  id          = "my-demo-cert-001"
  in_use      = ["wef-prod"]

  passphrase = "SecurityPassphrase"
}

data "criblio_certificates" "my_certificates" {

}

/*
output "my_certificates" {
  value     = data.criblio_certificates.my_certificates
  sensitive = true
}
*/


