
resource "criblio_group" "my_edge_subfleet" {
  //count required for cribl internal testing
  //count is not required for most customer implementations
  count = var.onprem == false ? 1 : 0

  estimated_ingest_rate = 1024
  id                    = "my-edge-subfleet"
  is_fleet              = true
  name                  = "my-edge-subfleet"
  on_prem               = false
  product               = "edge"
  provisioned           = false
  inherits              = "default_fleet"
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

output "edge_subfleet" {
  value = criblio_group.my_edge_subfleet
}
