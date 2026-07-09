# Create Worker Group
resource "criblio_group" "example" {
  id                    = var.worker_group_id
  name                  = "my-worker-group"
  description           = "Cribl.Cloud Worker Group"
  product               = "stream"
  on_prem               = false
  is_fleet              = false
  worker_remote_access  = true
  estimated_ingest_rate = 2048 # Equivalent to 24 MB/s with 9 Worker Processes
  provisioned           = false
  cloud = {
    provider = "aws"
    region   = "us-west-2"
  }
}

# Commit configuration
resource "criblio_commit" "example" {
  effective = true
  group     = criblio_group.example.id
  message   = "Commit for create Worker Group example"

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
