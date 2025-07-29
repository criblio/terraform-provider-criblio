resource "criblio_search_dataset" "my_searchdataset" {
  api_aws_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2022-10-29T14:30:18.805Z"
      enable_acceleration = true
      modified            = "2022-05-30T21:53:13.848Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    regions = [
      "..."
    ]
    type = "...my_type..."
  }
  api_azure_data_explorer_dataset = {
    cluster     = "...my_cluster..."
    database    = "...my_database..."
    description = "...my_description..."
    id          = "...my_id..."
    location    = "...my_location..."
    metadata = {
      created             = "2022-07-20T04:53:40.059Z"
      enable_acceleration = false
      modified            = "2022-11-08T20:24:44.740Z"
      tags = [
        "..."
      ]
    }
    provider_id              = "...my_provider_id..."
    table                    = "...my_table..."
    timestamp_field          = "...my_timestamp_field..."
    timestamp_field_contents = "...my_timestamp_field_contents..."
    type                     = "...my_type..."
  }
  api_azure_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2022-03-07T23:49:14.581Z"
      enable_acceleration = true
      modified            = "2022-02-07T03:37:04.121Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    subscription_ids = [
      "..."
    ]
    type = "...my_type..."
  }
  api_elastic_search_dataset = {
    description = "...my_description..."
    id          = "...my_id..."
    index       = "...my_index..."
    metadata = {
      created             = "2021-06-28T12:13:39.681Z"
      enable_acceleration = false
      modified            = "2022-04-18T23:19:36.636Z"
      tags = [
        "..."
      ]
    }
    provider_id     = "...my_provider_id..."
    timestamp_field = "...my_timestamp_field..."
    type            = "...my_type..."
  }
  api_gcp_dataset = {
    description = "...my_description..."
    endpoint_configs = [
      {
        endpoint_name = "...my_endpoint_name..."
        region        = "...my_region..."
      }
    ]
    id = "...my_id..."
    metadata = {
      created             = "2022-02-09T05:56:26.955Z"
      enable_acceleration = false
      modified            = "2021-03-15T00:41:17.850Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  api_google_workspace_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2021-10-11T07:46:25.909Z"
      enable_acceleration = false
      modified            = "2022-02-08T09:33:22.326Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  apihttp_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2021-07-10T14:32:53.487Z"
      enable_acceleration = true
      modified            = "2022-07-12T22:13:15.898Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  api_ms_graph_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2021-05-13T08:00:21.989Z"
      enable_acceleration = true
      modified            = "2020-10-03T02:51:12.191Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  api_okta_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2022-05-24T06:54:16.236Z"
      enable_acceleration = true
      modified            = "2022-01-26T18:51:46.178Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  api_open_search_dataset = {
    description = "...my_description..."
    id          = "...my_id..."
    index       = "...my_index..."
    metadata = {
      created             = "2020-05-21T05:08:14.917Z"
      enable_acceleration = true
      modified            = "2022-12-04T14:39:26.054Z"
      tags = [
        "..."
      ]
    }
    provider_id     = "...my_provider_id..."
    timestamp_field = "...my_timestamp_field..."
    type            = "...my_type..."
  }
  api_tailscale_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2022-05-01T05:02:01.580Z"
      enable_acceleration = false
      modified            = "2022-06-06T11:24:52.881Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  api_zoom_dataset = {
    description = "...my_description..."
    enabled_endpoints = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2020-04-05T02:30:33.102Z"
      enable_acceleration = true
      modified            = "2022-03-03T21:04:42.705Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  aws_security_lake_dataset = {
    description = "...my_description..."
    filter      = "...my_filter..."
    id          = "...my_id..."
    metadata = {
      created             = "2020-05-23T05:15:35.526Z"
      enable_acceleration = false
      modified            = "2022-11-20T18:24:50.959Z"
      tags = [
        "..."
      ]
    }
    path        = "...my_path..."
    provider_id = "...my_provider_id..."
    selected_buckets = [
      {
        name   = "...my_name..."
        region = "...my_region..."
      }
    ]
    type = "...my_type..."
  }
  azure_blob_dataset = {
    container_name = "...my_container_name..."
    description    = "...my_description..."
    extra_paths = [
      {
        container_name = "...my_container_name..."
        filter         = "...my_filter..."
        path           = "...my_path..."
      }
    ]
    filter = "...my_filter..."
    id     = "...my_id..."
    metadata = {
      created             = "2022-01-25T19:43:50.033Z"
      enable_acceleration = true
      modified            = "2022-11-05T01:20:11.254Z"
      tags = [
        "..."
      ]
    }
    path                   = "...my_path..."
    provider_id            = "...my_provider_id..."
    skip_event_time_filter = true
    storage_classes = [
      "..."
    ]
    type = "...my_type..."
  }
  click_house_dataset = {
    database    = "...my_database..."
    description = "...my_description..."
    id          = "...my_id..."
    metadata = {
      created             = "2022-11-26T03:12:16.392Z"
      enable_acceleration = false
      modified            = "2020-08-24T23:52:27.143Z"
      tags = [
        "..."
      ]
    }
    provider_id     = "...my_provider_id..."
    table           = "...my_table..."
    timestamp_field = "...my_timestamp_field..."
    type            = "...my_type..."
  }
  cribl_leader_dataset = {
    description = "...my_description..."
    extra_paths = [
      {
        filter = "...my_filter..."
        path   = "...my_path..."
      }
    ]
    filter = "...my_filter..."
    id     = "...my_id..."
    metadata = {
      created             = "2022-02-07T19:14:43.710Z"
      enable_acceleration = false
      modified            = "2022-01-10T14:50:59.619Z"
      tags = [
        "..."
      ]
    }
    path        = "...my_path..."
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  edge_dataset = {
    description = "...my_description..."
    filter      = "...my_filter..."
    fleets = [
      "..."
    ]
    id = "...my_id..."
    metadata = {
      created             = "2022-02-27T21:57:20.492Z"
      enable_acceleration = true
      modified            = "2022-10-15T11:38:53.888Z"
      tags = [
        "..."
      ]
    }
    path        = "...my_path..."
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  gcs_dataset = {
    bucket      = "...my_bucket..."
    description = "...my_description..."
    extra_paths = [
      {
        bucket = "...my_bucket..."
        filter = "...my_filter..."
        region = "...my_region..."
      }
    ]
    filter = "...my_filter..."
    id     = "...my_id..."
    metadata = {
      created             = "2022-12-28T19:28:14.296Z"
      enable_acceleration = false
      modified            = "2021-10-12T18:14:25.523Z"
      tags = [
        "..."
      ]
    }
    provider_id            = "...my_provider_id..."
    region                 = "...my_region..."
    skip_event_time_filter = true
    storage_classes = [
      "..."
    ]
    type = "...my_type..."
  }
  meta_dataset = {
    datasets = [
      "..."
    ]
    description = "...my_description..."
    id          = "...my_id..."
    metadata = {
      created             = "2022-11-13T06:28:12.532Z"
      enable_acceleration = false
      modified            = "2022-10-25T23:22:09.388Z"
      tags = [
        "..."
      ]
    }
    provider_id = "...my_provider_id..."
    type        = "...my_type..."
  }
  prometheus_dataset = {
    description             = "...my_description..."
    desired_num_data_points = 7957.12
    id                      = "...my_id..."
    metadata = {
      created             = "2022-06-15T23:01:13.239Z"
      enable_acceleration = true
      modified            = "2022-09-30T02:38:42.892Z"
      tags = [
        "..."
      ]
    }
    metric_name_pattern = "...my_metric_name_pattern..."
    provider_id         = "...my_provider_id..."
    step_mode           = "nearest_unit"
    type                = "...my_type..."
  }
  s3_dataset = {
    auto_detect_region = false
    bucket             = "...my_bucket..."
    description        = "...my_description..."
    extra_paths = [
      {
        auto_detect_region = true
        bucket             = "...my_bucket..."
        filter             = "...my_filter..."
        path               = "...my_path..."
        region             = "...my_region..."
      }
    ]
    filter = "...my_filter..."
    id     = "...my_id..."
    metadata = {
      created             = "2021-10-01T19:20:31.326Z"
      enable_acceleration = false
      modified            = "2021-02-10T22:10:18.036Z"
      tags = [
        "..."
      ]
    }
    path                   = "...my_path..."
    provider_id            = "...my_provider_id..."
    region                 = "...my_region..."
    skip_event_time_filter = true
    storage_classes = [
      "..."
    ]
    type = "...my_type..."
  }
  snowflake_dataset = {
    database    = "...my_database..."
    description = "...my_description..."
    id          = "...my_id..."
    metadata = {
      created             = "2021-03-22T17:16:32.304Z"
      enable_acceleration = false
      modified            = "2021-12-25T19:38:56.916Z"
      tags = [
        "..."
      ]
    }
    provider_id     = "...my_provider_id..."
    role            = "...my_role..."
    schema          = "...my_schema..."
    table           = "...my_table..."
    timestamp_field = "...my_timestamp_field..."
    type            = "...my_type..."
    warehouse       = "...my_warehouse..."
  }
}