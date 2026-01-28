terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id    = "main"
  cloud_domain    = "cribl-playground.cloud"
}

resource "criblio_pack_pipeline" "my_packpipeline" {
  group_id = "default"
  id       = "my_id"
  pack     = criblio_pack.pipeline_pack.id
  conf = {
    async_func_timeout = 9066
    description        = "my_description"
    functions = [
      {
        id     = "serde"
        filter = "true"
        conf = jsonencode({
          mode      = "extract"
          type      = "json"
          src_field = "_raw"
        })
      },
      {
        id     = "eval"
        filter = "true"
        conf = jsonencode({
          add = [
            {
              disabled = false
              name     = "_value"
              value    = "1"
            }
          ]
          remove = [
            "host",
            "_raw",
            "source",
            "cribl_breaker",
          ]
        })
      },
      {
        id     = "publish_metrics"
        filter = "true"
        conf = jsonencode({
          overwrite = false
          dimensions = [
            "!_*",
            "*",
          ],
          removeDimensions = ["cribl_pipe"]
          fields = [
            {
              metricType   = "gauge"
              inFieldName  = "_value"
              outFieldExpr = "saas_env"
            }
          ]
        })
      },
    ]
    groups = {
      default = {
        name = "default"
      }
    }
    streamtags = [
      "tags"
    ]
  }
}

resource "criblio_pack" "pipeline_pack" {
  id           = "pack-with-pipeline"
  group_id     = "default"
  description  = "Pack with pipeline"
  disabled     = true
  display_name = "Pack with pipeline"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

