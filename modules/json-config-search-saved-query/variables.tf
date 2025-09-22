variable "json_file_path" {
  description = "Path to JSON file containing the saved query configuration"
  type        = string
  default     = ""
}

variable "json_content" {
  description = "JSON content for the saved query configuration (alternative to json_file_path)"
  type        = string
  default     = ""
}

variable "query_id_override" {
  description = "Override the query ID from the JSON configuration"
  type        = string
  default     = ""
}

variable "additional_config" {
  description = "Additional configuration to merge with the parsed JSON"
  type        = any
  default     = {}
}
