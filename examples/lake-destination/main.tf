# Cribl Lake Destination Configuration
resource "criblio_destination" "cribl_lake" {
  id       = "cribl-lake-11"
  group_id = "default"

  output_cribl_lake = {
    id                                = "cribl-lake-11"
    type                              = "cribl_lake"
    description                       = "Cribl Lake destination"
    disabled                          = false
    streamtags                        = ["test", "lake"]
    dest_path                         = "default_logs"
    format                            = "json"
    compress                          = "gzip"
    add_id_to_stage_path              = true
    aws_authentication_method         = "auto"
    base_file_name                    = "CriblOut"
    max_file_size_mb                  = 32
    max_open_files                    = 100
    write_high_water_mark             = 64
    on_backpressure                   = "block"
    deadletter_enabled                = false
    on_disk_full_backpressure         = "block"
    max_file_open_time_sec            = 300
    max_file_idle_time_sec            = 30
    verify_permissions                = true
    max_closing_files_to_backpressure = 100
    max_concurrent_file_parts         = 1
    empty_dir_cleanup_sec             = 300
    max_retry_num                     = 20
  }
}