resource "criblio_regex" "my_regex" {
  description = "Regex for IPv4 addresses"
  group_id    = "default"
  id          = "regexEntry1"
  lib         = "custom"
  regex       = "\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b"
  sample_data = "192.168.1.1, 10.0.0.1"
  tags        = "network,ip"
}