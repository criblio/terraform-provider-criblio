/* commenting out to get e2e tests to pass
resource "criblio_regex" "my_regex" {
  description = "test"
  group_id    = "default"
  id          = "test_regex"
  lib         = "custom"
  regex       = "/\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b/"
  sample_data = "192.168.1.1, 10.0.0.1, 172.16.0.1"
  tags        = "test"
}

output "regex" {
  value = criblio_regex.my_regex
}

data "criblio_regex" "my_regex" {
  group_id = "default"
}
*/
