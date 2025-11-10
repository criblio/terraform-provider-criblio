# Using SDK for Field Name Conversion

## The Insight

The SDK already handles the conversion from API camelCase fields to Terraform snake_case fields! We don't need manual field mapping - we can reuse the provider's conversion methods.

## How It Works

### Provider's Conversion Flow

The provider uses two types of conversion methods:

1. **`RefreshFrom*` methods**: Convert from API types (camelCase) → Terraform types (snake_case)
   - Example: `RefreshFromSharedGlobalVar(ctx, resp *shared.GlobalVar)`
   - Takes SDK `shared.GlobalVar` (camelCase fields) 
   - Populates `GlobalVarResourceModel` (snake_case fields via `tfsdk` tags)

2. **`ToShared*` methods**: Convert from Terraform types (snake_case) → API types (camelCase)
   - Example: `ToSharedGlobalVar(ctx) *shared.GlobalVar`
   - Takes `GlobalVarResourceModel` (snake_case fields)
   - Returns `shared.GlobalVar` (camelCase fields)

### Field Name Mapping

The mapping is defined in the ResourceModel structs using `tfsdk` tags:

```go
type GlobalVarResourceModel struct {
    ID          types.String `tfsdk:"id"`           // Terraform: id
    GroupID     types.String `tfsdk:"group_id"`     // Terraform: group_id
    Description types.String `tfsdk:"description"`  // Terraform: description
    // ...
}
```

The SDK's `shared.GlobalVar` has JSON tags for API fields:

```go
type GlobalVar struct {
    ID          string  `json:"id"`           // API: id
    GroupID     *string `json:"groupId"`      // API: groupId (camelCase)
    Description *string `json:"description"`  // API: description
    // ...
}
```

The `RefreshFromSharedGlobalVar` method handles the conversion:

```go
func (r *GlobalVarResourceModel) RefreshFromSharedGlobalVar(ctx context.Context, resp *shared.GlobalVar) {
    r.ID = types.StringValue(resp.ID)                    // id → id
    r.GroupID = types.StringValue(*resp.GroupID)          // groupId → group_id
    r.Description = types.StringPointerValue(resp.Description) // description → description
    // ...
}
```

## How CLI Tool Can Use This

### Approach: Reuse RefreshFrom Methods

Instead of manual field mapping, the CLI tool can:

1. **Parse YAML** to get the configuration structure
2. **Convert YAML to SDK types** (shared.* structs with camelCase fields)
   - YAML files from diag bundle have camelCase (same as API)
   - Unmarshal YAML into `shared.GlobalVar`, `shared.Source`, etc.
3. **Use RefreshFrom methods** to convert to Terraform ResourceModel types
   - Call `RefreshFromSharedGlobalVar(ctx, &sharedGlobalVar)`
   - This automatically converts camelCase → snake_case
4. **Generate HCL** from ResourceModel types
   - ResourceModel has `tfsdk` tags defining Terraform field names
   - Generate HCL using these field names

### Example Flow

```go
// 1. Parse YAML (camelCase fields)
yamlData := parseYAML("global-var.yaml")
// yamlData: { "id": "unixtime", "groupId": "default", "description": "..." }

// 2. Convert to SDK type (camelCase)
var sharedVar shared.GlobalVar
yaml.Unmarshal(yamlData, &sharedVar)
// sharedVar.GroupID = "default" (camelCase field)

// 3. Use RefreshFrom method (automatic conversion)
var model GlobalVarResourceModel
model.RefreshFromSharedGlobalVar(ctx, &sharedVar)
// model.GroupID = "default" (snake_case field via tfsdk:"group_id")

// 4. Generate HCL from model
hcl := generateHCL(&model)
// hcl: resource "criblio_global_var" { group_id = "default" ... }
```

## Benefits

### ✅ No Manual Field Mapping
- Don't need to maintain `field_mapper.go` with manual mappings
- Field names are defined once in provider ResourceModel structs
- Automatic conversion via `RefreshFrom*` methods

### ✅ Always in Sync
- When provider schemas change, conversions automatically update
- No risk of field mapping drift
- Same conversion logic as provider uses

### ✅ Type Safety
- Works with generated SDK types
- Compile-time type checking
- No string-based field name lookups

### ✅ Handles Complex Types
- Nested structures automatically handled
- Arrays, maps, objects all converted correctly
- Type conversions (string → number, etc.) handled by SDK

## Implementation Strategy

### Option 1: Direct SDK Type Usage (Recommended)

```go
// Parse YAML to SDK shared types
var sharedVar shared.GlobalVar
yaml.Unmarshal(data, &sharedVar)

// Use RefreshFrom method
var model GlobalVarResourceModel
model.RefreshFromSharedGlobalVar(ctx, &sharedVar)

// Generate HCL from model (using tfsdk tags)
hcl := generateHCLFromModel(&model)
```

### Option 2: Hybrid Approach

For resources where YAML structure doesn't exactly match SDK types:
1. Parse YAML to `map[string]interface{}`
2. Convert to SDK types manually (for complex cases)
3. Use `RefreshFrom*` methods for field name conversion
4. Generate HCL from ResourceModel

## Current State

The current `field_mapper.go` implementation does manual camelCase → snake_case conversion. This can be replaced with:

1. **YAML → SDK types**: Unmarshal YAML into `shared.*` structs
2. **SDK → Terraform**: Use `RefreshFrom*` methods
3. **Terraform → HCL**: Generate from ResourceModel using `tfsdk` tags

This approach:
- ✅ Eliminates manual field mapping
- ✅ Reuses provider's conversion logic
- ✅ Ensures consistency with provider
- ✅ Automatically stays in sync with schema changes

## Migration Path

1. Update converters to use SDK types instead of `map[string]interface{}`
2. Use `RefreshFrom*` methods for field conversion
3. Remove or simplify `field_mapper.go` (only needed for edge cases)
4. Generate HCL from ResourceModel structs using reflection on `tfsdk` tags

This is a more robust, maintainable approach that maximizes code reuse!

