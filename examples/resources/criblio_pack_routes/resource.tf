resource "criblio_pack_routes" "my_packroutes" {
  comments = [
    {
      additional_properties = "{ \"see\": \"documentation\" }"
      comment               = "...my_comment..."
    }
  ]
  group_id = "Cribl"
  groups = {
    key = {
      description = "...my_description..."
      disabled    = true
      name        = "...my_name..."
    }
  }
  id   = "default"
  pack = "observability-pack"
  routes = [
    {
      additional_properties    = "{ \"see\": \"documentation\" }"
      description              = "Route application errors to Splunk output"
      disabled                 = false
      enable_output_expression = false
      filter                   = "level == 'error'"
      final                    = true
      name                     = "Errors to Splunk"
      output = {
        key = jsonencode("value")
      }
      output_expression = {
        key = jsonencode("value")
      }
      pipeline = "main"
    }
  ]
}