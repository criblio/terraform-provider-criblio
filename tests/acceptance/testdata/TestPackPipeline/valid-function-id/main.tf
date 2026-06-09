resource "criblio_pack_pipeline" "valid_function_id" {
  id       = "valid-function-id-test"
  pack     = "test-pack"
  group_id = "default"
  conf = {
    functions = [
      {
        id     = "publish_metrics"
        filter = "true"
        conf   = jsonencode({})
      }
    ]
  }
}
