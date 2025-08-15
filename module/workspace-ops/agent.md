# Cribl Workspace Operations - Agent Documentation

## Overview
This directory contains Terraform configurations for managing Cribl Cloud workspace operations, including Stream processing, Edge fleet management, Lake storage, and various data sources and destinations.

## Project Goals
- Create a comprehensive Terraform module for Cribl Cloud workspace management
- Demonstrate integration between Stream, Edge, and Lake components
- Provide reusable modules for sources, destinations, and pipelines
- Enable easy deployment of Edge fleets with proper bootstrap token management

## Current Directory Structure

### Core Configuration Files
- **provider.tf** - Cribl Cloud provider configuration with OAuth2 authentication
- **variables.tf** - Global variables including tenant, workspace, credentials, and infrastructure settings
- **stream.tf** - Stream worker group configuration with sources, pipelines, and destinations
- **edge.tf** - Edge fleet configuration (currently experiencing dependency cycles)
- **lake.tf** - Cribl Lake dataset and search configurations
- **search.tf** - Search configurations for data analysis
- **pack.tf** - Cribl pack installations

### Data Files
- **cribl-palo-alto-networks-source-1.0.0.crbl** - Palo Alto Networks source pack
- **cribl-search-aws-vpc-flow-logs_0.1.1.crbl** - AWS VPC Flow Logs search pack

### Generated Files
- **edge_get_bootstrap_token.sh** - Generated script for retrieving Edge bootstrap tokens
- **templates/edge-docker-compose.yml.tpl** - Docker Compose template for Edge deployment

### State Files
- **terraform.tfstate** - Current Terraform state
- **terraform.tfstate.backup** - Backup of previous state

## Current Implementation Status

### ✅ Working Components
1. **Stream Configuration** (`stream.tf`)
   - OpenTelemetry (OTLP) source on port 20007
   - Syslog source on port 20005 with firewall logs processing
   - HTTP source on port 20003 for API ingestion
   - Data processing pipeline with eval, serde (CSV parsing for Palo Alto logs)
   - Multiple destinations: Cribl Lake, OpenTelemetry, CrowdStrike SIEM
   - Proper commit and deploy workflow

2. **Lake Configuration** (`lake.tf`)
   - Dataset creation for log storage
   - Search configurations for data analysis

3. **Modular Architecture**
   - Separate modules for sources (`../source-module`)
   - Separate modules for destinations (`../destination-module`)
   - Reusable bootstrap token module (`../cribl-bootstrap-token-module`)

### ✅ Completed: Edge Fleet Configuration (`edge.tf`)

#### Successfully Resolved All Dependencies
The Edge configuration is now fully operational with all dependency cycles resolved through proper staging and import blocks.

#### Edge Architecture Implementation
1. **Edge Worker Group**
   - Product: edge
   - Fleet configuration: `is_fleet = true`
   - On-premises deployment: `on_prem = true`
   - Remote access enabled: `worker_remote_access = true`
   - Group ID: `edge-fleet-tf2`

2. **Bootstrap Token Management**
   - Fully integrated `../cribl-bootstrap-token-module`
   - Automatic token retrieval with OAuth2 authentication
   - Bootstrap token securely embedded in Docker Compose files
   - Master URL: `tls://TOKEN@WORKSPACE-TENANT.cribl.cloud:4200?group=EDGE_GROUP`

3. **Complete Edge Configuration**
   - **Pipeline**: `edge-main-processing` with metadata enrichment
   - **Sources**: 
     - Imported existing Edge sources: File Monitor, Journal Files, System Metrics, System State, Kubernetes sources
     - Custom sources: Syslog (515/10515), HTTP (8088), TCP Metrics (8125)
   - **Destinations**: S3 archive and Cribl HTTP (port 10200)
   - **Docker Deployment**: Complete 3-node containerized fleet with proper networking

#### Edge Deployment Workflow (Implemented)
1. ✅ Create Edge group in Cribl Cloud
2. ✅ Generate bootstrap token via OAuth2 module
3. ✅ Create Docker Compose configuration with bootstrap token
4. ✅ Deploy pipelines and custom sources/destinations
5. ✅ Import existing Edge sources into Terraform management
6. ✅ Generate deployment scripts and Docker configuration
7. ✅ Commit and deploy all configurations

#### Import Strategy for Built-in Sources
Successfully implemented import blocks for existing Edge sources:
- File Monitor (`in_file_varlog`)
- Journal Files (`in_journal_local`) 
- System Metrics (`in_system_metrics`)
- System State (`in_system_state`)

