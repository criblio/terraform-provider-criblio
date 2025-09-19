# Cribl Cloud Workspace Operations - Complete Deployment Example

## 🎯 Overview

This directory provides a comprehensive Terraform-based example for deploying a full Cribl Cloud workspace with integrated Stream processing, Edge fleet management, Lake storage, and Search capabilities. It demonstrates best practices for production-ready deployments in the cloud.

## 🚀 Quick Start

### Prerequisites

- Terraform >= 1.0
- Docker & Docker Compose (for Edge deployment)
- Cribl Cloud account with OAuth2 credentials
- AWS account (for S3 and other integrations)

### Using This Module

Create a new Terraform configuration that references this module:

```hcl
# main.tf
module "cribl_workspace" {
  # For GitHub repository via SSH
  source = "git::ssh://git@github.com/criblio/terraform-provider-criblio.git//modules/workspace-ops"
  
  # Or for GitHub via HTTPS
  # source = "git::https://github.com/criblio/terraform-provider-criblio.git//modules/workspace-ops"
  
  # Or for local development
  # source = "../terraform-provider-criblio/modules/workspace-ops"
  
  # Required variables
  cribl_client_id     = var.cribl_client_id
  cribl_client_secret = var.cribl_client_secret
  cloud_tenant        = "your-tenant-id"
  workspace           = "main"
  
  # Optional: Override defaults
  group-cloud         = "production-stream"
  group-hybrid        = "production-hybrid"
  edge_group          = "production-edge"
  cloud_region        = "us-east-1"
  instance_type       = "t3.large"
  edge_instance_count = 5
  edge_instance_type  = "t3.medium"
}

# variables.tf
variable "cribl_client_id" {
  description = "Cribl OAuth2 client ID"
  type        = string
  sensitive   = true
}

variable "cribl_client_secret" {
  description = "Cribl OAuth2 client secret"
  type        = string
  sensitive   = true
}
```

Deploy the infrastructure:

```bash
# Set credentials
export TF_VAR_cribl_client_id="your-client-id"
export TF_VAR_cribl_client_secret="your-client-secret"

# Initialize and deploy
terraform init
terraform plan
terraform apply


```

## 📊 Architecture Components

### Stream Processing (`stream.tf`)
- **Sources**: OTLP (port 20007), Syslog (port 20005), HTTP (port 20003)
- **Pipelines**: Data enrichment, CSV parsing for firewall logs, field extraction
- **Destinations**: Cribl Lake, OpenTelemetry, CrowdStrike SIEM
- **Automatic commit and deploy workflow**

### Edge Fleet Management (`edge.tf`)
- **3-node containerized Edge fleet**
- **Automatic bootstrap token generation**
- **Built-in source imports**: File Monitor, Journal, System Metrics
- **Custom sources**: Syslog, HTTP, TCP Metrics
- **Dual destinations**: S3 archive, Cribl HTTP forwarding

### Lake Storage (`lake.tf`)
- **Dataset creation with 30-day retention**
- **JSON format with gzip compression**
- **Lakehouse deployment (medium tier)**
- **Automatic dataset-lakehouse connection**

### Search Capabilities (`search.tf`)
- **S3 dataset provider with AssumeRole authentication**
- **Web logs dataset with time-based partitioning**
- **Saved queries with scheduled execution**
- **Pack integration for VPC Flow Logs analysis**

## 🔧 Configuration Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `cloud_tenant` | Cribl Cloud organization ID | `zamora-berries-wild8` |
| `workspace` | Cribl workspace name | `main` |
| `cribl_client_id` | OAuth2 client ID | Required (sensitive) |
| `cribl_client_secret` | OAuth2 client secret | Required (sensitive) |
| `group-cloud` | Stream worker group | `notdefault` |
| `group-hybrid` | Hybrid worker group | `notdefaulthybrid` |
| `edge_group` | Edge fleet group | `edge-fleet-tf2` |
| `edge_instance_count` | Edge fleet size | `3` |
| `edge_instance_type` | Instance type for Edge nodes | `t3.small` |
| `cribl_version` | Cribl software version | `4.12.0` |
| `cloud_region` | AWS region | `ca-central-1` |
| `instance_type` | EC2 instance type | `t3.medium` |
| `instance_count` | Number of instances | `2` |
| `environment` | Environment name | N/A |
| `estimated_ingest_rate` | Estimated data ingest rate | N/A |

## 📁 Module Structure

```
workspace-ops/
├── provider.tf              # Provider configuration
├── variables.tf             # Variable definitions
├── stream.tf               # Stream processing setup
├── edge.tf                 # Edge fleet configuration
├── lake.tf                 # Lake storage setup
├── search.tf               # Search configuration
├── README.md               # This documentation
├── agent.md                # Agent documentation
├── *.crbl                  # Pack files (Palo Alto, VPC Flow Logs)
└── edge-docker-compose.yml # Generated Docker config
```

