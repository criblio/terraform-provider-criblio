resource "criblio_pack_source" "my_packsource" {
  group_id = "Cribl"
  id       = "pack-input-hec"
  input_appscope = {
    auth_token = "***REDACTED***"
    auth_type  = "secret"
    breaker_rulesets = [
      "appscope-lines",
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "Receive AppScope telemetry over TCP or UNIX socket"
    disabled            = false
    enable_proxy_header = false
    enable_unix_path    = false
    environment         = "main"
    filter = {
      allow = [
        {
          arg      = "-c /etc/nginx/nginx.conf"
          config   = "default"
          procname = "nginx"
        }
      ]
      transport_url = "unix:///var/run/appscope.sock"
    }
    host               = "0.0.0.0"
    id                 = "appscope-ingest"
    ip_whitelist_regex = "^10\\."
    max_active_cxn     = 2000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "/var/lib/cribl/state/appscope"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    port     = 57000
    pq = {
      commit_frequency      = 5.85
      compress              = "gzip"
      max_buffer_size       = 47.94
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 30
    socket_idle_timeout    = 60
    socket_max_lifespan    = 3600
    stale_channel_flush_ms = 1500
    streamtags = [
      "appscope",
      "observability",
    ]
    text_secret = "appscope-auth-secret"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    type              = "appscope"
    unix_socket_path  = "/var/run/appscope.sock"
    unix_socket_perms = "770"
  }
  input_azure_blob = {
    auth_type   = "manual"
    azure_cloud = "...my_azure_cloud..."
    breaker_rulesets = [
      "access-logs-v1",
      "json-breaker",
    ]
    certificate = {
      certificate_name = "...my_certificate_name..."
    }
    client_id          = "...my_client_id..."
    client_text_secret = "...my_client_text_secret..."
    connection_string  = "$$${{secret:azure_storage_connection_string}"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description     = "Azure Blob queue events ingester"
    disabled        = false
    endpoint_suffix = "...my_endpoint_suffix..."
    environment     = "main"
    file_filter     = "^logs/.*\\.json$"
    id              = "azure-blob-queue"
    max_messages    = 16
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 900
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    pq = {
      commit_frequency      = 5.88
      compress              = "none"
      max_buffer_size       = 48.23
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    queue_name             = "my-blob-notify-queue"
    send_to_routes         = true
    service_period_secs    = 5
    skip_on_error          = true
    stale_channel_flush_ms = 15000
    storage_account_name   = "...my_storage_account_name..."
    streamtags = [
      "prod",
      "azure",
    ]
    tenant_id          = "...my_tenant_id..."
    text_secret        = "...my_text_secret..."
    type               = "azure_blob"
    visibility_timeout = 300
  }
  input_cloudflare_hec = {
    access_control_allow_headers = [
    ]
    access_control_allow_origin = [
    ]
    activity_log_sample_rate = 3.88
    allowed_indexes = [
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
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
        token_secret = "...my_token_secret..."
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
    disabled                = false
    emit_token_metrics      = false
    enable_health_check     = true
    enable_proxy_header     = true
    environment             = "...my_environment..."
    hec_api                 = "...my_hec_api..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 570.14
    max_active_req          = 2.45
    max_requests_per_socket = 7
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 33269.88
    pq = {
      commit_frequency      = 3.71
      compress              = "none"
      max_buffer_size       = 45.55
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    request_timeout        = 1.68
    send_to_routes         = true
    socket_timeout         = 7.9
    stale_channel_flush_ms = 22807582.17
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = true
    }
    type = "cloudflare_hec"
  }
  input_collection = {
    breaker_rulesets = [
      "access-logs-v1",
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    disabled    = false
    environment = "main"
    id          = "collect-nginx"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    output   = "datalake"
    pipeline = "default"
    pq = {
      commit_frequency      = 3.08
      compress              = "none"
      max_buffer_size       = 42.51
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = false
    }
    send_to_routes         = true
    stale_channel_flush_ms = 15000
    streamtags = [
      "prod",
      "nginx",
    ]
    throttle_rate_per_sec = "10 MB"
    type                  = "collection"
  }
  input_confluent_cloud = {
    authentication_timeout = 15000
    auto_commit_interval   = 5000
    auto_commit_threshold  = 1000
    backoff_rate           = 3
    brokers = [
      "pkc-12345.us-central1.gcp.confluent.cloud:9092",
    ]
    connection_timeout = 15000
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "Confluent Cloud consumer for nginx access logs"
    disabled           = false
    environment        = "main"
    from_beginning     = true
    group_id           = "web-team"
    heartbeat_interval = 3000
    id                 = "ccloud-nginx"
    initial_backoff    = 500
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
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 7.13
      compress              = "none"
      max_buffer_size       = 49.55
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    request_timeout            = 60000
    sasl = {
      auth_type            = "manual"
      broker_service_class = "...my_broker_service_class..."
      client_id            = "...my_client_id..."
      client_text_secret   = "...my_client_text_secret..."
      credentials_secret   = "...my_credentials_secret..."
      disabled             = true
      keytab_location      = "...my_keytab_location..."
      mechanism            = "scram-sha-256"
      oauth_enabled        = false
      oauth_params = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      oauth_secret_type = "...my_oauth_secret_type..."
      password          = "...my_password..."
      principal         = "...my_principal..."
      sasl_extensions = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      token_url = "...my_token_url..."
      username  = "...my_username..."
    }
    send_to_routes  = true
    session_timeout = 30000
    streamtags = [
      "prod",
      "ccloud",
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
      "nginx_access",
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
    description = "Internal Cribl-generated events"
    disabled    = false
    environment = "main"
    filter      = "host=\"edge-*\" AND sourcetype!=\"metrics\""
    id          = "cribl-internal"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 9.9
      compress              = "gzip"
      max_buffer_size       = 43.66
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    send_to_routes = true
    streamtags = [
      "internal",
      "cribl",
    ]
    type = "cribl"
  }
  input_cribl_http = {
    activity_log_sample_rate = 10
    auth_tokens = [
      {
        description  = "...my_description..."
        enabled      = true
        token_secret = "...my_token_secret..."
      }
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Cribl HTTP-compatible ingestion endpoint"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "cribl-http-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency      = 5.87
      compress              = "none"
      max_buffer_size       = 49.68
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "prod",
      "cribl_http",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = false
    }
    type = "cribl_http"
  }
  input_cribl_lake_http = {
    activity_log_sample_rate = 10
    auth_tokens = [
      "lake-token-1",
      "lake-token-2",
    ]
    auth_tokens_ext = [
      {
        description = "...my_description..."
        elasticsearch_metadata = {
          default_dataset = "...my_default_dataset..."
          enabled         = false
        }
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        splunk_hec_metadata = {
          allowed_indexes_at_token = [
          ]
          default_dataset = "...my_default_dataset..."
          enabled         = false
        }
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
    description             = "Cribl Lake HTTP ingestion endpoint"
    disabled                = false
    elastic_api             = "...my_elastic_api..."
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "lake-http-ingest"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "lake-default"
    port     = 9088
    pq = {
      commit_frequency      = 5.05
      compress              = "none"
      max_buffer_size       = 47.49
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    splunk_hec_acks = false
    splunk_hec_api  = "...my_splunk_hec_api..."
    streamtags = [
      "lake",
      "ingest",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = false
    }
    type = "cribl_lake_http"
  }
  input_cribl_tcp = {
    auth_tokens = [
      {
        description  = "...my_description..."
        enabled      = true
        token_secret = "...my_token_secret..."
      }
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "This is the field used for description for this input"
    disabled              = false
    enable_load_balancing = true
    enable_proxy_header   = false
    environment           = "main"
    host                  = "0.0.0.0"
    id                    = "cribl-tcp-listener"
    max_active_cxn        = 2000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 9514
    pq = {
      commit_frequency      = 4.41
      compress              = "none"
      max_buffer_size       = 42.78
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 15
    socket_idle_timeout    = 60
    socket_max_lifespan    = 3600
    streamtags = [
      "prod",
      "cribl_tcp",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = false
    }
    type = "cribl_tcp"
  }
  input_criblmetrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description   = "Emit Cribl internal metrics"
    disabled      = false
    environment   = "main"
    full_fidelity = true
    id            = "cribl-metrics"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 1.52
      compress              = "none"
      max_buffer_size       = 50.12
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    prefix         = "cribl.logstream."
    send_to_routes = true
    streamtags = [
      "cribl",
      "internal",
    ]
    type = "criblmetrics"
  }
  input_crowdstrike = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-s3-access"
    assume_role_external_id   = "cribl-external-123"
    aws_account_id            = "123456789012"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    breaker_rulesets = [
      "crowdstrike-breaker",
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
    description            = "Ingest CrowdStrike S3 notifications and objects"
    disabled               = false
    duration_seconds       = 3600
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "utf-8"
    endpoint               = "https://s3.us-east-1.amazonaws.com"
    environment            = "main"
    file_filter            = ".*\\.json(\\.gz)?$"
    id                     = "crowdstrike-sqs"
    include_sqs_metadata   = false
    max_messages           = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers = 4
    pipeline      = "default"
    poll_timeout  = 10
    pq = {
      commit_frequency      = 4.32
      compress              = "gzip"
      max_buffer_size       = 42.07
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = true
    }
    processed_tag_key      = "processed-by"
    processed_tag_value    = "cribl-processed"
    queue_name             = "https://sqs.us-east-1.amazonaws.com/123456789012/crowdstrike-events"
    region                 = "us-east-1"
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v2"
    skip_on_error          = true
    socket_timeout         = 600
    stale_channel_flush_ms = 1500
    streamtags = [
      "crowdstrike",
      "edr",
    ]
    tag_after_processing = "true"
    type                 = "crowdstrike"
    visibility_timeout   = 300
  }
  input_datadog_agent = {
    activity_log_sample_rate = 100
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Accept Datadog Agent intake and forward to destinations"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    extract_metrics         = false
    host                    = "0.0.0.0"
    id                      = "datadog-agent-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 10518
    pq = {
      commit_frequency      = 10.88
      compress              = "gzip"
      max_buffer_size       = 45.59
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    proxy_mode = {
      enabled             = true
      reject_unauthorized = true
    }
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "datadog",
      "metrics",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
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
    description = "Generate synthetic log events for testing"
    disabled    = false
    environment = "main"
    id          = "datagen-synthetic"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 8.51
      compress              = "none"
      max_buffer_size       = 44.45
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    samples = [
      {
        events_per_sec = 200
        sample         = "apache_common.log"
      }
    ]
    send_to_routes = true
    streamtags = [
      "synthetic",
      "test",
    ]
    type = "datagen"
  }
  input_edge_prometheus = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/edge-prom-discovery"
    assume_role_external_id   = "external-123"
    auth_type                 = "kubernetes"
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "$$${{secret:aws_secret_access_key}"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret = "edge-prom-credentials"
    description        = "Edge Prometheus scraper with K8s discovery"
    dimension_list = [
      "host",
      "source",
      "region",
    ]
    disabled           = false
    discovery_type     = "k8s-pods"
    duration_seconds   = 3600
    enable_assume_role = false
    endpoint           = "https://ec2.us-east-1.amazonaws.com"
    environment        = "main"
    id                 = "edge-prom-scraper"
    interval           = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    name_list = [
      "web-*.example.com",
      "node-exporter.internal.example.com",
    ]
    password = "$$${{secret:edge_prom_password}"
    persistence = {
      compress      = "gzip"
      enable        = false
      max_data_size = "...my_max_data_size..."
      max_data_time = "...my_max_data_time..."
      time_window   = "...my_time_window..."
    }
    pipeline = "default"
    pod_filter = [
      {
        description = "Scrape pod if annotation is true"
        filter      = "metadata.annotations['prometheus.io/scrape']"
      }
    ]
    pq = {
      commit_frequency      = 8.56
      compress              = "gzip"
      max_buffer_size       = 48.47
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled           = false
    record_type          = "AAAA"
    region               = "us-east-1"
    reject_unauthorized  = true
    reuse_connections    = true
    scrape_path          = "/metrics"
    scrape_path_expr     = "metadata.annotations['prometheus.io/path'] || '/metrics'"
    scrape_port          = 9100
    scrape_port_expr     = "metadata.annotations['prometheus.io/port'] || 9100"
    scrape_protocol      = "http"
    scrape_protocol_expr = "metadata.annotations['prometheus.io/scheme'] || 'http'"
    search_filter = [
      {
        name = "...my_name..."
        values = [
        ]
      }
    ]
    send_to_routes    = true
    signature_version = "v4"
    streamtags = [
      "edge",
      "prometheus",
    ]
    targets = [
      {
        host     = "node-exporter"
        path     = "/metrics"
        port     = 9100
        protocol = "http"
      }
    ]
    timeout       = 5000
    type          = "edge_prometheus"
    use_public_ip = true
    username      = "edge_user"
  }
  input_elastic = {
    activity_log_sample_rate = 10
    api_version              = "8.3.2"
    auth_tokens = [
      "es-api-token-1",
      "es-api-token-2",
    ]
    auth_type       = "basic"
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret  = "elastic-proxy-credentials"
    custom_api_version  = "{ \\n\n    \"name\": \"Cribl Elastic Proxy\", \\n\n    \"cluster_name\": \"cribl\", \\n\n    \"cluster_uuid\": \"abcd1234efgh5678ijkl9012\", \\n\n    \"version\": { \\n\n        \"number\": \"8.11.1\", \\n\n        \"build_type\": \"tar\", \\n\n        \"build_hash\": \"1a2b3c4\", \\n\n        \"build_date\": \"2025-09-01T00:00:00.000Z\", \\n\n        \"build_snapshot\": false, \\n\n        \"lucene_version\": \"9.10.0\", \\n\n        \"minimum_wire_compatibility_version\": \"7.17.0\", \\n\n        \"minimum_index_compatibility_version\": \"7.0.0\" \\n\n    }, \\n\n    \"tagline\": \"You Know, for Search\" \\n\n}"
    description         = "Elasticsearch bulk listener with proxy for non-bulk APIs"
    disabled            = false
    elastic_api         = "/ingest"
    enable_health_check = true
    enable_proxy_header = false
    environment         = "main"
    extra_http_headers = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    host                    = "0.0.0.0"
    id                      = "elastic-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    password = "$$${{secret:elastic_proxy_password}"
    pipeline = "default"
    port     = 9200
    pq = {
      commit_frequency      = 7.29
      compress              = "gzip"
      max_buffer_size       = 46.69
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    proxy_mode = {
      auth_type           = "manual"
      credentials_secret  = "...my_credentials_secret..."
      enabled             = true
      password            = "...my_password..."
      reject_unauthorized = false
      remove_headers = [
        "Authorization",
        "Content-Length",
      ]
      timeout_sec = 60
      url         = "https://elastic.example.com:9200"
      username    = "...my_username..."
    }
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "prod",
      "elastic",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = false
    }
    type     = "elastic"
    username = "elastic"
  }
  input_eventhub = {
    authentication_timeout = 15000
    auto_commit_interval   = 5000
    auto_commit_threshold  = 1000
    backoff_rate           = 3
    brokers = [
      "yourspace.servicebus.windows.net:9093",
    ]
    connection_timeout = 15000
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Azure Event Hubs Kafka consumer"
    disabled                = false
    environment             = "main"
    from_beginning          = true
    group_id                = "web-team"
    heartbeat_interval      = 3000
    id                      = "eventhub-nginx"
    initial_backoff         = 500
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    minimize_duplicates = true
    pipeline            = "default"
    pq = {
      commit_frequency      = 6.67
      compress              = "none"
      max_buffer_size       = 50.4
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    request_timeout            = 60000
    sasl = {
      auth_type               = "manual"
      cert_path               = "...my_cert_path..."
      certificate_name        = "...my_certificate_name..."
      client_id               = "...my_client_id..."
      client_secret           = "...my_client_secret..."
      client_secret_auth_type = "secret"
      client_text_secret      = "...my_client_text_secret..."
      disabled                = false
      mechanism               = "oauthbearer"
      oauth_endpoint          = "https://login.partner.microsoftonline.cn"
      passphrase              = "...my_passphrase..."
      password                = "...my_password..."
      priv_key_path           = "...my_priv_key_path..."
      scope                   = "...my_scope..."
      tenant_id               = "...my_tenant_id..."
      text_secret             = "...my_text_secret..."
      username                = "...my_username..."
    }
    send_to_routes  = true
    session_timeout = 30000
    streamtags = [
      "prod",
      "eventhub",
    ]
    tls = {
      disabled            = true
      reject_unauthorized = false
    }
    topics = [
      "logs",
    ]
    type = "eventhub"
  }
  input_exec = {
    breaker_rulesets = [
      "access-logs-v1",
    ]
    command = "tail -F /var/log/nginx/access.log"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    cron_schedule = "*/5 * * * *"
    description   = "Exec tail of nginx access logs"
    disabled      = false
    environment   = "main"
    id            = "exec-tail-logs"
    interval      = 60
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 4.61
      compress              = "none"
      max_buffer_size       = 44.51
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    retries                = 5
    schedule_type          = "interval"
    script                 = "...my_script..."
    send_to_routes         = true
    stale_channel_flush_ms = 15000
    streamtags = [
      "prod",
      "exec",
    ]
    type = "exec"
  }
  input_file = {
    breaker_rulesets = [
      "multiline-java",
    ]
    check_file_mod_time = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    delete_files = false
    depth        = 2
    description  = "Watch local files and tail new content"
    disabled     = false
    environment  = "main"
    filenames = [
      "/var/log/*.log",
      "/opt/app/logs/*log",
    ]
    filter_archived_files         = true
    force_text                    = false
    hash_len                      = 256
    id                            = "file-watcher"
    idle_timeout                  = 600
    include_unidentifiable_binary = true
    interval                      = 10
    max_age_dur                   = "3d"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    min_age_dur = "...my_min_age_dur..."
    mode        = "auto"
    path        = "/var/log"
    pipeline    = "default"
    pq = {
      commit_frequency      = 10.43
      compress              = "none"
      max_buffer_size       = 43.03
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    salt_hash              = true
    send_to_routes         = true
    stale_channel_flush_ms = 1500
    streamtags = [
      "filesystem",
      "logs",
    ]
    suppress_missing_path_errors = true
    tail_only                    = true
    type                         = "file"
  }
  input_firehose = {
    activity_log_sample_rate = 10
    auth_tokens = [
      "secret-token-1",
      "secret-token-2",
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Kinesis Firehose-compatible HTTP listener"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "firehose-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 9000
    pq = {
      commit_frequency      = 9.1
      compress              = "none"
      max_buffer_size       = 50.79
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "prod",
      "firehose",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = false
    }
    type = "firehose"
  }
  input_google_pubsub = {
    concurrency = 10
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    create_subscription = true
    create_topic        = false
    description         = "Google Pub/Sub pull subscription consumer"
    disabled            = false
    environment         = "main"
    google_auth_method  = "auto"
    id                  = "gpubsub-nginx"
    max_backlog         = 2000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    monitor_subscription = true
    ordered_delivery     = false
    pipeline             = "default"
    pq = {
      commit_frequency      = 8.56
      compress              = "none"
      max_buffer_size       = 51.86
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled                  = false
    region                      = "us-central1"
    request_timeout             = 45000
    secret                      = "gcp-service-account"
    send_to_routes              = true
    service_account_credentials = "$$${{file:/secrets/gcp-service-account.json}"
    streamtags = [
      "prod",
      "gpubsub",
    ]
    subscription_name = "projects/my-project/subscriptions/nginx-logs-sub"
    topic_name        = "projects/my-project/topics/nginx-logs"
    type              = "google_pubsub"
  }
  input_grafana = {
    activity_log_sample_rate = 10
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "Grafana listener supporting Prom remote write and Loki logs"
    disabled            = false
    enable_health_check = true
    enable_proxy_header = false
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "grafana-listener"
    ip_allowlist_regex  = "^10\\."
    ip_denylist_regex   = "^192\\.168\\.0\\."
    keep_alive_timeout  = 30
    loki_api            = "/loki/api/v1/push"
    loki_auth = {
      auth_header_expr   = "`Bearer $${token}`"
      auth_type          = "textSecret"
      credentials_secret = "loki-credentials"
      login_url          = "https://loki.example.com/oauth/token"
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
      password             = "$$${{secret:loki_password}"
      secret               = "$$${{secret:loki_oauth_secret}"
      secret_param_name    = "client_secret"
      text_secret          = "loki-token-secret"
      token                = "$$${{secret:loki_token}"
      token_attribute_name = "access_token"
      token_timeout_secs   = 3600
      username             = "loki_user"
    }
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 4318
    pq = {
      commit_frequency      = 1.47
      compress              = "gzip"
      max_buffer_size       = 44.05
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    prometheus_api = "/api/prom/push"
    prometheus_auth = {
      auth_header_expr   = "`Bearer $${token}`"
      auth_type          = "none"
      credentials_secret = "prom-credentials"
      login_url          = "https://grafana.example.com/oauth/token"
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
      password             = "$$${{secret:prom_password}"
      secret               = "$$${{secret:prom_oauth_secret}"
      secret_param_name    = "client_secret"
      text_secret          = "prom-token-secret"
      token                = "$$${{secret:prom_token}"
      token_attribute_name = "access_token"
      token_timeout_secs   = 3600
      username             = "grafana"
    }
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "prod",
      "grafana",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    type = "grafana"
  }
  input_http = {
    activity_log_sample_rate = 10
    auth_tokens = [
      "secret-token-1",
      "secret-token-2",
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
    cribl_api               = "/cribl"
    description             = "HTTP listener for webhook events"
    disabled                = false
    elastic_api             = "/elastic"
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "http-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency      = 5.79
      compress              = "gzip"
      max_buffer_size       = 51.73
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    splunk_hec_acks = false
    splunk_hec_api  = "/services/collector"
    streamtags = [
      "prod",
      "http",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    type = "http"
  }
  input_http_raw = {
    activity_log_sample_rate = 100
    allowed_methods = [
      "POST",
      "PUT",
    ]
    allowed_paths = [
      "/api/v1/hook",
      "/webhook/*",
    ]
    auth_tokens = [
      "supersecrettoken",
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
      "http-raw-breaker",
      "multiline-json",
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Accept raw HTTP payloads"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "http-raw-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency      = 6.83
      compress              = "none"
      max_buffer_size       = 44.79
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    request_timeout        = 30
    send_to_routes         = true
    socket_timeout         = 60
    stale_channel_flush_ms = 1500
    streamtags = [
      "http",
      "raw",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
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
    current_boot = true
    description  = "Ingest systemd journal files from disk"
    disabled     = false
    environment  = "main"
    id           = "journal-files"
    interval     = 10
    journals = [
      "system",
      "user-*.journal",
    ]
    max_age_dur = "24h"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    path     = "/var/log/journal"
    pipeline = "default"
    pq = {
      commit_frequency      = 6.19
      compress              = "gzip"
      max_buffer_size       = 50.96
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    rules = [
      {
        description = "Allow warnings or higher, exclude authpriv"
        filter      = "severity <= 4 && facility != 'authpriv'"
      }
    ]
    send_to_routes = true
    streamtags = [
      "systemd",
      "journald",
    ]
    type = "journal_files"
  }
  input_kafka = {
    authentication_timeout = 15000
    auto_commit_interval   = 5000
    auto_commit_threshold  = 1000
    backoff_rate           = 3
    brokers = [
      "kafka-1:9092",
      "kafka-2:9092",
    ]
    connection_timeout = 15000
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "My Kafka Input description for documentation"
    disabled           = false
    environment        = "main"
    from_beginning     = true
    group_id           = "web-team"
    heartbeat_interval = 3000
    id                 = "kafka-nginx"
    initial_backoff    = 500
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
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 2.69
      compress              = "none"
      max_buffer_size       = 44.12
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    request_timeout            = 60000
    sasl = {
      auth_type            = "secret"
      broker_service_class = "...my_broker_service_class..."
      client_id            = "...my_client_id..."
      client_text_secret   = "...my_client_text_secret..."
      credentials_secret   = "...my_credentials_secret..."
      disabled             = false
      keytab_location      = "...my_keytab_location..."
      mechanism            = "plain"
      oauth_enabled        = true
      oauth_params = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      oauth_secret_type = "...my_oauth_secret_type..."
      password          = "...my_password..."
      principal         = "...my_principal..."
      sasl_extensions = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      token_url = "...my_token_url..."
      username  = "...my_username..."
    }
    send_to_routes  = true
    session_timeout = 30000
    streamtags = [
      "prod",
      "kafka",
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
      "nginx_access",
    ]
    type = "kafka"
  }
  input_kinesis = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-kinesis-access"
    assume_role_external_id   = "cribl-external-123"
    avoid_duplicates          = false
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description              = "Ingest AWS Kinesis stream records"
    disabled                 = false
    duration_seconds         = 3600
    enable_assume_role       = true
    endpoint                 = "https://kinesis.us-east-1.amazonaws.com"
    environment              = "main"
    get_records_limit        = 8000
    get_records_limit_total  = 30000
    id                       = "kinesis-stream-ingest"
    load_balancing_algorithm = "ConsistentHashing"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    payload_format = "ndjson"
    pipeline       = "default"
    pq = {
      commit_frequency      = 5.87
      compress              = "gzip"
      max_buffer_size       = 49.57
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled          = false
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    send_to_routes      = true
    service_interval    = 1
    shard_expr          = "shardId.endsWith('1')"
    shard_iterator_type = "LATEST"
    signature_version   = "v2"
    stream_name         = "app-logs-stream"
    streamtags = [
      "aws",
      "kinesis",
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
    description = "Collect Kubernetes cluster events"
    disabled    = false
    environment = "main"
    id          = "kube-events"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 1.98
      compress              = "gzip"
      max_buffer_size       = 51.69
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
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
      "kubernetes",
      "events",
    ]
    type = "kube_events"
  }
  input_kube_logs = {
    breaker_rulesets = [
      "kube-logs-breaker",
      "multiline-java",
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "the Description for KubeLogs type inputs"
    disabled              = false
    enable_load_balancing = true
    environment           = "main"
    id                    = "kube-logs"
    interval              = 15
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
    pipeline = "default"
    pq = {
      commit_frequency      = 1.41
      compress              = "gzip"
      max_buffer_size       = 43.46
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    rules = [
      {
        description = "Include default namespace"
        filter      = "metadata.namespace == 'default'"
      }
    ]
    send_to_routes         = true
    stale_channel_flush_ms = 1500
    streamtags = [
      "kubernetes",
      "logs",
    ]
    timestamps = true
    type       = "kube_logs"
  }
  input_kube_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description = "Collect Kubernetes metrics from the API server"
    disabled    = false
    environment = "main"
    id          = "kube-metrics"
    interval    = 15
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      dest_path     = "/var/lib/cribl/state/kube_metrics"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency      = 7.59
      compress              = "none"
      max_buffer_size       = 48.01
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
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
      "kubernetes",
      "prod",
    ]
    type = "kube_metrics"
  }
  input_loki = {
    activity_log_sample_rate = 10
    auth_header_expr         = "`Bearer $${token}`"
    auth_type                = "none"
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "loki-credentials"
    description             = "Loki logs listener"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "loki-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    login_url               = "https://loki.example.com/oauth/token"
    loki_api                = "/loki/api/v1/push"
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
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
    password = "$$${{secret:loki_password}"
    pipeline = "default"
    port     = 3100
    pq = {
      commit_frequency      = 7.66
      compress              = "none"
      max_buffer_size       = 48.78
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled        = false
    request_timeout   = 30
    secret            = "$$${{secret:loki_oauth_secret}"
    secret_param_name = "client_secret"
    send_to_routes    = true
    socket_timeout    = 60
    streamtags = [
      "prod",
      "loki",
    ]
    text_secret = "loki-token-secret"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = false
    }
    token                = "$$${{secret:loki_token}"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "loki"
    username             = "loki_user"
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
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "metrics-listener"
    ip_whitelist_regex  = "^10\\."
    max_buffer_size     = 20000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 7.61
      compress              = "none"
      max_buffer_size       = 49.38
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    send_to_routes = true
    streamtags = [
      "metrics",
      "udp",
    ]
    tcp_port = 8126
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = false
    }
    type                   = "metrics"
    udp_port               = 8125
    udp_socket_rx_buf_size = 2905545357.22
  }
  input_microsoft_graph = {
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
    interval                = 22
    job_timeout             = "...my_job_timeout..."
    keep_alive_time         = 10.89
    log_level               = "info"
    max_missed_keep_alives  = 9.03
    max_task_reschedule     = 1.7
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
      commit_frequency      = 9.74
      compress              = "none"
      max_buffer_size       = 50.1
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled               = false
    reschedule_dropped_tasks = false
    resource                 = "...my_resource..."
    retry_rules = {
      codes = [
        9.91
      ]
      enable_header         = true
      interval              = 1734.82
      limit                 = 8.61
      multiplier            = 17.8
      retry_connect_reset   = false
      retry_connect_timeout = true
      type                  = "static"
    }
    send_to_routes = true
    start_date     = "...my_start_date..."
    streamtags = [
      "..."
    ]
    tenant_id   = "...my_tenant_id..."
    text_secret = "...my_text_secret..."
    timeout     = 1347.03
    ttl         = "...my_ttl..."
    type        = "microsoft_graph"
    url         = "...my_url..."
    username    = "...my_username..."
  }
  input_model_driven_telemetry = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description    = "Receive MDT telemetry over gRPC"
    disabled       = false
    environment    = "main"
    host           = "0.0.0.0"
    id             = "mdt-grpc"
    max_active_cxn = 2000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 57000
    pq = {
      commit_frequency      = 1.29
      compress              = "none"
      max_buffer_size       = 49.43
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled          = false
    send_to_routes      = true
    shutdown_timeout_ms = 5000
    streamtags = [
      "mdt",
      "grpc",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = false
    }
    type = "model_driven_telemetry"
  }
  input_msk = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/msk-readonly"
    assume_role_external_id   = "external-123"
    authentication_timeout    = 15000
    auto_commit_interval      = 5000
    auto_commit_threshold     = 1000
    aws_api_key               = "$$${{secret:aws_access_key_id}"
    aws_authentication_method = "auto"
    aws_secret                = "aws-msk-credentials"
    aws_secret_key            = "$$${{secret:aws_secret_access_key}"
    backoff_rate              = 3
    brokers = [
      "b-1.msk-cluster.a1b2c3d4.e1.kafka.us-east-1.amazonaws.com:9092",
      "b-2.msk-cluster.a1b2c3d4.e1.kafka.us-east-1.amazonaws.com:9092",
    ]
    connection_timeout = 15000
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "MSK consumer for nginx access logs"
    disabled           = false
    duration_seconds   = 3600
    enable_assume_role = false
    endpoint           = "https://kafka.us-east-1.amazonaws.com"
    environment        = "main"
    from_beginning     = true
    group_id           = "web-team"
    heartbeat_interval = 3000
    id                 = "msk-nginx"
    initial_backoff    = 500
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
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 10.02
      compress              = "none"
      max_buffer_size       = 49.16
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    region                     = "us-east-1"
    reject_unauthorized        = true
    request_timeout            = 60000
    reuse_connections          = true
    send_to_routes             = true
    session_timeout            = 30000
    signature_version          = "v2"
    streamtags = [
      "prod",
      "msk",
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
      "nginx_access",
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
    description         = "Receive NetFlow v5/v9/IPFIX on UDP 2055"
    disabled            = false
    enable_pass_through = false
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "netflow-listener"
    ip_allowlist_regex  = "^10\\."
    ip_denylist_regex   = "^192\\.168\\.1\\.\\d{1,3}$"
    ipfix_enabled       = true
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 2055
    pq = {
      commit_frequency      = 2.83
      compress              = "gzip"
      max_buffer_size       = 51.9
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    send_to_routes = true
    streamtags = [
      "netflow",
      "network",
    ]
    template_cache_minutes = 60
    type                   = "netflow"
    udp_socket_rx_buf_size = 4194304
    v5_enabled             = true
    v9_enabled             = true
  }
  input_office365_mgmt = {
    app_id        = "99999999-aaaa-bbbb-cccc-111111111111"
    auth_type     = "secret"
    client_secret = "$$${{secret:o365_client_secret}"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        content_type = "Exchange"
        description  = "...my_description..."
        enabled      = true
        interval     = 5
        log_level    = "warn"
      }
    ]
    description             = "Office 365 Management API collector"
    disabled                = false
    environment             = "main"
    id                      = "o365-mgmt"
    ignore_group_jobs_limit = false
    ingestion_lag           = 90
    job_timeout             = "15m"
    keep_alive_time         = 30
    max_missed_keep_alives  = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline  = "default"
    plan_type = "dod"
    pq = {
      commit_frequency      = 8.47
      compress              = "none"
      max_buffer_size       = 45.85
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled           = false
    publisher_identifier = "my-company-cribl"
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
      "prod",
      "o365",
    ]
    tenant_id   = "11111111-2222-3333-4444-555555555555"
    text_secret = "o365-client-secret"
    timeout     = 300
    ttl         = "4h"
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
    client_id     = "99999999-aaaa-bbbb-cccc-111111111111"
    client_secret = "$$${{secret:o365_client_secret}"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "o365-credentials"
    description             = "Office 365 Message Trace collector"
    disable_time_filter     = true
    disabled                = false
    end_date                = "-2h@h"
    environment             = "main"
    id                      = "o365-msg-trace"
    ignore_group_jobs_limit = false
    interval                = 15
    job_timeout             = "15m"
    keep_alive_time         = 30
    log_level               = "silly"
    max_missed_keep_alives  = 3
    max_task_reschedule     = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    password  = "$$${{secret:o365_password}"
    pipeline  = "default"
    plan_type = "gcc_high"
    pq = {
      commit_frequency      = 9.74
      compress              = "none"
      max_buffer_size       = 46.49
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled               = false
    reschedule_dropped_tasks = true
    resource                 = "https://outlook.office365.com"
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
    start_date     = "-3h@h"
    streamtags = [
      "prod",
      "o365",
    ]
    tenant_id   = "11111111-2222-3333-4444-555555555555"
    text_secret = "o365-client-secret"
    timeout     = 300
    ttl         = "4h"
    type        = "office365_msg_trace"
    url         = "https://reports.office365.com/ecp/reportingwebservice/reporting.svc/MessageTrace"
    username    = "o365_user"
  }
  input_office365_service = {
    app_id        = "99999999-aaaa-bbbb-cccc-111111111111"
    auth_type     = "manual"
    client_secret = "$$${{secret:o365_client_secret}"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        content_type = "Messages"
        description  = "Poll interval minutes (1-60)"
        enabled      = true
        interval     = 5
        log_level    = "info"
      }
    ]
    description             = "Office 365 Service Health collector"
    disabled                = false
    environment             = "main"
    id                      = "o365-service"
    ignore_group_jobs_limit = false
    job_timeout             = "15m"
    keep_alive_time         = 30
    max_missed_keep_alives  = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline  = "default"
    plan_type = "gcc_high"
    pq = {
      commit_frequency      = 8.14
      compress              = "none"
      max_buffer_size       = 42.92
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
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
      "prod",
      "o365",
    ]
    tenant_id   = "11111111-2222-3333-4444-555555555555"
    text_secret = "o365-client-secret"
    timeout     = 300
    ttl         = "4h"
    type        = "office365_service"
  }
  input_open_telemetry = {
    activity_log_sample_rate = 10
    auth_header_expr         = "`Bearer $${token}`"
    auth_type                = "token"
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "otel-credentials-secret"
    description             = "Receive OpenTelemetry traces, metrics, and logs"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    extract_logs            = true
    extract_metrics         = false
    extract_spans           = true
    host                    = "0.0.0.0"
    id                      = "otel-grpc"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    login_url               = "https://auth.example.com/oauth/token"
    max_active_cxn          = 2000
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
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
    otlp_version = "1.3.1"
    password     = "***REDACTED***"
    pipeline     = "default"
    port         = 4317
    pq = {
      commit_frequency      = 10.7
      compress              = "none"
      max_buffer_size       = 48.88
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled        = false
    protocol          = "grpc"
    request_timeout   = 30
    secret            = "s3cr3t"
    secret_param_name = "client_secret"
    send_to_routes    = true
    socket_timeout    = 60
    streamtags = [
      "otel",
      "grpc",
    ]
    text_secret = "otel-token-secret"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = false
    }
    token                = "***REDACTED***"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "open_telemetry"
    username             = "otel-user"
  }
  input_openai = {
    api_key = "...my_api_key..."
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        cron_schedule = "...my_cron_schedule..."
        disabled      = false
        earliest      = "...my_earliest..."
        endpoint_metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        job_timeout = "...my_job_timeout..."
        latest      = "...my_latest..."
        log_level   = "debug"
        manage_state = {
          # ...
        }
        max_pages = 1.22
        pagination_attribute = [
          "..."
        ]
        pagination_cur_relation_attribute  = "...my_pagination_cur_relation_attribute..."
        pagination_last_page_expr          = "...my_pagination_last_page_expr..."
        pagination_next_relation_attribute = "...my_pagination_next_relation_attribute..."
        pagination_type                    = "response_body"
        request_params = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        state_merge_expression  = "...my_state_merge_expression..."
        state_tracking          = true
        state_update_expression = "...my_state_update_expression..."
      }
    ]
    description             = "...my_description..."
    disabled                = true
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    keep_alive_time         = 13.11
    max_missed_keep_alives  = 7.29
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    openai_organization = "...my_openai_organization..."
    openai_project      = "...my_openai_project..."
    pipeline            = "...my_pipeline..."
    pq = {
      commit_frequency      = 10.25
      compress              = "gzip"
      max_buffer_size       = 47.24
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 162.4
    retry_rules = {
      codes = [
        4.25
      ]
      enable_header         = false
      interval              = 2065.27
      limit                 = 15.25
      multiplier            = 18.57
      retry_connect_reset   = true
      retry_connect_timeout = false
      type                  = "backoff"
    }
    send_to_routes = true
    streamtags = [
      "..."
    ]
    text_secret = "...my_text_secret..."
    ttl         = "...my_ttl..."
    type        = "openai"
  }
  input_prometheus = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/prometheus-discovery"
    assume_role_external_id   = "external-123"
    auth_type                 = "manual"
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "secret"
    aws_secret                = "...my_aws_secret..."
    aws_secret_key            = "$$${{secret:aws_secret_access_key}"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret = "prom-credentials"
    description        = "Prometheus pull-based scraper"
    dimension_list = [
      "host",
      "source",
      "region",
    ]
    disabled                = false
    discovery_type          = "static"
    duration_seconds        = 3600
    enable_assume_role      = false
    endpoint                = "https://ec2.us-east-1.amazonaws.com"
    environment             = "main"
    id                      = "prometheus-scraper"
    ignore_group_jobs_limit = false
    interval                = 5
    job_timeout             = "15m"
    keep_alive_time         = 30
    log_level               = "info"
    max_missed_keep_alives  = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    name_list = [
      "web-*.example.com",
      "db-1.internal.example.com",
    ]
    password = "$$${{secret:prom_password}"
    pipeline = "default"
    pq = {
      commit_frequency      = 4.85
      compress              = "none"
      max_buffer_size       = 46.51
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled          = false
    record_type         = "SRV"
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    scrape_path         = "/metrics"
    scrape_port         = 9100
    scrape_protocol     = "http"
    search_filter = [
      {
        name = "...my_name..."
        values = [
        ]
      }
    ]
    send_to_routes    = true
    signature_version = "v4"
    streamtags = [
      "prod",
      "prometheus",
    ]
    target_list = [
      "http://localhost:9090/metrics",
      "node-exporter:9100",
      "db:9200/metrics",
    ]
    timeout       = 6.2
    ttl           = "4h"
    type          = "prometheus"
    use_public_ip = true
    username      = "prom_user"
  }
  input_prometheus_rw = {
    activity_log_sample_rate = 10
    auth_header_expr         = "`Bearer $${token}`"
    auth_type                = "basic"
    capture_headers          = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret      = "prom-credentials"
    description             = "Prometheus Remote Write listener"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "prom-rw-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    login_url               = "https://prom.example.com/oauth/token"
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
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
    password = "$$${{secret:prom_password}"
    pipeline = "default"
    port     = 9090
    pq = {
      commit_frequency      = 2.2
      compress              = "gzip"
      max_buffer_size       = 45.62
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled        = false
    prometheus_api    = "/write"
    request_timeout   = 30
    secret            = "$$${{secret:prom_oauth_secret}"
    secret_param_name = "client_secret"
    send_to_routes    = true
    socket_timeout    = 60
    streamtags = [
      "prod",
      "prometheus",
    ]
    text_secret = "prom-token-secret"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    token                = "$$${{secret:prom_token}"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "prometheus_rw"
    username             = "prom_user"
  }
  input_raw_udp = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "Receive raw UDP datagrams and split on newlines"
    disabled           = false
    environment        = "main"
    host               = "0.0.0.0"
    id                 = "raw-udp-listener"
    ingest_raw_bytes   = false
    ip_whitelist_regex = "^10\\."
    max_buffer_size    = 20000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 1514
    pq = {
      commit_frequency      = 10.54
      compress              = "gzip"
      max_buffer_size       = 47
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    send_to_routes         = true
    single_msg_udp_packets = true
    streamtags = [
      "udp",
      "raw",
    ]
    type                   = "raw_udp"
    udp_socket_rx_buf_size = 4194304
  }
  input_s3 = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-s3-access"
    assume_role_external_id   = "cribl-external-123"
    aws_account_id            = "123456789012"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "manual"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    breaker_rulesets = [
      "s3-breaker",
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
    description            = "Ingest S3 notifications and objects"
    disabled               = false
    duration_seconds       = 3600
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "utf-8"
    endpoint               = "https://s3.us-east-1.amazonaws.com"
    environment            = "main"
    file_filter            = ".*\\.json(\\.gz)?$"
    id                     = "s3-notifications"
    include_sqs_metadata   = true
    max_messages           = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 300
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    poll_timeout                   = 10
    pq = {
      commit_frequency      = 3.11
      compress              = "gzip"
      max_buffer_size       = 48.81
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = true
    }
    processed_tag_key      = "processed-by"
    processed_tag_value    = "cribl-processed"
    queue_name             = "https://sqs.us-east-1.amazonaws.com/123456789012/my-s3-queue"
    region                 = "us-east-1"
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v4"
    skip_on_error          = true
    socket_timeout         = 600
    stale_channel_flush_ms = 1500
    streamtags = [
      "aws",
      "s3",
    ]
    tag_after_processing = true
    type                 = "s3"
    visibility_timeout   = 300
  }
  input_s3_inventory = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-s3-access"
    assume_role_external_id   = "cribl-external-123"
    aws_account_id            = "123456789012"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    breaker_rulesets = [
      "s3-inventory-breaker",
    ]
    checkpointing = {
      enabled = false
      retries = 92.48
    }
    checksum_suffix = "checksum"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "Ingest S3 inventory manifests and listed objects"
    disabled               = false
    duration_seconds       = 3600
    enable_assume_role     = true
    enable_sqs_assume_role = true
    endpoint               = "https://s3.us-east-1.amazonaws.com"
    environment            = "main"
    file_filter            = "^.*inventory.*\\.csv(\\.gz)?$"
    id                     = "s3-inventory"
    include_sqs_metadata   = true
    max_manifest_size_kb   = 4096
    max_messages           = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 300
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    poll_timeout                   = 10
    pq = {
      commit_frequency      = 5.96
      compress              = "none"
      max_buffer_size       = 46.48
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = true
    }
    processed_tag_key      = "processed-by"
    processed_tag_value    = "cribl-processed"
    queue_name             = "https://sqs.us-east-1.amazonaws.com/123456789012/inventory-queue"
    region                 = "us-east-1"
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v4"
    skip_on_error          = true
    socket_timeout         = 600
    stale_channel_flush_ms = 1500
    streamtags = [
      "aws",
      "s3-inventory",
    ]
    tag_after_processing     = "false"
    type                     = "s3_inventory"
    validate_inventory_files = true
    visibility_timeout       = 300
  }
  input_security_lake = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-security-lake-access"
    assume_role_external_id   = "cribl-external-123"
    aws_account_id            = "123456789012"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    breaker_rulesets = [
      "security-lake-breaker",
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
    description            = "Ingest AWS Security Lake notifications and objects"
    disabled               = false
    duration_seconds       = 3600
    enable_assume_role     = true
    enable_sqs_assume_role = true
    encoding               = "utf-8"
    endpoint               = "https://s3.us-east-1.amazonaws.com"
    environment            = "main"
    file_filter            = ".*\\.json(\\.gz)?$"
    id                     = "security-lake-ingest"
    include_sqs_metadata   = true
    max_messages           = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 300
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    poll_timeout                   = 10
    pq = {
      commit_frequency      = 3.79
      compress              = "none"
      max_buffer_size       = 42.44
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = false
    }
    processed_tag_key      = "processed-by"
    processed_tag_value    = "cribl-processed"
    queue_name             = "https://sqs.us-east-1.amazonaws.com/123456789012/security-lake-queue"
    region                 = "us-east-1"
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v4"
    skip_on_error          = true
    socket_timeout         = 600
    stale_channel_flush_ms = 1500
    streamtags = [
      "aws",
      "security-lake",
    ]
    tag_after_processing = "false"
    type                 = "security_lake"
    visibility_timeout   = 300
  }
  input_snmp = {
    best_effort_parsing = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description        = "Receive SNMP traps and forward to destinations"
    disabled           = false
    environment        = "main"
    host               = "0.0.0.0"
    id                 = "snmp-traps"
    ip_whitelist_regex = "^10\\."
    max_buffer_size    = 20000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 162
    pq = {
      commit_frequency      = 2.53
      compress              = "gzip"
      max_buffer_size       = 44.33
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    send_to_routes = true
    snmp_v3_auth = {
      allow_unmatched_trap = false
      v3_auth_enabled      = true
      v3_users = [
        {
          auth_key      = "...my_auth_key..."
          auth_protocol = "sha256"
          name          = "snmp-user"
          priv_key      = "...my_priv_key..."
          priv_protocol = "aes256b"
        }
      ]
    }
    streamtags = [
      "network",
      "snmp",
    ]
    type                   = "snmp"
    udp_socket_rx_buf_size = 4194304
    varbinds_with_types    = true
  }
  input_splunk = {
    auth_tokens = [
      {
        description = "Token for prod universal forwarders"
        token       = "UF-secret-1"
      }
    ]
    breaker_rulesets = [
      "access-logs-v1",
      "syslog-breaker",
    ]
    compress = "auto"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "Splunk S2S listener for UF/HF"
    disabled            = false
    drop_control_fields = true
    enable_proxy_header = false
    environment         = "main"
    extract_metrics     = true
    host                = "0.0.0.0"
    id                  = "splunk-listener"
    ip_whitelist_regex  = "^10\\."
    max_active_cxn      = 2000
    max_s2_sversion     = "v4"
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 8089
    pq = {
      commit_frequency      = 9.29
      compress              = "none"
      max_buffer_size       = 46.84
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 15
    socket_idle_timeout    = 60
    socket_max_lifespan    = 3600
    stale_channel_flush_ms = 15000
    streamtags = [
      "prod",
      "splunk",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = true
    }
    type             = "splunk"
    use_fwd_timezone = true
  }
  input_splunk_hec = {
    access_control_allow_headers = [
      "Authorization",
      "Content-Type",
    ]
    access_control_allow_origin = [
      "https://app.example.com",
      "https://grafana.example.com",
    ]
    activity_log_sample_rate = 10
    allowed_indexes = [
      "main",
      "metrics",
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
          "main",
          "metrics",
        ]
        auth_type   = "secret"
        description = "Token for HEC webhooks"
        enabled     = true
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        token        = "...my_token..."
        token_secret = "...my_token_secret..."
      }
    ]
    breaker_rulesets = [
      "access-logs-v1",
      "syslog-breaker",
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Splunk HEC listener for webhooks"
    disabled                = false
    drop_control_fields     = true
    emit_token_metrics      = true
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    extract_metrics         = true
    host                    = "0.0.0.0"
    id                      = "splunk-hec-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.0\\."
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency      = 7.55
      compress              = "none"
      max_buffer_size       = 48.46
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    request_timeout        = 30
    send_to_routes         = true
    socket_timeout         = 60
    splunk_hec_acks        = false
    splunk_hec_api         = "/services/collector"
    stale_channel_flush_ms = 15000
    streamtags = [
      "prod",
      "splunk",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = true
    }
    type             = "splunk_hec"
    use_fwd_timezone = true
  }
  input_splunk_search = {
    auth_header_expr = "`Bearer $${token}`"
    auth_type        = "basic"
    breaker_rulesets = [
      "Splunk Search Ruleset",
      "access-logs-v1",
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    credentials_secret = "splunk-credentials"
    cron_schedule      = "*/5 * * * *"
    description        = "Scheduled Splunk search for error rates"
    disabled           = false
    earliest           = "-1h@h"
    encoding           = "UTF-8"
    endpoint           = "/services/search/v2/jobs/export"
    endpoint_headers = [
      {
        name  = "Authorization"
        value = "\"Bearer $$${{secret:splunk_token}\""
      }
    ]
    endpoint_params = [
      {
        name  = "earliest_time"
        value = "$${earliest}"
      }
    ]
    environment             = "main"
    id                      = "splunk-search-errors"
    ignore_group_jobs_limit = false
    job_timeout             = "15m"
    keep_alive_time         = 30
    latest                  = "now"
    log_level               = "info"
    login_url               = "https://splunk.example.com:8089/services/auth/login"
    max_missed_keep_alives  = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
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
    output_mode = "csv"
    password    = "$$${{secret:splunk_password}"
    pipeline    = "default"
    pq = {
      commit_frequency      = 4.05
      compress              = "none"
      max_buffer_size       = 45.31
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled          = false
    reject_unauthorized = false
    request_timeout     = 120
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
    search                 = "search index=main sourcetype=access_combined status>=500 | stats count by host"
    search_head            = "https://splunk.example.com:8089"
    secret                 = "$$${{secret:splunk_oauth_secret}"
    secret_param_name      = "password"
    send_to_routes         = true
    stale_channel_flush_ms = 15000
    streamtags = [
      "prod",
      "splunk",
    ]
    text_secret          = "splunk-token-secret"
    token                = "$$${{secret:splunk_token}"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    ttl                  = "4h"
    type                 = "splunk_search"
    use_round_robin_dns  = true
    username             = "splunk_user"
  }
  input_sqs = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/cribl-sqs-access"
    assume_role_external_id   = "cribl-external-123"
    aws_account_id            = "123456789012"
    aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    create_queue       = false
    description        = "Ingest SQS messages"
    disabled           = false
    duration_seconds   = 3600
    enable_assume_role = true
    endpoint           = "https://sqs.us-east-1.amazonaws.com"
    environment        = "main"
    id                 = "sqs-events"
    max_messages       = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    num_receivers = 4
    pipeline      = "default"
    poll_timeout  = 10
    pq = {
      commit_frequency      = 4.03
      compress              = "none"
      max_buffer_size       = 46.43
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled          = false
    queue_name          = "https://sqs.us-east-1.amazonaws.com/123456789012/events-queue"
    queue_type          = "standard"
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    send_to_routes      = true
    signature_version   = "v4"
    streamtags = [
      "aws",
      "sqs",
    ]
    type               = "sqs"
    visibility_timeout = 300
  }
  input_syslog = {
    allow_non_standard_app_name = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description                          = "Receive syslog over UDP/TCP with framing detection"
    disabled                             = false
    enable_enhanced_proxy_header_parsing = true
    enable_load_balancing                = true
    enable_proxy_header                  = false
    environment                          = "main"
    host                                 = "0.0.0.0"
    id                                   = "syslog-listener"
    infer_framing                        = true
    ip_whitelist_regex                   = "^10\\."
    keep_fields_list = [
      "host",
      "app",
    ]
    max_active_cxn  = 2000
    max_buffer_size = 20000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    octet_counting = false
    pipeline       = "default"
    pq = {
      commit_frequency      = 9.87
      compress              = "gzip"
      max_buffer_size       = 51.7
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    send_to_routes         = true
    single_msg_udp_packets = true
    socket_ending_max_wait = 30
    socket_idle_timeout    = 60
    socket_max_lifespan    = 3600
    streamtags = [
      "syslog",
      "network",
    ]
    strictly_infer_octet_counting = true
    tcp_port                      = 514
    timestamp_timezone            = "UTC"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    type                   = "syslog"
    udp_port               = 514
    udp_socket_rx_buf_size = 4194304
  }
  input_system_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    container = {
      all_containers = true
      detail         = true
      docker_socket = [
        "/var/run/docker.sock",
      ]
      docker_timeout = 10
      filters = [
        {
          expr = "container.name =~ /nginx|redis/"
        }
      ]
      mode       = "all"
      per_device = true
    }
    description = "Host, CPU, memory, network, disk, process and container metrics"
    disabled    = false
    environment = "main"
    host = {
      custom = {
        cpu = {
          detail  = true
          mode    = "all"
          per_cpu = true
          time    = true
        }
        disk = {
          detail = true
          devices = [
            "!loop*",
            "sda*",
          ]
          fstypes = [
            "ext4",
            "!*tmpfs",
          ]
          inodes = true
          mode   = "all"
          mountpoints = [
            "/",
            "/var",
            "!/proc*",
          ]
          per_device = true
        }
        memory = {
          detail = true
          mode   = "all"
        }
        network = {
          detail = true
          devices = [
            "!lo",
            "eth0",
          ]
          mode          = "custom"
          per_interface = true
          protocols     = true
        }
        system = {
          mode      = "basic"
          processes = true
        }
      }
      mode = "basic"
    }
    id       = "sysmetrics"
    interval = 15
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "/opt/cribl/state/system_metrics"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "48h"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency      = 6.83
      compress              = "none"
      max_buffer_size       = 48.26
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
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
      "sys",
      "metrics",
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
        enable = true
      }
      hostsfile = {
        enable = true
      }
      interfaces = {
        enable = true
      }
      login_users = {
        enable = true
      }
      metadata = {
        enable = true
      }
      ports = {
        enable = true
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
    description                    = "Collect system state metrics and spool to disk"
    disable_native_last_log_module = false
    disable_native_module          = false
    disabled                       = false
    environment                    = "main"
    id                             = "system-state"
    interval                       = 600
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "/var/lib/cribl/state/system_state"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency      = 3.31
      compress              = "none"
      max_buffer_size       = 50.09
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    send_to_routes = true
    streamtags = [
      "sys",
      "state",
    ]
    type = "system_state"
  }
  input_tcp = {
    auth_token = "...my_auth_token..."
    auth_type  = "manual"
    breaker_rulesets = [
      "multiline-json",
      "tcp-syslog-breaker",
    ]
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description         = "Receive generic TCP payloads"
    disabled            = false
    enable_header       = false
    enable_proxy_header = false
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "tcp-listener"
    ip_whitelist_regex  = "^10\\."
    max_active_cxn      = 2000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 9000
    pq = {
      commit_frequency      = 7.04
      compress              = "none"
      max_buffer_size       = 51.06
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
    preprocess = {
      args = [
        "..."
      ]
      command  = "...my_command..."
      disabled = true
    }
    send_to_routes         = true
    socket_ending_max_wait = 30
    socket_idle_timeout    = 60
    socket_max_lifespan    = 3600
    stale_channel_flush_ms = 1500
    streamtags = [
      "tcp",
      "ingest",
    ]
    text_secret = "...my_text_secret..."
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.1"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = false
    }
    type = "tcp"
  }
  input_tcpjson = {
    auth_token = "$$${{secret:tcpjson_token}"
    auth_type  = "secret"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "TCP JSON listener for app logs"
    disabled              = false
    enable_load_balancing = true
    enable_proxy_header   = false
    environment           = "main"
    host                  = "0.0.0.0"
    id                    = "tcpjson-listener"
    ip_whitelist_regex    = "^10\\."
    max_active_cxn        = 2000
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 9001
    pq = {
      commit_frequency      = 4.14
      compress              = "none"
      max_buffer_size       = 43.17
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    send_to_routes         = true
    socket_ending_max_wait = 15
    socket_idle_timeout    = 60
    socket_max_lifespan    = 3600
    streamtags = [
      "prod",
      "tcpjson",
    ]
    text_secret = "tcpjson-token-secret"
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = true
      request_cert        = true
    }
    type = "tcpjson"
  }
  input_wef = {
    allow_machine_id_mismatch = false
    auth_method               = "clientCert"
    ca_fingerprint            = "9A:4F:2B:8E:1D:3C:A7:5B:9E:0F:11:22:33:44:55:66:77:88:99:AA"
    capture_headers           = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description              = "Receive Windows Event Forwarding (WEF) over HTTPS"
    disabled                 = false
    enable_health_check      = true
    enable_proxy_header      = true
    environment              = "main"
    host                     = "0.0.0.0"
    id                       = "wef-listener"
    ip_allowlist_regex       = "^10\\."
    ip_denylist_regex        = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout       = 60
    keytab                   = "/etc/krb5.keytab"
    log_fingerprint_mismatch = true
    max_active_req           = 512
    max_requests_per_socket  = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 5986
    pq = {
      commit_frequency      = 7.64
      compress              = "none"
      max_buffer_size       = 43.08
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    principal      = "HTTP/wef.example.com@EXAMPLE.COM"
    send_to_routes = true
    socket_timeout = 60
    streamtags = [
      "windows",
      "wef",
    ]
    subscriptions = [
      {
        batch_timeout      = 30
        compress           = true
        content_format     = "Raw"
        heartbeat_interval = 60
        id                 = "default-subscription"
        locale             = "en-US"
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        queries = [
          {
            path             = "...my_path..."
            query_expression = "...my_query_expression..."
          }
        ]
        query_selector       = "simple"
        read_existing_events = false
        send_bookmarks       = true
        subscription_name    = "Security"
        targets = [
          "wef1.corp.local",
          "*.corp.local",
        ]
        version   = "7f0c2f2e-1c3b-4d2a-9d6e-5a1b2c3d4e5f"
        xml_query = "...my_xml_query..."
      }
    ]
    tls = {
      ca_path               = "/etc/ssl/certs/ca-bundle.crt"
      cert_path             = "/etc/ssl/certs/server.crt"
      certificate_name      = "wef-cert"
      common_name_regex     = "^WEF-CLIENT-.*$"
      disabled              = false
      keytab                = "{ \"see\": \"documentation\" }"
      max_version           = "TLSv1.3"
      min_version           = "TLSv1.1"
      ocsp_check            = false
      ocsp_check_fail_close = false
      passphrase            = "***REDACTED***"
      principal             = "{ \"see\": \"documentation\" }"
      priv_key_path         = "/etc/ssl/private/server.key"
      reject_unauthorized   = true
      request_cert          = true
    }
    type = "wef"
  }
  input_win_event_logs = {
    batch_size = 500
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description            = "Collect Windows Event Logs from local system"
    disable_json_rendering = true
    disable_native_module  = false
    disable_xml_rendering  = true
    disabled               = false
    environment            = "main"
    event_format           = "json"
    id                     = "win-event-logs"
    interval               = 10
    log_names = [
      "Application",
      "Security",
      "System",
      "Microsoft-Windows-Sysmon/Operational",
    ]
    max_event_bytes = 131072
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 8.03
      compress              = "gzip"
      max_buffer_size       = 50.86
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled     = false
    read_mode      = "newest"
    send_to_routes = true
    streamtags = [
      "windows",
      "eventlogs",
    ]
    type = "win_event_logs"
  }
  input_windows_metrics = {
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description           = "Collect Windows performance counters and spool to disk"
    disable_native_module = false
    disabled              = false
    environment           = "main"
    host = {
      custom = {
        cpu = {
          detail  = true
          mode    = "basic"
          per_cpu = true
          time    = true
        }
        disk = {
          detail     = false
          mode       = "all"
          per_volume = true
          volumes = [
            "...",
            "!E:",
          ]
        }
        memory = {
          detail = true
          mode   = "basic"
        }
        network = {
          detail = true
          devices = [
            "Ethernet*",
            "!Loopback*",
          ]
          mode          = "basic"
          per_interface = true
          protocols     = false
        }
        system = {
          detail = true
          mode   = "basic"
        }
      }
      mode = "custom"
    }
    id       = "windows-metrics"
    interval = 10
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    persistence = {
      compress      = "none"
      dest_path     = "/var/lib/cribl/state/windows_metrics"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency      = 4.01
      compress              = "gzip"
      max_buffer_size       = 47.63
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled = false
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
      "windows",
      "metrics",
    ]
    type = "windows_metrics"
  }
  input_wiz = {
    auth_audience_override = "wiz-api"
    auth_type              = "manual"
    auth_url               = "https://auth.app.wiz.io/oauth/token"
    client_id              = "123e4567-client-id"
    client_secret          = "***REDACTED***"
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    content_config = [
      {
        content_description = "...my_content_description..."
        content_query       = "...my_content_query..."
        content_type        = "...my_content_type..."
        cron_schedule       = "...my_cron_schedule..."
        earliest            = "...my_earliest..."
        enabled             = true
        job_timeout         = "...my_job_timeout..."
        latest              = "...my_latest..."
        log_level           = "debug"
        manage_state = {
          # ...
        }
        max_pages               = 0.47
        state_merge_expression  = "...my_state_merge_expression..."
        state_tracking          = true
        state_update_expression = "...my_state_update_expression..."
      }
    ]
    description             = "Ingest Wiz content via GraphQL API"
    disabled                = false
    endpoint                = "https://api.us1.app.wiz.io/graphql"
    environment             = "main"
    id                      = "wiz-ingest"
    ignore_group_jobs_limit = false
    keep_alive_time         = 30
    max_missed_keep_alives  = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency      = 10.85
      compress              = "gzip"
      max_buffer_size       = 46.92
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 300
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
    send_to_routes = true
    streamtags = [
      "wiz",
      "security",
    ]
    text_secret = "wiz-client-secret"
    ttl         = "4h"
    type        = "wiz"
  }
  input_wiz_webhook = {
    activity_log_sample_rate = 3.29
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
    disabled                = true
    enable_health_check     = false
    enable_proxy_header     = true
    environment             = "...my_environment..."
    host                    = "...my_host..."
    id                      = "...my_id..."
    ip_allowlist_regex      = "...my_ip_allowlist_regex..."
    ip_denylist_regex       = "...my_ip_denylist_regex..."
    keep_alive_timeout      = 507.02
    max_active_req          = 5.16
    max_requests_per_socket = 3
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "...my_pipeline..."
    port     = 17894.17
    pq = {
      commit_frequency      = 5.71
      compress              = "none"
      max_buffer_size       = 46.61
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "smart"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled             = false
    request_timeout        = 0.37
    send_to_routes         = true
    socket_timeout         = 6.98
    stale_channel_flush_ms = 17177358.87
    streamtags = [
      "..."
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = true
      max_version         = "TLSv1"
      min_version         = "TLSv1.3"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    type = "wiz_webhook"
  }
  input_zscaler_hec = {
    access_control_allow_headers = [
      "Authorization",
      "Content-Type",
    ]
    access_control_allow_origin = [
      "https://*.zscaler.com",
    ]
    activity_log_sample_rate = 100
    allowed_indexes = [
      "zscaler-*",
    ]
    auth_tokens = [
      {
        allowed_indexes_at_token = [
          "zscaler-*",
        ]
        auth_type   = "secret"
        description = "Zscaler Collector"
        enabled     = true
        metadata = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        token        = "...my_token..."
        token_secret = "...my_token_secret..."
      }
    ]
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Receive Zscaler HEC events over HTTP(S)"
    disabled                = false
    emit_token_metrics      = true
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    hec_acks                = false
    hec_api                 = "/services/collector"
    host                    = "0.0.0.0"
    id                      = "zscaler-hec-listener"
    ip_allowlist_regex      = "^10\\."
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "...my_name..."
        value = "...my_value..."
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency      = 5.64
      compress              = "none"
      max_buffer_size       = 43.74
      max_buffer_size_bytes = "...my_max_buffer_size_bytes..."
      max_file_size         = "...my_max_file_size..."
      max_size              = "...my_max_size..."
      mode                  = "always"
      path                  = "...my_path..."
      pq_controls = {
        # ...
      }
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "zscaler",
      "hec",
    ]
    tls = {
      ca_path             = "...my_ca_path..."
      cert_path           = "...my_cert_path..."
      certificate_name    = "...my_certificate_name..."
      common_name_regex   = "...my_common_name_regex..."
      disabled            = false
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.2"
      passphrase          = "...my_passphrase..."
      priv_key_path       = "...my_priv_key_path..."
      reject_unauthorized = false
      request_cert        = true
    }
    type = "zscaler_hec"
  }
  pack = "observability-pack"
}