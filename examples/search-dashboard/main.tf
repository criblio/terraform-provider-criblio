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

resource "criblio_search_dashboard_category" "sre" {
  description = "SRE dashboards"
  id          = "SRE"
  is_pack     = false
  name        = "SRE"
}

resource "criblio_search_dashboard" "ecs_tasks" {
  name              = "ECS tasks"
  id                = "ecs_tasks"
  category          = criblio_search_dashboard_category.sre.id
  cache_ttl_seconds = 0
  refresh_rate      = 60000
  created           = 1733366400
  created_by        = "saas-operations"
  modified          = 1733366400
  modified_by       = "saas-operations"

  elements = [
    {
      element = {
        id      = "2tga2fara"
        inputId = "time"
        type    = "input.timerange"
        title   = "Time"
        layout = {
          x = 2
          y = 0
          w = 2
          h = 2
        }
        config = {
          defaultValue = jsonencode({
            earliest = "-3h"
            latest   = "now"
            timezone = "local"
          })
        }
      }
    },
    {
      element = {
        id    = "c6ydxqpde"
        type  = "chart.pie"
        title = "Tasks by group"
        layout = {
          x = 0
          y = 0
          w = 6
          h = 4
        }
        config = {
          colorPalette         = 0
          colorPaletteReversed = false
          customData = jsonencode({
            trellis      = false
            connectNulls = "Leave gaps"
            stack        = false
            dataFields   = ["group", "count_"]
            seriesCount  = 1
          })
          xAxis = jsonencode({
            labelOrientation = 0
            position         = "Bottom"
          })
          yAxis = jsonencode({
            position  = "Left"
            scale     = "Linear"
            splitLine = true
          })
          legend = jsonencode({
            position = "Right"
            truncate = true
          })
          onClickAction = jsonencode({
            type = "None"
          })
        }
        search = {
          search_query_inline = {
            type  = "inline"
            query = "dataset=\"ecs_event\" lastStatus=\"RUNNING\" | summarize count() by group"
            earliest = {
              str = "$time.earliest$"
            }
            latest = {
              str = "$time.latest$"
            }
            timezone = "$time.timezone$"
          }
        }
      }
    },
    {
      element = {
        id    = "678uwb26a-copy-copy-copy"
        type  = "chart.column"
        title = "Failed tasks by group/status "
        layout = {
          x = 6
          y = 0
          w = 6
          h = 4
        }
        config = {
          colorPalette         = 0
          colorPaletteReversed = false
          customData = jsonencode({
            trellis      = false
            connectNulls = "Leave gaps"
            stack        = false
            dataFields   = []
            seriesCount  = 4
          })
          xAxis = jsonencode({
            labelOrientation = 0
            position         = "Bottom"
          })
          yAxis = jsonencode({
            position  = "Left"
            scale     = "Linear"
            splitLine = true
          })
          legend = jsonencode({
            position = "Bottom"
            truncate = true
          })
          onClickAction = jsonencode({
            type = "None"
          })
        }
        search = {
          search_query_inline = {
            type  = "inline"
            query = "dataset=\"ecs_event\" containers.0.exitCode!=\"0\" lastStatus in (\"DEPROVISIONING\", \"STOPPED\") | timestats count() by group, lastStatus"
            earliest = {
              str = "$time.earliest$"
            }
            latest = {
              str = "$time.latest$"
            }
            timezone = "$time.timezone$"
          }
        }
      }
    },
    {
      element = {
        id    = "678uwb26a"
        type  = "chart.column"
        title = "Tenant-terraform tasks by status"
        layout = {
          x = 0
          y = 4
          w = 6
          h = 4
        }
        config = {
          colorPalette         = 0
          colorPaletteReversed = false
          customData = jsonencode({
            trellis      = false
            connectNulls = "Leave gaps"
            stack        = false
            dataFields   = ["_time", "RUNNING", "PROVISIONING", "PENDING", "STOPPED", "DEPROVISIONING"]
            seriesCount  = 5
          })
          xAxis = jsonencode({
            labelOrientation = 0
            position         = "Bottom"
          })
          yAxis = jsonencode({
            position  = "Left"
            scale     = "Linear"
            splitLine = true
          })
          legend = jsonencode({
            position = "Bottom"
            truncate = true
          })
          onClickAction = jsonencode({
            type = "None"
          })
        }
        search = {
          search_query_inline = {
            type  = "inline"
            query = "dataset=\"ecs_event\" group=\"family:tenant-terraform\" | timestats count() by lastStatus"
            earliest = {
              str = "$time.earliest$"
            }
            latest = {
              str = "$time.latest$"
            }
            timezone = "$time.timezone$"
          }
        }
      }
    },
    {
      element = {
        id    = "678uwb26a-copy"
        type  = "chart.column"
        title = "Cleanup tasks by status"
        layout = {
          x = 6
          y = 4
          w = 6
          h = 4
        }
        config = {
          colorPalette         = 0
          colorPaletteReversed = false
          customData = jsonencode({
            trellis      = false
            connectNulls = "Leave gaps"
            stack        = false
            dataFields   = ["_time", "PENDING", "PROVISIONING", "RUNNING", "DEPROVISIONING", "STOPPED"]
            seriesCount  = 5
          })
          xAxis = jsonencode({
            labelOrientation = 0
            position         = "Bottom"
          })
          yAxis = jsonencode({
            position  = "Left"
            scale     = "Linear"
            splitLine = true
          })
          legend = jsonencode({
            position = "Bottom"
            truncate = true
          })
          onClickAction = jsonencode({
            type = "None"
          })
        }
        search = {
          search_query_inline = {
            type  = "inline"
            query = "dataset=\"ecs_event\" group=\"family:cleanup\" | timestats count() by lastStatus"
            earliest = {
              str = "$time.earliest$"
            }
            latest = {
              str = "$time.latest$"
            }
            timezone = "$time.timezone$"
          }
        }
      }
    },
    {
      element = {
        id    = "678uwb26a-copy-copy"
        type  = "chart.column"
        title = "Tenant-exec tasks by status"
        layout = {
          x = 0
          y = 8
          w = 6
          h = 4
        }
        config = {
          colorPalette         = 0
          colorPaletteReversed = false
          customData = jsonencode({
            trellis      = false
            connectNulls = "Leave gaps"
            stack        = false
            dataFields   = []
            seriesCount  = 5
          })
          xAxis = jsonencode({
            labelOrientation = 0
            position         = "Bottom"
          })
          yAxis = jsonencode({
            position  = "Left"
            scale     = "Linear"
            splitLine = true
          })
          legend = jsonencode({
            position = "Bottom"
            truncate = true
          })
          onClickAction = jsonencode({
            type = "None"
          })
        }
        search = {
          search_query_inline = {
            type  = "inline"
            query = "dataset=\"ecs_event\" group=\"family:tenant-exec\" | timestats count() by lastStatus"
            earliest = {
              str = "$time.earliest$"
            }
            latest = {
              str = "$time.latest$"
            }
            timezone = "$time.timezone$"
          }
        }
      }
    }
  ]
}