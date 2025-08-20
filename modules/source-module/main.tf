resource "criblio_source" "this" {
  id       = var.source_id
  group_id = var.group_id

  # Syslog
  input_syslog = var.source_type == "syslog" ? {
    input_syslog_syslog1 = merge(
      {
        id             = var.source_id
        host           = "0.0.0.0"
        tcp_port       = var.port
        udp_port       = var.port
        type           = "syslog"
        disabled       = var.disabled
        pq_enabled     = var.pq_enabled
        send_to_routes = length(var.connections) == 0
      },
      var.description != "" ? { description = var.description } : {},
      length(var.connections) > 0 ? { connections = var.connections } : {},
      var.pipeline != null ? { pipeline = var.pipeline } : {},
      length(var.streamtags) > 0 ? { streamtags = var.streamtags } : {},
      var.custom_config
    )
  } : null

  # Cribl HTTP
  input_cribl_http = var.source_type == "cribl_http" ? merge(
    {
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
    var.description != "" ? { description = var.description } : {},
    length(var.connections) > 0 ? { connections = var.connections } : {},
    var.pipeline != null ? { pipeline = var.pipeline } : {},
    var.custom_config
  ) : null

  # Regular HTTP
  input_http = var.source_type == "http" ? local.http_source_config : null

  # Regular TCP  
  input_tcp = var.source_type == "tcp" ? local.tcp_source_config : null

  # Cribl TCP
  input_cribl_tcp = var.source_type == "cribl_tcp" ? local.tcp_config : null

  # OpenTelemetry
  input_open_telemetry = var.source_type == "otlp" ? local.otlp_config : null
}