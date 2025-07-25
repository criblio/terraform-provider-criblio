resource "criblio_pack_breakers" "my_packbreakers" {
  description    = "...my_description..."
  group_id       = "...my_group_id..."
  id             = "...my_id..."
  lib            = "custom"
  min_raw_length = 94618.96
  pack           = "...my_pack..."
  rules = [
    {
      condition           = "...my_condition..."
      disabled            = true
      event_breaker_regex = "...my_event_breaker_regex..."
      fields = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      max_event_bytes     = 101343288.08
      name                = "...my_name..."
      parser_enabled      = true
      should_use_data_raw = false
      timestamp = {
        format = "...my_format..."
        length = 9.13
        type   = "current"
      }
      timestamp_anchor_regex = "...my_timestamp_anchor_regex..."
      timestamp_earliest     = "...my_timestamp_earliest..."
      timestamp_latest       = "...my_timestamp_latest..."
      timestamp_timezone     = "...my_timestamp_timezone..."
      type                   = "aws_vpcflow"
    }
  ]
  tags = "...my_tags..."
}