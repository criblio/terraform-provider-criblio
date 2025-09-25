resource "criblio_pack_source" "my_packsource" {
  group_id = "...my_group_id..."
  id       = "...my_id..."
  input_appscope = {
    auth_token = "...my_auth_token..."
    auth_type  = "secret"
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
    enable_proxy_header = true
    enable_unix_path    = false
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
    max_active_cxn     = 1.37
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "...my_dest_path..."
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    port     = 28161.43
    pq = {
      commit_frequency = 5.85
      compress         = "gzip"
      max_buffer_size  = 47.94
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 1.65
    socket_idle_timeout    = 1.4
    socket_max_lifespan    = 9.58
    stale_channel_flush_ms = 41730255.43
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
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
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
    disabled        = false
    endpoint_suffix = "...my_endpoint_suffix..."
    environment     = "...my_environment..."
    file_filter     = "...my_file_filter..."
    id              = "...my_id..."
    max_messages    = 14.25
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 13.36
    parquet_chunk_download_timeout = 3208.34
    parquet_chunk_size_mb          = 75.95
    pipeline                       = "...my_pipeline..."
    pq = {
      commit_frequency = 5.88
      compress         = "none"
      max_buffer_size  = 48.23
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    queue_name             = "...my_queue_name..."
    send_to_routes         = true
    service_period_secs    = 6.29
    skip_on_error          = true
    stale_channel_flush_ms = 19720963.38
    storage_account_name   = "...my_storage_account_name..."
    streamtags = [
      "..."
    ]
    tenant_id          = "...my_tenant_id..."
    text_secret        = "...my_text_secret..."
    type               = "azure_blob"
    visibility_timeout = 231727.65
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
    disabled    = true
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
      commit_frequency = 3.08
      compress         = "none"
      max_buffer_size  = 42.51
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
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
    send_to_routes         = true
    stale_channel_flush_ms = 21269562.73
    streamtags = [
      "..."
    ]
    throttle_rate_per_sec = "...my_throttle_rate_per_sec..."
    type                  = "collection"
  }
  input_confluent_cloud = {
    authentication_timeout = 2937234.14
    auto_commit_interval   = 659880.71
    auto_commit_threshold  = 5879.65
    backoff_rate           = 11.16
    brokers = [
      "..."
    ]
    connection_timeout = 2914293.03
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = false
    environment        = "...my_environment..."
    from_beginning     = false
    group_id           = "...my_group_id..."
    heartbeat_interval = 3166951.53
    id                 = "...my_id..."
    initial_backoff    = 419145.49
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = true
      }
      connection_timeout  = 3962.74
      disabled            = false
      max_retries         = 46.18
      request_timeout     = 31813.85
      schema_registry_url = "...my_schema_registry_url..."
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
    }
    max_back_off            = 134471.78
    max_bytes               = 393370490.01
    max_bytes_per_partition = 9387222.15
    max_retries             = 91.99
    max_socket_errors       = 53.98
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 7.13
      compress         = "none"
      max_buffer_size  = 49.55
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                 = true
    reauthentication_threshold = 1675705.38
    rebalance_timeout          = 422499.3
    request_timeout            = 3058268.79
    sasl = {
      disabled  = true
      mechanism = "scram-sha-256"
    }
    send_to_routes  = false
    session_timeout = 3060629.99
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      disabled            = true
      max_version         = "TLSv1"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
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
      commit_frequency = 9.9
      compress         = "gzip"
      max_buffer_size  = 43.66
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
    type = "cribl"
  }
  input_cribl_http = {
    activity_log_sample_rate = 3.47
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
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 585.48
    max_active_req          = 4.61
    max_requests_per_socket = 1
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 7127.76
    pq = {
      commit_frequency = 5.87
      compress         = "none"
      max_buffer_size  = 49.68
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 2.4
    send_to_routes  = true
    socket_timeout  = 5.36
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "cribl_http"
  }
  input_cribl_lake_http = {
    activity_log_sample_rate = 1.57
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
    disabled                = true
    enable_health_check     = false
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 447.43
    max_active_req          = 5.44
    max_requests_per_socket = 9
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 26909.28
    pq = {
      commit_frequency = 5.05
      compress         = "none"
      max_buffer_size  = 47.49
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled      = false
    request_timeout = 5.06
    send_to_routes  = false
    socket_timeout  = 0.26
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
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
    full_fidelity = false
    id            = "...my_id..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 1.52
      compress         = "none"
      max_buffer_size  = 50.12
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
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
    disabled              = true
    enable_load_balancing = false
    enable_proxy_header   = true
    environment           = "...my_environment..."
    host                  = "...my_host..."
    id                    = "...my_id..."
    max_active_cxn        = 1.27
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 52135.7
    pq = {
      commit_frequency = 4.41
      compress         = "none"
      max_buffer_size  = 42.78
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = false
    socket_ending_max_wait = 7.99
    socket_idle_timeout    = 6.78
    socket_max_lifespan    = 3.77
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
      retries = 14.14
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = true
    duration_seconds       = 3017.54
    enable_assume_role     = false
    enable_sqs_assume_role = true
    encoding               = "...my_encoding..."
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_messages           = 6
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers = 34.07
    pipeline      = "...my_pipeline..."
    poll_timeout  = 14.09
    pq = {
      commit_frequency = 4.32
      compress         = "gzip"
      max_buffer_size  = 42.07
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
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = false
    reuse_connections      = false
    send_to_routes         = false
    signature_version      = "v2"
    skip_on_error          = true
    socket_timeout         = 22004.07
    stale_channel_flush_ms = 14290632.13
    streamtags = [
      "..."
    ]
    tag_after_processing = "true"
    type                 = "crowdstrike"
    visibility_timeout   = 662.93
  }
  input_datadog_agent = {
    activity_log_sample_rate = 2.98
    capture_headers          = true
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
    extract_metrics         = false
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 57.71
    max_active_req          = 5.65
    max_requests_per_socket = 2
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 14790.78
    pq = {
      commit_frequency = 10.88
      compress         = "gzip"
      max_buffer_size  = 45.59
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = true
    proxy_mode = {
      enabled             = false
      reject_unauthorized = true
    }
    request_timeout = 3.2
    send_to_routes  = true
    socket_timeout  = 6.91
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
      commit_frequency = 8.51
      compress         = "none"
      max_buffer_size  = 44.45
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled = false
    samples = [
      {
        events_per_sec = 7.37
        sample         = "...my_sample..."
      }
    ]
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "datagen"
  }
  input_edge_prometheus = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    auth_type                 = "manual"
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
      "..."
    ]
    disabled           = false
    discovery_type     = "k8s-pods"
    duration_seconds   = 4666.4
    enable_assume_role = true
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    id                 = "...my_id..."
    interval           = 3.21
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
      enable        = false
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
      commit_frequency = 8.56
      compress         = "gzip"
      max_buffer_size  = 48.47
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled           = true
    record_type          = "AAAA"
    region               = "...my_region..."
    reject_unauthorized  = true
    reuse_connections    = false
    scrape_path          = "...my_scrape_path..."
    scrape_path_expr     = "...my_scrape_path_expr..."
    scrape_port          = 37487.6
    scrape_port_expr     = "...my_scrape_port_expr..."
    scrape_protocol      = "http"
    scrape_protocol_expr = "...my_scrape_protocol_expr..."
    search_filter = [
      {
        name = "...my_name..."
        values = [
          "..."
        ]
      }
    ]
    send_to_routes    = true
    signature_version = "v4"
    streamtags = [
      "..."
    ]
    targets = [
      {
        host     = "...my_host..."
        path     = "...my_path..."
        port     = 9226.3
        protocol = "http"
      }
    ]
    timeout       = 11361.91
    type          = "edge_prometheus"
    use_public_ip = true
    username      = "...my_username..."
  }
  input_elastic = {
    activity_log_sample_rate = 3.46
    api_version              = "custom"
    auth_tokens = [
      "..."
    ]
    auth_type       = "credentialsSecret"
    capture_headers = true
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
    keep_alive_timeout      = 519.14
    max_active_req          = 9.4
    max_requests_per_socket = 7
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    password = "...my_password..."
    pipeline = "...my_pipeline..."
    port     = 63358.75
    pq = {
      commit_frequency = 7.29
      compress         = "gzip"
      max_buffer_size  = 46.69
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = true
    proxy_mode = {
      auth_type           = "secret"
      enabled             = true
      reject_unauthorized = true
      remove_headers = [
        "..."
      ]
      timeout_sec = 3207862514268799.5
      url         = "...my_url..."
    }
    request_timeout = 1
    send_to_routes  = false
    socket_timeout  = 8.18
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
      request_cert        = false
    }
    type     = "elastic"
    username = "...my_username..."
  }
  input_eventhub = {
    authentication_timeout = 2856176.19
    auto_commit_interval   = 1319807.47
    auto_commit_threshold  = 9173.68
    backoff_rate           = 10.91
    brokers = [
      "..."
    ]
    connection_timeout = 3004423.35
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = true
    environment             = "...my_environment..."
    from_beginning          = false
    group_id                = "...my_group_id..."
    heartbeat_interval      = 1973630.58
    id                      = "...my_id..."
    initial_backoff         = 435093.88
    max_back_off            = 131350
    max_bytes               = 228637439.84
    max_bytes_per_partition = 6765683.06
    max_retries             = 90.41
    max_socket_errors       = 90.71
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    minimize_duplicates = true
    pipeline            = "...my_pipeline..."
    pq = {
      commit_frequency = 6.67
      compress         = "none"
      max_buffer_size  = 50.4
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                 = true
    reauthentication_threshold = 408790.3
    rebalance_timeout          = 2000665.85
    request_timeout            = 77186.25
    sasl = {
      disabled  = false
      mechanism = "oauthbearer"
    }
    send_to_routes  = true
    session_timeout = 198201.45
    streamtags = [
      "..."
    ]
    tls = {
      disabled            = true
      reject_unauthorized = false
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
    interval      = 2.17
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 4.61
      compress         = "none"
      max_buffer_size  = 44.51
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    retries                = 2.16
    schedule_type          = "interval"
    send_to_routes         = true
    stale_channel_flush_ms = 3982835.68
    streamtags = [
      "..."
    ]
    type = "exec"
  }
  input_file = {
    breaker_rulesets = [
      "..."
    ]
    check_file_mod_time = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    delete_files = false
    depth        = 0.56
    description  = "...my_description..."
    disabled     = true
    environment  = "...my_environment..."
    filenames = [
      "..."
    ]
    force_text                    = false
    hash_len                      = 9.59
    id                            = "...my_id..."
    idle_timeout                  = 2.16
    include_unidentifiable_binary = false
    interval                      = 10.53
    max_age_dur                   = "...my_max_age_dur..."
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    mode     = "manual"
    path     = "...my_path..."
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 10.43
      compress         = "none"
      max_buffer_size  = 43.03
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = true
    send_to_routes         = false
    stale_channel_flush_ms = 7766576.84
    streamtags = [
      "..."
    ]
    suppress_missing_path_errors = false
    tail_only                    = true
    type                         = "file"
  }
  input_firehose = {
    activity_log_sample_rate = 1.98
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
    disabled                = true
    enable_health_check     = false
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 287.57
    max_active_req          = 9.93
    max_requests_per_socket = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 53600.14
    pq = {
      commit_frequency = 9.1
      compress         = "none"
      max_buffer_size  = 50.79
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 2.73
    send_to_routes  = false
    socket_timeout  = 5.78
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
    type = "firehose"
  }
  input_google_pubsub = {
    concurrency = 49.72
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    create_subscription = false
    create_topic        = true
    description         = "...my_description..."
    disabled            = true
    environment         = "...my_environment..."
    google_auth_method  = "auto"
    id                  = "...my_id..."
    max_backlog         = 4.09
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    ordered_delivery = false
    pipeline         = "...my_pipeline..."
    pq = {
      commit_frequency = 8.56
      compress         = "none"
      max_buffer_size  = 51.86
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                  = false
    region                      = "...my_region..."
    request_timeout             = 10003
    secret                      = "...my_secret..."
    send_to_routes              = false
    service_account_credentials = "...my_service_account_credentials..."
    streamtags = [
      "..."
    ]
    subscription_name = "...my_subscription_name..."
    topic_name        = "...my_topic_name..."
    type              = "google_pubsub"
  }
  input_grafana = {
    input_grafana_grafana1 = {
      activity_log_sample_rate = 4.98
      capture_headers          = true
      connections = [
        {
          output   = "...my_output..."
          pipeline = "...my_pipeline..."
        }
      ]
      description         = "...my_description..."
      disabled            = true
      enable_health_check = false
      enable_proxy_header = true
      environment         = "...my_environment..."
      host                = "...my_host..."
      id                  = "...my_id..."
      ip_allowlist_regex  = "...my_ip_allowlist_regex..."
      ip_denylist_regex   = "...my_ip_denylist_regex..."
      keep_alive_timeout  = 78.37
      loki_api            = "...my_loki_api..."
      loki_auth = {
        auth_header_expr   = "...my_auth_header_expr..."
        auth_type          = "textSecret"
        credentials_secret = "...my_credentials_secret..."
        login_url          = "...my_login_url..."
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
        password             = "...my_password..."
        secret               = "...my_secret..."
        secret_param_name    = "...my_secret_param_name..."
        text_secret          = "...my_text_secret..."
        token                = "...my_token..."
        token_attribute_name = "...my_token_attribute_name..."
        token_timeout_secs   = 153201.44
        username             = "...my_username..."
      }
      max_active_req          = 2.9
      max_requests_per_socket = 8
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      pipeline = "...my_pipeline..."
      port     = 50121.26
      pq = {
        commit_frequency = 7.49
        compress         = "gzip"
        max_buffer_size  = 46.8
        max_file_size    = "...my_max_file_size..."
        max_size         = "...my_max_size..."
        mode             = "smart"
        path             = "...my_path..."
      }
      pq_enabled     = true
      prometheus_api = "...my_prometheus_api..."
      prometheus_auth = {
        auth_header_expr   = "...my_auth_header_expr..."
        auth_type          = "credentialsSecret"
        credentials_secret = "...my_credentials_secret..."
        login_url          = "...my_login_url..."
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
        password             = "...my_password..."
        secret               = "...my_secret..."
        secret_param_name    = "...my_secret_param_name..."
        text_secret          = "...my_text_secret..."
        token                = "...my_token..."
        token_attribute_name = "...my_token_attribute_name..."
        token_timeout_secs   = 226188.99
        username             = "...my_username..."
      }
      request_timeout = 9.92
      send_to_routes  = true
      socket_timeout  = 0.76
      streamtags = [
        "..."
      ]
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        common_name_regex   = "{ \"see\": \"documentation\" }"
        disabled            = false
        max_version         = "TLSv1.2"
        min_version         = "TLSv1.3"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = "{ \"see\": \"documentation\" }"
        request_cert        = false
      }
      type = "grafana"
    }
    input_grafana_grafana2 = {
      activity_log_sample_rate = 8.39
      capture_headers          = true
      connections = [
        {
          output   = "...my_output..."
          pipeline = "...my_pipeline..."
        }
      ]
      description         = "...my_description..."
      disabled            = true
      enable_health_check = false
      enable_proxy_header = false
      environment         = "...my_environment..."
      host                = "...my_host..."
      id                  = "...my_id..."
      ip_allowlist_regex  = "...my_ip_allowlist_regex..."
      ip_denylist_regex   = "...my_ip_denylist_regex..."
      keep_alive_timeout  = 481.89
      loki_api            = "...my_loki_api..."
      loki_auth = {
        auth_header_expr   = "...my_auth_header_expr..."
        auth_type          = "oauth"
        credentials_secret = "...my_credentials_secret..."
        login_url          = "...my_login_url..."
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
        password             = "...my_password..."
        secret               = "...my_secret..."
        secret_param_name    = "...my_secret_param_name..."
        text_secret          = "...my_text_secret..."
        token                = "...my_token..."
        token_attribute_name = "...my_token_attribute_name..."
        token_timeout_secs   = 71960.83
        username             = "...my_username..."
      }
      max_active_req          = 0.92
      max_requests_per_socket = 0
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      pipeline = "...my_pipeline..."
      port     = 22887.21
      pq = {
        commit_frequency = 7.32
        compress         = "none"
        max_buffer_size  = 48.08
        max_file_size    = "...my_max_file_size..."
        max_size         = "...my_max_size..."
        mode             = "smart"
        path             = "...my_path..."
      }
      pq_enabled     = false
      prometheus_api = "...my_prometheus_api..."
      prometheus_auth = {
        auth_header_expr   = "...my_auth_header_expr..."
        auth_type          = "none"
        credentials_secret = "...my_credentials_secret..."
        login_url          = "...my_login_url..."
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
        password             = "...my_password..."
        secret               = "...my_secret..."
        secret_param_name    = "...my_secret_param_name..."
        text_secret          = "...my_text_secret..."
        token                = "...my_token..."
        token_attribute_name = "...my_token_attribute_name..."
        token_timeout_secs   = 72520.14
        username             = "...my_username..."
      }
      request_timeout = 1.09
      send_to_routes  = true
      socket_timeout  = 6.68
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
        min_version         = "TLSv1.2"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = "{ \"see\": \"documentation\" }"
        request_cert        = true
      }
      type = "grafana"
    }
  }
  input_http = {
    activity_log_sample_rate = 4.81
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
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 253.68
    max_active_req          = 2.7
    max_requests_per_socket = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 45701.82
    pq = {
      commit_frequency = 5.79
      compress         = "gzip"
      max_buffer_size  = 51.73
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 4.91
    send_to_routes  = true
    socket_timeout  = 6.77
    splunk_hec_acks = true
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
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "http"
  }
  input_http_raw = {
    activity_log_sample_rate = 9.27
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
    capture_headers = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "...my_description..."
    disabled                = true
    enable_health_check     = false
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 108.11
    max_active_req          = 2.4
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 36582.64
    pq = {
      commit_frequency = 6.83
      compress         = "none"
      max_buffer_size  = 44.79
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = true
    request_timeout        = 4.77
    send_to_routes         = false
    socket_timeout         = 5.89
    stale_channel_flush_ms = 15090330.96
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
    type = "http_raw"
  }
  input_journal_files = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    current_boot = false
    description  = "...my_description..."
    disabled     = false
    environment  = "...my_environment..."
    id           = "...my_id..."
    interval     = 6.22
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
      commit_frequency = 6.19
      compress         = "gzip"
      max_buffer_size  = 50.96
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
    send_to_routes = false
    streamtags = [
      "..."
    ]
    type = "journal_files"
  }
  input_kafka = {
    authentication_timeout = 2652921.72
    auto_commit_interval   = 3478352.95
    auto_commit_threshold  = 6316.75
    backoff_rate           = 17.82
    brokers = [
      "..."
    ]
    connection_timeout = 2290406.81
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    environment        = "...my_environment..."
    from_beginning     = false
    group_id           = "...my_group_id..."
    heartbeat_interval = 2743984.37
    id                 = "...my_id..."
    initial_backoff    = 196819.38
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = false
      }
      connection_timeout  = 37034.38
      disabled            = true
      max_retries         = 67.48
      request_timeout     = 47589.14
      schema_registry_url = "...my_schema_registry_url..."
      tls = {
        ca_path             = "...my_ca_path..."
        cert_path           = "...my_cert_path..."
        certificate_name    = "...my_certificate_name..."
        disabled            = false
        max_version         = "TLSv1.3"
        min_version         = "TLSv1"
        passphrase          = "...my_passphrase..."
        priv_key_path       = "...my_priv_key_path..."
        reject_unauthorized = false
        servername          = "...my_servername..."
      }
    }
    max_back_off            = 164531.11
    max_bytes               = 248872825.96
    max_bytes_per_partition = 9344616.68
    max_retries             = 34.78
    max_socket_errors       = 68.55
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 2.69
      compress         = "none"
      max_buffer_size  = 44.12
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled                 = true
    reauthentication_threshold = 1606181.62
    rebalance_timeout          = 2248660.36
    request_timeout            = 2721416.33
    sasl = {
      disabled  = false
      mechanism = "plain"
    }
    send_to_routes  = false
    session_timeout = 362328.16
    streamtags = [
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
    topics = [
      "..."
    ]
    type = "kafka"
  }
  input_kinesis = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    avoid_duplicates          = false
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
    duration_seconds         = 29491.68
    enable_assume_role       = true
    endpoint                 = "...my_endpoint..."
    environment              = "...my_environment..."
    get_records_limit        = 6889.36
    get_records_limit_total  = 20005.23
    id                       = "...my_id..."
    load_balancing_algorithm = "RoundRobin"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    payload_format = "ndjson"
    pipeline       = "...my_pipeline..."
    pq = {
      commit_frequency = 5.87
      compress         = "gzip"
      max_buffer_size  = 49.57
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = true
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = true
    send_to_routes      = true
    service_interval    = 4.81
    shard_expr          = "...my_shard_expr..."
    shard_iterator_type = "TRIM_HORIZON"
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
      commit_frequency = 1.98
      compress         = "gzip"
      max_buffer_size  = 51.69
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
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
    disabled              = true
    enable_load_balancing = false
    environment           = "...my_environment..."
    id                    = "...my_id..."
    interval              = 2.81
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      enable        = false
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 1.41
      compress         = "gzip"
      max_buffer_size  = 43.46
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
    send_to_routes         = false
    stale_channel_flush_ms = 9527792.95
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
    interval    = 9.33
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
      commit_frequency = 7.59
      compress         = "none"
      max_buffer_size  = 48.01
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
    type = "kube_metrics"
  }
  input_loki = {
    activity_log_sample_rate = 2.82
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "none"
    capture_headers          = false
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "...my_credentials_secret..."
    description             = "...my_description..."
    disabled                = false
    enable_health_check     = false
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 258.94
    login_url               = "...my_login_url..."
    loki_api                = "...my_loki_api..."
    max_active_req          = 2.08
    max_requests_per_socket = 5
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
    port     = 4252.49
    pq = {
      commit_frequency = 7.66
      compress         = "none"
      max_buffer_size  = 48.78
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled        = true
    request_timeout   = 4.82
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_to_routes    = false
    socket_timeout    = 8.29
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
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 207624.84
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
    max_buffer_size     = 3.84
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 7.61
      compress         = "none"
      max_buffer_size  = 49.38
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
    tcp_port = 11197.86
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type                   = "metrics"
    udp_port               = 35923.63
    udp_socket_rx_buf_size = 2905545357.22
  }
  input_model_driven_telemetry = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description    = "...my_description..."
    disabled       = true
    environment    = "...my_environment..."
    host           = "...my_host..."
    id             = "...my_id..."
    max_active_cxn = 3.31
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 23604.97
    pq = {
      commit_frequency = 1.29
      compress         = "none"
      max_buffer_size  = 49.43
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled          = false
    send_to_routes      = true
    shutdown_timeout_ms = 8.1
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
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "model_driven_telemetry"
  }
  input_msk = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    authentication_timeout    = 2337527.72
    auto_commit_interval      = 2808196.34
    auto_commit_threshold     = 3614.49
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "...my_aws_secret_key..."
    backoff_rate              = 2.05
    brokers = [
      "..."
    ]
    connection_timeout = 675079.56
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "...my_description..."
    disabled           = true
    duration_seconds   = 32002.83
    enable_assume_role = false
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    from_beginning     = false
    group_id           = "...my_group_id..."
    heartbeat_interval = 2373003.56
    id                 = "...my_id..."
    initial_backoff    = 363083.16
    kafka_schema_registry = {
      auth = {
        credentials_secret = "...my_credentials_secret..."
        disabled           = true
      }
      connection_timeout  = 32986.14
      disabled            = false
      max_retries         = 80.94
      request_timeout     = 18578.18
      schema_registry_url = "...my_schema_registry_url..."
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
    }
    max_back_off            = 137439.93
    max_bytes               = 520222029.03
    max_bytes_per_partition = 2007697
    max_retries             = 53.77
    max_socket_errors       = 58.9
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 10.02
      compress         = "none"
      max_buffer_size  = 49.16
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled                 = true
    reauthentication_threshold = 681492.4
    rebalance_timeout          = 2684656.55
    region                     = "...my_region..."
    reject_unauthorized        = true
    request_timeout            = 811027.57
    reuse_connections          = true
    send_to_routes             = true
    session_timeout            = 1819829.34
    signature_version          = "v2"
    streamtags = [
      "..."
    ]
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
    disabled            = false
    enable_pass_through = true
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
    port     = 39962.65
    pq = {
      commit_frequency = 2.83
      compress         = "gzip"
      max_buffer_size  = 51.9
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
    template_cache_minutes = 269.34
    type                   = "netflow"
    udp_socket_rx_buf_size = 423502369.93
    v5_enabled             = true
    v9_enabled             = true
  }
  input_office365_mgmt = {
    app_id        = "...my_app_id..."
    auth_type     = "secret"
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
        interval     = 39.73
        log_level    = "warn"
      }
    ]
    description             = "...my_description..."
    disabled                = true
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    ingestion_lag           = 5131.61
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 16.08
    max_missed_keep_alives  = 5.52
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline  = "...my_pipeline..."
    plan_type = "dod"
    pq = {
      commit_frequency = 8.47
      compress         = "none"
      max_buffer_size  = 45.85
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled           = true
    publisher_identifier = "...my_publisher_identifier..."
    retry_rules = {
      codes = [
        5.84
      ]
      enable_header         = true
      interval              = 16425.8
      limit                 = 3.42
      multiplier            = 6.3
      retry_connect_reset   = false
      retry_connect_timeout = false
      type                  = "none"
    }
    send_to_routes = true
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 1210.68
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
    ignore_group_jobs_limit = false
    interval                = 1.29
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 17.57
    log_level               = "silly"
    max_missed_keep_alives  = 7.45
    max_task_reschedule     = 5.6
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    password  = "...my_password..."
    pipeline  = "...my_pipeline..."
    plan_type = "gcc_high"
    pq = {
      commit_frequency = 9.74
      compress         = "none"
      max_buffer_size  = 46.49
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled               = false
    reschedule_dropped_tasks = false
    resource                 = "...my_resource..."
    retry_rules = {
      codes = [
        2.83
      ]
      enable_header         = false
      interval              = 17938.05
      limit                 = 0.54
      multiplier            = 6.89
      retry_connect_reset   = false
      retry_connect_timeout = false
      type                  = "none"
    }
    send_to_routes = true
    start_date     = "...my_start_date..."
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 2229.18
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
        interval     = 59.2
        log_level    = "info"
      }
    ]
    description             = "...my_description..."
    disabled                = true
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 10.67
    max_missed_keep_alives  = 10.35
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline  = "...my_pipeline..."
    plan_type = "gcc_high"
    pq = {
      commit_frequency = 8.14
      compress         = "none"
      max_buffer_size  = 42.92
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = false
    retry_rules = {
      codes = [
        3.08
      ]
      enable_header         = true
      interval              = 9947.1
      limit                 = 17.64
      multiplier            = 14.96
      retry_connect_reset   = false
      retry_connect_timeout = false
      type                  = "none"
    }
    send_to_routes = true
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 1175.99
    ttl         = "...my_ttl..."
    type        = "office365_service"
  }
  input_open_telemetry = {
    activity_log_sample_rate = "{ \"see\": \"documentation\" }"
    auth_header_expr         = "...my_auth_header_expr..."
    auth_type                = "textSecret"
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
    keep_alive_timeout      = 125.91
    login_url               = "...my_login_url..."
    max_active_cxn          = 5.36
    max_active_req          = 7.62
    max_requests_per_socket = 4
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
    otlp_version = "1.3.1"
    password     = "...my_password..."
    pipeline     = "...my_pipeline..."
    port         = 24313.31
    pq = {
      commit_frequency = 10.7
      compress         = "none"
      max_buffer_size  = 48.88
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled        = false
    protocol          = "grpc"
    request_timeout   = 0.09
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_to_routes    = true
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
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 292620.71
    type                 = "open_telemetry"
    username             = "...my_username..."
  }
  input_prometheus = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    auth_type                 = "manual"
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
      "..."
    ]
    disabled                = true
    discovery_type          = "ec2"
    duration_seconds        = 40827.31
    enable_assume_role      = true
    endpoint                = "...my_endpoint..."
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    interval                = 36.27
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 17.9
    log_level               = "error"
    max_missed_keep_alives  = 8.4
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
      commit_frequency = 4.85
      compress         = "none"
      max_buffer_size  = 46.51
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled          = false
    record_type         = "SRV"
    region              = "...my_region..."
    reject_unauthorized = false
    reuse_connections   = false
    scrape_path         = "...my_scrape_path..."
    scrape_port         = 32569.25
    scrape_protocol     = "https"
    search_filter = [
      {
        name = "...my_name..."
        values = [
          "..."
        ]
      }
    ]
    send_to_routes    = false
    signature_version = "v4"
    streamtags = [
      "..."
    ]
    target_list = [
      "..."
    ]
    ttl           = "...my_ttl..."
    type          = "prometheus"
    use_public_ip = false
    username      = "...my_username..."
  }
  input_prometheus_rw = {
    activity_log_sample_rate = 8.82
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
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 475.59
    login_url               = "...my_login_url..."
    max_active_req          = 2.98
    max_requests_per_socket = 1
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
    port     = 5688.01
    pq = {
      commit_frequency = 2.2
      compress         = "gzip"
      max_buffer_size  = 45.62
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled        = false
    prometheus_api    = "...my_prometheus_api..."
    request_timeout   = 5.95
    secret            = "...my_secret..."
    secret_param_name = "...my_secret_param_name..."
    send_to_routes    = true
    socket_timeout    = 6.86
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
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 176398.89
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
    disabled           = false
    environment        = "...my_environment..."
    host               = "...my_host..."
    id                 = "...my_id..."
    ingest_raw_bytes   = true
    ip_whitelist_regex = "...my_ip_whitelist_regex..."
    max_buffer_size    = 2.36
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 41946.99
    pq = {
      commit_frequency = 10.54
      compress         = "gzip"
      max_buffer_size  = 47
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = true
    send_to_routes         = true
    single_msg_udp_packets = true
    streamtags = [
      "..."
    ]
    type                   = "raw_udp"
    udp_socket_rx_buf_size = 2207313577.18
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
      retries = 93.53
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = true
    duration_seconds       = 7333.11
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "...my_encoding..."
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_messages           = 9.83
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 10.24
    parquet_chunk_download_timeout = 630.07
    parquet_chunk_size_mb          = 47.7
    pipeline                       = "...my_pipeline..."
    poll_timeout                   = 17.37
    pq = {
      commit_frequency = 3.11
      compress         = "gzip"
      max_buffer_size  = 48.81
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
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
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v4"
    skip_on_error          = true
    socket_timeout         = 1535.43
    stale_channel_flush_ms = 4285303.66
    streamtags = [
      "..."
    ]
    tag_after_processing = false
    type                 = "s3"
    visibility_timeout   = 25858.09
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
      enabled = false
      retries = 92.48
    }
    checksum_suffix = "...my_checksum_suffix..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = false
    duration_seconds       = 32892.34
    enable_assume_role     = true
    enable_sqs_assume_role = true
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_manifest_size_kb   = 8
    max_messages           = 9.96
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 31.6
    parquet_chunk_download_timeout = 969.93
    parquet_chunk_size_mb          = 91.52
    pipeline                       = "...my_pipeline..."
    poll_timeout                   = 8.14
    pq = {
      commit_frequency = 5.96
      compress         = "none"
      max_buffer_size  = 46.48
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
    processed_tag_key      = "...my_processed_tag_key..."
    processed_tag_value    = "...my_processed_tag_value..."
    queue_name             = "...my_queue_name..."
    region                 = "...my_region..."
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = false
    signature_version      = "v4"
    skip_on_error          = false
    socket_timeout         = 31558.99
    stale_channel_flush_ms = 41327873.14
    streamtags = [
      "..."
    ]
    tag_after_processing     = "false"
    type                     = "s3_inventory"
    validate_inventory_files = true
    visibility_timeout       = 8804.03
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
      enabled = true
      retries = 9.49
    }
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "...my_description..."
    disabled               = false
    duration_seconds       = 21288.08
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "...my_encoding..."
    endpoint               = "...my_endpoint..."
    environment            = "...my_environment..."
    file_filter            = "...my_file_filter..."
    id                     = "...my_id..."
    max_messages           = 3.31
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 48.52
    parquet_chunk_download_timeout = 1197.61
    parquet_chunk_size_mb          = 74.21
    pipeline                       = "...my_pipeline..."
    poll_timeout                   = 11.12
    pq = {
      commit_frequency = 3.79
      compress         = "none"
      max_buffer_size  = 42.44
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
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = false
    signature_version      = "v4"
    skip_on_error          = false
    socket_timeout         = 32105.81
    stale_channel_flush_ms = 27690302.35
    streamtags = [
      "..."
    ]
    tag_after_processing = "false"
    type                 = "security_lake"
    visibility_timeout   = 8868.64
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
    max_buffer_size    = 7.7
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 37824.81
    pq = {
      commit_frequency = 2.53
      compress         = "gzip"
      max_buffer_size  = 44.33
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled     = true
    send_to_routes = false
    snmp_v3_auth = {
      allow_unmatched_trap = true
      v3_auth_enabled      = false
      v3_users = [
        {
          auth_key      = "{ \"see\": \"documentation\" }"
          auth_protocol = "sha256"
          name          = "...my_name..."
          priv_protocol = "{ \"see\": \"documentation\" }"
        }
      ]
    }
    streamtags = [
      "..."
    ]
    type                   = "snmp"
    udp_socket_rx_buf_size = 918882948.42
    varbinds_with_types    = false
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
    compress = "disabled"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "...my_description..."
    disabled            = false
    drop_control_fields = false
    enable_proxy_header = true
    environment         = "...my_environment..."
    extract_metrics     = true
    host                = "...my_host..."
    id                  = "...my_id..."
    ip_whitelist_regex  = "...my_ip_whitelist_regex..."
    max_active_cxn      = 3.36
    max_s2_sversion     = "v4"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 30244.35
    pq = {
      commit_frequency = 9.29
      compress         = "none"
      max_buffer_size  = 46.84
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 1.4
    socket_idle_timeout    = 0.89
    socket_max_lifespan    = 1.32
    stale_channel_flush_ms = 4831222.63
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type             = "splunk"
    use_fwd_timezone = false
  }
  input_splunk_hec = {
    access_control_allow_headers = [
      "..."
    ]
    access_control_allow_origin = [
      "..."
    ]
    activity_log_sample_rate = 7.44
    allowed_indexes = [
      "..."
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
          "..."
        ]
        auth_type   = "secret"
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
    breaker_rulesets = [
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
    disabled                = true
    drop_control_fields     = false
    emit_token_metrics      = true
    enable_health_check     = "{ \"see\": \"documentation\" }"
    enable_proxy_header     = false
    environment             = "...my_environment..."
    extract_metrics         = true
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 480.42
    max_active_req          = 7.27
    max_requests_per_socket = 8
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 45730.24
    pq = {
      commit_frequency = 7.55
      compress         = "none"
      max_buffer_size  = 48.46
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled             = false
    request_timeout        = 7.68
    send_to_routes         = false
    socket_timeout         = 0.26
    splunk_hec_acks        = true
    splunk_hec_api         = "...my_splunk_hec_api..."
    stale_channel_flush_ms = 1355570.83
    streamtags = [
      "..."
    ]
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
      request_cert        = true
    }
    type             = "splunk_hec"
    use_fwd_timezone = false
  }
  input_splunk_search = {
    auth_header_expr = "...my_auth_header_expr..."
    auth_type        = "oauth"
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
    disabled           = true
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
    ignore_group_jobs_limit = false
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 12.35
    latest                  = "...my_latest..."
    log_level               = "info"
    login_url               = "...my_login_url..."
    max_missed_keep_alives  = 4.86
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
    output_mode = "csv"
    password    = "...my_password..."
    pipeline    = "...my_pipeline..."
    pq = {
      commit_frequency = 4.05
      compress         = "none"
      max_buffer_size  = 45.31
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled          = false
    reject_unauthorized = true
    request_timeout     = 586.16
    retry_rules = {
      codes = [
        9.39
      ]
      enable_header         = true
      interval              = 6165.82
      limit                 = 15.85
      multiplier            = 10.91
      retry_connect_reset   = false
      retry_connect_timeout = true
      type                  = "backoff"
    }
    search                 = "...my_search..."
    search_head            = "...my_search_head..."
    secret                 = "...my_secret..."
    secret_param_name      = "...my_secret_param_name..."
    send_to_routes         = false
    stale_channel_flush_ms = 27372045.99
    streamtags = [
      "..."
    ]
    text_secret          = "...my_text_secret..."
    token                = "...my_token..."
    token_attribute_name = "...my_token_attribute_name..."
    token_timeout_secs   = 245461.79
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
    duration_seconds   = 34384.07
    enable_assume_role = false
    endpoint           = "...my_endpoint..."
    environment        = "...my_environment..."
    id                 = "...my_id..."
    max_messages       = 5.23
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers = 45.55
    pipeline      = "...my_pipeline..."
    poll_timeout  = 19.8
    pq = {
      commit_frequency = 4.03
      compress         = "none"
      max_buffer_size  = 46.43
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled          = false
    queue_name          = "...my_queue_name..."
    queue_type          = "fifo"
    region              = "...my_region..."
    reject_unauthorized = true
    reuse_connections   = false
    send_to_routes      = true
    signature_version   = "v4"
    streamtags = [
      "..."
    ]
    type               = "sqs"
    visibility_timeout = 16398.99
  }
  input_syslog = {
    input_syslog_syslog1 = {
      allow_non_standard_app_name = false
      connections = [
        {
          output   = "...my_output..."
          pipeline = "...my_pipeline..."
        }
      ]
      description                          = "...my_description..."
      disabled                             = false
      enable_enhanced_proxy_header_parsing = true
      enable_load_balancing                = false
      enable_proxy_header                  = false
      environment                          = "...my_environment..."
      host                                 = "...my_host..."
      id                                   = "...my_id..."
      infer_framing                        = false
      ip_whitelist_regex                   = "...my_ip_whitelist_regex..."
      keep_fields_list = [
        "..."
      ]
      max_active_cxn  = 9.27
      max_buffer_size = 4.44
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      octet_counting = true
      pipeline       = "...my_pipeline..."
      pq = {
        commit_frequency = 2.22
        compress         = "none"
        max_buffer_size  = 46.68
        max_file_size    = "...my_max_file_size..."
        max_size         = "...my_max_size..."
        mode             = "smart"
        path             = "...my_path..."
      }
      pq_enabled             = false
      send_to_routes         = true
      single_msg_udp_packets = false
      socket_ending_max_wait = 4.48
      socket_idle_timeout    = 8.62
      socket_max_lifespan    = 8.74
      streamtags = [
        "..."
      ]
      strictly_infer_octet_counting = true
      tcp_port                      = 54010.44
      timestamp_timezone            = "...my_timestamp_timezone..."
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
      type                   = "syslog"
      udp_port               = 45541.65
      udp_socket_rx_buf_size = 2227773755.28
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
      detail         = false
      docker_socket = [
        "..."
      ]
      docker_timeout = 5.24
      filters = [
        {
          expr = "...my_expr..."
        }
      ]
      mode       = "basic"
      per_device = true
    }
    description = "...my_description..."
    disabled    = false
    environment = "...my_environment..."
    host = {
      custom = {
        cpu = {
          detail  = false
          mode    = "disabled"
          per_cpu = true
          time    = false
        }
        disk = {
          detail = false
          devices = [
            "..."
          ]
          fstypes = [
            "..."
          ]
          mode = "basic"
          mountpoints = [
            "..."
          ]
          per_device = true
        }
        memory = {
          detail = true
          mode   = "disabled"
        }
        network = {
          detail = false
          devices = [
            "..."
          ]
          mode          = "basic"
          per_interface = false
        }
        system = {
          mode      = "basic"
          processes = true
        }
      }
      mode = "basic"
    }
    id       = "...my_id..."
    interval = 9.24
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "...my_dest_path..."
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 6.83
      compress         = "none"
      max_buffer_size  = 48.26
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
        enable = false
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
        enable = true
      }
      routes = {
        enable = true
      }
      services = {
        enable = false
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
    disabled              = false
    environment           = "...my_environment..."
    id                    = "...my_id..."
    interval              = 5.13
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "...my_dest_path..."
      enable        = true
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 3.31
      compress         = "none"
      max_buffer_size  = 50.09
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
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
    disabled            = true
    enable_header       = false
    enable_proxy_header = true
    environment         = "...my_environment..."
    host                = "...my_host..."
    id                  = "...my_id..."
    ip_whitelist_regex  = "...my_ip_whitelist_regex..."
    max_active_cxn      = 8.41
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 7592.47
    pq = {
      commit_frequency = 7.04
      compress         = "none"
      max_buffer_size  = 51.06
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
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
    socket_ending_max_wait = 8.02
    socket_idle_timeout    = 5.47
    socket_max_lifespan    = 9.38
    stale_channel_flush_ms = 8313205.9
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
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "tcp"
  }
  input_tcpjson = {
    auth_token = "...my_auth_token..."
    auth_type  = "secret"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disabled              = true
    enable_load_balancing = true
    enable_proxy_header   = false
    environment           = "...my_environment..."
    host                  = "...my_host..."
    id                    = "...my_id..."
    ip_whitelist_regex    = "...my_ip_whitelist_regex..."
    max_active_cxn        = 1.87
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 21110.14
    pq = {
      commit_frequency = 4.14
      compress         = "none"
      max_buffer_size  = 43.17
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled             = false
    send_to_routes         = false
    socket_ending_max_wait = 6.74
    socket_idle_timeout    = 6.02
    socket_max_lifespan    = 6.18
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
      request_cert        = true
    }
    type = "tcpjson"
  }
  input_wef = {
    allow_machine_id_mismatch = false
    auth_method               = "clientCert"
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
    enable_health_check      = false
    enable_proxy_header      = true
    environment              = "...my_environment..."
    host                     = "...my_host..."
    id                       = "...my_id..."
    ip_allowlist_regex       = "...my_ip_allowlist_regex..."
    ip_denylist_regex        = "...my_ip_denylist_regex..."
    keep_alive_timeout       = 338.86
    keytab                   = "...my_keytab..."
    log_fingerprint_mismatch = true
    max_active_req           = 0.44
    max_requests_per_socket  = 5
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 906.42
    pq = {
      commit_frequency = 7.64
      compress         = "none"
      max_buffer_size  = 43.08
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "smart"
      path             = "...my_path..."
    }
    pq_enabled     = false
    principal      = "...my_principal..."
    send_to_routes = false
    socket_timeout = 8.64
    streamtags = [
      "..."
    ]
    subscriptions = [
      {
        batch_timeout      = 4.05
        compress           = false
        content_format     = "RenderedText"
        heartbeat_interval = 2.46
        id                 = "...my_id..."
        locale             = "...my_locale..."
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        query_selector       = "simple"
        read_existing_events = false
        send_bookmarks       = true
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
      disabled              = true
      keytab                = "{ \"see\": \"documentation\" }"
      max_version           = "TLSv1.3"
      min_version           = "TLSv1.1"
      ocsp_check            = true
      ocsp_check_fail_close = true
      passphrase            = "...my_passphrase..."
      principal             = "{ \"see\": \"documentation\" }"
      priv_key_path         = "...my_priv_key_path..."
      reject_unauthorized   = false
      request_cert          = true
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
    disabled              = true
    environment           = "...my_environment..."
    host = {
      custom = {
        cpu = {
          detail  = false
          mode    = "custom"
          per_cpu = true
          time    = false
        }
        disk = {
          mode       = "disabled"
          per_volume = false
          volumes = [
            "..."
          ]
        }
        memory = {
          detail = true
          mode   = "all"
        }
        network = {
          detail = false
          devices = [
            "..."
          ]
          mode          = "basic"
          per_interface = false
        }
        system = {
          detail = true
          mode   = "custom"
        }
      }
      mode = "custom"
    }
    id       = "...my_id..."
    interval = 8.96
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
      commit_frequency = 4.01
      compress         = "gzip"
      max_buffer_size  = 47.63
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled = true
    process = {
      sets = [
        {
          filter           = "...my_filter..."
          include_children = false
          name             = "...my_name..."
        }
      ]
    }
    send_to_routes = true
    streamtags = [
      "..."
    ]
    type = "windows_metrics"
  }
  input_win_event_logs = {
    batch_size = 8.86
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "...my_description..."
    disable_native_module = false
    disabled              = true
    environment           = "...my_environment..."
    event_format          = "xml"
    id                    = "...my_id..."
    interval              = 8.53
    log_names = [
      "..."
    ]
    max_event_bytes = 60321986.88
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 8.03
      compress         = "gzip"
      max_buffer_size  = 50.86
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled     = true
    read_mode      = "newest"
    send_to_routes = true
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
    disabled                = true
    endpoint                = "...my_endpoint..."
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    keep_alive_time         = 14.78
    max_missed_keep_alives  = 5.49
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    pq = {
      commit_frequency = 10.85
      compress         = "gzip"
      max_buffer_size  = 46.92
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 1826.58
    retry_rules = {
      codes = [
        5.48
      ]
      enable_header         = false
      interval              = 10483.06
      limit                 = 15.93
      multiplier            = 14.76
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "none"
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
      "..."
    ]
    access_control_allow_origin = [
      "..."
    ]
    activity_log_sample_rate = 10.12
    allowed_indexes = [
      "..."
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
          "..."
        ]
        auth_type   = "secret"
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
    enable_proxy_header     = false
    environment             = "...my_environment..."
    hec_acks                = false
    hec_api                 = "...my_hec_api..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 361.81
    max_active_req          = 7.46
    max_requests_per_socket = 7
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 56574.01
    pq = {
      commit_frequency = 5.64
      compress         = "none"
      max_buffer_size  = 43.74
      max_file_size    = "...my_max_file_size..."
      max_size         = "...my_max_size..."
      mode             = "always"
      path             = "...my_path..."
    }
    pq_enabled      = true
    request_timeout = 0.75
    send_to_routes  = false
    socket_timeout  = 5.43
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
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = true
    }
    type = "zscaler_hec"
  }
  pack = "...my_pack..."
}