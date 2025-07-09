resource "criblio_destination" "my_destination" {
  group_id = "...my_group_id..."
  id       = "...my_id..."
  output_azure_blob = {
    add_id_to_stage_path = false
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
    compression_level       = "normal"
    connection_string       = "...my_connection_string..."
    container_name          = "...my_container_name..."
    create_container        = true
    deadletter_enabled      = false
    deadletter_path         = "...my_deadletter_path..."
    description             = "...my_description..."
    dest_path               = "...my_dest_path..."
    empty_dir_cleanup_sec   = 21121.14
    enable_page_checksum    = false
    enable_statistics       = true
    enable_write_page_index = true
    endpoint_suffix         = "...my_endpoint_suffix..."
    environment             = "...my_environment..."
    file_name_suffix        = "...my_file_name_suffix..."
    format                  = "json"
    header_line             = "...my_header_line..."
    id                      = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_concurrent_file_parts = 6
    max_file_idle_time_sec    = 319.67
    max_file_open_time_sec    = 30.58
    max_file_size_mb          = 714.31
    max_open_files            = 1383.44
    max_retry_num             = 10.66
    on_backpressure           = "drop"
    on_disk_full_backpressure = "drop"
    parquet_data_page_version = "DATA_PAGE_V1"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 17284892.23
    parquet_version           = "PARQUET_1_0"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    remove_empty_dirs         = true
    should_log_invalid_rows   = true
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
    write_high_water_mark = 468.94
  }
  output_azure_data_explorer = {
    add_id_to_stage_path = false
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
    concurrency        = 3.92
    database           = "...my_database..."
    deadletter_enabled = false
    description        = "...my_description..."
    environment        = "...my_environment..."
    extent_tags = [
      {
        prefix = "ingestBy"
        value  = "...my_value..."
      }
    ]
    file_name_suffix  = "...my_file_name_suffix..."
    flush_immediately = true
    flush_period_sec  = 2.57
    format            = "parquet"
    id                = "...my_id..."
    ingest_if_not_exists = [
      {
        value = "...my_value..."
      }
    ]
    ingest_mode               = "streaming"
    ingest_url                = "...my_ingest_url..."
    is_mapping_obj            = true
    keep_alive                = true
    mapping_ref               = "...my_mapping_ref..."
    max_concurrent_file_parts = 9.29
    max_file_idle_time_sec    = 988.97
    max_file_open_time_sec    = 350.72
    max_file_size_mb          = 512.37
    max_open_files            = 1540.24
    max_payload_events        = 4.87
    max_payload_size_kb       = 3428.8
    oauth_endpoint            = "https://login.microsoftonline.com"
    oauth_type                = "certificate"
    on_backpressure           = "drop"
    on_disk_full_backpressure = "block"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    remove_empty_dirs                 = true
    report_level                      = "failuresAndSuccesses"
    report_method                     = "queue"
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 4.86
        http_status     = 181.57
        initial_backoff = 252760.12
        max_backoff     = 123282.61
      }
    ]
    retain_blob_on_success = false
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
      backoff_rate    = 10.54
      initial_backoff = 253703.31
      max_backoff     = 26520.97
      timeout_retry   = false
    }
    timeout_sec                = 8241692246205138
    type                       = "azure_data_explorer"
    use_round_robin_dns        = false
    validate_database_settings = false
  }
  output_azure_eventhub = {
    ack                    = 8
    authentication_timeout = 1447914.42
    backoff_rate           = 2.1
    brokers = [
      "..."
    ]
    connection_timeout = 1229939.68
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_event_count  = 562.24
    flush_period_sec   = 2.06
    format             = "raw"
    id                 = "...my_id..."
    initial_backoff    = 136215.67
    max_back_off       = 42762.07
    max_record_size_kb = 10.23
    max_retries        = 72.04
    on_backpressure    = "drop"
    pipeline           = "...my_pipeline..."
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "error"
    pq_on_backpressure         = "block"
    pq_path                    = "...my_pq_path..."
    reauthentication_threshold = 1570187.28
    request_timeout            = 804154.32
    sasl = {
      disabled  = true
      mechanism = "plain"
    }
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    tls = {
      disabled            = true
      reject_unauthorized = true
    }
    topic = "...my_topic..."
    type  = "azure_eventhub"
  }
  output_azure_logs = {
    api_url     = "...my_api_url..."
    auth_type   = "secret"
    compress    = true
    concurrency = 28.04
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 5.15
    id                          = "...my_id..."
    keypair_secret              = "...my_keypair_secret..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 9.67
    max_payload_size_kb         = 8299.52
    on_backpressure             = "drop"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    resource_id                       = "...my_resource_id..."
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 8.21
        http_status     = 288.12
        initial_backoff = 109132.84
        max_backoff     = 34037.32
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
      backoff_rate    = 1.85
      initial_backoff = 465514.23
      max_backoff     = 83472.79
      timeout_retry   = false
    }
    timeout_sec         = 1749220463386416.2
    type                = "azure_logs"
    use_round_robin_dns = false
    workspace_id        = "...my_workspace_id..."
    workspace_key       = "...my_workspace_key..."
  }
  output_click_house = {
    async_inserts    = true
    auth_header_expr = "...my_auth_header_expr..."
    auth_type        = "basic"
    column_mappings = [
      {
        column_name             = "...my_column_name..."
        column_type             = "...my_column_type..."
        column_value_expression = "...my_column_value_expression..."
      }
    ]
    compress                   = true
    concurrency                = 15.53
    credentials_secret         = "...my_credentials_secret..."
    database                   = "...my_database..."
    describe_table             = "...my_describe_table..."
    description                = "...my_description..."
    dump_format_errors_to_disk = true
    environment                = "...my_environment..."
    exclude_mapping_fields = [
    ]
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 7.19
    format                      = "json-compact-each-row-with-names"
    id                          = "...my_id..."
    login_url                   = "...my_login_url..."
    mapping_type                = "automatic"
    max_payload_events          = 0.05
    max_payload_size_kb         = 10188.82
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
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 19.23
        http_status     = 260.57
        initial_backoff = 37931.99
        max_backoff     = 31355.47
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
      backoff_rate    = 11.99
      initial_backoff = 270487.69
      max_backoff     = 44857.19
      timeout_retry   = true
    }
    timeout_sec = 857692718093333.8
    tls = {
      ca_path          = "...my_ca_path..."
      cert_path        = "...my_cert_path..."
      certificate_name = "...my_certificate_name..."
      disabled         = true
      max_version      = "TLSv1.3"
      min_version      = "TLSv1.2"
      passphrase       = "...my_passphrase..."
      priv_key_path    = "...my_priv_key_path..."
      servername       = "...my_servername..."
    }
    token                  = "...my_token..."
    token_attribute_name   = "...my_token_attribute_name..."
    token_timeout_secs     = 297875.2
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
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    description               = "...my_description..."
    duration_seconds          = 14032.13
    enable_assume_role        = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    flush_period_sec          = 7.05
    id                        = "...my_id..."
    log_group_name            = "...my_log_group_name..."
    log_stream_name           = "...my_log_stream_name..."
    max_queue_size            = 30.9
    max_record_size_kb        = 2483.83
    on_backpressure           = "block"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "backpressure"
    pq_on_backpressure  = "block"
    pq_path             = "...my_pq_path..."
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = false
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "cloudwatch"
  }
  output_confluent_cloud = {
    ack                    = 1
    authentication_timeout = 2796245.38
    backoff_rate           = 17.03
    brokers = [
      "..."
    ]
    compression        = "lz4"
    connection_timeout = 1694680.11
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_event_count  = 8492.17
    flush_period_sec   = 4.97
    format             = "raw"
    id                 = "...my_id..."
    initial_backoff    = 379354.8
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout      = 3442.55
      default_key_schema_id   = 1.09
      default_value_schema_id = 3.27
      disabled                = true
      max_retries             = 15.86
      request_timeout         = 49378.51
      schema_registry_url     = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = false
        max_version         = "TLSv1.1"
        min_version         = "TLSv1"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = true
        servername          = "...my_servername..."
      }
    }
    max_back_off       = 61861.1
    max_record_size_kb = 8.17
    max_retries        = 5.85
    on_backpressure    = "drop"
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
    reauthentication_threshold = 1433127.23
    request_timeout            = 992593.26
    sasl = {
      disabled  = false
      mechanism = "scram-sha-256"
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
      disabled            = true
      max_version         = "TLSv1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    topic = "...my_topic..."
    type  = "confluent_cloud"
  }
  output_cribl_http = {
    compression            = "gzip"
    concurrency            = 11.56
    description            = "...my_description..."
    dns_resolve_period_sec = 5488.51
    environment            = "...my_environment..."
    exclude_fields = [
      "..."
    ]
    exclude_self = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "none"
    flush_period_sec              = 2.35
    id                            = "...my_id..."
    load_balance_stats_period_sec = 12.74
    load_balanced                 = false
    max_payload_events            = 6.83
    max_payload_size_kb           = 2827.75
    on_backpressure               = "block"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 18.82
        http_status     = 531.39
        initial_backoff = 117330.79
        max_backoff     = 52596.23
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
      backoff_rate    = 19.69
      initial_backoff = 173747.34
      max_backoff     = 163675.23
      timeout_retry   = true
    }
    timeout_sec = 8157421429180492
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    token_ttl_minutes = 38.63
    type              = "cribl_http"
    url               = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 0.48
      }
    ]
    use_round_robin_dns = false
  }
  output_cribl_lake = {
    add_id_to_stage_path              = true
    assume_role_arn                   = "...my_assume_role_arn..."
    assume_role_external_id           = "...my_assume_role_external_id..."
    aws_authentication_method         = "auto"
    aws_secret_key                    = "...my_aws_secret_key..."
    base_file_name                    = "...my_base_file_name..."
    bucket                            = "...my_bucket..."
    deadletter_enabled                = true
    deadletter_path                   = "...my_deadletter_path..."
    description                       = "...my_description..."
    dest_path                         = "...my_dest_path..."
    duration_seconds                  = 35105.34
    empty_dir_cleanup_sec             = 67648.06
    enable_assume_role                = false
    endpoint                          = "...my_endpoint..."
    environment                       = "...my_environment..."
    file_name_suffix                  = "...my_file_name_suffix..."
    format                            = "json"
    header_line                       = "...my_header_line..."
    id                                = "...my_id..."
    kms_key_id                        = "...my_kms_key_id..."
    max_closing_files_to_backpressure = 434.23
    max_concurrent_file_parts         = 5.6
    max_file_idle_time_sec            = 80824.88
    max_file_open_time_sec            = 9898.55
    max_file_size_mb                  = 57
    max_open_files                    = 1917.9
    max_retry_num                     = 4.32
    object_acl                        = "aws-exec-read"
    on_backpressure                   = "block"
    on_disk_full_backpressure         = "drop"
    pipeline                          = "...my_pipeline..."
    region                            = "...my_region..."
    reject_unauthorized               = true
    remove_empty_dirs                 = true
    reuse_connections                 = true
    server_side_encryption            = "aws:kms"
    signature_version                 = "v4"
    stage_path                        = "...my_stage_path..."
    storage_class                     = "STANDARD"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "cribl_lake"
    verify_permissions    = true
    write_high_water_mark = 1881.2
  }
  output_cribl_tcp = {
    compression            = "none"
    connection_timeout     = 8.15
    description            = "...my_description..."
    dns_resolve_period_sec = 42922.72
    environment            = "...my_environment..."
    exclude_fields = [
      "..."
    ]
    exclude_self = true
    host         = "...my_host..."
    hosts = [
      {
        host       = "...my_host..."
        port       = 48249.06
        servername = "...my_servername..."
        tls        = "off"
        weight     = 3.03
      }
    ]
    id                            = "...my_id..."
    load_balance_stats_period_sec = 12.07
    load_balanced                 = false
    log_failed_requests           = true
    max_concurrent_senders        = 8.1
    on_backpressure               = "drop"
    pipeline                      = "...my_pipeline..."
    port                          = 56940.97
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "always"
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
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    token_ttl_minutes = 33.86
    type              = "cribl_tcp"
    write_timeout     = 5.99
  }
  output_crowdstrike_next_gen_siem = {
    auth_type   = "manual"
    compress    = true
    concurrency = 6.28
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 9.03
    format                      = "JSON"
    id                          = "...my_id..."
    max_payload_events          = 3.6
    max_payload_size_kb         = 6393.69
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
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 2.5
        http_status     = 127.7
        initial_backoff = 468169.59
        max_backoff     = 73357.1
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
      backoff_rate    = 17.69
      initial_backoff = 381223.65
      max_backoff     = 67882.01
      timeout_retry   = false
    }
    timeout_sec         = 6739864097378947
    token               = "...my_token..."
    type                = "crowdstrike_next_gen_siem"
    url                 = "...my_url..."
    use_round_robin_dns = true
  }
  output_datadog = {
    allow_api_key_from_events = true
    api_key                   = "...my_api_key..."
    auth_type                 = "manual"
    batch_by_tags             = true
    compress                  = false
    concurrency               = 22.17
    content_type              = "json"
    custom_url                = "...my_custom_url..."
    description               = "...my_description..."
    environment               = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 6.96
    host                        = "...my_host..."
    id                          = "...my_id..."
    max_payload_events          = 4.96
    max_payload_size_kb         = 4725.05
    message                     = "...my_message..."
    on_backpressure             = "drop"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 3.46
        http_status     = 392.17
        initial_backoff = 126652.97
        max_backoff     = 36651.45
      }
    ]
    safe_headers = [
      "..."
    ]
    send_counters_as_count = false
    service                = "...my_service..."
    severity               = "error"
    site                   = "us"
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
      backoff_rate    = 11.84
      initial_backoff = 172125.36
      max_backoff     = 13416.36
      timeout_retry   = true
    }
    timeout_sec           = 8902767411011333
    total_memory_limit_kb = 1.59
    type                  = "datadog"
    use_round_robin_dns   = true
  }
  output_dataset = {
    api_key          = "...my_api_key..."
    auth_type        = "manual"
    compress         = false
    concurrency      = 14.32
    custom_url       = "...my_custom_url..."
    default_severity = "info"
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
    flush_period_sec            = 4.17
    id                          = "...my_id..."
    max_payload_events          = 3.04
    max_payload_size_kb         = 1716.8
    message_field               = "...my_message_field..."
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "none"
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
        backoff_rate    = 2.98
        http_status     = 338.21
        initial_backoff = 98153.12
        max_backoff     = 50787.77
      }
    ]
    safe_headers = [
      "..."
    ]
    server_host_field = "...my_server_host_field..."
    site              = "custom"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    timeout_retry_settings = {
      backoff_rate    = 2.29
      initial_backoff = 417618.62
      max_backoff     = 53881.05
      timeout_retry   = true
    }
    timeout_sec           = 8209735142372086
    timestamp_field       = "...my_timestamp_field..."
    total_memory_limit_kb = 4.5
    type                  = "dataset"
    use_round_robin_dns   = true
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
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    compress                  = "gzip"
    compression_level         = "best_compression"
    deadletter_enabled        = true
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    dest_path                 = "...my_dest_path..."
    duration_seconds          = 43113.72
    empty_dir_cleanup_sec     = 82151.17
    enable_assume_role        = true
    enable_page_checksum      = true
    enable_statistics         = true
    enable_write_page_index   = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    format                    = "json"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    kms_key_id                        = "...my_kms_key_id..."
    max_closing_files_to_backpressure = 3000.28
    max_concurrent_file_parts         = 9.42
    max_file_idle_time_sec            = 56781.01
    max_file_open_time_sec            = 28709.8
    max_file_size_mb                  = 134.87
    max_open_files                    = 1956.83
    max_retry_num                     = 1.71
    object_acl                        = "bucket-owner-full-control"
    on_backpressure                   = "block"
    on_disk_full_backpressure         = "block"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "...my_parquet_page_size..."
    parquet_row_group_length          = 32406920.2
    parquet_version                   = "PARQUET_2_6"
    partitioning_fields = [
      "..."
    ]
    pipeline                = "...my_pipeline..."
    region                  = "...my_region..."
    reject_unauthorized     = false
    remove_empty_dirs       = true
    reuse_connections       = true
    server_side_encryption  = "aws:kms"
    should_log_invalid_rows = true
    signature_version       = "v4"
    stage_path              = "...my_stage_path..."
    storage_class           = "GLACIER_IR"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "dl_s3"
    verify_permissions    = true
    write_high_water_mark = 1640.35
  }
  output_dynatrace_http = {
    active_gate_domain = "...my_active_gate_domain..."
    auth_type          = "token"
    compress           = false
    concurrency        = 4.7
    description        = "...my_description..."
    endpoint           = "manual"
    environment        = "...my_environment..."
    environment_id     = "...my_environment_id..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 2.73
    format                      = "plaintext"
    id                          = "...my_id..."
    keep_alive                  = false
    max_payload_events          = 5590.91
    max_payload_size_kb         = 1064.93
    method                      = "PUT"
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
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 12.84
        http_status     = 507.81
        initial_backoff = 443637.25
        max_backoff     = 159340.75
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
      backoff_rate    = 18.92
      initial_backoff = 538199.44
      max_backoff     = 148849.5
      timeout_retry   = true
    }
    timeout_sec           = 8159783654596632
    token                 = "...my_token..."
    total_memory_limit_kb = 2.15
    type                  = "dynatrace_http"
    url                   = "...my_url..."
    use_round_robin_dns   = false
  }
  output_dynatrace_otlp = {
    auth_token_name    = "...my_auth_token_name..."
    compress           = "gzip"
    concurrency        = 11.94
    connection_timeout = 7.73
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
    failed_request_logging_mode    = "payloadAndHeaders"
    flush_period_sec               = 5.07
    http_compress                  = "none"
    http_logs_endpoint_override    = "...my_http_logs_endpoint_override..."
    http_metrics_endpoint_override = "...my_http_metrics_endpoint_override..."
    http_traces_endpoint_override  = "...my_http_traces_endpoint_override..."
    id                             = "...my_id..."
    keep_alive                     = false
    keep_alive_time                = 6.09
    max_payload_size_kb            = 2517.53
    metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "drop"
    otlp_version    = "1.3.1"
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
    protocol                          = "http"
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 19.86
        http_status     = 592.64
        initial_backoff = 434894.76
        max_backoff     = 109352.67
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
      backoff_rate    = 8.1
      initial_backoff = 169047.42
      max_backoff     = 164491.44
      timeout_retry   = false
    }
    timeout_sec         = 8431236050803757
    token_secret        = "...my_token_secret..."
    type                = "dynatrace_otlp"
    use_round_robin_dns = false
  }
  output_elastic = {
    auth = {
      auth_type = "manualAPIKey"
      disabled  = false
    }
    compress               = false
    concurrency            = 10.25
    description            = "...my_description..."
    dns_resolve_period_sec = 43899.02
    doc_type               = "...my_doc_type..."
    elastic_pipeline       = "...my_elastic_pipeline..."
    elastic_version        = "6"
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
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 8.2
    id                            = "...my_id..."
    include_doc_id                = false
    index                         = "...my_index..."
    load_balance_stats_period_sec = 19.24
    load_balanced                 = true
    max_payload_events            = 6.89
    max_payload_size_kb           = 3136.35
    on_backpressure               = "block"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "backpressure"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 11.31
        http_status     = 354.9
        initial_backoff = 1863.4
        max_backoff     = 45979.83
      }
    ]
    retry_partial_errors = true
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
      backoff_rate    = 13.57
      initial_backoff = 116512.92
      max_backoff     = 133601.54
      timeout_retry   = true
    }
    timeout_sec = 6650174910156927
    type        = "elastic"
    url         = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 1.87
      }
    ]
    use_round_robin_dns = false
    write_action        = "index"
  }
  output_elastic_cloud = {
    auth = {
      auth_type = "manualAPIKey"
      disabled  = false
    }
    compress         = false
    concurrency      = 17.26
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
    flush_period_sec            = 0.99
    id                          = "...my_id..."
    include_doc_id              = true
    index                       = "...my_index..."
    max_payload_events          = 0.72
    max_payload_size_kb         = 87952.81
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 10.64
        http_status     = 577.87
        initial_backoff = 562503.66
        max_backoff     = 39502.38
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
      backoff_rate    = 18.8
      initial_backoff = 304893.17
      max_backoff     = 179310.37
      timeout_retry   = false
    }
    timeout_sec = 7586109465427404
    type        = "elastic_cloud"
    url         = "...my_url..."
  }
  output_exabeam = {
    add_id_to_stage_path      = true
    aws_api_key               = "...my_aws_api_key..."
    aws_secret_key            = "...my_aws_secret_key..."
    bucket                    = "...my_bucket..."
    collector_instance_id     = "...my_collector_instance_id..."
    deadletter_enabled        = false
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    empty_dir_cleanup_sec     = 74366.88
    encoded_configuration     = "...my_encoded_configuration..."
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    id                        = "...my_id..."
    max_file_idle_time_sec    = 1120.36
    max_file_open_time_sec    = 394.37
    max_file_size_mb          = 746.3
    max_open_files            = 932.4
    max_retry_num             = 5.25
    object_acl                = "bucket-owner-read"
    on_backpressure           = "drop"
    on_disk_full_backpressure = "block"
    pipeline                  = "...my_pipeline..."
    region                    = "...my_region..."
    reject_unauthorized       = false
    remove_empty_dirs         = true
    reuse_connections         = true
    signature_version         = "v4"
    site_id                   = "...my_site_id..."
    site_name                 = "...my_site_name..."
    stage_path                = "...my_stage_path..."
    storage_class             = "NEARLINE"
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
    compress                = "none"
    compression_level       = "best_compression"
    deadletter_enabled      = true
    deadletter_path         = "...my_deadletter_path..."
    description             = "...my_description..."
    dest_path               = "...my_dest_path..."
    empty_dir_cleanup_sec   = 13326.35
    enable_page_checksum    = true
    enable_statistics       = false
    enable_write_page_index = false
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
    max_file_idle_time_sec    = 80.78
    max_file_open_time_sec    = 639.59
    max_file_size_mb          = 544.3
    max_open_files            = 1295.51
    max_retry_num             = 4.57
    on_backpressure           = "drop"
    on_disk_full_backpressure = "drop"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 55900201.28
    parquet_version           = "PARQUET_2_6"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    remove_empty_dirs         = true
    should_log_invalid_rows   = false
    stage_path                = "...my_stage_path..."
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "filesystem"
    write_high_water_mark = 2700.71
  }
  output_google_chronicle = {
    api_key               = "...my_api_key..."
    api_key_secret        = "...my_api_key_secret..."
    api_version           = "v2"
    authentication_method = "serviceAccount"
    compress              = false
    concurrency           = 32
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
    failed_request_logging_mode = "none"
    flush_period_sec            = 7.12
    id                          = "...my_id..."
    log_format_type             = "unstructured"
    log_text_field              = "...my_log_text_field..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 0.88
    max_payload_size_kb         = 526.81
    namespace                   = "...my_namespace..."
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    region                            = "...my_region..."
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 4.69
        http_status     = 594.36
        initial_backoff = 421441.88
        max_backoff     = 111988.73
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
      backoff_rate    = 18.65
      initial_backoff = 151416.39
      max_backoff     = 44173.49
      timeout_retry   = true
    }
    timeout_sec           = 2580522638877335.5
    total_memory_limit_kb = 0.27
    type                  = "google_chronicle"
    use_round_robin_dns   = false
  }
  output_google_cloud_logging = {
    cache_fill_bytes_expression = "...my_cache_fill_bytes_expression..."
    cache_hit_expression        = "...my_cache_hit_expression..."
    cache_lookup_expression     = "...my_cache_lookup_expression..."
    cache_validated_expression  = "...my_cache_validated_expression..."
    concurrency                 = 4.3
    connection_timeout          = 0.71
    description                 = "...my_description..."
    environment                 = "...my_environment..."
    file_expression             = "...my_file_expression..."
    first_expression            = "...my_first_expression..."
    flush_period_sec            = 6.37
    function_expression         = "...my_function_expression..."
    google_auth_method          = "auto"
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
    log_location_type       = "organization"
    log_name_expression     = "...my_log_name_expression..."
    max_payload_events      = 8.23
    max_payload_size_kb     = 1437.62
    on_backpressure         = "drop"
    payload_expression      = "...my_payload_expression..."
    payload_format          = "json"
    pipeline                = "...my_pipeline..."
    pq_compress             = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size          = "...my_pq_max_file_size..."
    pq_max_size               = "...my_pq_max_size..."
    pq_mode                   = "always"
    pq_on_backpressure        = "block"
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
    throttle_rate_req_per_sec = 223
    timeout_sec               = 7182771915545540
    total_memory_limit_kb     = 7.11
    total_splits_expression   = "...my_total_splits_expression..."
    trace_expression          = "...my_trace_expression..."
    trace_sampled_expression  = "...my_trace_sampled_expression..."
    type                      = "google_cloud_logging"
    uid_expression            = "...my_uid_expression..."
    user_agent_expression     = "...my_user_agent_expression..."
  }
  output_google_cloud_storage = {
    add_id_to_stage_path      = false
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    compress                  = "none"
    compression_level         = "best_compression"
    deadletter_enabled        = true
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    dest_path                 = "...my_dest_path..."
    empty_dir_cleanup_sec     = 39665.96
    enable_page_checksum      = false
    enable_statistics         = true
    enable_write_page_index   = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    file_name_suffix          = "...my_file_name_suffix..."
    format                    = "json"
    header_line               = "...my_header_line..."
    id                        = "...my_id..."
    key_value_metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    max_file_idle_time_sec    = 1412.08
    max_file_open_time_sec    = 1398.16
    max_file_size_mb          = 129.8
    max_open_files            = 589.48
    max_retry_num             = 9.36
    object_acl                = "private"
    on_backpressure           = "block"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 33720757.45
    parquet_version           = "PARQUET_2_4"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    region                    = "...my_region..."
    reject_unauthorized       = true
    remove_empty_dirs         = false
    reuse_connections         = true
    should_log_invalid_rows   = true
    signature_version         = "v2"
    stage_path                = "...my_stage_path..."
    storage_class             = "NEARLINE"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "google_cloud_storage"
    verify_permissions    = true
    write_high_water_mark = 2313.39
  }
  output_google_pubsub = {
    batch_size         = 2745.71
    batch_timeout      = 87866.97
    create_topic       = false
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_period_sec   = "{ \"see\": \"documentation\" }"
    google_auth_method = "secret"
    id                 = "...my_id..."
    max_in_progress    = 96.45
    max_queue_size     = 10.49
    max_record_size_kb = 63.95
    on_backpressure    = "drop"
    ordered_delivery   = false
    pipeline           = "...my_pipeline..."
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size            = "...my_pq_max_file_size..."
    pq_max_size                 = "...my_pq_max_size..."
    pq_mode                     = "backpressure"
    pq_on_backpressure          = "drop"
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
    output_grafana_cloud_grafana_cloud2 = {
      compress    = false
      concurrency = 15.33
      description = "...my_description..."
      environment = "...my_environment..."
      extra_http_headers = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      failed_request_logging_mode = "payload"
      flush_period_sec            = 9.15
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
      max_payload_events  = 0.02
      max_payload_size_kb = 1467.85
      message             = "...my_message..."
      message_format      = "protobuf"
      metric_rename_expr  = "...my_metric_rename_expr..."
      on_backpressure     = "block"
      pipeline            = "...my_pipeline..."
      pq_compress         = "none"
      pq_controls = {
        # ...
      }
      pq_max_file_size   = "...my_pq_max_file_size..."
      pq_max_size        = "...my_pq_max_size..."
      pq_mode            = "backpressure"
      pq_on_backpressure = "drop"
      pq_path            = "...my_pq_path..."
      prometheus_auth = {
        auth_type          = "textSecret"
        credentials_secret = "...my_credentials_secret..."
        password           = "...my_password..."
        text_secret        = "...my_text_secret..."
        token              = "...my_token..."
        username           = "...my_username..."
      }
      prometheus_url                    = "...my_prometheus_url..."
      reject_unauthorized               = false
      response_honor_retry_after_header = true
      response_retry_settings = [
        {
          backoff_rate    = 2.13
          http_status     = 267.11
          initial_backoff = 422832.39
          max_backoff     = 148572.04
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
        backoff_rate    = 9.08
        initial_backoff = 337102.69
        max_backoff     = 59208.53
        timeout_retry   = true
      }
      timeout_sec         = 8481559183316778
      type                = "grafana_cloud"
      use_round_robin_dns = true
    }
  }
  output_graphite = {
    connection_timeout     = 8.2
    description            = "...my_description..."
    dns_resolve_period_sec = 71846.6
    environment            = "...my_environment..."
    flush_period_sec       = 0.82
    host                   = "...my_host..."
    id                     = "...my_id..."
    mtu                    = 6951.17
    on_backpressure        = "queue"
    pipeline               = "...my_pipeline..."
    port                   = 23741.86
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "backpressure"
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
    type                  = "graphite"
    write_timeout         = 6.25
  }
  output_honeycomb = {
    auth_type   = "secret"
    compress    = true
    concurrency = 22.28
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
    flush_period_sec            = 1.69
    id                          = "...my_id..."
    max_payload_events          = 5.21
    max_payload_size_kb         = 2486.93
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
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 15.95
        http_status     = 397.2
        initial_backoff = 415209.78
        max_backoff     = 50314.15
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
      backoff_rate    = 4.06
      initial_backoff = 261968.08
      max_backoff     = 159362.76
      timeout_retry   = true
    }
    timeout_sec         = 4600176163311220
    type                = "honeycomb"
    use_round_robin_dns = true
  }
  output_humio_hec = {
    auth_type   = "secret"
    compress    = true
    concurrency = 16.01
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 3.72
    format                      = "raw"
    id                          = "...my_id..."
    max_payload_events          = 2.66
    max_payload_size_kb         = 28200.73
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
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 4.49
        http_status     = 424.66
        initial_backoff = 409932.86
        max_backoff     = 115777.5
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
      backoff_rate    = 1.8
      initial_backoff = 53671.67
      max_backoff     = 170178.81
      timeout_retry   = true
    }
    timeout_sec         = 1047839844333980.6
    token               = "...my_token..."
    type                = "humio_hec"
    url                 = "...my_url..."
    use_round_robin_dns = false
  }
  output_influxdb = {
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "none"
    bucket                   = "...my_bucket..."
    compress                 = false
    concurrency              = 7.51
    credentials_secret       = "...my_credentials_secret..."
    database                 = "...my_database..."
    description              = "...my_description..."
    dynamic_value_field_name = true
    environment              = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 4.41
    id                          = "...my_id..."
    login_url                   = "...my_login_url..."
    max_payload_events          = 3.49
    max_payload_size_kb         = 28327.39
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
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 11.22
        http_status     = 259.83
        initial_backoff = 498223.3
        max_backoff     = 105504.37
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
      backoff_rate    = 12.68
      initial_backoff = 586999.02
      max_backoff     = 93135.19
      timeout_retry   = true
    }
    timeout_sec          = 7419738081071437
    timestamp_precision  = "ns"
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 114248.52
    type                 = "influxdb"
    url                  = "...my_url..."
    use_round_robin_dns  = false
    use_v2_api           = false
    username             = "...my_username..."
    value_field_name     = "...my_value_field_name..."
  }
  output_kafka = {
    ack                    = 0
    authentication_timeout = 696282.49
    backoff_rate           = 15.14
    brokers = [
      "..."
    ]
    compression        = "none"
    connection_timeout = 1903008.14
    description        = "...my_description..."
    environment        = "...my_environment..."
    flush_event_count  = 8789.01
    flush_period_sec   = 2.16
    format             = "json"
    id                 = "...my_id..."
    initial_backoff    = 325726.51
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = true
      }
      connection_timeout      = 45010.18
      default_key_schema_id   = 4.68
      default_value_schema_id = 9.01
      disabled                = true
      max_retries             = 59.54
      request_timeout         = 8414.4
      schema_registry_url     = "...my_schema_registry_url..."
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
    }
    max_back_off       = 68585.65
    max_record_size_kb = 2.51
    max_retries        = 83.79
    on_backpressure    = "block"
    pipeline           = "...my_pipeline..."
    pq_compress        = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "always"
    pq_on_backpressure         = "drop"
    pq_path                    = "...my_pq_path..."
    protobuf_library_id        = "...my_protobuf_library_id..."
    reauthentication_threshold = 124816.72
    request_timeout            = 1296275.02
    sasl = {
      disabled  = true
      mechanism = "plain"
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
      disabled            = true
      max_version         = "TLSv1.2"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
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
    concurrency               = 2.22
    description               = "...my_description..."
    duration_seconds          = 43147.47
    enable_assume_role        = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    flush_period_sec          = 8.09
    id                        = "...my_id..."
    max_record_size_kb        = 4786.04
    on_backpressure           = "drop"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "backpressure"
    pq_on_backpressure  = "block"
    pq_path             = "...my_pq_path..."
    region              = "...my_region..."
    reject_unauthorized = false
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
    auth_type          = "none"
    compress           = true
    concurrency        = 24.67
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 6.69
    id                          = "...my_id..."
    labels = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    max_payload_events  = 8.16
    max_payload_size_kb = 3180.75
    message             = "...my_message..."
    message_format      = "json"
    on_backpressure     = "drop"
    password            = "...my_password..."
    pipeline            = "...my_pipeline..."
    pq_compress         = "none"
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
        backoff_rate    = 11.29
        http_status     = 517.07
        initial_backoff = 2825.74
        max_backoff     = 83387.46
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
      backoff_rate    = 17.44
      initial_backoff = 384929.59
      max_backoff     = 13948.64
      timeout_retry   = true
    }
    timeout_sec           = 6338440315996480
    token                 = "...my_token..."
    total_memory_limit_kb = 6.52
    type                  = "loki"
    url                   = "...my_url..."
    use_round_robin_dns   = false
    username              = "...my_username..."
  }
  output_minio = {
    add_id_to_stage_path      = false
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "manual"
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
    empty_dir_cleanup_sec     = 25780.23
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
    max_concurrent_file_parts = 3.07
    max_file_idle_time_sec    = 17452.63
    max_file_open_time_sec    = 12033.42
    max_file_size_mb          = 347.47
    max_open_files            = 1121.64
    max_retry_num             = 1.24
    object_acl                = "public-read-write"
    on_backpressure           = "drop"
    on_disk_full_backpressure = "block"
    parquet_data_page_version = "DATA_PAGE_V2"
    parquet_page_size         = "...my_parquet_page_size..."
    parquet_row_group_length  = 18956028.18
    parquet_version           = "PARQUET_1_0"
    partition_expr            = "...my_partition_expr..."
    pipeline                  = "...my_pipeline..."
    region                    = "...my_region..."
    reject_unauthorized       = true
    remove_empty_dirs         = false
    reuse_connections         = true
    server_side_encryption    = "AES256"
    should_log_invalid_rows   = false
    signature_version         = "v4"
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
    write_high_water_mark = 1249.29
  }
  output_msk = {
    ack                       = 7
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    authentication_timeout    = 2076871.67
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    backoff_rate              = 8.41
    brokers = [
      "..."
    ]
    compression        = "none"
    connection_timeout = 2898984.28
    description        = "...my_description..."
    duration_seconds   = 3273.91
    enable_assume_role = false
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    flush_event_count  = 2775.55
    flush_period_sec   = 0.11
    format             = "raw"
    id                 = "...my_id..."
    initial_backoff    = 50366.77
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout      = 24426.31
      default_key_schema_id   = 3.14
      default_value_schema_id = 7.51
      disabled                = false
      max_retries             = 70.06
      request_timeout         = 38575.48
      schema_registry_url     = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = true
        max_version         = "TLSv1"
        min_version         = "TLSv1.2"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = true
        servername          = "...my_servername..."
      }
    }
    max_back_off       = 118252.89
    max_record_size_kb = 10.3
    max_retries        = 72.16
    on_backpressure    = "drop"
    pipeline           = "...my_pipeline..."
    pq_compress        = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size           = "...my_pq_max_file_size..."
    pq_max_size                = "...my_pq_max_size..."
    pq_mode                    = "error"
    pq_on_backpressure         = "drop"
    pq_path                    = "...my_pq_path..."
    protobuf_library_id        = "...my_protobuf_library_id..."
    reauthentication_threshold = 659875.3
    region                     = "...my_region..."
    reject_unauthorized        = true
    request_timeout            = 2709174.61
    reuse_connections          = true
    signature_version          = "v2"
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
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    topic = "...my_topic..."
    type  = "msk"
  }
  output_netflow = {
    description            = "...my_description..."
    dns_resolve_period_sec = 16590.34
    environment            = "...my_environment..."
    hosts = [
      {
        host = "...my_host..."
        port = 40081.81
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
    compress    = true
    concurrency = 6.46
    custom_url  = "...my_custom_url..."
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payloadAndHeaders"
    flush_period_sec            = 1.88
    id                          = "...my_id..."
    log_type                    = "...my_log_type..."
    max_payload_events          = 6.06
    max_payload_size_kb         = 680.88
    message_field               = "...my_message_field..."
    metadata = [
      {
        name  = "service"
        value = "...my_value..."
      }
    ]
    on_backpressure = "block"
    pipeline        = "...my_pipeline..."
    pq_compress     = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    region                            = "EU"
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 16.51
        http_status     = 425.25
        initial_backoff = 171075.58
        max_backoff     = 37898
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
      backoff_rate    = 14.66
      initial_backoff = 175241.79
      max_backoff     = 59881.47
      timeout_retry   = false
    }
    timeout_sec           = 4368732314350300
    total_memory_limit_kb = 7.76
    type                  = "newrelic"
    use_round_robin_dns   = true
  }
  output_newrelic_events = {
    account_id  = "...my_account_id..."
    api_key     = "...my_api_key..."
    auth_type   = "secret"
    compress    = false
    concurrency = 26.87
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
    failed_request_logging_mode = "payload"
    flush_period_sec            = 4.94
    id                          = "...my_id..."
    max_payload_events          = 4.03
    max_payload_size_kb         = 192.87
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "drop"
    pq_path                           = "...my_pq_path..."
    region                            = "Custom"
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 6.24
        http_status     = 571.62
        initial_backoff = 407273.39
        max_backoff     = 91757.45
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
      backoff_rate    = 15.9
      initial_backoff = 67047.09
      max_backoff     = 39664.88
      timeout_retry   = false
    }
    timeout_sec         = 3142152222818271.5
    type                = "newrelic_events"
    use_round_robin_dns = false
  }
  output_open_telemetry = {
    auth_header_expr   = "...my_auth_header_expr..."
    auth_type          = "oauth"
    compress           = "gzip"
    concurrency        = 29.63
    connection_timeout = 0.13
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
    failed_request_logging_mode    = "payloadAndHeaders"
    flush_period_sec               = 2.52
    http_compress                  = "gzip"
    http_logs_endpoint_override    = "...my_http_logs_endpoint_override..."
    http_metrics_endpoint_override = "...my_http_metrics_endpoint_override..."
    http_traces_endpoint_override  = "...my_http_traces_endpoint_override..."
    id                             = "...my_id..."
    keep_alive                     = false
    keep_alive_time                = 7.49
    login_url                      = "...my_login_url..."
    max_payload_size_kb            = 2009.53
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
    on_backpressure = "queue"
    otlp_version    = "1.3.1"
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
    protocol                          = "http"
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 16.93
        http_status     = 545.85
        initial_backoff = 54922.96
        max_backoff     = 145087.51
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
      backoff_rate    = 6.4
      initial_backoff = 215425.68
      max_backoff     = 118165.54
      timeout_retry   = true
    }
    timeout_sec = 7703655941061551
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 219683.13
    type                 = "open_telemetry"
    use_round_robin_dns  = true
    username             = "...my_username..."
  }
  output_prometheus = {
    auth_header_expr   = "...my_auth_header_expr..."
    auth_type          = "none"
    concurrency        = 13.24
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 7.15
    id                          = "...my_id..."
    login_url                   = "...my_login_url..."
    max_payload_events          = 1.38
    max_payload_size_kb         = 3560.68
    metric_rename_expr          = "...my_metric_rename_expr..."
    metrics_flush_period_sec    = 2.57
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
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 6.42
        http_status     = 577.58
        initial_backoff = 539774.25
        max_backoff     = 164557.04
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
      backoff_rate    = 11.67
      initial_backoff = 572427.79
      max_backoff     = 28087.22
      timeout_retry   = true
    }
    timeout_sec          = 8357600847636273
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 130270.62
    type                 = "prometheus"
    url                  = "...my_url..."
    use_round_robin_dns  = false
    username             = "...my_username..."
  }
  output_ring = {
    compress        = "gzip"
    description     = "...my_description..."
    dest_path       = "...my_dest_path..."
    environment     = "...my_environment..."
    format          = "json"
    id              = "...my_id..."
    max_data_size   = "...my_max_data_size..."
    max_data_time   = "...my_max_data_time..."
    on_backpressure = "drop"
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
    automatic_schema          = false
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
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
    duration_seconds          = 38731.79
    empty_dir_cleanup_sec     = 17958.77
    enable_assume_role        = true
    enable_page_checksum      = true
    enable_statistics         = true
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
    max_closing_files_to_backpressure = 4050.63
    max_concurrent_file_parts         = 3.72
    max_file_idle_time_sec            = 30875.25
    max_file_open_time_sec            = 4976.89
    max_file_size_mb                  = 281.07
    max_open_files                    = 448.99
    max_retry_num                     = 2.7
    object_acl                        = "bucket-owner-read"
    on_backpressure                   = "drop"
    on_disk_full_backpressure         = "drop"
    parquet_data_page_version         = "DATA_PAGE_V2"
    parquet_page_size                 = "...my_parquet_page_size..."
    parquet_row_group_length          = 58778119.74
    parquet_version                   = "PARQUET_1_0"
    partition_expr                    = "...my_partition_expr..."
    pipeline                          = "...my_pipeline..."
    region                            = "...my_region..."
    reject_unauthorized               = false
    remove_empty_dirs                 = true
    reuse_connections                 = false
    server_side_encryption            = "aws:kms"
    should_log_invalid_rows           = true
    signature_version                 = "v4"
    stage_path                        = "...my_stage_path..."
    storage_class                     = "STANDARD_IA"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "s3"
    verify_permissions    = true
    write_high_water_mark = 2926.5
  }
  output_security_lake = {
    account_id                = "...my_account_id..."
    add_id_to_stage_path      = true
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    automatic_schema          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "manual"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    base_file_name            = "...my_base_file_name..."
    bucket                    = "...my_bucket..."
    custom_source             = "...my_custom_source..."
    deadletter_enabled        = true
    deadletter_path           = "...my_deadletter_path..."
    description               = "...my_description..."
    duration_seconds          = 3293.42
    empty_dir_cleanup_sec     = 17679.05
    enable_assume_role        = false
    enable_page_checksum      = false
    enable_statistics         = true
    enable_write_page_index   = false
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
    max_closing_files_to_backpressure = 4013.47
    max_concurrent_file_parts         = 2.43
    max_file_idle_time_sec            = 66555.1
    max_file_open_time_sec            = 69415.39
    max_file_size_mb                  = 589.48
    max_open_files                    = 166.27
    max_retry_num                     = 8.34
    object_acl                        = "authenticated-read"
    on_backpressure                   = "drop"
    on_disk_full_backpressure         = "drop"
    parquet_data_page_version         = "DATA_PAGE_V1"
    parquet_page_size                 = "...my_parquet_page_size..."
    parquet_row_group_length          = 18227287.19
    parquet_schema                    = "...my_parquet_schema..."
    parquet_version                   = "PARQUET_2_6"
    pipeline                          = "...my_pipeline..."
    region                            = "...my_region..."
    reject_unauthorized               = true
    remove_empty_dirs                 = false
    reuse_connections                 = false
    server_side_encryption            = "AES256"
    should_log_invalid_rows           = false
    signature_version                 = "v2"
    stage_path                        = "...my_stage_path..."
    storage_class                     = "STANDARD_IA"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type                  = "security_lake"
    verify_permissions    = false
    write_high_water_mark = 2215.45
  }
  output_sentinel = {
    advanced_content_type      = "...my_advanced_content_type..."
    auth_type                  = "oauth"
    client_id                  = "...my_client_id..."
    compress                   = true
    concurrency                = 3.35
    custom_content_type        = "...my_custom_content_type..."
    custom_drop_when_null      = false
    custom_event_delimiter     = "...my_custom_event_delimiter..."
    custom_payload_expression  = "...my_custom_payload_expression..."
    custom_source_expression   = "...my_custom_source_expression..."
    dce_endpoint               = "...my_dce_endpoint..."
    dcr_id                     = "...my_dcr_id..."
    description                = "...my_description..."
    endpoint_url_configuration = "ID"
    environment                = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "payload"
    flush_period_sec            = 1.35
    format                      = "advanced"
    format_event_code           = "...my_format_event_code..."
    format_payload_code         = "...my_format_payload_code..."
    id                          = "...my_id..."
    keep_alive                  = false
    login_url                   = "...my_login_url..."
    max_payload_events          = 9.72
    max_payload_size_kb         = 174.22
    on_backpressure             = "block"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
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
        backoff_rate    = 17.07
        http_status     = 253.85
        initial_backoff = 289884.68
        max_backoff     = 111728.25
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
      backoff_rate    = 15.12
      initial_backoff = 89681.06
      max_backoff     = 40934.45
      timeout_retry   = false
    }
    timeout_sec           = 4083333491358569
    total_memory_limit_kb = 4.85
    type                  = "sentinel"
    url                   = "...my_url..."
    use_round_robin_dns   = false
  }
  output_service_now = {
    auth_token_name    = "...my_auth_token_name..."
    compress           = "deflate"
    concurrency        = 29.02
    connection_timeout = 1.68
    description        = "...my_description..."
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode    = "payloadAndHeaders"
    flush_period_sec               = 2.62
    http_compress                  = "gzip"
    http_logs_endpoint_override    = "...my_http_logs_endpoint_override..."
    http_metrics_endpoint_override = "...my_http_metrics_endpoint_override..."
    http_traces_endpoint_override  = "...my_http_traces_endpoint_override..."
    id                             = "...my_id..."
    keep_alive                     = false
    keep_alive_time                = 3.42
    max_payload_size_kb            = 7580.24
    metadata = [
      {
        key   = "...my_key..."
        value = "...my_value..."
      }
    ]
    on_backpressure = "block"
    otlp_version    = "1.3.1"
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
    protocol                          = "grpc"
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 3.18
        http_status     = 376.33
        initial_backoff = 515752.4
        max_backoff     = 133676.75
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
      backoff_rate    = 8.34
      initial_backoff = 91476.44
      max_backoff     = 150501.62
      timeout_retry   = true
    }
    timeout_sec = 3110546845803040.5
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
    }
    token_secret        = "...my_token_secret..."
    type                = "service_now"
    use_round_robin_dns = true
  }
  output_signalfx = {
    auth_type   = "manual"
    compress    = false
    concurrency = 18.63
    description = "...my_description..."
    environment = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode = "none"
    flush_period_sec            = 0.74
    id                          = "...my_id..."
    max_payload_events          = 9.85
    max_payload_size_kb         = 4662.99
    on_backpressure             = "queue"
    pipeline                    = "...my_pipeline..."
    pq_compress                 = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    realm                             = "...my_realm..."
    reject_unauthorized               = true
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 18.64
        http_status     = 107.68
        initial_backoff = 45304.03
        max_backoff     = 162323.4
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
      backoff_rate    = 16.46
      initial_backoff = 309073.74
      max_backoff     = 177854.6
      timeout_retry   = true
    }
    timeout_sec         = 4162869804131207
    token               = "...my_token..."
    type                = "signalfx"
    use_round_robin_dns = false
  }
  output_snmp = {
    description            = "...my_description..."
    dns_resolve_period_sec = 67309.88
    environment            = "...my_environment..."
    hosts = [
      {
        host = "...my_host..."
        port = 4562.26
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
    duration_seconds          = 33859.9
    enable_assume_role        = false
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    id                        = "...my_id..."
    max_retries               = 2.37
    message_group_id          = "...my_message_group_id..."
    on_backpressure           = "drop"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "error"
    pq_on_backpressure  = "drop"
    pq_path             = "...my_pq_path..."
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
    topic_arn = "...my_topic_arn..."
    type      = "sns"
  }
  output_splunk = {
    auth_token               = "...my_auth_token..."
    auth_type                = "secret"
    compress                 = "always"
    connection_timeout       = 3.87
    description              = "...my_description..."
    enable_ack               = true
    enable_multi_metrics     = false
    environment              = "...my_environment..."
    host                     = "...my_host..."
    id                       = "...my_id..."
    log_failed_requests      = true
    max_failed_health_checks = 8.62
    max_s2_sversion          = "v3"
    nested_fields            = "json"
    on_backpressure          = "drop"
    pipeline                 = "...my_pipeline..."
    port                     = 9699.07
    pq_compress              = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "always"
    pq_on_backpressure = "block"
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
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    type          = "splunk"
    write_timeout = 5.14
  }
  output_splunk_hec = {
    auth_type              = "manual"
    compress               = true
    concurrency            = 31.06
    description            = "...my_description..."
    dns_resolve_period_sec = 84122.87
    enable_multi_metrics   = true
    environment            = "...my_environment..."
    exclude_self           = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 8.12
    id                            = "...my_id..."
    load_balance_stats_period_sec = 13.07
    load_balanced                 = false
    max_payload_events            = 1.94
    max_payload_size_kb           = 995482.79
    next_queue                    = "...my_next_queue..."
    on_backpressure               = "block"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "always"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 18.77
        http_status     = 103.29
        initial_backoff = 394109.29
        max_backoff     = 43786.54
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
      backoff_rate    = 13.27
      initial_backoff = 282045.13
      max_backoff     = 40548.41
      timeout_retry   = true
    }
    timeout_sec = 4120082867142921
    token       = "...my_token..."
    type        = "splunk_hec"
    url         = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 0.36
      }
    ]
    use_round_robin_dns = false
  }
  output_splunk_lb = {
    auth_token             = "...my_auth_token..."
    auth_type              = "manual"
    compress               = "auto"
    connection_timeout     = 6.82
    description            = "...my_description..."
    dns_resolve_period_sec = 42455.46
    enable_ack             = true
    enable_multi_metrics   = true
    environment            = "...my_environment..."
    exclude_self           = false
    hosts = [
      {
        host       = "...my_host..."
        port       = 39320.05
        servername = "...my_servername..."
        tls        = "off"
        weight     = 0.71
      }
    ]
    id                = "...my_id..."
    indexer_discovery = false
    indexer_discovery_configs = {
      auth_token = "...my_auth_token..."
      auth_tokens = [
        {
          auth_type = "secret"
        }
      ]
      auth_type            = "secret"
      master_uri           = "...my_master_uri..."
      refresh_interval_sec = 34502.49
      reject_unauthorized  = false
      site                 = "...my_site..."
      text_secret          = "...my_text_secret..."
    }
    load_balance_stats_period_sec = 15.17
    log_failed_requests           = false
    max_concurrent_senders        = 2.63
    max_failed_health_checks      = 1.7
    max_s2_sversion               = "v3"
    nested_fields                 = "json"
    on_backpressure               = "drop"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                = "...my_pq_max_file_size..."
    pq_max_size                     = "...my_pq_max_size..."
    pq_mode                         = "backpressure"
    pq_on_backpressure              = "block"
    pq_path                         = "...my_pq_path..."
    sender_unhealthy_time_allowance = 20354.85
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
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    type          = "splunk_lb"
    write_timeout = 1.89
  }
  output_sqs = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    create_queue              = true
    description               = "...my_description..."
    duration_seconds          = 13455.59
    enable_assume_role        = true
    endpoint                  = "...my_endpoint..."
    environment               = "...my_environment..."
    flush_period_sec          = 2.67
    id                        = "...my_id..."
    max_in_progress           = 39.9
    max_queue_size            = 10.4
    max_record_size_kb        = 254.35
    message_group_id          = "...my_message_group_id..."
    on_backpressure           = "queue"
    pipeline                  = "...my_pipeline..."
    pq_compress               = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size    = "...my_pq_max_file_size..."
    pq_max_size         = "...my_pq_max_size..."
    pq_mode             = "always"
    pq_on_backpressure  = "drop"
    pq_path             = "...my_pq_path..."
    queue_name          = "...my_queue_name..."
    queue_type          = "standard"
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = true
    signature_version   = "v2"
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    type = "sqs"
  }
  output_statsd = {
    connection_timeout     = 8.08
    description            = "...my_description..."
    dns_resolve_period_sec = 41630.76
    environment            = "...my_environment..."
    flush_period_sec       = 9.91
    host                   = "...my_host..."
    id                     = "...my_id..."
    mtu                    = 16935.61
    on_backpressure        = "block"
    pipeline               = "...my_pipeline..."
    port                   = 56247.95
    pq_compress            = "none"
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
    type                  = "statsd"
    write_timeout         = 1.5
  }
  output_statsd_ext = {
    connection_timeout     = 5.72
    description            = "...my_description..."
    dns_resolve_period_sec = 11011.37
    environment            = "...my_environment..."
    flush_period_sec       = 2.05
    host                   = "...my_host..."
    id                     = "...my_id..."
    mtu                    = 43787.67
    on_backpressure        = "block"
    pipeline               = "...my_pipeline..."
    port                   = 7342.64
    pq_compress            = "gzip"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "error"
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
    type                  = "statsd_ext"
    write_timeout         = 5.49
  }
  output_sumo_logic = {
    compress        = true
    concurrency     = 29.75
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
    flush_period_sec            = 7.9
    format                      = "raw"
    id                          = "...my_id..."
    max_payload_events          = 8.27
    max_payload_size_kb         = 598.71
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
        backoff_rate    = 16.93
        http_status     = 540.95
        initial_backoff = 120724.02
        max_backoff     = 10794.57
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
      backoff_rate    = 11.2
      initial_backoff = 48916.08
      max_backoff     = 81782.11
      timeout_retry   = false
    }
    timeout_sec           = 6917100564962345
    total_memory_limit_kb = 0.6
    type                  = "sumo_logic"
    url                   = "...my_url..."
    use_round_robin_dns   = false
  }
  output_syslog = {
    app_name            = "...my_app_name..."
    connection_timeout  = 7.06
    description         = "...my_description..."
    environment         = "...my_environment..."
    facility            = 6
    host                = "...my_host..."
    id                  = "...my_id..."
    load_balanced       = false
    log_failed_requests = false
    max_record_size     = 46285.92
    message_format      = "rfc5424"
    octet_count_framing = false
    on_backpressure     = "drop"
    pipeline            = "...my_pipeline..."
    port                = 53044.62
    pq_compress         = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "always"
    pq_on_backpressure = "block"
    pq_path            = "...my_pq_path..."
    protocol           = "tcp"
    severity           = 4
    streamtags = [
      "..."
    ]
    system_fields = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    timestamp_format      = "syslog"
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
    type                       = "syslog"
    udp_dns_resolve_period_sec = 42130.55
    write_timeout              = 0.73
  }
  output_tcpjson = {
    auth_token             = "...my_auth_token..."
    auth_type              = "manual"
    compression            = "gzip"
    connection_timeout     = 1.99
    description            = "...my_description..."
    dns_resolve_period_sec = 2753.35
    environment            = "...my_environment..."
    exclude_self           = true
    host                   = "...my_host..."
    hosts = [
      {
        host       = "...my_host..."
        port       = 21739.5
        servername = "...my_servername..."
        tls        = "inherit"
        weight     = 7.62
      }
    ]
    id                            = "...my_id..."
    load_balance_stats_period_sec = 13.63
    load_balanced                 = true
    log_failed_requests           = true
    max_concurrent_senders        = 0.26
    on_backpressure               = "drop"
    pipeline                      = "...my_pipeline..."
    port                          = 27648.54
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size   = "...my_pq_max_file_size..."
    pq_max_size        = "...my_pq_max_size..."
    pq_mode            = "always"
    pq_on_backpressure = "block"
    pq_path            = "...my_pq_path..."
    send_header        = true
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
      max_version         = "TLSv1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    token_ttl_minutes = 43.32
    type              = "tcpjson"
    write_timeout     = 5.91
  }
  output_wavefront = {
    auth_type   = "secret"
    compress    = false
    concurrency = 10.29
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
    flush_period_sec            = 7.08
    id                          = "...my_id..."
    max_payload_events          = 3.19
    max_payload_size_kb         = 5322.37
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
    reject_unauthorized               = true
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 2.91
        http_status     = 440.17
        initial_backoff = 552631.89
        max_backoff     = 161549.44
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
      backoff_rate    = 10.59
      initial_backoff = 279076.37
      max_backoff     = 39734.1
      timeout_retry   = false
    }
    timeout_sec         = 3176122588065585.5
    token               = "...my_token..."
    type                = "wavefront"
    use_round_robin_dns = true
  }
  output_webhook = {
    advanced_content_type     = "...my_advanced_content_type..."
    auth_header_expr          = "...my_auth_header_expr..."
    auth_type                 = "credentialsSecret"
    compress                  = true
    concurrency               = 28.59
    credentials_secret        = "...my_credentials_secret..."
    custom_content_type       = "...my_custom_content_type..."
    custom_drop_when_null     = true
    custom_event_delimiter    = "...my_custom_event_delimiter..."
    custom_payload_expression = "...my_custom_payload_expression..."
    custom_source_expression  = "...my_custom_source_expression..."
    description               = "...my_description..."
    dns_resolve_period_sec    = 17531.87
    environment               = "...my_environment..."
    exclude_self              = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 5.36
    format                        = "ndjson"
    format_event_code             = "...my_format_event_code..."
    format_payload_code           = "...my_format_payload_code..."
    id                            = "...my_id..."
    keep_alive                    = false
    load_balance_stats_period_sec = 19.23
    load_balanced                 = false
    login_url                     = "...my_login_url..."
    max_payload_events            = 2.31
    max_payload_size_kb           = 418826.72
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
    on_backpressure = "drop"
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
    reject_unauthorized               = false
    response_honor_retry_after_header = true
    response_retry_settings = [
      {
        backoff_rate    = 5.68
        http_status     = 419.31
        initial_backoff = 114095.05
        max_backoff     = 142192.26
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
      backoff_rate    = 13.32
      initial_backoff = 491754.08
      max_backoff     = 28265.56
      timeout_retry   = true
    }
    timeout_sec = 6402570494355455
    tls = {
      ca_path          = "...my_ca_path..."
      cert_path        = "...my_cert_path..."
      certificate_name = "...my_certificate_name..."
      disabled         = false
      max_version      = "TLSv1.3"
      min_version      = "TLSv1.3"
      passphrase       = "...my_passphrase..."
      priv_key_path    = "...my_priv_key_path..."
      servername       = "...my_servername..."
    }
    token                 = "...my_token..."
    token_attribute_name  = "...my_token_attribute_name..."
    token_timeout_secs    = 235590.21
    total_memory_limit_kb = 5.04
    type                  = "webhook"
    url                   = "...my_url..."
    urls = [
      {
        url    = "...my_url..."
        weight = 6.35
      }
    ]
    use_round_robin_dns = false
    username            = "...my_username..."
  }
  output_xsiam = {
    auth_type              = "secret"
    concurrency            = 1.74
    description            = "...my_description..."
    dns_resolve_period_sec = 24588.01
    environment            = "...my_environment..."
    exclude_self           = true
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    failed_request_logging_mode   = "payload"
    flush_period_sec              = 7.14
    id                            = "...my_id..."
    load_balance_stats_period_sec = 12.36
    load_balanced                 = true
    max_payload_events            = 7.97
    max_payload_size_kb           = 2387.65
    on_backpressure               = "drop"
    pipeline                      = "...my_pipeline..."
    pq_compress                   = "none"
    pq_controls = {
      # ...
    }
    pq_max_file_size                  = "...my_pq_max_file_size..."
    pq_max_size                       = "...my_pq_max_size..."
    pq_mode                           = "error"
    pq_on_backpressure                = "block"
    pq_path                           = "...my_pq_path..."
    reject_unauthorized               = false
    response_honor_retry_after_header = false
    response_retry_settings = [
      {
        backoff_rate    = 8.54
        http_status     = 298.73
        initial_backoff = 57567.81
        max_backoff     = 59850.71
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
    throttle_rate_req_per_sec = 909
    timeout_retry_settings = {
      backoff_rate    = 10.5
      initial_backoff = 118601.76
      max_backoff     = 85354.16
      timeout_retry   = true
    }
    timeout_sec           = 7621312263336938
    token                 = "...my_token..."
    total_memory_limit_kb = 4.19
    type                  = "xsiam"
    url                   = "...my_url..."
    urls = [
      {
        url    = "{ \"see\": \"documentation\" }"
        weight = 3.29
      }
    ]
    use_round_robin_dns = true
  }
}