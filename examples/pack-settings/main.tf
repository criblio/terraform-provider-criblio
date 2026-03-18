resource "criblio_pack" "edge_supervisor_health" {
  id           = "EdgeSupervisorHealth"
  group_id     = "default_fleet"
  description  = "EdgeSupervisorHealth - Monitors Supervisor process health on Search Deux instances"
  disabled     = false
  display_name = "EdgeSupervisorHealth"
  version      = "1.0.0"
  author       = "Observability Team"
}

# Pipeline for EdgeSupervisorHealth - processes Supervisor health check data
resource "criblio_pack_pipeline" "edge_supervisor_health_main" {
  id       = "main"
  pack     = criblio_pack.edge_supervisor_health.id
  group_id = "default_fleet"

  depends_on = [
    criblio_pack.edge_supervisor_health
  ]

  conf = {
    streamtags = []

    functions = [
      {
        id       = "serde"
        filter   = "true"
        disabled = false
        conf = jsonencode({
          mode     = "extract"
          type     = "json"
          srcField = "_raw"
        })
      },
      {
        id       = "eval"
        filter   = "true"
        disabled = false
        conf = jsonencode({
          add = [
            {
              name     = "_time"
              value    = "Date.now()/1000"
              disabled = false
            },
            {
              name     = "value"
              value    = "Number(return_code)"
              disabled = false
            },
            {
              name     = "service"
              value    = "'search-supervisor'"
              disabled = false
            }
          ]
          remove = [
            "cribl_*",
            "source",
            "_raw",
            "return_code"
          ]
        })
      },
      {
        id       = "publish_metrics"
        filter   = "true"
        disabled = false
        conf = jsonencode({
          overwrite = false
          dimensions = [
            "status",
            "saas_domain",
            "service",
            "fleet",
            "tenantId",
            "workspace",
            "instance",
            "cloudProvider",
            "instance_type",
            "region",
            "container_tag",
            "container_state"
          ]
          removeMetrics    = []
          removeDimensions = []
          fields = [
            {
              metricType   = "gauge"
              inFieldName  = "value"
              outFieldExpr = "'health_check'"
            }
          ]
        })
      },
    ]
  }
}
