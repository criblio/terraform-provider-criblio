# JSON Config Router Destination Module

This Terraform module allows you to create Cribl Router destinations directly from JSON configuration files exported from the Cribl Stream UI.

## Features

- **Direct JSON Import**: Use JSON configurations exported from Cribl Stream UI
- **Flexible Input**: Support for both file-based and inline JSON configurations
- **Field Transformation**: Automatic handling of field name differences between UI and Terraform
- **Override Support**: Ability to override destination ID and add additional configuration
- **Schema Compliant**: Filters out unknown fields and handles all router-specific attributes

## Usage

### Basic Usage with JSON File

```terraform
module "router_destination" {
  source = "../../modules/json-config-destination-router"

  group_id       = "default"
  json_file_path = "./configs/router_destination.json"
}
```

### With Overrides

```terraform
module "router_destination_custom" {
  source = "../../modules/json-config-destination-router"

  group_id                = "default"
  json_file_path          = "./configs/router_destination.json"
  destination_id_override = "custom-router-dest"
  
  additional_config = {
    description = "Custom router destination"
    environment = "production"
  }
}
```

### Inline JSON Configuration

```terraform
module "router_destination_inline" {
  source = "../../modules/json-config-destination-router"

  group_id = "default"
  json_content = jsonencode({
    id          = "inline-router-dest"
    type        = "router"
    description = "Router destination from inline JSON"
    rules = [
      {
        description = "Route to S3"
        filter      = "source=='web-logs'"
        output      = "s3-destination"
        final       = true
      },
      {
        description = "Route to Splunk"
        filter      = "source=='app-logs'"
        output      = "splunk-destination"
        final       = false
      }
    ]
    systemFields = ["cribl_pipe"]
    streamtags   = ["router", "json-config"]
  })
}
```

## Example JSON Configuration

```json
{
  "id": "router-destination-example",
  "type": "router",
  "description": "Routes events to different outputs based on conditions",
  "environment": "production",
  "pipeline": "main",
  "rules": [
    {
      "description": "High priority events to Splunk",
      "filter": "severity >= 8",
      "output": "splunk-hec-destination",
      "final": true
    },
    {
      "description": "All other events to S3",
      "filter": "true",
      "output": "s3-archive",
      "final": false
    }
  ],
  "systemFields": ["cribl_pipe", "_time"],
  "streamtags": ["routing", "production"]
}
```

## Supported Fields

### Top-level Fields
- `id` (required) - Unique identifier for the router destination
- `type` (required) - Must be "router"
- `description` - Description of the destination
- `environment` - Git branch restriction
- `pipeline` - Processing pipeline
- `systemFields` - Automatically added fields
- `streamtags` - Tags for filtering and grouping

### Rules Configuration
- `rules` (required) - Array of routing rules
  - `filter` (required) - JavaScript expression for event selection
  - `output` (required) - Target destination ID
  - `description` - Rule description
  - `final` - Whether to stop processing other rules (default: true)

## Outputs

- `destination_id` - The ID of the created router destination
- `destination_type` - The type of the destination ("router")

## Requirements

- Terraform >= 0.13
- Criblio provider >= 1.0.0

## Field Transformation

The module automatically handles field name transformations:
- `systemFields` → `system_fields`

All other fields are passed through directly as they match the Terraform schema.

## Notes

- Router destinations require at least one rule in the `rules` array
- The `filter` field in rules uses JavaScript expressions
- Rules are processed in order, and `final: true` stops further rule processing
- Unknown fields from the JSON are automatically filtered out
