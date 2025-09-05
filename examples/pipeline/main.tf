resource "criblio_pipeline" "my_pipeline" {
  conf = {
    async_func_timeout = 60
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
  group_id = "default"
  id       = "pipeline-1"
}
