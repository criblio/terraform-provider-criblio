resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  accelerated_fields = [
    "..."
  ]
  bucket_name              = "...my_bucket_name..."
  description              = "...my_description..."
  format                   = "json"
  id                       = "myLakeDatasetToCRUD"
  lake_id                  = "myUniqueLakeIdToCRUD"
  retention_period_in_days = 3.75
  search_config = {
    datatypes = [
      "..."
    ]
    metadata = {
      created             = "2021-06-18T21:07:29.756Z"
      enable_acceleration = false
      modified            = "2022-10-01T07:28:47.966Z"
      tags = [
        "..."
      ]
    }
  }
}