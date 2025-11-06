# Solution: Use Speakeasy Overlay to Force Typed Schema

## The Problem

Speakeasy generates `Items []map[string]any` for the single endpoint because the response uses `additionalProperties: true` instead of a schema reference.

## Solution: Add Overlay Rule

You can use Speakeasy's **overlay** feature (similar to `entity-mapping-overlay.yml`) to override the response schema for the single endpoint.

### Option 1: Create Response Overlay

Create or update `.speakeasy/response-overlay.yml`:

```yaml
overlay: 1.0.0
info:
  title: Override response schema for GetGlobalVariableByID
  version: 1.0.0
actions:
  # Override the response schema for the single endpoint
  - target: "$.paths['/m/{groupId}/lib/vars/{id}'].get.responses['200'].content.application/json.schema.properties.items.items"
    update:
      $ref: "#/components/schemas/GlobalVar"
      # Remove additionalProperties
      type: null
      additionalProperties: null
```

### Option 2: Use Entity Mapping Overlay

Add to `.speakeasy/entity-mapping-overlay.yml`:

```yaml
  # Force typed schema for GetGlobalVariableByID response
  - target: "$.paths['/m/{groupId}/lib/vars/{id}'].get.responses['200'].content.application/json.schema.properties.items.items"
    update:
      x-speakeasy-entity: GlobalVar
      $ref: "#/components/schemas/GlobalVar"
```

### Option 3: Modify OpenAPI Directly

If overlays don't work, modify `openapi.yml` directly:

**Find the `*a103` anchor definition** (around line 67952-67968) and change:

```yaml
# Before
items:
  type: array
  items:
    type: object
    additionalProperties: true

# After  
items:
  type: array
  items:
    $ref: "#/components/schemas/GlobalVar"
```

Or directly in the single endpoint:

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

## After Modification

1. **Regenerate SDK:**
   ```bash
   speakeasy generate terraform
   # or
   make generate  # if you have a Makefile target
   ```

2. **Verify Generated Code:**
   Check `internal/sdk/models/operations/getglobalvariablebyid.go`:
   ```go
   // Should now be:
   Items []shared.GlobalVar  // ✅ Typed structs!
   
   // Instead of:
   Items []map[string]any    // ❌ Generic maps
   ```

3. **Update Refresh Function:**
   Update `RefreshFromOperationsGetGlobalVariableByIDResponseBody` to use typed structs:
   ```go
   if len(resp.Items) > 0 {
       gv := resp.Items[0]  // ✅ Now it's shared.GlobalVar!
       r.ID = types.StringValue(gv.ID)
       if gv.Description != nil {
           r.Description = types.StringValue(*gv.Description)
       }
       // ... etc
   }
   ```

## Testing

After making changes:
1. Regenerate SDK
2. Update provider code
3. Test import:
   ```bash
   terraform import criblio_global_var.unixtime '{"group_id": "default", "id": "unixtime"}'
   terraform refresh
   terraform state show criblio_global_var.unixtime
   # Should show: id, group_id, description, lib, tags, type, value (all fields!)
   ```

## Note

If the OpenAPI spec is auto-generated from Cribl's API, you may need to:
- Use an overlay (Option 1 or 2) to avoid modifying the base spec
- Or create a post-processing script to modify the spec before SDK generation
- Or request Cribl to update the OpenAPI spec to use schema references

The overlay approach is best if you can't modify the base OpenAPI spec directly.

