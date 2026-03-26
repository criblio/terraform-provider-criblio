resource "criblio_search_dashboard" "my_searchdashboard" {
  id          = "sample_test_dashboard"
  name        = "Sample Test Dashboard"
  description = "A sample dashboard with several panels"

  elements = [
    # 1) Single Value Visualization (counter.single)
    {
      dashboard_element_visualization = {
        id    = "uhyck3nbk"
        type  = "counter.single"
        title = "Single Value Visualization"

        layout = {
          x = 0
          y = 0
          w = 4
          h = 3
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = "dataset=\"$vt_dummy\" event<42 | count"

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "style": false,
            "applyThreshold": false,
            "colorThresholds": {
              "thresholds": [
                { "color": "#45850B", "threshold": 30 },
                { "color": "#EFDB23", "threshold": 70 },
                { "color": "#B20000", "threshold": 100 }
              ]
            },
            "colorPalette": 0,
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "legend": {
              "position": "Right",
              "truncate": true
            },
            "series": [],
            "color": "#0091eb",
            "decimals": 0,
            "label": " The answer to life, the universe, and everything",
            "timestats": false
          }
          EOT
        }
      }
    },

    # 2) Donut Chart Visualization (chart.pie)
    {
      dashboard_element_visualization = {
        id    = "arr3nh2me"
        type  = "chart.pie"
        title = "Donut Chart Visualization"

        layout = {
          x = 4
          y = 0
          w = 4
          h = 3
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<100 
            | extend method=iif(event%3==0, 'POST', 'GET') 
            | summarize count() by method
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "colorPalette": 0,
            "colorPaletteReversed": false,
            "customData": {
              "summarizeOthers": false,
              "seriesCount": 1
            },
            "legend": {
              "position": "Right",
              "truncate": true
            },
            "onClickAction": {
              "type": "None"
            },
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "series": [
              {
                "yAxisField": "count_",
                "name": "count_",
                "color": "#00CCCC"
              }
            ],
            "timestats": false
          }
          EOT
        }
      }
    },

    # 3) Gauge Chart Visualization (chart.gauge)
    {
      dashboard_element_visualization = {
        id    = "x8878143y"
        type  = "chart.gauge"
        title = "Gauge Chart Visualization"

        layout = {
          x = 8
          y = 0
          w = 4
          h = 3
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<42
            | count 
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "colorThresholds": {
              "thresholds": [
                { "color": "#45850B", "threshold": 30 },
                { "color": "#EFDB23", "threshold": 70 },
                { "color": "#B20000", "threshold": 100 }
              ]
            },
            "legend": {
              "position": "None",
              "truncate": true
            },
            "colorPalette": 10,
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "timestats": false
          }
          EOT
        }
      }
    },

    # 4) Area Chart Visualization (chart.area)
    {
      dashboard_element_visualization = {
        id    = "ndkk3w9ph"
        type  = "chart.area"
        title = "Area Chart Visualization"

        layout = {
          x = 0
          y = 3
          w = 6
          h = 3
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<600 
            | extend _time=_time-rand(600), method=iif(event%2>0, "GET", "POST") 
            | timestats span=1m count() by method
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "colorPalette": 0,
            "colorPaletteReversed": false,
            "customData": {
              "trellis": false,
              "connectNulls": "Leave gaps",
              "stack": false,
              "seriesCount": 2
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "legend": {
              "position": "None",
              "truncate": true
            },
            "onClickAction": {
              "type": "None"
            },
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "timestats": true
          }
          EOT
        }
      }
    },

    # 5) Bar Chart Visualization (chart.column)
    {
      dashboard_element_visualization = {
        id    = "0rfhfiufp"
        type  = "chart.column"
        title = "Bar Chart Visualization"

        layout = {
          x = 6
          y = 3
          w = 6
          h = 3
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<600 
            | extend _time=_time-rand(600), method=iif(event%2>0, "GET", "POST") 
            | timestats span=1m count()
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "colorPalette": 1,
            "colorPaletteReversed": false,
            "customData": {
              "trellis": false,
              "connectNulls": "Leave gaps",
              "stack": false,
              "seriesCount": 1
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "legend": {
              "position": "None",
              "truncate": true
            },
            "onClickAction": {
              "type": "None"
            },
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "timestats": true,
            "series": [
              {
                "yAxisField": "count_",
                "name": "count_",
                "color": "#FF8042"
              }
            ]
          }
          EOT
        }
      }
    },

    # 6) Funnel Chart Visualization (chart.funnel)
    {
      dashboard_element_visualization = {
        id    = "dbkmmais5"
        type  = "chart.funnel"
        title = "Funnel Chart Visualization"

        layout = {
          x = 0
          y = 6
          w = 4
          h = 4
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<100 
            | extend method=iif(event%3==0, 'POST', 'GET') 
            | summarize count() by method
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "colorPalette": 9,
            "colorPaletteReversed": false,
            "customData": {
              "summarizeOthers": false,
              "seriesCount": 1
            },
            "legend": {
              "position": "None",
              "truncate": true
            },
            "onClickAction": {
              "type": "None"
            },
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "series": [
              {
                "yAxisField": "count_",
                "name": "count_",
                "color": "#9E0142"
              }
            ],
            "timestats": false
          }
          EOT
        }
      }
    },

    # 7) Line Chart Visualization (chart.line)
    {
      dashboard_element_visualization = {
        id    = "qtifqfly4"
        type  = "chart.line"
        title = "Line Chart Visualization"

        layout = {
          x = 4
          y = 6
          w = 8
          h = 4
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<600 
            | extend _time=_time-rand(600), method=iif(event%2>0, "GET", "POST") 
            | timestats span=1m count() by method
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "colorPalette": 12,
            "colorPaletteReversed": false,
            "customData": {
              "trellis": false,
              "connectNulls": "Leave gaps",
              "stack": false,
              "seriesCount": 2
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "legend": {
              "position": "Right",
              "truncate": true
            },
            "onClickAction": {
              "type": "None"
            },
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "timestats": true,
            "series": [
              {
                "yAxisField": "POST",
                "name": "POST",
                "color": "#56B4E9"
              },
              {
                "yAxisField": "GET",
                "name": "GET",
                "color": "#000000"
              }
            ]
          }
          EOT
        }
      }
    },

    # 8) Raw Events Visualization (list.events)
    {
      dashboard_element_visualization = {
        id    = "uxwdqxfsa"
        type  = "list.events"
        title = "Raw Events Visualization"

        layout = {
          x = 0
          y = 10
          w = 12
          h = 4
        }

        search = {
          search_query_inline = {
            type  = "inline"
            query = <<-EOT
            dataset="$vt_dummy" event<20
              | extend bytes = rand(10000), user = iif(event%3==0, 'admin', 'guest'), method=iif(event%3==0, 'POST', 'GET'), url = "/api/v1/m/default_search/search/query?"
              | project-away dataset
            EOT

            earliest = {
              str = "-1h"
            }

            latest = {
              str = "now"
            }
          }
        }

        config = {
          json = <<-EOT
          {
            "onClickAction": {
              "type": "None"
            },
            "colorPalette": 0,
            "data": {
              "connectNulls": "Leave gaps",
              "stack": false
            },
            "xAxis": {
              "labelOrientation": 0,
              "position": "Bottom"
            },
            "yAxis": {
              "position": "Left",
              "scale": "Linear",
              "splitLine": true
            },
            "legend": {
              "position": "Right",
              "truncate": true
            },
            "series": [
              {
                "yAxisField": "status",
                "name": "status",
                "color": "#00CCCC"
              },
              {
                "yAxisField": "response_time",
                "name": "response_time",
                "color": "#ffa600"
              }
            ],
            "axis": {
              "xAxis": "time",
              "yAxis": ["status", "response_time"]
            }
          }
          EOT
        }
      }
    }
  ]
}

