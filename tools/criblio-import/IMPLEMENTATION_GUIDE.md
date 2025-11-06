# Implementation Guide for criblio-import CLI Tool

This guide provides step-by-step instructions for implementing the CLI tool based on the design document.

## Phase 1: Foundation Setup

### Step 1.1: Initialize the Project
```bash
cd tools/criblio-import
go mod init github.com/criblio/terraform-provider-criblio/tools/criblio-import
go mod edit -replace github.com/criblio/terraform-provider-criblio=../../
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get gopkg.in/yaml.v3@latest
go get github.com/hashicorp/hcl/v2@latest
go get github.com/fatih/color@latest
```

### Step 1.2: SDK Client Initialization (Authentication Handled Automatically)
**Important:** Authentication is handled automatically by `CriblTerraformHook` which is registered in the SDK.

Create `internal/client/client.go`:
- Initialize SDK with options
- The hook automatically handles:
  - Bearer token authentication (`CRIBL_BEARER_TOKEN`)
  - OAuth authentication (client_id/client_secret)
  - On-prem authentication (username/password)
  - Credentials file reading (`~/.cribl/credentials`)
  - URL construction for cloud and on-prem

```go
package client

import (
    "context"
    "github.com/criblio/terraform-provider-criblio/internal/sdk"
    "github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
)

// NewClient creates a new SDK client with authentication handled automatically
// The CriblTerraformHook (registered in internal/sdk/internal/hooks/registration.go)
// handles all authentication methods based on environment variables or credentials file
func NewClient(ctx context.Context, opts ...sdk.SDKOption) (*sdk.CriblIo, error) {
    // Create security source from environment/config
    security := shared.Security{
        // Hook will read from environment variables or credentials file
        // No need to manually set these - hook handles it
    }
    
    // Initialize SDK - hook automatically handles auth
    client := sdk.New(
        sdk.WithSecurity(security),
        // Add other options as needed
    )
    
    return client, nil
}
```

**Note:** The `CriblTerraformHook` is automatically registered and handles:
- Reading `CRIBL_BEARER_TOKEN`, `CRIBL_CLIENT_ID`, `CRIBL_CLIENT_SECRET`, etc.
- Reading from `~/.cribl/credentials` file
- OAuth token acquisition and caching
- On-prem authentication via username/password
- URL construction for workspace and gateway endpoints

### Step 1.3: Implement API Client Wrapper
**Note:** Since authentication is handled by the hook, we just need to use the SDK directly.
Create `internal/client/client.go`:
- Wrap SDK client
- Add method for `/bulk/diag/download` endpoint
- Handle archive download and extraction

```go
func (c *Client) DownloadBulkDiag(ctx context.Context) ([]byte, error) {
    // Call API endpoint
    // Return tar.gz bytes
}

func (c *Client) ExtractArchive(data []byte, outputDir string) error {
    // Extract tar.gz to output directory
    // Return list of extracted files
}
```

### Step 1.4: Implement Archive Extraction
Create `internal/parser/archive.go`:
- Use `archive/tar` and `compress/gzip`
- Extract files to temporary directory
- Organize by resource type (sources/, destinations/, etc.)

## Phase 2: Parsing & Mapping

### Step 2.1: Implement YAML Parser
Create `internal/parser/yaml.go`:
- Parse individual YAML files
- Detect resource type from file path
- Return structured configuration maps

### Step 2.2: Resource Type Detection
Map file paths to resource types:
- `sources/*.yaml` → `criblio_source`
- `destinations/*.yaml` → `criblio_destination`
- `routes/*.yaml` → `criblio_routes`
- `pipelines/*.yaml` → `criblio_pipeline`
- `packs/*.yaml` → `criblio_pack`

### Step 2.3: Field Mapper Implementation
The `field_mapper.go` is already created - enhance it:
- Add more field mappings from existing modules
- Handle edge cases (acronyms, special names)
- Support nested field mapping

## Phase 3: Converters

### Step 3.1: Source Converter
Create `internal/generators/source.go`:
- Implement `ConvertSource()` method
- Handle all source types (http, tcp, syslog, etc.)
- Map type-specific input attributes
- Handle nested structures (TLS, PQ, etc.)

Reference:
- `internal/provider/source_resource.go` for schema
- `modules/json-config-source-*/locals.tf` for field mappings

