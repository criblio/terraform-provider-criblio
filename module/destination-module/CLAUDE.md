# Cribl Terraform Provider - Destination Module

## Overview
This module provides a reusable Terraform configuration for creating Cribl destinations with various types.

## Project Structure
```
destination-module/
├── locals.tf      # Base configurations for each destination type
├── main.tf        # Main resource definition
├── outputs.tf     # Output values
├── provider.tf    # Provider configuration
├── variables.tf   # Input variables
└── CLAUDE.md      # This reference file
```

## Supported Destination Types
- `cribl_lake` - Cribl Lake storage
- `s3` - Amazon S3
- `splunk_hec` - Splunk HTTP Event Collector
- `cribl_http` - HTTP endpoint
- `cribl_tcp` - TCP connection
- `syslog` - Syslog output
- `kafka` - Apache Kafka
- `kinesis` - AWS Kinesis
- `elasticsearch` - Elasticsearch/Elastic
- `webhook` - Generic webhook
- `open_telemetry` - OpenTelemetry Protocol (OTLP)
- `crowdstrike_next_gen_siem` - CrowdStrike Next-Gen SIEM (ng-siem)

## Key Implementation Details

### Adding New Destinations
1. Add base configuration in `locals.tf`:
   - Define `base_<destination_type>_config` with default settings
   - Add transformation config that merges base with custom_config
   
2. Update `main.tf`:
   - Add `output_<destination_type>` attribute to the resource

3. Update `variables.tf`:
   - Add the new type to the validation condition
   - Update the error message

### Resource Naming Convention
The Terraform resource uses specific output names:
- `output_elastic` (not `output_elasticsearch`)
- `output_crowdstrike_next_gen_siem` (for ng-siem)
- `output_open_telemetry` (for OTLP)

### Common Variables
- `destination_id` - Unique identifier
- `group_id` - Worker group ID
- `destination_type` - Type of destination
- `description` - Optional description
- `disabled` - Enable/disable flag
- `streamtags` - Stream tags list
- `pipeline` - Pipeline configuration
- `custom_config` - Override any default settings

### Testing Commands
```bash
# Format check
terraform fmt -check

# Validation
terraform validate

# Plan
terraform plan
```

## Example Usage
```hcl
module "otlp_destination" {
  source = "./destination-module"
  
  destination_id   = "otlp-1"
  group_id        = "default"
  destination_type = "open_telemetry"
  description     = "OpenTelemetry destination"
  url             = "https://otlp-collector.example.com:4318"
  
  custom_config = {
    otlp_version = "1.0.0"
    protocol     = "grpc"
  }
}

module "crowdstrike_destination" {
  source = "./destination-module"
  
  destination_id   = "cs-siem-1"
  group_id        = "default"
  destination_type = "crowdstrike_next_gen_siem"
  description     = "CrowdStrike SIEM destination"
  url             = "https://api.crowdstrike.com/siem/v1"
  token           = var.crowdstrike_token
}
```

## Reference Files
- Examples: `/examples/resources/criblio_destination/resource.tf`
- Provider docs: Check the terraform-provider-criblio repository

## Notes
- The module uses a pattern where base configurations are defined in `locals.tf` and merged with `custom_config`
- Each destination type has its own specific settings but shares common fields
- The `custom_config` variable allows overriding any default setting
- Sensitive data like tokens should be passed as variables, not hardcoded