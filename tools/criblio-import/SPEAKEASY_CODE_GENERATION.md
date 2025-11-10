# Speakeasy Code Generation Process

## Overview

The Terraform provider uses **Speakeasy** to automatically generate both the SDK and Terraform provider code from the OpenAPI schema (`openapi.yml`).

## Code Generation Flow

```
┌─────────────────┐
│  openapi.yml    │  ← Base OpenAPI specification
└────────┬────────┘
         │
         │ + Overlays
         ▼
┌─────────────────────────────┐
│ entity-mapping-overlay.yml  │  ← Modifies schema before generation
└────────┬────────────────────┘
         │
         │ Speakeasy Generation
         ▼
┌─────────────────────────────────────────┐
│  Generated Code                         │
│  ├── internal/sdk/                      │  ← API client SDK
│  │   ├── models/                        │
│  │   │   ├── operations/                │  ← API operations
│  │   │   └── shared/                    │  ← Shared types/models
│  │   └── criblio.go                     │  ← Main SDK client
│  │                                       │
│  └── internal/provider/                 │  ← Terraform provider
│      ├── *_resource.go                  │  ← Resource definitions
│      ├── *_resource_sdk.go              │  ← SDK integration
│      └── *_data_source.go               │  ← Data source definitions
└─────────────────────────────────────────┘
```

## Key Components

### 1. **Base Schema**: `openapi.yml`
- OpenAPI 3.0 specification
- Defines all API endpoints, request/response schemas
- Source of truth for API structure

### 2. **Overlays**: `.speakeasy/entity-mapping-overlay.yml`
- Modifies the OpenAPI schema before code generation
- Used to:
  - Add `x-speakeasy-entity` tags for deterministic generation
  - Override response schemas (like the GlobalVar fix we added)
  - Resolve non-deterministic type generation

### 3. **Generation Config**: `.speakeasy/gen.yaml`
- Speakeasy configuration
- Defines generation settings:
  - SDK class name
  - Terraform provider settings
  - Authentication configuration
  - Code generation options

### 4. **Generated Code**

#### SDK (`internal/sdk/`)
- **API Client**: `criblio.go` - Main client with all API operations
- **Operations**: `models/operations/*.go` - Individual API endpoint handlers
- **Shared Types**: `models/shared/*.go` - Reusable data structures

#### Terraform Provider (`internal/provider/`)
- **Resources**: `*_resource.go` - Terraform resource definitions
- **SDK Integration**: `*_resource_sdk.go` - Code to convert between Terraform and SDK types
- **Data Sources**: `*_data_source.go` - Read-only data sources

## How Overlays Work

The overlay file modifies the OpenAPI schema **before** Speakeasy generates code:

### Example: GlobalVar Typed Schema Fix

**Problem**: Single endpoint returned `Items []map[string]any` instead of `Items []shared.GlobalVar`

**Solution**: Added overlay rule:
```yaml
- target: "$.paths['/m/{groupId}/lib/vars/{id}'].get.responses['200'].content.application/json.schema.properties.items.items"
  update:
    $ref: "#/components/schemas/GlobalVar"
    x-speakeasy-entity: GlobalVar
```

**Result**: After regeneration, `GetGlobalVariableByIDResponseBody` now has:
```go
Items []shared.GlobalVar  // ✅ Typed structs!
```

## Generation Process

### Step 1: Update OpenAPI Schema
- Modify `openapi.yml` or add overlay rules

### Step 2: Run Speakeasy Generation
```bash
# Typically via Makefile or direct command
make generate
# or
speakeasy generate terraform
```

### Step 3: Review Generated Code
- Check `internal/sdk/` for SDK changes
- Check `internal/provider/` for provider changes
- Verify types match expectations

### Step 4: Update Provider Code (if needed)
- If SDK types change, update provider's `RefreshFrom*` functions
- Update resource schemas if needed
- Test with `terraform plan/apply`

## Benefits for CLI Tool

### 1. **Code Reuse**
The CLI tool can directly import the generated SDK:
```go
import "github.com/criblio/terraform-provider-criblio/internal/sdk"
```

### 2. **Type Safety**
- Same types used in provider and CLI
- Consistent API interaction
- Compile-time type checking

### 3. **Automatic Updates**
- When OpenAPI schema updates, SDK regenerates
- CLI tool automatically gets new API endpoints/types
- No manual API client maintenance

### 4. **Authentication**
- SDK includes `CriblTerraformHook` for automatic auth
- CLI tool reuses same authentication logic
- Supports all auth methods (Bearer, OAuth, on-prem, credentials file)

## Overlay Patterns

### Pattern 1: Entity Tagging
```yaml
- target: "$.components.schemas.GlobalVar"
  update:
    x-speakeasy-entity: GlobalVar
```
**Purpose**: Ensures deterministic type generation

### Pattern 2: Response Schema Override
```yaml
- target: "$.paths['/m/{groupId}/lib/vars/{id}'].get.responses['200'].content.application/json.schema.properties.items.items"
  update:
    $ref: "#/components/schemas/GlobalVar"
```
**Purpose**: Force typed schema instead of generic maps

### Pattern 3: Path-Level Entity Mapping
```yaml
- target: "$.paths['/m/{groupId}/lib/grok/{id}'].get.responses['200'].content.application/json.schema"
  update:
    x-speakeasy-entity: Grok
```
**Purpose**: Map endpoint response to specific entity type

## CLI Tool Integration

The CLI tool benefits from this generation process:

1. **Import Generated SDK**:
   ```go
   import (
       "github.com/criblio/terraform-provider-criblio/internal/sdk"
       "github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
       "github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
   )
   ```

2. **Use Generated Types**:
   ```go
   client := sdk.New(...)
   resp, err := client.GlobalVariables.GetGlobalVariable(ctx, request)
   // resp.Object.Items is []shared.GlobalVar (typed!)
   ```

3. **Reuse Authentication**:
   ```go
   // SDK automatically handles auth via CriblTerraformHook
   // No manual auth code needed
   ```

## Maintenance Workflow

### When OpenAPI Schema Updates:

1. **Update `openapi.yml`** (or receive updated version)
2. **Add overlay rules** if needed (in `entity-mapping-overlay.yml`)
3. **Regenerate code**: `make generate`
4. **Review changes** in `internal/sdk/` and `internal/provider/`
5. **Update provider code** if SDK types changed
6. **Test provider**: `terraform plan/apply`
7. **CLI tool automatically benefits** from updated SDK

### When Adding New Resource Types:

1. **Ensure OpenAPI schema** includes the resource
2. **Add overlay rule** if needed for entity mapping
3. **Regenerate**: `make generate`
4. **Provider code auto-generated** in `internal/provider/`
5. **CLI tool can use** new resource types immediately

## Key Takeaways

✅ **Single Source of Truth**: `openapi.yml` defines everything  
✅ **Automatic Generation**: SDK and provider code auto-generated  
✅ **Type Safety**: Consistent types across provider and CLI  
✅ **Code Reuse**: CLI tool directly imports generated SDK  
✅ **Easy Updates**: Schema changes propagate automatically  
✅ **Overlay Flexibility**: Can fix schema issues without modifying base OpenAPI  

This architecture makes the CLI tool development much simpler - we don't need to maintain a separate API client, we just import the generated SDK!

