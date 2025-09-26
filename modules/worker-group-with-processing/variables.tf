# Worker Group Variables
variable "group_id" {
  description = "Unique identifier for the worker group"
  type        = string
}

variable "group_name" {
  description = "Display name for the worker group (must match pattern: lowercase letters, numbers, and hyphens only)"
  type        = string
  validation {
    condition     = can(regex("^[a-z0-9-]+$", var.group_name))
    error_message = "Group name must contain only lowercase letters, numbers, and hyphens."
  }
}

variable "cloud_provider" {
  description = "Cloud provider (aws or azure)"
  type        = string
  validation {
    condition     = contains(["aws", "azure"], var.cloud_provider)
    error_message = "Cloud provider must be either 'aws' or 'azure'."
  }
}

variable "cloud_region" {
  description = "Cloud region where the worker group will be deployed"
  type        = string
}

variable "estimated_ingest_rate" {
  description = "Estimated data ingest rate in MB/s"
  type        = number
  default     = 1024
  validation {
    condition     = contains([1024, 2048, 3072, 4096, 5120, 7168, 10240, 13312, 15360], var.estimated_ingest_rate)
    error_message = "Estimated ingest rate must be one of: 1024, 2048, 3072, 4096, 5120, 7168, 10240, 13312, 15360 MB/s."
  }
}

variable "description" {
  description = "Description for the worker group"
  type        = string
  default     = ""
}

variable "is_fleet" {
  description = "Whether this is a fleet worker group"
  type        = bool
  default     = false
}

variable "on_prem" {
  description = "Whether this is an on-premises worker group"
  type        = bool
  default     = false
}

variable "product" {
  description = "Cribl product (stream or edge)"
  type        = string
  default     = "stream"
  validation {
    condition     = contains(["stream", "edge"], var.product)
    error_message = "Product must be either 'stream' or 'edge'."
  }
}

variable "provisioned" {
  description = "Whether the worker group is provisioned"
  type        = bool
  default     = true
}

variable "worker_remote_access" {
  description = "Whether to enable remote access for workers"
  type        = bool
  default     = false
}

variable "streamtags" {
  description = "Stream tags for the worker group"
  type        = list(string)
  default     = []
}

# Palo Alto Pack Variables
variable "install_palo_alto_pack" {
  description = "Whether to install the Palo Alto pack on this worker group"
  type        = bool
  default     = true
}

variable "palo_alto_pack_filename" {
  description = "Filename of the Palo Alto pack to install"
  type        = string
  default     = ""
}

variable "palo_alto_pack_description" {
  description = "Description for the Palo Alto pack"
  type        = string
  default     = ""
}

variable "palo_alto_pack_display_name" {
  description = "Display name for the Palo Alto pack"
  type        = string
  default     = ""
}

variable "palo_alto_pack_version" {
  description = "Version of the Palo Alto pack"
  type        = string
  default     = "1.0.0"
}

variable "palo_alto_pack_disabled" {
  description = "Whether the Palo Alto pack should be disabled"
  type        = bool
  default     = false
}

# CrowdStrike Pack Variables
variable "install_crowdstrike_pack" {
  description = "Whether to install the CrowdStrike pack on this worker group"
  type        = bool
  default     = true
}

variable "crowdstrike_pack_description" {
  description = "Description for the CrowdStrike pack"
  type        = string
  default     = "CrowdStrike Security Pack with Event Breaker Rulesets"
}

variable "crowdstrike_pack_display_name" {
  description = "Display name for the CrowdStrike pack"
  type        = string
  default     = "CrowdStrike Security Pack"
}

variable "crowdstrike_pack_version" {
  description = "Version of the CrowdStrike pack"
  type        = string
  default     = "1.0.0"
}

variable "crowdstrike_pack_disabled" {
  description = "Whether the CrowdStrike pack should be disabled"
  type        = bool
  default     = false
}

# Pipeline Variables
variable "create_pipeline" {
  description = "Whether to create a processing pipeline"
  type        = bool
  default     = true
}

variable "pipeline_description" {
  description = "Description for the processing pipeline"
  type        = string
  default     = "Standard data processing pipeline"
}

variable "pipeline_timeout" {
  description = "Async function timeout for the pipeline"
  type        = number
  default     = 60
}

variable "pipeline_output" {
  description = "Output destination for the pipeline"
  type        = string
  default     = "default"
}

variable "pipeline_streamtags" {
  description = "Stream tags for the pipeline"
  type        = list(string)
  default     = ["processing"]
}

variable "region_identifier" {
  description = "Region identifier to add to processed data"
  type        = string
}

# Commit and Deploy Variables
variable "auto_commit" {
  description = "Whether to automatically commit configuration changes"
  type        = bool
  default     = true
}

variable "auto_deploy" {
  description = "Whether to automatically deploy committed configurations"
  type        = bool
  default     = true
}

variable "commit_message" {
  description = "Commit message for configuration changes"
  type        = string
  default     = "Deploy pack and processing pipeline"
}
