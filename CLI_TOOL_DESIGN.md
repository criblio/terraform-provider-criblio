# Cribl Config to Terraform CLI Tool - Design Document

## Overview

A CLI tool to fetch existing Cribl configuration from `/bulk/diag/download` endpoint and convert YAML configuration files into Terraform modules/resources.

## Recommended Technology Stack

### Programming Language: **Go** ✅

**Why Go?**

1. **Code Reusability**: The provider is already written in Go and has a complete SDK (`internal/sdk`) that can be directly imported
2. **Single Binary Deployment**: Cross-platform static binaries (no runtime dependencies)
3. **Excellent CLI Libraries**: Cobra for CLI structure, Viper for configuration
4. **YAML/JSON Parsing**: Built-in support with `gopkg.in/yaml.v3` and `encoding/json`
5. **Terraform HCL Generation**: Libraries like `github.com/hashicorp/hcl/v2/hclwrite` for generating valid Terraform
6. **Performance**: Fast parsing and file I/O operations
7. **Consistency**: Same language as the Terraform provider means shared types and validation logic

### Alternative Considerations

- **Python**: Good for rapid prototyping but requires runtime, dependency management complexity
- **Node.js/TypeScript**: Not ideal for system tools, runtime overhead
- **Rust**: Excellent but steeper learning curve, less code reuse

## Architecture Design

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Tool (criblio-import)              │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │   Auth       │  │  API Client  │  │  Config Parser   │  │
│  │  Handler     │→ │  (SDK Reuse) │→ │  (YAML/JSON)     │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
│         │                  │                    │            │
│         └──────────────────┴────────────────────┘            │
│                            │                                  │
│                  ┌─────────▼─────────┐                        │
│                  │  Resource Mapper  │                        │
│                  │  (YAML → Terraform)                       │
│                  └─────────┬─────────┘                        │
│                            │                                  │
│         ┌──────────────────┼──────────────────┐              │
│         │                  │                  │              │
│  ┌──────▼──────┐  ┌───────▼──────┐  ┌───────▼──────┐      │
│  │   Source    │  │  Destination │  │     Pack     │      │
│  │  Converter  │  │   Converter  │  │   Converter  │      │
│  └─────────────┘  └──────────────┘  └──────────────┘      │
│         │                  │                  │              │
│         └──────────────────┴──────────────────┘              │
│                            │                                  │
│                  ┌─────────▼─────────┐                        │
│                  │  HCL Generator    │                        │
│                  │  (terraform fmt)  │                        │
│                  └─────────┬─────────┘                        │
│                            │                                  │
│                  ┌─────────▼─────────┐                        │
│                  │   File Writer     │                        │
│                  │  (output/)        │                        │
│                  └───────────────────┘                        │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
tools/criblio-import/
├── cmd/
│   └── criblio-import/
│       └── main.go              # CLI entry point
├── internal/
│   ├── auth/
│   │   └── auth.go              # Authentication (reuse provider auth)
│   ├── client/
│   │   └── client.go            # API client wrapper
│   ├── converter/
│   │   ├── converter.go         # Main conversion orchestrator
│   │   ├── resource_map.go      # Resource type mappings
│   │   └── field_mapper.go      # Field name transformations (camelCase → snake_case)
│   ├── generators/
│   │   ├── source.go            # Generate criblio_source resources
│   │   ├── destination.go       # Generate criblio_destination resources
│   │   ├── pipeline.go          # Generate criblio_pipeline resources
│   │   ├── pack.go              # Generate criblio_pack resources
│   │   ├── routes.go            # Generate criblio_routes resources
│   │   └── common.go            # Shared generation utilities
│   ├── parser/
│   │   ├── yaml.go              # YAML parsing
│   │   ├── archive.go           # Extract tar.gz from API
│   │   └── config.go            # Configuration file structure
│   └── output/
│       ├── writer.go            # File writing utilities
│       └── formatter.go         # Terraform formatting
├── pkg/
│   └── hcl/
│       └── hcl.go               # HCL generation helpers
├── go.mod
├── go.sum
└── README.md
```

## Key Components

### 1. Authentication Module
**Authentication is handled automatically** by `CriblTerraformHook` which is registered in the SDK.

When you create an SDK client with `sdk.New()`, the hook automatically:
- Reads Bearer token from `CRIBL_BEARER_TOKEN`
- Reads OAuth credentials from `CRIBL_CLIENT_ID`/`CRIBL_CLIENT_SECRET`
- Reads credentials file (`~/.cribl/credentials`)
- Handles on-prem authentication (username/password)
- Constructs URLs for workspace and gateway endpoints
- Caches and refreshes OAuth tokens

**No authentication code needed** - just use the SDK!

### 2. API Client
Leverage existing SDK:
```go
import "github.com/criblio/terraform-provider-criblio/internal/sdk"

