resource "criblio_pipeline" "my_pipeline" {
  conf = {
    async_func_timeout = 5000
    description        = "Main pipeline for app logs"
    functions = [
      {
        conf = {
          add = [
            {
              disabled = true
              name     = "...my_name..."
              value    = "...my_value..."
            }
          ]
          remove = [
            "..."
          ]
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
        disabled    = false
        name        = "...my_name..."
      }
    }
    output = "OutputSplunk"
    streamtags = [
      "prod",
      "app",
    ]
  }
  group_id = "Cribl"
  id       = "main"
}