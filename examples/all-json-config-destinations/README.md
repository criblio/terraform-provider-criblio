# Complete JSON Config Destinations Example

This comprehensive example demonstrates **all 13 JSON-config destination modules** available in the Terraform provider. You can copy JSON configurations directly from Cribl Stream UI and use them with these modules.

## 🚀 Features

- **Complete Coverage**: All destination types from the provider
- **Direct JSON Import**: Copy configurations straight from Cribl UI
- **Field Mapping**: Automatic `camelCase` to `snake_case` conversion
- **Unknown Field Filtering**: Ignores non-schema fields automatically
- **Production Ready**: All modules validated against provider schemas

## 📦 Included Destination Modules

### Core Data Destinations
1. **Cribl Lake** - `json-config-destination-cribl-lake`
2. **Amazon S3** - `json-config-destination-s3`
3. **Elasticsearch** - `json-config-destination-elasticsearch`

### Message Queue & Streaming
4. **Apache Kafka** - `json-config-destination-kafka`
5. **AWS Kinesis** - `json-config-destination-kinesis`

### HTTP & API Endpoints
6. **Webhook** - `json-config-destination-webhook`
7. **Splunk HEC** - `json-config-destination-splunk-hec`

### Cribl-Specific
8. **Cribl HTTP** - `json-config-destination-cribl-http`
9. **Cribl TCP** - `json-config-destination-cribl-tcp`

### Network & Protocol
10. **Syslog** - `json-config-destination-syslog`

### Monitoring & Observability
11. **OpenTelemetry** - `json-config-destination-open-telemetry`
12. **Prometheus** - `json-config-destination-prometheus`

### Security
13. **CrowdStrike Next-Gen SIEM** - `json-config-destination-crowdstrike-next-gen-siem`

## 📄 Configuration Files

Each destination includes a complete example JSON configuration in the `configs/` directory:

```
configs/
├── cribl_lake_destination.json
├── s3_destination.json
├── elasticsearch_destination.json
├── kafka_destination.json
├── kinesis_destination.json
├── webhook_destination.json
├── splunk_hec_destination.json
├── cribl_http_destination.json
├── cribl_tcp_destination.json
├── syslog_destination.json
├── open_telemetry_destination.json
├── prometheus_destination.json
└── crowdstrike_destination.json
```

## 🛠 Usage

### Deploy All Destinations

```bash
# Initialize
terraform init

# Validate configuration
terraform validate

# Plan deployment
terraform plan

# Deploy all 13 destinations
terraform apply
```

### Using Individual Modules

You can use any module individually:

```hcl
module "my_s3_destination" {
  source = "../../modules/json-config-destination-s3"

  group_id       = "default"
  json_file_path = "./my-s3-config.json"
}
```

### Using Inline JSON

```hcl
module "my_kafka_destination" {
  source = "../../modules/json-config-destination-kafka"

  group_id     = "default"
  json_content = jsonencode({
    id      = "my-kafka-dest"
    type    = "kafka"
    brokers = ["kafka:9092"]
    topic   = "events"
  })
}
```

### With Configuration Overrides

```hcl
module "custom_webhook" {
  source = "../../modules/json-config-destination-webhook"

  group_id              = "default"
  json_file_path        = "./webhook.json"
  destination_id_override = "custom-webhook-id"
  
  additional_config = {
    description = "Custom webhook with overrides"
    concurrency = 10
  }
}
```

## ✅ Field Mapping Examples

All modules automatically handle field transformations:

| JSON (camelCase) | Terraform (snake_case) |
|------------------|------------------------|
| `flushPeriodSec` | `flush_period_sec` |
| `maxPayloadSizeKb` | `max_payload_size_kb` |
| `onBackpressure` | `on_backpressure` |
| `rejectUnauthorized` | `reject_unauthorized` |
| `dnsResolvePeriodSec` | `dns_resolve_period_sec` |
| `responseRetrySettings` | `response_retry_settings` |

## 🎯 Output Information

After deployment, you'll get comprehensive output showing all created destinations:

```hcl
all_destinations_created = {
  group_id = "default"
  destinations = [
    { id = "cribl-lake-dest", type = "cribl_lake" },
    { id = "s3-dest", type = "s3" },
    { id = "kafka-dest", type = "kafka" },
    # ... all 13 destinations
  ]
}
```

## 🔧 Customization

### Adding Your Own JSON Configs

1. Export JSON from Cribl Stream UI
2. Save to `configs/my-destination.json`
3. Update `main.tf` to reference your file
4. Run `terraform plan` to verify

### Module-Specific Features

Each module supports:
- ✅ **File or inline JSON input**
- ✅ **ID overrides**
- ✅ **Additional configuration merging**
- ✅ **Type validation**
- ✅ **Schema compliance**

## 📚 Documentation

For detailed documentation on each module, see:
- Individual module README files in `modules/json-config-destination-*/README.md`
- [Cribl Terraform Provider Documentation](https://registry.terraform.io/providers/criblio/criblio/latest/docs)

## 🐛 Troubleshooting

### Common Issues

1. **Schema Validation Errors**: Ensure your JSON matches the expected schema
2. **Missing Fields**: Check that required fields are present in your JSON
3. **Type Mismatches**: Verify the `type` field matches the module

### Debug Mode

```bash
# Enable detailed logging
export TF_LOG=DEBUG
terraform plan
```

## 🎉 Success!

Once deployed, you'll have **13 different destination types** ready to receive data from Cribl Stream, all configured directly from your JSON exports!

---

**Note**: This example uses development provider overrides. For production use, specify the appropriate provider version in your `terraform` block.
