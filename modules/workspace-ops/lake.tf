resource "criblio_cribl_lake_dataset" "my_cribllakedataset" {
  id                       = "my_lake_dataset_tf5"
  description              = "My Cribl Lake Dataset"
  lake_id                  = "default"
  bucket_name              = "lake-${var.workspace}-${var.cloud_tenant}"
  format                   = "json"
  retention_period_in_days = 30

  search_config = {
    datatypes = [
      # Add your datatypes here if needed
    ]
    metadata = {
      enable_acceleration = false
      tags = [
        "test_tag"
      ]
    }
  }
}

resource "criblio_destination" "cribl_lake" {
  id = "cribl-lake-3"
  #   group_id = criblio_group.syslog_worker_group.id
  group_id = "default"
  output_cribl_lake = {
    id          = "cribl-lake-3"
    type        = "cribl_lake"
    description = "Cribl Lake destination for syslog data"
    dest_path   = "default_logs"
  }

  lifecycle {
    create_before_destroy = true
  }
}


resource "criblio_cribl_lake_house" "my_cribllakehouse" {
  description = "My Lakehouse for dataset"
  tier_size   = "medium"
  id          = "test-lakehouse-5"
}

# Check lakehouse status periodically
data "criblio_cribl_lake_house" "lakehouse_status" {
  id = criblio_cribl_lake_house.my_cribllakehouse.id
}

# Add a 10-minute delay before creating the dataset connection
resource "null_resource" "delay_before_connection" {
  provisioner "local-exec" {
    command = "sleep 600" # 600 seconds = 10 minutes
  }

  depends_on = [data.criblio_cribl_lake_house.lakehouse_status]
}

resource "criblio_lakehouse_dataset_connection" "my_cribllakehouse_dataset_connection" {
  lake_dataset_id = criblio_cribl_lake_dataset.my_cribllakedataset.id
  lakehouse_id    = criblio_cribl_lake_house.my_cribllakehouse.id

  depends_on = [null_resource.delay_before_connection]
}