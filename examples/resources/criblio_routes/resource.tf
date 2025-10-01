resource "criblio_routes" "my_routes" {
  comments = [
    {
      additional_properties = "{ \"see\": \"documentation\" }"
      comment               = "...my_comment..."
    }
  ]
  group_id = "myExistingGroupId"
  groups = {
    key = {
      description = "...my_description..."
      disabled    = false
      name        = "...my_name..."
    }
  }
  id = "...my_id..."
  routes = [
    {
      additional_properties    = "{ \"see\": \"documentation\" }"
      description              = "...my_description..."
      disabled                 = false
      enable_output_expression = false
      filter                   = "...my_filter..."
      final                    = false
      name                     = "...my_name..."
      output                   = "{ \"see\": \"documentation\" }"
      output_expression        = "{ \"see\": \"documentation\" }"
      pipeline                 = "...my_pipeline..."
    }
  ]
}