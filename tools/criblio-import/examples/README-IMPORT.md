# Importing Groups to Terraform State

After generating `groups.tf`, you need to import the existing groups into Terraform state.

## Quick Start

```bash
# 1. Create provider.tf if it doesn't exist
cat > provider.tf << 'EOF'
terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Authentication via environment variables or ~/.cribl/credentials
}
EOF

# 2. Initialize Terraform
terraform init

# 3. Import groups
./import-groups-simple.sh

# 4. Review the plan
terraform plan

# 5. Apply if everything looks good
terraform apply
```

## Manual Import

If you prefer to import manually:

```bash
terraform import criblio_group.default "default"
terraform import criblio_group.defaultHybrid "defaultHybrid"
terraform import criblio_group.default_fleet "default_fleet"
terraform import criblio_group.edge_test "edge-test"
terraform import criblio_group.edge_filemonitor_fleet "edge_filemonitor_fleet"
```

## Import Format

For `criblio_group`, the import ID is just the group ID:
```
terraform import criblio_group.<resource_name> <group_id>
```

The `product` field is already set in the `groups.tf` file, so it doesn't need to be in the import ID.

## Troubleshooting

### Already in state
If a group is already imported, you'll see a message. You can remove it from state if needed:
```bash
terraform state rm criblio_group.<resource_name>
```

### Import errors
If import fails, check:
1. Authentication is working (credentials file or environment variables)
2. Group IDs match what's in Cribl
3. Product field is set correctly in groups.tf

### After import
Run `terraform plan` to see if there are any differences between the state and the actual configuration in Cribl.

