# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract destination type
  destination_type = local.raw_config.type

  # Destination ID - use override if provided, otherwise from JSON
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to open_telemetry (comprehensive from schema)
  field_mappings = {
    # OpenTelemetry specific mappings
    "authHeaderExpr"                = "auth_header_expr"
    "authType"                      = "auth_type"
    "connectionTimeout"             = "connection_timeout"
    "credentialsSecret"             = "credentials_secret"
    "extraHttpHeaders"              = "extra_http_headers"
    "failedRequestLoggingMode"      = "failed_request_logging_mode"
    "flushPeriodSec"                = "flush_period_sec"
    "httpCompress"                  = "http_compress"
    "httpLogsEndpointOverride"      = "http_logs_endpoint_override"
    "httpMetricsEndpointOverride"   = "http_metrics_endpoint_override"
    "httpTracesEndpointOverride"    = "http_traces_endpoint_override"
    "keepAlive"                     = "keep_alive"
    "keepAliveTime"                 = "keep_alive_time"
    "loginUrl"                      = "login_url"
    "maxPayloadSizeKb"              = "max_payload_size_kb"
    "maxPayloadSizeKB"              = "max_payload_size_kb" # Handle both cases
    "oauthHeaders"                  = "oauth_headers"
    "oauthParams"                   = "oauth_params"
    "onBackpressure"                = "on_backpressure"
    "otlpVersion"                   = "otlp_version"
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
    "textSecret"                    = "text_secret"
    "timeoutRetrySettings"          = "timeout_retry_settings"
    "timeoutSec"                    = "timeout_sec"

    # Nested object field mappings
    "backoffRate"    = "backoff_rate"
    "httpStatus"     = "http_status"
    "initialBackoff" = "initial_backoff"
    "maxBackoff"     = "max_backoff"
    "timeoutRetry"   = "timeout_retry"
    "disabled"       = "disabled"
    "name"           = "name"
    "value"          = "value"
    "key"            = "key"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "disabled", "environment", "pipeline", "streamtags", "type", "endpoint", "compress", "concurrency", "metadata", "password", "protocol", "username"]

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
