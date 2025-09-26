# Multi Worker Groups with Dual Pack Installation Example

This example demonstrates how to create multiple worker groups across different cloud providers (AWS and Azure) using a **reusable module** that installs **both CrowdStrike and Palo Alto packs** with processing pipelines. This modular approach eliminates code duplication and makes the configuration highly maintainable.

## Architecture Overview

This configuration creates:

1. **Six Worker Groups (via Module):**
   - **AWS East**: Worker group in `us-east-1` region
   - **AWS West**: Worker group in `us-west-2` region
   - **AWS Central**: Worker group in `us-central-1` region
   - **Azure East**: Worker group in `eastus` region
   - **Azure West**: Worker group in `westus2` region
   - **Azure Europe**: Worker group in `westeurope` region

2. **Modular Design**: Uses the `worker-group-with-processing` module for consistent deployment

3. **Dual Pack Installation**:
   - **Palo Alto Networks Pack**: Installed on all worker groups from local `.crbl` file
   - **CrowdStrike Security Pack**: Installed on all worker groups with event breaker rulesets

4. **Advanced Processing**:
   - **Event Breakers**: CrowdStrike event breaker rulesets for proper event parsing and timestamping
   - **Processing Pipelines**: Standard data processing pipelines with region metadata

5. **Automated Deployment**: Each worker group has automatic commit and deploy cycles via the module

## Module Benefits

### 🎯 **Code Reduction**: From 696 lines to 183 lines (73% reduction)
### 🔄 **Reusability**: Single module definition used 6 times
### 🛠️ **Maintainability**: Changes to the module automatically apply to all worker groups
### 📦 **Consistency**: Identical configuration across all deployments
### 🚀 **Scalability**: Easy to add more worker groups by adding more module calls

## Configuration Details

### Module Usage Pattern

Each worker group is created using the `worker-group-with-processing` module with dual pack support:

```hcl
module "aws_east_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "aws-east-workers"
  group_name            = "AWS East Workers"
  cloud_provider        = "aws"
  cloud_region          = "us-east-1"
  estimated_ingest_rate = 1024
  
  # Palo Alto Pack Configuration
  install_palo_alto_pack        = true
  palo_alto_pack_filename       = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description    = "Palo Alto Networks pack for AWS East workers"
  palo_alto_pack_display_name   = "Palo Alto Networks Source Pack - AWS East"
  palo_alto_pack_version        = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for AWS East workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - AWS East"
  
  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for AWS East region"
  region_identifier    = "aws-east"
  
  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to AWS East workers"
}
```

### Worker Groups

| Group | Cloud Provider | Region | Estimated Ingest Rate | Stream Tags |
|-------|---------------|--------|----------------------|-------------|
| AWS East Workers | AWS | us-east-1 | 1024 MB/s | aws, east, production |
| AWS West Workers | AWS | us-west-2 | 2048 MB/s | aws, west, production |
| AWS Central Workers | AWS | us-central-1 | 1536 MB/s | aws, central, production |
| Azure East Workers | Azure | eastus | 1536 MB/s | azure, east, production |
| Azure West Workers | Azure | westus2 | 1024 MB/s | azure, west, production |
| Azure Europe Workers | Azure | westeurope | 2048 MB/s | azure, europe, production |

### Module Features

- **Dual Pack Installation**: 
  - **Palo Alto Networks Pack**: From local `.crbl` file with network security processing
  - **CrowdStrike Security Pack**: With security data processing and metadata enrichment
- **Advanced Processing**: 
  - **Event Breakers**: CrowdStrike event breaker rulesets with timestamp extraction
  - **Processing Pipeline**: Standard data pipeline with region-specific metadata
- **Automatic Deployment**: Built-in commit and deploy cycle for all resources
- **Flexible Configuration**: Extensive customization through module variables

### Deployment Process

For each worker group, the following sequence occurs:

1. **Worker Group Creation**: Creates the worker group with cloud-specific configuration
2. **Pack Installation**: Installs the pack from the local file
3. **Pipeline Creation**: Deploys the data processing pipeline with region-specific metadata
4. **Commit**: Commits the configuration changes with a descriptive message
5. **Config Version**: Retrieves the committed configuration version
6. **Deploy**: Deploys the configuration to the worker group

## Prerequisites

1. **Cribl Stream Instance**: Access to a Cribl Stream instance with appropriate permissions
2. **Cloud Credentials**: Proper credentials and permissions for AWS and Azure
3. **Pack File**: The `cribl-palo-alto-networks-source-1.0.0.crbl` file must be present in the same directory
4. **Terraform Provider**: Cribl Terraform Provider configured with valid credentials

## Required Files

Ensure these files are in the same directory:

```
multi-worker-groups-with-pack/
├── main.tf
├── README.md
├── cribl-palo-alto-networks-1.1.5.crbl
└── CrowdstrikeEventBreakerRuleset.json (reference)
```

**Note**: The CrowdStrike JSON configuration is converted to a `pack_pipeline` resource directly in the Terraform module, so only the Palo Alto `.crbl` file needs to be present for pack installation.

## Usage

### 1. Provider Configuration

First, configure the Cribl Terraform Provider in your `terraform.tf` or `provider.tf` file:

```hcl
terraform {
  required_providers {
    criblio = {
      source = "registry.terraform.io/cribldata/criblio"
      version = "~> 1.0"
    }
  }
}

provider "criblio" {
  # Configure with your organization and workspace details
  # organization_id = "your-org-id"
  # workspace_id = "your-workspace-id"
  # base_url = "https://your-instance.cribl.cloud"
  # token = "your-api-token"
}
```

### 2. Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Review the planned changes
terraform plan