// Add new method to SDK or create custom endpoint call
client := sdk.New(opts...)
response, err := client.Diag.GetBulkDiagDownload(ctx)
```

### 3. Configuration Parser
- Extract tar.gz archive
- Parse YAML files by resource type
- Map file paths to resource types:
  - `sources/*.yaml` → `criblio_source`
  - `destinations/*.yaml` → `criblio_destination`
  - `routes/*.yaml` → `criblio_routes`
  - `pipelines/*.yaml` → `criblio_pipeline`
  - `packs/*.yaml` → `criblio_pack`

### 4. Field Mapper
Convert Cribl API field names to Terraform attribute names:
- `maxBufferSize` → `max_buffer_size`
- `clientId` → `client_id`
- Nested structures (e.g., `tls`, `pq` objects)

### 5. Resource Converters
Each converter knows how to:
- Read YAML structure
- Map fields to Terraform schema
- Handle type-specific logic
- Generate valid HCL

### 6. HCL Generator
Use `github.com/hashicorp/hcl/v2/hclwrite`:
- Generate `.tf` files with proper formatting
- Support complex nested structures
- Generate variables file for sensitive data (optional)

## CLI Usage Design

```bash
# Basic usage
criblio-import --output ./terraform-configs

# With authentication options
criblio-import \
  --bearer-token $CRIBL_TOKEN \
  --workspace-id main \
  --organization-id my-org \
  --output ./terraform-configs

# Using OAuth
criblio-import \
  --client-id $CLIENT_ID \
  --client-secret $CLIENT_SECRET \
  --workspace-id main \
  --organization-id my-org \
  --output ./terraform-configs

# For on-prem
criblio-import \
  --onprem-server-url http://localhost:9000 \
  --onprem-username admin \
  --onprem-password admin \
  --output ./terraform-configs

# Options
criblio-import \
  --output ./terraform-configs \
  --format resources       # Generate individual resource files
  --format modules         # Generate Terraform modules (default)
  --exclude sources        # Exclude resource types
  --include destinations   # Only include specific types
  --var-file vars.tf      # Generate variables.tf for sensitive data
  --dry-run               # Show what would be generated
```

## Output Structure

### Option 1: Individual Resources (Simple)
```
output/
├── main.tf                 # Provider configuration
├── sources.tf              # All source resources
├── destinations.tf         # All destination resources
├── pipelines.tf            # All pipeline resources
├── routes.tf               # All routes resources
├── packs.tf                # All pack resources
└── variables.tf             # Sensitive variables (optional)
```

### Option 2: Modules Structure (Recommended)
```
output/
├── main.tf                 # Root module with provider
├── modules/
│   ├── sources/
│   │   ├── main.tf         # Source resources
│   │   ├── variables.tf    # Module variables
│   │   └── outputs.tf     # Module outputs
│   ├── destinations/
│   │   └── ...
│   ├── pipelines/
│   │   └── ...
│   └── packs/
│       └── ...
└── variables.tf            # Root variables
```

## Field Mapping Strategy

### 1. Use Existing Resource Schemas
Inspect `internal/provider/*_resource.go` to get exact Terraform attribute names.

### 2. Field Name Conversion
Create a mapping table based on existing modules (see `modules/json-config-*/locals.tf`):
```go
var FieldMappings = map[string]string{
    "maxBufferSize": "max_buffer_size",
    "maxFileSize": "max_file_size",
    "clientId": "client_id",
    "commonNameRegex": "common_name_regex",
    // ... etc
}
```

### 3. Type-Specific Logic
Some resources need special handling:
- **Sources**: Map `type` field to `input_*` attributes based on source type
- **Destinations**: Similar type-based mapping
- **Packs**: Handle pack files, dependencies, versions
- **Routes**: Route expressions and filters

## Implementation Steps

### Phase 1: Foundation (Week 1)
1. Set up Go module structure
2. Implement authentication (reuse provider code)
3. Create API client wrapper for `/bulk/diag/download`
4. Implement archive extraction (tar.gz)

### Phase 2: Parsing & Mapping (Week 2)
1. YAML parser for configuration files
2. Resource type detection from file paths
3. Basic field name mapping (camelCase → snake_case)
4. Schema validation against Terraform resource definitions

### Phase 3: Converters (Week 3-4)
1. Implement converters for each resource type:
   - Sources (most complex, many types)
   - Destinations
   - Pipelines
   - Routes
   - Packs
2. Handle nested structures (TLS, PQ, etc.)
3. Generate basic HCL

### Phase 4: Generation & Output (Week 5)
1. HCL file generation
2. Module structure creation
3. Variable generation for sensitive data
4. Terraform formatting (`terraform fmt` integration)

### Phase 5: Polish & Testing (Week 6)
1. Error handling and logging
2. Progress indicators
3. Dry-run mode
4. Integration tests
5. Documentation

## Dependencies

```go
require (
    github.com/criblio/terraform-provider-criblio v0.0.0
    github.com/spf13/cobra v1.8.0          // CLI framework
    github.com/spf13/viper v1.18.0         // Configuration management
    github.com/hashicorp/hcl/v2 v2.19.0    // HCL generation
    github.com/hashicorp/hcl/v2/hclwrite  // HCL writing
    gopkg.in/yaml.v3 v3.0.1               // YAML parsing
    github.com/fatih/color v1.16.0         // Colored output
)
```

## Deployment Strategy

### Option 1: Standalone Binary (Recommended)
- Build single static binary for each platform
- GitHub Releases: `criblio-import-darwin-amd64`, `criblio-import-linux-amd64`, `criblio-import-windows-amd64.exe`
- Users download and add to PATH

### Option 2: Go Install
```bash
go install github.com/criblio/criblio-import@latest
```

### Option 3: Package Managers
- Homebrew (macOS)
- Chocolatey (Windows)
- AUR (Linux)

### Option 4: Docker Container
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o criblio-import ./cmd/criblio-import

FROM alpine:latest
COPY --from=builder /build/criblio-import /usr/local/bin/
ENTRYPOINT ["criblio-import"]
```

