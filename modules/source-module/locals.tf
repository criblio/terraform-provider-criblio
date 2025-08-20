locals {
  # Base syslog configuration
  base_syslog_config = {
    id             = var.source_id
    host           = "0.0.0.0"
    tcp_port       = var.port
    udp_port       = var.port
    type           = "syslog"
    disabled       = var.disabled
    pq_enabled     = var.pq_enabled
    send_to_routes = length(var.connections) == 0

    description = var.description != "" ? var.description : null
    connections = length(var.connections) > 0 ? var.connections : null
    pipeline    = var.pipeline
    streamtags  = length(var.streamtags) > 0 ? var.streamtags : null
  }

  # Cribl HTTP configuration
  base_http_config = {
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

  # Regular HTTP source configuration
  base_http_source_config = {
    id             = var.source_id
    host           = "0.0.0.0"
    port           = var.port
    type           = "http"
    disabled       = var.disabled
    pq_enabled     = var.pq_enabled
    send_to_routes = length(var.connections) == 0

    description = var.description != "" ? var.description : null
    connections = length(var.connections) > 0 ? var.connections : null
    pipeline    = var.pipeline
    streamtags  = length(var.streamtags) > 0 ? var.streamtags : null

    max_active_req = 256
    auth_type      = "none"
  }

  # Regular TCP configuration
  base_tcp_source_config = {
    id             = var.source_id
    host           = "0.0.0.0"
    port           = var.port
    type           = "tcp"
    disabled       = var.disabled
    pq_enabled     = var.pq_enabled
    send_to_routes = length(var.connections) == 0

    description = var.description != "" ? var.description : null
    connections = length(var.connections) > 0 ? var.connections : null
    pipeline    = var.pipeline
    streamtags  = length(var.streamtags) > 0 ? var.streamtags : null

    tls = {
      disabled = true
    }
  }

  # Cribl TCP configuration
  base_tcp_config = {
    id             = var.source_id
    host           = "0.0.0.0"
    port           = var.port
    type           = "cribl_tcp"
    disabled       = var.disabled
    pq_enabled     = var.pq_enabled
    send_to_routes = length(var.connections) == 0

    description = var.description != "" ? var.description : null
    connections = length(var.connections) > 0 ? var.connections : null
    pipeline    = var.pipeline
    streamtags  = length(var.streamtags) > 0 ? var.streamtags : null

    tls = {
      disabled = true
    }
  }

  # OpenTelemetry configuration
  # In locals.tf, update the base_otlp_config:
  # In locals.tf, update the base_otlp_config:
  base_otlp_config = {
    id           = var.source_id
    type         = "open_telemetry"
    protocol     = "grpc" # or "http"
    port         = var.port
    otlp_version = "1.3.1"
    host         = "0.0.0.0"

    # Authentication
    auth_type        = "none" # or "basic" if using credentials
    auth_header_expr = "true"

    # Connection settings
    send_to_routes = length(var.connections) == 0
    connections    = length(var.connections) > 0 ? var.connections : null

    # OTLP specific settings - snake_case!
    extract_logs        = false
    extract_metrics     = true
    extract_spans       = true
    enable_health_check = false

    # Network settings
    ip_allowlist_regex      = "/.*/"
    ip_denylist_regex       = "/^$/"
    keep_alive_timeout      = 15
    max_active_cxn          = 1000
    max_active_req          = 256
    max_requests_per_socket = 0
    request_timeout         = 0
    socket_timeout          = 0

    # Other settings
    disabled           = var.disabled
    description        = var.description != "" ? var.description : null
    pipeline           = var.pipeline
    streamtags         = var.streamtags
    token_timeout_secs = 3600

    # TLS (if needed)
    # tls = {
    #   disabled = true
    # }
  }


  # REST Collector configuration  
  # In locals.tf, update the base_rest_collector_config:
  # Update locals.tf for collection

  # Transform configs - remove nulls and merge with custom
  syslog_config = var.source_type == "syslog" ? {
    for k, v in merge(local.base_syslog_config, var.custom_config) : k => v
    if v != null
  } : null

  http_config = var.source_type == "cribl_http" ? {
    for k, v in merge(local.base_http_config, var.custom_config) : k => v
    if v != null
  } : null

  http_source_config = var.source_type == "http" ? {
    for k, v in merge(local.base_http_source_config, var.custom_config) : k => v
    if v != null
  } : null

  tcp_source_config = var.source_type == "tcp" ? {
    for k, v in merge(local.base_tcp_source_config, var.custom_config) : k => v
    if v != null
  } : null

  tcp_config = var.source_type == "cribl_tcp" ? {
    for k, v in merge(local.base_tcp_config, var.custom_config) : k => v
    if v != null
  } : null

  otlp_config = var.source_type == "otlp" ? {
    for k, v in merge(local.base_otlp_config, var.custom_config) : k => v
    if v != null
  } : null


}