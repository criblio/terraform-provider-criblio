# Why Read Function Doesn't Populate Individual Fields

## The Issue

The Read function **IS using the correct endpoint** (`GetGlobalVariableByID` → `/m/{groupId}/lib/vars/{id}`), but the SDK generated **different response types** for list vs single endpoints.

## The Problem: Response Type Mismatch

### List Endpoint (`GetGlobalVariable`)
**Endpoint:** `/m/{groupId}/lib/vars`  
**Response Type:** `GetGlobalVariableResponseBody`
```go
Items []shared.GlobalVar  // ✅ Typed structs!
```

**Our import tool uses this:**
```go
res, err := client.GlobalVariables.GetGlobalVariable(ctx, request)
globalVars := res.Object.Items  // []shared.GlobalVar - typed!
for _, gv := range globalVars {
    id := gv.ID                    // ✅ Direct access
    desc := gv.Description         // ✅ Direct access
    lib := gv.Lib                 // ✅ Direct access
}
```

### Single Endpoint (`GetGlobalVariableByID`)
**Endpoint:** `/m/{groupId}/lib/vars/{id}`  
**Response Type:** `GetGlobalVariableByIDResponseBody`
```go
Items []map[string]any  // ❌ Generic maps!
```

**Provider's Read function uses this:**
```go
res, err := r.client.GlobalVariables.GetGlobalVariableByID(ctx, *request)
resp := res.Object.Items  // []map[string]any - generic!
// Should extract: resp[0]["id"], resp[0]["description"], etc.
// But RefreshFromOperationsGetGlobalVariableByIDResponseBody doesn't do this
```

## Why This Happens

Both endpoints return the same OpenAPI schema structure:
```yaml
responses:
  "200":
    schema:
      type: object
      properties:
        items:
          type: array
          items:
            type: object
            additionalProperties: true
```

But the SDK generator (Speakeasy) created:
- **List endpoint**: Typed `shared.GlobalVar` structs (because it can infer from examples/usage)
- **Single endpoint**: Generic `map[string]any` (because `additionalProperties: true` is more generic)

## The Root Cause

The `RefreshFromOperationsGetGlobalVariableByIDResponseBody` function only populates the `items` field as a generic map:

```go
func (r *GlobalVarResourceModel) RefreshFromOperationsGetGlobalVariableByIDResponseBody(...) {
    // Only populates r.Items with generic maps
    // Does NOT extract individual fields from items[0]
    // Does NOT populate r.Description, r.Lib, r.Tags, r.Type, r.Value
}
```

It should:
1. Extract the first item: `variable := resp.Items[0]`
2. Parse fields from the map:
   - `r.ID = variable["id"]`
   - `r.Description = variable["description"]`
   - `r.Lib = variable["lib"]`
   - etc.

## The Fix Needed

Update `RefreshFromOperationsGetGlobalVariableByIDResponseBody` to extract individual fields from `items[0]`:

```go
func (r *GlobalVarResourceModel) RefreshFromOperationsGetGlobalVariableByIDResponseBody(...) {
    if resp != nil && len(resp.Items) > 0 {
        item := resp.Items[0]  // Get first item from array
        
        // Extract and populate individual fields
        if id, ok := item["id"].(string); ok {
            r.ID = types.StringValue(id)
        }
        if desc, ok := item["description"].(string); ok {
            r.Description = types.StringValue(desc)
        }
        if lib, ok := item["lib"].(string); ok {
            r.Lib = types.StringValue(lib)
        }
        // ... etc for all fields
        
        // Also populate items (for backwards compatibility)
        // ... existing items population code
    }
}
```

## Summary

- ✅ **Read function uses correct endpoint** (`/m/{groupId}/lib/vars/{id}`)
- ❌ **SDK generates generic maps** instead of typed structs for single endpoint
- ❌ **Refresh function doesn't extract fields** from the map
- ✅ **Our import tool works** because it uses the list endpoint with typed structs

The fix is to update the `RefreshFromOperationsGetGlobalVariableByIDResponseBody` function to extract individual fields from `items[0]` map.

