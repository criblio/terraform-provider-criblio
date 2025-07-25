resource "criblio_pack_breakers" "my_packbreakers" {
  description    = "my_description"
  group_id       = "default"
  id             = "my_id"
  lib            = "custom"
  min_raw_length = 94618.96
  pack           = criblio_pack.breakers_pack.id
  rules = [
    {
      condition           = "my_condition"
      disabled            = true
      event_breaker_regex = "my_event_breaker_regex"
      fields = [
        {
          name  = "my_name"
          value = "my_value"
        }
      ]
      max_event_bytes     = 101343288.08
      name                = "my_name"
      parser_enabled      = true
      should_use_data_raw = false
      timestamp = {
        format = "my_format"
        length = 9.13
        type   = "current"
      }
      timestamp_anchor_regex = "/\\s(?=\\d{10}\\s\\d{10}\\s\\w)/"
      timestamp_timezone     = "utc"
      type                   = "aws_vpcflow"
    }
  ]
  tags = "my_tags"
}

resource "criblio_pack" "breakers_pack" {
  id           = "pack-with-breakers"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

# Output the pack details to see the read-only attributes
output "pack_breakers_details" {
  value = criblio_pack.my_pack
}
