# Hybrid Group Configuration
resource "criblio_group" "my_group_defaulthybrid" {
  estimated_ingest_rate = var.estimated_ingest_rate
  id                    = var.group-hybrid
  is_fleet              = false
  name                  = var.group-hybrid
  on_prem               = true
  product               = "stream"
  provisioned           = false
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}
# Bootstrap Token Module - gets auth token after hybrid group is built
module "bootstrap_token" {
  source = "../cribl-bootstrap-token-module"
  
  client_id     = var.cribl_client_id
  client_secret = var.cribl_client_secret
  organization  = var.cloud_tenant
  workspace     = var.workspace
  group         = var.group-hybrid
  depends_on = [criblio_group.my_group_defaulthybrid]
}

# Use the worker installation module
module "cribl_worker" {
  source = "../cribl-worker-module"
  
  organization = var.cloud_tenant
  workspace    = var.workspace
  group        = var.group-hybrid
  auth_token   = module.bootstrap_token.bootstrap_token
  depends_on = [criblio_group.my_group_defaulthybrid]
}




data "criblio_config_version" "my_pipelineconfigversion_hybrid" {
  id = var.group-hybrid
  depends_on = [criblio_commit.my_pipecommit]
}

resource "criblio_commit" "my_pipecommit_hybrid" {
  effective = true
  group     = var.group-hybrid
  message   = "terraform commit hybrid"
}
resource "criblio_deploy" "my_pipedeploy_safe_hybrid" {
  id = var.group-hybrid
  version = length(data.criblio_config_version.my_pipelineconfigversion_hybrid.items) > 0 ? data.criblio_config_version.my_pipelineconfigversion_hybrid.items[length(data.criblio_config_version.my_pipelineconfigversion_hybrid.items) - 1] : "default"
}

# Read the auth token from the file (for backward compatibility)
data "local_file" "auth_token" {
  filename = module.bootstrap_token.token_file_path
  depends_on = [module.bootstrap_token]
}

