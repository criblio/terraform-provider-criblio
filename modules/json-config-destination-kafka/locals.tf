# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract destination type
  destination_type = local.raw_config.type

  # Destination ID - use override if provided, otherwise from JSON
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to kafka (comprehensive from schema)
  field_mappings = {
    # Kafka specific mappings
    "authenticationTimeout"    = "authentication_timeout"
    "backoffRate"              = "backoff_rate"
    "connectionTimeout"        = "connection_timeout"
    "flushEventCount"          = "flush_event_count"
    "flushPeriodSec"           = "flush_period_sec"
    "initialBackoff"           = "initial_backoff"
    "kafkaSchemaRegistry"      = "kafka_schema_registry"
    "maxBackOff"               = "max_back_off"
    "maxRecordSizeKb"          = "max_record_size_kb"
    "maxRecordSizeKB"          = "max_record_size_kb" # Handle both cases
    "maxRetries"               = "max_retries"
    "onBackpressure"           = "on_backpressure"
    "pqCompress"               = "pq_compress"
    "pqControls"               = "pq_controls"
    "pqMaxFileSize"            = "pq_max_file_size"
    "pqMaxSize"                = "pq_max_size"
    "pqMode"                   = "pq_mode"
    "pqOnBackpressure"         = "pq_on_backpressure"
    "pqPath"                   = "pq_path"
    "protobufLibraryId"        = "protobuf_library_id"
    "protobufMessageType"      = "protobuf_message_type"
    "protobufTypeNameStrategy" = "protobuf_type_name_strategy"
    "requestTimeout"           = "request_timeout"
    "saslMechanism"            = "sasl_mechanism"
    "schemaRegistryCompress"   = "schema_registry_compress"
    "schemaRegistryFormat"     = "schema_registry_format"
    "sendBatch"                = "send_batch"
    "systemFields"             = "system_fields"

    # SASL nested fields
    "username"          = "username"
    "password"          = "password"
    "authToken"         = "auth_token"
    "authType"          = "auth_type"
    "credentialsSecret" = "credentials_secret"
    "textSecret"        = "text_secret"

    # Schema registry nested fields
    "schemaRegistryUrl"    = "schema_registry_url"
    "defaultKeySchemaId"   = "default_key_schema_id"
    "defaultValueSchemaId" = "default_value_schema_id"
    "connectionTimeout"    = "connection_timeout"
    "requestTimeout"       = "request_timeout"
    "maxRetries"           = "max_retries"

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
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "disabled", "environment", "pipeline", "streamtags", "brokers", "topic", "ack", "compression", "format"]

  # First pass: transform top-level keys (only include known fields)
  base_config = {
    for key, value in local.raw_config :
    (contains(keys(local.field_mappings), key) ? local.field_mappings[key] : key) => value
    if key != "type" && (contains(keys(local.field_mappings), key) || contains(local.direct_fields, key))
  }

  # Handle nested objects for SASL, Schema Registry, and TLS transformations
  # This would be expanded if needed for complex nested transformations

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
