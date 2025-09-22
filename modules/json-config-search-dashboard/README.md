# JSON Config Search Dashboard Module

This Terraform module allows you to create Cribl Search dashboards directly from JSON configurations exported from the Cribl Stream UI.

## Features

- ✅ **Direct JSON Import**: Copy dashboard JSON from Cribl UI
- ✅ **Field Mapping**: Automatic `camelCase` to `snake_case` conversion
- ✅ **Complex Elements**: Support for charts, tables, inputs, and markdown elements
- ✅ **Flexible Input**: Support for both file path and inline JSON content
- ✅ **Configuration Merging**: Additional Terraform-specific overrides
- ✅ **Schedule Support**: Dashboard scheduling with cron expressions
- ✅ **Layout Management**: Grid-based positioning system

## Usage

### Basic Usage with JSON File

```hcl
module "my_dashboard" {
  source = "../../modules/json-config-search-dashboard"

  json_file_path = "./configs/my_dashboard.json"
}
```

### Inline JSON Content

```hcl
module "simple_dashboard" {
  source = "../../modules/json-config-search-dashboard"

  json_content = jsonencode({
    id         = "simple-dash"
    name       = "Simple Dashboard"
    created    = 1695403200
    created_by = "admin@example.com"
    modified   = 1695403200
    elements   = []
  })
}
```

### With Configuration Overrides

```hcl
module "custom_dashboard" {
  source = "../../modules/json-config-search-dashboard"

  json_file_path         = "./dashboard.json"
  dashboard_id_override  = "custom-dashboard-id"
  
  additional_config = {
    description  = "Customized dashboard"
    category     = "Operations"
    refresh_rate = 60
  }
}
```

## JSON Configuration Example

```json
{
  "id": "operations-dashboard",
  "name": "Operations Dashboard",
  "description": "Real-time operations monitoring",
  "category": "Operations",
  "created": 1695403200,
  "createdBy": "admin@example.com",
  "modified": 1695403260,
  "displayCreatedBy": "Admin User",
  "cacheTtlSeconds": 300,
  "refreshRate": 30,
  "elements": [
    {
      "element": {
        "id": "chart1",
        "title": "Event Count Over Time",
        "type": "chart.line",
        "variant": "visualization",
        "layout": {
          "x": 0, "y": 0, "w": 12, "h": 6
        },
        "search": {
          "search_query_inline": {
            "type": "inline",
            "query": "dataset=\"events\" | timechart count",
            "earliest": { "str": "-1h" },
            "latest": { "str": "now" }
          }
        },
        "value": {
          "xAxis": "\"_time\"",
          "yAxis": "\"count\""
        }
      }
    }
  ],
  "schedule": {
    "cronSchedule": "0 8 * * *",
    "enabled": false,
    "keepLastN": 5,
    "notifications": { "disabled": true },
    "tz": "UTC"
  }
}
```

## Field Mapping

| JSON Field (camelCase) | Terraform Attribute (snake_case) |
|----------------------|----------------------------------|
| `createdBy` | `created_by` |
| `modifiedBy` | `modified_by` |
| `displayCreatedBy` | `display_created_by` |
| `displayModifiedBy` | `display_modified_by` |
| `cacheTtlSeconds` | `cache_ttl_seconds` |
| `refreshRate` | `refresh_rate` |
| `resolvedDatasetIds` | `resolved_dataset_ids` |
| `packId` | `pack_id` |

## Variables

- `json_file_path` - Path to JSON file containing the dashboard configuration
- `json_content` - JSON content for the dashboard configuration (alternative to json_file_path)
- `dashboard_id_override` - Override the dashboard ID from the JSON configuration
- `additional_config` - Additional configuration to merge with the parsed JSON

## Outputs

- `dashboard_id` - The ID of the created search dashboard
- `dashboard_name` - The name of the created search dashboard
- `final_config` - The final merged configuration (sensitive)

## Required JSON Fields

- `id` - Unique identifier for the dashboard
- `name` - Name of the dashboard
- `created` - Creation timestamp
- `created_by` - Creator identifier
- `modified` - Modification timestamp
- `elements` - Dashboard elements array

## Optional JSON Fields

- `description` - Description of the dashboard
- `category` - Dashboard category
- `cache_ttl_seconds` - Cache TTL in seconds
- `display_created_by` - Display name of creator
- `display_modified_by` - Display name of modifier
- `modified_by` - Modifier identifier
- `pack_id` - Associated pack ID
- `refresh_rate` - Auto-refresh rate in seconds
- `resolved_dataset_ids` - Array of dataset IDs
- `schedule` - Scheduling configuration

## Element Types

Supported dashboard element types:
- `chart.line` - Line charts
- `chart.column` - Column charts
- `chart.pie` - Pie charts
- `list.table` - Data tables
- `list.events` - Event lists
- `counter.single` - Single value counters
- `input.timerange` - Time range inputs
- `input.dropdown` - Dropdown inputs
- `markdown.default` - Markdown elements

## Layout System

Elements use a grid-based layout:
```json
"layout": {
  "x": 0,  // X position (0-11)
  "y": 0,  // Y position
  "w": 6,  // Width (1-12)
  "h": 4   // Height
}
```

## Search Query Types

1. **Inline Queries**: Embedded search queries
2. **Saved Queries**: References to saved query IDs
3. **Values Queries**: Static value arrays

## Value Field Format

The `value` field in elements must contain JSON-formatted strings:
```json
"value": {
  "xAxis": "\"_time\"",
  "yAxis": "\"count\"",
  "columns": "\"auto\""
}
```

## Example Output

After successful deployment:

```hcl
dashboard_id   = "operations-dashboard"
dashboard_name = "Operations Dashboard"
```
