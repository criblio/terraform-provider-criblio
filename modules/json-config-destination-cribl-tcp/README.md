# JSON Config Cribl TCP Destination Module

This Terraform module creates a Cribl TCP destination from JSON configuration exported from Cribl Stream.

## Features

- **Direct JSON Import**: Copy JSON configurations directly from Cribl Stream UI
- **Field Mapping**: Automatic conversion from `camelCase` (JSON) to `snake_case` (Terraform)
- **Unknown Field Filtering**: Ignores fields not defined in the provider schema
- **Flexible Input**: Supports both file path and inline JSON content
- **Configuration Merging**: Allows additional Terraform-specific overrides

## Usage

### From JSON File

```hcl
module "cribl_tcp_destination" {
  source = "./modules/json-config-destination-cribl-tcp"

  group_id       = "default"
  json_file_path = "./configs/cribl_tcp_destination.json"
}
```

### From Inline JSON

```hcl
module "cribl_tcp_destination" {
  source = "./modules/json-config-destination-cribl-tcp"

  group_id     = "default"
  json_content = jsonencode({
    id          = "my-cribl-tcp-dest"
    type        = "cribl_tcp"
    description = "Cribl TCP destination"
    host        = "cribl-receiver"
    port        = 10300
    compression = "gzip"
  })
}
```

### With Load Balancing

```hcl
module "cribl_tcp_destination" {
  source = "./modules/json-config-destination-cribl-tcp"

  group_id       = "default"
  json_file_path = "./configs/cribl_tcp_destination.json"
  
  additional_config = {
    load_balanced = true
    hosts = [
      {
        host = "cribl-receiver-1"
        port = 10300
        weight = 1
      },
      {
        host = "cribl-receiver-2"
        port = 10300
        weight = 1
      }
    ]
  }
}
```

## Field Mappings

The module automatically handles field name transformations:

- `connectionTimeout` → `connection_timeout`
- `dnsResolvePeriodSec` → `dns_resolve_period_sec`
- `loadBalanced` → `load_balanced`
- `maxConcurrentSenders` → `max_concurrent_senders`
- `throttleRatePerSec` → `throttle_rate_per_sec`
- And many more...

## Supported Fields

All fields from the Cribl TCP destination schema are supported, including:

- Basic settings (compression, connection timeout)
- Load balancing configuration
- Host and port settings
- TLS configuration
- Persistent queue settings
- System fields and tags
- Throttling settings

## Example JSON Configuration

```json
{
  "id": "cribl-tcp-dest",
  "type": "cribl_tcp",
  "description": "Cribl TCP destination",
  "host": "cribl-receiver",
  "port": 10300,
  "compression": "gzip",
  "connectionTimeout": 10000,
  "loadBalanced": false,
  "onBackpressure": "block",
  "systemFields": ["cribl_pipe"],
  "streamtags": ["cribl-tcp"],
  "tls": {
    "disabled": true
  }
}
```

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| group_id | The group ID to create the destination in | `string` | n/a | yes |
| json_file_path | Path to JSON file containing destination configuration | `string` | `""` | no |
| json_content | JSON content for destination configuration | `string` | `""` | no |
| destination_id_override | Override the destination ID from JSON | `string` | `""` | no |
| additional_config | Additional configuration to merge | `any` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| destination_id | The ID of the created Cribl TCP destination |
| destination_type | The type of the created destination |
| group_id | The group ID where the destination was created |
| final_config | The final merged configuration (sensitive) |
