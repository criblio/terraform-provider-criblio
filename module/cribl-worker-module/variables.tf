# Variables for Cribl Worker Installation Module

variable "organization" {
  description = "Cribl Cloud organization ID"
  type        = string
}

variable "workspace" {
  description = "Cribl workspace name"
  type        = string
}

variable "workspace_url" {
  description = "Full workspace URL (optional - if provided, overrides organization/workspace)"
  type        = string
  default     = ""
}

variable "group" {
  description = "Worker group name"
  type        = string
}

variable "auth_token" {
  description = "Authentication token for worker installation"
  type        = string
  sensitive   = true
}

variable "domain" {
  description = "Cribl domain"
  type        = string
  default     = "cribl.cloud"
}

variable "cribl_user" {
  description = "User for Cribl installation"
  type        = string
  default     = "cribl"
}

variable "cribl_group" {
  description = "Group for Cribl installation"
  type        = string
  default     = "cribl"
}

variable "install_dir" {
  description = "Installation directory (URL encoded)"
  type        = string
  default     = "%2Fopt%2Fcribl"
} 