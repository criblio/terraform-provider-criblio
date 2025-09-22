# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract destination type
  destination_type = local.raw_config.type

  # Destination ID - use override if provided, otherwise from JSON
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to kinesis (comprehensive from schema)
  field_mappings = {
    # Kinesis specific mappings
    "asNdjson"                = "as_ndjson"
    "assumeRoleArn"           = "assume_role_arn"
    "assumeRoleExternalId"    = "assume_role_external_id"
    "awsApiKey"               = "aws_api_key"
    "awsAuthenticationMethod" = "aws_authentication_method"
    "awsSecret"               = "aws_secret"
    "awsSecretKey"            = "aws_secret_key"
    "durationSeconds"         = "duration_seconds"
    "enableAssumeRole"        = "enable_assume_role"
    "flushPeriodSec"          = "flush_period_sec"
    "maxRecordSizeKb"         = "max_record_size_kb"
    "maxRecordSizeKB"         = "max_record_size_kb" # Handle both cases
    "onBackpressure"          = "on_backpressure"
    "partitionKey"            = "partition_key"
    "pqCompress"              = "pq_compress"
    "pqControls"              = "pq_controls"
    "pqMaxFileSize"           = "pq_max_file_size"
    "pqMaxSize"               = "pq_max_size"
    "pqMode"                  = "pq_mode"
    "pqOnBackpressure"        = "pq_on_backpressure"
    "pqPath"                  = "pq_path"
    "rejectUnauthorized"      = "reject_unauthorized"
    "reuseConnections"        = "reuse_connections"
    "signatureVersion"        = "signature_version"
    "streamName"              = "stream_name"
    "systemFields"            = "system_fields"
    "textSecret"              = "text_secret"
    "useAssumeRoleExternalId" = "use_assume_role_external_id"
    "useListShards"           = "use_list_shards"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "disabled", "environment", "pipeline", "streamtags", "region", "compression", "concurrency", "endpoint", "type"]

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
