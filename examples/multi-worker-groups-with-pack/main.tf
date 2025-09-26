# Multi Worker Groups with Pack and Pipeline Installation Example
# This example creates 6 worker groups (3 AWS, 3 Azure) using a reusable module

# AWS Worker Groups
module "aws_east_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "aws-east-workers"
  group_name            = "aws-east-workers"
  cloud_provider        = "aws"
  cloud_region          = "us-east-1"
  estimated_ingest_rate = 1024
  description           = "AWS East region worker group for production workloads"
  streamtags            = ["aws", "east", "production"]

  # Palo Alto Pack Configuration
  install_palo_alto_pack      = true
  palo_alto_pack_filename     = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description  = "Palo Alto Networks pack for AWS East workers"
  palo_alto_pack_display_name = "Palo Alto Networks Source Pack - AWS East"
  palo_alto_pack_version      = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for AWS East workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - AWS East"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for AWS East region"
  region_identifier    = "aws-east"
  pipeline_streamtags  = ["aws", "east", "processing"]

  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to AWS East workers"
}

module "aws_west_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "aws-west-workers"
  group_name            = "aws-west-workers"
  cloud_provider        = "aws"
  cloud_region          = "us-west-2"
  estimated_ingest_rate = 2048
  description           = "AWS West region worker group for production workloads"
  streamtags            = ["aws", "west", "production"]

  # Palo Alto Pack Configuration
  install_palo_alto_pack      = true
  palo_alto_pack_filename     = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description  = "Palo Alto Networks pack for AWS West workers"
  palo_alto_pack_display_name = "Palo Alto Networks Source Pack - AWS West"
  palo_alto_pack_version      = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for AWS West workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - AWS West"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for AWS West region"
  region_identifier    = "aws-west"
  pipeline_streamtags  = ["aws", "west", "processing"]

  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to AWS West workers"
}

module "aws_central_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "aws-central-workers"
  group_name            = "aws-central-workers"
  cloud_provider        = "aws"
  cloud_region          = "us-central-1"
  estimated_ingest_rate = 2048
  description           = "AWS Central region worker group for production workloads"
  streamtags            = ["aws", "central", "production"]

  # Palo Alto Pack Configuration
  install_palo_alto_pack      = true
  palo_alto_pack_filename     = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description  = "Palo Alto Networks pack for AWS Central workers"
  palo_alto_pack_display_name = "Palo Alto Networks Source Pack - AWS Central"
  palo_alto_pack_version      = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for AWS Central workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - AWS Central"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for AWS Central region"
  region_identifier    = "aws-central"
  pipeline_streamtags  = ["aws", "central", "processing"]

  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to AWS Central workers"
}

# Azure Worker Groups
module "azure_east_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "azure-east-workers"
  group_name            = "azure-east-us-workers"
  cloud_provider        = "azure"
  cloud_region          = "eastus"
  estimated_ingest_rate = 2048
  description           = "Azure East US region worker group for production workloads"
  streamtags            = ["azure", "east", "production"]

  # Palo Alto Pack Configuration
  install_palo_alto_pack      = true
  palo_alto_pack_filename     = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description  = "Palo Alto Networks pack for Azure East workers"
  palo_alto_pack_display_name = "Palo Alto Networks Source Pack - Azure East"
  palo_alto_pack_version      = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for Azure East workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - Azure East"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for Azure East region"
  region_identifier    = "azure-east"
  pipeline_streamtags  = ["azure", "east", "processing"]

  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to Azure East workers"
}

module "azure_west_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "azure-west-workers"
  group_name            = "azure-west-workers"
  cloud_provider        = "azure"
  cloud_region          = "westus2"
  estimated_ingest_rate = 1024
  description           = "Azure West US region worker group for production workloads"
  streamtags            = ["azure", "west", "production"]

  # Palo Alto Pack Configuration
  install_palo_alto_pack      = true
  palo_alto_pack_filename     = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description  = "Palo Alto Networks pack for Azure West workers"
  palo_alto_pack_display_name = "Palo Alto Networks Source Pack - Azure West"
  palo_alto_pack_version      = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for Azure West workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - Azure West"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for Azure West region"
  region_identifier    = "azure-west"
  pipeline_streamtags  = ["azure", "west", "processing"]

  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to Azure West workers"
}

module "azure_europe_workers" {
  source = "../../modules/worker-group-with-processing"

