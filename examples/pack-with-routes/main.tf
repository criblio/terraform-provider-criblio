resource "criblio_pack_routes" "my_packroutes" {
  group_id = "default"
  pack     = criblio_pack.my_pack.id
  routes = [
    {
      name     = "my_name"
      pipeline = "main"
    }
  ]
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
