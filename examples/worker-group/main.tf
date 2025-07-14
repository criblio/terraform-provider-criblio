resource "criblio_group" "my_group" {
  #
  cloud = {
      # options: aws, azure
      # region: us-west-2, eastus, etc.
    provider = "aws"
    region   = "us-east-2"
  }
  estimated_ingest_rate = 1024
  id                    = "my-group"
  is_fleet              = false
  name                  = "my-group"
  on_prem               = false
  product               = "stream"
  provisioned           = false
  streamtags = [
    "test",
    "network"
  ]
  worker_remote_access = false
}

output "group" {
  value = criblio_group.my_group
}
