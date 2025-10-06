resource "criblio_lakehouse_dataset_connection" "my_lakehousedatasetconnection" {
  lake_dataset_id = "web-logs"
  lakehouse_id    = "lakehouse-01"
  request_body = {
    # ...
  }
}