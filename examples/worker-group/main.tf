resource "criblio_group" "azure_eastus_stream_group" {
  cloud = {
    provider = "azure"
    region   = "eastus"
  }
  estimated_ingest_rate = 1024
  id                    = "my-group"
  name                  = "my-group"
  is_fleet              = false
  on_prem               = false
  product               = "stream"
  provisioned           = false
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

