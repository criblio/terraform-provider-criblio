resource "criblio_pack_pipeline" "my_packpipeline" {
  conf = {
    async_func_timeout = 5000
    description        = "Main pipeline for app logs"
    functions = [
      {
        conf = {
          key = jsonencode("value")
        }
        description = "Parse and enrich fields"
        disabled    = false
        filter      = "_source == \"app\""
        final       = false
        group_id    = "default"
        id          = "eval"
      }
    ]
    groups = {
      key = {
        description = "...my_description..."
        disabled    = true
        name        = "...my_name..."
      }
    }
    output = "OutputSplunk"
    streamtags = [
      "prod",
      "app",
    ]
  }
  group_id = "myExistingGroupId"
  id       = "main"
  pack     = "myExistingPackId"
}