## End User Usage Examples

### Example 1: Basic Import
```bash
# Set up authentication via environment
export CRIBL_CLIENT_ID="..."
export CRIBL_CLIENT_SECRET="..."
export CRIBL_ORGANIZATION_ID="..."
export CRIBL_WORKSPACE_ID="..."

# Run import
criblio-import --output ./my-terraform-configs

# Review generated files
cd my-terraform-configs
terraform init
terraform plan
```

### Example 2: Selective Import
```bash
# Only import sources and destinations
criblio-import \
  --output ./configs \
  --include sources,destinations
```

### Example 3: With Variables
```bash
# Generate with variables for sensitive data
criblio-import \
  --output ./configs \
  --var-file variables.tf \
  --sensitive-fields password,api_key,token
```

### Example 4: Preview Before Import
```bash
# See what would be generated
criblio-import --dry-run --output ./preview
```

## Error Handling & Edge Cases

1. **Missing Fields**: Log warnings, use defaults from schema
2. **Unknown Resource Types**: Skip with warning, log to file
3. **Invalid YAML**: Report file and line number
4. **API Errors**: Retry logic, clear error messages
5. **Name Conflicts**: Sanitize resource names, handle duplicates
6. **Large Configs**: Progress bars, streaming where possible

## Testing Strategy

1. **Unit Tests**: Each converter function
2. **Integration Tests**: End-to-end with mock API responses
3. **Golden Tests**: Compare generated HCL against expected output
4. **Snapshot Tests**: Track changes in generated files

## Future Enhancements

1. **Import State**: Generate `terraform import` commands
2. **Diff Mode**: Compare existing config with generated Terraform
3. **Update Mode**: Update existing Terraform files instead of overwriting
4. **Backup**: Create backup of original configs
5. **Validation**: Run `terraform validate` on generated files
6. **CI/CD Integration**: GitHub Action, GitLab CI template

## Comparison with Alternatives

### Why not use Terraform Import directly?
- Terraform import requires manual `terraform import` for each resource
- This tool automates bulk import
- Generates full configuration, not just state

### Why not use Terraformer?
- Terraformer is generic and may not understand Cribl-specific structures
- This tool can leverage existing provider schemas
- Better handling of nested Cribl configurations

## Open Questions

1. **Endpoint Availability**: Confirm `/bulk/diag/download` exists or needs to be implemented
2. **YAML Structure**: What's the exact structure of config files in the bundle?
3. **Resource Dependencies**: How to handle resource relationships (e.g., pack dependencies)?
4. **Idempotency**: Ensure generated Terraform can be applied multiple times
5. **Versioning**: Track Cribl API version compatibility

## Next Steps

1. ✅ Confirm `/bulk/diag/download` endpoint specification
2. ✅ Get sample YAML config structure from Cribl team
3. ✅ Review existing resource schemas for field mappings
4. ✅ Prototype basic YAML → HCL conversion for one resource type
5. ✅ Get feedback and iterate

