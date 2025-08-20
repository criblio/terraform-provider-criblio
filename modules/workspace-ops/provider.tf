# Shared Provider Configuration for Examples
# This file contains the standard Cribl provider setup used across examples

terraform {
  # required_version = ">= 1.0"

  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "~> 1.9.2"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9"
    }
  }
}

# Cribl Provider Configuration
provider "criblio" {
  client_id       = var.cribl_client_id
  client_secret   = var.cribl_client_secret
  organization_id = var.cloud_tenant
  workspace_id    = var.workspace
}

