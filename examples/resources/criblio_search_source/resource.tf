resource "criblio_search_source" "my_searchsource" {
  auth_tokens = [
    {
      description = "...my_description..."
      enabled     = false
      token       = "...my_token..."
    }
  ]
  cribl_api       = "/cribl/_bulk"
  description     = "...my_description..."
  disabled        = false
  elastic_api     = "/elastic/_bulk"
  host            = "...my_host..."
  id              = "...my_id..."
  port            = 8088
  prometheus_api  = "/write"
  splunk_hec_acks = true
  splunk_hec_api  = "/services/collector"
  subscriptions = [
    {
      batch_timeout      = 5
      content_format     = "RenderedText"
      heartbeat_interval = 60
      id                 = "default-subscription"
      queries = [
        {
          path             = "Security"
          query_expression = "*"
        }
      ]
      query_selector    = "simple"
      subscription_name = "subscription-1"
      targets = [
      ]
      version   = "7f0c2f2e-1c3b-4d2a-9d6e-5a1b2c3d4e5f"
      xml_query = "//Event/System"
    }
  ]
  tcp_port = 36722
  tls = {
    cert_path     = "...my_cert_path..."
    disabled      = true
    min_version   = "TLSv1.1"
    priv_key_path = "...my_priv_key_path..."
  }
  type     = "splunk_hec"
  udp_port = 38619
}