resource "criblio_source" "my_source" {
  group_id = "default"
  id       = "input-hec-1"
  input_appscope = {
    auth_token = "***REDACTED***"
    auth_type  = "manual"
    breaker_rulesets = [
      "appscope-lines",
    ]
    connections = [
      {
        output   = "s3-appscope"
        pipeline = "default"
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
    ip_whitelist_regex = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_active_cxn     = 2000
    metadata = [
      {
        name  = "source"
        value = "\"appscope\""
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
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "appscope-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type              = "appscope"
    unix_socket_path  = "/var/run/appscope.sock"
    unix_socket_perms = "770"
  }
  input_azure_blob = {
    auth_type   = "clientSecret"
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
    connection_string  = "$${{secret:azure_storage_connection_string}"
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
        name  = "source"
        value = "\"azure_blob\""
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 900
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
        name  = "source"
        value = "\"collection\""
      }
    ]
    output   = "datalake"
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    preprocess = {
      args = [
        "-c",
        ".message |= upcase",
      ]
      command  = "jq"
      disabled = true
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
        credentials_secret = "ccloud-schema-basic"
        disabled           = true
      }
      connection_timeout  = 30000
      disabled            = true
      max_retries         = 3
      request_timeout     = 30000
      schema_registry_url = "https://psrc-12345.us-central1.gcp.confluent.cloud"
      tls = {
        ca_path             = "/etc/ssl/certs/ca.pem"
        cert_path           = "/etc/ssl/certs/client.crt"
        certificate_name    = "ccloud-schema-cert"
        disabled            = true
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "$${{secret:ccloud_key_pass}"
        priv_key_path       = "/etc/ssl/private/client.key"
        reject_unauthorized = true
        servername          = "psrc-12345.us-central1.gcp.confluent.cloud"
      }
    }
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 0
    metadata = [
      {
        name  = "source"
        value = "\"ccloud\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    request_timeout            = 60000
    sasl = {
      disabled  = true
      mechanism = "scram-sha-512"
    }
    send_to_routes  = true
    session_timeout = 30000
    streamtags = [
      "prod",
      "ccloud",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "ccloud-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:ccloud_key_pass}"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "pkc-12345.us-central1.gcp.confluent.cloud"
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
        name  = "source"
        value = "\"cribl_internal\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
    description             = "Cribl HTTP-compatible ingestion endpoint"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "cribl-http-listener"
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"cribl_http\""
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "cribl-http-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:cribl_http_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
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
    capture_headers = true
    connections = [
      {
        output   = "...my_output..."
        pipeline = "...my_pipeline..."
      }
    ]
    description             = "Cribl Lake HTTP ingestion endpoint"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "lake-http-ingest"
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"cribl_lake_http\""
      }
    ]
    pipeline = "lake-default"
    port     = 9088
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "lake",
      "ingest",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "lake-http-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:lake_http_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "cribl_lake_http"
  }
  input_cribl_tcp = {
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
        name  = "source"
        value = "\"cribl_tcp\""
      }
    ]
    pipeline = "default"
    port     = 9514
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "cribl-tcp-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:cribl_tcp_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "cribl_tcp"
  }
  input_criblmetrics = {
    connections = [
      {
        output   = "s3-metrics"
        pipeline = "default"
      }
    ]
    description   = "Emit Cribl internal metrics"
    disabled      = false
    environment   = "main"
    full_fidelity = true
    id            = "cribl-metrics"
    metadata = [
      {
        name  = "source"
        value = "\"criblmetrics\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
      retries = 3
    }
    connections = [
      {
        output   = "s3-crowdstrike"
        pipeline = "default"
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
    max_messages           = 10
    metadata = [
      {
        name  = "source"
        value = "\"crowdstrike\""
      }
    ]
    num_receivers = 4
    pipeline      = "default"
    poll_timeout  = 10
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    preprocess = {
      args = [
        "--compact-output",
      ]
      command  = "jq -r .message"
      disabled = true
    }
    processed_tag_key      = "processed-by"
    processed_tag_value    = "cribl-processed"
    queue_name             = "https://sqs.us-east-1.amazonaws.com/123456789012/crowdstrike-events"
    region                 = "us-east-1"
    reject_unauthorized    = true
    reuse_connections      = true
    send_to_routes         = true
    signature_version      = "v4"
    skip_on_error          = true
    socket_timeout         = 600
    stale_channel_flush_ms = 1500
    streamtags = [
      "crowdstrike",
      "edr",
    ]
    tag_after_processing = "...my_tag_after_processing..."
    type                 = "crowdstrike"
    visibility_timeout   = 300
  }
  input_datadog_agent = {
    activity_log_sample_rate = 100
    capture_headers          = true
    connections = [
      {
        output   = "s3-dd"
        pipeline = "default"
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
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "source"
        value = "\"datadog_agent\""
      }
    ]
    pipeline = "default"
    port     = 10518
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "dd-agent-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "datadog_agent"
  }
  input_datagen = {
    connections = [
      {
        output   = "s3-synthetic"
        pipeline = "default"
      }
    ]
    description = "Generate synthetic log events for testing"
    disabled    = false
    environment = "main"
    id          = "datagen-synthetic"
    metadata = [
      {
        name  = "source"
        value = "\"datagen\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
    aws_authentication_method = "auto"
    aws_secret_key            = "$${{secret:aws_secret_access_key}"
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
        name  = "source"
        value = "\"edge_prometheus\""
      }
    ]
    name_list = [
      "web-*.example.com",
      "node-exporter.internal.example.com",
    ]
    password = "$${{secret:edge_prom_password}"
    persistence = {
      compress      = "gzip"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "48h"
      time_window   = "10m"
    }
    pipeline = "default"
    pod_filter = [
      {
        description = "Scrape pod if annotation is true"
        filter      = "metadata.annotations['prometheus.io/scrape']"
      }
    ]
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled           = false
    record_type          = "SRV"
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
          "..."
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
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"elastic\""
      }
    ]
    password = "$${{secret:elastic_proxy_password}"
    pipeline = "default"
    port     = 9200
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    proxy_mode = {
      auth_type           = "manual"
      enabled             = true
      reject_unauthorized = false
      remove_headers = [
        "Authorization",
        "Content-Length",
      ]
      timeout_sec = 60
      url         = "https://elastic.example.com:9200"
    }
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    streamtags = [
      "prod",
      "elastic",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "elastic-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:elastic_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
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
        name  = "source"
        value = "\"eventhub\""
      }
    ]
    minimize_duplicates = true
    pipeline            = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    request_timeout            = 60000
    sasl = {
      disabled  = false
      mechanism = "plain"
    }
    send_to_routes  = true
    session_timeout = 30000
    streamtags = [
      "prod",
      "eventhub",
    ]
    tls = {
      disabled            = false
      reject_unauthorized = true
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
        name  = "source"
        value = "\"exec\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled             = false
    retries                = 5
    schedule_type          = "interval"
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
        output   = "s3-logs"
        pipeline = "default"
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
    force_text                    = false
    hash_len                      = 256
    id                            = "file-watcher"
    idle_timeout                  = 600
    include_unidentifiable_binary = true
    interval                      = 10
    max_age_dur                   = "3d"
    metadata = [
      {
        name  = "source"
        value = "\"file\""
      }
    ]
    mode     = "auto"
    path     = "/var/log"
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled             = false
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
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"firehose\""
      }
    ]
    pipeline = "default"
    port     = 9000
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "firehose-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:firehose_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
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
    google_auth_method  = "secret"
    id                  = "gpubsub-nginx"
    max_backlog         = 2000
    metadata = [
      {
        name  = "source"
        value = "\"gpubsub\""
      }
    ]
    ordered_delivery = false
    pipeline         = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled                  = false
    region                      = "us-central1"
    request_timeout             = 45000
    secret                      = "gcp-service-account"
    send_to_routes              = true
    service_account_credentials = "$${{file:/secrets/gcp-service-account.json}"
    streamtags = [
      "prod",
      "gpubsub",
    ]
    subscription_name = "projects/my-project/subscriptions/nginx-logs-sub"
    topic_name        = "projects/my-project/topics/nginx-logs"
    type              = "google_pubsub"
  }
  input_grafana = {
    # ...
  }
  input_http = {
    activity_log_sample_rate = 10
    auth_tokens = [
      "secret-token-1",
      "secret-token-2",
    ]
    auth_tokens_ext = [
      {
        description = "Token for webhook ingestion"
        metadata = [
          {
            name  = "sourcetype"
            value = "\"http:webhook\""
          }
        ]
        token = "secret-token-1"
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
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"http\""
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "http-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:http_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
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
        description = "CI webhook"
        metadata = [
          {
            name  = "team"
            value = "\"platform\""
          }
        ]
        token = "supersecrettoken"
      }
    ]
    breaker_rulesets = [
      "http-raw-breaker",
      "multiline-json",
    ]
    capture_headers = true
    connections = [
      {
        output   = "s3-raw"
        pipeline = "default"
      }
    ]
    description             = "Accept raw HTTP payloads"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "http-raw-listener"
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "source"
        value = "\"http_raw\""
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "http-raw-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "http_raw"
  }
  input_journal_files = {
    connections = [
      {
        output   = "s3-journald"
        pipeline = "default"
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
        name  = "source"
        value = "\"journal_files\""
      }
    ]
    path     = "/var/log/journal"
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
        credentials_secret = "schema-basic-auth"
        disabled           = true
      }
      connection_timeout  = 30000
      disabled            = true
      max_retries         = 3
      request_timeout     = 30000
      schema_registry_url = "http://schema-registry:8081"
      tls = {
        ca_path             = "/etc/ssl/certs/ca.pem"
        cert_path           = "/etc/ssl/certs/client.crt"
        certificate_name    = "schema-registry-cert"
        disabled            = true
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "$${{secret:kafka_key_pass}"
        priv_key_path       = "/etc/ssl/private/client.key"
        reject_unauthorized = true
        servername          = "schema-registry"
      }
    }
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 10
    metadata = [
      {
        name  = "source"
        value = "\"kafka\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled                 = false
    reauthentication_threshold = 300000
    rebalance_timeout          = 60000
    request_timeout            = 60000
    sasl = {
      disabled  = true
      mechanism = "scram-sha-512"
    }
    send_to_routes  = true
    session_timeout = 30000
    streamtags = [
      "prod",
      "kafka",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "kafka-client-cert"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:kafka_key_pass}"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "kafka-1"
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
        output   = "s3-kinesis"
        pipeline = "default"
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
        name  = "source"
        value = "\"kinesis\""
      }
    ]
    payload_format = "ndjson"
    pipeline       = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled          = false
    region              = "us-east-1"
    reject_unauthorized = true
    reuse_connections   = true
    send_to_routes      = true
    service_interval    = 1
    shard_expr          = "shardId.endsWith('1')"
    shard_iterator_type = "LATEST"
    signature_version   = "v4"
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
        output   = "s3-events"
        pipeline = "default"
      }
    ]
    description = "Collect Kubernetes cluster events"
    disabled    = false
    environment = "main"
    id          = "kube-events"
    metadata = [
      {
        name  = "source"
        value = "\"kube_events\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    rules = [
      {
        description = "Only warning events"
        filter      = "type == 'Warning'"
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
        output   = "s3-logs"
        pipeline = "default"
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
        name  = "source"
        value = "\"kube_logs\""
      }
    ]
    persistence = {
      compress      = "gzip"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
        output   = "s3-main"
        pipeline = "default"
      }
    ]
    description = "Collect Kubernetes metrics from the API server"
    disabled    = false
    environment = "main"
    id          = "kube-metrics"
    interval    = 15
    metadata = [
      {
        name  = "source"
        value = "\"kube_metrics\""
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "/var/lib/cribl/state/kube_metrics"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    rules = [
      {
        description = "Include default namespace"
        filter      = "metadata.namespace == 'default'"
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
    auth_header_expr         = "`Bearer ${token}`"
    auth_type                = "token"
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
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    login_url               = "https://loki.example.com/oauth/token"
    loki_api                = "/loki/api/v1/push"
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"loki\""
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
    password = "$${{secret:loki_password}"
    pipeline = "default"
    port     = 3100
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled        = false
    request_timeout   = 30
    secret            = "$${{secret:loki_oauth_secret}"
    secret_param_name = "client_secret"
    send_to_routes    = true
    socket_timeout    = 60
    streamtags = [
      "prod",
      "loki",
    ]
    text_secret = "loki-token-secret"
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "loki-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:loki_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "$${{secret:loki_token}"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "loki"
    username             = "loki_user"
  }
  input_metrics = {
    connections = [
      {
        output   = "s3-metrics"
        pipeline = "default"
      }
    ]
    description         = "...my_description..."
    disabled            = false
    enable_proxy_header = false
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "metrics-listener"
    ip_whitelist_regex  = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_buffer_size     = 20000
    metadata = [
      {
        name  = "source"
        value = "\"metrics\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled     = false
    send_to_routes = true
    streamtags = [
      "metrics",
      "udp",
    ]
    tcp_port = 8126
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "metrics-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type                   = "metrics"
    udp_port               = 8125
    udp_socket_rx_buf_size = 2294508638.44
  }
  input_model_driven_telemetry = {
    connections = [
      {
        output   = "s3-network"
        pipeline = "default"
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
        name  = "source"
        value = "\"model_driven_telemetry\""
      }
    ]
    pipeline = "default"
    port     = 57000
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled          = false
    send_to_routes      = true
    shutdown_timeout_ms = 5000
    streamtags = [
      "mdt",
      "grpc",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "mdt-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
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
    aws_api_key               = "$${{secret:aws_access_key_id}"
    aws_authentication_method = "secret"
    aws_secret                = "aws-msk-credentials"
    aws_secret_key            = "$${{secret:aws_secret_access_key}"
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
        credentials_secret = "schema-basic-auth"
        disabled           = true
      }
      connection_timeout  = 30000
      disabled            = true
      max_retries         = 3
      request_timeout     = 30000
      schema_registry_url = "http://schema-registry:8081"
      tls = {
        ca_path             = "/etc/ssl/certs/ca.pem"
        cert_path           = "/etc/ssl/certs/client.crt"
        certificate_name    = "schema-registry-cert"
        disabled            = true
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "$${{secret:kafka_key_pass}"
        priv_key_path       = "/etc/ssl/private/client.key"
        reject_unauthorized = true
        servername          = "schema-registry"
      }
    }
    max_back_off            = 120000
    max_bytes               = 10485760
    max_bytes_per_partition = 1048576
    max_retries             = 10
    max_socket_errors       = 0
    metadata = [
      {
        name  = "source"
        value = "\"msk\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
    signature_version          = "v4"
    streamtags = [
      "prod",
      "msk",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/client.crt"
      certificate_name    = "msk-client-cert"
      disabled            = false
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:kafka_key_pass}"
      priv_key_path       = "/etc/ssl/private/client.key"
      reject_unauthorized = true
      servername          = "msk-broker-1"
    }
    topics = [
      "nginx_access",
    ]
    type = "msk"
  }
  input_netflow = {
    connections = [
      {
        output   = "s3-netflow"
        pipeline = "default"
      }
    ]
    description         = "Receive NetFlow v5/v9/IPFIX on UDP 2055"
    disabled            = false
    enable_pass_through = false
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "netflow-listener"
    ip_allowlist_regex  = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex   = "^192\\.168\\.1\\.\\d{1,3}$"
    ipfix_enabled       = true
    metadata = [
      {
        name  = "source"
        value = "\"netflow\""
      }
    ]
    pipeline = "default"
    port     = 2055
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
    auth_type     = "manual"
    client_secret = "$${{secret:o365_client_secret}"
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
        log_level    = "info"
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
        name  = "source"
        value = "\"office365_mgmt\""
      }
    ]
    pipeline  = "default"
    plan_type = "enterprise_gcc"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled           = false
    publisher_identifier = "my-company-cribl"
    retry_rules = {
      codes = [
        429,
        500,
        503,
      ]
      enable_header         = true
      interval              = 1000
      limit                 = 5
      multiplier            = 2
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "backoff"
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
    auth_type = "oauth"
    cert_options = {
      cert_path        = "/etc/ssl/certs/client.crt"
      certificate_name = "o365-cert"
      passphrase       = "$${{secret:o365_cert_pass}"
      priv_key_path    = "/etc/ssl/private/client.key"
    }
    client_id     = "99999999-aaaa-bbbb-cccc-111111111111"
    client_secret = "$${{secret:o365_client_secret}"
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
    log_level               = "info"
    max_missed_keep_alives  = 3
    max_task_reschedule     = 3
    metadata = [
      {
        name  = "source"
        value = "\"office365_msg_trace\""
      }
    ]
    password  = "$${{secret:o365_password}"
    pipeline  = "default"
    plan_type = "enterprise_gcc"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled               = false
    reschedule_dropped_tasks = true
    resource                 = "https://outlook.office365.com"
    retry_rules = {
      codes = [
        429,
        500,
        503,
      ]
      enable_header         = true
      interval              = 1000
      limit                 = 5
      multiplier            = 2
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "backoff"
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
    client_secret = "$${{secret:o365_client_secret}"
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
        name  = "source"
        value = "\"office365_service\""
      }
    ]
    pipeline  = "default"
    plan_type = "enterprise_gcc"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    retry_rules = {
      codes = [
        429,
        500,
        503,
      ]
      enable_header         = true
      interval              = 1000
      limit                 = 5
      multiplier            = 2
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "backoff"
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
    activity_log_sample_rate = "{ \"see\": \"documentation\" }"
    auth_header_expr         = "`Bearer ${token}`"
    auth_type                = "token"
    capture_headers          = "{ \"see\": \"documentation\" }"
    connections = [
      {
        output   = "s3-otel"
        pipeline = "default"
      }
    ]
    credentials_secret      = "otel-credentials-secret"
    description             = "Receive OpenTelemetry traces, metrics, and logs"
    disabled                = false
    enable_health_check     = true
    enable_proxy_header     = "{ \"see\": \"documentation\" }"
    environment             = "main"
    extract_logs            = true
    extract_metrics         = false
    extract_spans           = true
    host                    = "0.0.0.0"
    id                      = "otel-grpc"
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    login_url               = "https://auth.example.com/oauth/token"
    max_active_cxn          = 2000
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "source"
        value = "\"open_telemetry\""
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
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "otel-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "***REDACTED***"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "open_telemetry"
    username             = "otel-user"
  }
  input_prometheus = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/prometheus-discovery"
    assume_role_external_id   = "external-123"
    auth_type                 = "manual"
    aws_authentication_method = "auto"
    aws_secret_key            = "$${{secret:aws_secret_access_key}"
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
        name  = "source"
        value = "\"prometheus\""
      }
    ]
    name_list = [
      "web-*.example.com",
      "db-1.internal.example.com",
    ]
    password = "$${{secret:prom_password}"
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
          "..."
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
    ttl           = "4h"
    type          = "prometheus"
    use_public_ip = true
    username      = "prom_user"
  }
  input_prometheus_rw = {
    activity_log_sample_rate = 10
    auth_header_expr         = "`Bearer ${token}`"
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
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    login_url               = "https://prom.example.com/oauth/token"
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"prometheus_rw\""
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
    password = "$${{secret:prom_password}"
    pipeline = "default"
    port     = 9090
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled        = false
    prometheus_api    = "/write"
    request_timeout   = 30
    secret            = "$${{secret:prom_oauth_secret}"
    secret_param_name = "client_secret"
    send_to_routes    = true
    socket_timeout    = 60
    streamtags = [
      "prod",
      "prometheus",
    ]
    text_secret = "prom-token-secret"
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "prom-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:prom_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    token                = "$${{secret:prom_token}"
    token_attribute_name = "access_token"
    token_timeout_secs   = 3600
    type                 = "prometheus_rw"
    username             = "prom_user"
  }
  input_raw_udp = {
    connections = [
      {
        output   = "s3-raw"
        pipeline = "default"
      }
    ]
    description        = "Receive raw UDP datagrams and split on newlines"
    disabled           = false
    environment        = "main"
    host               = "0.0.0.0"
    id                 = "raw-udp-listener"
    ingest_raw_bytes   = false
    ip_whitelist_regex = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_buffer_size    = 20000
    metadata = [
      {
        name  = "source"
        value = "\"raw_udp\""
      }
    ]
    pipeline = "default"
    port     = 1514
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
    aws_authentication_method = "auto"
    aws_secret                = "aws-credentials-secret"
    aws_secret_key            = "***REDACTED***"
    breaker_rulesets = [
      "s3-breaker",
    ]
    checkpointing = {
      enabled = true
      retries = 3
    }
    connections = [
      {
        output   = "s3-main"
        pipeline = "default"
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
    max_messages           = 10
    metadata = [
      {
        name  = "source"
        value = "\"s3\""
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 300
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    poll_timeout                   = 10
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    preprocess = {
      args = [
        "--compact-output",
      ]
      command  = "jq -r .message"
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
      enabled = true
      retries = 3
    }
    checksum_suffix = "checksum"
    connections = [
      {
        output   = "s3-archive"
        pipeline = "default"
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
    max_manifest_size_kb   = 4096
    max_messages           = 10
    metadata = [
      {
        name  = "source"
        value = "\"s3_inventory\""
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 300
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    poll_timeout                   = 10
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    preprocess = {
      args = [
        "--compact-output",
      ]
      command  = "jq -r .message"
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
    tag_after_processing     = "...my_tag_after_processing..."
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
      retries = 3
    }
    connections = [
      {
        output   = "s3-security-lake"
        pipeline = "default"
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
    max_messages           = 10
    metadata = [
      {
        name  = "source"
        value = "\"security_lake\""
      }
    ]
    num_receivers                  = 4
    parquet_chunk_download_timeout = 300
    parquet_chunk_size_mb          = 10
    pipeline                       = "default"
    poll_timeout                   = 10
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    preprocess = {
      args = [
        "--compact-output",
      ]
      command  = "jq -r .message"
      disabled = true
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
    tag_after_processing = "...my_tag_after_processing..."
    type                 = "security_lake"
    visibility_timeout   = 300
  }
  input_snmp = {
    best_effort_parsing = true
    connections = [
      {
        output   = "s3-network"
        pipeline = "default"
      }
    ]
    description        = "Receive SNMP traps and forward to destinations"
    disabled           = false
    environment        = "main"
    host               = "0.0.0.0"
    id                 = "snmp-traps"
    ip_whitelist_regex = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_buffer_size    = 20000
    metadata = [
      {
        name  = "source"
        value = "\"snmp\""
      }
    ]
    pipeline = "default"
    port     = 162
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled     = false
    send_to_routes = true
    snmp_v3_auth = {
      allow_unmatched_trap = false
      v3_auth_enabled      = true
      v3_users = [
        {
          auth_key      = "{ \"see\": \"documentation\" }"
          auth_protocol = "sha256"
          name          = "snmp-user"
          priv_protocol = "none"
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
    ip_whitelist_regex  = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_active_cxn      = 2000
    max_s2_sversion     = "v4"
    metadata = [
      {
        name  = "source"
        value = "\"splunk\""
      }
    ]
    pipeline = "default"
    port     = 8089
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "splunk-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:splunk_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
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
        auth_type   = "manual"
        description = "Token for HEC webhooks"
        enabled     = true
        metadata = [
          {
            name  = "sourcetype"
            value = "\"http:hec\""
          }
        ]
        token        = "{ \"see\": \"documentation\" }"
        token_secret = "{ \"see\": \"documentation\" }"
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
    enable_health_check     = "{ \"see\": \"documentation\" }"
    enable_proxy_header     = false
    environment             = "main"
    extract_metrics         = true
    host                    = "0.0.0.0"
    id                      = "splunk-hec-listener"
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.0\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"hec\""
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "splunk-hec-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:hec_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type             = "splunk_hec"
    use_fwd_timezone = true
  }
  input_splunk_search = {
    auth_header_expr = "`Bearer ${token}`"
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
        value = "\"Bearer $${{secret:splunk_token}\""
      }
    ]
    endpoint_params = [
      {
        name  = "earliest_time"
        value = "${earliest}"
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
        name  = "source"
        value = "\"splunk_search\""
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
    output_mode = "json"
    password    = "$${{secret:splunk_password}"
    pipeline    = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled          = false
    reject_unauthorized = false
    request_timeout     = 120
    retry_rules = {
      codes = [
        429,
        503,
        502,
      ]
      enable_header         = true
      interval              = 1000
      limit                 = 5
      multiplier            = 2
      retry_connect_reset   = true
      retry_connect_timeout = true
      type                  = "backoff"
    }
    search                 = "search index=main sourcetype=access_combined status>=500 | stats count by host"
    search_head            = "https://splunk.example.com:8089"
    secret                 = "$${{secret:splunk_oauth_secret}"
    secret_param_name      = "password"
    send_to_routes         = true
    stale_channel_flush_ms = 15000
    streamtags = [
      "prod",
      "splunk",
    ]
    text_secret          = "splunk-token-secret"
    token                = "$${{secret:splunk_token}"
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
        output   = "s3-raw"
        pipeline = "default"
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
        name  = "source"
        value = "\"sqs\""
      }
    ]
    num_receivers = 4
    pipeline      = "default"
    poll_timeout  = 10
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
    input_syslog_syslog2 = {
      allow_non_standard_app_name = true
      connections = [
        {
          output   = "s3-syslog"
          pipeline = "default"
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
      ip_whitelist_regex                   = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
      keep_fields_list = [
        "host",
        "app",
      ]
      max_active_cxn  = 2000
      max_buffer_size = 20000
      metadata = [
        {
          name  = "source"
          value = "\"syslog\""
        }
      ]
      octet_counting = false
      pipeline       = "default"
      pq = {
        commit_frequency = 100
        compress         = "gzip"
        max_buffer_size  = 5000
        max_file_size    = "100 MB"
        max_size         = "10GB"
        mode             = "always"
        path             = "/opt/cribl/state/queues"
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
        ca_path             = "/etc/ssl/certs/ca-bundle.crt"
        cert_path           = "/etc/ssl/certs/server.crt"
        certificate_name    = "syslog-cert"
        common_name_regex   = "{ \"see\": \"documentation\" }"
        disabled            = true
        max_version         = "TLSv1.3"
        min_version         = "TLSv1.2"
        passphrase          = "***REDACTED***"
        priv_key_path       = "/etc/ssl/private/server.key"
        reject_unauthorized = "{ \"see\": \"documentation\" }"
        request_cert        = false
      }
      type                   = "syslog"
      udp_port               = 514
      udp_socket_rx_buf_size = 4194304
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
          mode = "all"
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
        name  = "source"
        value = "\"system_metrics\""
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
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    process = {
      sets = [
        {
          filter           = "proc.name == 'nginx'"
          include_children = true
          name             = "nginx-workers"
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
    description           = "Collect system state metrics and spool to disk"
    disable_native_module = false
    disabled              = false
    environment           = "main"
    id                    = "system-state"
    interval              = 600
    metadata = [
      {
        name  = "source"
        value = "\"system_state\""
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
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
    auth_type = "manual"
    breaker_rulesets = [
      "multiline-json",
      "tcp-syslog-breaker",
    ]
    connections = [
      {
        output   = "s3-raw"
        pipeline = "default"
      }
    ]
    description         = "Receive generic TCP payloads"
    disabled            = false
    enable_header       = false
    enable_proxy_header = false
    environment         = "main"
    host                = "0.0.0.0"
    id                  = "tcp-listener"
    ip_whitelist_regex  = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_active_cxn      = 2000
    metadata = [
      {
        name  = "source"
        value = "\"tcp\""
      }
    ]
    pipeline = "default"
    port     = 9000
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    preprocess = {
      args = [
        "--compact-output",
      ]
      command  = "jq -r .message"
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
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "tcp-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "tcp"
  }
  input_tcpjson = {
    auth_token = "$${{secret:tcpjson_token}"
    auth_type  = "manual"
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
    ip_whitelist_regex    = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    max_active_cxn        = 2000
    metadata = [
      {
        name  = "source"
        value = "\"tcpjson\""
      }
    ]
    pipeline = "default"
    port     = 9001
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "tcpjson-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:tcpjson_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
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
        output   = "s3-wef"
        pipeline = "default"
      }
    ]
    description              = "Receive Windows Event Forwarding (WEF) over HTTPS"
    disabled                 = false
    enable_health_check      = true
    enable_proxy_header      = true
    environment              = "main"
    host                     = "0.0.0.0"
    id                       = "wef-listener"
    ip_allowlist_regex       = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex        = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout       = 60
    keytab                   = "/etc/krb5.keytab"
    log_fingerprint_mismatch = true
    max_active_req           = 512
    max_requests_per_socket  = 0
    metadata = [
      {
        name  = "source"
        value = "\"wef\""
      }
    ]
    pipeline = "default"
    port     = 5986
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
            name  = "subscription"
            value = "\"wef-security\""
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
        version = "7f0c2f2e-1c3b-4d2a-9d6e-5a1b2c3d4e5f"
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
      min_version           = "TLSv1.2"
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
        output   = "s3-logs"
        pipeline = "default"
      }
    ]
    description           = "Collect Windows Event Logs from local system"
    disable_native_module = false
    disabled              = false
    environment           = "main"
    event_format          = "json"
    id                    = "win-event-logs"
    interval              = 10
    log_names = [
      "Application",
      "Security",
      "System",
      "Microsoft-Windows-Sysmon/Operational",
    ]
    max_event_bytes = 131072
    metadata = [
      {
        name  = "source"
        value = "\"win_event_logs\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
        output   = "s3-main"
        pipeline = "default"
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
        }
        system = {
          detail = true
          mode   = "basic"
        }
      }
      mode = "basic"
    }
    id       = "windows-metrics"
    interval = 10
    metadata = [
      {
        name  = "source"
        value = "\"windows_metrics\""
      }
    ]
    persistence = {
      compress      = "gzip"
      dest_path     = "/var/lib/cribl/state/windows_metrics"
      enable        = true
      max_data_size = "4GB"
      max_data_time = "4d"
      time_window   = "10m"
    }
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled = false
    process = {
      sets = [
        {
          filter           = "name == \"w3wp.exe\""
          include_children = true
          name             = "IIS worker processes"
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
        output   = "s3-wiz"
        pipeline = "default"
      }
    ]
    content_config = [
      {
        content_description = "...my_content_description..."
        content_type        = "...my_content_type..."
        enabled             = true
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
        name  = "source"
        value = "\"wiz\""
      }
    ]
    pipeline = "default"
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled      = false
    request_timeout = 300
    retry_rules = {
      codes = [
        429,
        503,
      ]
      enable_header         = true
      interval              = 1000
      limit                 = 5
      multiplier            = 2
      retry_connect_reset   = false
      retry_connect_timeout = false
      type                  = "backoff"
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
        auth_type   = "manual"
        description = "Zscaler Collector"
        enabled     = true
        metadata = [
          {
            name  = "source"
            value = "\"zscaler_hec\""
          }
        ]
        token        = "{ \"see\": \"documentation\" }"
        token_secret = "{ \"see\": \"documentation\" }"
      }
    ]
    capture_headers = true
    connections = [
      {
        output   = "s3-zscaler"
        pipeline = "default"
      }
    ]
    description             = "Receive Zscaler HEC events over HTTP(S)"
    disabled                = false
    emit_token_metrics      = true
    enable_health_check     = "{ \"see\": \"documentation\" }"
    enable_proxy_header     = false
    environment             = "main"
    hec_acks                = false
    hec_api                 = "/services/collector"
    host                    = "0.0.0.0"
    id                      = "zscaler-hec-listener"
    ip_allowlist_regex      = "^10\\.0\\.\\d{1,3}\\.\\d{1,3}$"
    ip_denylist_regex       = "^192\\.168\\.1\\.\\d{1,3}$"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 0
    metadata = [
      {
        name  = "source"
        value = "\"zscaler_hec\""
      }
    ]
    pipeline = "default"
    port     = 8088
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "100 MB"
      max_size         = "10GB"
      mode             = "always"
      path             = "/opt/cribl/state/queues"
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
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "zscaler-hec-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "***REDACTED***"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "zscaler_hec"
  }
}