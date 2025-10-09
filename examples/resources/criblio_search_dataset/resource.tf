resource "criblio_search_dataset" "my_searchdataset" {
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
      created             = "2025-10-06T12:00:00Z"
      enable_acceleration = true
      modified            = "2025-10-06T12:34:56Z"
      tags = [
        "prod",
        "pii",
      ]
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
}