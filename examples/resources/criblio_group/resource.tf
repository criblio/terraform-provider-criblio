resource "criblio_group" "my_group" {
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
  config_version         = "v12"
  deploying_worker_count = 3
  description            = "Production Edge configuration group"
  estimated_ingest_rate  = 500000
  git = {
    commit        = "1a2b3c4d5e6f"
    local_changes = 0
    log = [
      {
        author_email = "dev@acme.com"
        author_name  = "Acme Dev"
        date         = "2025-10-06T12:00:00Z"
        hash         = "1a2b3c4d5e6f7g8h9i0j"
        message      = "Deploy config for prod edge group"
        short        = "1a2b3c4"
      }
    ]
  }
  id                        = "config-group-prod-edge"
  incompatible_worker_count = 0
  inherits                  = "base-config"
  is_fleet                  = false
  is_search                 = false
  lookup_deployments = [
    {
      context = "prod"
      lookups = [
        {
          deployed_version = "v3"
          file             = "lookups/ip_to_region.csv"
          version          = "v4"
        }
      ]
    }
  ]
  name        = "Prod Edge"
  on_prem     = true
  product     = "stream"
  provisioned = true
  streamtags = [
    "prod",
    "edge",
  ]
  tags                 = "environment=prod,team=platform"
  type                 = "lake_access"
  upgrade_version      = "6.0.1"
  worker_count         = 12
  worker_remote_access = true
}