resource "criblio_lakehouse_dataset_connection" "my_lakehousedatasetconnection" {
  lake_dataset_id = "...my_lake_dataset_id..."
  lakehouse_id    = "...my_lakehouse_id..."
  request_body = {
    # ...
  }
}