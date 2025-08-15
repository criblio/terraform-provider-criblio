# Cribl Source Module

## Overview
This is a Terraform module for creating Cribl sources. It abstracts the complexity of different source types and provides a unified interface with sensible defaults.

## Requirements
- Terraform >= 1.0
- Cribl provider (`criblio/criblio`)

## Module Location
- Path: `./modules/cribl-source`
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

## Input Variables

### Required Variables
| Name | Description | Type |
|------|-------------|------|
| `source_id` | Unique identifier for the source | string |
| `group_id` | Worker group ID | string |
| `source_type` | Type of source (see supported types) | string |
| `port` | Port number (cannot be default ports: 9514, 10200, 10001, 10080, 10060, 4317, 4318) | number |

### Optional Variables
| Name | Description | Type | Default |
|------|-------------|------|---------|
| `description` | Source description | string | "" |
| `disabled` | Disable source | bool | false |
| `connections` | Output connections array | list(object) | [] |
| `pipeline` | Default pipeline | string | null |
| `pq_enabled` | Enable persistent queue | bool | false |
| `streamtags` | Stream tags array | list(string) | [] |
| `custom_config` | Override any default settings | any | {} |

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