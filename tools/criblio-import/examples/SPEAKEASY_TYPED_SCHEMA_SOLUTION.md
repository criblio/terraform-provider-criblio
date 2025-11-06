# Solution: Force Speakeasy to Use Typed Schema for Single Endpoint

## The Problem

The single endpoint (`/m/{groupId}/lib/vars/{id}`) uses `additionalProperties: true` which causes Speakeasy to generate generic `map[string]any` instead of typed `shared.GlobalVar` structs.

## Current OpenAPI Schema

### List Endpoint (✅ Works - Typed)
```yaml
"/m/{groupId}/lib/vars":
  get:
    responses:
      "200":
        content:
          application/json:
            schema:
              type: object
              properties:
                items:
                  type: array
                  items:
                    $ref: "#/components/schemas/GlobalVar"  # ← Typed reference!
```

### Single Endpoint (❌ Doesn't Work - Generic)
```yaml
"/m/{groupId}/lib/vars/{id}":
  get:
    responses: *a103  # ← Uses same schema but...
    
# Where *a103 is:
responses:
  "200":
    schema:
      type: object
      properties:
        items:
          type: array
          items:
            type: object
            additionalProperties: true  # ← Generic map!
```

## Solution Options

### Option 1: Modify OpenAPI Schema (Recommended)

Change the single endpoint response to reference the `GlobalVar` schema directly:

```yaml
"/m/{groupId}/lib/vars/{id}":
  get:
    responses:
      "200":
        description: a list of Global Variable objects
        content:
          application/json:
            schema:
              type: object
              properties:
                items:
                  type: array
                  items:
                    $ref: "#/components/schemas/GlobalVar"  # ← Use typed schema!
```

This will make Speakeasy generate:
```go
Items []shared.GlobalVar  // ✅ Typed structs instead of map[string]any
```

### Option 2: Use Speakeasy Annotation

Speakeasy supports `x-speakeasy-type-override` annotation:

```yaml
"/m/{groupId}/lib/vars/{id}":
  get:
    responses:
      "200":
        content:
          application/json:
            schema:
              type: object
              properties:
                items:
                  type: array
                  items:
                    x-speakeasy-type-override: GlobalVar  # ← Force type
                    type: object
                    additionalProperties: true
```

**Note:** This annotation may not be supported in all Speakeasy versions.

### Option 3: Fix Provider Code (Current Workaround)

Instead of modifying OpenAPI, fix the provider's `RefreshFromOperationsGetGlobalVariableByIDResponseBody` to extract fields from the map:

```go
func (r *GlobalVarResourceModel) RefreshFromOperationsGetGlobalVariableByIDResponseBody(
    ctx context.Context, 
    resp *operations.GetGlobalVariableByIDResponseBody,
) diag.Diagnostics {
    var diags diag.Diagnostics
    
    if resp != nil && len(resp.Items) > 0 {
        item := resp.Items[0]  // Get first item from array
        
        // Extract individual fields from map
        if idVal, ok := item["id"].(string); ok {
            r.ID = types.StringValue(idVal)
        }
        
        if descVal, ok := item["description"].(string); ok && descVal != "" {
            r.Description = types.StringValue(descVal)
        }
        
        if libVal, ok := item["lib"].(string); ok && libVal != "" {
            r.Lib = types.StringValue(libVal)
        }
        
        if tagsVal, ok := item["tags"].(string); ok && tagsVal != "" {
            r.Tags = types.StringValue(tagsVal)
        }
        
        if typeVal, ok := item["type"].(string); ok && typeVal != "" {
            r.Type = types.StringValue(typeVal)
        }
        
        if valueVal, ok := item["value"].(string); ok && valueVal != "" {
            r.Value = types.StringValue(valueVal)
        }
        
        // Also populate items (for backwards compatibility)
        r.Items = nil
        for _, itemsItem := range resp.Items {
            var items map[string]jsontypes.Normalized
            if len(itemsItem) > 0 {
                items = make(map[string]jsontypes.Normalized, len(itemsItem))
                for key, value := range itemsItem {
                    result, _ := json.Marshal(value)
                    items[key] = jsontypes.NewNormalizedValue(string(result))
                }
            }
            r.Items = append(r.Items, items)
        }
    }
    
    return diags
}
```

## Recommended Approach

**Option 1** (Modify OpenAPI) is the best long-term solution because:
- ✅ Fixes the root cause
- ✅ Makes SDK generate correct types automatically
- ✅ No manual field extraction needed
- ✅ Consistent with list endpoint

**Option 3** (Fix Provider Code) is the immediate workaround:
- ✅ No OpenAPI changes needed
- ✅ Works with current SDK
- ⚠️ Requires manual field extraction
- ⚠️ Needs to be maintained

## Implementation Steps

### To Fix OpenAPI Schema:

1. Find the `*a103` response definition
2. Change `items` from `additionalProperties: true` to `$ref: "#/components/schemas/GlobalVar"`
3. Regenerate SDK with Speakeasy
4. The generated `GetGlobalVariableByIDResponseBody` will have `Items []shared.GlobalVar`
5. Update `RefreshFromOperationsGetGlobalVariableByIDResponseBody` to use typed structs

### To Fix Provider Code (Current):

1. Update `RefreshFromOperationsGetGlobalVariableByIDResponseBody` to extract fields from map
2. Test that Read function populates all fields
3. Verify state matches config after import

## Testing

After implementing either solution:

```bash
# Import a resource
terraform import criblio_global_var.unixtime '{"group_id": "default", "id": "unixtime"}'

# Refresh state
terraform refresh

# Check state has all fields
terraform state show criblio_global_var.unixtime
# Should show: id, group_id, description, lib, tags, type, value

# Plan should show no changes
terraform plan
# Should show: "No changes" (or only computed field changes)
```

