# User-supplied parameters
#
# Supply values for the declared variables in this file in terraform.tfvars.
#
# Required for the Worker Group:
# - worker_group_id: ID of the existing Worker Group to scale
#
# Before first apply, import existing system settings into Terraform state:
#   terraform import criblio_group_system_settings.scaled <worker_group_id>
#
# Replace the placeholder values in the api, ssl, and other blocks in main.tf
# with values that match your deployment.

variable "worker_group_id" {
  type        = string
  description = "ID of the existing Worker Group to scale."
}
