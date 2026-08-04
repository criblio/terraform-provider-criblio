resource "criblio_search_datatype" "example" {
  id              = "test"
  lib             = "custom"
  data_format     = "ndjson"
  max_event_bytes = 65536
  min_raw_length  = 0
  search_version  = "v2"
  tags            = ""

  automatic_extraction = {
    extraction_type = "json"
  }

  timestamp_extraction = {
    type         = "auto"
    anchor_regex = "/^/"
    timezone     = "UTC"
    earliest     = "0"
    latest       = "+10years"
    scan_depth   = 150
  }
}
