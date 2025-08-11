# modules/cribl-source/variables.tf - COMPLETE FILE
variable "source_id" {
  description = "Unique identifier for the source"
  type        = string
}

variable "group_id" {
  description = "Worker group ID"
  type        = string
}

variable "source_type" {
  description = "Type of source (syslog, cribl_http, tcp, etc.)"
  type        = string
  validation {
    condition = contains([
      "syslog",
      "cribl_http",
      # Add more as needed
    ], var.source_type)
    error_message = "Invalid source type."
  }
}

variable "port" {
  description = "Port number (required) - For syslog, this sets both TCP and UDP ports. Cannot use default ports (9514 for syslog, 10200 for HTTP)"
  type        = number
  validation {
    condition = (
      var.port != 10200 && # HTTP default
      var.port != 9514     # Syslog default
    )
    error_message = "Port cannot be the default port (10200 for HTTP, 9514 for syslog). Please choose a different port."
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