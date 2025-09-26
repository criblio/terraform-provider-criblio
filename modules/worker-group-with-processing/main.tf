# Worker Group with Pack and Pipeline Processing Module
# This module creates a complete worker group with pack installation, processing pipeline, and deployment

terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}
# Worker Group
resource "criblio_group" "worker_group" {
  cloud = {
    provider = var.cloud_provider
    region   = var.cloud_region
  }
  estimated_ingest_rate = var.estimated_ingest_rate
  id                    = var.group_id
  is_fleet              = var.is_fleet
  name                  = var.group_name
  on_prem               = var.on_prem
  product               = var.product
  provisioned           = var.provisioned
  streamtags            = var.streamtags
  worker_remote_access  = var.worker_remote_access
  description           = var.description
}

# Pack Installations
resource "criblio_pack" "palo_alto_pack" {
  count        = var.install_palo_alto_pack ? 1 : 0
  id           = "${var.group_id}-palo-alto-pack"
  group_id     = criblio_group.worker_group.id
  filename     = var.palo_alto_pack_filename
  description  = var.palo_alto_pack_description
  disabled     = var.palo_alto_pack_disabled
  display_name = var.palo_alto_pack_display_name
  version      = var.palo_alto_pack_version
}

resource "criblio_pack" "crowdstrike_pack" {
  count        = var.install_crowdstrike_pack ? 1 : 0
  id           = "${var.group_id}-crowdstrike-pack"
  group_id     = criblio_group.worker_group.id
  description  = var.crowdstrike_pack_description
  disabled     = var.crowdstrike_pack_disabled
  display_name = var.crowdstrike_pack_display_name
  version      = var.crowdstrike_pack_version
  # CrowdStrike pack created from pack_pipeline configuration
}

# CrowdStrike Pack Event Breaker
resource "criblio_pack_breakers" "crowdstrike_event_breaker" {
  count          = var.install_crowdstrike_pack ? 1 : 0
  group_id       = criblio_group.worker_group.id
  id             = "crowdstrike-event-breaker"
  pack           = criblio_pack.crowdstrike_pack[0].id
  description    = "CrowdStrike Event Breaker & Timestamp Processing"
  lib            = "custom"
  min_raw_length = 256
  tags           = "Crowdstrike,Security,fdr"

  rules = [
    {
      condition              = "/\"ContextTimeStamp\":/.test(_raw)"
      type                   = "regex"
      timestamp_anchor_regex = "/\"ContextTimeStamp\":\"/"
      timestamp = {
        type   = "format"
        length = 50
        format = "%s.%L"
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 768000
      disabled            = false
      event_breaker_regex = "/[\\r\\n]+(?=\\{\\\")/",
      name                = "ContextTimeStamp Breaker"
      fields              = []
      parser_enabled      = false
      should_use_data_raw = false
    },
    {
      condition              = "/\"timestamp\":\"\\d{4}-/.test(_raw) && !/\"ContextTimeStamp\":/.test(_raw)"
      type                   = "regex"
      timestamp_anchor_regex = "/\"timestamp\":\"/"
      timestamp = {
        type   = "format"
        length = 50
        format = "%Y-%m-%dT%H:%M:%SZ"
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 768000
      disabled            = false
      event_breaker_regex = "/[\\n\\r]+(?=\\{\\\")/",
      name                = "timestamp human readable breaker"
      fields              = []
      parser_enabled      = false
      should_use_data_raw = false
    },
    {
      condition              = "/\"timestamp\":\"\\d{10}/.test(_raw) && !_raw.includes('ContextTimeStamp')"
      type                   = "regex"
      timestamp_anchor_regex = "/\"timestamp\":\"/"
      timestamp = {
        type   = "format"
        length = 50
        format = "%Q"
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 768000
      disabled            = false
      event_breaker_regex = "/[\\n\\r]+(?=\\{\\\")/",
      name                = "timestamp field epoch"
      fields              = []
      parser_enabled      = false
      should_use_data_raw = false
    },
    {
      condition              = "(/\"_time\":/.test(_raw) || /\"Time\":/.test(_raw) ) && !_raw.includes('ContextTimeStamp') && !/\"timestamp\":/.test(_raw)"
      type                   = "regex"
      timestamp_anchor_regex = "/(\"_time\":\")|(\"Time\":\")/"
      timestamp = {
        type   = "format"
        length = 50
        format = "%s.%L|%s.%f"
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 768000
      disabled            = false
      event_breaker_regex = "/[\\n\\r]+(?=\\{\\\")/",
      name                = "_time or Time fields"
      fields              = []
      parser_enabled      = false
      should_use_data_raw = false
    },
    {
      condition              = "true"
      type                   = "regex"
      timestamp_anchor_regex = "/^/"
      timestamp = {
        type   = "current"
        length = 150
      }
      timestamp_timezone  = "local"
      timestamp_earliest  = "-420weeks"
      timestamp_latest    = "+1week"
      max_event_bytes     = 768000
      disabled            = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/",
      name                = "crowdstrike fallback"
      fields              = []
      parser_enabled      = false
      should_use_data_raw = false
    }
  ]

  depends_on = [criblio_pack.crowdstrike_pack]
}

# Processing Pipeline
resource "criblio_pipeline" "data_processing" {
  count    = var.create_pipeline ? 1 : 0
  group_id = criblio_group.worker_group.id
  id       = "${var.group_id}-processing"
  conf = {
    async_func_timeout = var.pipeline_timeout
    description        = var.pipeline_description
    functions = [
      {
        id          = "eval"
        filter      = "true"
        disabled    = false
        description = "Add region metadata"
        conf = {
          add = jsonencode([
            {
              "name"  = "_region"
              "value" = var.region_identifier
            },
            {
              "name"  = "_processing_time"
              "value" = "C.Time.now()"
            }
          ])
        }
      }
    ]
    output     = var.pipeline_output
    streamtags = var.pipeline_streamtags
  }
}

# Commit Configuration
resource "criblio_commit" "configuration_commit" {
  count     = var.auto_commit ? 1 : 0
  effective = true
  group     = criblio_group.worker_group.id
  message   = var.commit_message
  depends_on = [
    criblio_pack.palo_alto_pack,
    criblio_pack.crowdstrike_pack,
    criblio_pack_breakers.crowdstrike_event_breaker,
    criblio_pipeline.data_processing
  ]
}

# Config Version Data Source
data "criblio_config_version" "committed_config" {
  count      = var.auto_commit ? 1 : 0
  id         = criblio_group.worker_group.id
  depends_on = [criblio_commit.configuration_commit]
}

# Deploy Configuration
resource "criblio_deploy" "configuration_deploy" {
  count   = var.auto_deploy ? 1 : 0
  id      = criblio_group.worker_group.id
  version = data.criblio_config_version.committed_config[0].items[0]
}
