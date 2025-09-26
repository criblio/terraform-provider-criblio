# Worker Group Outputs
output "worker_group" {
  description = "Complete worker group resource details"
  value       = criblio_group.worker_group
}

output "worker_group_id" {
  description = "ID of the created worker group"
  value       = criblio_group.worker_group.id
}

output "worker_group_name" {
  description = "Name of the created worker group"
  value       = criblio_group.worker_group.name
}

output "cloud_details" {
  description = "Cloud provider and region details"
  value = {
    provider = criblio_group.worker_group.cloud.provider
    region   = criblio_group.worker_group.cloud.region
  }
}

# Palo Alto Pack Outputs
output "palo_alto_pack" {
  description = "Palo Alto pack installation details (if pack was installed)"
  value       = var.install_palo_alto_pack ? criblio_pack.palo_alto_pack[0] : null
}

output "palo_alto_pack_id" {
  description = "ID of the installed Palo Alto pack (if pack was installed)"
  value       = var.install_palo_alto_pack ? criblio_pack.palo_alto_pack[0].id : null
}

output "palo_alto_pack_details" {
  description = "Palo Alto pack installation summary"
  value = var.install_palo_alto_pack ? {
    id           = criblio_pack.palo_alto_pack[0].id
    group_id     = criblio_pack.palo_alto_pack[0].group_id
    display_name = criblio_pack.palo_alto_pack[0].display_name
    version      = criblio_pack.palo_alto_pack[0].version
    disabled     = criblio_pack.palo_alto_pack[0].disabled
  } : null
}

# CrowdStrike Pack Outputs
output "crowdstrike_pack" {
  description = "CrowdStrike pack installation details (if pack was installed)"
  value       = var.install_crowdstrike_pack ? criblio_pack.crowdstrike_pack[0] : null
}

output "crowdstrike_pack_id" {
  description = "ID of the installed CrowdStrike pack (if pack was installed)"
  value       = var.install_crowdstrike_pack ? criblio_pack.crowdstrike_pack[0].id : null
}

output "crowdstrike_pack_details" {
  description = "CrowdStrike pack installation summary"
  value = var.install_crowdstrike_pack ? {
    id           = criblio_pack.crowdstrike_pack[0].id
    group_id     = criblio_pack.crowdstrike_pack[0].group_id
    display_name = criblio_pack.crowdstrike_pack[0].display_name
    version      = criblio_pack.crowdstrike_pack[0].version
    disabled     = criblio_pack.crowdstrike_pack[0].disabled
  } : null
}

# CrowdStrike Pack Event Breaker Outputs
output "crowdstrike_pack_breakers" {
  description = "CrowdStrike pack event breaker details (if event breaker was created)"
  value       = var.install_crowdstrike_pack ? criblio_pack_breakers.crowdstrike_event_breaker[0] : null
}

output "crowdstrike_pack_breakers_id" {
  description = "ID of the CrowdStrike pack event breaker (if event breaker was created)"
  value       = var.install_crowdstrike_pack ? criblio_pack_breakers.crowdstrike_event_breaker[0].id : null
}

# Pipeline Outputs
output "pipeline" {
  description = "Processing pipeline details (if pipeline was created)"
  value       = var.create_pipeline ? criblio_pipeline.data_processing[0] : null
}

output "pipeline_id" {
  description = "ID of the processing pipeline (if pipeline was created)"
  value       = var.create_pipeline ? criblio_pipeline.data_processing[0].id : null
}

output "pipeline_details" {
  description = "Pipeline configuration summary"
  value = var.create_pipeline ? {
    id          = criblio_pipeline.data_processing[0].id
    group_id    = criblio_pipeline.data_processing[0].group_id
    description = criblio_pipeline.data_processing[0].conf.description
  } : null
}

# Deployment Outputs
output "commit_details" {
  description = "Commit details (if auto_commit is enabled)"
  value = var.auto_commit ? {
    group   = criblio_commit.configuration_commit[0].group
    message = criblio_commit.configuration_commit[0].message
  } : null
}

output "config_version" {
  description = "Configuration version details (if auto_commit is enabled)"
  value       = var.auto_commit ? data.criblio_config_version.committed_config[0].items[0] : null
}

output "deployment_status" {
  description = "Complete deployment status"
  value = var.auto_deploy ? {
    group_id       = criblio_deploy.configuration_deploy[0].id
    config_version = data.criblio_config_version.committed_config[0].items[0]
    commit_message = criblio_commit.configuration_commit[0].message
  } : null
}

# Summary Output
output "summary" {
  description = "Summary of all created resources"
  value = {
    worker_group = {
      id                    = criblio_group.worker_group.id
      name                  = criblio_group.worker_group.name
      provider              = criblio_group.worker_group.cloud.provider
      region                = criblio_group.worker_group.cloud.region
      estimated_ingest_rate = criblio_group.worker_group.estimated_ingest_rate
    }
    palo_alto_pack_installed   = var.install_palo_alto_pack
    crowdstrike_pack_installed = var.install_crowdstrike_pack
    pipeline_created           = var.create_pipeline
    auto_deployed              = var.auto_deploy
  }
}
