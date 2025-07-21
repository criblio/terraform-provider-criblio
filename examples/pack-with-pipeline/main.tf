resource "criblio_pack_pipeline" "my_packpipeline" {
  conf = {
    async_func_timeout = 9066
    description        = "my_description"
    functions = [
      {
        conf = {
          description = "testing config"
        }
        description = "my_description"
        disabled    = false
        filter      = "my_filter"
        final       = true
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
  id           = "pack-with-pipeline"
  group_id     = "default"
  description  = "Pack from source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"

}

