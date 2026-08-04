resource "criblio_search_dataset" "engine_dataset" {
  cribl_search_dataset = {
    auto_detect_region           = true
    breaker_rulesets             = ["Cribl Search"]
    engine                       = "tf_testing" #should exists before creating the dataset
    event_storage_schema_version = 1
    expected_relative_time_range = {
      latest = "1d"
    }
    filter = "true"
    id     = "test_engine_dataset_1"
    metadata = {
      enable_acceleration = false
    }
    partitioning_scheme    = "none"
    provider_id            = "lakehouse"
    retention_period       = 3650
    search_version         = "v1"
    skip_event_time_filter = false
    stale_channel_flush_ms = 10000
    storage_classes        = []
  }
}
