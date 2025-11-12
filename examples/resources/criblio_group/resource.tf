resource "criblio_group" "my_group" {
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
  description           = "Production Edge configuration group"
  estimated_ingest_rate = 500000
  id                    = "config-group-prod-edge"
  is_fleet              = false
  max_worker_age        = "1h"
  name                  = "Prod Edge"
  on_prem               = true
  product               = "stream"
  provisioned           = true
  streamtags = [
    "prod",
    "edge",
  ]
  tags                 = "environment=prod,team=platform"
  type                 = "lake_access"
  worker_remote_access = true
}