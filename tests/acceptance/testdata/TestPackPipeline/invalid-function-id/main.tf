resource "criblio_pack_pipeline" "invalid_function_id" {
  id       = "invalid-function-id-test"
  pack     = "test-pack"
  group_id = "default"
  conf = {
    functions = [
      {
        id     = "publish_metrics_mistaken_unique_name"
        filter = "true"
        conf   = jsonencode({})
      }
    ]
  }
}
