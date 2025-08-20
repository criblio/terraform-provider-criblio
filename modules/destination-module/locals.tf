locals {
  # Base Cribl Lake configuration
  base_cribl_lake_config = {
    id          = var.destination_id
    type        = "cribl_lake"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Cribl Lake specific settings
    dest_path                         = var.dest_path
    format                            = var.format
    compress                          = var.compress
    max_file_size_mb                  = 32
    max_open_files                    = 100
    write_high_water_mark             = 64
    on_backpressure                   = "block"
    deadletter_enabled                = false
    on_disk_full_backpressure         = "block"
    max_file_open_time_sec            = 300
    max_file_idle_time_sec            = 30
    verify_permissions                = true
    max_closing_files_to_backpressure = 100
    max_concurrent_file_parts         = 1
    empty_dir_cleanup_sec             = 300
    max_retry_num                     = 20
  }

  # Base S3 configuration
  base_s3_config = {
    id          = var.destination_id
    type        = "s3"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # S3 specific settings
    bucket                    = var.bucket
    region                    = var.region
    dest_path                 = var.dest_path
    format                    = var.format
    compress                  = var.compress
    aws_authentication_method = "auto"
    base_file_name            = "CriblOut"
    file_name_suffix          = var.format == "json" ? "'.json.gz'" : "'.gz'"
    max_file_size_mb          = 32
    max_open_files            = 100
    write_high_water_mark     = 64
    on_backpressure           = "block"
    max_file_open_time_sec    = 300
    max_file_idle_time_sec    = 30
    verify_permissions        = true
    max_retry_num             = 20
  }

  # Base Splunk HEC configuration
  base_splunk_hec_config = {
    id          = var.destination_id
    type        = "splunk_hec"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Splunk HEC specific settings
    token                 = var.token
    validate_server_certs = true
    compress              = var.compress == "gzip" ? true : false
    format                = var.format
    max_events_per_req    = 0
    flush_period_sec      = 1
    max_retry_events      = 1000
    on_backpressure       = "block"
    timeout               = 30
    honor_keep_alive      = true
    request_concurrency   = 5
    max_retry_num         = 10
    ack_enabled           = false
    load_balanced         = false          # Set to false to use single URL
    url                   = var.splunk_url # Single URL when load_balanced = false
  }

  # Base Cribl HTTP configuration
  base_cribl_http_config = {
    id          = var.destination_id
    type        = "cribl_http"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Cribl HTTP specific settings
    method                = var.method
    headers               = {}
    compress              = var.compress == "gzip" ? true : false
    format                = var.format
    flush_period_sec      = 1
    max_events_per_req    = 0
    on_backpressure       = "block"
    timeout               = 30
    honor_keep_alive      = true
    request_concurrency   = 5
    max_retry_num         = 10
    validate_server_certs = true
    auth_type             = "none"
    load_balanced         = false   # Set to false to use single URL
    url                   = var.url # Single URL when load_balanced = false
  }

  # Base Cribl TCP configuration
  base_cribl_tcp_config = {
    id          = var.destination_id
    type        = "cribl_tcp"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Cribl TCP specific settings
    host                     = "localhost"
    port                     = 9999
    tls                      = { disabled = true }
    connection_timeout       = 10000
    write_timeout            = 60000
    on_backpressure          = "block"
    idle_timeout_sec         = 300
    connection_stats_enabled = false
  }

  # Base Syslog configuration
  base_syslog_config = {
    id          = var.destination_id
    type        = "syslog"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Syslog specific settings
    host               = "localhost"
    port               = 514
    protocol           = "udp"
    facility           = 16
    severity           = 6
    app_name           = "cribl"
    format             = "rfc5424"
    on_backpressure    = "block"
    connection_timeout = 10000
    write_timeout      = 60000
  }


  # Base Kafka configuration
  base_kafka_config = {
    id          = var.destination_id
    type        = "kafka"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Kafka specific settings
    brokers            = ["localhost:9092"]
    topic              = "cribl"
    format             = var.format
    compress           = var.compress == "gzip" ? "gzip" : "none"
    flush_period_sec   = 1
    max_events_per_req = 0
    on_backpressure    = "block"
    timeout            = 30
    max_retry_num      = 10
    acks               = "all"
    batch_size         = 16384
    linger_ms          = 0
  }

  # Base Kinesis configuration
  base_kinesis_config = {
    id          = var.destination_id
    type        = "kinesis"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Kinesis specific settings
    stream_name               = "cribl-stream"
    region                    = var.region
    format                    = var.format
    compress                  = var.compress == "gzip" ? true : false
    flush_period_sec          = 1
    max_events_per_req        = 500
    on_backpressure           = "block"
    timeout                   = 30
    max_retry_num             = 10
    aws_authentication_method = "auto"
    partition_key             = "_partition"
  }

  # Base Elasticsearch configuration
  base_elasticsearch_config = {
    id          = var.destination_id
    type        = "elasticsearch"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Elasticsearch specific settings
    url                   = var.url
    index                 = "cribl"
    format                = "json"
    flush_period_sec      = 1
    max_events_per_req    = 0
    on_backpressure       = "block"
    timeout               = 30
    honor_keep_alive      = true
    request_concurrency   = 5
    max_retry_num         = 10
    validate_server_certs = true
    auth_type             = "none"
    version               = "7.x"
  }

  # Base Webhook configuration
  base_webhook_config = {
    id          = var.destination_id
    type        = "webhook"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Webhook specific settings from examples
    auth_type                         = "none"
    compress                          = var.compress == "gzip" ? true : false
    concurrency                       = 5
    dns_resolve_period_sec            = 600
    exclude_self                      = false
    extra_http_headers                = []
    failed_request_logging_mode       = "none"
    flush_period_sec                  = 1
    format                            = var.format
    keep_alive                        = true
    load_balance_stats_period_sec     = 300
    load_balanced                     = true
    max_payload_events                = 0
    max_payload_size_kb               = 4096
    method                            = var.method
    on_backpressure                   = "block"
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings           = []
    safe_headers                      = []
    timeout_sec                       = 30
    urls = [
      {
        url = var.url # Required - will fail if empty
      }
    ]
    use_round_robin_dns = false
  }

  # Transform configs - remove nulls and merge with custom
  cribl_lake_config = var.destination_type == "cribl_lake" ? {
    for k, v in merge(local.base_cribl_lake_config, var.custom_config) : k => v
    if v != null
  } : null

  s3_config = var.destination_type == "s3" ? {
    for k, v in merge(local.base_s3_config, var.custom_config) : k => v
    if v != null
  } : null

  splunk_hec_config = var.destination_type == "splunk_hec" ? {
    for k, v in merge(local.base_splunk_hec_config, var.custom_config) : k => v
    if v != null
  } : null

  cribl_http_config = var.destination_type == "cribl_http" ? {
    for k, v in merge(local.base_cribl_http_config, var.custom_config) : k => v
    if v != null
  } : null

  cribl_tcp_config = var.destination_type == "cribl_tcp" ? {
    for k, v in merge(local.base_cribl_tcp_config, var.custom_config) : k => v
    if v != null
  } : null

  syslog_config = var.destination_type == "syslog" ? {
    for k, v in merge(local.base_syslog_config, var.custom_config) : k => v
    if v != null
  } : null

  kafka_config = var.destination_type == "kafka" ? {
    for k, v in merge(local.base_kafka_config, var.custom_config) : k => v
    if v != null
  } : null

  kinesis_config = var.destination_type == "kinesis" ? {
    for k, v in merge(local.base_kinesis_config, var.custom_config) : k => v
    if v != null
  } : null

  elasticsearch_config = var.destination_type == "elasticsearch" ? {
    for k, v in merge(local.base_elasticsearch_config, var.custom_config) : k => v
    if v != null
  } : null

  webhook_config = var.destination_type == "webhook" ? {
    for k, v in merge(local.base_webhook_config, var.custom_config) : k => v
    if v != null
  } : null

  # Base Open Telemetry configuration
  base_open_telemetry_config = {
    id          = var.destination_id
    type        = "open_telemetry"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Open Telemetry specific settings
    endpoint                          = var.url
    protocol                          = "http"
    otlp_version                      = "1.3.1"
    auth_type                         = "none"
    compress                          = var.compress == "gzip" ? "gzip" : "none"
    http_compress                     = var.compress == "gzip" ? "gzip" : "none"
    concurrency                       = 5
    flush_period_sec                  = 1
    max_payload_size_kb               = 4096
    timeout_sec                       = 30
    keep_alive                        = true
    keep_alive_time                   = 30
    connection_timeout                = 10
    on_backpressure                   = "block"
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    failed_request_logging_mode       = "none"
    metadata                          = []
    extra_http_headers                = []
    response_retry_settings           = []
    safe_headers                      = []
  }

  # Base CrowdStrike Next-Gen SIEM configuration
  base_crowdstrike_next_gen_siem_config = {
    id          = var.destination_id
    type        = "crowdstrike_next_gen_siem"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # CrowdStrike Next-Gen SIEM specific settings
    url                               = var.url
    token                             = var.token
    auth_type                         = "manual"
    format                            = "JSON"
    compress                          = var.compress == "gzip" ? true : false
    concurrency                       = 5
    flush_period_sec                  = 1
    max_payload_events                = 0
    max_payload_size_kb               = 4096
    on_backpressure                   = "block"
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    failed_request_logging_mode       = "none"
    timeout_sec                       = 30
    use_round_robin_dns               = false
    extra_http_headers                = []
    response_retry_settings           = []
    safe_headers                      = []
  }

  open_telemetry_config = var.destination_type == "open_telemetry" ? {
    for k, v in merge(local.base_open_telemetry_config, var.custom_config) : k => v
    if v != null
  } : null

  crowdstrike_next_gen_siem_config = var.destination_type == "crowdstrike_next_gen_siem" ? {
    for k, v in merge(local.base_crowdstrike_next_gen_siem_config, var.custom_config) : k => v
    if v != null
  } : null

  # Base Prometheus configuration
  base_prometheus_config = {
    id          = var.destination_id
    type        = "prometheus"
    description = var.description != "" ? var.description : null
    disabled    = var.disabled
    streamtags  = var.streamtags
    pipeline    = var.pipeline

    # Prometheus specific settings
    url                       = var.url
    auth_type                 = "none"
    concurrency               = 5
    flush_period_sec          = 10
    max_payload_size_kb       = 1024
    max_payload_events        = 1000
    compress                  = "none"
    timeout_sec               = 30
    on_backpressure           = "queue"
    use_round_robin_dns       = false
    keep_alive                = true
    add_timestamp             = true
    write_timestamp_precision = "ms"
  }

  # Merge custom config for Prometheus
  prometheus_config = var.destination_type == "prometheus" ? {
    for k, v in merge(local.base_prometheus_config, var.custom_config) : k => v
    if v != null
  } : null
}