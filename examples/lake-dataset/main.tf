resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  bucket_name              = "lake-tfprovider-beautiful-nguyen-y8y4azd"
  description              = "my_description"
  format                   = "json"
  id                       = "my_lake_dataset"
  lake_id                  = "default"
  retention_period_in_days = 30
  search_config = {
    metadata = {
      tags = [
        "test_tag"
      ]
    }
  }
}
