resource "criblio_pack_pipeline" "my_packpipeline" {
  conf = {
    async_func_timeout = 9066
    description        = "...my_description..."
    functions = [
      {
        conf = {
          key = jsonencode("value")
        }
        description = "...my_description..."
        disabled    = false
        filter      = "...my_filter..."
        final       = true
        group_id    = "...my_group_id..."
        id          = "...my_id..."
      }
    ]
    groups = {
      key = {
        description = "...my_description..."
        disabled    = true
        name        = "...my_name..."
      }
    }
    output = "...my_output..."
    streamtags = [
      "..."
    ]
  }
  group_id = "...my_group_id..."
  id       = "...my_id..."
  pack     = "...my_pack..."
}