# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract destination type
  destination_type = local.raw_config.type

  # Destination ID - use override if provided, otherwise from JSON
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to cribl_http (comprehensive from schema)
  field_mappings = {
    # Cribl HTTP specific mappings
    "dnsResolvePeriodSec"           = "dns_resolve_period_sec"
    "excludeFields"                 = "exclude_fields"
    "excludeSelf"                   = "exclude_self"
    "extraHttpHeaders"              = "extra_http_headers"
    "failedRequestLoggingMode"      = "failed_request_logging_mode"
    "flushPeriodSec"                = "flush_period_sec"
    "loadBalanceStatsPeriodSec"     = "load_balance_stats_period_sec"
    "loadBalanced"                  = "load_balanced"
    "maxPayloadEvents"              = "max_payload_events"
    "maxPayloadSizeKb"              = "max_payload_size_kb"
    "maxPayloadSizeKB"              = "max_payload_size_kb" # Handle both cases
    "onBackpressure"                = "on_backpressure"
    "pqCompress"                    = "pq_compress"
    "pqControls"                    = "pq_controls"
    "pqMaxFileSize"                 = "pq_max_file_size"
    "pqMaxSize"                     = "pq_max_size"
    "pqMode"                        = "pq_mode"
    "pqOnBackpressure"              = "pq_on_backpressure"
    "pqPath"                        = "pq_path"
    "rejectUnauthorized"            = "reject_unauthorized"
    "responseHonorRetryAfterHeader" = "response_honor_retry_after_header"
    "responseRetrySettings"         = "response_retry_settings"
    "safeHeaders"                   = "safe_headers"
    "systemFields"                  = "system_fields"
    "timeoutRetrySettings"          = "timeout_retry_settings"
    "timeoutSec"                    = "timeout_sec"
    "useRoundRobinDns"              = "use_round_robin_dns"

    # Auth nested fields
    "authType"          = "auth_type"
    "username"          = "username"
    "password"          = "password"
    "credentialsSecret" = "credentials_secret"
    "textSecret"        = "text_secret"
    "disabled"          = "disabled"

    # TLS nested fields
    "caPath"          = "ca_path"
    "certPath"        = "cert_path"
    "certificateName" = "certificate_name"
    "maxVersion"      = "max_version"
    "minVersion"      = "min_version"
    "passphrase"      = "passphrase"
    "privKeyPath"     = "priv_key_path"
    "servername"      = "servername"

    # Response retry settings nested fields
    "backoffRate"    = "backoff_rate"
    "httpStatus"     = "http_status"
    "initialBackoff" = "initial_backoff"
    "maxBackoff"     = "max_backoff"

    # Timeout retry settings nested fields
    "timeoutRetry" = "timeout_retry"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "disabled", "environment", "pipeline", "streamtags", "type", "compression", "concurrency", "url", "urls"]

  # First pass: transform top-level keys (only include known fields)
  base_config = {
    for key, value in local.raw_config :
    (contains(keys(local.field_mappings), key) ? local.field_mappings[key] : key) => value
    if key != "type" && (contains(keys(local.field_mappings), key) || contains(local.direct_fields, key))
  }

  # Final configuration with additional config merged
  final_config = merge(
    {
      id   = local.destination_id
      type = local.destination_type
    },
    local.base_config,
    var.additional_config
  )
}
