variable "group_id" {
  description = "Worker group ID where the source will be created"
  type        = string
}

variable "json_file_path" {
  description = "Path to JSON configuration file for the source"
  type        = string
  default     = ""
}

variable "json_content" {
  description = "Direct JSON content as string (alternative to json_file_path)"
  type        = string
  default     = ""
}

variable "source_id_override" {
  description = "Override the source ID from JSON config"
  type        = string
  default     = ""
}

variable "additional_config" {
  description = "Additional configuration to merge with JSON config"
  type        = any
  default     = {}
}

# Validation
locals {
  has_file    = var.json_file_path != ""
  has_content = var.json_content != ""

  validation_check       = local.has_file || local.has_content ? true : tobool("Either json_file_path or json_content must be provided")
  mutual_exclusive_check = local.has_file && local.has_content ? tobool("Only one of json_file_path or json_content should be provided") : true
}
