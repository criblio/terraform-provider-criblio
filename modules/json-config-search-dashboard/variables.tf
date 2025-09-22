variable "json_file_path" {
  description = "Path to JSON file containing the dashboard configuration"
  type        = string
  default     = ""
}

variable "json_content" {
  description = "JSON content for the dashboard configuration (alternative to json_file_path)"
  type        = string
  default     = ""
}

variable "dashboard_id_override" {
  description = "Override the dashboard ID from the JSON configuration"
  type        = string
  default     = ""
}

variable "additional_config" {
  description = "Additional configuration to merge with the parsed JSON"
  type        = any
  default     = {}
}
