### Sources ###


module "otlp_source" {
  source = "../source-module"
  
  source_id   = "otel-collector"
  group_id    = var.group-cloud
  source_type = "otlp"
  port        = 20007
  description = "OpenTelemetry receiver"
  
  # Optional: Override defaults
  custom_config = {
    extractSpans   = true   # camelCase!
    extractMetrics = true
    extractLogs    = false
    protocol       = "grpc"
  }
}


module "my_syslog" {
  source = "../source-module"
  
  source_id   = "firewall-logstf"
  group_id    = var.group-cloud
  source_type = "syslog"
  port        = 20005  
  connections = [
    {
      output   = module.cribl_lake_dest.destination_id
      pipeline = ""
    },
    {
      output   = module.crowdstrike_siem_dest.destination_id
      pipeline = criblio_pipeline.my_pipeline.id
    }
  ]
    depends_on = [criblio_destination.cribl_lake, criblio_pipeline.my_pipeline, module.crowdstrike_siem_dest]
}

module "http_source" {
  source = "../source-module" 
  source_id   = "http-input"
  group_id    = var.group-cloud
  source_type = "cribl_http"
  port        = 20003
  description = "HTTP input for receiving events"
  
  # Connections
  connections = [
    {
      output   = criblio_destination.cribl_lake.id
      pipeline = criblio_pipeline.my_pipeline.id
    }
  ]
  
  # Pipeline override
  pipeline = criblio_pipeline.my_pipeline.id
  
  # Tags
  streamtags = ["http", "api"]
  
  # Persistent queue disabled
  pq_enabled = false
  
  # Custom configuration for all other settings
  custom_config = {
    # Auth
    auth_tokens = [""]  # Empty token
    
    # Network settings
    capture_headers         = true
    enable_health_check     = true
    enable_proxy_header     = false
    ip_allowlist_regex      = "/.*/"
    ip_denylist_regex       = "/^$/"
    keep_alive_timeout      = 5
    max_active_req          = 1000
    max_requests_per_socket = 0
    request_timeout         = 30
    socket_timeout          = 60
    
    # Environment
    environment = "production"
    
    # Metadata
    metadata = [
      {
        name  = "source_type"
        value = "http"
      },
      {
        name  = "data_center"
        value = "us-west-2"
      }
    ]
    
    # TLS override (disable it)
    tls = {
      disabled = true
    }
  }
  
  depends_on = [criblio_destination.cribl_lake, criblio_pipeline.my_pipeline]
}

### Pipelines ###
#################
resource "criblio_pipeline" "my_pipeline" {
  id = "my-first-tf-pipeline2"
  group_id = var.group-cloud
  conf = {
    description = "Test pipeline for Terraform provider"
    functions = [
      {
        id          = "eval"
        description = "Set organizationId test11"
        disabled    = false
        filter      = "true"
        final       = false
        conf = {
          add = jsonencode([
            {
              disabled = false
              name     = "daaa"
              value    = "ssaadfas"
            }
          ])
          keep = jsonencode([
            "aedgaa"
          ])
          remove = jsonencode([
            "adfgsdfgaaccsd"
          ])
        }
      },
{
  id     = "serde"
  filter = "true"
  conf = {
    mode     = jsonencode("extract")
    type     = jsonencode("csv")
    srcField = jsonencode("_raw")
    keep = jsonencode([
      "*"
    ])
    remove = jsonencode([
      "_raw"
    ])
    fields = jsonencode([
      "future_use_0",
      "receive_time",
      "serial_number",
      "type",
      "threat_content_type",
      "future_use_1",
      "generated_time",
      "source_ip",
      "destination_ip",
      "nat_source_ip",
      "nat_destination_ip",
      "rule_name",
      "source_user",
      "destination_user",
      "application",
      "virtual_system",
      "source_zone",
      "destination_zone",
      "inbound_interface",
      "outbound_interface",
      "log_action",
      "future_use_2",
      "session_id",
      "repeat_count",
      "source_port",
      "destination_port",
      "nat_source_port",
      "nat_destination_port",
      "flags",
      "protocol",
      "action",
      "bytes",
      "bytes_sent",
      "bytes_received",
      "packets",
      "start_time",
      "elapsed_time",
      "category",
      "future_use_3",
      "sequence_number",
      "action_flags",
      "source_location",
      "destination_location",
      "future_use_4",
      "packets_sent",
      "packets_received",
      "session_end_reason",
      "device_group_hierarchy_level_1",
      "device_group_hierarchy_level_2",
      "device_group_hierarchy_level_3",
      "device_group_hierarchy_level_4",
      "virtual_system_name",
      "device_name",
      "action_source",
      "source_vm_uuid",
      "destination_vm_uuid",
      "tunnel_id_imsi",
      "monitor_tag_imei",
      "parent_session_id",
      "parent_start_time",
      "tunnel_type",
      "sctp_association_id",
      "sctp_chunks",
      "sctp_chunks_sent",
      "sctp_chunks_received"
    ])
  }
}
            

    ]
    output = "default"
    streamtags = [""]
  }
  # depends_on = [criblio_group.workergroup]
}
####################
### Destinations ###
####################

