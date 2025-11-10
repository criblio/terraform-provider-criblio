# [ERD] Cribl Configuration Import CLI Tool

**By** Infrastructure Team

**2 min**

**0**

**Add a reaction**

## PRD

CLI Tool for Importing Cribl Configuration to Terraform

## Epic

INFRA-7290: Cribl Config to Terraform Import Tool

**IN PROGRESS**

## Initiative Ticket (TCTI)

TCTI-XXX: Infrastructure as Code Migration

**IN PROGRESS**

## Slack

https://cribl.enterprise.slack.com/archives/[CHANNEL_ID]

Connect your Slack account

## Customer Use Case

Many Cribl customers have been managing their Cribl configurations through the UI for an extended period, building up complex configurations with hundreds of resources including sources, destinations, pipelines, routes, packs, and global variables. With the recent availability of the Terraform provider for Cribl, these customers now want to migrate from UI-managed configurations to Infrastructure as Code (IaC) to gain benefits such as:

- Version control for all configuration changes
- Consistency across environments (dev, staging, production)
- GitOps workflows and automated deployments
- Configuration drift detection and remediation
- Disaster recovery through code-based infrastructure

However, manually recreating hundreds of existing UI-configured resources in Terraform would be time-consuming, error-prone, and impractical. The `criblio-import` CLI tool solves this by providing a **bulk export and migration solution** that:

- Downloads all existing Cribl configuration from a running instance via the `/bulk/diag/download` API endpoint (which exports the complete configuration bundle)
- Automatically converts UI-configured YAML/JSON resources to Terraform HCL format
- Generates organized Terraform modules ready for version control
- Supports automatic import of all exported resources into Terraform state
- Enables seamless one-time migration from UI-managed configurations to IaC-managed infrastructure

This allows customers to quickly bring their existing UI-configured infrastructure under Terraform management without manual recreation, enabling them to adopt Infrastructure as Code practices going forward.

## MVP Requirements

### Core Functionality

**Configuration Download**
- Download configuration bundle from Cribl API `/bulk/diag/download` endpoint
- Support for both Cribl Cloud and on-premises deployments
- Handle archive extraction (.tar.gz format)
- Parse YAML configuration files from the diag bundle
- **Important Limitation**: The `/bulk/diag/download` endpoint does not export secrets or sensitive authentication credentials (auth tokens, API keys, passwords, etc.). Users will need to manually add these sensitive values to the generated Terraform configuration after import. This affects:
  - Source authentication tokens (e.g., Splunk HEC tokens, Datadog API keys)
  - Destination authentication credentials (e.g., AWS access keys, OAuth tokens)
  - Collector authentication settings
  - Any other sensitive fields that are redacted in the diag bundle
- The tool will generate placeholder comments or variable references for sensitive fields to guide users on what needs to be populated

**Resource Conversion**
- Convert Cribl YAML configurations to Terraform HCL resource definitions
- Support all major resource types:
  - Sources (criblio_source)
  - Destinations (criblio_destination)
  - Pipelines (criblio_pipeline)
  - Routes (criblio_routes)
  - Packs (criblio_pack)
  - Global Variables (criblio_global_var)
  - Groups (criblio_group)
  - And other supported resource types
- **Reuse SDK Conversion Methods**: Instead of manual field mapping, the CLI tool leverages the provider's `RefreshFrom*` methods (e.g., `RefreshFromSharedGlobalVar`, `RefreshFromSharedSource`) which automatically handle the conversion from API camelCase fields to Terraform snake_case fields. This ensures:
  - **Exact field name matching**: Uses the same conversion logic as the provider
  - **No manual mapping maintenance**: Field mappings are defined once in the provider's ResourceModel structs with `tfsdk` tags
  - **Type safety**: Works with generated SDK types, ensuring consistency
  - **Automatic updates**: When provider schemas change, conversions automatically stay in sync
- Handle nested structures and complex types using the same SDK conversion patterns
- Preserve resource relationships and dependencies

**Output Generation**
- Generate organized Terraform module structure (default)
- Support flat resource file structure (optional)
- Create provider configuration (main.tf)
- Generate variables.tf for sensitive fields (optional) - **Note**: Since the diag bundle doesn't export secrets, this helps users identify which fields need manual population
- Include README.md with usage instructions and warnings about missing secrets
- Add placeholder comments in generated HCL for fields that require manual secret input (e.g., `# TODO: Add authentication token - not exported by diag bundle`)

### Authentication and Configuration

**Multi-Source Configuration**
- Support command-line flags (highest priority)
- Read from environment variables (CRIBL_*)
- Load from credentials file (~/.cribl/credentials)
- Priority resolution: Flags > Env > File > Defaults