**Note**: Built-in Edge sources don't support Terraform-managed connections. Connections must be configured manually in the UI by dragging from sources to destinations.

## Integration Points

### Authentication
- OAuth2 client credentials flow
- Client ID and secret stored in variables
- Environment-specific endpoints (production/staging)

### Data Flow
```
Data Sources → Stream/Edge Processing → Pipelines → Destinations
```

### Destinations
- **Cribl Lake**: Long-term storage and search
- **S3**: Archive storage with partitioning
- **OpenTelemetry**: Observability platforms
- **CrowdStrike SIEM**: Security monitoring
- **Prometheus**: Metrics collection (Edge)

## Variables Configuration
Key variables defined in `variables.tf`:
- `cloud_tenant`: Cribl Cloud organization ID
- `workspace`: Cribl workspace name
- `cribl_client_id`/`cribl_client_secret`: OAuth credentials
- `group-hybrid`: Worker group for hybrid deployment
- `group-cloud`: Stream group (default)
- Infrastructure: instance types, counts, regions
- Cribl version: Currently 4.12.0

## Dependencies
- **External Modules**:
  - `../source-module`: Reusable source configurations
  - `../destination-module`: Reusable destination configurations (now includes Prometheus)
  - `../cribl-bootstrap-token-module`: Bootstrap token retrieval

## Recent Achievements ✅

1. **Resolved All Dependency Cycles**: Successfully implemented proper staging and dependencies
2. **Complete Bootstrap Token Integration**: Automated OAuth2 token retrieval and Docker integration
3. **Import Strategy**: Successfully imported built-in Edge sources into Terraform management
4. **Full End-to-End Deployment**: Complete Edge fleet configuration ready for deployment

## Current Deployment Status

### Ready for Deployment
The Edge configuration includes:
- **3-node Docker fleet** with proper port mapping and networking
- **Automatic bootstrap token** embedded in containers
- **Complete pipeline processing** with metadata enrichment
- **Multiple data sources** (imported built-ins + custom sources)
- **Dual destinations** (S3 archive + Cribl HTTP forwarding)
- **Generated deployment scripts** for easy execution

### Manual Configuration Required
- **UI Connections**: Built-in Edge sources require manual drag-and-drop connections in Cribl Cloud UI
- **Connection targets**: Connect imported sources to `edge-s3-archive` and `edge-cribl-http` destinations

## Next Steps

1. **Deploy Edge Fleet**:
   ```bash
   terraform apply
   ./deploy-edge-fleet.sh
   ```

2. **Configure Connections in UI**:
   - Open Cribl Cloud: `https://WORKSPACE-TENANT.cribl.cloud/?group=edge-fleet-tf2`
   - Drag connections from imported sources to destinations:
     - File Monitor → S3 Archive + Cribl HTTP
     - Journal Files → S3 Archive
     - System Metrics → Cribl HTTP
     - System State → S3 Archive

## Usage Notes

### Bootstrap Token Retrieval
Manual command for getting Edge bootstrap tokens:
```bash
../../get_access_token.sh -i CLIENT_ID -s CLIENT_SECRET -o ORG -w WORKSPACE -e ENVIRONMENT -g EDGE_GROUP
```

### Edge Fleet Deployment
Automated deployment using generated Docker Compose:
```bash
# Deploy the complete 3-node Edge fleet
./deploy-edge-fleet.sh

# Or manually using Docker Compose
docker-compose -f edge-docker-compose.yml up -d
```

### Individual Node Access
- Node 1: http://localhost:9001 (Syslog: 516/10516, HTTP: 8089, TCP: 8126)
- Node 2: http://localhost:9002 (Syslog: 517/10517, HTTP: 8090, TCP: 8127)  
- Node 3: http://localhost:9003 (Syslog: 518/10518, HTTP: 8091, TCP: 8128)

### Cribl Cloud Management
Access your Edge fleet: `https://WORKSPACE-TENANT.cribl.cloud/?group=edge-fleet-tf2`

## Architecture Benefits
- **Modularity**: Reusable components for sources/destinations
- **Scalability**: Edge fleet can scale across multiple nodes
- **Observability**: Integration with multiple monitoring platforms
- **Security**: Proper authentication and data masking
- **Flexibility**: Support for multiple deployment patterns (cloud, edge, hybrid)