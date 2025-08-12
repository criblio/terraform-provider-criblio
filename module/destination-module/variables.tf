variable "destination_id" {
  description = "Unique identifier for the destination"
  type        = string
}

variable "group_id" {
  description = "Worker group ID"
  type        = string
}

variable "destination_type" {
  description = "Type of destination"
  type        = string
  validation {
    condition = contains([
      "cribl_lake",
      "s3",
      "splunk_hec",
      "cribl_http",
      "cribl_tcp",
      "syslog",
      "kafka",
      "kinesis",
      "elasticsearch",
      "webhook",
      # Add more as needed
    ], var.destination_type)
    error_message = "Invalid destination type. Must be one of: cribl_lake, s3, splunk_hec, cribl_http, cribl_tcp, syslog, kafka, kinesis, elasticsearch, webhook"
  }
}

variable "description" {
  description = "Destination description"
  type        = string
  default     = ""
}

variable "disabled" {
  description = "Disable this destination"
  type        = bool
  default     = false
}

variable "streamtags" {
  description = "Stream tags"
  type        = list(string)
  default     = []
}

variable "pipeline" {
  description = "Pipeline for this destination"
  type        = string
  default     = null
}

variable "custom_config" {
  description = "Custom configuration to override defaults"
  type        = any
  default     = {}
}

# Cribl Lake specific variables
variable "dest_path" {
  description = "Destination path for Cribl Lake"
  type        = string
  default     = "default_logs"
}

variable "format" {
  description = "Output format (json, ndjson, parquet, etc.)"
  type        = string
  default     = "json"
}

variable "compress" {
  description = "Compression type (gzip, none, etc.)"
  type        = string
  default     = "gzip"
}

# S3 specific variables
variable "bucket" {
  description = "S3 bucket name"
  type        = string
  default     = ""
}

variable "region" {
  description = "AWS region for S3"
  type        = string
  default     = "us-east-1"
}

# HTTP specific variables
variable "url" {
  description = "HTTP URL for HTTP destinations"
  type        = string
  default     = ""
}

variable "method" {
  description = "HTTP method"
  type        = string
  default     = "POST"
}

# Splunk specific variables
variable "splunk_url" {
  description = "Splunk HEC URL"
  type        = string
  default     = ""
}

variable "token" {
  description = "Authentication token"
  type        = string
  default     = ""
  sensitive   = true
}