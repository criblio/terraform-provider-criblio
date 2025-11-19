# Pack with Full Configuration Example

This example demonstrates how to create a **Cribl Pack with a comprehensive pipeline configuration** showcasing various pipeline functions and their configurations.

## 🎯 **What This Example Shows**

This configuration creates:

1. **A Cribl Pack** (`pack-with-full-config`)
2. **A Pipeline** with multiple functions demonstrating:
   - **Eval functions** for adding and removing fields
   - **Serde function** for extracting key-value pairs
   - **Serialize function** for re-serializing data (disabled by default)

## 📋 **Pipeline Functions Demonstrated**

| **Function** | **Purpose** | **Description** |
|--------------|-------------|-----------------|
| `eval` | Field Addition | Exposes metadata environments (domain, organizationId, accountId, etc.) |
| `serde` | Data Extraction | Extracts `_raw` data to key-value pairs |
| `serialize` | Data Serialization | Re-serializes events (disabled in this example) |
| `eval` | Field Removal | Removes the `_raw` field after processing |

## 📁 **File Structure**

```
examples/pack-with-full-config/
├── main.tf                  # Main Terraform configuration
└── README.md               # This file
```

## 🚀 **How to Use**

### **1. Update Provider Configuration**

Edit the provider configuration in `main.tf` to match your Cribl environment:

```hcl
provider "criblio" {
  organization_id = "your-org-id"
  workspace_id    = "your-workspace"
  cloud_domain    = "your-domain.cloud"
}
```

### **2. Initialize Terraform**

```bash
cd examples/pack-with-full-config
terraform init
```

### **3. Review the Plan**

```bash
terraform plan
```

### **4. Apply the Configuration**

```bash
terraform apply
```

## 🔧 **Key Features Demonstrated**

### **✅ Complex Pipeline Configuration**

The pipeline demonstrates how to use `jsonencode()` for complex nested configurations:

```hcl
conf = jsonencode({
  add = [
    {
      name = "domain"
      value = "__metadata.env.CRIBL_CLOUD_DOMAIN"
      disabled = false
    },
    # ... more fields
  ]
})
```

### **✅ Multiple Function Types**

Shows various function types in sequence:
- **Eval**: Add computed fields from metadata
- **Serde**: Parse structured data
- **Serialize**: Re-format data (optional)
- **Eval**: Clean up fields

### **✅ Metadata Access**

Demonstrates accessing various metadata sources:
- Cloud environment variables (`CRIBL_CLOUD_DOMAIN`)
- Cribl metadata (`__metadata.cribl.group`)
- AWS metadata (`__metadata.aws.identity.instanceId`)

### **✅ Disabled Functions**

Shows how to include functions in configuration but keep them disabled:

```hcl
{
  id       = "serialize"
  filter   = "true"
  disabled = true
  # ...
}
```

## 📝 **Configuration Details**

### **Pack Configuration**

```hcl
resource "criblio_pack" "full_config_pack" {
  id           = "pack-with-full-config"
  group_id     = "default"
  description  = "Pack with full pipeline configuration"
  disabled     = false
  display_name = "Pack with Full Config"
  version      = "1.0.0"
}
```

### **Pipeline Configuration**

The pipeline is configured with:
- Empty streamtags array
- Four functions in sequence
- Each function with appropriate filters and configurations

## 🎯 **Real-World Usage**

This example is useful for:

1. **Understanding complex pipeline configurations** in Terraform
2. **Learning how to use different function types** together
3. **Managing metadata enrichment** in pipelines
4. **Structuring pipeline functions** for auditd or similar log processing

## 🔗 **Related Examples**

- `../pack-with-pipeline/` - Basic pipeline configuration
- `../pack-with-destination/` - Pack with destination configuration
- `../pack-with-source/` - Pack with source configuration

## 📊 **Expected Behavior**

After applying, the pack and pipeline will be created in your Cribl environment. You can verify:

```bash
terraform show
```

The pipeline will process events by:
1. Adding metadata fields (domain, organizationId, etc.)
2. Parsing the `_raw` field into key-value pairs
3. Removing the original `_raw` field

This creates a structured event with all relevant metadata for downstream processing.

