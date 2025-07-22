resource "criblio_pack_routes" "my_packroutes" {
  comments = [
    {
      additional_properties = "{ \"see\": \"documentation\" }"
      comment               = "my_comment"
    }
  ]
  group_id = "default"
  groups = {
    key = {
      description = "my_description"
      disabled    = true
      name        = "my_name"
    }
  }
  id   = "my_id"
  pack = criblio_pack.my_pack.id
  routes = [
    {
      additional_properties    = "{ \"see\": \"documentation\" }"
      description              = "my_description"
      disabled                 = true
      enable_output_expression = true
      filter                   = "my_filter"
      final                    = false
      name                     = "my_name"
      output                   = "{ \"see\": \"documentation\" }"
      output_expression        = "{ \"see\": \"documentation\" }"
      pipeline                 = "my_pipeline"
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
