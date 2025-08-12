# Cribl Destination Module

A reusable Terraform module for creating and managing Cribl destinations with various types and configurations.

## Usage

```hcl
module "my_destination" {
  source = "./destination-module"
  
  destination_id   = "unique-destination-id"
  group_id         = "worker-group-id"
  destination_type = "s3"  # See supported types below
  description      = "My destination description"
  
  # Type-specific variables
  bucket    = "my-s3-bucket"
  region    = "us-east-1"
  dest_path = "logs/"
  
  # Optional custom configuration
  custom_config = {
    max_file_size_mb = 64
  }
}
```

## Supported Destination Types

| Type | Description | Key Variables |
|------|-------------|---------------|
| `cribl_lake` | Cribl Lake storage | `dest_path`, `format`, `compress` |
| `s3` | Amazon S3 | `bucket`, `region`, `dest_path` |
| `splunk_hec` | Splunk HTTP Event Collector | `splunk_url`, `token` |
| `cribl_http` | HTTP endpoint | `url`, `method` |
| `cribl_tcp` | TCP connection | Use `custom_config` for host/port |
| `syslog` | Syslog output | Use `custom_config` for host/port/protocol |
| `kafka` | Apache Kafka | Use `custom_config` for brokers/topic |
| `kinesis` | AWS Kinesis | `region`, use `custom_config` for stream_name |
| `elasticsearch` | Elasticsearch/Elastic | `url` |
| `webhook` | Generic webhook | `url`, `method` |
| `open_telemetry` | OpenTelemetry Protocol (OTLP) | `url` |
| `crowdstrike_next_gen_siem` | CrowdStrike Next-Gen SIEM | `url`, `token` |

## Input Variables

### Required Variables

| Name | Type | Description |
|------|------|-------------|
| `destination_id` | `string` | Unique identifier for the destination |
| `group_id` | `string` | Worker group ID |
| `destination_type` | `string` | Type of destination (see supported types) |

### Common Optional Variables

| Name | Type | Default | Description |
|------|------|---------|-------------|
| `description` | `string` | `""` | Destination description |
| `disabled` | `bool` | `false` | Disable this destination |
| `streamtags` | `list(string)` | `[]` | Stream tags for filtering |
| `pipeline` | `string` | `null` | Pipeline to process events |
| `custom_config` | `any` | `{}` | Custom configuration to override defaults |

### Type-Specific Variables

| Name | Type | Default | Used By | Description |
|------|------|---------|---------|-------------|
| `dest_path` | `string` | `"default_logs"` | cribl_lake, s3 | Destination path |
| `format` | `string` | `"json"` | Multiple | Output format (json, ndjson, parquet) |
| `compress` | `string` | `"gzip"` | Multiple | Compression type (gzip, none) |
| `bucket` | `string` | `""` | s3 | S3 bucket name |
| `region` | `string` | `"us-east-1"` | s3, kinesis | AWS region |
| `url` | `string` | `""` | HTTP-based | Endpoint URL |
| `method` | `string` | `"POST"` | HTTP-based | HTTP method |
| `splunk_url` | `string` | `""` | splunk_hec | Splunk HEC URL |
| `token` | `string` | `""` | splunk_hec, crowdstrike | Authentication token (sensitive) |

## Outputs

| Name | Type | Description |
|------|------|-------------|
| `destination_id` | `string` | The ID of the created destination |
| `destination_type` | `string` | The type of the created destination |

## Examples

### S3 Destination

```hcl
module "s3_logs" {
  source = "./destination-module"
  
  destination_id   = "production-logs-s3"
  group_id         = "prod-workers"
  destination_type = "s3"
  description      = "Production logs to S3"
  
  bucket    = "my-company-logs"
  region    = "us-west-2"
  dest_path = "cribl/production/"
  format    = "json"
  compress  = "gzip"
  
  streamtags = ["production", "logs"]
  
  custom_config = {
    max_file_size_mb = 128
    storage_class    = "STANDARD_IA"
  }
}
```

