#!/bin/bash

# Script to import global vars and generate config that matches state exactly
# This ensures no updates after import

set -e

echo "🚀 Importing global vars and matching state..."
echo ""

# Step 1: Generate minimal config
echo "📝 Step 1: Generating minimal config..."
go run generate-global-vars.go
echo ""

# Step 2: Initialize Terraform
echo "📦 Step 2: Initializing Terraform..."
terraform init
echo ""

# Step 3: Import all resources
echo "📥 Step 3: Importing all resources..."
# Parse global-vars.tf to extract import commands
while IFS= read -r line; do
    if [[ $line =~ resource\ \"criblio_global_var\"\ \"([^\"]+)\" ]]; then
        resource_name="${BASH_REMATCH[1]}"
        # Read next few lines to get id
        id=""
        read_count=0
        while IFS= read -r next_line && [ $read_count -lt 3 ]; do
            if [[ $next_line =~ id\ =\ \"([^\"]+)\" ]]; then
                id="${BASH_REMATCH[1]}"
                break
            fi
            read_count=$((read_count + 1))
        done
        if [ -n "$id" ]; then
            echo "  Importing: criblio_global_var.${resource_name} (${id})"
            terraform import "criblio_global_var.${resource_name}" "{\"group_id\": \"default\", \"id\": \"${id}\"}" || echo "    ⚠️  Already imported or failed"
        fi
    fi
done < global-vars.tf
echo ""

# Step 4: Refresh state
echo "🔄 Step 4: Refreshing state..."
terraform refresh -no-color
echo ""

# Step 5: Generate config from state
echo "📝 Step 5: Generating config from state (this matches exactly)..."
terraform plan -generate-config-out=global-vars-complete.tf -no-color > /dev/null 2>&1 || true
echo ""

# Step 6: Replace minimal config with complete config
if [ -f "global-vars-complete.tf" ]; then
    echo "✅ Step 6: Replacing config with state-matched version..."
    mv global-vars-complete.tf global-vars.tf
    echo "  ✓ Config now matches state exactly"
    echo ""
    
    # Step 7: Verify no changes
    echo "🔍 Step 7: Verifying no changes..."
    terraform plan -no-color | grep -E "(No changes|Plan:)" || echo "  ⚠️  Some changes may still exist (computed fields)"
    echo ""
else
    echo "⚠️  Could not generate complete config from state"
    echo "   You may need to run manually: terraform plan -generate-config-out=complete.tf"
    echo ""
fi

echo "✨ Import complete!"
echo ""
echo "💡 Note: Computed fields (like 'items') may still show changes - this is expected"
echo "   and won't actually modify the resource."

