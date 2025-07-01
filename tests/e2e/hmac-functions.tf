resource "criblio_hmac_function" "my_hmacfunction" {
  description       = "test hmac function"
  group_id          = "default"
  header_expression = "'hmac sha256 ' + C.Crypto.createHmac('test', C.Secret('yourSecret','text').value, 'sha256','hex')"
  header_name       = "signature"
  id                = "my_hmacfunction"
  lib               = "cribl"
  string_builders = [
    "true"
  ]
  string_delim = "true"
}

output "hmac_function" {
  value = criblio_hmac_function.my_hmacfunction
}
