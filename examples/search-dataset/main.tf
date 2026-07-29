resource "criblio_search_dataset" "v1" {
  s3_dataset = {
    auto_detect_region = false
    bucket             = "lake-main-beautiful-nguyen-y8y4azd"
    description        = "Search v1 dataset"
    id                 = "search_dataset_v1"
    provider_id        = "S3"
    region             = "us-west-2"
    search_version     = "v1"
    type               = "s3"
  }
}

resource "criblio_search_dataset" "v2" {
  s3_dataset = {
    description    = "Search v2 dataset"
    id             = "search_dataset_v2"
    provider_id    = "S3"
    search_version = "v2"
    storage_classes = [
      "STANDARD"
    ]
    type = "s3"

    paths = [
      {
        auto_detect_region  = true
        bucket              = "lake-main-beautiful-nguyen-y8y4azd"
        partitioning_scheme = "none"
        region              = "us-west-2"
        filters = [
          {
            data_type_id = "generic_ndjson"
            filter       = "**"
          }
        ]
      }
    ]
  }
}
