# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract source type
  source_type = local.raw_config.type

  # Source ID - use override if provided, otherwise from JSON
  source_id = var.source_id_override != "" ? var.source_id_override : local.raw_config.id

  # Field mappings specific to http (comprehensive from schema)
  field_mappings = {
    # HTTP specific mappings
    "activityLogSampleRate" = "activity_log_sample_rate"
    "authTokens"            = "auth_tokens"
    "authTokensExt"         = "auth_tokens_ext"
    "captureHeaders"        = "capture_headers"
    "criblApi"              = "cribl_api"
    "elasticApi"            = "elastic_api"
    "enableHealthCheck"     = "enable_health_check"
    "enableProxyHeader"     = "enable_proxy_header"
    "ipAllowlistRegex"      = "ip_allowlist_regex"
    "ipDenylistRegex"       = "ip_denylist_regex"
    "keepAliveTimeout"      = "keep_alive_timeout"
    "maxActiveReq"          = "max_active_req"
    "maxRequestsPerSocket"  = "max_requests_per_socket"
    "pqEnabled"             = "pq_enabled"
    "requestTimeout"        = "request_timeout"
    "sendToRoutes"          = "send_to_routes"
    "socketTimeout"         = "socket_timeout"
    "splunkHecAcks"         = "splunk_hec_acks"
    "splunkHecApi"          = "splunk_hec_api"

    # TLS specific mappings
    "caPath"             = "ca_path"
    "certPath"           = "cert_path"
    "certificateName"    = "certificate_name"
    "commonNameRegex"    = "common_name_regex"
    "maxVersion"         = "max_version"
    "minVersion"         = "min_version"
    "passphrase"         = "passphrase"
    "privKeyPath"        = "priv_key_path"
    "rejectUnauthorized" = "reject_unauthorized"
    "requestCert"        = "request_cert"

    # PQ mappings
    "commitFrequency" = "commit_frequency"
    "maxBufferSize"   = "max_buffer_size"
    "maxFileSize"     = "max_file_size"
    "maxSize"         = "max_size"
  }

  # First pass: transform top-level keys
  base_config = {
    for key, value in local.raw_config :
    lookup(local.field_mappings, key, key) => value
    if key != "type" # Exclude the type field
  }

  # Second pass: transform nested objects separately
  pq_transformed = can(local.base_config.pq) ? {
    for key, value in local.base_config.pq :
    lookup(local.field_mappings, key, key) => value
  } : null

  tls_transformed = can(local.base_config.tls) ? {
    for key, value in local.base_config.tls :
    lookup(local.field_mappings, key, key) => value
  } : null

  # Final terraform config with nested objects replaced
  terraform_config = merge(
    local.base_config,
    local.pq_transformed != null ? { pq = local.pq_transformed } : {},
    local.tls_transformed != null ? { tls = local.tls_transformed } : {}
  )

  # Final configuration with additional config merged
  final_config = merge(
    local.terraform_config,
    {
      id   = local.source_id
      type = local.source_type
    },
    var.additional_config
  )
}
