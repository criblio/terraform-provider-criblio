
locals {
  streamtags_hybrid = [
    "datacenter1",
    "someothertag"
  ]
}
resource "criblio_group" "hybrid_worker_group" {
  id                   = "my-hybrid-group"
  name                 = "my-hybrid-group"
  is_fleet             = false
  on_prem              = true
  product              = "stream"
  streamtags           = local.streamtags_hybrid
  worker_remote_access = false
}

module "hybrid_worker_group_bootstrap" {
  source          = "../../modules/criblio-hybrid-worker-bootstrap"
  group_id        = criblio_group.hybrid_worker_group.id
  group_tags      = local.streamtags_hybrid
  organization_id = var.organization_id
  workspace_id    = var.workspace_id
  cloud_domain    = var.cloud_domain
}
