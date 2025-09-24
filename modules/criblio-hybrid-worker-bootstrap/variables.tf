# Variables for Cribl Worker Installation Module

variable "organization_id" {
  description = "Cribl Cloud organization ID"
  type        = string
}

variable "workspace_id" {
  description = "Cribl workspace name"
  type        = string
}

variable "group_id" {
  description = "Worker group name"
  type        = string
}
variable "group_tags" {
  description = "Worker group tags"
  type        = list(string)
  default     = []
}

variable "cloud_domain" {
  description = "Cribl domain"
  type        = string
  default     = "cribl.cloud"
}

variable "os_cribl_user" {
  description = "Operating system username for the Cribl installation."
  type        = string
  default     = "cribl"
}

variable "os_cribl_group" {
  description = "Operating system group name for the Cribl installation."
  type        = string
  default     = "cribl"
}

variable "install_dir" {
  description = "Installation path"
  type        = string
  default     = "/opt/cribl"
} 