# ✅ Working Import Solution

## The Key: Minimal Config (Only Required Fields)

After testing with the examples, here's what works:

### Step 1: Create Minimal Config

Only specify **required fields** (`id` and `product`):

```hcl
resource "criblio_group" "default" {
  id      = "default"
  product = "stream"  # Required field
}
```

### Step 2: Import

```bash
terraform import criblio_group.default "default"
```

### Step 3: Refresh

```bash
terraform refresh
```

This populates all fields in state from the API.

### Step 4: Verify

```bash
terraform plan
```

Should show **"No changes"** if everything matches!

## Why This Works

1. **Minimal config**: Only required fields (`id`, `product`) are specified
2. **No replacement triggers**: Fields with `RequiresReplaceIfConfigured()` are NOT in config
3. **Refresh populates state**: After refresh, state has all actual values from API
4. **State matches reality**: Since config only has required fields, Terraform sees no changes

## Comparison

### ❌ Full Config (Causes Replacement)
```hcl
resource "criblio_group" "default" {
  id      = "default"
  product = "stream"
  is_fleet = false      # ← Causes replacement if doesn't match
  on_prem = false       # ← Causes replacement if doesn't match
  cloud = { ... }      # ← Causes replacement if doesn't match
}
```

### ✅ Minimal Config (Works!)
```hcl
resource "criblio_group" "default" {
  id      = "default"
  product = "stream"  # Only required fields
}
```

## Updated Generator

The `generate-groups.go` script should generate minimal configs:

```go
terraformResources.WriteString(fmt.Sprintf("resource \"criblio_group\" \"%s\" {\n", resourceName))
terraformResources.WriteString(fmt.Sprintf("  id      = %q\n", group.ID))
terraformResources.WriteString(fmt.Sprintf("  product = %q\n", groupWithProduct.Product))
terraformResources.WriteString("}\n\n")
```

Then users can:
1. Import the resources
2. Run `terraform refresh` to populate state
3. Optionally run `terraform plan -generate-config-out=complete.tf` to get full config

## Example: Complete Flow

```bash
# 1. Generate minimal config
go run generate-groups.go

# 2. Import all resources
terraform import criblio_group.default "default"
terraform import criblio_group.defaultHybrid "defaultHybrid"
# ... etc

# 3. Refresh state
terraform refresh

# 4. Verify no changes
terraform plan
# Should show: "No changes"

# 5. (Optional) Generate complete config
terraform plan -generate-config-out=groups-complete.tf
```

This approach ensures **no replacement warnings**! 🎉

