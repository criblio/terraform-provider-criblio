
variable "client_id" {
  description = "The client ID for authentication."
  type        = string
}

variable "client_secret" {
  description = "The client secret for authentication."
  type        = string
  sensitive   = true
}

variable "organization_id" {
  description = "The organization ID."
  type        = string
}

variable "workspace_id" {
  description = "The workspace ID."
  type        = string
}

variable "cloud_domain" {
  description = "The cloud domain."
  type        = string
  default     = "cribl.cloud"
}
