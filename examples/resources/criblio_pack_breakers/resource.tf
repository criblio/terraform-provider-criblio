resource "criblio_pack_breakers" "my_packbreakers" {
  description    = "test"
  group_id       = "default"
  id             = "test_packbreakers"
  lib            = "custom"
  min_raw_length = 256
  pack           = criblio_pack.breakers_pack.id
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
    },
    {
      condition              = "true"
      type                   = "csv"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone = "utc"
      max_event_bytes    = 1024000
      disabled           = false
      delimiter          = ","
      quote_char         = "\""
      escape_char        = "\""
      name               = "csv"
    },
    {
      name                   = "header"
      condition              = "true"
      type                   = "header"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 51200
      disabled            = false
      parser_enabled      = false
      should_use_data_raw = false
      delimiter_regex     = "/\\t/"
      fields_line_regex   = "/^#[Ff]ields[:]?\\s+(.*)/"
      header_line_regex   = "/^#/"
      null_field_val      = "-"
      clean_fields        = true
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
    },
    {
      name                   = "json_array"
      condition              = "true"
      type                   = "json_array"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 51200
      disabled            = false
      parser_enabled      = false
      should_use_data_raw = false
      json_extract_all    = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
    },
    {
      name                   = "json"
      condition              = "true"
      type                   = "json"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 51200
      disabled            = false
      parser_enabled      = false
      should_use_data_raw = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
    },
    {
      name                   = "timestamp"
      condition              = "true"
      type                   = "timestamp"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 51200
      disabled            = false
      parser_enabled      = false
      should_use_data_raw = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
    },
    {
      name                   = "aws_cloudtrail"
      condition              = "true"
      type                   = "aws_cloudtrail"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 51200
      disabled            = false
      parser_enabled      = false
      should_use_data_raw = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
    },
    {
      name                   = "aws_vpcflow"
      condition              = "true"
      type                   = "aws_vpcflow"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "auto"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 51200
      disabled            = false
      parser_enabled      = false
      should_use_data_raw = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
    }
  ]
  tags = "test"
}

resource "criblio_pack" "breakers_pack" {
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
  value = criblio_pack_breakers.my_packbreakers
}
