# JSON Config Cribl HTTP Destination Module

This Terraform module creates a Cribl HTTP destination from JSON configuration exported from Cribl Stream.

## Features

- **Direct JSON Import**: Copy JSON configurations directly from Cribl Stream UI
- **Field Mapping**: Automatic conversion from `camelCase` (JSON) to `snake_case` (Terraform)
- **Unknown Field Filtering**: Ignores fields not defined in the provider schema
- **Flexible Input**: Supports both file path and inline JSON content
- **Configuration Merging**: Allows additional Terraform-specific overrides

## Usage

### From JSON File

```hcl
module "cribl_http_destination" {
  source = "./modules/json-config-destination-cribl-http"

  group_id       = "default"
  json_file_path = "./configs/cribl_http_destination.json"
}
```

### From Inline JSON

```hcl
module "cribl_http_destination" {
  source = "./modules/json-config-destination-cribl-http"

  group_id     = "default"
  json_content = jsonencode({
    id          = "my-cribl-http-dest"
    type        = "cribl_http"
    description = "Cribl HTTP destination"
    url         = "http://cribl-receiver:10080"
    compression = "gzip"
  })
}
```

### With Overrides

```hcl
module "cribl_http_destination" {
  source = "./modules/json-config-destination-cribl-http"

  group_id              = "default"
  json_file_path        = "./configs/cribl_http_destination.json"
  destination_id_override = "custom-cribl-http-id"
  
  additional_config = {
    description = "Overridden description"
    compression = "none"
  }
}
```

## Field Mappings

The module automatically handles field name transformations:

- `flushPeriodSec` → `flush_period_sec`
- `maxPayloadSizeKb` → `max_payload_size_kb`
- `loadBalanced` → `load_balanced`
- `dnsResolvePeriodSec` → `dns_resolve_period_sec`
- And many more...

## Supported Fields

All fields from the Cribl HTTP destination schema are supported, including:

- Basic settings (compression, concurrency, timeout)
- Load balancing configuration
- Authentication settings
- TLS configuration
- Persistent queue settings
- Retry settings
- Headers and metadata

## Example JSON Configuration

```json
{
  "id": "cribl-http-dest",
  "type": "cribl_http",
  "description": "Cribl HTTP destination",
  "compression": "gzip",
  "concurrency": 5,
  "url": "http://cribl-receiver:10080",
  "flushPeriodSec": 1,
  "maxPayloadSizeKb": 4096,
  "loadBalanced": false,
  "onBackpressure": "block",
  "systemFields": ["cribl_pipe"],
  "streamtags": ["cribl-http"]
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
| destination_id | The ID of the created Cribl HTTP destination |
| destination_type | The type of the created destination |
| group_id | The group ID where the destination was created |
| final_config | The final merged configuration (sensitive) |
