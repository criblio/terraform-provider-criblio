# Cribl Source Module - Claude Documentation

## Overview
This is a Terraform module for creating Cribl sources. It abstracts the complexity of different source types and provides a unified interface with sensible defaults.

## Module Location
- Path: `./modules/cribl-source` or `../source-module`
- Provider: `criblio/criblio`

## Supported Source Types
- `syslog` - Syslog input (TCP/UDP on same port)
- `cribl_http` - Cribl HTTP receiver
- `http` - Regular HTTP source
- `tcp` - Regular TCP source
- `cribl_tcp` - Cribl TCP receiver
- `otlp` - OpenTelemetry (internally maps to "open_telemetry")

## Key Features
1. **Unified Interface**: Single module handles all source types
2. **Smart Defaults**: Each source type has appropriate defaults
3. **Port Validation**: Prevents using default Cribl ports
4. **Flexible Configuration**: `custom_config` allows overriding any setting

## Required Variables
- `source_id` - Unique identifier for the source
- `group_id` - Worker group ID
- `source_type` - Type of source (see supported types)
- `port` - Port number (cannot be default ports: 9514, 10200, 10001, 10080, 10060, 4317, 4318)

## Optional Variables
- `description` - Source description (default: "")
- `disabled` - Disable source (default: false)
- `connections` - Output connections array (default: [])
- `pipeline` - Default pipeline (default: null)
- `pq_enabled` - Enable persistent queue (default: false)
- `streamtags` - Stream tags array (default: [])
- `custom_config` - Override any default settings (default: {})

## Usage Examples

### Simple Syslog Source
```hcl
module "syslog" {
  source = "./modules/cribl-source"
  
  source_id   = "firewall-logs"
  group_id    = "default"
  source_type = "syslog"
  port        = 20005
}