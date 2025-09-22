# Example 1: Cribl HTTP Source from JSON file
module "cribl_http_source" {
  source = "../../modules/json-config-source-cribl-http"

  group_id       = "default"
  json_file_path = "./configs/cribl_http_source.json"
}

# Example 2: Syslog Source from JSON file
module "syslog_source" {
  source = "../../modules/json-config-source-syslog"

  group_id       = "default"
  json_file_path = "./configs/syslog_source.json"
}

# Example 3: Regular HTTP Source from JSON file
module "http_source" {
  source = "../../modules/json-config-source-http"

  group_id       = "default"
  json_file_path = "./configs/http_source.json"
}

# Example 4: TCP Source from JSON file
module "tcp_source" {
  source = "../../modules/json-config-source-tcp"

  group_id       = "default"
  json_file_path = "./configs/tcp_source.json"
}

# Example 5: Cribl TCP Source from JSON file
module "cribl_tcp_source" {
  source = "../../modules/json-config-source-cribl-tcp"

  group_id       = "default"
  json_file_path = "./configs/cribl_tcp_source.json"
}

# Example 6: OpenTelemetry Source from JSON file
module "opentelemetry_source" {
  source = "../../modules/json-config-source-open-telemetry"

  group_id       = "default"
  json_file_path = "./configs/opentelemetry_source.json"
}

# Outputs showing all created sources
output "all_sources_created" {
  value = {
    group_id = "default"
    sources = [
      {
        id   = module.cribl_http_source.source_id
        type = module.cribl_http_source.source_type
      },
      {
        id   = module.syslog_source.source_id
        type = module.syslog_source.source_type
      },
      {
        id   = module.http_source.source_id
        type = module.http_source.source_type
      },
      {
        id   = module.tcp_source.source_id
        type = module.tcp_source.source_type
      },
      {
        id   = module.cribl_tcp_source.source_id
        type = module.cribl_tcp_source.source_type
      },
      {
        id   = module.opentelemetry_source.source_id
        type = module.opentelemetry_source.source_type
      }
    ]
  }
}
