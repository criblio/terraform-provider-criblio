resource "criblio_search_dataset" "my_searchdataset" {
  api_aws_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "CloudTrail",
      "CloudWatch",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    regions = [
      "us-east-1",
      "us-west-2",
    ]
    type = "cribl_lake"
  }
  api_azure_data_explorer_dataset = {
    cluster     = "myadxcluster"
    database    = "logsdb"
    description = "This is a generic dataset"
    id          = "myGenericDatasetId"
    location    = "eastus2"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id              = "myProviderId"
    table                    = "logs"
    timestamp_field          = "event_time"
    timestamp_field_contents = "kusto_datetime"
    type                     = "cribl_lake"
  }
  api_azure_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "ActivityLogs",
      "SignInLogs",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    subscription_ids = [
      "sub-12345",
      "sub-67890",
    ]
    type = "cribl_lake"
  }
  api_elastic_search_dataset = {
    description = "This is a generic dataset"
    id          = "myGenericDatasetId"
    index       = "metrics-*"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id     = "myProviderId"
    timestamp_field = "@timestamp"
    type            = "cribl_lake"
  }
  api_gcp_dataset = {
    description = "This is a generic dataset"
    endpoint_configs = [
      {
        endpoint_name = "pubsub"
        region        = "us-central1"
      }
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  api_google_workspace_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "login",
      "admin",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  api_ms_graph_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "auditLogs",
      "signIns",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  api_okta_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "users",
      "events",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  api_open_search_dataset = {
    description = "This is a generic dataset"
    id          = "myGenericDatasetId"
    index       = "logs-*"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id     = "myProviderId"
    timestamp_field = "@timestamp"
    type            = "cribl_lake"
  }
  api_tailscale_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "devices",
      "connections",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  api_zoom_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "users",
      "meetings",
      "recordings",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  apihttp_dataset = {
    description = "This is a generic dataset"
    enabled_endpoints = [
      "getUsers",
      "getEvents",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  aws_security_lake_dataset = {
    description = "This is a generic dataset"
    filter      = "path.includes('CloudTrail')"
    id          = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    path        = "AWS/CloudTrail/v1/region=us-east-1/accountId=123456789012/eventDay=20251001/"
    provider_id = "myProviderId"
    selected_buckets = [
      {
        name   = "aws-securitylake-us-east-1"
        region = "us-east-1"
      }
    ]
    type = "cribl_lake"
  }
  azure_blob_dataset = {
    container_name = "my-container"
    description    = "This is a generic dataset"
    extra_paths = [
      {
        container_name = "my-other-container"
        filter         = "truthy"
        path           = "archive/2025/"
      }
    ]
    filter = "path.endsWith('.json')"
    id     = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    path                   = "logs/2025/10/01/"
    provider_id            = "myProviderId"
    skip_event_time_filter = true
    storage_classes = [
      "glacier",
      "standard",
    ]
    type = "cribl_lake"
  }
  click_house_dataset = {
    database    = "analytics_db"
    description = "This is a generic dataset"
    id          = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id     = "myProviderId"
    table           = "logs"
    timestamp_field = "event_time"
    type            = "cribl_lake"
  }
  cribl_leader_dataset = {
    description = "This is a generic dataset"
    extra_paths = [
      {
        filter = "truthy"
        path   = "path/to/dir"
      }
    ]
    filter = "filename.endsWith('.log')"
    id     = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    path        = "/var/log/cribl"
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  edge_dataset = {
    description = "This is a generic dataset"
    filter      = "size > 0"
    fleets = [
      "fleetA",
      "fleetB",
    ]
    id = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    path        = "/data/edge"
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  gcs_dataset = {
    bucket      = "my-gcs-bucket"
    description = "This is a generic dataset"
    extra_paths = [
      {
        bucket = "my-extra-bucket"
        filter = "path.includes('archive')"
        region = "europe-west1"
      }
    ]
    filter = "path.endsWith('.json')"
    id     = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id            = "myProviderId"
    region                 = "us-central1"
    skip_event_time_filter = false
    storage_classes = [
      "STANDARD",
      "NEARLINE",
    ]
    type = "cribl_lake"
  }
  meta_dataset = {
    datasets = [
      "dataset1",
      "dataset2",
    ]
    description = "This is a generic dataset"
    id          = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
  prometheus_dataset = {
    description             = "This is a generic dataset"
    desired_num_data_points = 500
    id                      = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    metric_name_pattern = "http_requests_total"
    provider_id         = "myProviderId"
    step_mode           = "nearest_unit"
    type                = "cribl_lake"
  }
  s3_dataset = {
    auto_detect_region = false
    bucket             = "...my_bucket..."
    description        = "This is a generic dataset"
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
    id     = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    path                   = "...my_path..."
    provider_id            = "myProviderId"
    region                 = "...my_region..."
    skip_event_time_filter = true
    storage_classes = [
      "..."
    ]
    type = "cribl_lake"
  }
  snowflake_dataset = {
    database    = "analytics_db"
    description = "This is a generic dataset"
    id          = "myGenericDatasetId"
    metadata = {
      earliest            = "2025-09-30T13:41:44Z"
      enable_acceleration = true
      field_list = [
        "field1",
        "field2",
      ]
      latest_run_info = {
        earliest_scanned_time = 1759324416
        finished_at           = 1759325416
        latest_scanned_time   = 1759326416
        object_count          = 5000
      }
      scan_mode = "detailed"
    }
    provider_id     = "myProviderId"
    role            = "SYSADMIN"
    schema          = "public"
    table           = "logs"
    timestamp_field = "event_time"
    type            = "cribl_lake"
    warehouse       = "COMPUTE_WH"
  }
}