resource "criblio_search_dashboard" "search_dashboard_default_search_l3wro0ghb" {
  cache_ttl_seconds = 0
  description       = ""
  elements = [{
    dashboard_element_input = {
      config = {
        defaultValue = "{\"earliest\":\"-6h\",\"latest\":\"now\",\"timezone\":\"local\"}"
      }
      hide_panel       = false
      horizontal_chart = false
      id               = "7z8tf73fk"
      input_id         = "timerange"
      layout = {
        h = 2
        w = 2
        x = 0
        y = 0
      }
      title = "Time Range"
      type  = "input.timerange"
    }
    }, {
    dashboard_element_visualization = {
      config = {
        color = "\"#336666\""
      }
      hide_panel       = false
      horizontal_chart = false
      id               = "tx2w91r2y"
      layout = {
        h = 2
        w = 4
        x = 0
        y = 0
      }
      search = {
        search_query_inline = {
          earliest = {
            str = "-1h"
          }
          latest = {
            str = "now"
          }
          parent_search_id = "7vu8glgfy"
          query            = "\n | where message in (\"query started\") \n | summarize distinctTenantCount=dcount(tenantId)\n "
          sample_rate      = 1
          timezone         = null
          type             = "inline"
        }
      }
      title = "Num of Orgs Interacted with SI in the Time Period"
      type  = "counter.single"
    }
    }, {
    dashboard_element_visualization = {
      config = {
        color = "\"#2fcdcc\""
      }
      hide_panel       = false
      horizontal_chart = false
      id               = "tx2w91r2y-copy"
      layout = {
        h = 2
        w = 4
        x = 4
        y = 0
      }
      search = {
        search_query_inline = {
          earliest = {
            str = "-1h"
          }
          latest = {
            str = "now"
          }
          parent_search_id = "7vu8glgfy"
          query            = "\n | where message in (\"query started\") \n | summarize distinctWorkspaceCount=dcount(workspace) by tenantId\n | summarize totalDistinctWorkspaceCount=sum(distinctWorkspaceCount)\n \n "
          sample_rate      = 1
          timezone         = null
          type             = "inline"
        }
      }
      title = "Num of Workspace Interacted with SI in the Time Period"
      type  = "counter.single"
    }
    }, {
    dashboard_element_visualization = {
      config           = {}
      hide_panel       = false
      horizontal_chart = false
      id               = "g1izebyro"
      layout = {
        h = 4
        w = 12
        x = 0
        y = 2
      }
      search = {
        search_query_inline = {
          earliest = {
            str = "$timerange.earliest$"
          }
          latest = {
            str = "$timerange.latest$"
          }
          parent_search_id = "7vu8glgfy"
          query            = "\n | where message in (\"query started\", \"Query failed\", \"Query completed\")\n | extend status=case(message==\"Query completed\",\"success\",\n                      message==\"Query failed\",\"failure\",\n                      message==\"query started\",\"requested\", \n                      \"others\")\n | summarize requested=countif(status==\"requested\"), success=countif(status==\"success\"), failure=countif(status==\"failure\") by tenantId, workspace\n | extend error_rate=100*failure/requested"
          sample_rate      = 1
          timezone         = "$timerange.timezone$"
          type             = "inline"
        }
      }
      title = "Query Stats"
      type  = "list.table"
    }
    }, {
    dashboard_element_visualization = {
      config           = {}
      hide_panel       = true
      horizontal_chart = false
      id               = "7vu8glgfy"
      layout = {
        h = 4
        w = 12
        x = 0
        y = 6
      }
      search = {
        search_query_inline = {
          earliest = {
            str = "$timerange.earliest$"
          }
          latest = {
            str = "$timerange.latest$"
          }
          parent_search_id = ""
          query            = "dataset=\"$vt_dummy\" event<42 | count"
          sample_rate      = 1
          timezone         = "$timerange.timezone$"
          type             = "inline"
        }
      }
      title = "Parent Query"
      type  = "chart.column"
    }
  }]
  id           = "test_dashboard"
  name         = "Test Dashboard"
  refresh_rate = 0
  schedule = {
    cron_schedule = "0 0 * * *"
    enabled       = false
    keep_last_n   = 2
    notifications = {
      disabled = false
    }
    tz = "UTC"
  }
}