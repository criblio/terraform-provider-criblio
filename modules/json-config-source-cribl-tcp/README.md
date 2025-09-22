# JSON Config Cribl HTTP Source Module

A Terraform module that creates Cribl HTTP sources from JSON configuration files. This is a specialized module for `cribl_http` source types.

## Features

✅ **Cribl HTTP Specific** - Optimized for cribl_http source configuration  
✅ **JSON Configuration** - Use JSON config files directly from Cribl exports  
✅ **Field Name Transformation** - Automatically converts camelCase to snake_case  
✅ **Flexible Input** - Support file paths or direct JSON content  
✅ **Override Support** - Override specific fields while preserving JSON config  

## Why Source-Type-Specific Modules?

Due to Terraform's validation behavior, a single "universal" module that supports all source types isn't practically feasible. Each source type has different required attributes that are validated even when not used. Therefore, we provide dedicated modules for each source type to ensure proper validation and configuration.

## Usage

### Basic Usage with JSON File

```hcl
module "http_source" {
  source = "./modules/json-config-source-cribl-http"
  
  group_id       = "my-worker-group"
  json_file_path = "./configs/http_source.json"
}
```

### Usage with Direct JSON Content

```hcl
module "http_source" {
  source = "./modules/json-config-source-cribl-http"
  
  group_id     = "my-workers"  
  json_content = jsonencode({
    id   = "my-http-source"
    type = "cribl_http"
    port = 10081
    host = "0.0.0.0"
    authTokens = ["token1", "token2"]
    # ... other cribl_http config
  })
}
```

### Usage with Overrides

```hcl
module "http_source" {
  source = "./modules/json-config-source-cribl-http"
  
  group_id            = "my-group"
  json_file_path      = "./configs/http_source.json"
  source_id_override  = "custom-source-id"
  
  additional_config = {
    description = "Overridden description"
    disabled    = false
  }
}
```

## JSON Configuration Example

Your JSON configuration should match the cribl_http format:

```json
{
  "id": "cribl_http_source",
  "type": "cribl_http",
  "activityLogSampleRate": 100,
  "authTokens": [
    "token1",
    "token2"
  ],
  "captureHeaders": false,
  "description": "HTTP Source for Cribl",
  "disabled": false,
  "enableHealthCheck": true,
  "host": "0.0.0.0",
  "port": 10081,
  "pqEnabled": true,
  "requestTimeout": 30,
  "sendToRoutes": true,
  "streamtags": [
    "http",
    "cribl",
    "source"
  ],
  "tls": {
    "disabled": true,
    "maxVersion": "TLSv1.3",
    "minVersion": "TLSv1.2"
  }
}
```

## Field Name Transformation

The module automatically transforms JSON field names to match Terraform schema:

| JSON Field Name | Terraform Field Name |
|----------------|---------------------|
| `activityLogSampleRate` | `activity_log_sample_rate` |
| `authTokens` | `auth_tokens` |
| `captureHeaders` | `capture_headers` |
| `enableHealthCheck` | `enable_health_check` |
| `pqEnabled` | `pq_enabled` |
| `sendToRoutes` | `send_to_routes` |

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| `group_id` | Worker group ID where the source will be created | `string` | n/a | yes |
| `json_file_path` | Path to JSON configuration file | `string` | `""` | no* |
| `json_content` | Direct JSON content as string | `string` | `""` | no* |
| `source_id_override` | Override the source ID from JSON config | `string` | `""` | no |
| `additional_config` | Additional configuration to merge with JSON config | `any` | `{}` | no |

*Either `json_file_path` or `json_content` must be provided.

## Outputs

| Name | Description |
|------|-------------|
| `source_id` | The ID of the created source |
| `source_type` | The type of source that was created |
| `group_id` | The group ID where the source was created |
| `parsed_config` | The parsed and transformed JSON configuration (sensitive) |

## Related Modules

- `json-config-source-syslog` - For syslog sources
- `json-config-source-http` - For regular HTTP sources  
- `json-config-source-tcp` - For TCP sources

## Requirements

- Terraform >= 1.0
- Cribl provider (`criblio/criblio`) >= 1.0.0
