# Perfect State Matching After Import

## The Goal

After import, we want `terraform plan` to show **"No changes"** (or minimal changes).

## The Challenge

The provider's Read function only populates the `items` computed field, not individual fields. This means:
- State has: `id`, `group_id`, `items` (computed)
- Config can't have: `items` (it's computed/read-only)
- Config can have: `id`, `group_id`, `description`, `lib`, `tags`, `type`, `value`

## Solution: Use `terraform plan -generate-config-out`

Terraform's `-generate-config-out` flag generates a config file that **exactly matches what's in state**.

### Step-by-Step Process

1. **Generate minimal config** (only required fields)
   ```bash
   go run generate-global-vars.go
   ```

2. **Import resources**
   ```bash
   terraform init
   terraform import criblio_global_var.unixtime '{"group_id": "default", "id": "unixtime"}'
   # ... import all resources
   ```

3. **Refresh state**
   ```bash
   terraform refresh
   ```

4. **Generate config from state** (this matches exactly!)
   ```bash
   terraform plan -generate-config-out=global-vars-complete.tf
   ```

5. **Replace config with generated one**
   ```bash
   mv global-vars-complete.tf global-vars.tf
   ```

6. **Verify no changes**
   ```bash
   terraform plan
   # Should show: "No changes" or only computed field changes
   ```

### Automated Script

Use the provided script:
```bash
./import-and-match-state.sh
```

This script:
1. Generates minimal config
2. Imports all resources
3. Refreshes state
4. Generates config from state
5. Replaces config with state-matched version

## What Gets Generated

The generated config from state will have:
- Only fields that are actually in state (not just in `items`)
- Since state only has `id`, `group_id`, and `items`, the generated config will be minimal
- But it will match state exactly, so no updates!

## Computed Fields

The `items` field is computed and may still show changes. This is:
- ✅ **Expected behavior** - computed fields always show as changing
- ✅ **Harmless** - won't actually modify the resource
- ✅ **Normal** - this is how Terraform works with computed fields

## Result

After using `terraform plan -generate-config-out`:
- ✅ Config matches state exactly
- ✅ No individual field updates
- ⚠️ `items` computed field may still show (expected, harmless)

This is the best we can achieve without fixing the provider's Read function.

