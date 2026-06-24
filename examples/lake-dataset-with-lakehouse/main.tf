resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  bucket_name              = "lake-main-beautiful-nguyen-y8y4azd"
  description              = "my_description"
  format                   = "json"
  id                       = "my_lake_dataset_with_lakehouse_22"
  lake_id                  = "default"
  retention_period_in_days = 30
  search_config = {
    metadata = {
      tags = [
        "test_tag"
      ]
    }
  }
  depends_on = [
    criblio_cribl_lake_house.my_cribllakehouse
  ]
}

resource "criblio_cribl_lake_house" "my_cribllakehouse" {
  description = "My Lakehouse for dataset"
  tier_size   = "medium"
  id          = "test-lakehouse-10"
}

resource "criblio_lakehouse_dataset_connection" "my_cribllakehouse_dataset_connection" {
  lake_dataset_id = criblio_cribl_lake_dataset.my_cribllakedataset.id
  lakehouse_id    = criblio_cribl_lake_house.my_cribllakehouse.id
}
