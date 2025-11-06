# Best Practice for Importing criblio_group Resources

## The Problem

When you import a resource with `terraform import`, only the `id` field is set in state. If your config file has other fields (like `product`, `is_fleet`, `on_prem`) specified, Terraform sees them as "changing from null to a value" and wants to replace the resource.

## Solution: Minimal Config + Import

### Step 1: Create Minimal Resource Definitions

Only specify the `id` field in your resource config:

```hcl
resource "criblio_group" "default" {
  id = "default"
}
```

### Step 2: Import the Resource

```bash
terraform import criblio_group.default "default"
```

### Step 3: Refresh State

```bash
terraform refresh
```

This populates all fields in state from the API.

### Step 4: Generate Config (Optional)

```bash
terraform plan -generate-config-out=generated.tf
```

This generates a complete config file with all fields populated from state.

### Step 5: Verify No Changes

```bash
terraform plan
```

Should show "No changes" if everything is correct.

## Alternative: Using Import Blocks (Terraform 1.5+)

You can use Terraform's import blocks for a cleaner approach:

```hcl
import {
  to = criblio_group.default
  id = "default"
}

resource "criblio_group" "default" {
  id = "default"
}
```

Then run:
```bash
terraform plan -generate-config-out=generated.tf
terraform apply
```

This will:
1. Import the resource
2. Generate complete config from state
3. Apply (no changes should be needed)

## Why This Works

- **Minimal config**: Only `id` is specified, so no fields trigger replacement
- **Refresh populates state**: After import, `terraform refresh` calls the Read function to populate all fields from API
- **State matches reality**: State now has all actual values from Cribl
- **No replacement needed**: Since config only has `id` and state has all fields, Terraform sees no changes

## What NOT to Do

❌ **Don't specify fields in config before import:**
```hcl
resource "criblio_group" "default" {
  id = "default"
  product = "stream"  # ← This will cause replacement!
  is_fleet = false    # ← This will cause replacement!
}
```

If these values don't exactly match what's in Cribl, Terraform will want to replace.

## Recommended Import Flow

1. **Generate minimal config** (only `id` fields)
2. **Import resources** (`terraform import` or import blocks)
3. **Refresh state** (`terraform refresh`)
4. **Generate complete config** (`terraform plan -generate-config-out=generated.tf`)
5. **Review and apply** (should be no changes if everything matches)

## Example: Complete Import Script

```bash
#!/bin/bash

# 1. Create minimal groups.tf with only id fields
cat > groups.tf << 'EOF'
resource "criblio_group" "default" {
  id = "default"
}
EOF

# 2. Import
terraform import criblio_group.default "default"

# 3. Refresh
terraform refresh

# 4. Generate complete config
terraform plan -generate-config-out=groups-complete.tf

# 5. Review
terraform plan
```

This approach ensures no replacement warnings!

