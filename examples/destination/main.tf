locals {
  # Base Cribl HTTP configuration
  cribl_http_config = {
    id          = "cribl_http_prod"
    type        = "cribl_http"
    description = "Send events to Cribl Worker HTTP endpoint with retries"
    disabled    = false
    streamtags  = ["prod", "cribl"]
    pipeline    = "default"

    # Cribl HTTP specific settings
    method                = "POST"
    headers               = {}
    compress              = true
    format                = "json"
    flush_period_sec      = 2
    max_events_per_req    = 1000
    on_backpressure       = "block"
    timeout               = 30
    honor_keep_alive      = true
    request_concurrency   = 8
    max_retry_num         = 10
    validate_server_certs = true
    auth_type             = "none"
    load_balanced         = true
    url                   = "https://edge.example.com:10200"
    urls = [
      {
        url    = "https://edge01.example.com:10200"
        weight = 2
      }
    ]
    use_round_robin_dns = true

    # Additional settings
    compression            = "gzip"
    concurrency            = 8
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
    load_balance_stats_period_sec = 300
    max_payload_size_kb           = 2048
    pq_compress                   = "gzip"
    pq_controls = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
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
    token_ttl_minutes = 60
    tls = {
      ca_path             = "/etc/ssl/certs/ca-bundle.crt"
      cert_path           = "/opt/cribl/certs/client.crt"
      certificate_name    = "cribl-client"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "your_passphrase_load_it_from_the_secrets_manager_or_env_variable"
      priv_key_path       = "/opt/cribl/certs/client.key"
      reject_unauthorized = true
      servername          = "collector.cribl.example.com"
    }
  }
}

resource "criblio_destination" "my_cribl_http_destination" {
  group_id          = "default"
  id                = "cribl_http_prod"
  output_cribl_http = local.cribl_http_config
}

/*
data "criblio_destinations" "my_destinations" {
  group_id = "default"
}

output "my_destinations" {
  value = data.criblio_destinations.my_destinations
} 
*/