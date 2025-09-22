# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract destination type
  destination_type = local.raw_config.type

  # Destination ID - use override if provided, otherwise from JSON
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to cribl_tcp (comprehensive from schema)
  field_mappings = {
    # Cribl TCP specific mappings
    "connectionTimeout"         = "connection_timeout"
    "dnsResolvePeriodSec"       = "dns_resolve_period_sec"
    "excludeFields"             = "exclude_fields"
    "excludeSelf"               = "exclude_self"
    "loadBalanceStatsPeriodSec" = "load_balance_stats_period_sec"
    "loadBalanced"              = "load_balanced"
    "logFailedRequests"         = "log_failed_requests"
    "maxConcurrentSenders"      = "max_concurrent_senders"
    "onBackpressure"            = "on_backpressure"
    "pqCompress"                = "pq_compress"
    "pqControls"                = "pq_controls"
    "pqMaxFileSize"             = "pq_max_file_size"
    "pqMaxSize"                 = "pq_max_size"
    "pqMode"                    = "pq_mode"
    "pqOnBackpressure"          = "pq_on_backpressure"
    "pqPath"                    = "pq_path"
    "systemFields"              = "system_fields"
    "throttleRatePerSec"        = "throttle_rate_per_sec"

    # TLS nested fields
    "caPath"             = "ca_path"
    "certPath"           = "cert_path"
    "certificateName"    = "certificate_name"
    "disabled"           = "disabled"
    "maxVersion"         = "max_version"
    "minVersion"         = "min_version"
    "passphrase"         = "passphrase"
    "privKeyPath"        = "priv_key_path"
    "rejectUnauthorized" = "reject_unauthorized"
    "servername"         = "servername"

    # Hosts nested fields
    "weight" = "weight"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "disabled", "environment", "pipeline", "streamtags", "type", "compression", "host", "hosts", "port", "tls"]

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
