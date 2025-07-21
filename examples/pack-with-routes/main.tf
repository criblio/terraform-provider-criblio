resource "criblio_pack_routes" "my_packroutes" {
  description         = "my_description"
  disabled            = true
  display_name        = "my_display_name"
  group_id            = "default"
  id                  = "my_id"
  pack                = criblio_pack.my_pack.id
  pack_path_parameter = "my_pack_path_parameter"
  source              = "default"
  version             = "my_version"
}

resource "criblio_pack" "my_pack" {
  id           = "pack-with-routes"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}
