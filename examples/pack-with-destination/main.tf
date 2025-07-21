resource "criblio_pack_destination" "my_packdestination" {
  allow_custom_functions = false
  author                 = "...my_author..."
  description            = "...my_description..."
  display_name           = "...my_display_name..."
  exports = [
    "..."
  ]
  force                  = true
  group_id               = "...my_group_id..."
  id                     = "...my_id..."
  inputs                 = 1.12
  is_disabled            = true
  min_log_stream_version = "...my_min_log_stream_version..."
  outputs                = 3.79
  pack                   = criblio_pack.my_pack.id
  source                 = "...my_source..."
  spec                   = "...my_spec..."
  tags = {
    data_type = [
      "..."
    ]
    domain = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    technology = [
      "..."
    ]
  }
  version = "...my_version..."
}

resource "criblio_pack" "my_pack" {
  id           = "pack-from-source"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"

}

# Output the pack details to see the read-only attributes
output "pack_details" {
  value = criblio_pack.my_pack
}
