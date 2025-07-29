# AWS Provider Configuration
# This file can be extended to support other cloud providers like Azure, GCP, etc.

provider "aws" {
  region  = var.cloud_region  # Uses the cloud_region variable (default: us-west-2)
  profile = var.aws_profile   # Uses the aws_profile variable (default: Lab-Power)
}

# Example: Azure Provider (commented out)
# provider "azurerm" {
#   features {}
#   subscription_id = var.azure_subscription_id
#   tenant_id       = var.azure_tenant_id
# }

# Example: GCP Provider (commented out)
# provider "google" {
#   project = var.gcp_project_id
#   region  = var.gcp_region
# }