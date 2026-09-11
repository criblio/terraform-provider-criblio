resource "criblio_destination" "my_destination" {
  group_id = "default"
  id       = "out-s3-main"

  output_s3 = {
    id              = "out-s3-main"
    type            = "s3"
    bucket          = "`my-cribl-bucket`"
    region          = "us-east-1"
    aws_api_key     = "AKIAIOSFODNN7EXAMPLE"
    aws_secret_key  = var.aws_secret_key
    dest_path       = "`logs/$${C.Time.strftime(_time, '%Y/%m/%d')}`"
    stage_path      = "$CRIBL_HOME/state/outputs/out-s3-main"
    compress        = "gzip"
    format          = "json"
    on_backpressure = "block"
    pipeline        = "passthru"
  }
}

resource "criblio_destination" "filesystem" {
  group_id = "default"
  id       = "out-filesystem-main"

  output_filesystem = {
    id               = "out-filesystem-main"
    type             = "filesystem"
    dest_path        = "/var/log/cribl/out"
    stage_path       = "/var/log/cribl/stage"
    partition_expr   = "C.Time.strftime(_time ? _time : Date.now()/1000, '%Y/%m/%d')"
    base_file_name   = "`CriblOut`"
    file_name_suffix = "`.json.gz`"
    format           = "json"
    compress         = "gzip"
  }
}

variable "aws_secret_key" {
  type      = string
  sensitive = true
}
