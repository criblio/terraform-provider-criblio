resource "criblio_lakehouse_dataset_connection" "my_lakehousedatasetconnection" {
  lake_dataset_id = "myExistingLakeDatasetIdToCRUD"
  lakehouse_id    = "myExistingLakehouseId"
  request_body = {
    # ...
  }
}