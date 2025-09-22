# JSON Config Search Saved Query Module

This Terraform module allows you to create Cribl Search saved queries directly from JSON configurations exported from the Cribl Stream UI.

## Features

- ✅ **Direct JSON Import**: Copy saved query JSON from Cribl UI
- ✅ **Field Mapping**: Automatic `camelCase` to `snake_case` conversion
- ✅ **Flexible Input**: Support for both file path and inline JSON content
- ✅ **Configuration Merging**: Additional Terraform-specific overrides
- ✅ **Schedule Support**: Full scheduling with cron expressions and notifications
- ✅ **Privacy Controls**: Support for private/public saved queries

## Usage

### Basic Usage with JSON File

```hcl
module "my_saved_query" {
  source = "../../modules/json-config-search-saved-query"

  json_file_path = "./configs/my_query.json"
}
```

### Inline JSON Content

```hcl
module "inline_query" {
  source = "../../modules/json-config-search-saved-query"

  json_content = jsonencode({
    id          = "my-query"
    name        = "My Query"
    query       = "dataset=\"events\" | stats count"
    earliest    = "-1h"
    latest      = "now"
    isPrivate   = true
  })
}
```

### With Configuration Overrides

```hcl
module "custom_query" {
  source = "../../modules/json-config-search-saved-query"

  json_file_path      = "./query.json"
  query_id_override   = "custom-query-id"
  
  additional_config = {
    description = "Customized query description"
    is_private  = false
  }
}
```

## JSON Configuration Example

```json
{
  "id": "sample-query",
  "name": "Sample Events Query",
  "description": "Query for recent events",
  "query": "dataset=\"cribl_edge_appscope_metrics\" | limit 1000",
  "earliest": "-1h",
  "latest": "now",
  "isPrivate": true,
  "schedule": {
    "cronSchedule": "0 */6 * * *",
    "enabled": false,
    "keepLastN": 10,
    "notifications": {
      "disabled": true
    },
    "tz": "UTC"
  }
}
```

## Field Mapping

| JSON Field (camelCase) | Terraform Attribute (snake_case) |
|----------------------|----------------------------------|
| `isPrivate` | `is_private` |
| `cronSchedule` | `cron_schedule` |
| `keepLastN` | `keep_last_n` |

## Variables

- `json_file_path` - Path to JSON file containing the saved query configuration
- `json_content` - JSON content for the saved query configuration (alternative to json_file_path)
- `query_id_override` - Override the query ID from the JSON configuration
- `additional_config` - Additional configuration to merge with the parsed JSON

## Outputs

- `query_id` - The ID of the created search saved query
- `query_name` - The name of the created search saved query
- `final_config` - The final merged configuration (sensitive)

## Required JSON Fields

- `id` - Unique identifier for the saved query
- `name` - Name of the saved query
- `query` - The search query string

## Optional JSON Fields

- `description` - Description of the saved query
- `earliest` - Earliest time for the search range
- `latest` - Latest time for the search range
- `isPrivate` - Whether the saved query is private
- `schedule` - Scheduling configuration with cron expressions

## Schedule Configuration

The `schedule` object supports:
- `cronSchedule` - Cron expression for scheduling
- `enabled` - Whether scheduling is enabled
- `keepLastN` - Number of results to keep
- `notifications` - Notification settings
- `tz` - Timezone for scheduling

## Example Output

After successful deployment:

```hcl
query_id   = "sample-query"
query_name = "Sample Events Query"
```
