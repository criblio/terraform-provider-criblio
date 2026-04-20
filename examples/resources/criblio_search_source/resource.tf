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
      batch_timeout      = 3.51
      compress           = false
      content_format     = "Raw"
      heartbeat_interval = 4.94
      id                 = "...my_id..."
      locale             = "...my_locale..."
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      queries = [
        {
          path             = "...my_path..."
          query_expression = "...my_query_expression..."
        }
      ]
      query_selector       = "xml"
      read_existing_events = false
      send_bookmarks       = true
      subscription_name    = "...my_subscription_name..."
      targets = [
        "..."
      ]
      version   = "...my_version..."
      xml_query = "...my_xml_query..."
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