# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract source type
  source_type = local.raw_config.type

  # Source ID - use override if provided, otherwise from JSON
  source_id = var.source_id_override != "" ? var.source_id_override : local.raw_config.id

  # Field mappings specific to tcp
  field_mappings = {
    # TCP specific mappings
    "pqEnabled"           = "pq_enabled"
    "sendToRoutes"        = "send_to_routes"
    "maxActiveCxn"        = "max_active_cxn"
    "socketIdleTimeout"   = "socket_idle_timeout"
    "socketMaxLifespan"   = "socket_max_lifespan"
    "socketEndingMaxWait" = "socket_ending_max_wait"
    "staleChannelFlushMs" = "stale_channel_flush_ms"
    "enableHeader"        = "enable_header"
    "enableProxyHeader"   = "enable_proxy_header"
    "ipWhitelistRegex"    = "ip_whitelist_regex"
    "authType"            = "auth_type"
    "breakerRulesets"     = "breaker_rulesets"

    # TLS specific mappings
    "caPath"             = "ca_path"
    "certPath"           = "cert_path"
    "certificateName"    = "certificate_name"
    "maxVersion"         = "max_version"
    "minVersion"         = "min_version"
    "privKeyPath"        = "priv_key_path"
    "rejectUnauthorized" = "reject_unauthorized"
    "requestCert"        = "request_cert"

    # PQ mappings
    "commitFrequency" = "commit_frequency"
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
