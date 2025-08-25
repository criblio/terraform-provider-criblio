resource "criblio_pack_destination" "my_packdestination" {
  group_id = "...my_group_id..."
  id       = "...my_id..."
  output_azure_blob = {
    add_id_to_stage_path = true
    auth_type            = "manual"
    automatic_schema     = true
    azure_cloud          = "...my_azure_cloud..."
    base_file_name       = "...my_base_file_name..."
    certificate = {
      certificate_name = "...my_certificate_name..."
    }
    client_id               = "...my_client_id..."
    client_text_secret      = "...my_client_text_secret..."
    compress                = "gzip"
    compression_level       = "best_speed"
    connection_string       = "...my_connection_string..."
    container_name          = "...my_container_name..."
    create_container        = true
    deadletter_enabled      = true
    deadletter_path         = "...my_deadletter_path..."
    description             = "...my_description..."
    dest_path               = "...my_dest_path..."
    empty_dir_cleanup_sec   = 55323.17
    enable_page_checksum    = false
    enable_statistics       = false
    enable_write_page_index = false
    endpoint_suffix         = "...my_endpoint_suffix..."
    environment             = "...my_environment..."
    file_name_suffix        = "...my_file_name_suffix..."
    format                  = "raw"
    header_line             = "...my_header_line..."
    id                      = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_concurrent_file_parts = 8.23
    max_file_idle_time_sec    = 753.81
    max_file_open_time_sec    = 1092.11
    max_file_size_mb          = 408.27
    max_open_files            = 26.97
    max_retry_num             = 5.36
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V1"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 1327800.98
    parquet_version           = "PARQUET_1_0"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    remove_empty_dirs         = true
    should_log_invalid_rows   = false
    stage_path                = "...my_stage_path..."
    storage_account_name      = "...my_storage_account_name..."
    storage_class             = "Archive"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tenant_id             = "...my_tenant_id..."
    text_secret           = "...my_text_secret..."
    type                  = "azure_blob"
    write_high_water_mark = 1674.32
  }
  output_azure_data_explorer = {
    add_id_to_stage_path = true
    additional_properties = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    certificate = {
      certificate_name = "...my_certificate_name..."
    }
    client_id          = "...my_client_id..."
    client_secret      = "...my_client_secret..."
    cluster_url        = "...my_cluster_url..."
    compress           = "none"
    concurrency        = 25.3
    database           = "...my_database..."
    deadletter_enabled = true
    description        = "...my_description..."
    environment        = "...my_environment..."
    extent_tags = [
      {
        prefix = "dropBy"
        value  = "...my_value..."
      }
    ]
    file_name_suffix  = "...my_file_name_suffix..."
    flush_immediately = true
    flush_period_sec  = 5.98
    format            = "json"
    id                = "...my_id..."
    ingest_if_not_exists = [
      {
        value = "...my_value..."
      }
    ]
    ingest_mode               = "streaming"
    ingest_url                = "...my_ingest_url..."
    is_mapping_obj            = false
    keep_alive                = true
    mapping_ref               = "...my_mapping_ref..."
    max_concurrent_file_parts = 7.76
    max_file_idle_time_sec    = 636.51
    max_file_open_time_sec    = 1785.56
    max_file_size_mb          = 49.98
    max_open_files            = 913.7
    max_payload_events        = 4.65
    max_payload_size_kb       = 1906.66
    oauth_endpoint            = "https://login.microsoftonline.com"
    oauth_type                = "certificate"
    on_backpressure           = "queue"
    on_disk_full_backpressure = "drop"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    remove_empty_dirs                 = true
    report_level                      = "doNotReport"
    report_method                     = "queueAndTable"
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 11.6
        http_status     = 251.57
        initial_backoff = 223817.9
        max_backoff     = 166141.84
      }
    ]
    retain_blob_on_success = true
    scope                  = "...my_scope..."
    stage_path             = "...my_stage_path..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    table       = "...my_table..."
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 15.6
      initial_backoff = 71830.59
      max_backoff     = 118420.63
      timeout_retry   = true
    }
    timeout_sec                = 1821283157205735.2
    type                       = "azure_data_explorer"
    use_round_robin_dns        = true
    validate_database_settings = true
  }
  output_azure_eventhub = {
    ack                    = 10
    authentication_timeout = 430621.06
    backoff_rate           = 5.96
    brokers = [
      "..."
    ]
    connection_timeout = 180818.45
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_event_count  = 9013.17
    flush_period_sec   = 9.34
    format             = "raw"
    id                 = "...my_id..."
    initial_backoff    = 373583.55
    max_back_off       = 55727.44
    max_record_size_kb = 3.25
    max_retries        = 55.57
    on_backpressure    = "block"
    pipeline           = "...my_pipeline..."
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "always"
    pq_on_backpressure         = "block"
    pq_path                    = "...my_pq_path..."
    reauthentication_threshold = 230505.87
    request_timeout            = 2476956.46
    sasl = {
      disabled  = false
      mechanism = "plain"
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
    type  = "azure_eventhub"
  }
  output_azure_logs = {
    api_url     = "...my_api_url..."
    auth_type   = "secret"
    compress    = true
    concurrency = 23.37
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 2.44
    id                          = "...my_id..."
    keypair_secret              = "...my_keypair_secret..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 0.59
    max_payload_size_kb         = 4642.62
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    resource_id                       = "...my_resource_id..."
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 18.09
        http_status     = 234.32
        initial_backoff = 334064.95
        max_backoff     = 165677.12
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 1.28
      initial_backoff = 379283.94
      max_backoff     = 135310.21
      timeout_retry   = false
    }
    timeout_sec         = 4274665945637552
    type                = "azure_logs"
    use_round_robin_dns = true
    workspace_id        = "...my_workspace_id..."
    workspace_key       = "...my_workspace_key..."
  }
  output_click_house = {
    async_inserts    = true
    auth_header_expr = "...my_auth_header_expr..."
    auth_type        = "none"
    column_mappings = [
      {
        column_name             = "...my_column_name..."
        column_type             = "...my_column_type..."
        column_value_expression = "...my_column_value_expression..."
      }
    ]
    compress                   = true
    concurrency                = 9.96
    credentials_secret         = "...my_credentials_secret..."
    database                   = "...my_database..."
    describe_table             = "...my_describe_table..."
    description                = "...my_description..."
    dump_format_errors_to_disk = false
    environment                = "...my_environment..."
    exclude_mapping_fields = [
    ]
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 5.87
    format                      = "json-compact-each-row-with-names"
    id                          = "...my_id..."
    login_url                   = "...my_login_url..."
    mapping_type                = "automatic"
    max_payload_events          = 2.57
    max_payload_size_kb         = 9715.84
    oauth_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    oauth_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "queue"
    password        = "...my_password..."
    pipeline        = "...my_pipeline..."
    pq_compress     = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 16.8
        http_status     = 425.87
        initial_backoff = 273904.6
        max_backoff     = 85192.14
      }
    ]
    safe_headers = [
      "..."
    ]
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    sql_username      = "...my_sql_username..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    table_name  = "...my_table_name..."
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 11.96
      initial_backoff = 42734.45
      max_backoff     = 112843.03
      timeout_retry   = false
    }
    timeout_sec = 3517606791143139
    tls = {
      ca_path          = "...my_ca_path..."
      cert_path        = "...my_cert_path..."
      certificate_name = "...my_certificate_name..."
      disabled         = false
      max_version      = "TLSv1.2"
      min_version      = "TLSv1"
      passphrase       = "...my_passphrase..."
      priv_key_path    = "...my_priv_key_path..."
      servername       = "...my_servername..."
    }
    token                  = "...my_token..."
    token_attribute_name   = "...my_token_attribute_name..."
    token_timeout_secs     = 184047.77
    type                   = "click_house"
    url                    = "...my_url..."
    use_round_robin_dns    = true
    username               = "...my_username..."
    wait_for_async_inserts = true
  }
  output_cloudwatch = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "manual"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    description               = "...my_description..."
    duration_seconds          = 5242.24
    enable_assume_role        = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    flush_period_sec          = 3.27
    id                        = "...my_id..."
    log_group_name            = "...my_log_group_name..."
    log_stream_name           = "...my_log_stream_name..."
    max_queue_size            = 1.07
    max_record_size_kb        = 6210.86
    on_backpressure           = "block"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "always"
    pq_on_backpressure  = "block"
    pq_path             = "...my_pq_path..."
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = true
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "cloudwatch"
  }
  output_confluent_cloud = {
    ack                    = 8
    authentication_timeout = 1295471.86
    backoff_rate           = 18.05
    brokers = [
      "..."
    ]
    compression        = "none"
    connection_timeout = 1112740.74
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_event_count  = 5876.29
    flush_period_sec   = 8.16
    format             = "protobuf"
    id                 = "...my_id..."
    initial_backoff    = 106771.36
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout      = 46613.5
      default_key_schema_id   = 6.1
      default_value_schema_id = 7.4
      disabled                = true
      max_retries             = 24.37
      request_timeout         = 11546.95
      schema_registry_url     = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = true
        max_version         = "TLSv1.1"
        min_version         = "TLSv1"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = false
        servername          = "...my_servername..."
      }
    }
    max_back_off       = 110622.93
    max_record_size_kb = 9.39
    max_retries        = 8.46
    on_backpressure    = "drop"
    pipeline           = "...my_pipeline..."
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "backpressure"
    pq_on_backpressure         = "block"
    pq_path                    = "...my_pq_path..."
    protobuf_library_id        = "...my_protobuf_library_id..."
    reauthentication_threshold = 1366964.01
    request_timeout            = 1002698.61
    sasl = {
      disabled  = false
      mechanism = "scram-sha-512"
    }
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    topic = "...my_topic..."
    type  = "confluent_cloud"
  }
  output_cribl_http = {
    compression            = "gzip"
    concurrency            = 30.41
    description            = "...my_description..."
    dns_resolve_period_sec = 44049.28
    environment            = "...my_environment..."
    exclude_fields = [
      "..."
    ]
    exclude_self = false
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "none"
    flush_period_sec              = 8.75
    id                            = "...my_id..."
    load_balance_stats_period_sec = 19.71
    load_balanced                 = true
    max_payload_events            = 6.38
    max_payload_size_kb           = 1747.03
    on_backpressure               = "drop"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 2.3
        http_status     = 585.05
        initial_backoff = 464186.35
        max_backoff     = 76356.76
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 9.19
      initial_backoff = 292539.69
      max_backoff     = 161617.7
      timeout_retry   = true
    }
    timeout_sec = 8642518876171064
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    token_ttl_minutes = 17.49
    type              = "cribl_http"
    url               = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 0.17
      }
    ]
    use_round_robin_dns = true
  }
  output_cribl_lake = {
    description = "Cribl Lake destination"
    dest_path   = "test_lake_dataset"
    id          = "test_lake_output"
    type        = "cribl_lake"
  }
  output_cribl_tcp = {
    compression            = "none"
    connection_timeout     = 8
    description            = "...my_description..."
    dns_resolve_period_sec = 77972.9
    environment            = "...my_environment..."
    exclude_fields = [
      "..."
    ]
    exclude_self = false
    host         = "...my_host..."
    hosts = [
      {
        host       = "...my_host..."
        port       = 61317.22
        servername = "...my_servername..."
        tls        = "inherit"
        weight     = 1.12
      }
    ]
    id                            = "...my_id..."
    load_balance_stats_period_sec = 10.34
    load_balanced                 = false
    log_failed_requests           = false
    max_concurrent_senders        = 4.57
    on_backpressure               = "queue"
    pipeline                      = "...my_pipeline..."
    port                          = 3105.16
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "backpressure"
    pq_on_backpressure = "drop"
    pq_path            = "...my_pq_path..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.2"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    token_ttl_minutes = 25.77
    type              = "cribl_tcp"
    write_timeout     = 0.12
  }
  output_crowdstrike_next_gen_siem = {
    auth_type   = "manual"
    compress    = true
    concurrency = 1.21
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 5.55
    format                      = "raw"
    id                          = "...my_id..."
    max_payload_events          = 3.06
    max_payload_size_kb         = 24905.36
    on_backpressure             = "drop"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 1.99
        http_status     = 504.25
        initial_backoff = 415593.64
        max_backoff     = 69028.18
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 9.38
      initial_backoff = 492268.73
      max_backoff     = 40070.73
      timeout_retry   = true
    }
    timeout_sec         = 4075364216465747
    token               = "...my_token..."
    type                = "crowdstrike_next_gen_siem"
    url                 = "...my_url..."
    use_round_robin_dns = false
  }
  output_datadog = {
    allow_api_key_from_events = false
    api_key                   = "...my_api_key..."
    auth_type                 = "secret"
    batch_by_tags             = true
    compress                  = true
    concurrency               = 24.8
    content_type              = "text"
    custom_url                = "...my_custom_url..."
    description               = "...my_description..."
    environment               = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 7.78
    host                        = "...my_host..."
    id                          = "...my_id..."
    max_payload_events          = 6.1
    max_payload_size_kb         = 6408.66
    message                     = "...my_message..."
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 16.24
        http_status     = 143.86
        initial_backoff = 570591.46
        max_backoff     = 108293.13
      }
    ]
    safe_headers = [
      "..."
    ]
    send_counters_as_count = true
    service                = "...my_service..."
    severity               = "notice"
    site                   = "eu"
    source                 = "...my_source..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 10.99
      initial_backoff = 518300.67
      max_backoff     = 158957.39
      timeout_retry   = true
    }
    timeout_sec           = 5081753886598422
    total_memory_limit_kb = 1.34
    type                  = "datadog"
    use_round_robin_dns   = false
  }
  output_dataset = {
    api_key          = "...my_api_key..."
    auth_type        = "manual"
    compress         = false
    concurrency      = 13.68
    custom_url       = "...my_custom_url..."
    default_severity = "warning"
    description      = "...my_description..."
    environment      = "...my_environment..."
    exclude_fields = [
      "..."
    ]
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 9.72
    id                          = "...my_id..."
    max_payload_events          = 8.93
    max_payload_size_kb         = 5484.69
    message_field               = "...my_message_field..."
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 5.29
        http_status     = 396.18
        initial_backoff = 237642.4
        max_backoff     = 31652.38
      }
    ]
    safe_headers = [
      "..."
    ]
    server_host_field = "...my_server_host_field..."
    site              = "eu"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 12.05
      initial_backoff = 227296.99
      max_backoff     = 132128.05
      timeout_retry   = false
    }
    timeout_sec           = 5175919553799471
    timestamp_field       = "...my_timestamp_field..."
    total_memory_limit_kb = 8.24
    type                  = "dataset"
    use_round_robin_dns   = false
  }
  output_default = {
    default_id  = "...my_default_id..."
    environment = "...my_environment..."
    id          = "...my_id..."
    pipeline    = "...my_pipeline..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "default"
  }
  output_devnull = {
    environment = "...my_environment..."
    id          = "...my_id..."
    pipeline    = "...my_pipeline..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "devnull"
  }
  output_disk_spool = {
    compress       = "none"
    description    = "...my_description..."
    environment    = "...my_environment..."
    id             = "...my_id..."
    max_data_size  = "...my_max_data_size..."
    max_data_time  = "...my_max_data_time..."
    partition_expr = "...my_partition_expr..."
    pipeline       = "...my_pipeline..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    time_window = "...my_time_window..."
    type        = "disk_spool"
  }
  output_dl_s3 = {
    add_id_to_stage_path      = true
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    automatic_schema          = false
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    compress                  = "gzip"
    compression_level         = "best_speed"
    deadletter_enabled        = false
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    dest_path                 = "...my_dest_path..."
    duration_seconds          = 15880.81
    empty_dir_cleanup_sec     = 72778.15
    enable_assume_role        = true
    enable_page_checksum      = true
    enable_statistics         = false
    enable_write_page_index   = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    format                    = "parquet"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    kms_key_id                        = "...my_kms_key_id..."
    max_closing_files_to_backpressure = 4116.96
    max_concurrent_file_parts         = 3.62
    max_file_idle_time_sec            = 1646.63
    max_file_open_time_sec            = 1294.93
    max_file_size_mb                  = 291.56
    max_open_files                    = 1527.44
    max_retry_num                     = 6.4
    object_acl                        = "public-read-write"
    on_backpressure                   = "block"
    on_disk_full_backpressure         = "block"
    parquet_data_page_version         = "DATA_PAGE_V1"
    parquet_page_size                 = "...my_parquet_page_size..."
    parquet_row_group_length          = 36789571.94
    parquet_version                   = "PARQUET_2_4"
    partitioning_fields = [
      "..."
    ]
    pipeline                = "...my_pipeline..."
    region                  = "...my_region..."
    reject_unauthorized     = true
    remove_empty_dirs       = true
    reuse_connections       = false
    server_side_encryption  = "AES256"
    should_log_invalid_rows = true
    signature_version       = "v4"
    stage_path              = "...my_stage_path..."
    storage_class           = "DEEP_ARCHIVE"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "dl_s3"
    verify_permissions    = false
    write_high_water_mark = 1462.01
  }
  output_dynatrace_http = {
    active_gate_domain = "...my_active_gate_domain..."
    auth_type          = "token"
    compress           = true
    concurrency        = 20.98
    description        = "...my_description..."
    endpoint           = "cloud"
    environment        = "...my_environment..."
    environment_id     = "...my_environment_id..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 8.42
    format                      = "plaintext"
    id                          = "...my_id..."
    keep_alive                  = false
    max_payload_events          = 5679.23
    max_payload_size_kb         = 4322.11
    method                      = "PUT"
    on_backpressure             = "drop"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 4.8
        http_status     = 159.39
        initial_backoff = 148695.44
        max_backoff     = 53270.51
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    telemetry_type = "logs"
    text_secret    = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 12.37
      initial_backoff = 165905.65
      max_backoff     = 161909.34
      timeout_retry   = false
    }
    timeout_sec           = 5587058779806179
    token                 = "...my_token..."
    total_memory_limit_kb = 5.36
    type                  = "dynatrace_http"
    url                   = "...my_url..."
    use_round_robin_dns   = true
  }
  output_dynatrace_otlp = {
    auth_token_name    = "...my_auth_token_name..."
    compress           = "gzip"
    concurrency        = 23.13
    connection_timeout = 7.89
    description        = "...my_description..."
    endpoint           = "...my_endpoint..."
    endpoint_type      = "saas"
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode    = "none"
    flush_period_sec               = 1.29
    http_compress                  = "gzip"
    http_logs_endpoint_override    = "...my_http_logs_endpoint_override..."
    http_metrics_endpoint_override = "...my_http_metrics_endpoint_override..."
    http_traces_endpoint_override  = "...my_http_traces_endpoint_override..."
    id                             = "...my_id..."
    keep_alive                     = false
    keep_alive_time                = 3.06
    max_payload_size_kb            = 1992.01
    metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "block"
    otlp_version    = "1.3.1"
    pipeline        = "...my_pipeline..."
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    protocol                          = "http"
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 15.6
        http_status     = 541.15
        initial_backoff = 33374.24
        max_backoff     = 55428.96
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 18.36
      initial_backoff = 385892.16
      max_backoff     = 50280.25
      timeout_retry   = false
    }
    timeout_sec         = 4546289421435046
    token_secret        = "...my_token_secret..."
    type                = "dynatrace_otlp"
    use_round_robin_dns = false
  }
  output_elastic = {
    auth = {
      auth_type = "secret"
      disabled  = false
    }
    compress               = false
    concurrency            = 28.73
    description            = "...my_description..."
    dns_resolve_period_sec = 31696.4
    doc_type               = "...my_doc_type..."
    elastic_pipeline       = "...my_elastic_pipeline..."
    elastic_version        = "auto"
    environment            = "...my_environment..."
    exclude_self           = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    extra_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "none"
    flush_period_sec              = 3.65
    id                            = "...my_id..."
    include_doc_id                = false
    index                         = "...my_index..."
    load_balance_stats_period_sec = 12.6
    load_balanced                 = true
    max_payload_events            = 0.85
    max_payload_size_kb           = 91141.32
    on_backpressure               = "block"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 17.49
        http_status     = 371.74
        initial_backoff = 164706.64
        max_backoff     = 102099.58
      }
    ]
    retry_partial_errors = false
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 9.56
      initial_backoff = 522024.45
      max_backoff     = 33017.16
      timeout_retry   = false
    }
    timeout_sec = 739433064426379.6
    type        = "elastic"
    url         = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 3.56
      }
    ]
    use_round_robin_dns = false
    write_action        = "create"
  }
  output_elastic_cloud = {
    auth = {
      auth_type = "manualAPIKey"
      disabled  = true
    }
    compress         = false
    concurrency      = 27.58
    description      = "...my_description..."
    elastic_pipeline = "...my_elastic_pipeline..."
    environment      = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    extra_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 8.32
    id                          = "...my_id..."
    include_doc_id              = true
    index                       = "...my_index..."
    max_payload_events          = 3.95
    max_payload_size_kb         = 45448.07
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 17.11
        http_status     = 566.5
        initial_backoff = 261851.45
        max_backoff     = 86436.83
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 19.27
      initial_backoff = 487942.83
      max_backoff     = 119811.72
      timeout_retry   = true
    }
    timeout_sec = 1486952552692243.5
    type        = "elastic_cloud"
    url         = "...my_url..."
  }
  output_exabeam = {
    add_id_to_stage_path      = true
    aws_api_key               = "...my_aws_api_key..."
    aws_secret_key            = "...my_aws_secret_key..."
    bucket                    = "...my_bucket..."
    collector_instance_id     = "...my_collector_instance_id..."
    deadletter_enabled        = true
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    empty_dir_cleanup_sec     = 59606.21
    encoded_configuration     = "...my_encoded_configuration..."
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    id                        = "...my_id..."
    max_file_idle_time_sec    = 812.6
    max_file_open_time_sec    = 695.94
    max_file_size_mb          = 720.01
    max_open_files            = 1997.19
    max_retry_num             = 4.9
    object_acl                = "bucket-owner-read"
    on_backpressure           = "block"
    on_disk_full_backpressure = "drop"
    pipeline                  = "...my_pipeline..."
    region                    = "...my_region..."
    reject_unauthorized       = false
    remove_empty_dirs         = true
    reuse_connections         = false
    signature_version         = "v4"
    site_id                   = "...my_site_id..."
    site_name                 = "...my_site_name..."
    stage_path                = "...my_stage_path..."
    storage_class             = "STANDARD"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timezone_offset = "...my_timezone_offset..."
    type            = "exabeam"
  }
  output_filesystem = {
    add_id_to_stage_path    = false
    automatic_schema        = true
    base_file_name          = "...my_base_file_name..."
    compress                = "gzip"
    compression_level       = "best_speed"
    deadletter_enabled      = true
    deadletter_path         = "...my_deadletter_path..."
    description             = "...my_description..."
    dest_path               = "...my_dest_path..."
    empty_dir_cleanup_sec   = 61451.75
    enable_page_checksum    = false
    enable_statistics       = false
    enable_write_page_index = false
    environment             = "...my_environment..."
    file_name_suffix        = "...my_file_name_suffix..."
    format                  = "parquet"
    header_line             = "...my_header_line..."
    id                      = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_file_idle_time_sec    = 1289.27
    max_file_open_time_sec    = 304.92
    max_file_size_mb          = 338.63
    max_open_files            = 1766.53
    max_retry_num             = 1.83
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V1"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 24426080.91
    parquet_version           = "PARQUET_2_4"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    remove_empty_dirs         = false
    should_log_invalid_rows   = false
    stage_path                = "...my_stage_path..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "filesystem"
    write_high_water_mark = 516.14
  }
  output_google_chronicle = {
    api_key               = "...my_api_key..."
    api_key_secret        = "...my_api_key_secret..."
    api_version           = "v2"
    authentication_method = "manual"
    compress              = false
    concurrency           = 12.09
    custom_labels = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    customer_id = "...my_customer_id..."
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    extra_log_types = [
      {
        description = "...my_description..."
        log_type    = "...my_log_type..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 6.88
    id                          = "...my_id..."
    log_format_type             = "unstructured"
    log_text_field              = "...my_log_text_field..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 5.74
    max_payload_size_kb         = 124.53
    namespace                   = "...my_namespace..."
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    region                            = "...my_region..."
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 19.4
        http_status     = 324.58
        initial_backoff = 413439.35
        max_backoff     = 68567.56
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
      backoff_rate    = 4.63
      initial_backoff = 22588.87
      max_backoff     = 61071.63
      timeout_retry   = true
    }
    timeout_sec           = 4524053444580176
    total_memory_limit_kb = 3.9
    type                  = "google_chronicle"
    use_round_robin_dns   = true
  }
  output_google_cloud_logging = {
    cache_fill_bytes_expression = "...my_cache_fill_bytes_expression..."
    cache_hit_expression        = "...my_cache_hit_expression..."
    cache_lookup_expression     = "...my_cache_lookup_expression..."
    cache_validated_expression  = "...my_cache_validated_expression..."
    concurrency                 = 2.1
    connection_timeout          = 1.99
    description                 = "...my_description..."
    environment                 = "...my_environment..."
    file_expression             = "...my_file_expression..."
    first_expression            = "...my_first_expression..."
    flush_period_sec            = 6.41
    function_expression         = "...my_function_expression..."
    google_auth_method          = "secret"
    id                          = "...my_id..."
    id_expression               = "...my_id_expression..."
    index_expression            = "...my_index_expression..."
    insert_id_expression        = "...my_insert_id_expression..."
    last_expression             = "...my_last_expression..."
    latency_expression          = "...my_latency_expression..."
    line_expression             = "...my_line_expression..."
    log_labels = [
      {
        label            = "...my_label..."
        value_expression = "...my_value_expression..."
      }
    ]
    log_location_expression = "...my_log_location_expression..."
    log_location_type       = "project"
    log_name_expression     = "...my_log_name_expression..."
    max_payload_events      = 1.5
    max_payload_size_kb     = 8058.64
    on_backpressure         = "drop"
    payload_expression      = "...my_payload_expression..."
    payload_format          = "text"
    pipeline                = "...my_pipeline..."
    pq_compress             = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size          = "...my_pq_max_file_size..."
    pq_max_size               = "...my_pq_max_size..."
    pq_mode                   = "error"
    pq_on_backpressure        = "drop"
    pq_path                   = "...my_pq_path..."
    producer_expression       = "...my_producer_expression..."
    protocol_expression       = "...my_protocol_expression..."
    referer_expression        = "...my_referer_expression..."
    remote_ip_expression      = "...my_remote_ip_expression..."
    request_method_expression = "...my_request_method_expression..."
    request_size_expression   = "...my_request_size_expression..."
    request_url_expression    = "...my_request_url_expression..."
    resource_type_expression  = "...my_resource_type_expression..."
    resource_type_labels = [
      {
        label            = "...my_label..."
        value_expression = "...my_value_expression..."
      }
    ]
    response_size_expression    = "...my_response_size_expression..."
    secret                      = "...my_secret..."
    server_ip_expression        = "...my_server_ip_expression..."
    service_account_credentials = "...my_service_account_credentials..."
    severity_expression         = "...my_severity_expression..."
    span_id_expression          = "...my_span_id_expression..."
    status_expression           = "...my_status_expression..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_req_per_sec = 920
    timeout_sec               = 1797415024310893.5
    total_memory_limit_kb     = 7.07
    total_splits_expression   = "...my_total_splits_expression..."
    trace_expression          = "...my_trace_expression..."
    trace_sampled_expression  = "...my_trace_sampled_expression..."
    type                      = "google_cloud_logging"
    uid_expression            = "...my_uid_expression..."
    user_agent_expression     = "...my_user_agent_expression..."
  }
  output_google_cloud_storage = {
    add_id_to_stage_path      = true
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
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
    empty_dir_cleanup_sec     = 66967.85
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    format                    = "raw"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_file_idle_time_sec    = 1524.23
    max_file_open_time_sec    = 458.74
    max_file_size_mb          = 296.35
    max_open_files            = 987.61
    max_retry_num             = 2.79
    object_acl                = "public-read"
    on_backpressure           = "block"
    on_disk_full_backpressure = "drop"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 54910275.13
    parquet_version           = "PARQUET_2_4"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    region                    = "...my_region..."
    reject_unauthorized       = false
    remove_empty_dirs         = true
    reuse_connections         = true
    should_log_invalid_rows   = false
    signature_version         = "v2"
    stage_path                = "...my_stage_path..."
    storage_class             = "COLDLINE"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "google_cloud_storage"
    verify_permissions    = false
    write_high_water_mark = 1733.06
  }
  output_google_pubsub = {
    batch_size         = 7384.9
    batch_timeout      = 11786.48
    create_topic       = true
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_period_sec   = "{ \"see\": \"documentation\" }"
    google_auth_method = "secret"
    id                 = "...my_id..."
    max_in_progress    = 11.51
    max_queue_size     = 9.76
    max_record_size_kb = 85.85
    on_backpressure    = "drop"
    ordered_delivery   = false
    pipeline           = "...my_pipeline..."
    pq_compress        = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size            = "...my_pq_max_file_size..."
    pq_max_size                 = "...my_pq_max_size..."
    pq_mode                     = "backpressure"
    pq_on_backpressure          = "block"
    pq_path                     = "...my_pq_path..."
    region                      = "...my_region..."
    secret                      = "...my_secret..."
    service_account_credentials = "...my_service_account_credentials..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    topic_name = "...my_topic_name..."
    type       = "google_pubsub"
  }
  output_grafana_cloud = {
    output_grafana_cloud_grafana_cloud1 = {
      compress    = false
      concurrency = 26.26
      description = "...my_description..."
      environment = "...my_environment..."
      extra_http_headers = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      failed_request_logging_mode = "payload"
      flush_period_sec            = 9.23
      id                          = "...my_id..."
      labels = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      loki_auth = {
        auth_type          = "basic"
        credentials_secret = "...my_credentials_secret..."
        password           = "...my_password..."
        text_secret        = "...my_text_secret..."
        token              = "...my_token..."
        username           = "...my_username..."
      }
      loki_url            = "...my_loki_url..."
      max_payload_events  = 0.89
      max_payload_size_kb = 9058.73
      message             = "...my_message..."
      message_format      = "json"
      metric_rename_expr  = "...my_metric_rename_expr..."
      on_backpressure     = "block"
      pipeline            = "...my_pipeline..."
      pq_compress         = "gzip"
      pq_controls = {
        # ...
      }
      pq_max_file_size   = "...my_pq_max_file_size..."
      pq_max_size        = "...my_pq_max_size..."
      pq_mode            = "backpressure"
      pq_on_backpressure = "block"
      pq_path            = "...my_pq_path..."
      prometheus_auth = {
        auth_type          = "credentialsSecret"
        credentials_secret = "...my_credentials_secret..."
        password           = "...my_password..."
        text_secret        = "...my_text_secret..."
        token              = "...my_token..."
        username           = "...my_username..."
      }
      prometheus_url                    = "...my_prometheus_url..."
      reject_unauthorized               = true
      response_honor_retry_after_header = true
      response_retry_settings = [
        {
          backoff_rate    = 19.41
          http_status     = 287.25
          initial_backoff = 173045.8
          max_backoff     = 114128.22
        }
      ]
      safe_headers = [
        "..."
      ]
      streamtags = [
        "..."
      ]
      system_fields = [
        "..."
      ]
      timeout_retry_settings = {
        backoff_rate    = 19.16
        initial_backoff = 524925.01
        max_backoff     = 165357.85
        timeout_retry   = false
      }
      timeout_sec         = 8410084066962944
      type                = "grafana_cloud"
      use_round_robin_dns = true
    }
  }
  output_graphite = {
    connection_timeout     = 7.68
    description            = "...my_description..."
    dns_resolve_period_sec = 73838.71
    environment            = "...my_environment..."
    flush_period_sec       = 0
    host                   = "...my_host..."
    id                     = "...my_id..."
    mtu                    = 20549.24
    on_backpressure        = "drop"
    pipeline               = "...my_pipeline..."
    port                   = 26769.32
    pq_compress            = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "error"
    pq_on_backpressure = "drop"
    pq_path            = "...my_pq_path..."
    protocol           = "tcp"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    type                  = "graphite"
    write_timeout         = 0.71
  }
  output_honeycomb = {
    auth_type   = "manual"
    compress    = false
    concurrency = 7.42
    dataset     = "...my_dataset..."
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 4.27
    id                          = "...my_id..."
    max_payload_events          = 5.54
    max_payload_size_kb         = 4480.71
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 17.81
        http_status     = 574.66
        initial_backoff = 172747.51
        max_backoff     = 166242.49
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    team        = "...my_team..."
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 11.71
      initial_backoff = 183167.48
      max_backoff     = 17062.56
      timeout_retry   = false
    }
    timeout_sec         = 4127916144665701
    type                = "honeycomb"
    use_round_robin_dns = true
  }
  output_humio_hec = {
    auth_type   = "manual"
    compress    = true
    concurrency = 27.97
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 7.69
    format                      = "JSON"
    id                          = "...my_id..."
    max_payload_events          = 6.83
    max_payload_size_kb         = 32653.9
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 4.36
        http_status     = 134.47
        initial_backoff = 526947.65
        max_backoff     = 112451.32
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 19.48
      initial_backoff = 366619.61
      max_backoff     = 145965.16
      timeout_retry   = false
    }
    timeout_sec         = 8089370599158853
    token               = "...my_token..."
    type                = "humio_hec"
    url                 = "...my_url..."
    use_round_robin_dns = false
  }
  output_influxdb = {
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "credentialsSecret"
    bucket                   = "...my_bucket..."
    compress                 = false
    concurrency              = 7.17
    credentials_secret       = "...my_credentials_secret..."
    database                 = "...my_database..."
    description              = "...my_description..."
    dynamic_value_field_name = false
    environment              = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 4.51
    id                          = "...my_id..."
    login_url                   = "...my_login_url..."
    max_payload_events          = 5.21
    max_payload_size_kb         = 43679.21
    oauth_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    oauth_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "block"
    org             = "...my_org..."
    password        = "...my_password..."
    pipeline        = "...my_pipeline..."
    pq_compress     = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 13.22
        http_status     = 311.91
        initial_backoff = 490886.9
        max_backoff     = 124349.45
      }
    ]
    safe_headers = [
      "..."
    ]
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 1.3
      initial_backoff = 389060.46
      max_backoff     = 111384.89
      timeout_retry   = false
    }
    timeout_sec          = 8604065791755370
    timestamp_precision  = "s"
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 243059.19
    type                 = "influxdb"
    url                  = "...my_url..."
    use_round_robin_dns  = true
    use_v2_api           = false
    username             = "...my_username..."
    value_field_name     = "...my_value_field_name..."
  }
  output_kafka = {
    ack                    = 0
    authentication_timeout = 2381120.18
    backoff_rate           = 9.13
    brokers = [
      "..."
    ]
    compression        = "snappy"
    connection_timeout = 2861276.42
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_event_count  = 4347.43
    flush_period_sec   = 1
    format             = "json"
    id                 = "...my_id..."
    initial_backoff    = 285287.53
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout      = 29427.03
      default_key_schema_id   = 4.17
      default_value_schema_id = 7.43
      disabled                = false
      max_retries             = 52.95
      request_timeout         = 46927.92
      schema_registry_url     = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = true
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.1"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = true
        servername          = "...my_servername..."
      }
    }
    max_back_off       = 103584.41
    max_record_size_kb = 3.9
    max_retries        = 83.18
    on_backpressure    = "queue"
    pipeline           = "...my_pipeline..."
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "always"
    pq_on_backpressure         = "block"
    pq_path                    = "...my_pq_path..."
    protobuf_library_id        = "...my_protobuf_library_id..."
    reauthentication_threshold = 1112954.59
    request_timeout            = 3471769.08
    sasl = {
      disabled  = true
      mechanism = "kerberos"
    }
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    topic = "...my_topic..."
    type  = "kafka"
  }
  output_kinesis = {
    as_ndjson                 = true
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    compression               = "gzip"
    concurrency               = 20.46
    description               = "...my_description..."
    duration_seconds          = 30789.14
    enable_assume_role        = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    flush_period_sec          = 3.75
    id                        = "...my_id..."
    max_record_size_kb        = 4137.92
    on_backpressure           = "queue"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "always"
    pq_on_backpressure  = "block"
    pq_path             = "...my_pq_path..."
    region              = "...my_region..."
    reject_unauthorized = true
    reuse_connections   = true
    signature_version   = "v4"
    stream_name         = "...my_stream_name..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type            = "kinesis"
    use_list_shards = true
  }
  output_loki = {
    auth_type          = "credentialsSecret"
    compress           = true
    concurrency        = 25.43
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 3.83
    id                          = "...my_id..."
    labels = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    max_payload_events  = 1.25
    max_payload_size_kb = 4621.51
    message             = "...my_message..."
    message_format      = "json"
    on_backpressure     = "queue"
    password            = "...my_password..."
    pipeline            = "...my_pipeline..."
    pq_compress         = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 19.98
        http_status     = 298
        initial_backoff = 591584.04
        max_backoff     = 61340.13
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 2.98
      initial_backoff = 98153.12
      max_backoff     = 50787.77
      timeout_retry   = false
    }
    timeout_sec           = 8372597717400780
    token                 = "...my_token..."
    total_memory_limit_kb = 4
    type                  = "loki"
    url                   = "...my_url..."
    use_round_robin_dns   = true
    username              = "...my_username..."
  }
  output_minio = {
    add_id_to_stage_path      = true
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "manual"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    compress                  = "gzip"
    compression_level         = "best_speed"
    deadletter_enabled        = true
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    dest_path                 = "...my_dest_path..."
    empty_dir_cleanup_sec     = 31515.8
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    format                    = "parquet"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_concurrent_file_parts = 1.52
    max_file_idle_time_sec    = 81367.85
    max_file_open_time_sec    = 64617.65
    max_file_size_mb          = 709.95
    max_open_files            = 386.65
    max_retry_num             = 1.47
    object_acl                = "public-read-write"
    on_backpressure           = "drop"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V1"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 172832.5
    parquet_version           = "PARQUET_2_4"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    region                    = "...my_region..."
    reject_unauthorized       = true
    remove_empty_dirs         = true
    reuse_connections         = true
    server_side_encryption    = "AES256"
    should_log_invalid_rows   = false
    signature_version         = "v2"
    stage_path                = "...my_stage_path..."
    storage_class             = "STANDARD"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "minio"
    verify_permissions    = false
    write_high_water_mark = 41.09
  }
  output_msk = {
    ack                       = 1
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    authentication_timeout    = 1641453.71
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "manual"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    backoff_rate              = 8.43
    brokers = [
      "..."
    ]
    compression        = "lz4"
    connection_timeout = 1910015.54
    description        = "...my_description..."
    duration_seconds   = 27775.99
    enable_assume_role = false
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    flush_event_count  = 2774.07
    flush_period_sec   = 1.62
    format             = "json"
    id                 = "...my_id..."
    initial_backoff    = 241253.05
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout      = 3992.84
      default_key_schema_id   = 3.98
      default_value_schema_id = 8.19
      disabled                = true
      max_retries             = 73.69
      request_timeout         = 23465.67
      schema_registry_url     = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = true
        max_version         = "TLSv1.2"
        min_version         = "TLSv1.1"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = false
        servername          = "...my_servername..."
      }
    }
    max_back_off       = 108629.7
    max_record_size_kb = 8.67
    max_retries        = 8.4
    on_backpressure    = "drop"
    pipeline           = "...my_pipeline..."
    pq_compress        = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "always"
    pq_on_backpressure         = "drop"
    pq_path                    = "...my_pq_path..."
    protobuf_library_id        = "...my_protobuf_library_id..."
    reauthentication_threshold = 1077176.04
    region                     = "...my_region..."
    reject_unauthorized        = true
    request_timeout            = 1715384.27
    reuse_connections          = false
    signature_version          = "v4"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    topic = "...my_topic..."
    type  = "msk"
  }
  output_netflow = {
    description            = "...my_description..."
    dns_resolve_period_sec = 78639.6
    environment            = "...my_environment..."
    hosts = [
      {
        host = "...my_host..."
        port = 30084.75
      }
    ]
    id       = "...my_id..."
    pipeline = "...my_pipeline..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "netflow"
  }
  output_newrelic = {
    api_key     = "...my_api_key..."
    auth_type   = "secret"
    compress    = false
    concurrency = 22.63
    custom_url  = "...my_custom_url..."
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 4.8
    id                          = "...my_id..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 1.3
    max_payload_size_kb         = 136.53
    message_field               = "...my_message_field..."
    metadata = [
      {
        name  = "service"
        value = "...my_value..."
      }
    ]
    on_backpressure = "queue"
    pipeline        = "...my_pipeline..."
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    region                            = "Custom"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 15.56
        http_status     = 176.07
        initial_backoff = 399609.04
        max_backoff     = 150336.97
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 1.72
      initial_backoff = 577081.54
      max_backoff     = 148450.61
      timeout_retry   = true
    }
    timeout_sec           = 1657548890213634.5
    total_memory_limit_kb = 9.39
    type                  = "newrelic"
    use_round_robin_dns   = true
  }
  output_newrelic_events = {
    account_id  = "...my_account_id..."
    api_key     = "...my_api_key..."
    auth_type   = "manual"
    compress    = false
    concurrency = 2.8
    custom_url  = "...my_custom_url..."
    description = "...my_description..."
    environment = "...my_environment..."
    event_type  = "...my_event_type..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 0.63
    id                          = "...my_id..."
    max_payload_events          = 8.64
    max_payload_size_kb         = 661.76
    on_backpressure             = "drop"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    region                            = "EU"
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 18.84
        http_status     = 198.38
        initial_backoff = 536242.99
        max_backoff     = 34212.25
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 19.98
      initial_backoff = 99705.33
      max_backoff     = 142578.05
      timeout_retry   = false
    }
    timeout_sec         = 1878342669385474.5
    type                = "newrelic_events"
    use_round_robin_dns = true
  }
  output_open_telemetry = {
    auth_header_expr   = "...my_auth_header_expr..."
    auth_type          = "token"
    compress           = "gzip"
    concurrency        = 13.54
    connection_timeout = 0.01
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode    = "payload"
    flush_period_sec               = 5.39
    http_compress                  = "none"
    http_logs_endpoint_override    = "...my_http_logs_endpoint_override..."
    http_metrics_endpoint_override = "...my_http_metrics_endpoint_override..."
    http_traces_endpoint_override  = "...my_http_traces_endpoint_override..."
    id                             = "...my_id..."
    keep_alive                     = true
    keep_alive_time                = 6.7
    login_url                      = "...my_login_url..."
    max_payload_size_kb            = 3116.32
    metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    oauth_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    oauth_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "block"
    otlp_version    = "0.10.0"
    password        = "...my_password..."
    pipeline        = "...my_pipeline..."
    pq_compress     = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    protocol                          = "grpc"
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 13.74
        http_status     = 303.68
        initial_backoff = 349366.54
        max_backoff     = 40192.64
      }
    ]
    safe_headers = [
      "..."
    ]
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 15.63
      initial_backoff = 108886.99
      max_backoff     = 108211.16
      timeout_retry   = false
    }
    timeout_sec = 2211979713149453.5
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 221017.3
    type                 = "open_telemetry"
    use_round_robin_dns  = true
    username             = "...my_username..."
  }
  output_prometheus = {
    auth_header_expr   = "...my_auth_header_expr..."
    auth_type          = "none"
    concurrency        = 3.59
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 9.03
    id                          = "...my_id..."
    login_url                   = "...my_login_url..."
    max_payload_events          = 4.61
    max_payload_size_kb         = 9392.5
    metric_rename_expr          = "...my_metric_rename_expr..."
    metrics_flush_period_sec    = 9.03
    oauth_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    oauth_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "drop"
    password        = "...my_password..."
    pipeline        = "...my_pipeline..."
    pq_compress     = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 5.46
        http_status     = 151.65
        initial_backoff = 411730.38
        max_backoff     = 135846.3
      }
    ]
    safe_headers = [
      "..."
    ]
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_metadata     = true
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 9.91
      initial_backoff = 299842.71
      max_backoff     = 141362.36
      timeout_retry   = true
    }
    timeout_sec          = 761512813550484.8
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 12652.55
    type                 = "prometheus"
    url                  = "...my_url..."
    use_round_robin_dns  = true
    username             = "...my_username..."
  }
  output_ring = {
    compress        = "gzip"
    description     = "...my_description..."
    dest_path       = "...my_dest_path..."
    environment     = "...my_environment..."
    format          = "raw"
    id              = "...my_id..."
    max_data_size   = "...my_max_data_size..."
    max_data_time   = "...my_max_data_time..."
    on_backpressure = "block"
    partition_expr  = "...my_partition_expr..."
    pipeline        = "...my_pipeline..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "ring"
  }
  output_router = {
    description = "...my_description..."
    environment = "...my_environment..."
    id          = "...my_id..."
    pipeline    = "...my_pipeline..."
    rules = [
      {
        description = "...my_description..."
        filter      = "...my_filter..."
        final       = false
        output      = "...my_output..."
      }
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "router"
  }
  output_s3 = {
    add_id_to_stage_path      = false
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    compress                  = "none"
    compression_level         = "normal"
    deadletter_enabled        = false
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    dest_path                 = "...my_dest_path..."
    duration_seconds          = 13198.59
    empty_dir_cleanup_sec     = 31139.56
    enable_assume_role        = false
    enable_page_checksum      = false
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    format                    = "raw"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    kms_key_id                        = "...my_kms_key_id..."
    max_closing_files_to_backpressure = 870.15
    max_concurrent_file_parts         = 8.41
    max_file_idle_time_sec            = 45984.26
    max_file_open_time_sec            = 78997.61
    max_file_size_mb                  = 263.34
    max_open_files                    = 375.01
    max_retry_num                     = 8.66
    object_acl                        = "public-read-write"
    on_backpressure                   = "block"
    on_disk_full_backpressure         = "drop"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "...my_parquet_page_size..."
    parquet_row_group_length          = 19610353.91
    parquet_version                   = "PARQUET_2_4"
    partition_expr                    = "...my_partition_expr..."
    pipeline                          = "...my_pipeline..."
    region                            = "...my_region..."
    reject_unauthorized               = false
    remove_empty_dirs                 = false
    reuse_connections                 = true
    server_side_encryption            = "aws:kms"
    should_log_invalid_rows           = true
    signature_version                 = "v2"
    stage_path                        = "...my_stage_path..."
    storage_class                     = "INTELLIGENT_TIERING"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "s3"
    verify_permissions    = true
    write_high_water_mark = 2514.55
  }
  output_security_lake = {
    account_id                = "...my_account_id..."
    add_id_to_stage_path      = true
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    automatic_schema          = false
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    custom_source             = "...my_custom_source..."
    deadletter_enabled        = false
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    duration_seconds          = 36203.93
    empty_dir_cleanup_sec     = 13098.54
    enable_assume_role        = false
    enable_page_checksum      = false
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    kms_key_id                        = "...my_kms_key_id..."
    max_closing_files_to_backpressure = 148.09
    max_concurrent_file_parts         = 3.88
    max_file_idle_time_sec            = 63772.78
    max_file_open_time_sec            = 49744.3
    max_file_size_mb                  = 1021.37
    max_open_files                    = 558.11
    max_retry_num                     = 10.23
    object_acl                        = "bucket-owner-full-control"
    on_backpressure                   = "drop"
    on_disk_full_backpressure         = "drop"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "...my_parquet_page_size..."
    parquet_row_group_length          = 2808783.07
    parquet_schema                    = "...my_parquet_schema..."
    parquet_version                   = "PARQUET_2_4"
    pipeline                          = "...my_pipeline..."
    region                            = "...my_region..."
    reject_unauthorized               = true
    remove_empty_dirs                 = true
    reuse_connections                 = false
    server_side_encryption            = "aws:kms"
    should_log_invalid_rows           = false
    signature_version                 = "v2"
    stage_path                        = "...my_stage_path..."
    storage_class                     = "ONEZONE_IA"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "security_lake"
    verify_permissions    = true
    write_high_water_mark = 962.61
  }
  output_sentinel = {
    advanced_content_type      = "...my_advanced_content_type..."
    auth_type                  = "oauth"
    client_id                  = "...my_client_id..."
    compress                   = true
    concurrency                = 27.32
    custom_content_type        = "...my_custom_content_type..."
    custom_drop_when_null      = false
    custom_event_delimiter     = "...my_custom_event_delimiter..."
    custom_payload_expression  = "...my_custom_payload_expression..."
    custom_source_expression   = "...my_custom_source_expression..."
    dce_endpoint               = "...my_dce_endpoint..."
    dcr_id                     = "...my_dcr_id..."
    description                = "...my_description..."
    endpoint_url_configuration = "url"
    environment                = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 3.93
    format                      = "custom"
    format_event_code           = "...my_format_event_code..."
    format_payload_code         = "...my_format_payload_code..."
    id                          = "...my_id..."
    keep_alive                  = false
    login_url                   = "...my_login_url..."
    max_payload_events          = 6.43
    max_payload_size_kb         = 208.05
    on_backpressure             = "drop"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 3.78
        http_status     = 192.68
        initial_backoff = 569875.31
        max_backoff     = 63873.49
      }
    ]
    safe_headers = [
      "..."
    ]
    scope       = "...my_scope..."
    secret      = "...my_secret..."
    stream_name = "...my_stream_name..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 6.35
      initial_backoff = 55942.62
      max_backoff     = 151713.58
      timeout_retry   = false
    }
    timeout_sec           = 2862438398607763
    total_memory_limit_kb = 2.59
    type                  = "sentinel"
    url                   = "...my_url..."
    use_round_robin_dns   = true
  }
  output_service_now = {
    auth_token_name    = "...my_auth_token_name..."
    compress           = "deflate"
    concurrency        = 5.68
    connection_timeout = 1.35
    description        = "...my_description..."
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode    = "payload"
    flush_period_sec               = 2.66
    http_compress                  = "none"
    http_logs_endpoint_override    = "...my_http_logs_endpoint_override..."
    http_metrics_endpoint_override = "...my_http_metrics_endpoint_override..."
    http_traces_endpoint_override  = "...my_http_traces_endpoint_override..."
    id                             = "...my_id..."
    keep_alive                     = true
    keep_alive_time                = 9.46
    max_payload_size_kb            = 4911.02
    metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "drop"
    otlp_version    = "1.3.1"
    pipeline        = "...my_pipeline..."
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    protocol                          = "grpc"
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 7.43
        http_status     = 563.66
        initial_backoff = 173187.28
        max_backoff     = 122270.81
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 9.8
      initial_backoff = 4393.61
      max_backoff     = 35648.14
      timeout_retry   = false
    }
    timeout_sec = 1088602906594137.6
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
    }
    token_secret        = "...my_token_secret..."
    type                = "service_now"
    use_round_robin_dns = true
  }
  output_signalfx = {
    auth_type   = "manual"
    compress    = false
    concurrency = 5.75
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 1.66
    id                          = "...my_id..."
    max_payload_events          = 9.4
    max_payload_size_kb         = 9392.42
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    realm                             = "...my_realm..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 13.57
        http_status     = 586.58
        initial_backoff = 232593.45
        max_backoff     = 15175.48
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 10.19
      initial_backoff = 469024.48
      max_backoff     = 176316.39
      timeout_retry   = false
    }
    timeout_sec         = 3866173051202189.5
    token               = "...my_token..."
    type                = "signalfx"
    use_round_robin_dns = false
  }
  output_snmp = {
    description            = "...my_description..."
    dns_resolve_period_sec = 54402.89
    environment            = "...my_environment..."
    hosts = [
      {
        host = "...my_host..."
        port = 18749.9
      }
    ]
    id       = "...my_id..."
    pipeline = "...my_pipeline..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "snmp"
  }
  output_sns = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    description               = "...my_description..."
    duration_seconds          = 32972.43
    enable_assume_role        = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    id                        = "...my_id..."
    max_retries               = 5.95
    message_group_id          = "...my_message_group_id..."
    on_backpressure           = "drop"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "always"
    pq_on_backpressure  = "block"
    pq_path             = "...my_pq_path..."
    region              = "...my_region..."
    reject_unauthorized = true
    reuse_connections   = true
    signature_version   = "v4"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    topic_arn = "...my_topic_arn..."
    type      = "sns"
  }
  output_splunk = {
    auth_token               = "...my_auth_token..."
    auth_type                = "manual"
    compress                 = "auto"
    connection_timeout       = 3.78
    description              = "...my_description..."
    enable_ack               = false
    enable_multi_metrics     = true
    environment              = "...my_environment..."
    host                     = "...my_host..."
    id                       = "...my_id..."
    log_failed_requests      = false
    max_failed_health_checks = 6.17
    max_s2_sversion          = "v3"
    nested_fields            = "json"
    on_backpressure          = "queue"
    pipeline                 = "...my_pipeline..."
    port                     = 15386.29
    pq_compress              = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "error"
    pq_on_backpressure = "drop"
    pq_path            = "...my_pq_path..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret           = "...my_text_secret..."
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    type          = "splunk"
    write_timeout = 8.31
  }
  output_splunk_hec = {
    auth_type              = "secret"
    compress               = true
    concurrency            = 24.97
    description            = "...my_description..."
    dns_resolve_period_sec = 47008.18
    enable_multi_metrics   = false
    environment            = "...my_environment..."
    exclude_self           = false
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "none"
    flush_period_sec              = 1.13
    id                            = "...my_id..."
    load_balance_stats_period_sec = 17.43
    load_balanced                 = false
    max_payload_events            = 4.8
    max_payload_size_kb           = 911757.13
    next_queue                    = "...my_next_queue..."
    on_backpressure               = "block"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 5.37
        http_status     = 544.87
        initial_backoff = 152474.46
        max_backoff     = 37480.26
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tcp_routing = "...my_tcp_routing..."
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 12.64
      initial_backoff = 425987.93
      max_backoff     = 119739.45
      timeout_retry   = false
    }
    timeout_sec = 8923863287730264
    token       = "...my_token..."
    type        = "splunk_hec"
    url         = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 4.18
      }
    ]
    use_round_robin_dns = true
  }
  output_splunk_lb = {
    auth_token             = "...my_auth_token..."
    auth_type              = "secret"
    compress               = "auto"
    connection_timeout     = 9.43
    description            = "...my_description..."
    dns_resolve_period_sec = 54004.74
    enable_ack             = false
    enable_multi_metrics   = false
    environment            = "...my_environment..."
    exclude_self           = true
    hosts = [
      {
        host       = "...my_host..."
        port       = 45464.92
        servername = "...my_servername..."
        tls        = "off"
        weight     = 4.85
      }
    ]
    id                = "...my_id..."
    indexer_discovery = true
    indexer_discovery_configs = {
      auth_token = "...my_auth_token..."
      auth_tokens = [
        {
          auth_type = "secret"
        }
      ]
      auth_type            = "secret"
      master_uri           = "...my_master_uri..."
      refresh_interval_sec = 8238.57
      reject_unauthorized  = false
      site                 = "...my_site..."
      text_secret          = "...my_text_secret..."
    }
    load_balance_stats_period_sec = 19.58
    log_failed_requests           = false
    max_concurrent_senders        = 0.18
    max_failed_health_checks      = 0.42
    max_s2_sversion               = "v4"
    nested_fields                 = "json"
    on_backpressure               = "block"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                = "...my_pq_max_file_size..."
    pq_max_size                     = "...my_pq_max_size..."
    pq_mode                         = "backpressure"
    pq_on_backpressure              = "drop"
    pq_path                         = "...my_pq_path..."
    sender_unhealthy_time_allowance = 7268.46
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret           = "...my_text_secret..."
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    type          = "splunk_lb"
    write_timeout = 9.39
  }
  output_sqs = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    create_queue              = false
    description               = "...my_description..."
    duration_seconds          = 18284.05
    enable_assume_role        = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    flush_period_sec          = 6.7
    id                        = "...my_id..."
    max_in_progress           = 61.82
    max_queue_size            = 6.35
    max_record_size_kb        = 228.5
    message_group_id          = "...my_message_group_id..."
    on_backpressure           = "drop"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "always"
    pq_on_backpressure  = "block"
    pq_path             = "...my_pq_path..."
    queue_name          = "...my_queue_name..."
    queue_type          = "standard"
    region              = "...my_region..."
    reject_unauthorized = true
    reuse_connections   = false
    signature_version   = "v4"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "sqs"
  }
  output_statsd = {
    connection_timeout     = 1.29
    description            = "...my_description..."
    dns_resolve_period_sec = 74600.55
    environment            = "...my_environment..."
    flush_period_sec       = 9.96
    host                   = "...my_host..."
    id                     = "...my_id..."
    mtu                    = 44909.26
    on_backpressure        = "block"
    pipeline               = "...my_pipeline..."
    port                   = 31524.34
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "always"
    pq_on_backpressure = "block"
    pq_path            = "...my_pq_path..."
    protocol           = "tcp"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    type                  = "statsd"
    write_timeout         = 0.33
  }
  output_statsd_ext = {
    connection_timeout     = 5.54
    description            = "...my_description..."
    dns_resolve_period_sec = 80279.36
    environment            = "...my_environment..."
    flush_period_sec       = 4.45
    host                   = "...my_host..."
    id                     = "...my_id..."
    mtu                    = 34944.89
    on_backpressure        = "drop"
    pipeline               = "...my_pipeline..."
    port                   = 4559.78
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "...my_pq_path..."
    protocol           = "udp"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    type                  = "statsd_ext"
    write_timeout         = 0.09
  }
  output_sumo_logic = {
    compress        = true
    concurrency     = 8.16
    custom_category = "...my_custom_category..."
    custom_source   = "...my_custom_source..."
    description     = "...my_description..."
    environment     = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1.17
    format                      = "json"
    id                          = "...my_id..."
    max_payload_events          = 7.5
    max_payload_size_kb         = 951.56
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 10.83
        http_status     = 216.15
        initial_backoff = 258770.1
        max_backoff     = 77265.54
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    timeout_retry_settings = {
      backoff_rate    = 15.65
      initial_backoff = 147734.27
      max_backoff     = 18045.77
      timeout_retry   = true
    }
    timeout_sec           = 5935725774249079
    total_memory_limit_kb = 9.17
    type                  = "sumo_logic"
    url                   = "...my_url..."
    use_round_robin_dns   = true
  }
  output_syslog = {
    app_name            = "...my_app_name..."
    connection_timeout  = 8.39
    description         = "...my_description..."
    environment         = "...my_environment..."
    facility            = 8
    host                = "...my_host..."
    id                  = "...my_id..."
    load_balanced       = false
    log_failed_requests = false
    max_record_size     = 28086.84
    message_format      = "rfc3164"
    octet_count_framing = false
    on_backpressure     = "queue"
    pipeline            = "...my_pipeline..."
    port                = 23921.26
    pq_compress         = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "backpressure"
    pq_on_backpressure = "block"
    pq_path            = "...my_pq_path..."
    protocol           = "udp"
    severity           = 2
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    timestamp_format      = "iso8601"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    type                       = "syslog"
    udp_dns_resolve_period_sec = 66604.26
    write_timeout              = 2.78
  }
  output_tcpjson = {
    auth_token             = "...my_auth_token..."
    auth_type              = "manual"
    compression            = "none"
    connection_timeout     = 6.56
    description            = "...my_description..."
    dns_resolve_period_sec = 36485.94
    environment            = "...my_environment..."
    exclude_self           = false
    host                   = "...my_host..."
    hosts = [
      {
        host       = "...my_host..."
        port       = 8100.96
        servername = "...my_servername..."
        tls        = "off"
        weight     = 7.27
      }
    ]
    id                            = "...my_id..."
    load_balance_stats_period_sec = 18.26
    load_balanced                 = true
    log_failed_requests           = false
    max_concurrent_senders        = 9.76
    on_backpressure               = "drop"
    pipeline                      = "...my_pipeline..."
    port                          = 26330.99
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "always"
    pq_on_backpressure = "drop"
    pq_path            = "...my_pq_path..."
    send_header        = false
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret           = "...my_text_secret..."
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    token_ttl_minutes = 55.58
    type              = "tcpjson"
    write_timeout     = 5.56
  }
  output_wavefront = {
    auth_type   = "manual"
    compress    = false
    concurrency = 3.65
    description = "...my_description..."
    domain      = "...my_domain..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 5.92
    id                          = "...my_id..."
    max_payload_events          = 7.94
    max_payload_size_kb         = 2402.33
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2.05
        http_status     = 208.29
        initial_backoff = 376170.95
        max_backoff     = 109756.54
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 10.97
      initial_backoff = 439304.18
      max_backoff     = 136073.95
      timeout_retry   = false
    }
    timeout_sec         = 1604764437605072.8
    token               = "...my_token..."
    type                = "wavefront"
    use_round_robin_dns = false
  }
  output_webhook = {
    advanced_content_type     = "...my_advanced_content_type..."
    auth_header_expr          = "...my_auth_header_expr..."
    auth_type                 = "textSecret"
    compress                  = true
    concurrency               = 30.31
    credentials_secret        = "...my_credentials_secret..."
    custom_content_type       = "...my_custom_content_type..."
    custom_drop_when_null     = true
    custom_event_delimiter    = "...my_custom_event_delimiter..."
    custom_payload_expression = "...my_custom_payload_expression..."
    custom_source_expression  = "...my_custom_source_expression..."
    description               = "...my_description..."
    dns_resolve_period_sec    = 36601.64
    environment               = "...my_environment..."
    exclude_self              = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "none"
    flush_period_sec              = 5.05
    format                        = "json_array"
    format_event_code             = "...my_format_event_code..."
    format_payload_code           = "...my_format_payload_code..."
    id                            = "...my_id..."
    keep_alive                    = true
    load_balance_stats_period_sec = 12.05
    load_balanced                 = true
    login_url                     = "...my_login_url..."
    max_payload_events            = 4.79
    max_payload_size_kb           = 261448.74
    method                        = "PUT"
    oauth_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    oauth_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "block"
    password        = "...my_password..."
    pipeline        = "...my_pipeline..."
    pq_compress     = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2.17
        http_status     = 298.8
        initial_backoff = 330493.67
        max_backoff     = 125603.09
      }
    ]
    safe_headers = [
      "..."
    ]
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 18.85
      initial_backoff = 131843.06
      max_backoff     = 67206.95
      timeout_retry   = true
    }
    timeout_sec = 8339898100238354
    tls = {
      ca_path          = "...my_ca_path..."
      cert_path        = "...my_cert_path..."
      certificate_name = "...my_certificate_name..."
      disabled         = false
      max_version      = "TLSv1.2"
      min_version      = "TLSv1.2"
      passphrase       = "...my_passphrase..."
      priv_key_path    = "...my_priv_key_path..."
      servername       = "...my_servername..."
    }
    token                 = "...my_token..."
    token_attribute_name  = "...my_token_attribute_name..."
    token_timeout_secs    = 148758.23
    total_memory_limit_kb = 4.78
    type                  = "webhook"
    url                   = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 4.83
      }
    ]
    use_round_robin_dns = false
    username            = "...my_username..."
  }
  output_xsiam = {
    auth_type              = "token"
    concurrency            = 1.29
    description            = "...my_description..."
    dns_resolve_period_sec = 71005.67
    environment            = "...my_environment..."
    exclude_self           = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 7.7
    id                            = "...my_id..."
    load_balance_stats_period_sec = 14.39
    load_balanced                 = false
    max_payload_events            = 7.22
    max_payload_size_kb           = 385.16
    on_backpressure               = "queue"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 14.67
        http_status     = 166.25
        initial_backoff = 490950.74
        max_backoff     = 114257.9
      }
    ]
    safe_headers = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret               = "...my_text_secret..."
    throttle_rate_req_per_sec = 1802
    timeout_retry_settings = {
      backoff_rate    = 8.8
      initial_backoff = 487877.87
      max_backoff     = 150636.26
      timeout_retry   = true
    }
    timeout_sec           = 7356276809707653
    token                 = "...my_token..."
    total_memory_limit_kb = 9.15
    type                  = "xsiam"
    url                   = "...my_url..."
    urls = [
      {
        url    = "{ \"see\": \"documentation\" }"
        weight = 8.62
      }
    ]
    use_round_robin_dns = true
  }
  pack = "...my_pack..."
}