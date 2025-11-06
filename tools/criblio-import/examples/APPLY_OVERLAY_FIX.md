# Apply Overlay Fix for Typed Schema

## What Was Done

Added an overlay rule to `.speakeasy/entity-mapping-overlay.yml` to force the single endpoint response to use the typed `GlobalVar` schema instead of `additionalProperties: true`.

## The Overlay Rule

```yaml
  # Force typed schema for GetGlobalVariableByID response items
  - target: "$.paths['/m/{groupId}/lib/vars/{id}'].get.responses['200'].content.application/json.schema.properties.items.items"
    update:
      $ref: "#/components/schemas/GlobalVar"
      x-speakeasy-entity: GlobalVar
```

This tells Speakeasy:
- Target the `items.items` property in the response schema
- Replace `additionalProperties: true` with `$ref: "#/components/schemas/GlobalVar"`
- Use `GlobalVar` entity for type generation

## Next Steps

1. **Regenerate SDK:**
   ```bash
   # Run Speakeasy generation
   make generate
   # or
   speakeasy generate terraform
   ```

2. **Verify Generated Code:**
   Check `internal/sdk/models/operations/getglobalvariablebyid.go`:
   ```go
   // Should now be:
   Items []shared.GlobalVar  // ✅ Typed structs!
   ```

3. **Update Provider Code:**
   Update `RefreshFromOperationsGetGlobalVariableByIDResponseBody` in `internal/provider/globalvar_resource_sdk.go`:
   ```go
   func (r *GlobalVarResourceModel) RefreshFromOperationsGetGlobalVariableByIDResponseBody(...) {
       if resp != nil && len(resp.Items) > 0 {
           gv := resp.Items[0]  // ✅ Now it's shared.GlobalVar, not map[string]any!
           
           r.ID = types.StringValue(gv.ID)
           if gv.Description != nil {
               r.Description = types.StringValue(*gv.Description)
           }
           if gv.Lib != nil {
               r.Lib = types.StringValue(*gv.Lib)
           }
           if gv.Tags != nil {
               r.Tags = types.StringValue(*gv.Tags)
           }
           if gv.Type != nil {
               r.Type = types.StringValue(string(*gv.Type))
           }
           if gv.Value != nil {
               r.Value = types.StringValue(*gv.Value)
           }
           
           // Also populate items for backwards compatibility
           // ... existing items code
       }
   }
   ```

4. **Test:**
   ```bash
   terraform import criblio_global_var.unixtime '{"group_id": "default", "id": "unixtime"}'
   terraform refresh
   terraform state show criblio_global_var.unixtime
   # Should show all fields: id, group_id, description, lib, tags, type, value
   
   terraform plan
   # Should show "No changes" (or only computed field changes)
   ```

## Benefits

- ✅ **No OpenAPI spec changes needed** - overlay handles it
- ✅ **SDK generates typed structs automatically**
- ✅ **Type-safe code**
- ✅ **Consistent with list endpoint**
- ✅ **Works with auto-generated OpenAPI specs**

## How Overlays Work

Speakeasy overlays apply transformations to the OpenAPI spec before generation:
- They modify the spec in-memory during generation
- Don't require changing the base `openapi.yml` file
- Perfect for fixing auto-generated specs

The overlay we added will:
1. Find the response schema for `/m/{groupId}/lib/vars/{id}`
2. Replace `items.items` from `additionalProperties: true` to `$ref: "#/components/schemas/GlobalVar"`
3. Speakeasy generates `Items []shared.GlobalVar` instead of `Items []map[string]any`

