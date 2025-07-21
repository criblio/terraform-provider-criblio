resource "criblio_pack_source" "my_packsource" {
  conf = {
    async_func_timeout = 6187
    description        = "...my_description..."
    functions = [
      {
        conf = {
          # ...
        }
        description = "...my_description..."
        disabled    = false
        filter      = "...my_filter..."
        final       = false
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