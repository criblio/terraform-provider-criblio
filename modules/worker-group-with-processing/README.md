# Worker Group with Processing Module

This module creates a complete Cribl Stream worker group with optional pack installation, processing pipeline, and automatic deployment capabilities.

## Features

- **Worker Group Creation**: Creates a cloud-native worker group in AWS or Azure
- **Pack Installation**: Optional installation of processing packs from local files
- **Processing Pipeline**: Optional creation of data processing pipelines with region metadata
- **Automatic Deployment**: Optional automatic commit and deploy of configurations
- **Flexible Configuration**: Extensive customization options through variables

## Usage

### Basic Usage

```hcl
module "aws_east_workers" {
  source = "../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "aws-east-workers"
  group_name            = "AWS East Workers"
  cloud_provider        = "aws"
  cloud_region          = "us-east-1"
  estimated_ingest_rate = 1024
  description           = "AWS East region worker group"
  streamtags            = ["aws", "east", "production"]

  # Pack Configuration
  pack_filename     = "cribl-palo-alto-networks-source-1.0.0.crbl"
  pack_description  = "Palo Alto Networks processing pack"
  pack_display_name = "Palo Alto Networks Source Pack"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for AWS East"
  region_identifier    = "aws-east"
  pipeline_streamtags  = ["aws", "east", "processing"]

  # Deployment Configuration
  commit_message = "Deploy pack and pipeline to AWS East workers"
}
```

### Advanced Usage with Conditional Resources

```hcl
module "azure_europe_workers" {
  source = "../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "azure-europe-workers"
  group_name            = "Azure Europe Workers"
  cloud_provider        = "azure"
  cloud_region          = "westeurope"
  estimated_ingest_rate = 2048
  streamtags            = ["azure", "europe", "production"]

  # Conditional Pack Installation
  install_pack      = true
  pack_filename     = "custom-processing-pack-2.1.0.crbl"
  pack_description  = "Custom processing pack for European data"
  pack_display_name = "Custom European Processing Pack"
  pack_version      = "2.1.0"

  # Conditional Pipeline Creation
  create_pipeline      = true
  pipeline_description = "European GDPR-compliant processing pipeline"
  region_identifier    = "azure-europe"

  # Manual Deployment
  auto_commit = true
  auto_deploy = false  # Manual deploy for production safety
}
```

## Requirements

- Terraform >= 1.0
- Cribl Terraform Provider
- Pack files (.crbl) must be present in the same directory if `install_pack = true`

## Variables

### Worker Group Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `group_id` | Unique identifier for the worker group | `string` | n/a | yes |
| `group_name` | Display name for the worker group | `string` | n/a | yes |
| `cloud_provider` | Cloud provider (aws or azure) | `string` | n/a | yes |
| `cloud_region` | Cloud region for deployment | `string` | n/a | yes |
| `estimated_ingest_rate` | Estimated data ingest rate in MB/s | `number` | `1024` | no |
| `description` | Worker group description | `string` | `""` | no |
| `streamtags` | Stream tags for the worker group | `list(string)` | `[]` | no |

### Pack Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `install_pack` | Whether to install a pack | `bool` | `true` | no |
| `pack_filename` | Pack file to install | `string` | `""` | no |
| `pack_description` | Pack description | `string` | `""` | no |
| `pack_display_name` | Pack display name | `string` | `""` | no |
| `pack_version` | Pack version | `string` | `"1.0.0"` | no |

### Pipeline Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `create_pipeline` | Whether to create processing pipeline | `bool` | `true` | no |
| `pipeline_description` | Pipeline description | `string` | `"Standard data processing pipeline"` | no |
| `region_identifier` | Region identifier for metadata | `string` | n/a | yes |
| `pipeline_streamtags` | Pipeline stream tags | `list(string)` | `["processing"]` | no |

### Deployment Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `auto_commit` | Whether to auto-commit changes | `bool` | `true` | no |
| `auto_deploy` | Whether to auto-deploy | `bool` | `true` | no |
| `commit_message` | Commit message | `string` | `"Deploy pack and processing pipeline"` | no |

## Outputs

### Primary Outputs

| Name | Description |
|------|-------------|
| `worker_group_id` | ID of the created worker group |
| `worker_group_name` | Name of the created worker group |
| `cloud_details` | Cloud provider and region information |
| `summary` | Complete summary of all created resources |

### Conditional Outputs

| Name | Description |
|------|-------------|
| `pack_details` | Pack installation details (if pack installed) |
| `pipeline_details` | Pipeline configuration details (if pipeline created) |
| `deployment_status` | Deployment status (if auto-deployed) |

## Examples

See the parent example directory for complete usage examples with multiple worker groups.

## Best Practices

1. **Pack Files**: Ensure pack files are present in the working directory before applying
2. **Region Identifiers**: Use consistent region identifiers for better data tracking
3. **Stream Tags**: Use meaningful stream tags for organization and filtering
4. **Deployment Strategy**: Consider manual deployment (`auto_deploy = false`) for production environments
5. **Resource Naming**: Use descriptive and consistent naming conventions for group IDs

## Dependencies

This module creates resources in the following order:
1. Worker Group
2. Pack Installation (if enabled)
3. Processing Pipeline (if enabled)  
4. Configuration Commit (if enabled)
5. Configuration Deployment (if enabled)

All dependencies are handled automatically through Terraform's resource graph.
