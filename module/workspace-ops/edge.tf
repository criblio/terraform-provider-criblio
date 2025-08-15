################################
### Edge Fleet Configuration ###
################################

# Edge-specific variables
variable "edge_group" {
  description = "Edge worker group name"
  type        = string
  default     = "edge-fleet-tf2"
}

variable "edge_instance_count" {
  description = "Number of Edge instances"
  type        = number
  default     = 3
}

variable "edge_instance_type" {
  description = "Instance type for Edge nodes"
  type        = string
  default     = "t3.small"
}

################################
### Edge Worker Group ###
################################

resource "criblio_group" "edge_worker_group" {
  id                    = var.edge_group
  name                  = var.edge_group
  product               = "edge"
  description           = "Terraform-managed Edge fleet worker group"
  is_fleet              = true
  on_prem               = true
  provisioned           = true
  worker_remote_access  = true
  estimated_ingest_rate = var.estimated_ingest_rate
  
  streamtags = ["edge", "terraform", "distributed"]
}

# Start Docker containers after group creation
resource "null_resource" "start_edge_docker" {
  depends_on = [local_file.edge_docker_compose]
  
  provisioner "local-exec" {
    command = <<-EOT
      echo "Starting Edge Docker containers..."
      docker-compose -f edge-docker-compose.yml down 2>/dev/null || true
      docker-compose -f edge-docker-compose.yml up -d
      echo "Waiting for Edge nodes to connect (60 seconds)..."
      sleep 60
      echo "Edge nodes should now be connected"
    EOT
  }
  
  triggers = {
    docker_compose = local_file.edge_docker_compose.content
  }
}

# Add delay resource to ensure nodes are connected
resource "time_sleep" "wait_for_edge_nodes" {
  depends_on = [null_resource.start_edge_docker]
  
  create_duration = "30s"
}

################################
### Edge Pipelines ###
################################

