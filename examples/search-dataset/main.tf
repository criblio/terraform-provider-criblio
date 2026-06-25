resource "criblio_search_dataset" "my_s3_dataset" {
  s3_dataset = {
    auto_detect_region = false
    bucket             = "test_bucket"
    description        = "test"
    extra_paths = [
      {
        auto_detect_region = false
        bucket             = "test_bucket"
        filter             = "test"
        path               = "logs/*.log"
        region             = "us-east-1"
      }
    ]
    filter = "test"
    id     = "S3"
    metadata = {
      enable_acceleration = false
    }
    path        = "logs/*.log"
    provider_id = "S3"
    region      = "us-east-1"
    storage_classes = [
      "STANDARD"
    ]
    type = "s3"
  }
}

/*
resource "criblio_search_dataset" "my_cribl_lake_dataset" {
  cribl_lake_dataset = {
    description = "test"
    id          = "test_cribl_lake_dataset"
    metadata = {
      enable_acceleration = false
    }
    provider_id = "cribl_lake"
    type        = "cribl_lake"
  }
}
*/