## 🔄 Data Flow Architecture

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   Sources   │────▶│  Processing  │────▶│ Destinations │
├─────────────┤     ├──────────────┤     ├──────────────┤
│ • OTLP      │     │ • Pipelines  │     │ • Lake       │
│ • Syslog    │     │ • Filtering  │     │ • S3         │
│ • HTTP      │     │ • Enrichment │     │ • SIEM       │
│ • Metrics   │     │ • Parsing    │     │ • OTLP       │
└─────────────┘     └──────────────┘     └──────────────┘
```

## 🚢 Deployment Workflows

### Stream Deployment
1. Sources and pipelines are created
2. Destinations are configured
3. Connections are established
4. Configuration is committed
5. Automatic deployment to worker group

### Edge Fleet Deployment
1. Edge group created in Cribl Cloud
2. Bootstrap token generated via OAuth2
3. Docker Compose file generated with tokens
4. Pipelines and sources deployed
5. Built-in sources imported
6. Run deployment script to start containers

### Lake Setup
1. Dataset created with retention policy
2. Lakehouse provisioned (10-minute wait)
3. Dataset connected to lakehouse
4. Search queries configured

## 🔌 Integration Points

### Authentication Methods
- **OAuth2**: Client credentials flow for API access
- **AssumeRole**: AWS IAM role assumption for S3 access
- **Manual**: Token-based auth for specific destinations

### Supported Data Formats
- JSON (primary)
- CSV (Palo Alto firewall logs)
- OTLP (OpenTelemetry)
- Syslog (RFC3164/RFC5424)

### Network Ports
- **Stream Ports**: 20003 (HTTP), 20005 (Syslog), 20007 (OTLP)
- **Edge Ports**: 515 (Syslog UDP), 10515 (Syslog TCP), 8088 (HTTP), 8125 (Metrics)
- **UI Access**: 9001-9003 (Edge nodes)

## 📝 Usage Examples

### Send Data to Stream
```bash
# Send syslog data
echo "<14>Test message" | nc -u localhost 20005

# Send HTTP data
curl -X POST http://localhost:20003 \
  -H "Content-Type: application/json" \
  -d '{"event": "test", "source": "api"}'

# Send OTLP data
# Configure your OTLP exporter to point to localhost:20007
```

### Access Edge Nodes
```bash
# Node 1
open http://localhost:9001

# Node 2
open http://localhost:9002

# Node 3
open http://localhost:9003
```

### Monitor in Cribl Cloud
```bash
# Stream monitoring
open https://main-zamora-berries-wild8.cribl.cloud/?group=notdefault

# Edge monitoring
open https://main-zamora-berries-wild8.cribl.cloud/?group=edge-fleet-tf2
```

## 🛠️ Troubleshooting

### Common Issues

1. **Bootstrap Token Errors**
   - Verify OAuth2 credentials are correct
   - Check network connectivity to Cribl Cloud
   - Ensure proper workspace/tenant configuration

2. **Edge Nodes Not Connecting**
   - Check Docker is running
   - Verify bootstrap token is valid
   - Review container logs: `docker-compose logs`

3. **Destination Connection Failures**
   - Verify AWS credentials/roles
   - Check network security groups
   - Review destination endpoint URLs

### Debug Commands
```bash
# Check Terraform state
terraform state list

# View Edge container logs
docker-compose -f edge-docker-compose.yml logs -f

# Test connectivity
curl -k https://your-workspace-tenant.cribl.cloud:9000/api/v1/health

# Run the deployment script
./deploy-edge-fleet.sh
```

## 🔐 Security Considerations

- Store sensitive variables in environment variables or secure vaults
- Use IAM roles instead of static credentials where possible
- Enable TLS for all production deployments
- Regularly rotate OAuth2 credentials
- Configure proper network segmentation

## 📚 Additional Resources

- [Cribl Documentation](https://docs.cribl.io)
- [Terraform Provider Documentation](https://registry.terraform.io/providers/criblio/criblio)
- [Edge Deployment Guide](https://docs.cribl.io/edge)
- [Lake Configuration](https://docs.cribl.io/lake)

## 🤝 Contributing

To extend this example:

1. Add new sources in the `source-module`
2. Create custom pipelines for your data formats
3. Configure additional destinations
4. Extend the Edge fleet configuration

## 📄 License

This example is provided as-is for demonstration purposes. Modify as needed for your production deployments.

---

**Note**: Remember to configure manual connections in the Cribl Cloud UI for built-in Edge sources after deployment.
