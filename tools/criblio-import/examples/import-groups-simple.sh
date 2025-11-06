#!/bin/bash

# Simple script to import groups - manually specify the imports
# Based on the generated groups.tf file

set -e

echo "🚀 Importing groups into Terraform state..."
echo ""

# Initialize Terraform if needed
if [ ! -d ".terraform" ]; then
    echo "📦 Initializing Terraform..."
    terraform init
    echo ""
fi

# Import format: terraform import <resource_type>.<resource_name> <group_id>
# Note: The ImportState function only uses the ID, but product field must be set in the resource
# The groups.tf file already has product set, so we just need the group ID

echo "📥 Importing groups..."
echo ""

# Group 1: default (stream)
echo "  Importing criblio_group.default (default)"
terraform import criblio_group.default "default" || echo "    ⚠️  Already imported or failed"

# Group 2: defaultHybrid (stream)
echo "  Importing criblio_group.defaultHybrid (defaultHybrid)"
terraform import criblio_group.defaultHybrid "defaultHybrid" || echo "    ⚠️  Already imported or failed"

# Group 3: default_fleet (edge)
echo "  Importing criblio_group.default_fleet (default_fleet)"
terraform import criblio_group.default_fleet "default_fleet" || echo "    ⚠️  Already imported or failed"

# Group 4: edge_test (edge)
echo "  Importing criblio_group.edge_test (edge-test)"
terraform import criblio_group.edge_test "edge-test" || echo "    ⚠️  Already imported or failed"

# Group 5: edge_filemonitor_fleet (edge)
echo "  Importing criblio_group.edge_filemonitor_fleet (edge_filemonitor_fleet)"
terraform import criblio_group.edge_filemonitor_fleet "edge_filemonitor_fleet" || echo "    ⚠️  Already imported or failed"

echo ""
echo "✅ Import complete!"
echo ""
echo "💡 Next steps:"
echo "   1. Review the state: terraform show"
echo "   2. Check for differences: terraform plan"
echo "   3. Apply if needed: terraform apply"
echo ""

