variable "source_id" {
  description = "Unique identifier for the source"
  type        = string
}

variable "group_id" {
  description = "Worker group ID"
  type        = string
}

variable "source_type" {
  description = "Type of source"
  type        = string
  validation {
    condition = contains([
      "syslog",
      "cribl_http", # Cribl HTTP receiver
      "cribl_tcp",  # Cribl TCP receiver  
      "http",       # Regular HTTP source
      "tcp",        # Regular TCP source
      "otlp",       # OpenTelemetry
      # Add more as needed
    ], var.source_type)
    error_message = "Invalid source type. Must be one of: syslog, cribl_http, cribl_tcp, http, tcp, otlp"
  }
}

variable "port" {
  description = "Port number (required) - Cannot use default ports"
  type        = number
  validation {
    condition = (
      var.port != 10200 && # cribl_http default
      var.port != 9514 &&  # syslog default
      var.port != 10001 && # cribl_tcp default
      var.port != 10080 && # http default (assumed)
      var.port != 10060 && # tcp default (assumed)
      var.port != 4317 &&  # OTLP gRPC default
      var.port != 4318     # OTLP HTTP default
    )
    error_message = "Port cannot be a default port. Please choose a different port."
  }
}

variable "description" {
  description = "Source description"
  type        = string
  default     = ""
}

variable "disabled" {
  description = "Disable this source"
  type        = bool
  default     = false
}

variable "connections" {
  description = "Output connections"
  type = list(object({
    output   = string
    pipeline = string
  }))
  default = []
}

variable "pipeline" {
  description = "Default pipeline for this source"
  type        = string
  default     = null
}

variable "pq_enabled" {
  description = "Enable persistent queue"
  type        = bool
  default     = false
}

variable "streamtags" {
  description = "Stream tags"
  type        = list(string)
  default     = []
}

variable "custom_config" {
  description = "Custom configuration to override defaults"
  type        = any
  default     = {}
}