# Scale Worker Group
#
# Before first apply, import the existing Worker Group into Terraform state:
#   terraform import criblio_group.example <worker_group_id>
resource "criblio_group" "example" {
  id                    = var.worker_group_id
  name                  = var.worker_group_name
  description           = var.worker_group_description
  product               = "stream"
  on_prem               = false
  is_fleet              = false
  worker_remote_access  = true
  estimated_ingest_rate = 4096 # Equivalent to 48 MB/s with 21 Worker Processes
  provisioned           = true
  cloud = {
    provider = var.cloud_provider
    region   = var.cloud_region
  }
}

# Commit configuration
resource "criblio_commit" "example" {
  effective = true
  group     = criblio_group.example.id
  message   = "Commit for scale Worker Group example"

  depends_on = [criblio_group.example]
}

# Read config version
data "criblio_config_version" "latest" {
  id         = criblio_group.example.id
  depends_on = [criblio_commit.example]
}

# Deploy configuration
resource "criblio_deploy" "example" {
  id      = criblio_group.example.id
  version = data.criblio_config_version.latest.items[0]
}
