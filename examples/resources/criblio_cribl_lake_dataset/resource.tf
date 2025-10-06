resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  accelerated_fields = [
    "fields",
    "to",
    "accelerate",
  ]
  bucket_name              = "my-Lake-bucket-name"
  description              = "My description for this beautiful lake dataset"
  format                   = "json"
  id                       = "myLakeDatasetId"
  lake_id                  = "default"
  retention_period_in_days = 30
  search_config = {
    datatypes = [
      "json",
      "parquet",
    ]
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
  }
}