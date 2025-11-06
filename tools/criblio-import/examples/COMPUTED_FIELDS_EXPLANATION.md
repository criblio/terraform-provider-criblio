# Why Computed Fields Show Changes

## The Issue

After import, even with minimal config (only `id` and `group_id`), Terraform still shows changes because of the `items` computed field.

## Understanding Computed Fields

The `items` field is:
- `Computed: true` - Populated by the provider, not by config
- Read-only - You can't set it in config
- Always recalculated - Terraform shows it as changing

## What's Happening

### State After Import + Refresh
```
resource "criblio_global_var" "unixtime" {
    group_id = "default"
    id       = "unixtime"
    items    = [  # ← Computed field, always shown
        {
            "description" = "\"Current epoch time\""
            ...
        }
    ]
}
```

### Our Config
```hcl
resource "criblio_global_var" "unixtime" {
  id       = "unixtime"
  group_id = "default"
  type     = "expression"
}
```

### Terraform Plan Shows
```
~ resource "criblio_global_var" "unixtime" {
    ~ items = [...] -> (known after apply)  # ← Computed field always shows change
    + type  = "expression"                   # ← Adding because it's not in state
}
```

## Why This Happens

1. **`items` is computed**: Terraform will always show it as changing because it's recalculated on every plan
2. **Individual fields not in state**: The Read function doesn't populate `description`, `lib`, `tags`, `type`, `value` - they're only in `items`
3. **Config has `type`**: We added `type` to match the default, but state doesn't have it (it's in `items`)

## Solutions

### Option 1: Accept Computed Field Changes (Recommended)
The `items` field change is **expected behavior** for computed fields. It doesn't mean the resource will actually change - it's just Terraform showing that the computed field will be recalculated.

```bash
# The items change is harmless - it's just a computed field
terraform apply  # Will not actually change anything
```

### Option 2: Use `terraform plan -generate-config-out`
After import, generate config from state:
```bash
terraform import criblio_global_var.unixtime '{"group_id": "default", "id": "unixtime"}'
terraform refresh
terraform plan -generate-config-out=complete.tf
```

This generates a config that matches what's in state, which should minimize changes.

### Option 3: Suppress Diff in Provider (Provider Fix)
The provider could add a `SuppressDiff` plan modifier to the `items` field to prevent it from showing changes. But this requires provider changes.

## Current Behavior

With minimal config (only `id`, `group_id`, `type`):
- ✅ No individual field updates (except `type` which isn't in state)
- ⚠️ `items` computed field will always show as changing (expected)

The `items` change is **cosmetic** - it won't actually modify the resource. This is normal behavior for computed/read-only fields in Terraform.

## Best Practice

1. **Use minimal config** (only required fields)
2. **Accept that computed fields show changes** (this is normal)
3. **Use `terraform plan -generate-config-out`** if you want the full config
4. **Run `terraform apply`** - it won't actually change anything, just update state

