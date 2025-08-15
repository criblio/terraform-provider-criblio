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
  id       = "cribl-lake-3"
#   group_id = criblio_group.syslog_worker_group.id
    group_id = "default"
  output_cribl_lake = {
    id          = "cribl-lake-3"
    type        = "cribl_lake"
    description = "Cribl Lake destination for syslog data"
    disabled    = false
    streamtags  = ["syslog", "lake"]
    dest_path   = "default_logs"
    format      = "json"
    compress    = "gzip"
    add_id_to_stage_path = true
    aws_authentication_method = "auto"
    base_file_name = "CriblOut"
    file_name_suffix = "'.gz'"
    max_file_size_mb = 32
    max_open_files = 100
    write_high_water_mark = 64
    on_backpressure = "block"
    deadletter_enabled = false
    on_disk_full_backpressure = "block"
    max_file_open_time_sec = 300
    max_file_idle_time_sec = 30
    verify_permissions = true
    max_closing_files_to_backpressure = 100
    max_concurrent_file_parts = 1
    empty_dir_cleanup_sec = 300
    max_retry_num = 20
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
    command = "sleep 600"  # 600 seconds = 10 minutes
  }
  
  depends_on = [data.criblio_cribl_lake_house.lakehouse_status]
}

resource "criblio_lakehouse_dataset_connection" "my_cribllakehouse_dataset_connection" {
  lake_dataset_id = criblio_cribl_lake_dataset.my_cribllakedataset.id
  lakehouse_id    = criblio_cribl_lake_house.my_cribllakehouse.id
  
  depends_on = [null_resource.delay_before_connection]
}