**Technology Stack and Repository Structure**

The CLI tool is implemented in **Go (Golang)** and is part of the same repository as the Terraform provider (`terraform-provider-criblio`). This design decision was made for several critical reasons:

**Why Go?**
- **Maximum Code Reuse**: The Terraform provider is already written in Go with a complete SDK (`internal/sdk`) that can be directly imported. This eliminates the need to reimplement API clients, authentication logic, or data models.
- **Shared Authentication**: The provider's `CriblTerraformHook` automatically handles all authentication methods (Bearer token, OAuth, on-prem, credentials file) - the CLI tool reuses this without any additional code.
- **Type Safety**: The SDK is generated by Speakeasy from the OpenAPI schema, providing compile-time type safety and ensuring the CLI tool always uses the same types as the provider.
- **Single Binary Deployment**: Go compiles to a single static binary with no runtime dependencies, making distribution simple (GitHub Releases) and installation straightforward for end users.
- **Performance**: Fast execution for parsing large configuration bundles and converting hundreds of resources efficiently.

**Why Same Repository?**
- **Direct SDK Import**: The CLI tool can directly import `github.com/criblio/terraform-provider-criblio/internal/sdk` without any code duplication or translation layer.
- **Automatic Updates**: When the OpenAPI schema updates and the SDK regenerates, the CLI tool automatically benefits from new API endpoints and types without manual updates.
- **Schema Alignment**: The CLI tool can reference provider schemas (`internal/provider/*_resource.go`) for field mappings, validation, and ensuring generated Terraform matches provider expectations.
- **Simplified Development**: Single repository means single PR, single CI pipeline, and easier coordination between provider and tool development.
- **Consistent Stack**: Same language, same tooling, same patterns - easier for the team to maintain and extend.

This monorepo approach maximizes code reuse, ensures consistency, and simplifies long-term maintenance while providing the fastest path to a production-ready CLI tool.

**Flexible Filtering**
- Include specific resource types (--include)
- Exclude specific resource types (--exclude)
- Preview mode (--dry-run) to see what would be generated

### Terraform Integration

**Auto-Import Capability**
- Automatically import generated resources into Terraform state (--auto-import)
- Support for automatic apply after import (--auto-apply)
- Use terraform-exec library for programmatic Terraform control
- Handle import ID formatting for each resource type

**State Management**
- Refresh state after import to populate all fields
- Generate configuration that matches state (terraform plan -generate-config-out)
- Handle computed fields and state differences

### Code Reusability

**SDK Integration**
- Direct import of generated SDK (`internal/sdk`)
- Reuse authentication logic (CriblTerraformHook)
- Use shared types and models from provider
- Reference provider schemas for validation
- No separate API client maintenance required

**Type Safety**
- Use generated types from Speakeasy SDK
- Compile-time type checking
- Consistent types across provider and CLI tool

## System Design

The CLI tool is built in Go (Golang) to maximize code reuse with the existing Terraform provider. The provider uses Speakeasy to generate both the SDK and Terraform provider code from the OpenAPI schema (`openapi.yml`). The CLI tool directly imports the generated SDK, eliminating the need for a separate API client.

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    CODE GENERATION (Build Time)              │
└─────────────────────────────────────────────────────────────┘

    openapi.yml ──┐
                  │
    overlay.yml ──┼──> [Speakeasy] ──> internal/sdk/ (API Client)
                  │                    internal/provider/ (Terraform)
                  │
    gen.yaml ─────┘

┌─────────────────────────────────────────────────────────────┐
│                    CLI TOOL EXECUTION (Runtime)              │
└─────────────────────────────────────────────────────────────┘

    User Command
         │
         ├─> [Cobra CLI] ──> Parse flags/args
         │
         ├─> [Viper] ──> Load config (Flags > Env > File)
         │
         ├─> [CriblTerraformHook] ──> Auto-authenticate
         │
         ├─> [SDK Client] ──> Call API: /bulk/diag/download
         │
         ├─> [Archive Handler] ──> Extract .tar.gz
         │
         ├─> [YAML Parser] ──> Parse config files
         │
         ├─> [Resource Converters] ──> YAML → HCL
         │
         ├─> [HCL Generator] ──> Format with hclwrite
         │
         └─> [File Writer] ──> Write .tf files
                │
                └─> [Optional: Auto-Import] ──> terraform import
