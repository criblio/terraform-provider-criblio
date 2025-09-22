# JSON Config Search Examples

This example demonstrates the **JSON-config search modules** for creating Cribl search dashboards and saved queries directly from JSON configurations exported from Cribl Stream UI.

## 🎯 Features

- **Direct JSON Import**: Copy search configurations straight from Cribl UI
- **Field Mapping**: Automatic `camelCase` to `snake_case` conversion
- **Flexible Input**: Support for both file path and inline JSON content
- **Configuration Merging**: Allows additional Terraform-specific overrides
- **Production Ready**: Validated against provider schemas

## 📦 Included Modules

### 1. Search Saved Query
**Module**: `json-config-search-saved-query`

**Features**:
- Simple query definition with time ranges
- Optional scheduling with cron expressions
- Privacy controls
- Notification settings

**Required Fields**: `id`, `name`, `query`
**Optional Fields**: `description`, `earliest`, `latest`, `is_private`, `schedule`

### 2. Search Dashboard
**Module**: `json-config-search-dashboard`

**Features**:
- Complex dashboard layouts with multiple elements
- Chart visualizations (line, column, pie, etc.)
- Table and event list displays
- Inline and saved query integration
- Dashboard scheduling and caching

**Required Fields**: `id`, `name`, `created`, `created_by`, `modified`, `elements`
**Optional Fields**: `description`, `category`, `cache_ttl_seconds`, `refresh_rate`, `schedule`

## 📄 Configuration Examples

### Saved Query JSON

```json
{
  "id": "my-saved-query",
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

### Dashboard JSON

```json
{
  "id": "my-dashboard",
  "name": "Operations Dashboard",
  "description": "Real-time operations monitoring",
  "category": "Operations",
  "created": 1695403200,
  "createdBy": "admin@example.com",
  "modified": 1695403260,
  "elements": [
    {
      "element": {
        "id": "chart1",
        "title": "Event Count Over Time",
        "type": "chart.line",
        "layout": { "x": 0, "y": 0, "w": 12, "h": 6 },
        "search": {
          "search_query_inline": {
            "type": "inline",
            "query": "dataset=\"events\" | timechart count",
            "earliest": { "str": "-1h" },
            "latest": { "str": "now" }
          }
        }
      }
    }
  ]
}
```

## 🚀 Usage

### Deploy Examples

```bash
cd examples/json-config-search

# Initialize
terraform init

# Validate configuration
terraform validate

# Plan deployment
terraform plan

# Deploy search resources
terraform apply
```

### Using Individual Modules

#### Saved Query from File

```hcl
module "my_saved_query" {
  source = "../../modules/json-config-search-saved-query"

  json_file_path = "./my-query.json"
}
```

#### Dashboard with Inline JSON

```hcl
module "my_dashboard" {
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

#### With Configuration Overrides

```hcl
module "custom_dashboard" {
  source = "../../modules/json-config-search-dashboard"

  json_file_path        = "./dashboard.json"
  dashboard_id_override = "custom-dashboard-id"
  
  additional_config = {
    description  = "Customized dashboard"
    category     = "Custom"
    refresh_rate = 60
  }
}
```

## ✅ Field Mapping Examples

The modules automatically handle field transformations:

| JSON (camelCase) | Terraform (snake_case) |
|------------------|------------------------|
| `isPrivate` | `is_private` |
| `cronSchedule` | `cron_schedule` |
| `keepLastN` | `keep_last_n` |
| `cacheTtlSeconds` | `cache_ttl_seconds` |
| `displayCreatedBy` | `display_created_by` |
| `refreshRate` | `refresh_rate` |
| `resolvedDatasetIds` | `resolved_dataset_ids` |

## 🎯 Output Information

After deployment, you'll get comprehensive output:

```hcl
search_resources_created = {
  saved_queries = [
    { id = "json-config-saved-query", name = "Sample Events Query" },
    { id = "inline-saved-query", name = "Inline Query Example" }
  ]
  dashboards = [
    { id = "json-config-dashboard", name = "Sample Dashboard" },
    { id = "custom-dashboard-id", name = "Sample Dashboard" }
  ]
}
```

## 📚 Module Documentation

### Search Saved Query Module

**Variables**:
- `json_file_path` - Path to JSON file
- `json_content` - Inline JSON content
- `query_id_override` - Override query ID
- `additional_config` - Additional configuration

**Outputs**:
- `query_id` - Created query ID
- `query_name` - Created query name
- `final_config` - Merged configuration (sensitive)

### Search Dashboard Module

**Variables**:
- `json_file_path` - Path to JSON file
- `json_content` - Inline JSON content
- `dashboard_id_override` - Override dashboard ID
- `additional_config` - Additional configuration

**Outputs**:
- `dashboard_id` - Created dashboard ID
- `dashboard_name` - Created dashboard name
- `final_config` - Merged configuration (sensitive)

## 🐛 Troubleshooting

### Common Issues

1. **Missing Required Fields**: Ensure `id`, `name`, and other required fields are present
2. **Complex Elements**: Dashboard elements structure must match the provider schema
3. **Time Format**: Use proper time formats for `earliest`/`latest` fields
4. **Schedule Configuration**: Verify cron expressions and timezone settings

### Debug Mode

```bash
export TF_LOG=DEBUG
terraform plan
```

## 🎉 Success!

Once deployed, you'll have search dashboards and saved queries ready for use in Cribl Stream, all configured directly from your JSON exports!

---

**Note**: This example uses development provider overrides. For production use, specify the appropriate provider version.
