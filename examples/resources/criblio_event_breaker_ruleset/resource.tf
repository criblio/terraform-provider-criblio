resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  description    = "...my_description..."
  group_id       = "myExistingGroupId"
  id             = "myUniqueBreakerIdToCRUD"
  lib            = "custom"
  min_raw_length = 75773.87
  rules = [
    {
      condition           = "...my_condition..."
      disabled            = false
      event_breaker_regex = "...my_event_breaker_regex..."
      fields = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      max_event_bytes     = 24876655.41
      name                = "...my_name..."
      parser_enabled      = true
      should_use_data_raw = false
      timestamp = {
        format = "...my_format..."
        length = 3.11
        type   = "format"
      }
      timestamp_anchor_regex = "...my_timestamp_anchor_regex..."
      timestamp_earliest     = "...my_timestamp_earliest..."
      timestamp_latest       = "...my_timestamp_latest..."
      timestamp_timezone     = "...my_timestamp_timezone..."
      type                   = "csv"
    }
  ]
  tags = "...my_tags..."
}