resource "criblio_project_pipeline" "example" {
  group_id   = "default"
  project_id = "example"
  id         = "main"

  conf = {
    async_func_timeout = 1000
    description        = "Main project pipeline"
    output             = "default"
    streamtags         = []
  }
}
