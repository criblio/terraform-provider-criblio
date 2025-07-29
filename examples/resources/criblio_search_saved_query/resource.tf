resource "criblio_search_saved_query" "my_searchsavedquery" {
  chart_config = {
    apply_threshold = false
    axis = {
      x_axis = "...my_x_axis..."
      y_axis = [
        "..."
      ]
      y_axis_excluded = [
        "..."
      ]
    }
    color                  = "...my_color..."
    color_palette          = 0.59
    color_palette_reversed = true
    color_thresholds = {
      thresholds = [
        {
          color     = "...my_color..."
          threshold = 4.62
        }
      ]
    }
    custom_data = {
      connect_nulls = "...my_connect_nulls..."
      data_fields = [
        "..."
      ]
      is_point_color               = false
      limit_to_top_n               = 0.73
      lines                        = false
      name_field                   = "...my_name_field..."
      point_color_palette          = 8.61
      point_color_palette_reversed = false
      point_scale = {
        str = "...my_str..."
      }
      point_scale_data_field = "...my_point_scale_data_field..."
      series_count           = 0.84
      split_by               = "...my_split_by..."
      stack                  = true
      summarize_others       = true
      trellis                = true
    }
    decimals = 4.21
    label    = "...my_label..."
    legend = {
      position = "...my_position..."
      truncate = false
    }
    map_details = {
      latitude_field  = "...my_latitude_field..."
      longitude_field = "...my_longitude_field..."
      map_source_id   = "...my_map_source_id..."
      map_type        = "...my_map_type..."
      name_field      = "...my_name_field..."
      point_scale = {
        number = 9.31
      }
      value_field = "...my_value_field..."
    }
    on_click_action = {
      search                      = "...my_search..."
      selected_dashboard_id       = "...my_selected_dashboard_id..."
      selected_input_id           = "...my_selected_input_id..."
      selected_link_id            = "...my_selected_link_id..."
      selected_timerange_input_id = "...my_selected_timerange_input_id..."
      type                        = "...my_type..."
    }
    prefix    = "...my_prefix..."
    separator = false
    series = [
      {
        color = "...my_color..."
        data = [
          {
            # ...
          }
        ]
        map          = "...my_map..."
        name         = "...my_name..."
        type         = "gauge"
        y_axis_field = "...my_y_axis_field..."
      }
    ]
    should_apply_user_chart_settings = true
    style                            = false
    suffix                           = "...my_suffix..."
    type                             = "...my_type..."
    x_axis = {
      data_field        = "...my_data_field..."
      inverse           = false
      label_interval    = "...my_label_interval..."
      label_orientation = 0.15
      name              = "...my_name..."
      offset            = 1.03
      position          = "...my_position..."
      type              = "...my_type..."
    }
    y_axis = {
      data_field = [
        "..."
      ]
      interval   = 3.09
      max        = 0.68
      min        = 6.09
      position   = "...my_position..."
      scale      = "...my_scale..."
      split_line = false
      type       = "...my_type..."
    }
  }
  description      = "...my_description..."
  display_username = "...my_display_username..."
  earliest         = "...my_earliest..."
  id               = "...my_id..."
  is_private       = true
  is_system        = true
  latest           = "...my_latest..."
  lib              = "custom"
  name             = "...my_name..."
  query            = "...my_query..."
  resolved_dataset_ids = [
    "..."
  ]
  sample_rate = 2.25
  schedule = {
    cron_schedule = "...my_cron_schedule..."
    enabled       = true
    keep_last_n   = 7.78
    notifications = {
      disabled = false
      items = [
        {
          condition = "...my_condition..."
          conf = {
            # ...
          }
          disabled = false
          group    = "...my_group..."
          id       = "...my_id..."
          metadata = [
            {
              name  = "...my_name..."
              value = "...my_value..."
            }
          ]
          target_configs = [
            {
              conf = {
                body = "...my_body..."
                email_recipient = {
                  bcc = "Alec_Jenkins@yahoo.com"
                  cc  = "Mortimer.Von@hotmail.com"
                  to  = "Irwin5@gmail.com"
                }
                subject = "...my_subject..."
              }
              id = "...my_id..."
            }
          ]
          targets = [
            "..."
          ]
        }
      ]
    }
    resume_missed  = false
    resume_on_boot = true
    tz             = "...my_tz..."
  }
  table_config = "...my_table_config..."
  user         = "...my_user..."
}