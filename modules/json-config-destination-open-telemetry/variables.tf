variable "group_id" {
  description = "The group ID to create the destination in"
  type        = string
}

variable "json_file_path" {
  description = "Path to JSON file containing the destination configuration"
  type        = string
  default     = ""
}

variable "json_content" {
  description = "JSON content for the destination configuration (alternative to json_file_path)"
  type        = string
  default     = ""
}

variable "destination_id_override" {
  description = "Override the destination ID from the JSON configuration"
  type        = string
  default     = ""
}

variable "additional_config" {
  description = "Additional configuration to merge with the parsed JSON"
  type        = any
  default     = {}
}