```

### Component Structure

**CLI Framework (Cobra)**
- Command structure: `criblio-import [command] [flags]`
- Subcommands: `import`, `version`
- Flag management with Viper integration
- Help generation and shell completion support

**Configuration Management (Viper)**
- Multi-source configuration loading
- Environment variable support (CRIBL_* prefix)
- Config file support (~/.cribl/credentials.ini)
- Flag binding and precedence resolution

**Authentication Layer (CriblTerraformHook)**
- Automatic authentication method detection
- Support for Bearer token, OAuth, on-prem, credentials file
- URL construction for Cloud and on-prem deployments
- Reused from provider SDK (zero additional code)

**API Client (internal/sdk)**
- Generated by Speakeasy from OpenAPI schema
- All API endpoints available
- Type-safe operations and responses
- Automatic authentication via hook

**Resource Converters**
- Modular converter per resource type
- Field mapping (camelCase → snake_case)
- Type conversion and validation
- Schema-aware generation

**HCL Generation**
- Use `github.com/hashicorp/hcl/v2/hclwrite`
- Generate properly formatted Terraform files
- Support complex nested structures
- Optional formatting with `terraform fmt`

**Terraform Integration (terraform-exec)**
- Programmatic Terraform command execution
- Support for init, import, validate, plan, apply, refresh
- State management and validation

### Output Structure

```
output-directory/
├── main.tf                 # Provider configuration
├── variables.tf            # Input variables (if --var-file specified)
├── README.md               # Generated documentation
└── modules/
    ├── sources/
    │   ├── main.tf         # All source resources
    │   ├── variables.tf    # Module variables
    │   └── outputs.tf      # Module outputs
    ├── destinations/
    │   └── ...
    ├── pipelines/
    │   └── ...
    ├── routes/
    │   └── ...
    └── packs/
        └── ...
```

### Deployment

**Distribution Methods**
- GitHub Releases (primary): Pre-built binaries for Linux, macOS, Windows
- Go install: `go install github.com/criblio/terraform-provider-criblio/tools/criblio-import/cmd/criblio-import@latest`

**Binary Characteristics**
- Single static binary (no runtime dependencies)
- Cross-platform support (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
- Small footprint (~10-50MB)
- Fast startup and execution

## Further Considerations

### Future Enhancements

**Incremental Import**
- Support for importing only changed resources
- Compare existing Terraform state with current Cribl configuration
- Generate only delta changes

**Resource Dependency Resolution**
- Automatically detect and order resource dependencies
- Generate import blocks with proper dependency ordering
- Handle circular dependencies gracefully

**Validation and Testing**
- Validate generated Terraform before writing files
- Run `terraform validate` automatically
- Generate test files for imported resources

**Advanced Filtering**
- Filter by tags or metadata
- Filter by resource attributes
- Support for complex filter expressions

**Configuration Templates**
- Support for custom output templates
- Configurable module structure
- Custom variable generation strategies

**Additional Distribution Methods**
- Homebrew tap for macOS users (familiar package manager experience, easy updates)
- Docker container for CI/CD environments (no local installation, consistent environment, CI/CD friendly)

### Operational Considerations

**Error Handling**
- Clear error messages with resolution steps
- Progress indicators for long-running operations
- Detailed logging in verbose mode
- Graceful handling of partial failures

**Performance**
- Concurrent resource conversion (goroutines)
- Streaming archive extraction for large bundles
- Efficient YAML parsing
- Optimized HCL generation

**Security and Secrets Management**

**Secrets Export Limitation**
- The `/bulk/diag/download` endpoint does not export secrets or sensitive authentication credentials for security reasons
- Users must manually add authentication tokens and credentials after import
- Affected resource types:
  - **Sources**: Auth tokens (Splunk HEC tokens, Datadog API keys, etc.)
  - **Destinations**: API keys, OAuth tokens, AWS credentials, etc.
  - **Collectors**: Authentication settings and credentials
  - **Any resource** with sensitive fields (passwords, tokens, keys)

**Secrets Handling Strategy**
- Generate Terraform configuration with placeholder comments indicating missing secrets
- Option to generate `variables.tf` with variable declarations for sensitive fields
- Use Terraform variables (e.g., `var.splunk_hec_token`) instead of hardcoded values
- Support for Terraform Cloud/Enterprise variable management
- Provide clear documentation on which fields need manual population
- Generate warnings/notices in output about missing sensitive fields

**Best Practices for Users**
1. Review generated configuration for `# TODO: Add authentication token` comments
2. Use Terraform variables for all sensitive values (never commit secrets)
3. Store secrets in Terraform Cloud/Enterprise or secure secret management systems
4. Use `terraform plan` to identify missing required fields before apply
5. Consider using `--var-file` flag to generate variable templates


### Maintenance

**SDK Updates**
- Automatic benefit from SDK updates when OpenAPI schema changes
- No manual API client maintenance
- Type safety ensures compatibility

**Provider Alignment**
- CLI tool stays aligned with provider capabilities
- New resource types automatically supported
- Schema changes propagate automatically

