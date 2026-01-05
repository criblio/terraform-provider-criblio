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

# Chronicle destination
resource "criblio_destination" "my_chronicle_destination" {
  group_id = "default"
  id       = "chronicle_prod"
  output_chronicle = {
    id                                 = "chronicle_prod"
    type                               = "chronicle"
    description                        = "Send events to Google Chronicle"
    authentication_method              = "serviceAccountSecret"
    compress                           = true
    concurrency                        = 8
    gcp_project_id                     = "my-gcp-project"
    gcp_instance                       = "123e4567-e89b-12d3-a456-426614174000"
    log_type                           = "UNSTRUCTURED_DATA"
    ingestion_method                   = "BATCH"
    region                             = "us"
    on_backpressure                    = "block"
    pipeline                           = "default"
    flush_period_sec                   = 2
    max_payload_events                 = 1000
    max_payload_size_kb                = 2048
    timeout_sec                        = 30
    reject_unauthorized                = true
    service_account_credentials_secret = "chronicle-service-account"
  }
}

# Cloudflare R2 destination
resource "criblio_destination" "my_cloudflare_r2_destination" {
  group_id = "default"
  id       = "cloudflare_r2_prod"
  output_cloudflare_r2 = {
    id                        = "cloudflare_r2_prod"
    type                      = "cloudflare_r2"
    description               = "Write objects to Cloudflare R2"
    bucket                    = "my-r2-bucket"
    endpoint                  = "https://my-account-id.r2.cloudflarestorage.com"
    aws_authentication_method = "auto"
    compress                  = "gzip"
    format                    = "json"
    dest_path                 = "logs/ingest"
    on_backpressure           = "block"
    pipeline                  = "default"
    max_file_size_mb          = 64
    max_retry_num             = 10
  }
}



# Databricks destination
resource "criblio_destination" "my_databricks_destination" {
  group_id = "default"
  id       = "databricks_prod"
  output_databricks = {
    id                 = "databricks_prod"
    type               = "databricks"
    description        = "Write data to Databricks"
    workspace_id       = "https://my-workspace.cloud.databricks.com"
    catalog            = "main"
    schema             = "default"
    events_volume_name = "events_volume"
    compress           = "gzip"
    format             = "json"
    dest_path          = "logs/ingest"
    on_backpressure    = "block"
    pipeline           = "default"
    max_file_size_mb   = 64
    max_retry_num      = 10
    client_id          = "databricks-client-id"
    client_text_secret = "databricks-client-secret"
    scope              = "databricks-scope"
  }
}

# Microsoft Fabric destination
resource "criblio_destination" "my_microsoft_fabric_destination" {
  group_id = "default"
  id       = "microsoft_fabric_prod"
  output_microsoft_fabric = {
    id                 = "microsoft_fabric_prod"
    type               = "microsoft_fabric"
    description        = "Produce events to Microsoft Fabric"
    bootstrap_server   = "my-workspace.servicebus.windows.net:9093"
    topic              = "app-events"
    format             = "json"
    on_backpressure    = "block"
    pipeline           = "default"
    flush_period_sec   = 2
    flush_event_count  = 1000
    max_record_size_kb = 1024
    ack                = 1
    max_retries        = 5
    sasl = {
      disabled           = false
      mechanism          = "oauthbearer"
      client_id          = "fabric-client-id"
      client_text_secret = "fabric-client-secret"
      tenant_id          = "tenant-id"
      oauth_endpoint     = "https://login.microsoftonline.com"
      scope              = "https://servicebus.azure.net/.default"
    }
    tls = {
      disabled            = false
      reject_unauthorized = true
    }
  }
}

# SentinelOne AI SIEM destination
resource "criblio_destination" "my_sentinel_one_ai_siem_destination" {
  group_id = "default"
  id       = "sentinel_one_ai_siem_prod"
  output_sentinel_one_ai_siem = {
    id                   = "sentinel_one_ai_siem_prod"
    type                 = "sentinel_one_ai_siem"
    description          = "Send events to SentinelOne AI SIEM"
    base_url             = "https://api.sentinelone.com"
    endpoint             = "/services/collector/raw"
    auth_type            = "manual"
    token                = "sentinelone-token"
    compress             = false
    concurrency          = 8
    on_backpressure      = "block"
    pipeline             = "default"
    flush_period_sec     = 2
    max_payload_events   = 1000
    max_payload_size_kb  = 2048
    timeout_sec          = 30
    reject_unauthorized  = true
    region               = "US"
    data_source_vendor   = "Cribl"
    data_source_name     = "Cribl Stream"
    data_source_category = "Security"
    event_type           = "log"
    source               = "cribl"
    source_type          = "stream"
  }
}

/*
data "criblio_destinations" "my_destinations" {
  group_id = "default"
}

output "my_destinations" {
  value = data.criblio_destinations.my_destinations
} 
*/