resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  bucket_name = "lake-tfprovider-beautiful-nguyen-y8y4azd"
  description              = "my_description"
  format                   = "json"
  id                       = "my_lake_dataset"
  lake_id                  = "default"
  retention_period_in_days = 30
  search_config = {
    datatypes = [
    ]
    metadata = {
      created             = "2021-06-18T21:07:29.756Z"
      enable_acceleration = false
      modified            = "2022-10-01T07:28:47.966Z"
      tags = [
        "test_tag"
      ]
    }
  }
}