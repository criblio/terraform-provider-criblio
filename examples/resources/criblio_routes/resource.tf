resource "criblio_routes" "my_routes" {
  comments = [
    {
      additional_properties = "{ \"see\": \"documentation\" }"
      comment               = "...my_comment..."
    }
  ]
  group_id = "default"
  groups = {
    key = {
      description = "...my_description..."
      disabled    = false
      name        = "...my_name..."
    }
  }
  id = "default"
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