# Apply the configuration
terraform apply
```

### 3. Verification

After successful deployment, you can verify:

- Worker groups are created in the specified cloud regions
- Packs are installed and enabled on all worker groups
- Configuration is committed and deployed to all groups

## Outputs

The configuration provides comprehensive outputs:

### Worker Groups Output
Details about all 6 created worker groups including IDs, names, regions, providers, and estimated ingest rates.

### Palo Alto Pack Installations Output  
Information about Palo Alto pack installations on each worker group including pack IDs, versions, and status for all 6 groups.

### CrowdStrike Pack Installations Output
Information about CrowdStrike pack installations on each worker group including pack IDs, versions, and status for all 6 groups.

### CrowdStrike Pack Pipelines Output
Details about the CrowdStrike pack pipelines (event breaker rulesets) deployed to each worker group.

### Processing Pipelines Output
Details about the standard processing pipelines deployed to each worker group including pipeline IDs and descriptions.

### Deployment Status Output
Deployment status for all 6 worker groups including config versions and commit messages.

## Customization

### Adding More Worker Groups

Adding new worker groups is now extremely simple with the module:

```hcl
module "aws_south_workers" {
  source = "../../modules/worker-group-with-processing"

  group_id              = "aws-south-workers"
  group_name            = "AWS South Workers"
  cloud_provider        = "aws"
  cloud_region          = "sa-east-1"
  estimated_ingest_rate = 1024
  streamtags            = ["aws", "south", "production"]

  # Pack and pipeline configuration automatically handled
  pack_filename         = "cribl-palo-alto-networks-source-1.0.0.crbl"
  pack_description      = "Palo Alto Networks pack for AWS South workers"
  pack_display_name     = "Palo Alto Networks Source Pack - AWS South"
  region_identifier     = "aws-south"
  commit_message        = "Deploy pack and pipeline to AWS South workers"
}
```

That's it! No need to create individual resources or manage dependencies.

### Using Different Packs

To use a different pack:

1. Replace the `.crbl` file with your desired pack
2. Update the `filename` attribute in all `criblio_pack` resources
3. Update pack metadata (description, display_name, version) as needed

### Modifying Worker Group Settings

Common customizations:

- **Regions**: Change the `region` in the `cloud` block
- **Ingest Rate**: Modify `estimated_ingest_rate` based on expected load
- **Tags**: Update `streamtags` for better organization
- **Provisioning**: Set `provisioned = false` for on-demand scaling

## Troubleshooting

### Common Issues

1. **Pack File Not Found**: Ensure the `.crbl` file is in the same directory as `main.tf`
2. **Cloud Permissions**: Verify cloud provider credentials and permissions
3. **API Rate Limits**: If deploying to many groups, consider adding delays between deployments
4. **Dependency Conflicts**: Ensure proper dependency chains using `depends_on`

### Verification Commands

```bash
# Check all worker groups
terraform output worker_groups

# Check Palo Alto pack installations on all groups
terraform output palo_alto_pack_installations

# Check CrowdStrike pack installations on all groups
terraform output crowdstrike_pack_installations

# Check CrowdStrike pack pipelines on all groups
terraform output crowdstrike_pack_pipelines

# Check standard processing pipelines on all groups
terraform output pipelines

# Check deployment status for all groups
terraform output deployment_status

# Quick deployment summary
terraform output deployment_summary
```

## Clean Up

To remove all resources:

```bash
terraform destroy
```

**Note**: This will remove all worker groups, packs, and configurations. Ensure you have backups if needed.

## Advanced Configurations

### Conditional Pack Installation

You can make pack installation conditional:

```hcl
variable "install_pack" {
  description = "Whether to install the pack on worker groups"
  type        = bool
  default     = true
}

resource "criblio_pack" "aws_east_pack" {
  count = var.install_pack ? 1 : 0
  # ... rest of configuration
}
```

### Environment-Specific Configurations

Use variables to customize per environment:

```hcl
variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

resource "criblio_group" "aws_east_group" {
  id = "${var.environment}-aws-east-workers"
  streamtags = [
    "aws",
    "east", 
    var.environment
  ]
  # ... rest of configuration
}
```

## Module Structure

```
modules/worker-group-with-processing/
├── main.tf          # Core module logic
├── variables.tf     # Input variable definitions  
├── outputs.tf       # Output value definitions
└── README.md        # Module documentation
```

The module encapsulates all the complexity of creating worker groups, installing packs, creating pipelines, and managing deployments into a single, reusable component.

## Summary

This example provides a comprehensive foundation for managing multiple worker groups with consistent pack and pipeline installations across different cloud providers using Infrastructure as Code principles. The **modular approach** delivers significant benefits:

### 🚀 **Key Benefits**
- **73% Code Reduction**: From 696 to 183 lines
- **DRY Principle**: No repeated code
- **Easy Maintenance**: Update module once, affects all deployments
- **Rapid Scaling**: Add new regions in minutes
- **Consistent Deployments**: Identical configuration guaranteed

### 🎯 **Perfect For**
- **Enterprise Deployments**: Consistent processing across multiple regions
- **Multi-Cloud Strategy**: Seamless AWS and Azure integration  
- **DevOps Teams**: Infrastructure as Code best practices
- **Scalable Architecture**: Easy horizontal scaling
- **Compliance Requirements**: Consistent configuration and processing

### 🏗️ **Architecture Benefits**
- **Modular Design**: Reusable components
- **Dependency Management**: Automatic resource ordering
- **Error Handling**: Built-in validation and safety checks
- **Monitoring**: Comprehensive status outputs
- **Documentation**: Self-documenting infrastructure

This modular approach transforms complex multi-region deployments into simple, maintainable configurations that follow Infrastructure as Code best practices.
