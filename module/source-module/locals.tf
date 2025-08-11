# modules/cribl-source/locals.tf - Updated
locals {
  # Base syslog configuration
  base_syslog_config = {
    id             = var.source_id
    host           = "0.0.0.0"
    tcp_port       = var.port  # Now required, no coalesce needed
    udp_port       = var.port  # Now required, no coalesce needed
    type           = "syslog"
    disabled       = var.disabled
    pq_enabled     = var.pq_enabled
    send_to_routes = length(var.connections) == 0
    
    # Add optional fields only if they have values
    description = var.description != "" ? var.description : null
    connections = length(var.connections) > 0 ? var.connections : null
    pipeline    = var.pipeline
    streamtags  = length(var.streamtags) > 0 ? var.streamtags : null
  }
  
  base_http_config = {
    id                       = var.source_id
    port                     = var.port  # Now required, no coalesce needed
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
    
    description = var.description != "" ? var.description : null
    connections = length(var.connections) > 0 ? var.connections : null
    pipeline    = var.pipeline
    
    tls = {
      cert_path     = "$CRIBL_CLOUD_CRT"
      disabled      = false
      min_version   = "TLSv1.2"
      priv_key_path = "$CRIBL_CLOUD_KEY"
    }
  }
  
  # Use null instead of {} for false condition
  syslog_config = var.source_type == "syslog" ? {
    for k, v in merge(local.base_syslog_config, var.custom_config) : k => v
    if v != null
  } : null
  
  http_config = var.source_type == "cribl_http" ? {
    for k, v in merge(local.base_http_config, var.custom_config) : k => v
    if v != null
  } : null
}