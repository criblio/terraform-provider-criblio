resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  accelerated_fields = [
    "..."
  ]
  bucket_name = "...my_bucket_name..."
  cache_connection = {
    accelerated_fields = [
      "..."
    ]
    cache_ref          = "...my_cache_ref..."
    created_at         = 3.92
    migration_query_id = "...my_migration_query_id..."
    retention_in_days  = 8.82
  }
  deletion_started_at      = 5.22
  description              = "...my_description..."
  format                   = "json"
  id                       = "...my_id..."
  lake_id                  = "...my_lake_id..."
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
  view_name = "...my_view_name..."
}