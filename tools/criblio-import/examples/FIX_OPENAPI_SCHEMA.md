# Fix OpenAPI Schema to Use Typed Structs

## Current Problem

The single endpoint (`/m/{groupId}/lib/vars/{id}`) uses `*a103` response which has `additionalProperties: true`, causing Speakeasy to generate generic `map[string]any` instead of typed `shared.GlobalVar` structs.

## Solution: Modify OpenAPI Schema

### Step 1: Find Where `*a103` is Defined

Search for `a103:` in `openapi.yml` to find the anchor definition.

### Step 2: Update the Response Schema

Change the `items` array in `*a103` from `additionalProperties: true` to reference the `GlobalVar` schema:

**Before:**
```yaml
&a103
"200":
  schema:
    type: object
    properties:
      items:
        type: array
        items:
          type: object
          additionalProperties: true  # ← Generic map
```

**After:**
```yaml
&a103
"200":
  schema:
    type: object
    properties:
      items:
        type: array
        items:
          $ref: "#/components/schemas/GlobalVar"  # ← Typed schema!
```

### Step 3: Alternative - Direct Response Definition

Instead of using `*a103`, define the response directly in the single endpoint:

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
                    $ref: "#/components/schemas/GlobalVar"  # ← Use typed schema
```

### Step 4: Regenerate SDK

After modifying `openapi.yml`:

```bash
# Regenerate SDK with Speakeasy
speakeasy generate sdk --schema openapi.yml --lang terraform

# Or if using their CLI
speakeasy generate terraform
```

### Step 5: Update Provider Code

After regeneration, `GetGlobalVariableByIDResponseBody` will have:
```go
Items []shared.GlobalVar  // ✅ Typed structs!
```

Then update `RefreshFromOperationsGetGlobalVariableByIDResponseBody`:

```go
func (r *GlobalVarResourceModel) RefreshFromOperationsGetGlobalVariableByIDResponseBody(...) {
    if resp != nil && len(resp.Items) > 0 {
        gv := resp.Items[0]  // ✅ Now it's a shared.GlobalVar, not a map!
        
        r.ID = types.StringValue(gv.ID)
        if gv.Description != nil {
            r.Description = types.StringValue(*gv.Description)
        }
        if gv.Lib != nil {
            r.Lib = types.StringValue(*gv.Lib)
        }
        // ... etc for all fields
    }
}
```

## Benefits

- ✅ SDK generates typed structs automatically
- ✅ No manual field extraction from maps
- ✅ Type-safe code
- ✅ Consistent with list endpoint
- ✅ Matches OpenAPI schema structure

## Note

This requires modifying the OpenAPI spec and regenerating the SDK. If you can't modify the OpenAPI spec (it's auto-generated), you'll need to use the provider code fix approach instead.

