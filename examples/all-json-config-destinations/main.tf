# Example 1: Cribl Lake Destination from JSON file
module "cribl_lake_destination" {
  source = "../../modules/json-config-destination-cribl-lake"

  group_id       = "default"
  json_file_path = "./configs/cribl_lake_destination.json"
}

# Example 2: S3 Destination from JSON file  
module "s3_destination" {
  source = "../../modules/json-config-destination-s3"

  group_id       = "default"
  json_file_path = "./configs/s3_destination.json"
}

# Example 3: Splunk HEC Destination from JSON file
module "splunk_hec_destination" {
  source = "../../modules/json-config-destination-splunk-hec"

  group_id       = "default"
  json_file_path = "./configs/splunk_hec_destination.json"
}

# Example 4: Webhook Destination from JSON file
module "webhook_destination" {
  source = "../../modules/json-config-destination-webhook"

  group_id       = "default"
  json_file_path = "./configs/webhook_destination.json"
}

# Example 5: Kafka Destination from JSON file
module "kafka_destination" {
  source = "../../modules/json-config-destination-kafka"

  group_id       = "default"
  json_file_path = "./configs/kafka_destination.json"
}

# Example 6: Elasticsearch Destination from JSON file
module "elasticsearch_destination" {
  source = "../../modules/json-config-destination-elasticsearch"

  group_id       = "default"
  json_file_path = "./configs/elasticsearch_destination.json"
}

# Example 7: Kinesis Destination from JSON file
module "kinesis_destination" {
  source = "../../modules/json-config-destination-kinesis"

  group_id       = "default"
  json_file_path = "./configs/kinesis_destination.json"
}

# Example 8: Cribl HTTP Destination from JSON file
module "cribl_http_destination" {
  source = "../../modules/json-config-destination-cribl-http"

  group_id       = "default"
  json_file_path = "./configs/cribl_http_destination.json"
}

# Example 9: Cribl TCP Destination from JSON file
module "cribl_tcp_destination" {
  source = "../../modules/json-config-destination-cribl-tcp"

  group_id       = "default"
  json_file_path = "./configs/cribl_tcp_destination.json"
}

# Example 10: Syslog Destination from JSON file
module "syslog_destination" {
  source = "../../modules/json-config-destination-syslog"

  group_id       = "default"
  json_file_path = "./configs/syslog_destination.json"
}

# Example 11: OpenTelemetry Destination from JSON file
module "open_telemetry_destination" {
  source = "../../modules/json-config-destination-open-telemetry"

  group_id       = "default"
  json_file_path = "./configs/open_telemetry_destination.json"
}

# Example 12: CrowdStrike Next-Gen SIEM Destination from JSON file
module "crowdstrike_destination" {
  source = "../../modules/json-config-destination-crowdstrike-next-gen-siem"

  group_id       = "default"
  json_file_path = "./configs/crowdstrike_destination.json"
}

# Example 13: Prometheus Destination from JSON file
module "prometheus_destination" {
  source = "../../modules/json-config-destination-prometheus"

  group_id       = "default"
  json_file_path = "./configs/prometheus_destination.json"
}

# Outputs showing all created destinations
output "all_destinations_created" {
  value = {
    group_id = "default"
    destinations = [
      {
        id   = module.cribl_lake_destination.destination_id
        type = module.cribl_lake_destination.destination_type
      },
      {
        id   = module.s3_destination.destination_id
        type = module.s3_destination.destination_type
      },
      {
        id   = module.splunk_hec_destination.destination_id
        type = module.splunk_hec_destination.destination_type
      },
      {
        id   = module.webhook_destination.destination_id
        type = module.webhook_destination.destination_type
      },
      {
        id   = module.kafka_destination.destination_id
        type = module.kafka_destination.destination_type
      },
      {
        id   = module.elasticsearch_destination.destination_id
        type = module.elasticsearch_destination.destination_type
      },
      {
        id   = module.kinesis_destination.destination_id
        type = module.kinesis_destination.destination_type
      },
      {
        id   = module.cribl_http_destination.destination_id
        type = module.cribl_http_destination.destination_type
      },
      {
        id   = module.cribl_tcp_destination.destination_id
        type = module.cribl_tcp_destination.destination_type
      },
      {
        id   = module.syslog_destination.destination_id
        type = module.syslog_destination.destination_type
      },
      {
        id   = module.open_telemetry_destination.destination_id
        type = module.open_telemetry_destination.destination_type
      },
      {
        id   = module.crowdstrike_destination.destination_id
        type = module.crowdstrike_destination.destination_type
      },
      {
        id   = module.prometheus_destination.destination_id
        type = module.prometheus_destination.destination_type
      }
    ]
  }
}
