resource "criblio_pack_breakers" "my_packbreakers" {
  description    = "test"
  group_id       = "default"
  id             = "test_packbreakers"
  lib            = "custom"
  min_raw_length = 256
  pack           = criblio_pack.my_pack_breakers.id
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



resource "criblio_pack" "my_pack_breakers" {
  id           = "pack-breakers"
  group_id     = "default"
  description  = "Pack breakers"
  disabled     = true
  display_name = "Pack breakers"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

# Output the pack details to see the read-only attributes
output "pack_breakers_details" {
  value = criblio_pack.my_pack_breakers
}
