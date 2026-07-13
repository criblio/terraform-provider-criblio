# User-supplied parameters
#
# Supply values for the declared variables in this file in terraform.tfvars.
#
# Required for Cribl.Cloud authentication:
# - organization_id: Organization ID
# - workspace_id: Workspace name
# - client_id: Client ID for your API Credential
# - client_secret: Client Secret for your API Credential
#
# Required for the Worker Group:
# - worker_group_id: ID of the existing Worker Group to scale
# - worker_group_name: Name of the existing Worker Group
# - worker_group_description: Description of the existing Worker Group
# - cloud_provider: Cloud provider for the Worker Group (aws or azure)
# - cloud_region: Cloud region for the Worker Group
#
# To create an API Credential, see:
# https://docs.cribl.io/cribl-as-code/terraform-auth/#terraform-auth-cloud

variable "organization_id" {
  type        = string
  description = "Cribl.Cloud Organization ID"
}

variable "workspace_id" {
  type        = string
  description = "Cribl.Cloud Workspace name"
}

variable "worker_group_id" {
  type        = string
  description = "ID of the existing Worker Group to scale."
}

variable "worker_group_name" {
  type        = string
  description = "Name of the existing Worker Group."
}

variable "worker_group_description" {
  type        = string
  description = "Description of the existing Worker Group."
}

variable "cloud_provider" {
  type        = string
  description = "Cloud provider for the Worker Group (aws or azure)."
}

variable "cloud_region" {
  type        = string
  description = "Cloud region for the Worker Group."
}

variable "client_id" {
  type        = string
  description = "Client ID from your Cribl API Credential"
  sensitive   = true
}

variable "client_secret" {
  type        = string
  description = "Client Secret from your Cribl API Credential"
  sensitive   = true
}
