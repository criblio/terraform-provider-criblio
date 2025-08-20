variable "cloud_tenant" {
  description = "Cribl Cloud organization ID"
  type        = string
  default     = "zamora-berries-wild8"
}

variable "workspace" {
  description = "Cribl workspace name"
  type        = string
  default     = "main"
}

variable "group-hybrid" {
  description = "Worker group name"
  type        = string
  default     = "notdefaulthybrid"
}
variable "group-cloud" {
  description = "Stream group name"
  type        = string
  default     = "notdefault"
}

variable "cribl_version" {
  description = "Cribl version to use"
  type        = string
  default     = "4.12.0"
}

variable "environment" {
  description = "Environment (production, staging)"
  type        = string
  default     = "production"

  validation {
    condition     = contains(["production", "staging"], var.environment)
    error_message = "Environment must be either 'production' or 'staging'."
  }
}

variable "estimated_ingest_rate" {
  description = "Estimated data ingest rate"
  type        = number
  default     = 1024
}

variable "instance_count" {
  description = "Number of EC2 instances to create for the worker group"
  type        = number
  default     = 2
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.medium"
}

variable "cloud_region" {
  description = "AWS region for resources"
  type        = string
  default     = "ca-central-1"
}

# Sensitive variables (no defaults for security)
variable "cribl_client_id" {
  description = "Cribl OAuth2 client ID"
  type        = string
  sensitive   = true
  default     = "byebyebye"
  # No default - must be provided
}

variable "cribl_client_secret" {
  description = "Cribl OAuth2 client secret"
  type        = string
  sensitive   = true
  default     = "youcandancehavingthetimeofyourlife"
  # No default - must be provided
}

