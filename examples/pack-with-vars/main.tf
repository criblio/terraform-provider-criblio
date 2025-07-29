resource "criblio_pack_vars" "my_packvars" {
  description = "my_description"
  group_id    = "default"
  id          = "my_id"
  lib         = "my_lib"
  pack        = criblio_pack.vars_pack.id
  tags        = "my_tags"
  type        = "object"
  value       = "my_value"
}

resource "criblio_pack" "vars_pack" {
  id           = "pack-with-vars"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}
