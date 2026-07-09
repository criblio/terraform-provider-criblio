# User-supplied parameters
#
# Supply values for the declared variables in this file in terraform.tfvars.
#
# Required for the Worker Group:
# - worker_group_id: ID for the new Worker Group (must not already exist)

variable "worker_group_id" {
  type        = string
  description = "ID for the new Worker Group. Must not already exist in your deployment."
}
