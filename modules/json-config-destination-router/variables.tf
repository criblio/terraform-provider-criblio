variable "group_id" {
  description = "The ID of the group where the destination will be created"
  type        = string
}

variable "json_file_path" {
  description = "Path to the JSON file containing the destination configuration. Use either this or json_content."
  type        = string
  default     = ""
}

variable "json_content" {
  description = "JSON content as a string containing the destination configuration. Use either this or json_file_path."
  type        = string
  default     = ""
}

variable "destination_id_override" {
  description = "Override the destination ID from the JSON configuration"
  type        = string
  default     = ""
}

variable "additional_config" {
  description = "Additional configuration to merge with the JSON configuration"
  type        = any
  default     = {}
}
