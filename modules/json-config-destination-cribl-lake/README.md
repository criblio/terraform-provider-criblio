# JSON Config Destination - Cribl Lake

This Terraform module creates a Cribl Lake destination from JSON configuration, supporting both file-based and inline JSON.

## Features

- ✅ **JSON to Terraform**: Convert exported Cribl Lake destination JSON to Terraform resources
- ✅ **Field Mapping**: Automatic `camelCase` to `snake_case` field conversion
- ✅ **Flexible Input**: Support JSON files or inline JSON strings
- ✅ **Type Validation**: Ensures JSON contains `cribl_lake` destination type
- ✅ **ID Override**: Optional destination ID customization
- ✅ **Additional Config**: Merge extra configuration with JSON

## Usage

### From JSON File

```hcl
module "cribl_lake_destination" {
  source = "./modules/json-config-destination-cribl-lake"
  
  group_id       = "default"
  json_file_path = "./config/lake_destination.json"
}
```

### From Inline JSON

```hcl
module "cribl_lake_destination" {
  source = "./modules/json-config-destination-cribl-lake"
  
  group_id     = "default"
  json_content = jsonencode({
    id       = "my-lake-dest"
    type     = "cribl_lake"
    destPath = "my_dataset"
  })
}
```

### With Overrides

```hcl
module "cribl_lake_destination" {
  source = "./modules/json-config-destination-cribl-lake"
  
  group_id                 = "default" 
  json_file_path          = "./config/lake_destination.json"
  destination_id_override = "custom-lake-dest"
  
  additional_config = {
    description = "Custom Cribl Lake destination"
  }
}
```

## JSON Example

```json
{
  "id": "lake-destination",
  "type": "cribl_lake", 
  "description": "Send data to Cribl Lake",
  "destPath": "my_dataset"
}
```

## Field Mappings

This module automatically converts camelCase JSON fields to snake_case Terraform attributes:

| JSON Field | Terraform Attribute |
|------------|-------------------|
| `destPath` | `dest_path`       |

## Requirements

- Terraform >= 0.14
- Criblio Provider >= 1.0.0

## Outputs

- `destination_id`: The ID of the created destination
- `destination_type`: The type of the destination ("cribl_lake")
