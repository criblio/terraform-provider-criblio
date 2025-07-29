# Variables for Basic Example
# These can be populated from the shared provider configuration

variable "cribl_client_id" {
  description = "Cribl Cloud OAuth2 client ID"
  type        = string
  sensitive   = true
}

variable "cribl_client_secret" {
  description = "Cribl Cloud OAuth2 client secret"
  type        = string
  sensitive   = true
}

variable "cribl_organization" {
  description = "Cribl Cloud organization ID"
  type        = string
}

variable "cribl_workspace" {
  description = "Cribl workspace name"
  type        = string
}

variable "cribl_group" {
  description = "Cribl worker group name"
  type        = string
  default     = "defaultHybrid"
} 