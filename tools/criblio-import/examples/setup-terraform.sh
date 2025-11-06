#!/bin/bash

# Complete setup script: Initialize Terraform, Import Groups, and Apply

set -e

echo "🚀 Setting up Terraform for Cribl Groups"
echo "========================================"
echo ""

# Step 1: Create provider configuration if it doesn't exist
if [ ! -f "provider.tf" ]; then
    echo "📝 Creating provider.tf..."
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
  # See README.md for authentication setup
}
EOF
    echo "✅ Created provider.tf"
    echo ""
fi

# Step 2: Initialize Terraform
echo "📦 Initializing Terraform..."
terraform init
echo ""

# Step 3: Import groups
echo "📥 Importing groups into state..."
./import-groups.sh

# Step 4: Show plan
echo ""
echo "📊 Checking Terraform plan..."
terraform plan

echo ""
echo "✨ Setup complete!"
echo ""
echo "💡 If the plan looks good, run: terraform apply"
echo ""