# Main Edge processing pipeline
resource "criblio_pipeline" "edge_main_pipeline" {
  id       = "edge-main-processing"
  group_id = var.edge_group
  
  conf = {
    description = "Main Edge fleet data processing pipeline"
    functions = [
      {
        id          = "eval"
        description = "Add Edge metadata"
        filter      = "true"
        conf = {
          add = jsonencode([
            {
              name  = "edge_fleet"
              value = "'${var.edge_group}'"
            },
            {
              name  = "processing_time"
              value = "Date.now()"
            },
            {
              name  = "edge_location"
              value = "'${var.cloud_region}'"
            },
            {
              name  = "edge_node"
              value = "C.os.hostname()"
            }
          ])
        }
      }
    ]
    # output = "devnull"
    # streamtags = ["edge", "processed"]
  }
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# # Syslog processing pipeline
# resource "criblio_pipeline" "edge_syslog_pipeline" {
#   id       = "edge-syslog-processing"
#   group_id = var.edge_group
  
#   conf = {
#     description = "Edge syslog data processing"
#     functions = [
#       {
#         id          = "eval"
#         description = "Set organizationId test11"
#         disabled    = false
#         filter      = "true"
#         final       = false
#         conf = {
#           add = jsonencode([
#             {
#               disabled = false
#               name     = "daaa"
#               value    = "ssaadfas"
#             }
#           ])
#           keep = jsonencode([
#             "aedgaa"
#           ])
#           remove = jsonencode([
#             "adfgsdfgaaccsd"
#           ])
#         }
#       }
#     ]
#     output = "devnull"
#   }
  
#   depends_on = [time_sleep.wait_for_edge_nodes]
# }

################################
### Edge Sources ###
################################
# Use data sources to reference the automatically created Edge sources

# File Monitor source (automatically created by Edge)
data "criblio_source" "edge_file_monitor" {
  id       = "in_file_varlog"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# Journal Files source (automatically created by Edge)
data "criblio_source" "edge_journal" {
  id       = "in_journal_local"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# System Metrics source (automatically created by Edge)
data "criblio_source" "edge_system_metrics" {
  id       = "in_system_metrics"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# System State source (automatically created by Edge)
data "criblio_source" "edge_system_state" {
  id       = "in_system_state"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# Kubernetes Events source (automatically created by Edge)
data "criblio_source" "edge_kube_events" {
  id       = "in_kube_events"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# Kubernetes Logs source (automatically created by Edge)
data "criblio_source" "edge_kube_logs" {
  id       = "in_kube_logs"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# Kubernetes Metrics source (automatically created by Edge)
data "criblio_source" "edge_kube_metrics" {
  id       = "in_kube_metrics"
  group_id = var.edge_group
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# Syslog source
module "edge_syslog_source" {
  source = "../source-module"
  
  source_id   = "edge-syslog"
  group_id    = var.edge_group
  source_type = "syslog"
  port        = 515  # Non-default port as required by module
  description = "Edge fleet syslog receiver"
  
  custom_config = {
    udpPort       = 515
    tcpPort       = 10515
    mode          = "manual"
    facility      = 16
    timezone      = "UTC"
    enableProxyHeader = false
    maxActiveConns = 1000

    pipelines = ["edge-syslog-processing"]
  }
  connections = [
    {
      output   = module.edge_s3_dest.destination_id
      pipeline = criblio_pipeline.edge_main_pipeline.id
    },
    {
      output   = module.edge_cribl_http_dest.destination_id
      pipeline = criblio_pipeline.edge_main_pipeline.id
    }
  ]
  streamtags = ["syslog", "edge"]
  
  # depends_on = [criblio_pipeline.edge_syslog_pipeline]
}

# HTTP source for metrics and events
module "edge_http_source" {
  source = "../source-module"
  
  source_id   = "edge-http"
  group_id    = var.edge_group
  source_type = "http"
  port        = 8088  # Non-default port
  description = "Edge HTTP receiver for metrics and events"
  
  custom_config = {
    tls = {
      disabled = true
    }
    auth = {
      type = "manual"
      users = jsonencode([
        {
          username = "edge"
          password = "cribledge123"
        }
      ])
    }
    pipelines = ["edge-main-processing"]
  }
  
  connections = [
    {
      output   = module.edge_cribl_http_dest.destination_id
      pipeline = criblio_pipeline.edge_main_pipeline.id
    }
  ]
  streamtags = ["http", "edge"]
  
  depends_on = [criblio_pipeline.edge_main_pipeline]
}

# TCP source for metrics
module "edge_tcp_source" {
  source = "../source-module"
  
  source_id   = "edge-metrics"
  group_id    = var.edge_group
  source_type = "tcp"
  port        = 8125  # StatsD port
  description = "Edge metrics receiver"
  
  custom_config = {
    mode       = "server"
    protocol   = "tcp"
    pipelines  = ["edge-main-processing"]
  }
  
  connections = [
    {
      output   = module.edge_cribl_http_dest.destination_id
      pipeline = criblio_pipeline.edge_main_pipeline.id
    }
  ]
  streamtags = ["metrics", "edge"]
  
  depends_on = [criblio_pipeline.edge_main_pipeline]
}

################################
### Edge Destinations ###
################################

# S3 destination for archival
module "edge_s3_dest" {
  source = "../destination-module"
  
  destination_id   = "edge-s3-archive"
  group_id         = var.edge_group
  destination_type = "s3"
  description      = "Edge data S3 archive"
  
  dest_path = "edge-fleet"
  format    = "json"
  compress  = "gzip"
  
  custom_config = {
    bucket              = "cribl-edge-data-${var.cloud_tenant}"
    region              = var.cloud_region
    destPath            = "'edge-fleet/' + C.Time.strftime(_time, '%Y/%m/%d/%H')"
    objectACL           = "bucket-owner-full-control"
    storageClass        = "STANDARD_IA"
    serverSideEncryption = "aws:kms"
    partitionExpr       = "`$${edge_fleet}/$${C.Time.strftime(_time, '%Y/%m/%d')}`"
    maxFileSizeMB       = 128
    maxOpenFiles        = 10
    maxFileIdleTime     = "30s"
    onBackpressure      = "block"
  }
  
  streamtags = ["archive", "s3"]
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

# Cribl HTTP destination
module "edge_cribl_http_dest" {
  source = "../destination-module"
  
  destination_id   = "edge-cribl-http"
  group_id         = var.edge_group
  destination_type = "cribl_http"
  description      = "Edge Cribl HTTP forwarder"
  
  custom_config = {
    host               = "https://${var.group-cloud}.${var.workspace}-${var.cloud_tenant}.cribl.cloud:10200"  # Adjust to your Cribl Stream host
    port               = 10200
    authType           = "none"
    compress           = true
    rejectUnauthorized = false
    onBackpressure     = "queue"
    maxPayloadSizeKB   = 4096
    maxPayloadEvents   = 1000
    flushPeriodSec     = 1
    url                = "https://${var.workspace}-${var.cloud_tenant}.cribl.cloud:10200"
  }
  
  streamtags = ["cribl", "edge"]
  
  depends_on = [time_sleep.wait_for_edge_nodes]
}

################################
### Import Existing Edge Sources ###
################################

# Import blocks for system sources (commented out as placeholders)
# These sources are automatically created by Edge and cannot be managed via Terraform

# import {
#   to = criblio_source.file_monitor_managed
#   id = "{\"group_id\":\"${var.edge_group}\",\"id\":\"in_file_varlog\"}"
# }

# import {
#   to = criblio_source.journal_managed
#   id = "{\"group_id\":\"${var.edge_group}\",\"id\":\"in_journal_local\"}"
# }

# import {
#   to = criblio_source.system_metrics_managed
#   id = "{\"group_id\":\"${var.edge_group}\",\"id\":\"in_system_metrics\"}"
# }

# import {
#   to = criblio_source.system_state_managed
#   id = "{\"group_id\":\"${var.edge_group}\",\"id\":\"in_system_state\"}"
# }

# Note: System sources are automatically created by Edge and cannot be managed via Terraform
# The import blocks above are placeholders for future use when/if the provider supports it
# For now, use the UI to manually connect these sources to destinations

################################
### Edge Commit and Deploy ###
################################

# Commit all Edge configurations
resource "criblio_commit" "edge_group_commit" {
  effective = true
  group     = var.edge_group
  message   = "Edge group complete TF setup with pipelines - ${timestamp()}"
  
  depends_on = [
    criblio_group.edge_worker_group,
    criblio_pipeline.edge_main_pipeline,
    module.edge_syslog_source,
    module.edge_s3_dest,
    module.edge_cribl_http_dest
  ]
}

# Deploy happens automatically with effective commit
# No need for separate deploy resource

################################
### Edge Bootstrap Token ###
################################

# Use the bootstrap token module to get the token
module "edge_bootstrap_token" {
  source = "../cribl-bootstrap-token-module"
  
  client_id     = var.cribl_client_id
  client_secret = var.cribl_client_secret
  organization  = var.cloud_tenant
  workspace     = var.workspace
  workspace_url = "https://${var.workspace}-${var.cloud_tenant}.cribl.cloud"
  environment   = var.environment
  group         = var.edge_group
  
  # Auto-execute to get the token
  auto_execute       = true
  create_script_file = false
  token_output_file  = "${path.module}/.edge-bootstrap-token"
  
  depends_on = [criblio_group.edge_worker_group]
}

locals {
  edge_bootstrap_token = module.edge_bootstrap_token.bootstrap_token
  edge_master_url = "tls://${local.edge_bootstrap_token}@${var.workspace}-${var.cloud_tenant}.cribl.cloud:4200?group=${var.edge_group}&tls.rejectUnauthorized=false"
}

################################
### Docker Deployment Files ###
################################

# Generate Docker Compose file for Edge deployment
resource "local_file" "edge_docker_compose" {
  filename = "${path.module}/edge-docker-compose.yml"
  content  = <<-EOT
services:
%{ for i in range(1, var.edge_instance_count + 1) ~}
  edge-node-${i}:
    image: cribl/cribl:${var.cribl_version}
    container_name: edge-fleet-tf-node-${i}
    restart: unless-stopped
    environment:
      - CRIBL_DIST_MODE=managed-edge
      - CRIBL_DIST_MASTER_URL=${local.edge_master_url}
      - CRIBL_EDGE=1
    ports:
      - "${9000 + i}:9000"           # UI port
      - "${515 + i}:515/udp"         # Syslog UDP
      - "${10515 + i}:10515/tcp"     # Syslog TCP
      - "${8088 + i}:8088"           # HTTP input
      - "${8125 + i}:8125/tcp"       # TCP Metrics
    volumes:
      - edge-node-${i}-config:/opt/cribl/config
      - edge-node-${i}-data:/opt/cribl/data
    networks:
      - cribl-edge
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

%{ endfor ~}
networks:
  cribl-edge:
    driver: bridge

volumes:
%{ for i in range(1, var.edge_instance_count + 1) ~}
  edge-node-${i}-config:
  edge-node-${i}-data:
%{ endfor ~}
EOT
  
  depends_on = [module.edge_bootstrap_token]
}

# Generate deployment script
resource "local_file" "edge_deploy_script" {
  filename        = "${path.module}/deploy-edge-fleet.sh"
  file_permission = "0755"
  content         = <<-EOT
#!/bin/bash
# Cribl Edge Fleet Deployment Script
# Generated by Terraform

set -e

echo "========================================"
echo "Deploying Cribl Edge Fleet"
echo "========================================"
echo "Edge Group: ${var.edge_group}"
echo "Cloud Tenant: ${var.cloud_tenant}"
echo "Number of nodes: ${var.edge_instance_count}"
echo ""

# Stop and remove existing containers if any
echo "Cleaning up existing Edge containers..."
docker-compose -f edge-docker-compose.yml down 2>/dev/null || true

# Deploy Edge fleet
echo "Starting Edge fleet..."
docker-compose -f edge-docker-compose.yml up -d

# Wait for containers to start
echo "Waiting for Edge nodes to start..."
sleep 10

# Check container status
echo ""
echo "Edge Fleet Status:"
docker-compose -f edge-docker-compose.yml ps

echo ""
echo "Edge Fleet deployed successfully!"
echo "Access the Edge nodes at:"
%{ for i in range(1, var.edge_instance_count + 1) ~}
echo "  Node ${i}: http://localhost:${9000 + i}"
%{ endfor ~}
echo ""
echo "Monitor in Cribl Cloud:"
echo "  https://${var.workspace}-${var.cloud_tenant}.cribl.cloud/?group=${var.edge_group}"
EOT
  
  depends_on = [local_file.edge_docker_compose]
}

################################
### Helper Resources ###
################################

resource "random_string" "edge_token" {
  length  = 32
  special = false
}

################################
### Outputs ###
################################

output "edge_group_name" {
  value       = criblio_group.edge_worker_group.id
  description = "Edge worker group name"
}

output "edge_bootstrap_token" {
  value       = local.edge_bootstrap_token
  sensitive   = true
  description = "Bootstrap token for Edge nodes"
}

output "edge_master_url" {
  value       = local.edge_master_url
  sensitive   = true
  description = "Complete master URL for Edge nodes"
}

# Outputs for existing Edge sources - just the basic info
output "edge_sources_info" {
  value = {
    file_monitor = {
      id       = data.criblio_source.edge_file_monitor.id
      group_id = data.criblio_source.edge_file_monitor.group_id
    }
    journal = {
      id       = data.criblio_source.edge_journal.id
      group_id = data.criblio_source.edge_journal.group_id
    }
    system_metrics = {
      id       = data.criblio_source.edge_system_metrics.id
      group_id = data.criblio_source.edge_system_metrics.group_id
    }
    system_state = {
      id       = data.criblio_source.edge_system_state.id
      group_id = data.criblio_source.edge_system_state.group_id
    }
  }
  description = "Basic info for Edge data sources"
}

# Output all available attributes for the File Monitor source
output "edge_file_monitor_all_attributes" {
  value       = data.criblio_source.edge_file_monitor
  description = "All available attributes for the File Monitor source"
}

# Output all available attributes for the System Metrics source
output "edge_system_metrics_all_attributes" {
  value       = data.criblio_source.edge_system_metrics
  description = "All available attributes for the System Metrics source"
}

# Output all available attributes for the Journal source
output "edge_journal_all_attributes" {
  value       = data.criblio_source.edge_journal
  description = "All available attributes for the Journal source"
}

# Output all available attributes for the System State source
output "edge_system_state_all_attributes" {
  value       = data.criblio_source.edge_system_state
  description = "All available attributes for the System State source"
}




# Output all Edge sources as a map for easier reference
output "edge_sources_map" {
  value = {
    file_monitor    = data.criblio_source.edge_file_monitor.id
    journal         = data.criblio_source.edge_journal.id
    system_metrics  = data.criblio_source.edge_system_metrics.id
    system_state    = data.criblio_source.edge_system_state.id
    # kube_events     = data.criblio_source.edge_kube_events.id
    # kube_logs       = data.criblio_source.edge_kube_logs.id
    # kube_metrics    = data.criblio_source.edge_kube_metrics.id
    syslog_custom   = module.edge_syslog_source.source_id
    http_custom     = module.edge_http_source.source_id
    tcp_custom      = module.edge_tcp_source.source_id
  }
  description = "Map of all Edge source IDs for easy reference"
}

# Output destinations for reference
output "edge_destinations_map" {
  value = {
    s3_archive   = module.edge_s3_dest.destination_id
    cribl_http   = module.edge_cribl_http_dest.destination_id
  }
  description = "Map of all Edge destination IDs for easy reference"
}

output "edge_deployment_instructions" {
  value = <<-EOT
    ========================================
    Edge Fleet Setup Complete!
    ========================================
    
    Edge Group: ${var.edge_group}
    Cloud Tenant: ${var.cloud_tenant}
    Workspace URL: https://${var.workspace}-${var.cloud_tenant}.cribl.cloud
    
    DEPLOYMENT:
    
    1. Deploy Edge nodes using the generated Docker Compose:
       ./deploy-edge-fleet.sh
    
    Or manually:
       docker-compose -f edge-docker-compose.yml up -d
    
    2. Monitor your Edge fleet:
       https://${var.workspace}-${var.cloud_tenant}.cribl.cloud/?group=${var.edge_group}
    
    3. Edge nodes will automatically connect and receive configuration
    
    CONFIGURED COMPONENTS:
    - Pipelines: edge-main-processing, edge-syslog-processing
    - Sources: syslog (515/udp, 10515/tcp), http-in (8088), metrics-in (8125)
    - Destinations: s3-archive, cribl-http (10200)
    
    PORT MAPPINGS (per node):
    - Node 1: UI:9001, Syslog:516/10516, HTTP:8089, TCP:8126
    - Node 2: UI:9002, Syslog:517/10517, HTTP:8090, TCP:8127
    - Node 3: UI:9003, Syslog:518/10518, HTTP:8091, TCP:8128
  EOT
  description = "Instructions for deploying the Edge fleet"
}