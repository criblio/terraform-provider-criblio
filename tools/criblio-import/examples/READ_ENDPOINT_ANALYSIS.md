# Read Function Endpoint Analysis

## The Answer: Read IS Using the Correct Endpoint!

The Read function **IS using** `GetGlobalVariableByID` which maps to `/m/{groupId}/lib/vars/{id}`.

## The Real Problem: Response Type Difference

Both endpoints use the same response schema (`*a103`), but the SDK generates **different types**:

### List Endpoint (`/m/{groupId}/lib/vars`)
```go
// GetGlobalVariable returns:
Items []shared.GlobalVar  // ✅ Typed structs
```

**Why?** The SDK can infer the structure from the `GlobalVar` schema reference.

### Single Endpoint (`/m/{groupId}/lib/vars/{id}`)
```go
// GetGlobalVariableByID returns:
Items []map[string]any  // ❌ Generic maps
```

**Why?** The response schema uses `additionalProperties: true`, which the SDK interprets as a generic map.

## The OpenAPI Schema

Both endpoints reference `*a103`:

```yaml
responses: *a103
```

Where `*a103` is:
```yaml
"200":
  description: a list of Global Variable objects
  schema:
    type: object
    properties:
      items:
        type: array
        items:
          type: object
          additionalProperties: true  # ← This causes generic maps!
```

## Why This Matters

The `RefreshFromOperationsGetGlobalVariableByIDResponseBody` function receives:
- `resp.Items` = `[]map[string]any` (generic maps)

It should extract individual fields from `items[0]`:
```go
if len(resp.Items) > 0 {
    item := resp.Items[0]  // map[string]any
    
    // Extract fields
    if id, ok := item["id"].(string); ok {
        r.ID = types.StringValue(id)
    }
    if desc, ok := item["description"].(string); ok {
        r.Description = types.StringValue(desc)
    }
    // ... etc
}
```

But currently it only:
- Stores the whole map in `r.Items` (computed field)
- Doesn't extract individual fields

## The Fix

The `RefreshFromOperationsGetGlobalVariableByIDResponseBody` function needs to:
1. Extract `items[0]` (the variable data)
2. Parse individual fields from the map
3. Populate `r.Description`, `r.Lib`, `r.Tags`, `r.Type`, `r.Value`

This is a **provider code fix**, not an endpoint issue.

## Summary

- ✅ **Read uses correct endpoint**: `/m/{groupId}/lib/vars/{id}`
- ✅ **Endpoint is correct per OpenAPI spec**
- ❌ **Response parsing is incomplete**: Only populates `items`, not individual fields
- 🔧 **Fix needed**: Extract fields from `items[0]` map in `RefreshFromOperationsGetGlobalVariableByIDResponseBody`

