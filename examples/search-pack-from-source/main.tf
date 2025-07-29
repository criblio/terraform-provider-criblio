resource "criblio_pack" "my_search_pack" {
  id           = "pack-from-source"
  group_id     = "default_search"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

