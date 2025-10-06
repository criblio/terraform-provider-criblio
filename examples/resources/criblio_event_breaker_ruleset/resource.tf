resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  description    = "Break HTTP access logs into events"
  group_id       = "Cribl"
  id             = "access-logs-v1"
  lib            = "custom"
  min_raw_length = 512
  rules = [
    {
      condition           = "/GET|POST|PUT|DELETE/.test(_raw)"
      disabled            = false
      event_breaker_regex = "/\\n(?=\\S)/"
      fields = [
        {
          name  = "source"
          value = "\"nginx_access\""
        }
      ]
      max_event_bytes     = 65536
      name                = "nginx-access"
      parser_enabled      = false
      should_use_data_raw = false
      timestamp = {
        format = "%d/%b/%Y:%H:%M:%S %z"
        length = 150
        type   = "format"
      }
      timestamp_anchor_regex = "/\\d{2}\\/[A-Za-z]{3}\\/\\d{4}:\\d{2}:\\d{2}:\\d{2}/"
      timestamp_earliest     = "-90days"
      timestamp_latest       = "+1day"
      timestamp_timezone     = "UTC"
      type                   = "regex"
    }
  ]
  tags = "nginx,access,prod"
}