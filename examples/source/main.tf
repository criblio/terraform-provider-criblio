terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id    = "main"
  cloud_domain    = "cribl-playground.cloud"
}

resource "criblio_source" "my_http_source" {
  group_id = "default"
  id       = "http-listener"
  input_http = {
    activity_log_sample_rate = 10
    auth_tokens = [
      "secret-token-1",
      "secret-token-2",
    ]
    auth_tokens_ext = [
      {
        description = "Token for webhook ingestion"
        metadata = [
          {
            name  = "sourcetype"
            value = "\"http:webhook\""
          }
        ]
        token = "secret-token-1"
      }
    ]
    capture_headers = true
    connections = [
      {
        output   = "default"
        pipeline = "default"
      }
    ]
    cribl_api               = "/cribl"
    description             = "HTTP listener for webhook events"
    disabled                = false
    elastic_api             = "/elastic"
    enable_health_check     = true
    enable_proxy_header     = false
    environment             = "main"
    host                    = "0.0.0.0"
    id                      = "http-listener"
    keep_alive_timeout      = 30
    max_active_req          = 512
    max_requests_per_socket = 1000
    metadata = [
      {
        name  = "source"
        value = "\"http\""
      }
    ]
    pipeline = "default"
    port     = 10085
    pq = {
      commit_frequency = 100
      compress         = "gzip"
      max_buffer_size  = 5000
      max_file_size    = "128 MB"
      max_size         = "20GB"
      mode             = "smart"
      path             = "/opt/cribl/state/queues"
    }
    pq_enabled      = false
    request_timeout = 30
    send_to_routes  = true
    socket_timeout  = 60
    splunk_hec_acks = false
    splunk_hec_api  = "/services/collector"
    streamtags = [
      "prod",
      "http",
    ]
    tls = {
      ca_path             = "/etc/ssl/certs/ca.pem"
      cert_path           = "/etc/ssl/certs/server.crt"
      certificate_name    = "http-listener-cert"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.3"
      min_version         = "TLSv1.2"
      passphrase          = "$${{secret:http_key_pass}"
      priv_key_path       = "/etc/ssl/private/server.key"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "http"
  }
}

data "criblio_sources" "my_sources" {
  group_id = "default"
}

output "my_sources" {
  value = data.criblio_sources.my_sources
}

output "my_http_source_content" {
  value       = data.criblio_sources.my_sources.items[0].input_http
  description = "Content of the HTTP source from the data source"
}

output "all_http_sources" {
  value = [
    for item in data.criblio_sources.my_sources.items : item.input_http
    if item.input_http != null
  ]
  description = "All HTTP sources from the data source"
}

output "http_source_by_id" {
  value = [
    for item in data.criblio_sources.my_sources.items : item.input_http
    if item.input_http != null && try(item.input_http.id, "") == "cribl_http_source"
  ]
  description = "HTTP source with specific ID 'cribl_http_source'"
}