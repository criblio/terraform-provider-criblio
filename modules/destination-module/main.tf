resource "criblio_destination" "this" {
  id       = var.destination_id
  group_id = var.group_id

  # Cribl Lake
  output_cribl_lake = var.destination_type == "cribl_lake" ? local.cribl_lake_config : null

  # S3
  output_s3 = var.destination_type == "s3" ? local.s3_config : null

  # Splunk HEC
  output_splunk_hec = var.destination_type == "splunk_hec" ? local.splunk_hec_config : null

  # Cribl HTTP
  output_cribl_http = var.destination_type == "cribl_http" ? local.cribl_http_config : null

  # Cribl TCP
  output_cribl_tcp = var.destination_type == "cribl_tcp" ? local.cribl_tcp_config : null

  # Syslog
  output_syslog = var.destination_type == "syslog" ? local.syslog_config : null

  # Kafka
  output_kafka = var.destination_type == "kafka" ? local.kafka_config : null

  # Kinesis
  output_kinesis = var.destination_type == "kinesis" ? local.kinesis_config : null

  # Elasticsearch (correct name is elastic)
  output_elastic = var.destination_type == "elasticsearch" ? local.elasticsearch_config : null

  # Webhook
  output_webhook = var.destination_type == "webhook" ? local.webhook_config : null

  # Open Telemetry
  output_open_telemetry = var.destination_type == "open_telemetry" ? local.open_telemetry_config : null

  # CrowdStrike Next-Gen SIEM
  output_crowdstrike_next_gen_siem = var.destination_type == "crowdstrike_next_gen_siem" ? local.crowdstrike_next_gen_siem_config : null

  # Prometheus
  output_prometheus = var.destination_type == "prometheus" ? local.prometheus_config : null

  lifecycle {
    create_before_destroy = true
  }
}