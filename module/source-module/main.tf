# modules/cribl-source/main.tf - CORRECTED
terraform {
  # required_version = ">= 1.0"
  
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      # version = "~> 1.0"
    }
  }
}

resource "criblio_source" "this" {
  id       = var.source_id
  group_id = var.group_id
  
  # For syslog - direct assignment with conditional
  input_syslog = var.source_type == "syslog" ? {
    input_syslog_syslog1 = merge(
      {
        # Required fields
        id             = var.source_id
        host           = "0.0.0.0"
        tcp_port       = var.port
        udp_port       = var.port
        type           = "syslog"
        disabled       = var.disabled
        pq_enabled     = var.pq_enabled
        send_to_routes = length(var.connections) == 0
      },
      # Optional fields only if they have values
      var.description != "" ? { description = var.description } : {},
      length(var.connections) > 0 ? { connections = var.connections } : {},
      var.pipeline != null ? { pipeline = var.pipeline } : {},
      length(var.streamtags) > 0 ? { streamtags = var.streamtags } : {},
      # Custom overrides
      var.custom_config
    )
  } : null
  
  # For HTTP - direct assignment with conditional
  input_cribl_http = var.source_type == "cribl_http" ? merge(
    {
      # Required fields
      id                       = var.source_id
      port                     = var.port
      activity_log_sample_rate = 100
      capture_headers          = false
      disabled                 = var.disabled
      enable_proxy_header      = false
      host                     = "0.0.0.0"
      max_active_req           = 256
      pq_enabled               = var.pq_enabled
      request_timeout          = 0
      send_to_routes           = length(var.connections) == 0
      streamtags               = var.streamtags
      type                     = "cribl_http"
      tls = {
        cert_path     = "$CRIBL_CLOUD_CRT"
        disabled      = false
        min_version   = "TLSv1.2"
        priv_key_path = "$CRIBL_CLOUD_KEY"
      }
    },
    # Optional fields
    var.description != "" ? { description = var.description } : {},
    length(var.connections) > 0 ? { connections = var.connections } : {},
    var.pipeline != null ? { pipeline = var.pipeline } : {},
    # Custom overrides
    var.custom_config
  ) : null
}