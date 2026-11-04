resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  description    = "Break HTTP access logs and Cribl API JSON arrays into events"
  group_id       = "Cribl"
  id             = "access-logs-v1"
  lib            = "custom"
  min_raw_length = 512
  rules = [
    {
      condition           = "/GET|POST|PUT|DELETE/.test(_raw)"
      delimiter           = ","
      delimiter_regex     = "/\\t/"
      disabled            = false
      escape_char         = "\\"
      event_breaker_regex = "/\\n(?=\\S)/"
      fields = [
        {
          name  = "source"
          value = "\"nginx_access\""
        }
      ]
      fields_line_regex   = "/^#[Ff]ields[:]?\\s+(.*)/"
      header_line_regex   = "/^#/"
      max_event_bytes     = 65536
      name                = "nginx-access"
      parser_enabled      = false
      quote_char          = "\""
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
    },
    {
      name                   = "cribl-api-events"
      condition              = "source.includes(\"/api/v1/system/logs/\")"
      type                   = "json_array"
      json_array_field       = "items.events"
      parent_fields_to_copy  = ["source"]
      json_extract_all       = false
      timestamp_anchor_regex = "/\"time\":/"
      timestamp = {
        type   = "auto"
        length = 150
      }
    }
  ]
  tags = "nginx,access,prod"
}
