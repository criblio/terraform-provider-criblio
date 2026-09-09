variable "existing_lakehouse_id" {
  description = "ID of an existing Cribl Lake Lakehouse. New Lakehouses can no longer be created."
  type        = string
}

data "criblio_cribl_lake_house" "existing" {
  id = var.existing_lakehouse_id
}

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
}

resource "criblio_lakehouse_dataset_connection" "my_cribllakehouse_dataset_connection" {
  lake_dataset_id = criblio_cribl_lake_dataset.my_cribllakedataset.id
  lakehouse_id    = data.criblio_cribl_lake_house.existing.id
}
