# Create Worker Group
resource "criblio_group" "example" {
  id                    = var.worker_group_id
  name                  = var.worker_group_id
  product               = "stream"
  on_prem               = false
  is_fleet              = false
  worker_remote_access  = true
  estimated_ingest_rate = 2048 # Equivalent to 24 MB/s with 9 Worker Processes
  provisioned           = false
  cloud = {
    provider = "aws"
    region   = "us-east-1" # Replace with a different AWS region if desired
  }
}

# Create Syslog Source
resource "criblio_source" "syslog" {
  id       = "in-syslog-9021"
  group_id = criblio_group.example.id

  input_syslog = {
    id             = "in-syslog-9021"
    type           = "syslog"
    host           = "0.0.0.0"
    tcp_port       = 9021 # Replace with a different port number if desired
    disabled       = false
    send_to_routes = true
    tls = {
      disabled = true
    }
  }

  depends_on = [criblio_group.example]
}

# Create S3 Destination
resource "criblio_destination" "s3" {
  id       = "out_s3"
  group_id = criblio_group.example.id

  output_s3 = {
    id                    = "out_s3"
    type                  = "s3"
    bucket                = var.aws_bucket_name
    region                = var.aws_region
    aws_api_key           = var.aws_api_key
    aws_secret_key        = var.aws_secret_key
    stage_path            = "/tmp/cribl_stage"
    compress              = "gzip"
    compression_level     = "best_speed"
    empty_dir_cleanup_sec = 300
  }

  depends_on = [criblio_group.example]
}

# Create Pipeline
resource "criblio_pipeline" "filter" {
  id       = "my_pipeline"
  group_id = criblio_group.example.id

  conf = {
    async_func_timeout = 1000
    functions = [
      {
        id       = "eval"
        filter   = "true"
        disabled = false
        final    = true
        conf = jsonencode({
          remove = ["*"]
          keep   = ["eventSource", "eventID"]
        })
      }
    ]
  }

  depends_on = [criblio_group.example]
}

# Create Routes
resource "criblio_routes" "main" {
  group_id = criblio_group.example.id
  id       = "default" # Routing table ID (do not change; the supported value is default)

  routes = [
    {
      name        = "your_route"              # Replace with the name of the Route
      description = "This is my new Route"    # Replace with the desired Route description
      pipeline    = criblio_pipeline.filter.id
      output      = criblio_destination.s3.id
      filter      = "__inputId=='in-syslog-9021'"
      final       = true
      disabled    = false
    },
    {
      name     = "default"
      pipeline = "main"
      output   = "default"
      filter   = "true"
      final    = false
      disabled = false
    }
  ]

  depends_on = [
    criblio_source.syslog,
    criblio_destination.s3,
    criblio_pipeline.filter,
  ]
}

# Commit configuration
resource "criblio_commit" "example" {
  effective = true
  group     = criblio_group.example.id
  message   = "Commit for Cribl Stream example"

  depends_on = [criblio_routes.main]
}

# Read config version
data "criblio_config_version" "latest" {
  id         = criblio_group.example.id
  depends_on = [criblio_commit.example]
}

# Deploy configuration
resource "criblio_deploy" "example" {
  id      = criblio_group.example.id
  version = data.criblio_config_version.latest.items[0]
}
