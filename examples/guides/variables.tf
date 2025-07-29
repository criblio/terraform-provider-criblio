variable "cloud_tenant" {
  description = "Cribl Cloud organization ID"
  type        = string
  default     = ""
}

variable "workspace" {
  description = "Cribl workspace name"
  type        = string
  default     = "main"
}

variable "group-hybrid" {
  description = "Worker group name"
  type        = string
  default     = "my-group-tf-hybrid"
}
variable "group-cloud" {
  description = "Stream group name"
  type        = string
  default     = "my-group-tf"
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
  default     = 1
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = ""
}

variable "cloud_region" {
  description = "AWS region for resources"
  type        = string
  default     = ""
}

# Sensitive variables (no defaults for security)
variable "cribl_client_id" {
  description = "Cribl OAuth2 client ID"
  type        = string
  sensitive   = true
  default = ""
  # No default - must be provided
}

variable "cribl_client_secret" {
  description = "Cribl OAuth2 client secret"
  type        = string
  sensitive   = true
  default = ""
  # No default - must be provided
}

variable "aws_key_name" {
  description = "AWS key pair name for EC2 instances"
  type        = string
  default     = ""
}

variable "aws_profile" {
  description = "AWS profile to use"
  type        = string
  default     = ""
}

variable "ami_id" {
  description = "AMI ID for EC2 instances (optional - if not provided, will use latest Ubuntu)"
  type        = string
  default     = "" # Empty means use data source to find latest
}

 