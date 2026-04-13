resource "criblio_search_dataset" "my_searchdataset" {
  cribl_search_dataset = {
    additional_properties = "{ \"see\": \"documentation\" }"
    description           = "This is a generic dataset"
    id                    = "myGenericDatasetId"
    metadata = {
      created             = "2025-10-06T12:00:00Z"
      enable_acceleration = true
      modified            = "2025-10-06T12:34:56Z"
      tags = [
        "prod",
        "pii",
      ]
    }
    provider_id = "myProviderId"
    type        = "cribl_lake"
  }
}