resource "criblio_group" "workergroup" {
  cloud = {
    provider = "aws"
    region   = var.cloud_region
  }
  estimated_ingest_rate = var.estimated_ingest_rate
  id                    = var.group-cloud
  is_fleet              = false
  name                  = var.group-cloud
  on_prem               = false
  product               = "stream"
  provisioned           = true
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

output "group" {
  value = criblio_group.workergroup
}

data "criblio_config_version" "my_configversion" {
  id = var.group-cloud
  depends_on = [criblio_commit.my_commit]
}

resource "criblio_commit" "my_commit" {
  effective = true
  group   = var.group-cloud
  message = "terraform commit cloud commit"
}

resource "criblio_deploy" "my_deploy" {
  id      = var.group-cloud
  version = length(data.criblio_config_version.my_configversion.items) > 0 ? data.criblio_config_version.my_configversion.items[length(data.criblio_config_version.my_configversion.items) - 1] : "default"
}