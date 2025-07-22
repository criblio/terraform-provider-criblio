resource "criblio_pack_routes" "my_packroutes" {
  comments = [
    {
      additional_properties = "{ \"see\": \"documentation\" }"
      comment               = "...my_comment..."
    }
  ]
  group_id = "...my_group_id..."
  groups = {
    key = {
      description = "...my_description..."
      disabled    = true
      name        = "...my_name..."
    }
  }
  id   = "...my_id..."
  pack = "...my_pack..."
  routes = [
    {
      additional_properties    = "{ \"see\": \"documentation\" }"
      description              = "...my_description..."
      disabled                 = true
      enable_output_expression = true
      filter                   = "...my_filter..."
      final                    = false
      name                     = "...my_name..."
      output                   = "{ \"see\": \"documentation\" }"
      output_expression        = "{ \"see\": \"documentation\" }"
      pipeline                 = "...my_pipeline..."
    }
  ]
}