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
      created             = "2021-10-01T19:20:31.326Z"
      enable_acceleration = false
      modified            = "2021-02-10T22:10:18.036Z"
      tags = [
        "test"
      ]
    }
    path                   = "logs/*.log"
    provider_id            = "S3"
    region                 = "us-east-1"
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
      created             = "2021-09-10T09:02:49.190Z"
      enable_acceleration = false
      modified            = "2021-08-05T01:30:06.408Z"
      tags = [
        "test"
      ]
    }
    provider_id = "cribl_lake"
    type        = "cribl_lake"
  }
}
*/
