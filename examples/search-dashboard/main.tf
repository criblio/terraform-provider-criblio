resource "criblio_search_dashboard" "my_searchdashboard" {
  id          = "sample_test_dashboard"
  name        = "Sample Test Dashboard"
  description = "A sample dashboard with several panels"
  created     = 1703123456789
  created_by  = "terraform"
  modified    = 1703123456789
  elements = [
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<42 | count"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          style           = false
          apply_threshold = false
          color_thresholds = {
            thresholds = [
              {
                color     = "#45850B"
                threshold = 30
              },
              {
                color     = "#EFDB23"
                threshold = 70
              },
              {
                color     = "#B20000"
                threshold = 100
              }
            ]
          }
          color_palette = 0
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          legend = {
            position = "Right"
            truncate = true
          }
          series    = []
          color     = "#0091eb"
          decimals  = 0
          label     = " The answer to life, the universe, and everything"
          timestats = false
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<100 \n| extend method=iif(event%3==0, 'POST', 'GET') \n| summarize count() by method"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          color_palette          = 0
          color_palette_reversed = false
          custom_data = {
            summarize_others = false
            series_count     = 1
          }
          legend = {
            position = "Right"
            truncate = true
          }
          onClick_action = {
            type = "None"
          }
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          series = [
            {
              y_axis_field = "count_"
              name         = "count_"
              color        = "#00CCCC"
            }
          ]
          timestats = false
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<42\n| count \n"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          color_thresholds = {
            thresholds = [
              {
                color     = "#45850B"
                threshold = 30
              },
              {
                color     = "#EFDB23"
                threshold = 70
              },
              {
                color     = "#B20000"
                threshold = 100
              }
            ]
          }
          legend = {
            position = "None"
            truncate = true
          }
          color_palette = 10
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          timestats = false
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<600 \n| extend _time=_time-rand(600), method=iif(event%2>0, \"GET\", \"POST\") \n| timestats span=1m count() by method"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          color_palette          = 0
          color_palette_reversed = false
          custom_data = {
            trellis       = false
            connect_nulls = "Leave gaps"
            stack         = false
            series_count  = 2
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          legend = {
            position = "None"
            truncate = true
          }
          onClick_action = {
            type = "None"
          }
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          timestats = true
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<600 \n| extend _time=_time-rand(600), method=iif(event%2>0, \"GET\", \"POST\") \n| timestats span=1m count()"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          color_palette          = 1
          color_palette_reversed = false
          custom_data = {
            trellis       = false
            connect_nulls = "Leave gaps"
            stack         = false
            series_count  = 1
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          legend = {
            position = "None"
            truncate = true
          }
          onClick_action = {
            type = "None"
          }
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          timestats = true
          series = [
            {
              y_axis_field = "count_"
              name         = "count_"
              color        = "#FF8042"
            }
          ]
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<100 \n| extend method=iif(event%3==0, 'POST', 'GET') \n| summarize count() by method"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          color_palette          = 9
          color_palette_reversed = false
          custom_data = {
            summarize_others = false
            series_count     = 1
          }
          legend = {
            position = "None"
            truncate = true
          }
          onClick_action = {
            type = "None"
          }
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          series = [
            {
              y_axis_field = "count_"
              name         = "count_"
              color        = "#9E0142"
            }
          ]
          timestats = false
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<600 \n| extend _time=_time-rand(600), method=iif(event%2>0, \"GET\", \"POST\") \n| timestats span=1m count() by method"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          color_palette          = 12
          color_palette_reversed = false
          custom_data = {
            trellis       = false
            connect_nulls = "Leave gaps"
            stack         = false
            series_count  = 2
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          legend = {
            position = "Right"
            truncate = true
          }
          onClick_action = {
            type = "None"
          }
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          timestats = true
          series = [
            {
              y_axis_field = "POST"
              name         = "POST"
              color        = "#56B4E9"
            },
            {
              y_axis_field = "GET"
              name         = "GET"
              color        = "#000000"
            }
          ]
        }
      }
    },
    {
      element = {
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
            query = "dataset=\"$vt_dummy\" event<20\n  | extend bytes = rand(10000), user = iif(event%3==0, 'admin', 'guest'), method=iif(event%3==0, 'POST', 'GET'), url = \"/api/v1/m/default_search/search/query?\"\n  | project-away dataset"
            earliest = {
              number = 1703123456789
            }
            latest = {
              number = 1703127056789
            }
            type = "inline"
          }
        }
        config = {
          onClick_action = {
            type = "None"
          }
          color_palette = 0
          data = {
            connect_nulls = "Leave gaps"
            stack         = false
          }
          x_axis = {
            label_orientation = 0
            position          = "Bottom"
          }
          y_axis = {
            position   = "Left"
            scale      = "Linear"
            split_line = true
          }
          legend = {
            position = "Right"
            truncate = true
          }
          series = [
            {
              y_axis_field = "status"
              name         = "status"
              color        = "#00CCCC"
            },
            {
              y_axis_field = "response_time"
              name         = "response_time"
              color        = "#ffa600"
            }
          ]
          axis = {
            x_axis = "time"
            y_axis = ["status", "response_time"]
          }
        }
      }
    }
  ]
}
