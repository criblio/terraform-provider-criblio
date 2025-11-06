# Why Terraform Wants to Replace After Import

## The Problem

After importing resources, Terraform shows replacement warnings because:

1. **ImportState only sets `id`**: The `ImportState` function only sets the `id` field in state
2. **State has null values**: All other fields (like `product`, `is_fleet`, `on_prem`) are `null` in state
3. **Config has actual values**: Your `groups.tf` file has these fields set (e.g., `product = "edge"`)
4. **RequiresReplaceIfConfigured()**: Fields like `product`, `is_fleet`, `on_prem` have this plan modifier
5. **Terraform sees a change**: When a field is `null` in state but set in config, and it has `RequiresReplaceIfConfigured()`, Terraform wants to replace

## The Solution

After import, you need to refresh the state to populate all fields from the API. However, the `Read` function might need the `product` field to work correctly.

### Option 1: Run terraform refresh (Recommended)

After importing, run:
```bash
terraform refresh
```

This will call the Read function for each resource and populate all fields from the API.

### Option 2: Use terraform import with state manipulation

You can use `terraform state` commands to set the product field before refresh:

```bash
# After import, set product in state
terraform state show criblio_group.default_fleet
# Then manually update state if needed
```

### Option 3: Update the ImportState function

The provider's `ImportState` function could be enhanced to read the resource and set all fields. However, this requires changes to the provider itself.

## Current Workaround

The importer now includes a `Refresh` step after import. However, if the Read function fails because required fields are missing, you may need to:

1. Ensure your config file matches what's actually in Cribl
2. Run `terraform refresh` manually if needed
3. Or remove the fields from config that have `RequiresReplaceIfConfigured()` if they're not needed

## Fields with RequiresReplaceIfConfigured()

- `product` - Required for API calls
- `is_fleet` - Changes the group type
- `on_prem` - Changes deployment type
- `cloud` - Changes cloud configuration

These fields trigger replacement when they differ between state and config.

