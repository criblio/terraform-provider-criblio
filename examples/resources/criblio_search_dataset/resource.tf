resource "criblio_search_dataset" "my_searchdataset" {
  s3_dataset = {
    auto_detect_region = false
    bucket             = "example-search-logs"
    description        = "S3 search dataset"
    filter             = "true"
    id                 = "example_s3_dataset"
    metadata = {
      enable_acceleration = false
    }
    path        = "logs/*.log"
    provider_id = "S3"
    region      = "us-east-1"
    storage_classes = [
      "STANDARD",
    ]
    type = "s3"
  }
}
