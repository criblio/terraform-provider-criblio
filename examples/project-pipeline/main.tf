resource "criblio_project" "my_project" {
  group_id      = "default"
  id            = "project-1"
  description   = "Project pipeline example"
  destinations  = []
  subscriptions = []
}

resource "criblio_project_pipeline" "my_pipeline" {
  id         = "pipeline-1"
  group_id   = "default"
  project_id = criblio_project.my_project.id
  conf = {
    streamtags         = []
    description        = "my_description"
    output             = "my_output"
    async_func_timeout = 60
    functions = [
      {
        id       = "serde"
        filter   = "true"
        disabled = false
        conf = jsonencode({
          srcField = "_raw"
          mode     = "extract"
          type     = "json"
        })
      },
      {
        id       = "eval"
        filter   = "channel===\"ProcessMetrics\""
        disabled = false
        final    = true
        conf = jsonencode({
          remove = [
            "_raw"
          ]
        })
      },
      {
        id       = "drop"
        filter   = "!(/log\\/(?:cribl|access|audit)\\.log$/.test(source) || /service\\/(?:metrics|connection_proxy|lease|notifications|\\w*connections)/.test(source))"
        disabled = false
        conf     = jsonencode({})
      },
      {
        id       = "eval"
        filter   = "true"
        disabled = false
        conf = jsonencode({
          remove = [
            "_raw"
          ]
          add = [
            {
              name     = "url"
              value    = "url.match(\"(?:.*system\\/metrics)\")[0] || url"
              disabled = false
            }
          ]
        })
      },
    ]
  }
}

data "criblio_project_pipeline" "my_pipeline" {
  group_id   = criblio_project_pipeline.my_pipeline.group_id
  project_id = criblio_project_pipeline.my_pipeline.project_id
  id         = criblio_project_pipeline.my_pipeline.id
}

output "pipeline_conf" {
  value = data.criblio_project_pipeline.my_pipeline.conf
}
