

resource "criblio_group" "my_edge_fleet" {
  estimated_ingest_rate = 1024
  id                    = "my-edge-fleet"
  is_fleet              = true
  name                  = "my-edge-fleet"
  on_prem               = false
  product               = "edge"
  provisioned           = false
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

output "edge_fleet" {
  value = criblio_group.my_edge_fleet
}
