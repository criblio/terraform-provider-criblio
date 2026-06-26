resource "criblio_group" "my_group" {
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
  description           = "Production stream worker group"
  estimated_ingest_rate = 1024
  id                    = "config-group-prod-stream"
  is_fleet              = false
  name                  = "config-group-prod-stream"
  on_prem               = false
  product               = "stream"
  provisioned           = false
  streamtags = [
    "prod",
    "stream",
  ]
  tags                 = "environment=prod,team=platform"
  type                 = "stream"
  worker_remote_access = false
}
