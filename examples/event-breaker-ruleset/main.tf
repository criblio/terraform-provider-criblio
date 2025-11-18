resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  //count required for cribl internal testing
  //count is not required for most customer implementations
  count = var.onprem == false ? 1 : 0

  description    = "test"
  group_id       = "default"
  id             = "test_eventbreakerruleset"
  lib            = "custom"
  min_raw_length = 256
  rules = [
    {
      condition           = "PASS_THROUGH_SOURCE_TYPE"
      disabled            = false
      fields              = []
      max_event_bytes     = 51200
      name                = "test"
      parser_enabled      = false
      should_use_data_raw = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
      timestamp = {
        length = 150
        type   = "auto"
      }
      timestamp_anchor_regex = "/^/"
      timestamp_earliest     = "-420weeks"
      timestamp_latest       = "+1week"
      timestamp_timezone     = "local"
      type                   = "regex"
    }
  ]
  tags = "test"
}

output "event_breaker_ruleset" {
  value = criblio_event_breaker_ruleset.my_eventbreakerruleset
}
