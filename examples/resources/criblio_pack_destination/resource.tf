resource "criblio_pack_destination" "my_packdestination" {
  group_id = "Cribl"
  id       = "pack-out-s3"
  output_azure_blob = {
    add_id_to_stage_path = true
    auth_type            = "manual"
    automatic_schema     = true
    azure_cloud          = "AzurePublicCloud"
    base_file_name       = "`CriblOut`"
    certificate = {
      certificate_name = "azure-app-cert"
    }
    client_id               = "11111111-1111-1111-1111-111111111111"
    client_text_secret      = "azure-sp-client-secret"
    compress                = "gzip"
    compression_level       = "normal"
    connection_string       = "DefaultEndpointsProtocol=https;AccountName=criblstore;AccountKey=***REDACTED***;EndpointSuffix=core.windows.net"
    container_name          = "cribl-data"
    create_container        = true
    deadletter_enabled      = true
    deadletter_path         = "/var/lib/cribl/state/outputs/dead-letter"
    description             = "Write objects to Azure Blob Storage with date-based partitioning"
    dest_path               = "logs/ingest"
    empty_dir_cleanup_sec   = 600
    enable_page_checksum    = true
    enable_statistics       = true
    enable_write_page_index = true
    endpoint_suffix         = "core.windows.net"
    environment             = "main"
    file_name_suffix        = ".json.gz"
    format                  = "json"
    header_line             = "timestamp,host,message"
    id                      = "azure-blob-out"
    key_value_metadata = [
      {
        key   = "team"
        value = "platform"
      }
    ]
    max_concurrent_file_parts = 4
    max_file_idle_time_sec    = 30
    max_file_open_time_sec    = 300
    max_file_size_mb          = 64
    max_open_files            = 200
    max_retry_num             = 20
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "4MB"
    parquet_row_group_length  = 10000
    parquet_version           = "PARQUET_2_6"
    partition_expr            = "2024/01/15"
    pipeline                  = "default"
    remove_empty_dirs         = true
    should_log_invalid_rows   = true
    stage_path                = "/var/lib/cribl/state/outputs/staging"
    storage_account_name      = "criblstore"
    storage_class             = "Hot"
    streamtags = [
      "azure",
      "blob",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    tenant_id             = "00000000-0000-0000-0000-000000000000"
    text_secret           = "azure-connstr-secret"
    type                  = "azure_blob"
    write_high_water_mark = 256
  }
  output_azure_data_explorer = {
    add_id_to_stage_path = true
    additional_properties = [
      {
        key   = "format"
        value = "json"
      }
    ]
    certificate = {
      certificate_name = "adx-app-cert"
    }
    client_id          = "11111111-1111-1111-1111-111111111111"
    client_secret      = "***REDACTED***"
    cluster_url        = "https://mycluster.eastus.kusto.windows.net"
    compress           = "gzip"
    concurrency        = 8
    database           = "telemetry"
    deadletter_enabled = true
    description        = "Ingest data into Azure Data Explorer (Kusto)"
    environment        = "main"
    extent_tags = [
      {
        prefix = "ingestBy"
        value  = "source:app1"
      }
    ]
    file_name_suffix  = ".json.gz"
    flush_immediately = true
    flush_period_sec  = 1
    format            = "json"
    id                = "adx-out"
    ingest_if_not_exists = [
      {
        value = "batchId:2025-10-02T00:00Z"
      }
    ]
    ingest_mode               = "batching"
    ingest_url                = "https://ingest-mycluster.eastus.kusto.windows.net"
    is_mapping_obj            = false
    keep_alive                = true
    mapping_ref               = "my_table_mapping"
    max_concurrent_file_parts = 4
    max_file_idle_time_sec    = 30
    max_file_open_time_sec    = 300
    max_file_size_mb          = 64
    max_open_files            = 200
    max_payload_events        = 0
    max_payload_size_kb       = 4096
    oauth_endpoint            = "https://login.microsoftonline.com"
    oauth_type                = "clientSecret"
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    pipeline                  = "default"
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    remove_empty_dirs                 = true
    report_level                      = "failuresOnly"
    report_method                     = "queue"
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    retain_blob_on_success = true
    scope                  = "https://kusto.kusto.windows.net/.default"
    stage_path             = "/var/lib/cribl/state/outputs/staging"
    streamtags = [
      "azure",
      "adx",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    table       = "logs_raw"
    tenant_id   = "00000000-0000-0000-0000-000000000000"
    text_secret = "adx-client-secret"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec                = 30
    type                       = "azure_data_explorer"
    use_round_robin_dns        = true
    validate_database_settings = true
  }
  output_azure_eventhub = {
    ack                    = 1
    authentication_timeout = 10000
    backoff_rate           = 2
    brokers = [
      "myns.servicebus.windows.net:9093",
    ]
    connection_timeout = 10000
    description        = "Deliver events to Azure Event Hubs via Kafka protocol"
    environment        = "main"
    flush_event_count  = 1000
    flush_period_sec   = 1
    format             = "json"
    id                 = "eventhub-out"
    initial_backoff    = 1000
    max_back_off       = 60000
    max_record_size_kb = 768
    max_retries        = 5
    on_backpressure    = "block"
    pipeline           = "default"
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "100 MB"
    pq_max_size                = "10GB"
    pq_mode                    = "backpressure"
    pq_on_backpressure         = "block"
    pq_path                    = "/opt/cribl/state/queues"
    reauthentication_threshold = 60000
    request_timeout            = 60000
    sasl = {
      disabled  = false
      mechanism = "plain"
    }
    streamtags = [
      "azure",
      "eventhub",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    tls = {
      disabled            = false
      reject_unauthorized = true
    }
    topic = "app-events"
    type  = "azure_eventhub"
  }
  output_azure_logs = {
    api_url     = ".ods.opinsights.azure.com"
    auth_type   = "manual"
    compress    = true
    concurrency = 8
    description = "Send logs to Azure Log Analytics workspace"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    id                          = "azure-logs-out"
    keypair_secret              = "azure-log-analytics-keys"
    log_type                    = "Cribl"
    max_payload_events          = 0
    max_payload_size_kb         = 1024
    on_backpressure             = "block"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    resource_id                       = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm1"
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "azure",
      "loganalytics",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    type                = "azure_logs"
    use_round_robin_dns = true
    workspace_id        = "22222222-2222-2222-2222-222222222222"
    workspace_key       = "***REDACTED***"
  }
  output_chronicle = {
    api_version           = "...my_api_version..."
    authentication_method = "serviceAccount"
    compress              = true
    concurrency           = 25.61
    custom_labels = [
      {
        key          = "...my_key..."
        rbac_enabled = false
        value        = "...my_value..."
      }
    ]
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 0.11
    gcp_instance                = "...my_gcp_instance..."
    gcp_project_id              = "...my_gcp_project_id..."
    id                          = "...my_id..."
    ingestion_method            = "...my_ingestion_method..."
    log_text_field              = "...my_log_text_field..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 5.04
    max_payload_size_kb         = 2620.78
    namespace                   = "...my_namespace..."
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_backpressure_sec           = 8.78
    pq_max_buffer_size                = 961.48
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    pq_rate_per_sec                   = 5.18
    pq_strict_ordering                = false
    region                            = "...my_region..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 13.66
        http_status     = 575.57
        initial_backoff = 12382.43
        max_backoff     = 140551.67
      }
    ]
    safe_headers = [
      "..."
    ]
    service_account_credentials        = "...my_service_account_credentials..."
    service_account_credentials_secret = "...my_service_account_credentials_secret..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 19.58
      initial_backoff = 476612.3
      max_backoff     = 154167.96
      timeout_retry   = true
    }
    timeout_sec           = 4865019203906175
    total_memory_limit_kb = 6.13
    type                  = "chronicle"
    use_round_robin_dns   = true
  }
  output_click_house = {
    async_inserts    = true
    auth_header_expr = "`Bearer ${token}`"
    auth_type        = "basic"
    column_mappings = [
      {
        column_name             = "timestamp"
        column_type             = "DateTime64(3)"
        column_value_expression = "toDateTime64(ts, 3)"
      }
    ]
    compress                   = true
    concurrency                = 8
    credentials_secret         = "clickhouse_basic_auth"
    database                   = "logs"
    describe_table             = "DESCRIBE TABLE app_events"
    description                = "Ingest logs to ClickHouse with async inserts and TLS"
    dump_format_errors_to_disk = true
    environment                = "main"
    exclude_mapping_fields = [
      "_raw",
      "ts",
    ]
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    format                      = "json-each-row"
    id                          = "clickhouse_ingest_prod"
    login_url                   = "https://auth.example.com/oauth/token"
    mapping_type                = "automatic"
    max_payload_events          = 1000
    max_payload_size_kb         = 2048
    oauth_headers = [
      {
        name  = "Accept"
        value = "application/json"
      }
    ]
    oauth_params = [
      {
        name  = "grant_type"
        value = "client_credentials"
      }
    ]
    on_backpressure = "block"
    password        = "s3cr3tPass!"
    pipeline        = "main"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    secret            = "s3cr3tClientSecret"
    secret_param_name = "client_secret"
    sql_username      = "clickuser"
    streamtags = [
      "prod",
      "clickhouse",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    table_name  = "app_events"
    text_secret = "clickhouse_bearer_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    tls = {
      ca_path          = "/etc/ssl/certs/ca-bundle.crt"
      cert_path        = "/opt/cribl/certs/client.crt"
      certificate_name = "clickhouse-client"
      disabled         = false
      max_version      = "TLSv1.3"
      min_version      = "TLSv1.2"
      passphrase       = "s3cr3t"
      priv_key_path    = "/opt/cribl/certs/client.key"
      servername       = "clickhouse.example.com"
    }
    token                  = "chBearerToken_abc123xyz"
    token_attribute_name   = "access_token"
    token_timeout_secs     = 3600
    type                   = "click_house"
    url                    = "https://clickhouse.example.com:8443"
    use_round_robin_dns    = true
    username               = "clickuser"
    wait_for_async_inserts = true
  }
  output_cloudflare_r2 = {
    add_id_to_stage_path      = false
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    compress                  = "none"
    compression_level         = "best_speed"
    deadletter_enabled        = true
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    dest_path                 = "...my_dest_path..."
    directory_batch_size      = 0.49
    empty_dir_cleanup_sec     = 21545.04
    enable_page_checksum      = false
    enable_statistics         = true
    enable_write_page_index   = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    force_close_on_shutdown   = true
    format                    = "parquet"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_concurrent_file_parts = 3.4
    max_file_idle_time_sec    = 84949.45
    max_file_open_time_sec    = 59360.77
    max_file_size_mb          = 410.5
    max_open_files            = 1667.28
    max_retry_num             = 10.21
    object_acl                = "{ \"see\": \"documentation\" }"
    on_backpressure           = "block"
    on_disk_full_backpressure = "drop"
    parquet_data_page_version = "DATA_PAGE_V1"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 30271683.86
    parquet_schema            = "...my_parquet_schema..."
    parquet_version           = "PARQUET_2_4"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    region                    = "{ \"see\": \"documentation\" }"
    reject_unauthorized       = false
    remove_empty_dirs         = true
    reuse_connections         = true
    server_side_encryption    = "AES256"
    should_log_invalid_rows   = true
    signature_version         = "v2"
    stage_path                = "...my_stage_path..."
    storage_class             = "STANDARD"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "cloudflare_r2"
    verify_permissions    = false
    write_high_water_mark = 1607.6
  }
  output_cloudwatch = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/CloudWatchLogsWriter"
    assume_role_external_id   = "external-id-abc123"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws_cloudwatch_credentials"
    aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    description               = "Send application logs to Amazon CloudWatch Logs"
    duration_seconds          = 3600
    enable_assume_role        = true
    endpoint                  = "https://logs.us-east-1.amazonaws.com"
    environment               = "main"
    flush_period_sec          = 2
    id                        = "cloudwatch_logs_prod"
    log_group_name            = "/aws/eks/cluster-1/app-logs"
    log_stream_name           = "app-logs"
    max_queue_size            = 10
    max_record_size_kb        = 1024
    on_backpressure           = "block"
    pipeline                  = "main"
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "100 MB"
    pq_max_size         = "10GB"
    pq_mode             = "backpressure"
    pq_on_backpressure  = "block"
    pq_path             = "/opt/cribl/state/queues"
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    streamtags = [
      "prod",
      "cloudwatch",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type = "cloudwatch"
  }
  output_confluent_cloud = {
    ack                    = 1
    authentication_timeout = 10000
    backoff_rate           = 2
    brokers = [
      "mycluster.us-central1.gcp.confluent.cloud:9092",
    ]
    compression        = "gzip"
    connection_timeout = 10000
    description        = "Produce events to Confluent Cloud Kafka with Schema Registry"
    environment        = "main"
    flush_event_count  = 1000
    flush_period_sec   = 1
    format             = "json"
    id                 = "ccloud-out"
    initial_backoff    = 1000
    kafka_schema_registry = {
      auth = {
        credentials_secret = "ccloud-schema-registry-basic"
        disabled           = false
      }
      connection_timeout      = 30000
      default_key_schema_id   = 1
      default_value_schema_id = 100
      disabled                = false
      max_retries             = 3
      request_timeout         = 30000
      schema_registry_url     = "https://schema-registry.confluent.cloud"
      tls = {
        ca_path             = "/etc/ssl/certs/ca.pem"
        cert_path           = "/etc/ssl/certs/client.crt"
        certificate_name    = "ccloud-sr-cert"
        disabled            = false
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "***REDACTED***"
        priv_key_path       = "/etc/ssl/private/client.key"
        reject_unauthorized = true
        servername          = "schema-registry.confluent.cloud"
      }
    }
    max_back_off       = 60000
    max_record_size_kb = 768
    max_retries        = 5
    on_backpressure    = "block"
    pipeline           = "default"
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "100 MB"
    pq_max_size                = "10GB"
    pq_mode                    = "backpressure"
    pq_on_backpressure         = "block"
    pq_path                    = "/opt/cribl/state/queues"
    protobuf_library_id        = "user-events-protos"
    reauthentication_threshold = 60000
    request_timeout            = 60000
    sasl = {
      disabled  = false
      mechanism = "plain"
    }
    streamtags = [
      "confluent",
      "kafka",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "ccloud-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "mycluster.us-central1.gcp.confluent.cloud"
    }
    topic = "app-events"
    type  = "confluent_cloud"
  }
  output_cribl_http = {
    compression            = "gzip"
    concurrency            = 8
    description            = "Send events to Cribl Worker HTTP endpoint with retries"
    dns_resolve_period_sec = 300
    environment            = "main"
    exclude_fields = [
      "__kube_*",
      "__metadata",
    ]
    exclude_self = false
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode   = "payloadAndHeaders"
    flush_period_sec              = 2
    id                            = "cribl_http_prod"
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    max_payload_events            = 1000
    max_payload_size_kb           = 2048
    on_backpressure               = "block"
    pipeline                      = "main"
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "cribl",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/opt/cribl/certs/client.crt"
      certificate_name    = "cribl-client"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "s3cr3t"
      priv_key_path       = "/opt/cribl/certs/client.key"
      reject_unauthorized = true
      servername          = "collector.cribl.example.com"
    }
    token_ttl_minutes = 60
    type              = "cribl_http"
    url               = "https://edge.example.com:10200"
    urls = [
      {
        url    = "https://edge01.example.com:10200"
        weight = 2
      }
    ]
    use_round_robin_dns = true
  }
  output_cribl_lake = {
    description = "Cribl Lake destination"
    dest_path   = "security_logs"
    id          = "lake_ingest_prod"
    type        = "cribl_lake"
  }
  output_cribl_tcp = {
    compression            = "gzip"
    connection_timeout     = 10000
    description            = "Send events to Cribl Edge over TCP with TLS"
    dns_resolve_period_sec = 300
    environment            = "main"
    exclude_fields = [
      "__kube_*",
      "__metadata",
    ]
    exclude_self = false
    host         = "edge01.example.com"
    hosts = [
      {
        host       = "edge02.example.com"
        port       = 10300
        servername = "edge02.example.com"
        tls        = "inherit"
        weight     = 2
      }
    ]
    id                            = "cribl_tcp_prod"
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    log_failed_requests           = true
    max_concurrent_senders        = 4
    on_backpressure               = "block"
    pipeline                      = "main"
    port                          = 10300
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    streamtags = [
      "prod",
      "cribl",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    throttle_rate_per_sec = "10 MB"
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/opt/cribl/certs/client.crt"
      certificate_name    = "cribl-client"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "s3cr3t"
      priv_key_path       = "/opt/cribl/certs/client.key"
      reject_unauthorized = true
      servername          = "tcp.cribl.example.com"
    }
    token_ttl_minutes = 60
    type              = "cribl_tcp"
    write_timeout     = 30000
  }
  output_crowdstrike_next_gen_siem = {
    auth_type   = "secret"
    compress    = true
    concurrency = 8
    description = "Send events to CrowdStrike Next-Gen SIEM with token auth"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    format                      = "raw"
    id                          = "cs_nextgen_siem_prod"
    max_payload_events          = 1000
    max_payload_size_kb         = 8192
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "crowdstrike",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "crowdstrike_nextgen_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    token               = "csngs-0123456789abcdef0123456789abcdef"
    type                = "crowdstrike_next_gen_siem"
    url                 = "https://ingest.us-1.crowdstrike.com/api/ingest/hec/abcd1234/v1/services/collector"
    use_round_robin_dns = true
  }
  output_databricks = {
    add_id_to_stage_path    = false
    automatic_schema        = true
    base_file_name          = "...my_base_file_name..."
    catalog                 = "...my_catalog..."
    client_id               = "...my_client_id..."
    client_text_secret      = "...my_client_text_secret..."
    compress                = "gzip"
    compression_level       = "best_compression"
    deadletter_enabled      = true
    deadletter_path         = "...my_deadletter_path..."
    description             = "...my_description..."
    dest_path               = "...my_dest_path..."
    directory_batch_size    = 8.14
    empty_dir_cleanup_sec   = 72689.38
    enable_page_checksum    = false
    enable_statistics       = false
    enable_write_page_index = true
    environment             = "...my_environment..."
    events_volume_name      = "...my_events_volume_name..."
    file_name_suffix        = "...my_file_name_suffix..."
    force_close_on_shutdown = false
    format                  = "json"
    header_line             = "...my_header_line..."
    id                      = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_file_idle_time_sec    = 681.94
    max_file_open_time_sec    = 993.14
    max_file_size_mb          = 727.03
    max_open_files            = 389.22
    max_retry_num             = 6.59
    on_backpressure           = "block"
    on_disk_full_backpressure = "drop"
    parquet_data_page_version = "DATA_PAGE_V1"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 15677567.64
    parquet_schema            = "...my_parquet_schema..."
    parquet_version           = "PARQUET_2_6"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    remove_empty_dirs         = true
    schema                    = "...my_schema..."
    scope                     = "...my_scope..."
    should_log_invalid_rows   = true
    stage_path                = "...my_stage_path..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_sec           = 35.93
    type                  = "databricks"
    workspace_id          = "...my_workspace_id..."
    write_high_water_mark = 2087.98
  }
  output_datadog = {
    allow_api_key_from_events = false
    api_key                   = "0123456789abcdef0123456789abcdef"
    auth_type                 = "secret"
    batch_by_tags             = true
    compress                  = true
    concurrency               = 8
    content_type              = "json"
    custom_url                = "https://http-intake.logs.datadoghq.com/api/v2/logs"
    description               = "Send logs to Datadog Logs API with token auth"
    environment               = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    host                        = "web-01.example.com"
    id                          = "datadog_logs_prod"
    max_payload_events          = 1000
    max_payload_size_kb         = 2048
    message                     = "_raw"
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    send_counters_as_count = true
    service                = "web-app"
    severity               = "info"
    site                   = "us"
    source                 = "nginx"
    streamtags = [
      "prod",
      "datadog",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    tags = [
      "env:prod",
      "team:platform",
    ]
    text_secret = "datadog_api_key"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    total_memory_limit_kb = 51200
    type                  = "datadog"
    use_round_robin_dns   = true
  }
  output_dataset = {
    api_key          = "ds-0123456789abcdef0123456789abcdef"
    auth_type        = "secret"
    compress         = true
    concurrency      = 8
    custom_url       = "https://api.dataset.com/v1/logs"
    default_severity = "info"
    description      = "Send events to DataSet with API key authentication"
    environment      = "main"
    exclude_fields = [
      "sev",
      "_time",
    ]
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    id                          = "dataset_logs_prod"
    max_payload_events          = 1000
    max_payload_size_kb         = 2048
    message_field               = "_raw"
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    server_host_field = "host"
    site              = "us"
    streamtags = [
      "prod",
      "dataset",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "dataset_api_key"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    timestamp_field       = "ts"
    total_memory_limit_kb = 51200
    type                  = "dataset"
    use_round_robin_dns   = true
  }
  output_default = {
    default_id  = "http-default"
    environment = "main"
    id          = "default-output"
    pipeline    = "default"
    streamtags = [
      "default",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    type = "default"
  }
  output_devnull = {
    environment = "main"
    id          = "devnull-out"
    pipeline    = "default"
    streamtags = [
      "discard",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    type = "devnull"
  }
  output_disk_spool = {
    compress       = "gzip"
    description    = "Local disk spool for short-term buffering and replay"
    environment    = "main"
    id             = "disk_spool_buffer"
    max_data_size  = "100GB"
    max_data_time  = "7d"
    partition_expr = ""
    pipeline       = "main"
    streamtags = [
      "prod",
      "spool",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    time_window = "10m"
    type        = "disk_spool"
  }
  output_dl_s3 = {
    add_id_to_stage_path      = true
    assume_role_arn           = "arn:aws:iam::123456789012:role/S3Writer"
    assume_role_external_id   = "external-id-abc123"
    automatic_schema          = true
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws_s3_credentials"
    aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    base_file_name            = "app-logs"
    bucket                    = "logs-archive-prod"
    compress                  = "gzip"
    compression_level         = "normal"
    deadletter_enabled        = true
    deadletter_path           = "/opt/cribl/state/outputs/dead-letter"
    description               = "Archive logs to S3 in Parquet with field-based partitions"
    dest_path                 = "year=%Y/month=%m/day=%d/app=orders"
    duration_seconds          = 3600
    empty_dir_cleanup_sec     = 600
    enable_assume_role        = true
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "https://s3.us-east-1.amazonaws.com"
    environment               = "main"
    file_name_suffix          = ".parquet.gz"
    format                    = "parquet"
    header_line               = "timestamp,host,level,message"
    id                        = "dls3_archive_prod"
    key_value_metadata = [
      {
        key   = "OCSF Event Class"
        value = "9001"
      }
    ]
    kms_key_id                        = "arn:aws:kms:us-east-1:123456789012:key/abcd-1234-efgh-5678"
    max_closing_files_to_backpressure = 500
    max_concurrent_file_parts         = 5
    max_file_idle_time_sec            = 120
    max_file_open_time_sec            = 600
    max_file_size_mb                  = 256
    max_open_files                    = 200
    max_retry_num                     = 20
    object_acl                        = "private"
    on_disk_full_backpressure         = "block"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "128MB"
    parquet_row_group_length          = 100000
    parquet_version                   = "PARQUET_2_6"
    partitioning_fields = [
      "app",
      "env",
    ]
    pipeline                = "main"
    region                  = "us-east-1"
    reject_unauthorized     = true
    remove_empty_dirs       = true
    reuse_connections       = true
    server_side_encryption  = "aws:kms"
    should_log_invalid_rows = true
    signature_version       = "v4"
    stage_path              = "/opt/cribl/state/outputs/staging"
    storage_class           = "INTELLIGENT_TIERING"
    streamtags = [
      "prod",
      "archive",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type                  = "dl_s3"
    verify_permissions    = true
    write_high_water_mark = 256
  }
  output_dynatrace_http = {
    active_gate_domain = "https://activegate.example.com:9999/e/abc12345/api/v2/logs/ingest"
    auth_type          = "token"
    compress           = true
    concurrency        = 8
    description        = "Send logs to Dynatrace Logs Ingest API"
    endpoint           = "cloud"
    environment        = "main"
    environment_id     = "abc12345"
    extra_http_headers = [
      {
        name  = "Api-Token"
        value = "dt0c01.XXXX.YYYYZZZZ"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    format                      = "json_array"
    id                          = "dynatrace_http_metrics"
    keep_alive                  = true
    max_payload_events          = 10000
    max_payload_size_kb         = 4096
    method                      = "POST"
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "api-token",
    ]
    streamtags = [
      "prod",
      "dynatrace",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    telemetry_type = "logs"
    text_secret    = "dynatrace_api_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    token                 = "dt0c01.XXXX.YYYYZZZZ"
    total_memory_limit_kb = 51200
    type                  = "dynatrace_http"
    url                   = "https://abc.live.dynatrace.com/e/abc12345/api/v2/logs/ingest"
    use_round_robin_dns   = true
  }
  output_dynatrace_otlp = {
    auth_token_name    = "Authorization"
    compress           = "gzip"
    concurrency        = 5
    connection_timeout = 10000
    description        = "Send OTLP logs and metrics to Dynatrace SaaS"
    endpoint           = "https://abc123.live.dynatrace.com/api/v2/otlp"
    endpoint_type      = "saas"
    environment        = "main"
    extra_http_headers = [
      {
        name  = "Api-Token"
        value = "dt0c01.XXXX.YYYYZZZZ"
      }
    ]
    failed_request_logging_mode    = "payloadAndHeaders"
    flush_period_sec               = 2
    http_compress                  = "gzip"
    http_logs_endpoint_override    = "https://abc123.live.dynatrace.com/api/v2/otlp/v1/logs"
    http_metrics_endpoint_override = "https://abc123.live.dynatrace.com/api/v2/otlp/v1/metrics"
    http_traces_endpoint_override  = "https://abc123.live.dynatrace.com/api/v2/otlp/v1/traces"
    id                             = "dynatrace_otlp_export"
    keep_alive                     = true
    keep_alive_time                = 30
    max_payload_size_kb            = 2048
    metadata = [
      {
        key   = "x-tenant-id"
        value = "acme-prod"
      }
    ]
    on_backpressure = "block"
    otlp_version    = "1.3.1"
    pipeline        = "main"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    protocol                          = "http"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "api-token",
    ]
    streamtags = [
      "prod",
      "dynatrace",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    token_secret        = "dynatrace_otlp_token"
    type                = "dynatrace_otlp"
    use_round_robin_dns = true
  }
  output_elastic = {
    auth = {
      auth_type = "manualAPIKey"
      disabled  = false
    }
    compress               = true
    concurrency            = 8
    description            = "Send documents to Elasticsearch bulk API with retries and custom params"
    dns_resolve_period_sec = 300
    doc_type               = "_doc"
    elastic_pipeline       = "ingest-grok-pipeline"
    elastic_version        = "7"
    environment            = "main"
    exclude_self           = false
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    extra_params = [
      {
        name  = "filter_path"
        value = "errors,items.*.error,items.*._index,items.*.status"
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 1
    id                            = "elastic-out"
    include_doc_id                = true
    index                         = "\"logs-2024.01.15\""
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    max_payload_events            = 0
    max_payload_size_kb           = 4096
    on_backpressure               = "block"
    pipeline                      = "default"
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    retry_partial_errors = true
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "elastic",
      "es",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    type        = "elastic"
    url         = "https://es.example.com:9200/_bulk"
    urls = [
      {
        url    = "https://es-node-1.example.com:9200/_bulk"
        weight = 2
      }
    ]
    use_round_robin_dns = true
    write_action        = "create"
  }
  output_elastic_cloud = {
    auth = {
      auth_type = "manualAPIKey"
      disabled  = false
    }
    compress         = true
    concurrency      = 8
    description      = "Send documents to Elastic Cloud with retries and pipeline support"
    elastic_pipeline = "ingest-grok-pipeline"
    environment      = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    extra_params = [
      {
        name  = "filter_path"
        value = "errors,items.*.error,items.*._index,items.*.status"
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    id                          = "es-cloud-out"
    include_doc_id              = true
    index                       = "\"logs-2024.01.15\""
    max_payload_events          = 0
    max_payload_size_kb         = 4096
    on_backpressure             = "block"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "elastic",
      "cloud",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    type        = "elastic_cloud"
    url         = "my-deployment:ZXM0LmNsb3VkLmV... (truncated)"
  }
  output_exabeam = {
    add_id_to_stage_path      = true
    aws_api_key               = "***REDACTED***"
    aws_secret_key            = "***REDACTED***"
    bucket                    = "exabeam-data"
    collector_instance_id     = "11112222-3333-4444-5555-666677778888"
    deadletter_enabled        = true
    deadletter_path           = "/var/lib/cribl/state/outputs/dead-letter"
    description               = "Deliver logs to Exabeam Collector via GCS staging"
    empty_dir_cleanup_sec     = 600
    encoded_configuration     = "***REDACTED***"
    endpoint                  = "https://storage.googleapis.com"
    environment               = "main"
    id                        = "exabeam-out"
    max_file_idle_time_sec    = 30
    max_file_open_time_sec    = 300
    max_file_size_mb          = 64
    max_open_files            = 200
    max_retry_num             = 20
    object_acl                = "private"
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    pipeline                  = "default"
    region                    = "us-central1"
    reject_unauthorized       = true
    remove_empty_dirs         = true
    reuse_connections         = true
    signature_version         = "v4"
    site_id                   = "site-123"
    site_name                 = "\"corp-east\""
    stage_path                = "/var/lib/cribl/state/outputs/staging"
    storage_class             = "NEARLINE"
    streamtags = [
      "exabeam",
      "gcs",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    timezone_offset = "-07:00"
    type            = "exabeam"
  }
  output_filesystem = {
    add_id_to_stage_path    = true
    automatic_schema        = true
    base_file_name          = "`CriblOut`"
    compress                = "gzip"
    compression_level       = "normal"
    deadletter_enabled      = true
    deadletter_path         = "/var/lib/cribl/state/outputs/dead-letter"
    description             = "Write events to local filesystem with daily partitioning"
    dest_path               = "/var/log/cribl/out"
    empty_dir_cleanup_sec   = 600
    enable_page_checksum    = true
    enable_statistics       = true
    enable_write_page_index = true
    environment             = "main"
    file_name_suffix        = ".json.gz"
    format                  = "json"
    header_line             = "timestamp,host,message"
    id                      = "filesystem-out"
    key_value_metadata = [
      {
        key   = "team"
        value = "platform"
      }
    ]
    max_file_idle_time_sec    = 30
    max_file_open_time_sec    = 300
    max_file_size_mb          = 64
    max_open_files            = 200
    max_retry_num             = 20
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "4MB"
    parquet_row_group_length  = 10000
    parquet_version           = "PARQUET_2_6"
    partition_expr            = "C.Time.strftime(_time ? _time : Date.now()/1000, '%Y/%m/%d') + '/host=' + host"
    pipeline                  = "default"
    remove_empty_dirs         = true
    should_log_invalid_rows   = true
    stage_path                = "/var/log/cribl/stage"
    streamtags = [
      "filesystem",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    type                  = "filesystem"
    write_high_water_mark = 256
  }
  output_google_chronicle = {
    api_key               = "***REDACTED***"
    api_key_secret        = "chronicle-api-key"
    api_version           = "v1"
    authentication_method = "serviceAccount"
    compress              = true
    concurrency           = 8
    custom_labels = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    customer_id = "123e4567-e89b-12d3-a456-426614174000"
    description = "Send events to Google SecOps Chronicle"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    extra_log_types = [
      {
        description = "...my_description..."
        log_type    = "...my_log_type..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    id                          = "chronicle-out"
    log_format_type             = "unstructured"
    log_text_field              = "message"
    log_type                    = "CUSTOM_WEBLOG"
    max_payload_events          = 0
    max_payload_size_kb         = 1024
    namespace                   = "prod-us"
    on_backpressure             = "block"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    region                            = "us"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 30000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    service_account_credentials        = "***REDACTED***"
    service_account_credentials_secret = "chronicle-sa-credentials"
    streamtags = [
      "google",
      "chronicle",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 90
    total_memory_limit_kb = 5120
    type                  = "google_chronicle"
    use_round_robin_dns   = true
  }
  output_google_cloud_logging = {
    cache_fill_bytes_expression = "String(_raw.cache_fill_bytes)"
    cache_hit_expression        = "Boolean(_raw.cache_hit)"
    cache_lookup_expression     = "Boolean(_raw.cache_lookup)"
    cache_validated_expression  = "Boolean(_raw.cache_validated)"
    concurrency                 = 8
    connection_timeout          = 10000
    description                 = "Send logs to Google Cloud Logging with custom resource labels"
    environment                 = "main"
    file_expression             = "String(_raw.file)"
    first_expression            = "Boolean(_raw.operation_first)"
    flush_period_sec            = 1
    function_expression         = "String(_raw.function)"
    google_auth_method          = "secret"
    id                          = "gcl-out"
    id_expression               = "String(_raw.operation_id)"
    index_expression            = "Number(_raw.split_index)"
    insert_id_expression        = "Crypto.uuid()"
    last_expression             = "Boolean(_raw.operation_last)"
    latency_expression          = "(_raw.latency_ms/1000).toFixed(3) + \"s\""
    line_expression             = "String(_raw.line)"
    log_labels = [
      {
        label            = "environment"
        value_expression = "\"prod\""
      }
    ]
    log_location_expression = "\"projects/my-project\""
    log_location_type       = "project"
    log_name_expression     = "\"cribl_logs\""
    max_payload_events      = 0
    max_payload_size_kb     = 4096
    on_backpressure         = "block"
    payload_expression      = "{ message: _raw.message, severity: _raw.severity || \"DEFAULT\" }"
    payload_format          = "json"
    pipeline                = "default"
    pq_compress             = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size          = "100 MB"
    pq_max_size               = "10GB"
    pq_mode                   = "backpressure"
    pq_on_backpressure        = "block"
    pq_path                   = "/opt/cribl/state/queues"
    producer_expression       = "String(_raw.operation_producer)"
    protocol_expression       = "String(_raw.protocol)"
    referer_expression        = "String(_raw.referer)"
    remote_ip_expression      = "String(_raw.client_ip)"
    request_method_expression = "\"POST\""
    request_size_expression   = "String(length(_raw.request_body))"
    request_url_expression    = "\"https://example.com/api\""
    resource_type_expression  = "\"gce_instance\""
    resource_type_labels = [
      {
        label            = "instance_id"
        value_expression = "String(_raw.instance_id)"
      }
    ]
    response_size_expression    = "String(length(_raw.response_body))"
    secret                      = "gcl-service-account"
    server_ip_expression        = "String(_raw.server_ip)"
    service_account_credentials = "***REDACTED***"
    severity_expression         = "\"INFO\""
    span_id_expression          = "String(_raw.span_id)"
    status_expression           = "Number(_raw.status)"
    streamtags = [
      "gcp",
      "logging",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    throttle_rate_req_per_sec = 500
    timeout_sec               = 30
    total_memory_limit_kb     = 20480
    total_splits_expression   = "Number(_raw.split_total)"
    trace_expression          = "String(_raw.trace)"
    trace_sampled_expression  = "Boolean(_raw.trace_sampled)"
    type                      = "google_cloud_logging"
    uid_expression            = "String(_raw.split_uid)"
    user_agent_expression     = "String(_raw.user_agent)"
  }
  output_google_cloud_storage = {
    add_id_to_stage_path      = true
    automatic_schema          = true
    aws_api_key               = "***REDACTED***"
    aws_authentication_method = "manual"
    aws_secret                = "gcs-hmac-credentials"
    aws_secret_key            = "***REDACTED***"
    base_file_name            = "`CriblOut`"
    bucket                    = "cribl-data-bucket"
    compress                  = "gzip"
    compression_level         = "normal"
    deadletter_enabled        = true
    deadletter_path           = "/var/lib/cribl/state/outputs/dead-letter"
    description               = "Write objects to Google Cloud Storage with date-based partitioning"
    dest_path                 = "logs/ingest"
    empty_dir_cleanup_sec     = 600
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "https://storage.googleapis.com"
    environment               = "main"
    file_name_suffix          = ".json.gz"
    format                    = "json"
    header_line               = "timestamp,host,message"
    id                        = "gcs-out"
    key_value_metadata = [
      {
        key   = "team"
        value = "platform"
      }
    ]
    max_file_idle_time_sec    = 30
    max_file_open_time_sec    = 300
    max_file_size_mb          = 64
    max_open_files            = 200
    max_retry_num             = 20
    object_acl                = "private"
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "4MB"
    parquet_row_group_length  = 10000
    parquet_version           = "PARQUET_2_6"
    partition_expr            = "2024/01/15"
    pipeline                  = "default"
    region                    = "us-central1"
    reject_unauthorized       = true
    remove_empty_dirs         = true
    reuse_connections         = true
    should_log_invalid_rows   = true
    signature_version         = "v4"
    stage_path                = "/var/lib/cribl/state/outputs/staging"
    storage_class             = "NEARLINE"
    streamtags = [
      "gcp",
      "gcs",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    type                  = "google_cloud_storage"
    verify_permissions    = true
    write_high_water_mark = 256
  }
  output_google_pubsub = {
    batch_size         = 1000
    batch_timeout      = 200
    create_topic       = true
    description        = "Publish events to Google Pub/Sub with ordered delivery"
    environment        = "main"
    flush_period_sec   = 1
    google_auth_method = "secret"
    id                 = "gpubsub-out"
    max_in_progress    = 20
    max_queue_size     = 500
    max_record_size_kb = 256
    on_backpressure    = "block"
    ordered_delivery   = true
    pipeline           = "default"
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size            = "100 MB"
    pq_max_size                 = "10GB"
    pq_mode                     = "backpressure"
    pq_on_backpressure          = "block"
    pq_path                     = "/opt/cribl/state/queues"
    region                      = "us-central1"
    secret                      = "gcp-pubsub-sa"
    service_account_credentials = "***REDACTED***"
    streamtags = [
      "gcp",
      "pubsub",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    topic_name = "app-events"
    type       = "google_pubsub"
  }
  output_grafana_cloud = {
    compress    = true
    concurrency = 2
    description = "Send logs and metrics to Grafana Cloud Loki and Prometheus"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 10
    id                          = "grafana_cloud_logs_prod"
    labels = [
      {
        name  = "host"
        value = "web-01"
      }
    ]
    loki_auth = {
      auth_type          = "basic"
      credentials_secret = "grafana_loki_credentials"
      password           = "glc_efgh5678"
      text_secret        = "grafana_loki_token"
      token              = "12345:glc_efgh5678"
      username           = 12345
    }
    loki_url            = "https://logs-prod-us-central1.grafana.net"
    max_payload_events  = 1000
    max_payload_size_kb = 2048
    message             = "_raw"
    message_format      = "protobuf"
    metric_rename_expr  = "name.replace(/[^a-zA-Z0-9_]/g, '_')"
    on_backpressure     = "block"
    pipeline            = "main"
    pq_compress         = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    prometheus_auth = {
      auth_type          = "basic"
      credentials_secret = "grafana_prom_credentials"
      password           = "glc_abcd1234"
      text_secret        = "grafana_prom_token"
      token              = "12345:glc_abcd1234"
      username           = 12345
    }
    prometheus_url                    = "https://prometheus-blocks-prod-us-central1.grafana.net/api/prom/push"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "grafana",
    ]
    system_fields = [
      "cribl_host",
      "cribl_wp",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    type                = "grafana_cloud"
    use_round_robin_dns = true
  }
  output_graphite = {
    connection_timeout     = 10000
    description            = "Send metrics to Graphite in plaintext protocol"
    dns_resolve_period_sec = 300
    environment            = "main"
    flush_period_sec       = 1
    host                   = "graphite.example.com"
    id                     = "graphite_metrics_prod"
    mtu                    = 1400
    on_backpressure        = "block"
    pipeline               = "metrics"
    port                   = 2003
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    protocol           = "tcp"
    streamtags = [
      "prod",
      "graphite",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    throttle_rate_per_sec = "10 MB"
    type                  = "graphite"
    write_timeout         = 30000
  }
  output_honeycomb = {
    auth_type   = "manual"
    compress    = true
    concurrency = 8
    dataset     = "observability"
    description = "Send events to Honeycomb dataset"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    id                          = "honeycomb-out"
    max_payload_events          = 0
    max_payload_size_kb         = 4096
    on_backpressure             = "block"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "honeycomb",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    team        = "***REDACTED***"
    text_secret = "honeycomb-api-key"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    type                = "honeycomb"
    use_round_robin_dns = true
  }
  output_humio_hec = {
    auth_type   = "secret"
    compress    = true
    concurrency = 8
    description = "Send logs to CrowdStrike Falcon LogScale via HEC endpoint"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    format                      = "JSON"
    id                          = "humio_hec_prod"
    max_payload_events          = 1000
    max_payload_size_kb         = 8192
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "humio",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "humio_hec_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    token               = "humio-0123456789abcdef0123456789abcdef"
    type                = "humio_hec"
    url                 = "https://cloud.us.humio.com/api/v1/ingest/hec"
    use_round_robin_dns = true
  }
  output_influxdb = {
    auth_header_expr         = "`Bearer ${token}`"
    auth_type                = "token"
    bucket                   = "metrics_prod"
    compress                 = true
    concurrency              = 8
    credentials_secret       = "influxdb_basic_auth"
    database                 = "telegraf"
    description              = "Send metrics to InfluxDB with v2 API and token auth"
    dynamic_value_field_name = true
    environment              = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    id                          = "influxdb_metrics_prod"
    login_url                   = "https://influxdb.example.com/oauth/token"
    max_payload_events          = 5000
    max_payload_size_kb         = 8192
    oauth_headers = [
      {
        name  = "Accept"
        value = "application/json"
      }
    ]
    oauth_params = [
      {
        name  = "grant_type"
        value = "client_credentials"
      }
    ]
    on_backpressure = "block"
    org             = "acme-observability"
    password        = "s3cr3tPass!"
    pipeline        = "metrics"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    secret            = "s3cr3tClientSecret"
    secret_param_name = "client_secret"
    streamtags = [
      "prod",
      "influxdb",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "influxdb_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec          = 30
    timestamp_precision  = "ms"
    token                = "influxV2Token_abc123xyz"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "influxdb"
    url                  = "https://influxdb.example.com:8086/write"
    use_round_robin_dns  = true
    use_v2_api           = true
    username             = "influx_writer"
    value_field_name     = "value"
  }
  output_kafka = {
    ack                    = 1
    authentication_timeout = 10000
    backoff_rate           = 2
    brokers = [
      "kafka-1.example.com:9092",
      "kafka-2.example.com:9092",
    ]
    compression        = "gzip"
    connection_timeout = 10000
    description        = "Produce events to Kafka with retries and TLS"
    environment        = "main"
    flush_event_count  = 1000
    flush_period_sec   = 1
    format             = "json"
    id                 = "kafka-out"
    initial_backoff    = 1000
    kafka_schema_registry = {
      auth = {
        credentials_secret = "schema-registry-basic-auth"
        disabled           = false
      }
      connection_timeout      = 30000
      default_key_schema_id   = 1
      default_value_schema_id = 2
      disabled                = true
      max_retries             = 3
      request_timeout         = 30000
      schema_registry_url     = "https://schema-registry.example.com:8081"
      tls = {
        ca_path             = "/etc/ssl/certs/ca.pem"
        cert_path           = "/etc/ssl/certs/client.crt"
        certificate_name    = "kafka-client-cert"
        disabled            = false
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "***REDACTED***"
        priv_key_path       = "/etc/ssl/private/client.key"
        reject_unauthorized = true
        servername          = "schema-registry.example.com"
      }
    }
    max_back_off       = 60000
    max_record_size_kb = 768
    max_retries        = 5
    on_backpressure    = "block"
    pipeline           = "default"
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "100 MB"
    pq_max_size                = "10GB"
    pq_mode                    = "backpressure"
    pq_on_backpressure         = "block"
    pq_path                    = "/opt/cribl/state/queues"
    protobuf_library_id        = "user-events-protos"
    reauthentication_threshold = 60000
    request_timeout            = 60000
    sasl = {
      disabled  = false
      mechanism = "scram-sha-512"
    }
    streamtags = [
      "kafka",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "kafka-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "kafka-broker.example.com"
    }
    topic = "app-events"
    type  = "kafka"
  }
  output_kinesis = {
    as_ndjson                 = true
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-kinesis-writer"
    assume_role_external_id   = "cribl-external-123"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    compression               = "gzip"
    concurrency               = 8
    description               = "Deliver events to AWS Kinesis Data Streams"
    duration_seconds          = 3600
    enable_assume_role        = true
    endpoint                  = "https://kinesis.us-east-1.amazonaws.com"
    environment               = "main"
    flush_period_sec          = 1
    id                        = "kinesis-out"
    max_record_size_kb        = 1024
    on_backpressure           = "block"
    pipeline                  = "default"
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "100 MB"
    pq_max_size         = "10GB"
    pq_mode             = "backpressure"
    pq_on_backpressure  = "block"
    pq_path             = "/opt/cribl/state/queues"
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    signature_version   = "v4"
    stream_name         = "app-events-stream"
    streamtags = [
      "aws",
      "kinesis",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    type            = "kinesis"
    use_list_shards = true
  }
  output_loki = {
    auth_type          = "token"
    compress           = true
    concurrency        = 2
    credentials_secret = "grafana_loki_credentials"
    description        = "Send logs to Loki with labels and batching"
    environment        = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 10
    id                          = "loki_logs_prod"
    labels = [
      {
        name  = "host"
        value = "web-01"
      }
    ]
    max_payload_events  = 1000
    max_payload_size_kb = 2048
    message             = "_raw"
    message_format      = "protobuf"
    on_backpressure     = "block"
    password            = "glc_abcd1234"
    pipeline            = "main"
    pq_compress         = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "loki",
    ]
    system_fields = [
      "cribl_host",
      "cribl_wp",
    ]
    text_secret = "grafana_loki_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    token                 = "12345:glc_abcd1234"
    total_memory_limit_kb = 51200
    type                  = "loki"
    url                   = "https://loki.example.com/loki/api/v1/push"
    use_round_robin_dns   = true
    username              = 12345
  }
  output_microsoft_fabric = {
    ack                    = 9
    authentication_timeout = 3288750.42
    backoff_rate           = 5.11
    bootstrap_server       = "...my_bootstrap_server..."
    connection_timeout     = 1684445.58
    description            = "...my_description..."
    environment            = "...my_environment..."
    flush_event_count      = 1487.63
    flush_period_sec       = 1.31
    format                 = "json"
    id                     = "...my_id..."
    initial_backoff        = 29883.16
    max_back_off           = 77303.8
    max_record_size_kb     = 6.17
    max_retries            = 97.46
    on_backpressure        = "drop"
    pipeline               = "...my_pipeline..."
    pq_compress            = "none"
    pq_controls = {
      # ...
    }
    pq_max_backpressure_sec    = 0.99
    pq_max_buffer_size         = 672.08
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "error"
    pq_on_backpressure         = "block"
    pq_path                    = "...my_pq_path..."
    pq_rate_per_sec            = 8.31
    pq_strict_ordering         = true
    reauthentication_threshold = 288466.03
    request_timeout            = 1480874.35
    sasl = {
      cert_path               = "...my_cert_path..."
      certificate_name        = "...my_certificate_name..."
      client_id               = "...my_client_id..."
      client_secret_auth_type = "secret"
      client_text_secret      = "...my_client_text_secret..."
      disabled                = true
      mechanism               = "plain"
      oauth_endpoint          = "https://login.partner.microsoftonline.cn"
      passphrase              = "...my_passphrase..."
      priv_key_path           = "...my_priv_key_path..."
      scope                   = "...my_scope..."
      tenant_id               = "...my_tenant_id..."
      text_secret             = "...my_text_secret..."
      username                = "...my_username..."
    }
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tls = {
      disabled            = false
      reject_unauthorized = false
    }
    topic = "...my_topic..."
    type  = "microsoft_fabric"
  }
  output_minio = {
    add_id_to_stage_path      = true
    automatic_schema          = true
    aws_api_key               = "minio_access_key"
    aws_authentication_method = "manual"
    aws_secret                = "minio_credentials"
    aws_secret_key            = "minio_secret_key_123"
    base_file_name            = "app-logs"
    bucket                    = "logs-prod"
    compress                  = "gzip"
    compression_level         = "normal"
    deadletter_enabled        = true
    deadletter_path           = "/opt/cribl/state/outputs/dead-letter"
    description               = "Archive logs to MinIO in Parquet with date-based partitioning"
    dest_path                 = "year=%Y/month=%m/day=%d/app=orders"
    empty_dir_cleanup_sec     = 600
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "http://minio.example.com:9000"
    environment               = "main"
    file_name_suffix          = ".json.gz"
    format                    = "json"
    header_line               = "timestamp,host,level,message"
    id                        = "minio_archive_prod"
    key_value_metadata = [
      {
        key   = "OCSF Event Class"
        value = "9001"
      }
    ]
    max_concurrent_file_parts = 5
    max_file_idle_time_sec    = 120
    max_file_open_time_sec    = 600
    max_file_size_mb          = 128
    max_open_files            = 200
    max_retry_num             = 20
    object_acl                = "private"
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "128MB"
    parquet_row_group_length  = 100000
    parquet_version           = "PARQUET_2_6"
    partition_expr            = "2024/01/15"
    pipeline                  = "main"
    region                    = "us-east-1"
    reject_unauthorized       = true
    remove_empty_dirs         = true
    reuse_connections         = true
    server_side_encryption    = "AES256"
    should_log_invalid_rows   = true
    signature_version         = "v4"
    stage_path                = "/opt/cribl/state/outputs/staging"
    storage_class             = "STANDARD"
    streamtags = [
      "prod",
      "minio",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type                  = "minio"
    verify_permissions    = true
    write_high_water_mark = 256
  }
  output_msk = {
    ack                       = 1
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-msk-access"
    assume_role_external_id   = "cribl-external-123"
    authentication_timeout    = 10000
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    backoff_rate              = 2
    brokers = [
      "b-1.mskcluster.abcde.c2.kafka.us-east-1.amazonaws.com:9092",
      "b-2.mskcluster.abcde.c2.kafka.us-east-1.amazonaws.com:9092",
    ]
    compression        = "gzip"
    connection_timeout = 10000
    description        = "Produce events to Amazon MSK with retries and TLS"
    duration_seconds   = 3600
    enable_assume_role = true
    endpoint           = "https://kafka.us-east-1.amazonaws.com"
    environment        = "main"
    flush_event_count  = 1000
    flush_period_sec   = 1
    format             = "json"
    id                 = "msk-out"
    initial_backoff    = 1000
    kafka_schema_registry = {
      auth = {
        credentials_secret = "msk-schema-registry-basic"
        disabled           = false
      }
      connection_timeout      = 30000
      default_key_schema_id   = 1
      default_value_schema_id = 100
      disabled                = true
      max_retries             = 3
      request_timeout         = 30000
      schema_registry_url     = "https://schema-registry.example.com:8081"
      tls = {
        ca_path             = "/etc/ssl/certs/ca.pem"
        cert_path           = "/etc/ssl/certs/client.crt"
        certificate_name    = "msk-client-cert"
        disabled            = false
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "***REDACTED***"
        priv_key_path       = "/etc/ssl/private/client.key"
        reject_unauthorized = true
        servername          = "schema-registry.example.com"
      }
    }
    max_back_off       = 60000
    max_record_size_kb = 768
    max_retries        = 5
    on_backpressure    = "block"
    pipeline           = "default"
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "100 MB"
    pq_max_size                = "10GB"
    pq_mode                    = "backpressure"
    pq_on_backpressure         = "block"
    pq_path                    = "/opt/cribl/state/queues"
    protobuf_library_id        = "user-events-protos"
    reauthentication_threshold = 60000
    region                     = "us-east-1"
    reject_unauthorized        = true
    request_timeout            = 60000
    reuse_connections          = true
    signature_version          = "v4"
    streamtags = [
      "aws",
      "msk",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "msk-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "b-1.mskcluster.abcde.c2.kafka.us-east-1.amazonaws.com"
    }
    topic = "app-events"
    type  = "msk"
  }
  output_netflow = {
    description            = "Forward NetFlow v5/v9/IPFIX to downstream collectors"
    dns_resolve_period_sec = 300
    environment            = "main"
    hosts = [
      {
        host = "netflow-collector.example.com"
        port = 2055
      }
    ]
    id       = "netflow_export_prod"
    pipeline = "main"
    streamtags = [
      "prod",
      "netflow",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type = "netflow"
  }
  output_newrelic = {
    api_key     = "NRAK-0123456789abcdef0123456789abcdef"
    auth_type   = "manual"
    compress    = true
    concurrency = 8
    custom_url  = "https://log-api.newrelic.com/log/v1"
    description = "Send logs to New Relic Logs with custom endpoint"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    id                          = "newrelic_logs_prod"
    log_type                    = "access_log"
    max_payload_events          = 500
    max_payload_size_kb         = 512
    message_field               = "_raw"
    metadata = [
      {
        name  = "service"
        value = "`\"orders-service\"`"
      }
    ]
    on_backpressure = "block"
    pipeline        = "main"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    region                            = "US"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "newrelic",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "newrelic_api_key"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    total_memory_limit_kb = 51200
    type                  = "newrelic"
    use_round_robin_dns   = true
  }
  output_newrelic_events = {
    account_id  = "12345678"
    api_key     = "NRAK-0123456789abcdef0123456789abcdef"
    auth_type   = "secret"
    compress    = true
    concurrency = 8
    custom_url  = "https://insights-collector.newrelic.com/v1/accounts/12345678/events"
    description = "Send custom events to New Relic Events API"
    environment = "main"
    event_type  = "CriblCustomEvent"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    id                          = "newrelic_events_prod"
    max_payload_events          = 500
    max_payload_size_kb         = 512
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    region                            = "US"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "newrelic",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "newrelic_api_key"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    type                = "newrelic_events"
    use_round_robin_dns = true
  }
  output_open_telemetry = {
    auth_header_expr   = "`Bearer ${token}`"
    auth_type          = "token"
    compress           = "gzip"
    concurrency        = 5
    connection_timeout = 10000
    credentials_secret = "otel_basic_auth"
    description        = "Export telemetry to OTel Collector with OAuth and keepalive"
    endpoint           = "https://otel-collector.example.com:4317"
    environment        = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode    = "payloadAndHeaders"
    flush_period_sec               = 2
    http_compress                  = "gzip"
    http_logs_endpoint_override    = "https://otel-collector.example.com/v1/logs"
    http_metrics_endpoint_override = "https://otel-collector.example.com/v1/metrics"
    http_traces_endpoint_override  = "https://otel-collector.example.com/v1/traces"
    id                             = "otel_export_prod"
    keep_alive                     = true
    keep_alive_time                = 30
    login_url                      = "https://auth.example.com/oauth/token"
    max_payload_size_kb            = 2048
    metadata = [
      {
        key   = "x-tenant-id"
        value = "acme-prod"
      }
    ]
    oauth_headers = [
      {
        name  = "Accept"
        value = "application/json"
      }
    ]
    oauth_params = [
      {
        name  = "grant_type"
        value = "client_credentials"
      }
    ]
    on_backpressure = "block"
    otlp_version    = "1.3.1"
    password        = "s3cr3tPass!"
    pipeline        = "main"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    protocol                          = "grpc"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    secret            = "s3cr3tClientSecret"
    secret_param_name = "client_secret"
    streamtags = [
      "prod",
      "otel",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret = "otel_bearer_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/opt/cribl/certs/client.crt"
      certificate_name    = "otel-client"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "s3cr3t"
      priv_key_path       = "/opt/cribl/certs/client.key"
      reject_unauthorized = true
    }
    token                = "otelBearerToken_abc123xyz"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "open_telemetry"
    use_round_robin_dns  = true
    username             = "otel_user"
  }
  output_prometheus = {
    auth_header_expr   = "`Bearer ${token}`"
    auth_type          = "basic"
    concurrency        = 8
    credentials_secret = "prometheus_basic_auth"
    description        = "Send metrics to Prometheus remote_write with basic auth"
    environment        = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    id                          = "prometheus_metrics_prod"
    login_url                   = "https://auth.example.com/oauth/token"
    max_payload_events          = 1000
    max_payload_size_kb         = 2048
    metric_rename_expr          = "name.replace(/[^a-zA-Z0-9_]/g, '_')"
    metrics_flush_period_sec    = 60
    oauth_headers = [
      {
        name  = "Accept"
        value = "application/json"
      }
    ]
    oauth_params = [
      {
        name  = "grant_type"
        value = "client_credentials"
      }
    ]
    on_backpressure = "block"
    password        = "mimir_api_key_abcd1234"
    pipeline        = "metrics"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    secret            = "s3cr3tClientSecret"
    secret_param_name = "client_secret"
    send_metadata     = true
    streamtags = [
      "prod",
      "prometheus",
    ]
    system_fields = [
      "cribl_host",
      "cribl_wp",
    ]
    text_secret = "prometheus_bearer_token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec          = 30
    token                = "promBearerToken_abc123xyz"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "prometheus"
    url                  = "https://prometheus.example.com/api/v1/write"
    use_round_robin_dns  = true
    username             = "prometheus"
  }
  output_ring = {
    compress        = "gzip"
    description     = "Local ring buffer for short-term retention and replay"
    dest_path       = "/opt/cribl/state/ring_buffer_prod"
    environment     = "main"
    format          = "json"
    id              = "ring_buffer_prod"
    max_data_size   = "100GB"
    max_data_time   = "7d"
    on_backpressure = "block"
    partition_expr  = ""
    pipeline        = "main"
    streamtags = [
      "prod",
      "ring",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type = "ring"
  }
  output_router = {
    description = "Route events to outputs based on filter rules"
    environment = "main"
    id          = "router_main"
    pipeline    = "main"
    rules = [
      {
        description = "Route application errors to Splunk"
        filter      = "`_source == \"app\" && level == \"error\"`"
        final       = true
        output      = "OutputSplunk"
      }
    ]
    streamtags = [
      "prod",
      "routing",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type = "router"
  }
  output_s3 = {
    add_id_to_stage_path      = true
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-s3-writer"
    assume_role_external_id   = "cribl-external-123"
    automatic_schema          = true
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    base_file_name            = "`CriblOut`"
    bucket                    = "cribl-data-bucket"
    compress                  = "gzip"
    compression_level         = "normal"
    deadletter_enabled        = true
    deadletter_path           = "/var/lib/cribl/state/outputs/dead-letter"
    description               = "Write objects to S3 with date-based partitioning"
    dest_path                 = "logs/ingest"
    duration_seconds          = 3600
    empty_dir_cleanup_sec     = 600
    enable_assume_role        = true
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "https://s3.us-east-1.amazonaws.com"
    environment               = "main"
    file_name_suffix          = ".json.gz"
    format                    = "json"
    header_line               = "timestamp,host,message"
    id                        = "s3-out"
    key_value_metadata = [
      {
        key   = "team"
        value = "platform"
      }
    ]
    kms_key_id                        = "arn:aws:kms:us-east-1:123456789012:key/abcd1234-5678-90ab-cdef-EXAMPLEKEY"
    max_closing_files_to_backpressure = 500
    max_concurrent_file_parts         = 4
    max_file_idle_time_sec            = 30
    max_file_open_time_sec            = 300
    max_file_size_mb                  = 64
    max_open_files                    = 200
    max_retry_num                     = 20
    object_acl                        = "private"
    on_backpressure                   = "block"
    on_disk_full_backpressure         = "block"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "4MB"
    parquet_row_group_length          = 10000
    parquet_version                   = "PARQUET_2_6"
    partition_expr                    = "2024/01/15"
    pipeline                          = "default"
    region                            = "us-east-1"
    reject_unauthorized               = true
    remove_empty_dirs                 = true
    reuse_connections                 = true
    server_side_encryption            = "AES256"
    should_log_invalid_rows           = true
    signature_version                 = "v4"
    stage_path                        = "/var/lib/cribl/state/outputs/staging"
    storage_class                     = "STANDARD_IA"
    streamtags = [
      "s3",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    type                  = "s3"
    verify_permissions    = true
    write_high_water_mark = 256
  }
  output_security_lake = {
    account_id                = "123456789012"
    add_id_to_stage_path      = true
    assume_role_arn           = "arn:aws:iam::123456789012:role/SecurityLakeIngestRole"
    assume_role_external_id   = "external-id-abc123"
    automatic_schema          = true
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws_security_lake_credentials"
    aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    base_file_name            = "app-logs"
    bucket                    = "security-lake-us-east-1-123456789012"
    custom_source             = "cribl_custom_source"
    deadletter_enabled        = true
    deadletter_path           = "/opt/cribl/state/outputs/dead-letter"
    description               = "Deliver OCSF-compliant logs to Amazon Security Lake"
    duration_seconds          = 3600
    empty_dir_cleanup_sec     = 600
    enable_assume_role        = true
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "https://security-lake.us-east-1.amazonaws.com"
    environment               = "main"
    header_line               = "timestamp,host,level,message"
    id                        = "security_lake_export_prod"
    key_value_metadata = [
      {
        key   = "OCSF Event Class"
        value = "9001"
      }
    ]
    kms_key_id                        = "arn:aws:kms:us-east-1:123456789012:key/abcd-1234-efgh-5678"
    max_closing_files_to_backpressure = 500
    max_concurrent_file_parts         = 5
    max_file_idle_time_sec            = 120
    max_file_open_time_sec            = 600
    max_file_size_mb                  = 256
    max_open_files                    = 200
    max_retry_num                     = 20
    object_acl                        = "private"
    on_backpressure                   = "block"
    on_disk_full_backpressure         = "block"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "128MB"
    parquet_row_group_length          = 100000
    parquet_schema                    = "ocsf_1_1_0"
    parquet_version                   = "PARQUET_2_6"
    pipeline                          = "main"
    region                            = "us-east-1"
    reject_unauthorized               = true
    remove_empty_dirs                 = true
    reuse_connections                 = true
    server_side_encryption            = "aws:kms"
    should_log_invalid_rows           = true
    signature_version                 = "v4"
    stage_path                        = "/opt/cribl/state/outputs/staging"
    storage_class                     = "INTELLIGENT_TIERING"
    streamtags = [
      "prod",
      "securitylake",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type                  = "security_lake"
    verify_permissions    = true
    write_high_water_mark = 256
  }
  output_sentinel = {
    advanced_content_type      = "application/json"
    auth_type                  = "oauth"
    client_id                  = "11111111-2222-3333-4444-555555555555"
    compress                   = true
    concurrency                = 8
    custom_content_type        = "application/x-ndjson"
    custom_drop_when_null      = false
    custom_event_delimiter     = "\n"
    custom_payload_expression  = "`{ \"items\": [${events}] }`"
    custom_source_expression   = "raw=${_raw}"
    dce_endpoint               = "https://mydce-abc123.eastus.ingest.monitor.azure.com"
    dcr_id                     = "12345678-90ab-cdef-1234-567890abcdef"
    description                = "Send events to Microsoft Sentinel (DCR/DCE)"
    endpoint_url_configuration = "url"
    environment                = "main"
    extra_http_headers = [
      {
        name  = "Content-Type"
        value = "application/json"
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    format                      = "ndjson"
    format_event_code           = "if (__e.level === 'error') { __e.__eventOut = JSON.stringify(__e); }"
    format_payload_code         = "__e.__payloadOut = JSON.stringify({ records: __e.payload });"
    id                          = "sentinel-out"
    keep_alive                  = true
    login_url                   = "https://login.microsoftonline.com/<tenant>/oauth2/v2.0/token"
    max_payload_events          = 500
    max_payload_size_kb         = 1000
    on_backpressure             = "queue"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    scope       = "https://monitor.azure.com/.default"
    secret      = "***REDACTED***"
    stream_name = "Custom-MyTable_CL"
    streamtags = [
      "azure",
      "sentinel",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    total_memory_limit_kb = 20480
    type                  = "sentinel"
    url                   = "https://example.dce.ingest.monitor.azure.com"
    use_round_robin_dns   = true
  }
  output_sentinel_one_ai_siem = {
    auth_type                       = "secret"
    base_url                        = "...my_base_url..."
    compress                        = false
    concurrency                     = 31.25
    data_source_category            = "...my_data_source_category..."
    data_source_category_expression = "...my_data_source_category_expression..."
    data_source_name                = "...my_data_source_name..."
    data_source_name_expression     = "...my_data_source_name_expression..."
    data_source_vendor              = "...my_data_source_vendor..."
    data_source_vendor_expression   = "...my_data_source_vendor_expression..."
    description                     = "...my_description..."
    endpoint                        = "/services/collector/event"
    environment                     = "...my_environment..."
    event_type                      = "...my_event_type..."
    event_type_expression           = "...my_event_type_expression..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 9.1
    host                        = "...my_host..."
    host_expression             = "...my_host_expression..."
    id                          = "...my_id..."
    max_payload_events          = 2.86
    max_payload_size_kb         = 1688557.17
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_backpressure_sec           = 7.17
    pq_max_buffer_size                = 752.28
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    pq_rate_per_sec                   = 3.06
    pq_strict_ordering                = false
    region                            = "EMEA"
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 16.04
        http_status     = 446.64
        initial_backoff = 39944.69
        max_backoff     = 147908.36
      }
    ]
    safe_headers = [
      "..."
    ]
    source                 = "...my_source..."
    source_expression      = "...my_source_expression..."
    source_type            = "...my_source_type..."
    source_type_expression = "...my_source_type_expression..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 16.03
      initial_backoff = 80266.12
      max_backoff     = 22615.81
      timeout_retry   = false
    }
    timeout_sec = 4122508657543719
    token       = "...my_token..."
    type        = "sentinel_one_ai_siem"
  }
  output_service_now = {
    auth_token_name    = "lightstep-access-token"
    compress           = "gzip"
    concurrency        = 5
    connection_timeout = 10000
    description        = "Export telemetry to ServiceNow (Lightstep) OTLP ingest"
    endpoint           = "ingest.lightstep.com:443"
    environment        = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode    = "payloadAndHeaders"
    flush_period_sec               = 2
    http_compress                  = "gzip"
    http_logs_endpoint_override    = "https://ingest.lightstep.com/v1/logs"
    http_metrics_endpoint_override = "https://ingest.lightstep.com/v1/metrics"
    http_traces_endpoint_override  = "https://ingest.lightstep.com/v1/traces"
    id                             = "servicenow_otel_export"
    keep_alive                     = true
    keep_alive_time                = 30
    max_payload_size_kb            = 2048
    metadata = [
      {
        key   = "x-tenant-id"
        value = "acme-prod"
      }
    ]
    on_backpressure = "block"
    otlp_version    = "1.3.1"
    pipeline        = "main"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    protocol                          = "grpc"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "servicenow",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/opt/cribl/certs/client.crt"
      certificate_name    = "otel-client"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "s3cr3t"
      priv_key_path       = "/opt/cribl/certs/client.key"
      reject_unauthorized = true
    }
    token_secret        = "servicenow_access_token"
    type                = "service_now"
    use_round_robin_dns = true
  }
  output_signalfx = {
    auth_type   = "manual"
    compress    = true
    concurrency = 8
    description = "Send metrics to Splunk Observability (SignalFx)"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    id                          = "signalfx-out"
    max_payload_events          = 0
    max_payload_size_kb         = 4096
    on_backpressure             = "block"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    realm                             = "us1"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "signalfx",
      "metrics",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    text_secret = "signalfx-api-token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    token               = "***REDACTED***"
    type                = "signalfx"
    use_round_robin_dns = true
  }
  output_snmp = {
    description            = "Forward SNMP traps to network monitoring systems"
    dns_resolve_period_sec = 300
    environment            = "main"
    hosts = [
      {
        host = "snmp01.example.com"
        port = 162
      }
    ]
    id       = "snmp_trap_forwarder"
    pipeline = "main"
    streamtags = [
      "prod",
      "snmp",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type = "snmp"
  }
  output_sns = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/SNSPublisher"
    assume_role_external_id   = "external-id-abc123"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws_sns_credentials"
    aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    description               = "Publish alerts to Amazon SNS FIFO topic"
    duration_seconds          = 3600
    enable_assume_role        = true
    endpoint                  = "https://sns.us-east-1.amazonaws.com"
    environment               = "main"
    id                        = "sns_alerts_prod"
    max_retries               = 5
    message_group_id          = "`alerts-${C.vars.service}`"
    on_backpressure           = "block"
    pipeline                  = "main"
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "100 MB"
    pq_max_size         = "10GB"
    pq_mode             = "backpressure"
    pq_on_backpressure  = "block"
    pq_path             = "/opt/cribl/state/queues"
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    signature_version   = "v4"
    streamtags = [
      "prod",
      "alerts",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    topic_arn = "arn:aws:sns:us-east-1:123456789012:alerts-topic"
    type      = "sns"
  }
  output_splunk = {
    auth_token               = "***REDACTED***"
    auth_type                = "manual"
    compress                 = "auto"
    connection_timeout       = 10000
    description              = "Send events to Splunk indexers over S2S"
    enable_ack               = true
    enable_multi_metrics     = false
    environment              = "main"
    host                     = "splunk-indexer.example.com"
    id                       = "splunk-main"
    log_failed_requests      = false
    max_failed_health_checks = 1
    max_s2_sversion          = "v3"
    nested_fields            = "none"
    on_backpressure          = "block"
    pipeline                 = "default"
    port                     = 9997
    pq_compress              = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    streamtags = [
      "splunk",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    text_secret           = "splunk-indexer-token"
    throttle_rate_per_sec = "50 MB"
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "splunk-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "splunk-indexer.example.com"
    }
    type          = "splunk"
    write_timeout = 60000
  }
  output_splunk_hec = {
    auth_type              = "manual"
    compress               = true
    concurrency            = 8
    description            = "Send events to Splunk HEC"
    dns_resolve_period_sec = 300
    enable_multi_metrics   = false
    environment            = "main"
    exclude_self           = false
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 1
    id                            = "splunk-hec-main"
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    max_payload_events            = 0
    max_payload_size_kb           = 4096
    next_queue                    = "indexQueue"
    on_backpressure               = "block"
    pipeline                      = "default"
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "splunk",
      "hec",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    tcp_routing = "default_route"
    text_secret = "splunk-hec-token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    token       = "***REDACTED***"
    type        = "splunk_hec"
    url         = "http://splunk-hec.example.com:8088/services/collector/event"
    urls = [
      {
        url    = "http://splunk-hec-2.example.com:8088/services/collector/event"
        weight = 2
      }
    ]
    use_round_robin_dns = true
  }
  output_splunk_lb = {
    auth_token             = "***REDACTED***"
    auth_type              = "manual"
    compress               = "auto"
    connection_timeout     = 10000
    description            = "Load-balance events across Splunk indexers"
    dns_resolve_period_sec = 300
    enable_ack             = true
    enable_multi_metrics   = false
    environment            = "main"
    exclude_self           = false
    hosts = [
      {
        host       = "idx1.example.com"
        port       = 9997
        servername = "idx1.example.com"
        tls        = "inherit"
        weight     = 2
      }
    ]
    id                = "splunk-lb-main"
    indexer_discovery = true
    indexer_discovery_configs = {
      auth_token = "***REDACTED***"
      auth_tokens = [
        {
          auth_type = "secret"
        }
      ]
      auth_type            = "manual"
      master_uri           = "https://cm.example.com:8089"
      refresh_interval_sec = 300
      reject_unauthorized  = true
      site                 = "site1"
      text_secret          = "cluster-manager-token"
    }
    load_balance_stats_period_sec = 300
    log_failed_requests           = false
    max_concurrent_senders        = 8
    max_failed_health_checks      = 1
    max_s2_sversion               = "v3"
    nested_fields                 = "none"
    on_backpressure               = "block"
    pipeline                      = "default"
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                = "100 MB"
    pq_max_size                     = "10GB"
    pq_mode                         = "backpressure"
    pq_on_backpressure              = "block"
    pq_path                         = "/opt/cribl/state/queues"
    sender_unhealthy_time_allowance = 500
    streamtags = [
      "splunk",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    text_secret           = "splunk-indexer-token"
    throttle_rate_per_sec = "50 MB"
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "splunk-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "splunk-lb.example.com"
    }
    type          = "splunk_lb"
    write_timeout = 60000
  }
  output_sqs = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/SQSPublisher"
    assume_role_external_id   = "external-id-abc123"
    aws_account_id            = "123456789012"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws_sqs_credentials"
    aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    create_queue              = true
    description               = "Send events to Amazon SQS FIFO queue with batching"
    duration_seconds          = 3600
    enable_assume_role        = true
    endpoint                  = "https://sqs.us-east-1.amazonaws.com"
    environment               = "main"
    flush_period_sec          = 2
    id                        = "sqs_events_prod"
    max_in_progress           = 20
    max_queue_size            = 200
    max_record_size_kb        = 256
    message_group_id          = "logs"
    on_backpressure           = "block"
    pipeline                  = "main"
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "100 MB"
    pq_max_size         = "10GB"
    pq_mode             = "backpressure"
    pq_on_backpressure  = "block"
    pq_path             = "/opt/cribl/state/queues"
    queue_name          = "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue.fifo"
    queue_type          = "fifo"
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    signature_version   = "v4"
    streamtags = [
      "prod",
      "sqs",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    type = "sqs"
  }
  output_statsd = {
    connection_timeout     = 10000
    description            = "Send StatsD metrics to central aggregator"
    dns_resolve_period_sec = 300
    environment            = "main"
    flush_period_sec       = 1
    host                   = "statsd.example.com"
    id                     = "statsd_metrics_prod"
    mtu                    = 1400
    on_backpressure        = "block"
    pipeline               = "metrics"
    port                   = 8125
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    protocol           = "udp"
    streamtags = [
      "prod",
      "statsd",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    throttle_rate_per_sec = "10 MB"
    type                  = "statsd"
    write_timeout         = 30000
  }
  output_statsd_ext = {
    connection_timeout     = 10000
    description            = "Send extended StatsD metrics to external aggregator"
    dns_resolve_period_sec = 300
    environment            = "main"
    flush_period_sec       = 1
    host                   = "statsd-ext.example.com"
    id                     = "statsd_ext_metrics_prod"
    mtu                    = 1400
    on_backpressure        = "block"
    pipeline               = "metrics"
    port                   = 8125
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    protocol           = "udp"
    streamtags = [
      "prod",
      "statsd",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    throttle_rate_per_sec = "10 MB"
    type                  = "statsd_ext"
    write_timeout         = 30000
  }
  output_sumo_logic = {
    compress        = true
    concurrency     = 8
    custom_category = "prod/app/logs"
    custom_source   = "cribl-stream"
    description     = "Send logs to Sumo Logic with retries and batching"
    environment     = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 2
    format                      = "json"
    id                          = "sumologic_logs_prod"
    max_payload_events          = 500
    max_payload_size_kb         = 512
    on_backpressure             = "block"
    pipeline                    = "main"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "sumologic",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    total_memory_limit_kb = 51200
    type                  = "sumo_logic"
    url                   = "https://endpoint1.collection.us2.sumologic.com/receiver/v1/http/ABCDEFG1234567890"
    use_round_robin_dns   = true
  }
  output_syslog = {
    app_name            = "Cribl"
    connection_timeout  = 10000
    description         = "Send syslog to upstream collector"
    environment         = "main"
    facility            = 1
    host                = "syslog.receiver.example.com"
    id                  = "syslog-out"
    load_balanced       = true
    log_failed_requests = false
    max_record_size     = 1200
    message_format      = "rfc3164"
    octet_count_framing = true
    on_backpressure     = "block"
    pipeline            = "default"
    port                = 514
    pq_compress         = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    protocol           = "tcp"
    severity           = 5
    streamtags = [
      "syslog",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    throttle_rate_per_sec = "0"
    timestamp_format      = "syslog"
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "syslog-client-cert"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "syslog.example.com"
    }
    type                       = "syslog"
    udp_dns_resolve_period_sec = 300
    write_timeout              = 60000
  }
  output_tcpjson = {
    auth_token             = "***REDACTED***"
    auth_type              = "manual"
    compression            = "gzip"
    connection_timeout     = 10000
    description            = "Send JSON events over TCP to downstream services"
    dns_resolve_period_sec = 300
    environment            = "main"
    exclude_self           = false
    host                   = "tcp.receiver.example.com"
    hosts = [
      {
        host       = "tcp1.example.com"
        port       = 10300
        servername = "tcp1.example.com"
        tls        = "inherit"
        weight     = 2
      }
    ]
    id                            = "tcpjson-out"
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    log_failed_requests           = false
    max_concurrent_senders        = 8
    on_backpressure               = "block"
    pipeline                      = "default"
    port                          = 10300
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "100 MB"
    pq_max_size        = "10GB"
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "/opt/cribl/state/queues"
    send_header        = true
    streamtags = [
      "tcpjson",
      "prod",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    text_secret           = "tcpjson-auth-token"
    throttle_rate_per_sec = "50 MB"
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "tcpjson-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "tcp.receiver.example.com"
    }
    token_ttl_minutes = 60
    type              = "tcpjson"
    write_timeout     = 60000
  }
  output_wavefront = {
    auth_type   = "manual"
    compress    = true
    concurrency = 8
    description = "Send metrics to WaveFront"
    domain      = "longboard"
    environment = "main"
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "abc123"
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1
    id                          = "wavefront-out"
    max_payload_events          = 0
    max_payload_size_kb         = 4096
    on_backpressure             = "block"
    pipeline                    = "default"
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 503
        initial_backoff = 1000
        max_backoff     = 60000
      }
    ]
    safe_headers = [
      "X-Request-ID",
    ]
    streamtags = [
      "wavefront",
      "metrics",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    text_secret = "wavefront-api-token"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec         = 30
    token               = "***REDACTED***"
    type                = "wavefront"
    use_round_robin_dns = true
  }
  output_webhook = {
    advanced_content_type     = "application/json"
    auth_header_expr          = "`Bearer ${token}`"
    auth_type                 = "token"
    compress                  = true
    concurrency               = 10
    credentials_secret        = "webhook-credentials"
    custom_content_type       = "application/x-ndjson"
    custom_drop_when_null     = false
    custom_event_delimiter    = "\n"
    custom_payload_expression = "`{ \"items\": [${events}] }`"
    custom_source_expression  = "raw=${_raw}"
    description               = "Robust webhook delivery with backoff and retries"
    dns_resolve_period_sec    = 600
    environment               = "main"
    exclude_self              = true
    extra_http_headers = [
      {
        name  = "X-Custom-Header"
        value = "demo"
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 1
    format                        = "ndjson"
    format_event_code             = "if (__e.severity === 'error') { __e.__eventOut = JSON.stringify(__e); }"
    format_payload_code           = "__e.__payloadOut = JSON.stringify({ items: __e.payload });"
    id                            = "webhook-out"
    keep_alive                    = true
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    login_url                     = "https://auth.example.com/oauth/token"
    max_payload_events            = 1000
    max_payload_size_kb           = 8192
    method                        = "POST"
    oauth_headers = [
      {
        name  = "Accept"
        value = "application/json"
      }
    ]
    oauth_params = [
      {
        name  = "grant_type"
        value = "client_credentials"
      }
    ]
    on_backpressure = "queue"
    password        = "***REDACTED***"
    pipeline        = "default"
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "X-Trace-Id",
    ]
    secret            = "s3cr3t"
    secret_param_name = "client_secret"
    streamtags = [
      "webhook",
    ]
    system_fields = [
      "cribl_pipe",
    ]
    text_secret = "webhook-token-secret"
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec = 30
    tls = {
      ca_path          = "/etc/ssl/certs/ca-bundle.crt"
      cert_path        = "/etc/ssl/certs/client.crt"
      certificate_name = "webhook-client-cert"
      disabled         = true
      max_version      = "TLSv1.3"
      min_version      = "TLSv1.2"
      passphrase       = "***REDACTED***"
      priv_key_path    = "/etc/ssl/private/client.key"
      servername       = "api.example.com"
    }
    token                 = "***REDACTED***"
    token_attribute_name  = "access_token"
    token_timeout_secs    = 3600
    total_memory_limit_kb = 20480
    type                  = "webhook"
    url                   = "https://hooks.example.com/ingest"
    urls = [
      {
        url    = "https://hooks1.example.com/ingest"
        weight = 2
      }
    ]
    use_round_robin_dns = true
    username            = "api-user"
  }
  output_xsiam = {
    auth_type              = "secret"
    concurrency            = 8
    description            = "Send logs to Palo Alto Networks XSIAM with token auth"
    dns_resolve_period_sec = 300
    environment            = "main"
    exclude_self           = false
    extra_http_headers = [
      {
        name  = "X-Request-ID"
        value = "123e4567-e89b-12d3-a456-426614174000"
      }
    ]
    failed_request_logging_mode   = "payloadAndHeaders"
    flush_period_sec              = 2
    id                            = "xsiam_export_prod"
    load_balance_stats_period_sec = 300
    load_balanced                 = true
    max_payload_events            = 2000
    max_payload_size_kb           = 8192
    on_backpressure               = "block"
    pipeline                      = "main"
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "100 MB"
    pq_max_size                       = "10GB"
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "/opt/cribl/state/queues"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2
        http_status     = 429
        initial_backoff = 1000
        max_backoff     = 30000
      }
    ]
    safe_headers = [
      "content-type",
      "x-request-id",
    ]
    streamtags = [
      "prod",
      "xsiam",
    ]
    system_fields = [
      "cribl_pipe",
      "cribl_breaker",
    ]
    text_secret               = "xsiam_token"
    throttle_rate_req_per_sec = 500
    timeout_retry_settings = {
      backoff_rate    = 2
      initial_backoff = 1000
      max_backoff     = 30000
      timeout_retry   = true
    }
    timeout_sec           = 30
    token                 = "xsiam-0123456789abcdef0123456789abcdef"
    total_memory_limit_kb = 51200
    type                  = "xsiam"
    url                   = "https://api-tenant.paloaltonetworks.com/logs/v1/event"
    urls = [
      {
        url    = "{ \"see\": \"documentation\" }"
        weight = 8.62
      }
    ]
    use_round_robin_dns = true
  }
  pack = "observability-pack"
}