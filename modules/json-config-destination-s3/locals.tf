# Parse JSON configuration
locals {
  # Load JSON from file or content
  raw_config = var.json_file_path != "" ? jsondecode(file(var.json_file_path)) : jsondecode(var.json_content)

  # Extract destination type
  destination_type = local.raw_config.type

  # Destination ID - use override if provided, otherwise from JSON
  destination_id = var.destination_id_override != "" ? var.destination_id_override : local.raw_config.id

  # Field mappings specific to S3 (comprehensive from schema)
  field_mappings = {
    # S3 specific mappings
    "addIdToStagePath"              = "add_id_to_stage_path"
    "assumeRoleArn"                 = "assume_role_arn"
    "assumeRoleExternalId"          = "assume_role_external_id"
    "automaticSchema"               = "automatic_schema"
    "awsApiKey"                     = "aws_api_key"
    "awsAuthenticationMethod"       = "aws_authentication_method"
    "awsSecret"                     = "aws_secret"
    "awsSecretKey"                  = "aws_secret_key"
    "baseFileName"                  = "base_file_name"
    "compressionLevel"              = "compression_level"
    "deadletterEnabled"             = "deadletter_enabled"
    "deadletterPath"                = "deadletter_path"
    "destPath"                      = "dest_path"
    "durationSeconds"               = "duration_seconds"
    "emptyDirCleanupSec"            = "empty_dir_cleanup_sec"
    "enableAssumeRole"              = "enable_assume_role"
    "enablePageChecksum"            = "enable_page_checksum"
    "enableStatistics"              = "enable_statistics"
    "enableWritePageIndex"          = "enable_write_page_index"
    "fileNameSuffix"                = "file_name_suffix"
    "headerLine"                    = "header_line"
    "keyValueMetadata"              = "key_value_metadata"
    "kmsKeyId"                      = "kms_key_id"
    "maxClosingFilesToBackpressure" = "max_closing_files_to_backpressure"
    "maxConcurrentFileParts"        = "max_concurrent_file_parts"
    "maxFileIdleTimeSec"            = "max_file_idle_time_sec"
    "maxFileOpenTimeSec"            = "max_file_open_time_sec"
    "maxFileSizeMb"                 = "max_file_size_mb"
    "maxFileSizeMB"                 = "max_file_size_mb" # Handle both lowercase and uppercase B
    "maxOpenFiles"                  = "max_open_files"
    "maxRetryNum"                   = "max_retry_num"
    "objectAcl"                     = "object_acl"
    "objectACL"                     = "object_acl" # Handle both lowercase and uppercase ACL
    "onBackpressure"                = "on_backpressure"
    "onDiskFullBackpressure"        = "on_disk_full_backpressure"
    "parquetDataPageVersion"        = "parquet_data_page_version"
    "parquetPageSize"               = "parquet_page_size"
    "parquetRowGroupLength"         = "parquet_row_group_length"
    "parquetVersion"                = "parquet_version"
    "partitionExpr"                 = "partition_expr"
    "preprocess"                    = "preprocess"
    "region"                        = "region"
    "rejectUnauthorized"            = "reject_unauthorized"
    "removeEmptyDirs"               = "remove_empty_dirs"
    "reuseConnections"              = "reuse_connections"
    "roleSessionName"               = "role_session_name"
    "schemaRegistry"                = "schema_registry"
    "serverSideEncryption"          = "server_side_encryption"
    "signatureVersion"              = "signature_version"
    "stagePath"                     = "stage_path"
    "storageClass"                  = "storage_class"
    "systemFields"                  = "system_fields"
    "textSecret"                    = "text_secret"
    "useAssumeRoleExternalId"       = "use_assume_role_external_id"
    "verifyPermissions"             = "verify_permissions"
    "writeHighWaterMark"            = "write_high_water_mark"
  }

  # Known direct fields that don't need transformation
  direct_fields = ["id", "description", "disabled", "environment", "pipeline", "streamtags", "bucket", "region", "format", "compress"]

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
