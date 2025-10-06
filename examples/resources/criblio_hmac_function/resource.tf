resource "criblio_hmac_function" "my_hmacfunction" {
  description       = "Compute HMAC-SHA256 over method:path:timestamp"
  group_id          = "Cribl"
  header_expression = "hmac('sha256', secret, parts.join(':'))"
  header_name       = "X-Signature"
  id                = "hmac-default"
  lib               = "cribl"
  string_builders = [
    "request.method",
    "request.path",
    "request.headers['x-timestamp']",
  ]
  string_delim = ":"
}