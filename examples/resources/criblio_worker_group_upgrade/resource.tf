variable "worker_group_id" {
  type        = string
  description = "ID of an existing self-managed Worker Group, including a hybrid group, with running workers. Cribl-managed cloud groups cannot use this resource."
}

# Change this stable trigger deliberately to submit another upgrade.
resource "criblio_worker_group_upgrade" "example" {
  group_id           = var.worker_group_id
  trigger            = "upgrade-1"
  upgrade_mode       = "rolling"
  upgrade_percentage = 0.5
  worker_retries     = 5
  worker_retry_delay = 1000
}