# Cribl Lake Destination
module "cribl_lake_dest" {
  source = "../destination-module"
  
  destination_id   = "${criblio_cribl_lake_dataset.my_cribllakedataset.id}-dest"
  group_id         = var.group-cloud
  destination_type = "cribl_lake"
  description      = "Cribl Lake destination"
  
  # Cribl Lake specific
  dest_path = criblio_cribl_lake_dataset.my_cribllakedataset.id
  format    = "json"
  compress  = "gzip"
  
  # Tags
  streamtags = ["test", "lake"]
  
  # Custom configuration
  custom_config = {
    max_file_size_mb = 64
    base_file_name   = "TestOut"
  }
  depends_on = [criblio_pipeline.my_pipeline, criblio_cribl_lake_dataset.my_cribllakedataset]
}



# Open Telemetry Destination
module "open_telemetry_dest" {
  source = "../destination-module"
  
  destination_id   = "open-telemetry-tf"
  group_id         = var.group-cloud
  destination_type = "open_telemetry"
  description      = "Terraform OpenTelemetry destination"
  
  # OpenTelemetry specific
  url      = "https://otlp-collector.example.com:4318/v1/traces"
  format   = "json"
  compress = "gzip"
  
  # Tags
  streamtags = ["test", "opentelemetry", "otlp"]
  
  # Custom configuration
  custom_config = {
    protocol                    = "http"
    otlp_version               = "1.3.1"
    http_traces_endpoint_override = "/v1/traces"
    http_metrics_endpoint_override = "/v1/metrics"
    http_logs_endpoint_override = "/v1/logs"
    metadata = [
      {
        key   = "service.name"
        value = "cribl-otlp"
      },
      {
        key   = "environment"
        value = "terraform"
      }
    ]
    concurrency         = 10
    timeout_sec         = 30
    keep_alive          = true
    keep_alive_time     = 60
  }
}

# CrowdStrike Next-Gen SIEM Destination
module "crowdstrike_siem_dest" {
  source = "../destination-module"
  
  destination_id   = "test-crowdstrike-siem"
  group_id         = var.group-cloud
  destination_type = "crowdstrike_next_gen_siem"
  description      = "Test CrowdStrike Next-Gen SIEM destination"
  
  # CrowdStrike specific
  url      = "https://api.crowdstrike.com/siem/v1/events"
  token    = "CS-SWIFTIEST-TOKEN-12345678-1234-1234-1234-123456789abc"  # Use a variable in real usage
  format   = "JSON"
  compress = "gzip"
  
  # Tags
  streamtags = ["test", "crowdstrike", "siem"]
  
  # Custom configuration
  custom_config = {
    auth_type           = "manual"
    concurrency         = 5
    flush_period_sec    = 2
    max_payload_events  = 100
    max_payload_size_kb = 2048
    on_backpressure     = "queue"
    extra_http_headers = [
      {
        name  = "X-CS-Source"
        value = "cribl-test"
      }
    ]
    timeout_sec         = 45
    use_round_robin_dns = true
  }
}

################################
### Commit and Deploy ###
################################


resource "criblio_commit" "my_pipecommit" {
  effective = true
  group     = var.group-cloud
  message   = "terraform commit pipeline ${timestamp()}"
  depends_on = [
    module.otlp_source,
    module.my_syslog,
    module.http_source,
    criblio_pipeline.my_pipeline
  ]
}

# Step 2: Get config version after commit
data "criblio_config_version" "my_pipelineconfigversion" {
  id = var.group-cloud
  depends_on = [criblio_commit.my_pipecommit]
}

# Step 3: Deploy using the retrieved version
resource "criblio_deploy" "my_pipedeploy_safe" {
  id = var.group-cloud
  version = length(data.criblio_config_version.my_pipelineconfigversion.items) > 0 ? data.criblio_config_version.my_pipelineconfigversion.items[length(data.criblio_config_version.my_pipelineconfigversion.items) - 1] : "default"
  depends_on = [data.criblio_config_version.my_pipelineconfigversion]
}