### Step 3.2: Destination Converter
Create `internal/generators/destination.go`:
- Similar pattern to source converter
- Handle destination-specific types (s3, kafka, splunk, etc.)

### Step 3.3: Pipeline Converter
Create `internal/generators/pipeline.go`:
- Handle pipeline expressions
- Preserve function calls and expressions
- May need special handling for expressions vs. literal values

### Step 3.4: Routes Converter
Create `internal/generators/routes.go`:
- Handle route filters and expressions
- Map route rules to Terraform format

### Step 3.5: Pack Converter
Create `internal/generators/pack.go`:
- Handle pack metadata
- May need to handle pack files separately
- Map dependencies and versions

## Phase 4: HCL Generation & Output

### Step 4.1: HCL Writer Utilities
Create `internal/output/writer.go`:
- Write HCL files using `hclwrite`
- Organize resources into modules
- Generate `main.tf`, `variables.tf`, `outputs.tf`

### Step 4.2: Formatter
Create `internal/output/formatter.go`:
- Run `terraform fmt` on generated files (optional)
- Validate generated HCL syntax
- Provide formatting utilities

### Step 4.3: Module Structure Generator
Create `internal/output/module_generator.go`:
- Generate module structure based on format option
- Create module directories
- Generate module main.tf files
- Generate root main.tf with provider config

## Phase 5: Integration & CLI

### Step 5.1: Complete Main Command
Complete `cmd/criblio-import/main.go`:
- Wire up all components
- Handle CLI flags
- Implement dry-run mode
- Add progress indicators

### Step 5.2: Error Handling
- Graceful error handling
- User-friendly error messages
- Logging utilities
- Verbose mode

### Step 5.3: Testing
Create tests:
- Unit tests for converters
- Integration tests with mock data
- Golden file tests for HCL output

## Implementation Tips

### Reusing Provider Code
1. Import provider's SDK: `github.com/criblio/terraform-provider-criblio/internal/sdk`
2. Import provider types: `github.com/criblio/terraform-provider-criblio/internal/provider`
3. Reuse authentication logic from provider

### Field Mapping Reference
Check existing modules for field mappings:
- `modules/json-config-source-http/locals.tf`
- `modules/json-config-destination-s3/locals.tf`
- Similar patterns for all resource types

### Schema Reference
Use Terraform provider resource schemas:
- `internal/provider/*_resource.go` - schema definitions
- `docs/resources/*.md` - documentation with attribute names

### Handling Complex Types
- JSON fields: May need to preserve as JSON strings
- Expressions: Preserve as-is in Terraform
- Nested objects: Create nested blocks in HCL
- Lists: Convert to HCL list syntax

### Testing Strategy
1. Start with one resource type (e.g., HTTP source)
2. Test with sample YAML from actual Cribl instance
3. Compare generated Terraform with manual creation
4. Iterate and refine

## Sample YAML Structure (Expected)

```yaml
# sources/my-http-source.yaml
id: my-http-source
type: http
group: default
port: 8080
host: 0.0.0.0
maxBufferSize: 1024
tls:
  certPath: /path/to/cert
  rejectUnauthorized: true
```

Should convert to:

```hcl
resource "criblio_source" "my_http_source" {
  id       = "my-http-source"
  group_id = "default"

  input_http {
    port            = 8080
    host            = "0.0.0.0"
    max_buffer_size = 1024

    tls {
      cert_path           = "/path/to/cert"
      reject_unauthorized = true
    }
  }
}
```

## Next Steps

1. **Confirm API Endpoint**: Verify `/bulk/diag/download` endpoint exists and understand response format
2. **Get Sample Data**: Obtain sample YAML configs from Cribl instance
3. **Prototype**: Start with one resource type (HTTP source) to validate approach
4. **Iterate**: Expand to other resource types based on prototype
5. **Test**: Test with real Cribl instances
6. **Document**: Create user documentation and examples

## Questions to Resolve

1. What's the exact structure of the `/bulk/diag/download` response?
2. Are configs in YAML or JSON format?
3. What's the directory structure inside the archive?
4. Are there resource dependencies that need ordering?
5. How to handle sensitive fields (passwords, tokens)?
6. Should we generate `terraform import` commands as well?

