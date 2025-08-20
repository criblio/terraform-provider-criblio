resource "criblio_pack" "my_pack" {
  id           = "pack-with-dest"
  group_id     = "default"
  description  = "Pack from source with destination"
  disabled     = true
  display_name = "Pack from source with destination"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_destination" "my_packdest" {
  pack     = criblio_pack.my_pack.id
  group_id = "default"
  id       = "my_id"
  output_devnull = {
    environment = "my_environment"
    id          = "my_id"
    pipeline    = "my_pipeline"
    streamtags = [
      ""
    ]
    system_fields = [
      ""
    ]
    type = "devnull"
  }
}
