# Import Issue Summary

## The Problem

After importing `criblio_group` resources, Terraform wants to replace them because:

1. **ImportState only sets `id`**: State only has `id = "default"`, all other fields are `null`
2. **Config has `product`**: Your config has `product = "stream"` (required field)
3. **Read function doesn't populate state**: After `terraform refresh`, state still only has `id`
4. **Replacement triggered**: Since `product` has `RequiresReplaceIfConfigured()` and state has `null` but config has a value, Terraform wants to replace

## Root Cause

The provider's `Read` function should populate all fields from the API after refresh, but it's not working. This is a **provider issue** that needs to be fixed in the provider code.

## Workaround Solutions

### Option 1: Remove `product` from config after import

```bash
# 1. Import with minimal config (only id)
resource "criblio_group" "default" {
  id = "default"
}

terraform import criblio_group.default "default"

# 2. Remove product from config (it's in state after refresh)
# Just keep id in config

# 3. Verify
terraform plan
```

But this won't work if `product` is required.

### Option 2: Use Terraform's generate-config feature

```bash
# 1. Import
terraform import criblio_group.default "default"

# 2. Generate config from state
terraform plan -generate-config-out=generated.tf

# This should generate a complete config matching what's in state
```

### Option 3: Accept the replacement (if config values are correct)

If your config values match what you want in Cribl, just apply:

```bash
terraform apply
```

This will update the resource to match your config.

## Recommended Solution

The provider's `Read` function needs to be fixed to properly populate state after import. Until then:

1. **Import with minimal config** (only `id`)
2. **Accept that replacement warnings will occur**
3. **Review the plan carefully** to ensure values are correct
4. **Apply if you want to update the resource** to match your config

## Example: What SHOULD Work

```bash
# 1. Minimal config
resource "criblio_group" "default" {
  id = "default"
}

# 2. Import
terraform import criblio_group.default "default"

# 3. Refresh (SHOULD populate all fields but doesn't)
terraform refresh

# 4. State SHOULD have all fields but only has id
terraform state show criblio_group.default
# Shows: only id = "default"

# 5. This is the bug - Read function should populate state
```

The fix needs to be in the provider's `Read` function to properly fetch and populate all fields from the API.

