resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  accelerated_fields = [
    "fields",
    "to",
    "accelerate",
  ]
  bucket_name              = "my-Lake-bucket-name"
  description              = "My description for this beautiful lake dataset"
  format                   = "json"
  id                       = "myLakeDatasetToCRUD"
  lake_id                  = "myUniqueLakeIdToCRUD"
  retention_period_in_days = 30
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