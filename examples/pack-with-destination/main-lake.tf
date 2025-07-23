resource "criblio_pack" "my_pack" {
  id           = "pack-with-dest"
  group_id     = "default"
  description  = "Pack from source with destination"
  disabled     = false
  display_name = "Pack from source with destination"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_destination" "my_packdest" {
  id       = "test"
  pack  = criblio_pack.my_pack.id
  group_id = "default"

  output_s3 = {
    id                                = "test"
    type                              = "s3"
    system_fields                     = ["cribl_pipe"]
    streamtags                        = []
    aws_authentication_method         = "auto"
    signature_version                 = "v4"
    reuse_connections                 = true
    reject_unauthorized               = true
    enable_assume_role                = false
    duration_seconds                  = 3600
    stage_path                        = "$CRIBL_HOME/state/outputs/staging"
    add_id_to_stage_path              = true
    dest_path                         = ""
    object_acl                        = "private"
    remove_empty_dirs                 = true
    partition_expr                    = "C.Time.strftime(_time ? _time : Date.now()/1000, '%Y/%m/%d')"
    format                            = "json"
    base_file_name                    = "`CriblOut`"
    file_name_suffix                  = "test"
    max_file_size_mb                  = 32
    max_open_files                    = 100
    header_line                       = ""
    write_high_water_mark             = 64
    on_backpressure                   = "block"
    deadletter_enabled                = false
    on_disk_full_backpressure         = "block"
    max_file_open_time_sec            = 300
    max_file_idle_time_sec            = 30
    max_concurrent_file_parts         = 4
    verify_permissions                = true
    max_closing_files_to_backpressure = 100
    compress                          = "gzip"
    compression_level                 = "best_speed"
    compression_level                 = "best_speed"
    empty_dir_cleanup_sec             = 300
    bucket                            = "test"
  }
}


 