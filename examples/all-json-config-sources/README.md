# All JSON Config Sources Example

This example demonstrates how to use **all available JSON configuration source modules** to create Cribl sources from exported JSON configurations.

## 🎯 **What This Example Shows**

This configuration creates **6 different source types** using JSON files:

| **Source Type** | **Module** | **Port** | **JSON File** |
|-----------------|------------|----------|---------------|
| `cribl_http` | `json-config-source-cribl-http` | 10080 | `cribl_http_source.json` |
| `syslog` | `json-config-source-syslog` | 5514 | `syslog_source.json` |
| `http` | `json-config-source-http` | 8080 | `http_source.json` |
| `tcp` | `json-config-source-tcp` | 9999 | `tcp_source.json` |
| `cribl_tcp` | `json-config-source-cribl-tcp` | 9998 | `cribl_tcp_source.json` |
| `open_telemetry` | `json-config-source-open-telemetry` | 4317 | `opentelemetry_source.json` |

## 📁 **File Structure**

```
examples/all-json-config-sources/
├── main.tf                              # Main Terraform configuration
├── configs/                             # JSON configuration files
│   ├── cribl_http_source.json          # Cribl HTTP source config
│   ├── syslog_source.json              # Syslog source config  
│   ├── http_source.json                # Regular HTTP source config
│   ├── tcp_source.json                 # TCP source config
│   ├── cribl_tcp_source.json           # Cribl TCP source config
│   └── opentelemetry_source.json       # OpenTelemetry source config
└── README.md                           # This file
```

## 🚀 **How to Use**

### **1. Initialize Terraform**
```bash
cd examples/all-json-config-sources
terraform init
```

### **2. Review the Plan**
```bash
terraform plan
```

### **3. Apply the Configuration**
```bash
terraform apply
```

### **4. View Created Sources**
```bash
terraform output all_sources_created
```

## 🔧 **Key Features Demonstrated**

### **✅ JSON File Usage**
Each module uses `json_file_path` to load configuration from JSON files:
```hcl
module "cribl_http_source" {
  source = "../../modules/json-config-source-cribl-http"
  
  group_id       = "default"
  json_file_path = "./configs/cribl_http_source.json"
}
```

### **✅ Automatic Field Transformation**
JSON fields are automatically converted from `camelCase` to `snake_case`:
- `activityLogSampleRate` → `activity_log_sample_rate`
- `maxActiveReq` → `max_active_req`
- `pqEnabled` → `pq_enabled`

### **✅ Nested Object Support**
Nested configurations like `pq` and `tls` are properly transformed:
```json
{
  "pq": {
    "commitFrequency": 5,     // → commit_frequency
    "maxFileSize": "1 MB",    // → max_file_size
    "maxSize": "100 MB"       // → max_size
  }
}
```

### **✅ Complete Validation**
Each module validates that the JSON source type matches the expected type:
- `cribl_http` JSON must have `"type": "cribl_http"`
- `syslog` JSON must have `"type": "syslog"`
- etc.

## 📝 **JSON Configuration Examples**

Each JSON file in the `configs/` directory represents a real-world source configuration that could be exported from Cribl Stream UI. These files demonstrate:

- **Complete configurations** with all necessary fields
- **Realistic values** for production use
- **Nested objects** (pq, tls configurations)
- **Arrays** (connections, streamtags, metadata)
- **Source-specific fields** for each type

## 🎯 **Real-World Usage**

To use this in your environment:

1. **Export your sources** from Cribl Stream UI as JSON
2. **Place the JSON files** in the `configs/` directory
3. **Update the provider** configuration with your Cribl details
4. **Run `terraform apply`** to create the sources

## 📊 **Expected Output**

After applying, you'll see output similar to:
```hcl
all_sources_created = {
  "group_id" = "default"
  "sources" = [
    {
      "id" = "json-config-cribl-http"
      "type" = "cribl_http"
    },
    {
      "id" = "json-config-syslog"
      "type" = "syslog"
    },
    # ... and 4 more sources
  ]
}
```

## 🔗 **Related Examples**

- `../json-config-source/` - Individual source examples with overrides
- `../stream-syslog-to-lake/` - Complete pipeline example

This example shows how the JSON configuration modules can **scale to handle multiple source types** while maintaining **consistency and simplicity** in your Terraform configurations.
