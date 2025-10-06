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
      created             = "2025-10-06T12:00:00Z"
      enable_acceleration = true
      modified            = "2025-10-06T12:34:56Z"
      tags = [
        "prod",
        "pii",
      ]
    }
  }
}