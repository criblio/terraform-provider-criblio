# Global Var Import Issue: Root Cause Analysis

## The Problem

After importing `criblio_global_var.unixtime`, Terraform shows updates because:

1. **State only has `items` field** (computed/read-only)
2. **Individual fields are NOT in state** (`description`, `lib`, `tags`, `type`, `value`)
3. **Your config HAS these fields**
4. **Terraform sees: config has fields, state doesn't → wants to update**

## Root Cause: Provider Bug

The `RefreshFromOperationsGetGlobalVariableByIDResponseBody` function in `globalvar_resource_sdk.go` only populates the `items` field:

```go
func (r *GlobalVarResourceModel) RefreshFromOperationsGetGlobalVariableByIDResponseBody(...) {
    // Only populates r.Items
    // Does NOT extract individual fields from items[0]
    // Does NOT populate r.Description, r.Lib, r.Tags, r.Type, r.Value
}
```

### What Should Happen

The Read function should:
1. Get the variable from API (which returns it in `items` array)
2. Extract the first item: `variable := resp.Items[0]`
3. Populate individual fields:
   - `r.Description = variable["description"]`
   - `r.Lib = variable["lib"]`
   - `r.Tags = variable["tags"]`
   - `r.Type = variable["type"]`
   - `r.Value = variable["value"]`

### What Actually Happens

The Read function:
1. Gets the variable from API
2. Only populates `r.Items` with the entire response
3. Does NOT extract individual fields
4. Individual fields remain `null` in state

## Comparison

### Example File (`examples/global-var/main.tf`)
```hcl
resource "criblio_global_var" "my_globalvar" {
  description = "test"
  group_id    = "default"
  id          = "sample_globalvar"
  lib         = "test"
  tags        = "test"
  type        = "number"
  value       = 100  # ← Note: number, not string (might be documentation issue)
}
```

### Generated File (`global-vars.tf`)
```hcl
resource "criblio_global_var" "unixtime" {
  id          = "unixtime"
  group_id    = "default"
  description = "Current epoch time"
  lib         = "cribl"
  tags        = "cribl,sample"
  type        = "expression"
  value       = "Math.floor(Date.now()/1000)"  # ← String (correct per schema)
}
```

### State After Import
```
resource "criblio_global_var" "unixtime" {
    group_id = "default"
    id       = "unixtime"
    items    = [  # ← Only computed field populated
        {
            "description" = "\"Current epoch time\""
            "id"          = "\"unixtime\""
            "lib"         = "\"cribl\""
            "tags"        = "\"cribl,sample\""
            "type"        = "\"expression\""
            "value"       = "\"Math.floor(Date.now()/1000)\""
        },
    ]
    # description = null  ❌
    # lib = null          ❌
    # tags = null         ❌
    # type = null         ❌
    # value = null        ❌
}
```

### Terraform Plan Shows
```
~ resource "criblio_global_var" "unixtime" {
    + description = "Current epoch time"  # ← Adding (was null in state)
    ~ items       = [...] -> (known after apply)
    + lib         = "cribl"                # ← Adding (was null in state)
    + tags        = "cribl,sample"         # ← Adding (was null in state)
    + type        = "expression"           # ← Adding (was null in state)
    + value       = "Math.floor(...)"      # ← Adding (was null in state)
}
```

## Why This Happens

The API response structure is:
```json
{
  "items": [
    {
      "id": "unixtime",
      "description": "Current epoch time",
      "lib": "cribl",
      "tags": "cribl,sample",
      "type": "expression",
      "value": "Math.floor(Date.now()/1000)"
    }
  ]
}
```

The provider's Refresh function puts everything in `items` but doesn't extract individual fields.

## Solutions

### Option 1: Minimal Config (Workaround)
Only specify required fields in config:
```hcl
resource "criblio_global_var" "unixtime" {
  id       = "unixtime"
  group_id = "default"
}
```

Then:
```bash
terraform import criblio_global_var.unixtime '{"group_id": "default", "id": "unixtime"}'
terraform refresh
terraform plan -generate-config-out=complete.tf
```

### Option 2: Fix Provider (Recommended)
Update `RefreshFromOperationsGetGlobalVariableByIDResponseBody` to extract individual fields from `items[0]` and populate the model.

### Option 3: Accept Updates
If your config values are correct, just apply:
```bash
terraform apply
```

This will update the resource to match your config (which is fine if values are correct).

## Note About Value Type

The example file shows `value = 100` (number), but the schema shows `value` is a `StringAttribute`. This suggests:
- The example might be wrong, OR
- Terraform auto-converts numbers to strings

Based on the schema, `value` should always be a string in HCL.

