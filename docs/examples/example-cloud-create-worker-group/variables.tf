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
# - worker_group_id: ID for the new Worker Group (must not already exist)
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
  description = "ID for the new Worker Group. Must not already exist in your deployment."
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
