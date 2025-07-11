resource "criblio_source" "my_source" {
  group_id = "...my_group_id..."
  id       = "...my_id..."
  input_appscope = {
    auth_token = "...my_auth_token..."
    auth_type  = "manual"
    breaker_rulesets = [
      "..."
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "...my_description..."
    disabled            = false
    enable_proxy_header = false
    enable_unix_path    = true
    environment         = "...my_environment..."
    filter = {
      allow = [
        {
          arg      = "...my_arg..."
          config   = "...my_config..."
          procname = "...my_procname..."
        }
      ]
      transport_url = "...my_transport_url..."
    }
    host               = "...my_host..."
    id                 = "...my_id..."
    ip_whitelist_regex = "...my_ip_whitelist_regex..."
    max_active_cxn     = 9.71
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      dest_path     = "...my_dest_path..."
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    port     = 10821.34
    pq = {
      commit_frequency = 7.46
      compress         = "none"
      max_buffer_size  = 49.62
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = false
    socket_ending_max_wait = 1.44
    socket_idle_timeout    = 1.69
    socket_max_lifespan    = 4.11
    stale_channel_flush_ms = 34410614.39
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type              = "appscope"
    unix_socket_path  = "...my_unix_socket_path..."
    unix_socket_perms = "...my_unix_socket_perms..."
  }
  input_azure_blob = {
    auth_type   = "manual"
    azure_cloud = "...my_azure_cloud..."
    breaker_rulesets = [
      "..."
    ]
    certificate = {
      certificate_name = "...my_certificate_name..."
    }
    client_id          = "...my_client_id..."
    client_text_secret = "...my_client_text_secret..."
    connection_string  = "...my_connection_string..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description     = "...my_description..."
    disabled        = true
    endpoint_suffix = "...my_endpoint_suffix..."
    environment     = "...my_environment..."
    file_filter     = "...my_file_filter..."
    id              = "...my_id..."
    max_messages    = 29.53
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 76.24
    parquet_chunk_download_timeout = 2916.04
    parquet_chunk_size_mb          = 82.71
    pipeline                       = "...my_pipeline..."
    pq = {
      commit_frequency = 2.18
      compress         = "none"
      max_buffer_size  = 47.2
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = false
    queue_name             = "...my_queue_name..."
    send_to_routes         = false
    service_period_secs    = 8.37
    skip_on_error          = false
    stale_channel_flush_ms = 42632985.63
    storage_account_name   = "...my_storage_account_name..."
    streamtags = [
      "..."
    ]
    tenant_id          = "...my_tenant_id..."
    text_secret        = "...my_text_secret..."
    type               = "azure_blob"
    visibility_timeout = 246960.5
  }
  input_collection = {
    breaker_rulesets = [
      "..."
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    disabled    = false
    environment = "...my_environment..."
    id          = "...my_id..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    output   = "...my_output..."
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 4.72
      compress         = "gzip"
      max_buffer_size  = 42.52
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = true
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = true
    }
    send_to_routes         = false
    stale_channel_flush_ms = 15919396.49
    streamtags = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    type                  = "collection"
  }
  input_confluent_cloud = {
    authentication_timeout = 3228010.08
    auto_commit_interval   = 574056.62
    auto_commit_threshold  = 6088.43
    backoff_rate           = 19.92
    brokers = [
      "..."
    ]
    connection_timeout = 2896733.84
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    environment        = "...my_environment..."
    from_beginning     = true
    group_id           = "...my_group_id..."
    heartbeat_interval = 2878745.25
    id                 = "...my_id..."
    initial_backoff    = 385771.01
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = true
      }
      connection_timeout  = 2801.98
      disabled            = true
      max_retries         = 54.65
      request_timeout     = 2510.62
      schema_registry_url = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = false
        max_version         = "TLSv1"
        min_version         = "TLSv1.2"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = true
        servername          = "...my_servername..."
      }
    }
    max_back_off            = 138634.92
    max_bytes               = 153342787.93
    max_bytes_per_partition = 1122269.7
    max_retries             = 39.15
    max_socket_errors       = 26.42
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 10.5
      compress         = "none"
      max_buffer_size  = 48.51
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled                 = false
    reauthentication_threshold = 715684.29
    rebalance_timeout          = 655450.51
    request_timeout            = 1506422.44
    sasl = {
      disabled  = false
      mechanism = "plain"
    }
    send_to_routes  = false
    session_timeout = 865848.25
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      servername          = "...my_servername..."
    }
    topics = [
      "..."
    ]
    type = "confluent_cloud"
  }
  input_cribl = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description = "...my_description..."
    disabled    = false
    environment = "...my_environment..."
    filter      = "...my_filter..."
    id          = "...my_id..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 7.2
      compress         = "none"
      max_buffer_size  = 51.06
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled     = true
    send_to_routes = true
    streamtags = [
      "..."
    ]
    type = "cribl"
  }
  input_cribl_http = {
    activity_log_sample_rate = 2.54
    auth_tokens = [
      "..."
    ]
    capture_headers = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = false
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 303.71
    max_active_req          = 0.57
    max_requests_per_socket = 4
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 4736.58
    pq = {
      commit_frequency = 5.31
      compress         = "none"
      max_buffer_size  = 47.65
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 1.46
    send_to_routes  = true
    socket_timeout  = 0.8
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "cribl_http"
  }
  input_cribl_lake_http = {
    activity_log_sample_rate = 3.4
    auth_tokens = [
      "..."
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = true
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 51.28
    max_active_req          = 2.72
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 3239.03
    pq = {
      commit_frequency = 3.02
      compress         = "gzip"
      max_buffer_size  = 45.66
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 7.38
    send_to_routes  = false
    socket_timeout  = 5.22
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "cribl_lake_http"
  }
  input_criblmetrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description   = "...my_description..."
    disabled      = false
    environment   = "...my_environment..."
    full_fidelity = true
    id            = "...my_id..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 10.03
      compress         = "none"
      max_buffer_size  = 49.75
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled     = false
    prefix         = "...my_prefix..."
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "criblmetrics"
  }
  input_cribl_tcp = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disabled              = false
    enable_load_balancing = false
    enable_proxy_header   = true
    environment           = "...my_environment..."
    host                  = "...my_host..."
    id                    = "...my_id..."
    max_active_cxn        = 2.66
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 44126.06
    pq = {
      commit_frequency = 7.04
      compress         = "gzip"
      max_buffer_size  = 51.37
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = false
    socket_ending_max_wait = 7.39
    socket_idle_timeout    = 2.54
    socket_max_lifespan    = 7.93
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "cribl_tcp"
  }
  input_crowdstrike = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    breaker_rulesets = [
      "..."
    ]
    checkpointing = {
      enabled = true
      retries = 70.42
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = false
    duration_seconds       = 21127.93
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "...my_encoding..."
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_messages           = 8.34
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers = 51.77
    pipeline      = "...my_pipeline..."
    poll_timeout  = 7.64
    pq = {
      commit_frequency = 4.56
      compress         = "gzip"
      max_buffer_size  = 50.14
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = true
    }
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v4"
    skip_on_error          = true
    socket_timeout         = 39471.85
    stale_channel_flush_ms = 26602939.6
    streamtags = [
      "..."
    ]
    tag_after_processing = "false"
    type                 = "crowdstrike"
    visibility_timeout   = 38346.71
  }
  input_datadog_agent = {
    activity_log_sample_rate = 7.52
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = true
    environment             = "...my_environment..."
    extract_metrics         = true
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 336.23
    max_active_req          = 6.62
    max_requests_per_socket = 8
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 32867.01
    pq = {
      commit_frequency = 8.04
      compress         = "gzip"
      max_buffer_size  = 46.32
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = false
    proxy_mode = {
      enabled             = false
      reject_unauthorized = true
    }
    request_timeout = 0.77
    send_to_routes  = true
    socket_timeout  = 5.24
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "datadog_agent"
  }
  input_datagen = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description = "...my_description..."
    disabled    = false
    environment = "...my_environment..."
    id          = "...my_id..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 2.34
      compress         = "gzip"
      max_buffer_size  = 45.98
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = false
    samples = [
      {
        events_per_sec = 10.97
        sample         = "...my_sample..."
      }
    ]
    send_to_routes = true
    streamtags = [
      "..."
    ]
    type = "datagen"
  }
  input_edge_prometheus = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    auth_type                 = "secret"
    aws_authentication_method = "secret"
    aws_secret_key            = "...my_aws_secret_key..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    dimension_list = [
    ]
    disabled           = true
    discovery_type     = "k8s-pods"
    duration_seconds   = 6194.66
    enable_assume_role = true
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    id                 = "...my_id..."
    interval           = 5.37
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    name_list = [
      "..."
    ]
    password = "...my_password..."
    persistence = {
      compress      = "gzip"
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pod_filter = [
      {
        description = "...my_description..."
        filter      = "...my_filter..."
      }
    ]
    pq = {
      commit_frequency = 10.89
      compress         = "gzip"
      max_buffer_size  = 48.13
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled           = false
    record_type          = "SRV"
    region               = "...my_region..."
    reject_unauthorized  = false
    reuse_connections    = true
    scrape_path          = "...my_scrape_path..."
    scrape_path_expr     = "...my_scrape_path_expr..."
    scrape_port          = 37287.67
    scrape_port_expr     = "...my_scrape_port_expr..."
    scrape_protocol      = "https"
    scrape_protocol_expr = "...my_scrape_protocol_expr..."
    search_filter = [
      {
        name = "...my_name..."
        values = [
        ]
      }
    ]
    send_to_routes    = true
    signature_version = "v2"
    streamtags = [
      "..."
    ]
    targets = [
      {
        host     = "...my_host..."
        path     = "...my_path..."
        port     = 52043
        protocol = "http"
      }
    ]
    timeout       = 51072.71
    type          = "edge_prometheus"
    use_public_ip = false
    username      = "...my_username..."
  }
  input_elastic = {
    activity_log_sample_rate = 6.95
    api_version              = "8.3.2"
    auth_tokens = [
      "..."
    ]
    auth_type       = "credentialsSecret"
    capture_headers = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret  = "...my_credentials_secret..."
    custom_api_version  = "...my_custom_api_version..."
    description         = "...my_description..."
    disabled            = true
    elastic_api         = "...my_elastic_api..."
    enable_health_check = false
    enable_proxy_header = false
    environment         = "...my_environment..."
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 211.92
    max_active_req          = 8.3
    max_requests_per_socket = 6
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    password = "...my_password..."
    pipeline = "...my_pipeline..."
    port     = 63401.42
    pq = {
      commit_frequency = 4.42
      compress         = "gzip"
      max_buffer_size  = 44.97
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = true
    proxy_mode = {
      auth_type           = "secret"
      enabled             = false
      reject_unauthorized = false
      remove_headers = [
      ]
      timeout_sec = 1740239672156922.5
      url         = "...my_url..."
    }
    request_timeout = 9.56
    send_to_routes  = false
    socket_timeout  = 5.73
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type     = "elastic"
    username = "...my_username..."
  }
  input_eventhub = {
    authentication_timeout = 3263569.37
    auto_commit_interval   = 2913798.47
    auto_commit_threshold  = 3921.24
    backoff_rate           = 9.93
    brokers = [
      "..."
    ]
    connection_timeout = 952191.45
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    environment             = "...my_environment..."
    from_beginning          = true
    group_id                = "...my_group_id..."
    heartbeat_interval      = 1250654.13
    id                      = "...my_id..."
    initial_backoff         = 142635.29
    max_back_off            = 116091.33
    max_bytes               = 495245237.32
    max_bytes_per_partition = 9802527.45
    max_retries             = 96.14
    max_socket_errors       = 80.56
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    minimize_duplicates = false
    pipeline            = "...my_pipeline..."
    pq = {
      commit_frequency = 8.3
      compress         = "none"
      max_buffer_size  = 42.21
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                 = true
    reauthentication_threshold = 1753553.62
    rebalance_timeout          = 3367948.39
    request_timeout            = 8615.05
    sasl = {
      disabled  = true
      mechanism = "oauthbearer"
    }
    send_to_routes  = true
    session_timeout = 199043.15
    streamtags = [
      "..."
    ]
    tls = {
      disabled            = false
      reject_unauthorized = true
    }
    topics = [
      "..."
    ]
    type = "eventhub"
  }
  input_exec = {
    breaker_rulesets = [
      "..."
    ]
    command = "...my_command..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    cron_schedule = "...my_cron_schedule..."
    description   = "...my_description..."
    disabled      = false
    environment   = "...my_environment..."
    id            = "...my_id..."
    interval      = 4.15
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 5.86
      compress         = "none"
      max_buffer_size  = 44.63
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    retries                = 8.16
    schedule_type          = "interval"
    send_to_routes         = true
    stale_channel_flush_ms = 5570369.95
    streamtags = [
      "..."
    ]
    type = "exec"
  }
  input_file = {
    breaker_rulesets = [
      "..."
    ]
    check_file_mod_time = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    delete_files = true
    depth        = 3.86
    description  = "...my_description..."
    disabled     = true
    environment  = "...my_environment..."
    filenames = [
      "..."
    ]
    force_text                    = true
    hash_len                      = 10.09
    id                            = "...my_id..."
    idle_timeout                  = 9.52
    include_unidentifiable_binary = true
    interval                      = 9.61
    max_age_dur                   = "...my_max_age_dur..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    mode     = "auto"
    path     = "...my_path..."
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 4.04
      compress         = "none"
      max_buffer_size  = 43.27
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = false
    stale_channel_flush_ms = 38380929.61
    streamtags = [
      "..."
    ]
    suppress_missing_path_errors = true
    tail_only                    = true
    type                         = "file"
  }
  input_firehose = {
    activity_log_sample_rate = 3.21
    auth_tokens = [
      "..."
    ]
    capture_headers = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = false
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 14.72
    max_active_req          = 1.19
    max_requests_per_socket = 4
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 5325.25
    pq = {
      commit_frequency = 4.21
      compress         = "none"
      max_buffer_size  = 46.88
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled      = false
    request_timeout = 2.37
    send_to_routes  = true
    socket_timeout  = 1.73
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "firehose"
  }
  input_google_pubsub = {
    concurrency = 63.78
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    create_subscription = true
    create_topic        = false
    description         = "...my_description..."
    disabled            = false
    environment         = "...my_environment..."
    google_auth_method  = "secret"
    id                  = "...my_id..."
    max_backlog         = 8.38
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    ordered_delivery = false
    pipeline         = "...my_pipeline..."
    pq = {
      commit_frequency = 5.28
      compress         = "gzip"
      max_buffer_size  = 48.44
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled                  = false
    region                      = "...my_region..."
    request_timeout             = 10003.11
    secret                      = "...my_secret..."
    send_to_routes              = true
    service_account_credentials = "...my_service_account_credentials..."
    streamtags = [
      "..."
    ]
    subscription_name = "...my_subscription_name..."
    topic_name        = "...my_topic_name..."
    type              = "google_pubsub"
  }
  input_grafana = {
    # ...
  }
  input_http = {
    activity_log_sample_rate = 3.13
    auth_tokens = [
      "..."
    ]
    auth_tokens_ext = [
      {
        description = "...my_description..."
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        token = "...my_token..."
      }
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    cribl_api               = "...my_cribl_api..."
    description             = "...my_description..."
    disabled                = false
    elastic_api             = "...my_elastic_api..."
    enable_health_check     = true
    enable_proxy_header     = true
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 511.57
    max_active_req          = 2.9
    max_requests_per_socket = 8
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 42584.61
    pq = {
      commit_frequency = 2.24
      compress         = "none"
      max_buffer_size  = 45.05
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled      = false
    request_timeout = 1.82
    send_to_routes  = true
    socket_timeout  = 9.84
    splunk_hec_acks = false
    splunk_hec_api  = "...my_splunk_hec_api..."
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "http"
  }
  input_http_raw = {
    activity_log_sample_rate = 9.97
    allowed_methods = [
      "..."
    ]
    allowed_paths = [
      "..."
    ]
    auth_tokens = [
      "..."
    ]
    auth_tokens_ext = [
      {
        description = "...my_description..."
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        token = "...my_token..."
      }
    ]
    breaker_rulesets = [
      "..."
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = false
    enable_proxy_header     = true
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 315.47
    max_active_req          = 0.01
    max_requests_per_socket = 9
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 840.25
    pq = {
      commit_frequency = 1.36
      compress         = "gzip"
      max_buffer_size  = 46.1
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = true
    request_timeout        = 4.25
    send_to_routes         = false
    socket_timeout         = 3.94
    stale_channel_flush_ms = 17162456.39
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "http_raw"
  }
  input_journal_files = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    current_boot = true
    description  = "...my_description..."
    disabled     = false
    environment  = "...my_environment..."
    id           = "...my_id..."
    interval     = 8.87
    journals = [
      "..."
    ]
    max_age_dur = "...my_max_age_dur..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    path     = "...my_path..."
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 5.71
      compress         = "gzip"
      max_buffer_size  = 44.06
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = false
    rules = [
      {
        description = "...my_description..."
        filter      = "...my_filter..."
      }
    ]
    send_to_routes = true
    streamtags = [
      "..."
    ]
    type = "journal_files"
  }
  input_kafka = {
    authentication_timeout = 2612811.22
    auto_commit_interval   = 1980921.22
    auto_commit_threshold  = 6446.97
    backoff_rate           = 16.34
    brokers = [
      "..."
    ]
    connection_timeout = 2340192.34
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    environment        = "...my_environment..."
    from_beginning     = true
    group_id           = "...my_group_id..."
    heartbeat_interval = 2828908.53
    id                 = "...my_id..."
    initial_backoff    = 100808.83
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout  = 45495.99
      disabled            = true
      max_retries         = 97.38
      request_timeout     = 8439.49
      schema_registry_url = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = true
        max_version         = "TLSv1.1"
        min_version         = "TLSv1.1"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = true
        servername          = "...my_servername..."
      }
    }
    max_back_off            = 134271.89
    max_bytes               = 808758177.23
    max_bytes_per_partition = 8916916.4
    max_retries             = 6.98
    max_socket_errors       = 84.62
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 6.49
      compress         = "gzip"
      max_buffer_size  = 44.7
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                 = false
    reauthentication_threshold = 1231299.62
    rebalance_timeout          = 3002374.89
    request_timeout            = 2736560.86
    sasl = {
      disabled  = true
      mechanism = "scram-sha-512"
    }
    send_to_routes  = false
    session_timeout = 1308790.91
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      servername          = "...my_servername..."
    }
    topics = [
      "..."
    ]
    type = "kafka"
  }
  input_kinesis = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    avoid_duplicates          = true
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description              = "...my_description..."
    disabled                 = false
    duration_seconds         = 26349.88
    enable_assume_role       = true
    endpoint                 = "...my_endpoint..."
    environment              = "...my_environment..."
    get_records_limit        = 7947.22
    get_records_limit_total  = 20004.93
    id                       = "...my_id..."
    load_balancing_algorithm = "ConsistentHashing"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    payload_format = "line"
    pipeline       = "...my_pipeline..."
    pq = {
      commit_frequency = 7.77
      compress         = "none"
      max_buffer_size  = 43.27
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = true
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = false
    send_to_routes      = true
    service_interval    = 2.32
    shard_expr          = "...my_shard_expr..."
    shard_iterator_type = "LATEST"
    signature_version   = "v2"
    stream_name         = "...my_stream_name..."
    streamtags = [
      "..."
    ]
    type                  = "kinesis"
    verify_kpl_check_sums = true
  }
  input_kube_events = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description = "...my_description..."
    disabled    = false
    environment = "...my_environment..."
    id          = "...my_id..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 9.09
      compress         = "none"
      max_buffer_size  = 50.31
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = true
    rules = [
      {
        description = "...my_description..."
        filter      = "...my_filter..."
      }
    ]
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "kube_events"
  }
  input_kube_logs = {
    breaker_rulesets = [
      "..."
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disabled              = false
    enable_load_balancing = true
    environment           = "...my_environment..."
    id                    = "...my_id..."
    interval              = 7.59
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 8.47
      compress         = "gzip"
      max_buffer_size  = 45.63
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = true
    rules = [
      {
        description = "...my_description..."
        filter      = "...my_filter..."
      }
    ]
    send_to_routes         = true
    stale_channel_flush_ms = 25164066.67
    streamtags = [
      "..."
    ]
    timestamps = false
    type       = "kube_logs"
  }
  input_kube_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description = "...my_description..."
    disabled    = false
    environment = "...my_environment..."
    id          = "...my_id..."
    interval    = 2.3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "...my_dest_path..."
      enable        = false
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 4.98
      compress         = "none"
      max_buffer_size  = 44.58
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = true
    rules = [
      {
        description = "...my_description..."
        filter      = "...my_filter..."
      }
    ]
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "kube_metrics"
  }
  input_loki = {
    activity_log_sample_rate = 4.33
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "credentialsSecret"
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "...my_credentials_secret..."
    description             = "...my_description..."
    disabled                = true
    enable_health_check     = false
    enable_proxy_header     = true
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 23.89
    login_url               = "...my_login_url..."
    loki_api                = "...my_loki_api..."
    max_active_req          = 0.4
    max_requests_per_socket = 3
    metadata = [
      {
        name  = "...my_name..."
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
    password = "...my_password..."
    pipeline = "...my_pipeline..."
    port     = 27049.39
    pq = {
      commit_frequency = 2.01
      compress         = "none"
      max_buffer_size  = 43.95
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled        = false
    request_timeout   = 1.12
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_to_routes    = true
    socket_timeout    = 0.37
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 148445.66
    type                 = "loki"
    username             = "...my_username..."
  }
  input_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "...my_description..."
    disabled            = false
    enable_proxy_header = false
    environment         = "...my_environment..."
    host                = "...my_host..."
    id                  = "...my_id..."
    ip_whitelist_regex  = "...my_ip_whitelist_regex..."
    max_buffer_size     = 3.81
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 2.12
      compress         = "none"
      max_buffer_size  = 48.63
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled     = false
    send_to_routes = true
    streamtags = [
      "..."
    ]
    tcp_port = 1867.09
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type                   = "metrics"
    udp_port               = 22216.87
    udp_socket_rx_buf_size = 2294508638.44
  }
  input_model_driven_telemetry = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description    = "...my_description..."
    disabled       = false
    environment    = "...my_environment..."
    host           = "...my_host..."
    id             = "...my_id..."
    max_active_cxn = 7.7
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 27675.62
    pq = {
      commit_frequency = 4.63
      compress         = "none"
      max_buffer_size  = 43.89
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = false
    send_to_routes      = false
    shutdown_timeout_ms = 8.84
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1"
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "model_driven_telemetry"
  }
  input_msk = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    authentication_timeout    = 995810.31
    auto_commit_interval      = 481775.08
    auto_commit_threshold     = 1754.67
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    backoff_rate              = 3.39
    brokers = [
      "..."
    ]
    connection_timeout = 306053.8
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    duration_seconds   = 40267.7
    enable_assume_role = true
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    from_beginning     = true
    group_id           = "...my_group_id..."
    heartbeat_interval = 792040.64
    id                 = "...my_id..."
    initial_backoff    = 595314.32
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout  = 17882.16
      disabled            = true
      max_retries         = 1.41
      request_timeout     = 30438.96
      schema_registry_url = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = true
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = false
        servername          = "...my_servername..."
      }
    }
    max_back_off            = 140556.78
    max_bytes               = 658141657.62
    max_bytes_per_partition = 6885585.41
    max_retries             = 25.31
    max_socket_errors       = 4.21
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 1.63
      compress         = "none"
      max_buffer_size  = 50.29
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                 = false
    reauthentication_threshold = 52754.92
    rebalance_timeout          = 371225.77
    region                     = "...my_region..."
    reject_unauthorized        = false
    request_timeout            = 1734742.82
    reuse_connections          = true
    send_to_routes             = true
    session_timeout            = 656767.91
    signature_version          = "v4"
    streamtags = [
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
    topics = [
      "..."
    ]
    type = "msk"
  }
  input_netflow = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "...my_description..."
    disabled            = true
    enable_pass_through = false
    environment         = "...my_environment..."
    host                = "...my_host..."
    id                  = "...my_id..."
    ip_allowlist_regex  = "...my_ip_allowlist_regex..."
    ip_denylist_regex   = "...my_ip_denylist_regex..."
    ipfix_enabled       = false
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 55240.21
    pq = {
      commit_frequency = 5.86
      compress         = "gzip"
      max_buffer_size  = 46.86
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled     = true
    send_to_routes = true
    streamtags = [
      "..."
    ]
    template_cache_minutes = 405.69
    type                   = "netflow"
    udp_socket_rx_buf_size = 2487129551.16
    v5_enabled             = true
    v9_enabled             = true
  }
  input_office365_mgmt = {
    app_id        = "...my_app_id..."
    auth_type     = "manual"
    client_secret = "...my_client_secret..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        content_type = "...my_content_type..."
        description  = "...my_description..."
        enabled      = true
        interval     = 3.17
        log_level    = "error"
      }
    ]
    description             = "...my_description..."
    disabled                = false
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    ingestion_lag           = 3871.05
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 14.09
    max_missed_keep_alives  = 10.26
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline  = "...my_pipeline..."
    plan_type = "enterprise_gcc"
    pq = {
      commit_frequency = 5.53
      compress         = "gzip"
      max_buffer_size  = 47.16
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled           = true
    publisher_identifier = "...my_publisher_identifier..."
    retry_rules = {
      codes = [
        2.04
      ]
      enable_header         = false
      interval              = 6479.66
      limit                 = 9.17
      multiplier            = 7.61
      retry_connect_reset   = true
      retry_connect_timeout = false
      type                  = "static"
    }
    send_to_routes = false
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 2134.41
    ttl         = "...my_ttl..."
    type        = "office365_mgmt"
  }
  input_office365_msg_trace = {
    auth_type = "manual"
    cert_options = {
      cert_path        = "...my_cert_path..."
      certificate_name = "...my_certificate_name..."
      passphrase       = "...my_passphrase..."
      priv_key_path    = "...my_priv_key_path..."
    }
    client_id     = "...my_client_id..."
    client_secret = "...my_client_secret..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "...my_credentials_secret..."
    description             = "...my_description..."
    disable_time_filter     = false
    disabled                = true
    end_date                = "...my_end_date..."
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    interval                = 32.97
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 10.76
    log_level               = "debug"
    max_missed_keep_alives  = 10.83
    max_task_reschedule     = 3.22
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    password  = "...my_password..."
    pipeline  = "...my_pipeline..."
    plan_type = "enterprise_gcc"
    pq = {
      commit_frequency = 9.56
      compress         = "gzip"
      max_buffer_size  = 43.6
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled               = false
    reschedule_dropped_tasks = true
    resource                 = "...my_resource..."
    retry_rules = {
      codes = [
        3.69
      ]
      enable_header         = false
      interval              = 16223.83
      limit                 = 3.89
      multiplier            = 17.73
      retry_connect_reset   = false
      retry_connect_timeout = false
      type                  = "static"
    }
    send_to_routes = true
    start_date     = "...my_start_date..."
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 47.49
    ttl         = "...my_ttl..."
    type        = "office365_msg_trace"
    url         = "...my_url..."
    username    = "...my_username..."
  }
  input_office365_service = {
    app_id        = "...my_app_id..."
    auth_type     = "manual"
    client_secret = "...my_client_secret..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        content_type = "...my_content_type..."
        description  = "...my_description..."
        enabled      = false
        interval     = 52.06
        log_level    = "warn"
      }
    ]
    description             = "...my_description..."
    disabled                = true
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 15.91
    max_missed_keep_alives  = 4.02
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline  = "...my_pipeline..."
    plan_type = "dod"
    pq = {
      commit_frequency = 9.79
      compress         = "none"
      max_buffer_size  = 50.31
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = false
    retry_rules = {
      codes = [
        6.75
      ]
      enable_header         = true
      interval              = 13866.62
      limit                 = 16.32
      multiplier            = 15.13
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "backoff"
    }
    send_to_routes = false
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 1674.45
    ttl         = "...my_ttl..."
    type        = "office365_service"
  }
  input_open_telemetry = {
    activity_log_sample_rate = "{ \"see\": \"documentation\" }"
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "basic"
    capture_headers          = "{ \"see\": \"documentation\" }"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "...my_credentials_secret..."
    description             = "...my_description..."
    disabled                = true
    enable_health_check     = false
    enable_proxy_header     = "{ \"see\": \"documentation\" }"
    environment             = "...my_environment..."
    extract_logs            = false
    extract_metrics         = false
    extract_spans           = false
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 383.32
    login_url               = "...my_login_url..."
    max_active_cxn          = 6.15
    max_active_req          = 8.73
    max_requests_per_socket = 8
    metadata = [
      {
        name  = "...my_name..."
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
    otlp_version = "0.10.0"
    password     = "...my_password..."
    pipeline     = "...my_pipeline..."
    port         = 18349.27
    pq = {
      commit_frequency = 6.78
      compress         = "gzip"
      max_buffer_size  = 49.56
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled        = false
    protocol          = "http"
    request_timeout   = 3.78
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_to_routes    = true
    socket_timeout    = 1.89
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 200638.53
    type                 = "open_telemetry"
    username             = "...my_username..."
  }
  input_prometheus = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    auth_type                 = "manual"
    aws_authentication_method = "auto"
    aws_secret_key            = "...my_aws_secret_key..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret = "...my_credentials_secret..."
    description        = "...my_description..."
    dimension_list = [
    ]
    disabled                = false
    discovery_type          = "ec2"
    duration_seconds        = 4736.3
    enable_assume_role      = true
    endpoint                = "...my_endpoint..."
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    interval                = 51.04
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 17.36
    log_level               = "debug"
    max_missed_keep_alives  = 4
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    name_list = [
      "..."
    ]
    password = "...my_password..."
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 9.53
      compress         = "gzip"
      max_buffer_size  = 47.62
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = false
    record_type         = "A"
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = false
    scrape_path         = "...my_scrape_path..."
    scrape_port         = 61082.45
    scrape_protocol     = "http"
    search_filter = [
      {
        name = "...my_name..."
        values = [
        ]
      }
    ]
    send_to_routes    = false
    signature_version = "v2"
    streamtags = [
      "..."
    ]
    target_list = [
      "..."
    ]
    ttl           = "...my_ttl..."
    type          = "prometheus"
    use_public_ip = true
    username      = "...my_username..."
  }
  input_prometheus_rw = {
    activity_log_sample_rate = 9.93
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "none"
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "...my_credentials_secret..."
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 106.5
    login_url               = "...my_login_url..."
    max_active_req          = 6.12
    max_requests_per_socket = 3
    metadata = [
      {
        name  = "...my_name..."
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
    password = "...my_password..."
    pipeline = "...my_pipeline..."
    port     = 16545.91
    pq = {
      commit_frequency = 4.86
      compress         = "gzip"
      max_buffer_size  = 50.83
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled        = true
    prometheus_api    = "...my_prometheus_api..."
    request_timeout   = 8.19
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_to_routes    = false
    socket_timeout    = 8.67
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 293746.84
    type                 = "prometheus_rw"
    username             = "...my_username..."
  }
  input_raw_udp = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    environment        = "...my_environment..."
    host               = "...my_host..."
    id                 = "...my_id..."
    ingest_raw_bytes   = false
    ip_whitelist_regex = "...my_ip_whitelist_regex..."
    max_buffer_size    = 0.95
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 21267.65
    pq = {
      commit_frequency = 3.81
      compress         = "none"
      max_buffer_size  = 43
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = false
    single_msg_udp_packets = false
    streamtags = [
      "..."
    ]
    type                   = "raw_udp"
    udp_socket_rx_buf_size = 2133363074.21
  }
  input_s3 = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "manual"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    breaker_rulesets = [
      "..."
    ]
    checkpointing = {
      enabled = false
      retries = 61.98
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = false
    duration_seconds       = 34856.34
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "...my_encoding..."
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_messages           = 8.96
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 69.55
    parquet_chunk_download_timeout = 499.25
    parquet_chunk_size_mb          = 17.08
    pipeline                       = "...my_pipeline..."
    poll_timeout                   = 18.2
    pq = {
      commit_frequency = 6.29
      compress         = "gzip"
      max_buffer_size  = 47.48
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = false
    }
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = false
    reuse_connections      = false
    send_to_routes         = false
    signature_version      = "v4"
    skip_on_error          = false
    socket_timeout         = 29705.33
    stale_channel_flush_ms = 18587902.44
    streamtags = [
      "..."
    ]
    tag_after_processing = true
    type                 = "s3"
    visibility_timeout   = 16414.35
  }
  input_s3_inventory = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    breaker_rulesets = [
      "..."
    ]
    checkpointing = {
      enabled = true
      retries = 66.46
    }
    checksum_suffix = "...my_checksum_suffix..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = true
    duration_seconds       = 32507.46
    enable_assume_role     = false
    enable_sqs_assume_role = true
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_manifest_size_kb   = 8
    max_messages           = 9.2
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 88.47
    parquet_chunk_download_timeout = 589.64
    parquet_chunk_size_mb          = 57.64
    pipeline                       = "...my_pipeline..."
    poll_timeout                   = 5.74
    pq = {
      commit_frequency = 10.03
      compress         = "gzip"
      max_buffer_size  = 48.72
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = true
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = false
    }
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = true
    reuse_connections      = false
    send_to_routes         = false
    signature_version      = "v4"
    skip_on_error          = false
    socket_timeout         = 9696.9
    stale_channel_flush_ms = 26256833.64
    streamtags = [
      "..."
    ]
    tag_after_processing     = "false"
    type                     = "s3_inventory"
    validate_inventory_files = false
    visibility_timeout       = 26709.77
  }
  input_security_lake = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    breaker_rulesets = [
      "..."
    ]
    checkpointing = {
      enabled = false
      retries = 86.82
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = true
    duration_seconds       = 2917.34
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "...my_encoding..."
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_messages           = 8.06
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 6.73
    parquet_chunk_download_timeout = 2933.58
    parquet_chunk_size_mb          = 41.45
    pipeline                       = "...my_pipeline..."
    poll_timeout                   = 14.18
    pq = {
      commit_frequency = 9.12
      compress         = "none"
      max_buffer_size  = 46.18
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = false
    }
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = false
    reuse_connections      = true
    send_to_routes         = false
    signature_version      = "v2"
    skip_on_error          = false
    socket_timeout         = 26799.94
    stale_channel_flush_ms = 29469341.99
    streamtags = [
      "..."
    ]
    tag_after_processing = "false"
    type                 = "security_lake"
    visibility_timeout   = 16801.61
  }
  input_snmp = {
    best_effort_parsing = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    environment        = "...my_environment..."
    host               = "...my_host..."
    id                 = "...my_id..."
    ip_whitelist_regex = "...my_ip_whitelist_regex..."
    max_buffer_size    = 5.96
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 16809.21
    pq = {
      commit_frequency = 7.41
      compress         = "none"
      max_buffer_size  = 50.16
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled     = true
    send_to_routes = true
    snmp_v3_auth = {
      allow_unmatched_trap = true
      v3_auth_enabled      = true
      v3_users = [
        {
          auth_key      = "{ \"see\": \"documentation\" }"
          auth_protocol = "sha224"
          name          = "...my_name..."
          priv_protocol = "{ \"see\": \"documentation\" }"
        }
      ]
    }
    streamtags = [
      "..."
    ]
    type                   = "snmp"
    udp_socket_rx_buf_size = 2286477561.21
    varbinds_with_types    = true
  }
  input_splunk = {
    auth_tokens = [
      {
        description = "...my_description..."
        token       = "...my_token..."
      }
    ]
    breaker_rulesets = [
      "..."
    ]
    compress = "always"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "...my_description..."
    disabled            = true
    drop_control_fields = true
    enable_proxy_header = false
    environment         = "...my_environment..."
    extract_metrics     = false
    host                = "...my_host..."
    id                  = "...my_id..."
    ip_whitelist_regex  = "...my_ip_whitelist_regex..."
    max_active_cxn      = 8.13
    max_s2_sversion     = "v3"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 46299.32
    pq = {
      commit_frequency = 7.19
      compress         = "gzip"
      max_buffer_size  = 42.3
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 4.19
    socket_idle_timeout    = 8.04
    socket_max_lifespan    = 9.33
    stale_channel_flush_ms = 24735904.72
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type             = "splunk"
    use_fwd_timezone = false
  }
  input_splunk_hec = {
    access_control_allow_headers = [
    ]
    access_control_allow_origin = [
    ]
    activity_log_sample_rate = 5.25
    allowed_indexes = [
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
        ]
        auth_type   = "manual"
        description = "...my_description..."
        enabled     = false
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        token        = "{ \"see\": \"documentation\" }"
        token_secret = "{ \"see\": \"documentation\" }"
      }
    ]
    breaker_rulesets = [
      "..."
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = false
    drop_control_fields     = false
    emit_token_metrics      = true
    enable_health_check     = "{ \"see\": \"documentation\" }"
    enable_proxy_header     = true
    environment             = "...my_environment..."
    extract_metrics         = true
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 128.31
    max_active_req          = 9.63
    max_requests_per_socket = 5
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 58949.4
    pq = {
      commit_frequency = 10.82
      compress         = "gzip"
      max_buffer_size  = 45.5
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = true
    request_timeout        = 4.49
    send_to_routes         = false
    socket_timeout         = 4.38
    splunk_hec_acks        = true
    splunk_hec_api         = "...my_splunk_hec_api..."
    stale_channel_flush_ms = 42109803.47
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type             = "splunk_hec"
    use_fwd_timezone = false
  }
  input_splunk_search = {
    auth_header_expr = "...my_auth_header_expr..."
    auth_type        = "none"
    breaker_rulesets = [
      "..."
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret = "...my_credentials_secret..."
    cron_schedule      = "...my_cron_schedule..."
    description        = "...my_description..."
    disabled           = false
    earliest           = "...my_earliest..."
    encoding           = "...my_encoding..."
    endpoint           = "...my_endpoint..."
    endpoint_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    endpoint_params = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 13.53
    latest                  = "...my_latest..."
    log_level               = "debug"
    login_url               = "...my_login_url..."
    max_missed_keep_alives  = 3.83
    metadata = [
      {
        name  = "...my_name..."
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
    output_mode = "json"
    password    = "...my_password..."
    pipeline    = "...my_pipeline..."
    pq = {
      commit_frequency = 5.71
      compress         = "gzip"
      max_buffer_size  = 44.06
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = false
    reject_unauthorized = false
    request_timeout     = 1851.98
    retry_rules = {
      codes = [
        6.05
      ]
      enable_header         = false
      interval              = 7602.68
      limit                 = 18.39
      multiplier            = 16.13
      retry_connect_reset   = false
      retry_connect_timeout = false
      type                  = "static"
    }
    search                 = "...my_search..."
    search_head            = "...my_search_head..."
    secret                 = "...my_secret..."
    secret_param_name      = "...my_secret_param_name..."
    send_to_routes         = true
    stale_channel_flush_ms = 2870608.4
    streamtags = [
      "..."
    ]
    text_secret          = "...my_text_secret..."
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 219176.87
    ttl                  = "...my_ttl..."
    type                 = "splunk_search"
    use_round_robin_dns  = false
    username             = "...my_username..."
  }
  input_sqs = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_account_id            = "...my_aws_account_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    create_queue       = false
    description        = "...my_description..."
    disabled           = true
    duration_seconds   = 23574.97
    enable_assume_role = false
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    id                 = "...my_id..."
    max_messages       = 1.93
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers = 60.31
    pipeline      = "...my_pipeline..."
    poll_timeout  = 7.51
    pq = {
      commit_frequency = 1.7
      compress         = "none"
      max_buffer_size  = 47.38
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = false
    queue_name          = "...my_queue_name..."
    queue_type          = "fifo"
    region              = "...my_region..."
    reject_unauthorized = true
    reuse_connections   = false
    send_to_routes      = false
    signature_version   = "v2"
    streamtags = [
      "..."
    ]
    type               = "sqs"
    visibility_timeout = 42454.5
  }
  input_syslog = {
    input_syslog_syslog2 = {
      allow_non_standard_app_name = true
      connections = [
        {
          output   = "...my_output..."
          pipeline = "...my_pipeline..."
        }
      ]
      description                          = "...my_description..."
      disabled                             = false
      enable_enhanced_proxy_header_parsing = true
      enable_load_balancing                = true
      enable_proxy_header                  = true
      environment                          = "...my_environment..."
      host                                 = "...my_host..."
      id                                   = "...my_id..."
      infer_framing                        = true
      ip_whitelist_regex                   = "...my_ip_whitelist_regex..."
      keep_fields_list = [
      ]
      max_active_cxn  = 7.64
      max_buffer_size = 3.56
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      octet_counting = true
      pipeline       = "...my_pipeline..."
      pq = {
        commit_frequency = 6.96
        compress         = "none"
        max_buffer_size  = 45.65
        max_file_size    = "...my_max_file_size..."
        max_size         = "...my_max_size..."
        mode             = "always"
        path             = "...my_path..."
      }
      pq_enabled             = true
      send_to_routes         = false
      single_msg_udp_packets = false
      socket_ending_max_wait = 3.23
      socket_idle_timeout    = 2.03
      socket_max_lifespan    = 7.57
      streamtags = [
        "..."
      ]
      strictly_infer_octet_counting = true
      tcp_port                      = 62520.25
      timestamp_timezone            = "...my_timestamp_timezone..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        common_name_regex   = "{ \"see\": \"documentation\" }"
        disabled            = true
        max_version         = "TLSv1.2"
        min_version         = "TLSv1.2"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = "{ \"see\": \"documentation\" }"
        request_cert        = true
      }
      type                   = "syslog"
      udp_port               = 59552.67
      udp_socket_rx_buf_size = 635379092.96
    }
  }
  input_system_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    container = {
      all_containers = false
      detail         = true
      docker_socket = [
        "..."
      ]
      docker_timeout = 1.14
      filters = [
        {
          expr = "...my_expr..."
        }
      ]
      mode       = "custom"
      per_device = false
    }
    description = "...my_description..."
    disabled    = true
    environment = "...my_environment..."
    host = {
      custom = {
        cpu = {
          detail  = false
          mode    = "disabled"
          per_cpu = false
          time    = false
        }
        disk = {
          detail = true
          devices = [
            "..."
          ]
          fstypes = [
            "..."
          ]
          mode = "disabled"
          mountpoints = [
            "..."
          ]
          per_device = true
        }
        memory = {
          detail = false
          mode   = "disabled"
        }
        network = {
          detail = false
          devices = [
            "..."
          ]
          mode          = "all"
          per_interface = true
        }
        system = {
          mode      = "disabled"
          processes = false
        }
      }
      mode = "custom"
    }
    id       = "...my_id..."
    interval = 10.53
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      dest_path     = "...my_dest_path..."
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 7.36
      compress         = "gzip"
      max_buffer_size  = 48
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = true
    process = {
      sets = [
        {
          filter           = "...my_filter..."
          include_children = true
          name             = "...my_name..."
        }
      ]
    }
    send_to_routes = true
    streamtags = [
      "..."
    ]
    type = "system_metrics"
  }
  input_system_state = {
    collectors = {
      disk = {
        enable = true
      }
      dns = {
        enable = true
      }
      firewall = {
        enable = false
      }
      hostsfile = {
        enable = true
      }
      interfaces = {
        enable = false
      }
      login_users = {
        enable = true
      }
      metadata = {
        enable = false
      }
      ports = {
        enable = false
      }
      routes = {
        enable = true
      }
      services = {
        enable = true
      }
      user = {
        enable = true
      }
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disable_native_module = true
    disabled              = true
    environment           = "...my_environment..."
    id                    = "...my_id..."
    interval              = 9.73
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      dest_path     = "...my_dest_path..."
      enable        = false
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 1.29
      compress         = "none"
      max_buffer_size  = 48.67
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled     = false
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "system_state"
  }
  input_tcp = {
    auth_type = "manual"
    breaker_rulesets = [
      "..."
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "...my_description..."
    disabled            = false
    enable_header       = false
    enable_proxy_header = false
    environment         = "...my_environment..."
    host                = "...my_host..."
    id                  = "...my_id..."
    ip_whitelist_regex  = "...my_ip_whitelist_regex..."
    max_active_cxn      = 5.5
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 47674.49
    pq = {
      commit_frequency = 1.66
      compress         = "none"
      max_buffer_size  = 51.41
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = false
    }
    send_to_routes         = false
    socket_ending_max_wait = 6.18
    socket_idle_timeout    = 0.36
    socket_max_lifespan    = 5.19
    stale_channel_flush_ms = 8063309.13
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "tcp"
  }
  input_tcpjson = {
    auth_token = "...my_auth_token..."
    auth_type  = "manual"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disabled              = true
    enable_load_balancing = true
    enable_proxy_header   = true
    environment           = "...my_environment..."
    host                  = "...my_host..."
    id                    = "...my_id..."
    ip_whitelist_regex    = "...my_ip_whitelist_regex..."
    max_active_cxn        = 6.9
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 44098.72
    pq = {
      commit_frequency = 5.4
      compress         = "gzip"
      max_buffer_size  = 43.12
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = true
    send_to_routes         = true
    socket_ending_max_wait = 5.57
    socket_idle_timeout    = 8.47
    socket_max_lifespan    = 3.88
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "tcpjson"
  }
  input_wef = {
    allow_machine_id_mismatch = false
    auth_method               = "kerberos"
    ca_fingerprint            = "...my_ca_fingerprint..."
    capture_headers           = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description              = "...my_description..."
    disabled                 = true
    enable_health_check      = true
    enable_proxy_header      = false
    environment              = "...my_environment..."
    host                     = "...my_host..."
    id                       = "...my_id..."
    ip_allowlist_regex       = "...my_ip_allowlist_regex..."
    ip_denylist_regex        = "...my_ip_denylist_regex..."
    keep_alive_timeout       = 375.33
    keytab                   = "...my_keytab..."
    log_fingerprint_mismatch = true
    max_active_req           = 7.34
    max_requests_per_socket  = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 55639.99
    pq = {
      commit_frequency = 1.25
      compress         = "gzip"
      max_buffer_size  = 42.29
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled     = true
    principal      = "...my_principal..."
    send_to_routes = false
    socket_timeout = 7.92
    streamtags = [
      "..."
    ]
    subscriptions = [
      {
        batch_timeout      = 5.23
        compress           = false
        content_format     = "RenderedText"
        heartbeat_interval = 5.28
        id                 = "...my_id..."
        locale             = "...my_locale..."
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        query_selector       = "xml"
        read_existing_events = false
        send_bookmarks       = false
        subscription_name    = "...my_subscription_name..."
        targets = [
          "..."
        ]
        version = "...my_version..."
      }
    ]
    tls = {
      ca_path               = "...my_ca_path..."
      cert_path             = "...my_cert_path..."
      certificate_name      = "...my_certificate_name..."
      common_name_regex     = "...my_common_name_regex..."
      disabled              = false
      keytab                = "{ \"see\": \"documentation\" }"
      max_version           = "TLSv1"
      min_version           = "TLSv1.1"
      ocsp_check            = false
      ocsp_check_fail_close = true
      passphrase            = "...my_passphrase..."
      principal             = "{ \"see\": \"documentation\" }"
      priv_key_path         = "...my_priv_key_path..."
      reject_unauthorized   = false
      request_cert          = false
    }
    type = "wef"
  }
  input_windows_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disable_native_module = true
    disabled              = false
    environment           = "...my_environment..."
    host = {
      custom = {
        cpu = {
          detail  = true
          mode    = "disabled"
          per_cpu = true
          time    = false
        }
        disk = {
          mode       = "disabled"
          per_volume = true
          volumes = [
            "..."
          ]
        }
        memory = {
          detail = true
          mode   = "disabled"
        }
        network = {
          detail = true
          devices = [
            "..."
          ]
          mode          = "custom"
          per_interface = false
        }
        system = {
          detail = true
          mode   = "custom"
        }
      }
      mode = "disabled"
    }
    id       = "...my_id..."
    interval = 10.39
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      dest_path     = "...my_dest_path..."
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 8.52
      compress         = "none"
      max_buffer_size  = 48.56
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = false
    process = {
      sets = [
        {
          filter           = "...my_filter..."
          include_children = true
          name             = "...my_name..."
        }
      ]
    }
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "windows_metrics"
  }
  input_win_event_logs = {
    batch_size = 9.54
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disable_native_module = true
    disabled              = true
    environment           = "...my_environment..."
    event_format          = "json"
    id                    = "...my_id..."
    interval              = 2.52
    log_names = [
      "..."
    ]
    max_event_bytes = 82427367.56
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 8.61
      compress         = "gzip"
      max_buffer_size  = 49.04
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled     = true
    read_mode      = "oldest"
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "win_event_logs"
  }
  input_wiz = {
    auth_audience_override = "...my_auth_audience_override..."
    auth_type              = "manual"
    auth_url               = "...my_auth_url..."
    client_id              = "...my_client_id..."
    client_secret          = "...my_client_secret..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        content_description = "...my_content_description..."
        content_type        = "...my_content_type..."
        enabled             = true
      }
    ]
    description             = "...my_description..."
    disabled                = false
    endpoint                = "...my_endpoint..."
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    keep_alive_time         = 14
    max_missed_keep_alives  = 9.47
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 6.65
      compress         = "gzip"
      max_buffer_size  = 45.65
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 1833.22
    retry_rules = {
      codes = [
        4.2
      ]
      enable_header         = false
      interval              = 4792.78
      limit                 = 16.43
      multiplier            = 19.92
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "backoff"
    }
    send_to_routes = false
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    ttl         = "...my_ttl..."
    type        = "wiz"
  }
  input_zscaler_hec = {
    access_control_allow_headers = [
    ]
    access_control_allow_origin = [
    ]
    activity_log_sample_rate = 8.55
    allowed_indexes = [
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
        ]
        auth_type   = "manual"
        description = "...my_description..."
        enabled     = true
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        token        = "{ \"see\": \"documentation\" }"
        token_secret = "{ \"see\": \"documentation\" }"
      }
    ]
    capture_headers = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = true
    emit_token_metrics      = true
    enable_health_check     = "{ \"see\": \"documentation\" }"
    enable_proxy_header     = true
    environment             = "...my_environment..."
    hec_acks                = false
    hec_api                 = "...my_hec_api..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 162.87
    max_active_req          = 0.89
    max_requests_per_socket = 5
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 3454.48
    pq = {
      commit_frequency = 5.32
      compress         = "none"
      max_buffer_size  = 49.79
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 2.41
    send_to_routes  = false
    socket_timeout  = 5.66
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "zscaler_hec"
  }
}