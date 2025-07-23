resource "criblio_pack_pipeline" "my_packpipeline" {
  group_id = "default"
  id       = "my_id"
  pack     = criblio_pack.my_pack.id
  conf = {
    async_func_timeout = 9066
    description        = "my_description"
    functions = [
    ]
    groups = {
      key = {
        description = "my_description"
        disabled    = true
        name        = "my_name"
      }
    }
    output = "my_output"
    streamtags = [
      "tags"
    ]
  }
}

resource "criblio_pack" "my_pack" {
  id           = "pack-with-pipeline"
  group_id     = "default"
  description  = "Pack with pipeline"
  disabled     = true
  display_name = "Pack with pipeline"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

