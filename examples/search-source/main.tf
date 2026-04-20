# One block per Local Search source type. Each uses a distinct port in the 31100+ range; if apply fails with
# "port in use" or "outside pool", pick unused ports from your workspace UI dropdown and update here.

locals {
  # Base port — increment per resource so each listener is unique.
  p = 31100
}

resource "criblio_search_source" "cribl_http" {
  description = "Cribl HTTP ingest (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_cribl_http"
  port        = local.p + 0
  type        = "cribl_http"

  cribl_api = "/cribl/_bulk"
  tls = {
    disabled    = false
    min_version = "TLSv1.3"
  }
}

resource "criblio_search_source" "datadog_agent" {
  description = "Datadog Agent (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_datadog_agent"
  port        = local.p + 1
  type        = "datadog_agent"
}

resource "criblio_search_source" "elastic" {
  description = "Elastic bulk API (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_elastic"
  port        = local.p + 2
  type        = "elastic"

  elastic_api = "/elastic/_bulk"
}

resource "criblio_search_source" "http_raw" {
  description = "Raw HTTP (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_http_raw"
  port        = local.p + 3
  type        = "http_raw"
}

resource "criblio_search_source" "open_telemetry" {
  description = "OpenTelemetry (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_open_telemetry"
  port        = local.p + 4
  type        = "open_telemetry"
}

resource "criblio_search_source" "prometheus_rw" {
  description = "Prometheus remote write (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_prometheus_rw"
  port        = local.p + 5
  type        = "prometheus_rw"

  prometheus_api = "/write"
  udp_port       = local.p + 15
}

resource "criblio_search_source" "splunk" {
  description = "Splunk ingest (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_splunk"
  port        = local.p + 6
  type        = "splunk"
}

resource "criblio_search_source" "splunk_hec" {
  description = "Splunk HEC ingest for Local Search"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_splunk_hec"
  port        = local.p + 7
  type        = "splunk_hec"

  splunk_hec_api  = "/services/collector"
  splunk_hec_acks = false

  auth_tokens = [
    {
      description = "Primary HEC token"
      enabled     = true
      token       = "changeme-replace-with-your-hec-token"
    }
  ]
}

resource "criblio_search_source" "syslog" {
  description = "Syslog TCP+UDP (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_syslog"
  port        = local.p + 8
  type        = "syslog"

  tcp_port = local.p + 8
  udp_port = local.p + 9
}

resource "criblio_search_source" "tcp" {
  description = "TCP (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_tcp"
  port        = local.p + 10
  type        = "tcp"
}

resource "criblio_search_source" "tcpjson" {
  description = "TCP JSON (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_tcpjson"
  port        = local.p + 11
  type        = "tcpjson"
}

resource "criblio_search_source" "wef" {
  description = "Windows Event Forwarding (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_wef"
  port        = local.p + 12
  type        = "wef"

  subscriptions = [
    {
      subscription_name    = "Security"
      content_format       = "Raw"
      heartbeat_interval   = 60
      batch_timeout        = 30
      read_existing_events = false
      send_bookmarks       = true
      compress             = false
      targets              = ["*.example.wef.local"]
    }
  ]
}

resource "criblio_search_source" "wiz_webhook" {
  description = "Wiz webhook (example)"
  disabled    = false
  host        = "0.0.0.0"
  id          = "example_in_wiz_webhook"
  port        = local.p + 13
  type        = "wiz_webhook"
}

output "search_source_ids" {
  description = "Ids of the example Local Search sources (sensitive entries omitted)."
  value = {
    cribl_http     = criblio_search_source.cribl_http.id
    datadog_agent  = criblio_search_source.datadog_agent.id
    elastic        = criblio_search_source.elastic.id
    http_raw       = criblio_search_source.http_raw.id
    open_telemetry = criblio_search_source.open_telemetry.id
    prometheus_rw  = criblio_search_source.prometheus_rw.id
    splunk         = criblio_search_source.splunk.id
    splunk_hec     = criblio_search_source.splunk_hec.id
    syslog         = criblio_search_source.syslog.id
    tcp            = criblio_search_source.tcp.id
    tcpjson        = criblio_search_source.tcpjson.id
    wef            = criblio_search_source.wef.id
    wiz_webhook    = criblio_search_source.wiz_webhook.id
  }
}
