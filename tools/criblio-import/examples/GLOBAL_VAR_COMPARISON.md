# Global Var Comparison: Why Updates After Import

## Key Differences

### Example File (`examples/global-var/main.tf`)
```hcl
resource "criblio_global_var" "my_globalvar" {
  description = "test"
  group_id    = "default"
  id          = "sample_globalvar"
  lib         = "test"
  tags        = "test"
  type        = "number"
  value       = 100        # ← Number, not string!
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

### State After Import (`terraform state show`)
```
resource "criblio_global_var" "unixtime" {
    group_id = "default"
    id       = "unixtime"
    items    = [              # ← Computed field (read-only)
        {
            "description" = "\"Current epoch time\""  # ← Escaped/quoted
            "id"          = "\"unixtime\""
            "lib"         = "\"cribl\""
            "tags"        = "\"cribl,sample\""
            "type"        = "\"expression\""
            "value"       = "\"Math.floor(Date.now()/1000)\""
        },
    ]
}
```

## Why Updates Occur

### 1. **`items` Field is Computed (Read-Only)**
- `items` is a `Computed: true` field in the schema
- It contains the list representation of variables
- Should NOT be in config, only in state
- The Read function populates `items` but not individual fields directly

### 2. **Individual Fields Not Populated in State**
- After import, state only has `group_id`, `id`, and `items`
- Individual fields like `description`, `lib`, `tags`, `type`, `value` are NOT in state
- Your config HAS these fields
- Terraform sees: config has `description = "Current epoch time"` but state has `null` → wants to update

### 3. **Value Type Mismatch (Example File)**
- Example shows: `value = 100` (number)
- Schema shows: `value` is `StringAttribute` (should be string)
- This might be a documentation issue or the example might be wrong

## The Root Cause

The `Read` function populates the `items` field (computed), but **doesn't populate the individual fields** (`description`, `lib`, `tags`, `type`, `value`) in state. This is why:

1. After import: State has `id`, `group_id`, `items`
2. Your config: Has `id`, `group_id`, `description`, `lib`, `tags`, `type`, `value`
3. Terraform sees: Config has fields that state doesn't → wants to update

## Solution

### Option 1: Remove Fields from Config (Minimal Config)
Only specify required fields:
```hcl
resource "criblio_global_var" "unixtime" {
  id       = "unixtime"
  group_id = "default"
}
```

Then run `terraform refresh` to populate state, then use `terraform plan -generate-config-out=complete.tf` to get the full config.

### Option 2: Fix the Read Function (Provider Issue)
The provider's `Read` function should populate individual fields (`description`, `lib`, `tags`, `type`, `value`) from the API response, not just `items`.

### Option 3: Accept the Updates
If your config values are correct, just run `terraform apply` to update the resources to match your config.

## Comparison Summary

| Aspect | Example File | Generated File | State After Import |
|--------|-------------|----------------|-------------------|
| `value` type | `100` (number) | `"100"` (string) | In `items` as string |
| Individual fields | ✅ In config | ✅ In config | ❌ Only in `items` |
| `items` field | ❌ Not in config | ❌ Not in config | ✅ Computed field |

The issue is that the Read function doesn't extract individual fields from `items` and populate them in state.

