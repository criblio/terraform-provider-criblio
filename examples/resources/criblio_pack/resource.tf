resource "criblio_pack" "my_pack" {
  allow_custom_functions = false
  author                 = "Cribl"
  description            = "Parsers and routes for Apache logs"
  disabled               = false
  display_name           = "Apache Logs Pack"
  exports = [
    "pipelines/parse_apache",
    "routes/apache_to_clickhouse",
  ]
  filename               = "observability-pack-1.2.3.crbl"
  force                  = false
  group_id               = "Cribl"
  id                     = "pack-apache-logs"
  inputs                 = 2
  is_disabled            = false
  min_log_stream_version = "4.10.0"
  outputs                = 1
  source                 = "https://github.com/org/packs/apache"
  spec                   = "1.0.0"
  tags = {
    data_type = [
      "..."
    ]
    domain = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    technology = [
      "..."
    ]
  }
  version = "2.3.1"
}