resource "criblio_pack_breakers" "my_packbreakers" {
  description         = "my_description"
  disabled            = true
  display_name        = "my_display_name"
  group_id            = "my_group_id"
  id                  = "my_id"
  pack                = criblio_pack.breakers_pack.id
  pack_path_parameter = "my_pack_path_parameter"
  source              = "my_source"
  version             = "my_version"
}

resource "criblio_pack" "breakers_pack" {
  id           = "pack-with-breakers"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

# Output the pack details to see the read-only attributes
output "pack_breakers_details" {
  value = criblio_pack.my_pack
}
