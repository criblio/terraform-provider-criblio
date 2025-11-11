terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "0.1.0"
    }
  }
}

provider "criblio" {
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id    = "main"
  cloud_domain    = "cribl-playground.cloud"
}

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
