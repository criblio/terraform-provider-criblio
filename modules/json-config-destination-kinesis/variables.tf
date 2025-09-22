variable "group_id" {
  description = "The group ID where the destination will be created"
  type        = string
}

variable "json_file_path" {
  description = "Path to the JSON file containing the destination configuration"
  type        = string
  default     = ""
}

variable "json_content" {
  description = "JSON string containing the destination configuration"
  type        = string
  default     = ""
}

variable "destination_id_override" {
  description = "Override for the destination ID (optional)"
  type        = string
  default     = ""
}

variable "additional_config" {
  description = "Additional configuration to merge with JSON config"
  type        = any
  default     = {}
}

# Validation is handled in locals.tf - both json_file_path and json_content cannot be empty
