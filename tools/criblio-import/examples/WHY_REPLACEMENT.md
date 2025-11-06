# Why Terraform Wants to Replace After Import

## Root Cause

After importing, your Terraform state looks like this:
```hcl
resource "criblio_group" "default_fleet" {
    id = "default_fleet"
    # All other fields are null/missing
}
```

But your config file (`groups.tf`) has:
```hcl
resource "criblio_group" "default_fleet" {
  id = "default_fleet"
  product = "edge"              # ← Set in config
  on_prem = true                # ← Set in config  
  is_fleet = true               # ← Set in config
  worker_remote_access = true  # ← Set in config
}
```

## The Problem

1. **ImportState only sets `id`**: The provider's `ImportState` function only populates the `id` field
2. **Fields have `RequiresReplaceIfConfigured()`**: Fields like `product`, `is_fleet`, `on_prem` have this plan modifier
3. **Terraform sees a change**: When a field is `null` in state but set in config, Terraform sees it as "changing from null to a value"
4. **Replacement triggered**: Because of `RequiresReplaceIfConfigured()`, Terraform wants to destroy and recreate the resource

## The Solution

### Step 1: Run terraform refresh
After import, run:
```bash
terraform refresh
```

This should call the `Read` function for each resource and populate all fields from the API.

### Step 2: Verify state is populated
```bash
terraform state show criblio_group.default_fleet
```

You should see all fields populated with actual values from Cribl.

### Step 3: Check if config matches reality
If refresh worked but you still see replacement warnings, it means:
- Your config file has different values than what's actually in Cribl
- This is actually correct behavior - Terraform wants to update the resource to match your config

### Step 4: Remove fields from config (if not needed)
If you don't want to manage certain fields, remove them from `groups.tf`:
```hcl
resource "criblio_group" "default_fleet" {
  id = "default_fleet"
  # Don't specify product, is_fleet, etc. if you don't want to manage them
}
```

## Why Refresh Doesn't Always Work

The `Read` function needs to be able to fetch the resource from the API. If:
- The API call fails
- The resource doesn't exist
- Authentication issues

Then refresh won't populate the fields, and you'll see replacement warnings.

## Best Practice

1. **Import resources**
2. **Run `terraform refresh`** to populate state
3. **Review `terraform plan`** to see what would change
4. **Update your config** to match what's actually in Cribl (if needed)
5. **Run `terraform apply`** only if you want to make changes

## Current Status

The importer now includes a `Refresh` step after import. However, if you're still seeing replacement warnings:

1. Check if refresh actually populated the state
2. Compare state values with config values
3. Decide if you want Terraform to manage those fields or not