  # Worker Group Configuration
  group_id              = "azure-europe-workers"
  group_name            = "azure-europe-workers"
  cloud_provider        = "azure"
  cloud_region          = "westeurope"
  estimated_ingest_rate = 2048
  description           = "Azure West Europe region worker group for production workloads"
  streamtags            = ["azure", "europe", "production"]

  # Palo Alto Pack Configuration
  install_palo_alto_pack      = true
  palo_alto_pack_filename     = "cribl-palo-alto-networks-1.1.5.crbl"
  palo_alto_pack_description  = "Palo Alto Networks pack for Azure Europe workers"
  palo_alto_pack_display_name = "Palo Alto Networks Source Pack - Azure Europe"
  palo_alto_pack_version      = "1.1.5"

  # CrowdStrike Pack Configuration
  install_crowdstrike_pack      = true
  crowdstrike_pack_description  = "CrowdStrike Security Pack for Azure Europe workers"
  crowdstrike_pack_display_name = "CrowdStrike Security Pack - Azure Europe"

  # Pipeline Configuration
  pipeline_description = "Standard data processing pipeline for Azure Europe region"
  region_identifier    = "azure-europe"
  pipeline_streamtags  = ["azure", "europe", "processing"]

  # Deployment Configuration
  commit_message = "Deploy Palo Alto and CrowdStrike packs with processing pipeline to Azure Europe workers"
}

/*
# Outputs
output "worker_groups" {
  description = "Details of all created worker groups"
  value = {
    aws_east = module.aws_east_workers.summary
    # aws_west = module.aws_west_workers.summary
    # aws_central = module.aws_central_workers.summary
    azure_east = module.azure_east_workers.summary
    # azure_west = module.azure_west_workers.summary
    # azure_europe = module.azure_europe_workers.summary
  }
}

output "palo_alto_pack_installations" {
  description = "Details of Palo Alto pack installations on all worker groups"
  value = {
    aws_east_pack = module.aws_east_workers.palo_alto_pack_details
    # aws_west_pack = module.aws_west_workers.palo_alto_pack_details
    # aws_central_pack = module.aws_central_workers.palo_alto_pack_details
    azure_east_pack = module.azure_east_workers.palo_alto_pack_details
    # azure_west_pack = module.azure_west_workers.palo_alto_pack_details
    # azure_europe_pack = module.azure_europe_workers.palo_alto_pack_details
  }
}

output "crowdstrike_pack_installations" {
  description = "Details of CrowdStrike pack installations on all worker groups"
  value = {
    aws_east_pack = module.aws_east_workers.crowdstrike_pack_details
    # aws_west_pack = module.aws_west_workers.crowdstrike_pack_details
    # aws_central_pack = module.aws_central_workers.crowdstrike_pack_details
    azure_east_pack = module.azure_east_workers.crowdstrike_pack_details
    # azure_west_pack = module.azure_west_workers.crowdstrike_pack_details
    # azure_europe_pack = module.azure_europe_workers.crowdstrike_pack_details
  }
}

output "crowdstrike_pack_breakers" {
  description = "Details of CrowdStrike pack event breakers on all worker groups"
  value = {
    aws_east_breaker = module.aws_east_workers.crowdstrike_pack_breakers_id
    # aws_west_breaker = module.aws_west_workers.crowdstrike_pack_breakers_id
    # aws_central_breaker = module.aws_central_workers.crowdstrike_pack_breakers_id
    azure_east_breaker = module.azure_east_workers.crowdstrike_pack_breakers_id
    # azure_west_breaker = module.azure_west_workers.crowdstrike_pack_breakers_id
    # azure_europe_breaker = module.azure_europe_workers.crowdstrike_pack_breakers_id
  }
}

output "pipelines" {
  description = "Details of processing pipelines on all worker groups"
  value = {
    aws_east_pipeline = module.aws_east_workers.pipeline_details
    # aws_west_pipeline = module.aws_west_workers.pipeline_details
    # aws_central_pipeline = module.aws_central_workers.pipeline_details
    azure_east_pipeline = module.azure_east_workers.pipeline_details
    # azure_west_pipeline = module.azure_west_workers.pipeline_details
    # azure_europe_pipeline = module.azure_europe_workers.pipeline_details
  }
}

output "deployment_status" {
  description = "Deployment status for all worker groups"
  value = {
    aws_east = module.aws_east_workers.deployment_status
    # aws_west = module.aws_west_workers.deployment_status
    # aws_central = module.aws_central_workers.deployment_status
    azure_east = module.azure_east_workers.deployment_status
    # azure_west = module.azure_west_workers.deployment_status
    # azure_europe = module.azure_europe_workers.deployment_status
  }
}

*/