### Splunk HEC Destination

```hcl
module "splunk_events" {
  source = "./destination-module"
  
  destination_id   = "splunk-main"
  group_id         = "default"
  destination_type = "splunk_hec"
  description      = "Main Splunk indexer"
  
  splunk_url = "https://splunk.example.com:8088/services/collector"
  token      = var.splunk_hec_token  # Store token in variable
  format     = "json"
  
  custom_config = {
    index              = "main"
    source_type        = "cribl:events"
    max_events_per_req = 100
  }
}
```

### OpenTelemetry Destination

```hcl
module "otlp_traces" {
  source = "./destination-module"
  
  destination_id   = "otlp-collector"
  group_id         = "default"
  destination_type = "open_telemetry"
  description      = "OTLP traces collector"
  
  url      = "https://otlp.example.com:4318/v1/traces"
  compress = "gzip"
  
  custom_config = {
    protocol     = "http"
    otlp_version = "1.3.1"
    metadata = [
      {
        key   = "service.name"
        value = "my-service"
      },
      {
        key   = "environment"
        value = "production"
      }
    ]
    http_traces_endpoint_override  = "/v1/traces"
    http_metrics_endpoint_override = "/v1/metrics"
    http_logs_endpoint_override    = "/v1/logs"
  }
}
```

### CrowdStrike Next-Gen SIEM Destination

```hcl
module "crowdstrike_siem" {
  source = "./destination-module"
  
  destination_id   = "crowdstrike-siem"
  group_id         = "security-workers"
  destination_type = "crowdstrike_next_gen_siem"
  description      = "Security events to CrowdStrike"
  
  url   = "https://api.crowdstrike.com/siem/v1/events"
  token = var.crowdstrike_token  # Store token in variable
  
  streamtags = ["security", "siem"]
  
  custom_config = {
    max_payload_events  = 500
    max_payload_size_kb = 4096
    flush_period_sec    = 5
    extra_http_headers = [
      {
        name  = "X-CS-Source"
        value = "cribl-stream"
      }
    ]
  }
}
```

### Webhook Destination with Custom Headers

```hcl
module "webhook_alerts" {
  source = "./destination-module"
  
  destination_id   = "alert-webhook"
  group_id         = "default"
  destination_type = "webhook"
  description      = "Alert webhook endpoint"
  
  url      = "https://webhook.site/unique-url"
  method   = "POST"
  format   = "ndjson"
  compress = "gzip"
  
  custom_config = {
    extra_http_headers = [
      {
        name  = "Authorization"
        value = "Bearer ${var.webhook_token}"
      },
      {
        name  = "X-Custom-Header"
        value = "custom-value"
      }
    ]
    timeout_sec         = 30
    max_payload_size_kb = 8192
  }
}
```

## Custom Configuration

The `custom_config` variable allows you to override any default setting or add destination-specific configurations not exposed as module variables. This provides flexibility for advanced use cases.

Example of overriding defaults:
```hcl
custom_config = {
  # Override default values
  on_backpressure = "drop"
  max_retry_num   = 5
  timeout         = 60
  
  # Add type-specific settings
  tls = {
    disabled = false
    ca_path  = "/path/to/ca.pem"
  }
}
```

## Requirements

- Terraform >= 0.13
- Cribl Provider (criblio/criblio)

## Provider Configuration

The module expects the Cribl provider to be configured in the parent module:

```hcl
provider "criblio" {
  cloud_url     = "https://api.cribl.cloud"
  client_id     = var.cribl_client_id
  client_secret = var.cribl_client_secret
}
```

## Notes

- Sensitive values like tokens should be stored in variables or secret management systems
- The module uses `create_before_destroy` lifecycle to ensure zero downtime during updates
- Each destination type has specific default configurations optimized for that type
- Use `terraform plan` to preview changes before applying
- Some destination types require additional infrastructure (e.g., S3 buckets, Kafka clusters) to be set up separately

## License

This module is part of the terraform-provider-criblio project.