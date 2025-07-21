resource "criblio_pack_source" "my_packsource" {
  conf = {
    async_func_timeout = 6187
    description        = "my_description"
    functions = [
      {
        conf = {
          # ...
        }
        description = "my_description"
        disabled    = false
        filter      = "my_filter"
        final       = false
        group_id    = "my_group_id"
        id          = "my_id"
      }
    ]
    groups = {
      key = {
        description = "my_description"
        disabled    = true
        name        = "my_name"
      }
    }
    output = "my_output"
  }
  group_id = "my_group_id"
  id       = "my_id"
  pack     = criblio_pack.my_pack.id
}

resource "criblio_pack" "my_pack" {
  id           = "pack-with-source"
  group_id     = "default"
  description  = "Pack with source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"

}

