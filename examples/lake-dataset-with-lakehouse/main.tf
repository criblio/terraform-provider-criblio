resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  bucket_name = "lake-tfprovider-beautiful-nguyen-y8y4azd"
  description              = "my_description"
  format                   = "json"
  id                       = "my_lake_dataset_with_lakehouse_11"
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
  depends_on = [
    criblio_cribl_lake_house.my_cribllakehouse
  ]
}

resource "criblio_cribl_lake_house" "my_cribllakehouse" {
  description = "My Lakehouse for dataset"
  tier_size   = "medium"
  id          = "test-lakehouse-4"
}


# Add a 15-minute delay before creating the dataset connection, this is to ensure the lakehouse is ready
resource "null_resource" "delay_before_connection" {
  provisioner "local-exec" {
    command = "sleep 900"  # 900 seconds = 15 minutes
  }
}

resource "criblio_lakehouse_dataset_connection" "my_cribllakehouse_dataset_connection" {
  lake_dataset_id = criblio_cribl_lake_dataset.my_cribllakedataset.id
  lakehouse_id    = criblio_cribl_lake_house.my_cribllakehouse.id
  
  depends_on = [null_resource.delay_before_connection]
}


