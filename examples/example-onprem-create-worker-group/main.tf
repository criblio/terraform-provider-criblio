# Create Worker Group
resource "criblio_group" "example" {
  id                   = var.worker_group_id
  name                 = "my-worker-group"
  description          = "My Worker Group description"
  product              = "stream"
  on_prem              = true
  is_fleet             = false
  worker_remote_access = true
  provisioned          = false
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
