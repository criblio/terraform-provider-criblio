# Variables for Cribl Bootstrap Token Module

variable "client_id" {
  description = "Cribl Cloud OAuth2 client ID"
  type        = string
  sensitive   = true
}

variable "client_secret" {
  description = "Cribl Cloud OAuth2 client secret"
  type        = string
  sensitive   = true
}

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

variable "environment" {
  description = "Environment: production or staging"
  type        = string
  default     = "production"
  
  validation {
    condition     = contains(["production", "staging"], var.environment)
    error_message = "Environment must be either 'production' or 'staging'."
  }
}

variable "group" {
  description = "Group parameter for bootstrap token"
  type        = string
  default     = "defaultHybrid"
}

variable "auto_execute" {
  description = "Whether to automatically execute the script to get the token"
  type        = bool
  default     = true
}

variable "create_script_file" {
  description = "Whether to create a physical script file (useful for debugging)"
  type        = bool
  default     = false
}

variable "script_filename" {
  description = "Name of the script file to create (if create_script_file is true)"
  type        = string
  default     = "get_bootstrap_token.sh"
}

variable "token_output_file" {
  description = "File to store the retrieved bootstrap token"
  type        = string
  default     = "auth_token.txt"
}

# Note: Terraform's depends_on requires static lists, so dependencies
# should be handled at the module call level, not within the module 