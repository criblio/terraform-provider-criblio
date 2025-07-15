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

//import currently required for our testing due to provider not currently allowing deletion of criblio_group resource
//will be removed in a future revision when we enable resource deletion -- 15Jul2025 KT
import {
  to = criblio_group.my_edge_fleet
  id = "my-edge-fleet"
}
