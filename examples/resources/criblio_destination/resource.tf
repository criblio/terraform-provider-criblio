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

variable "aws_secret_key" {
  type      = string
  sensitive = true
}
