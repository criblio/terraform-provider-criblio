
resource "criblio_pack_source" "my_packsource" {
  group_id = "default"
  pack     = criblio_pack.source_pack.id
  id       = "my_id"

  input_tcp = {
    auth_type = "manual"
    breaker_rulesets = [
      ""
    ]
    connections = [
      {
        output   = "my_output"
        pipeline = "my_pipeline"
      }
    ]
    description         = "my_description"
    disabled            = true
    enable_header       = false
    enable_proxy_header = true
    environment         = "my_environment"
    host                = "my_host"
    id                  = "my_id"
    max_active_cxn      = 8.41
    metadata = [
      {
        name  = "my_name"
        value = "my_value"
      }
    ]
    pipeline = "my_pipeline"
    port     = 7592.47
    pq = {
      commit_frequency = 7.04
      compress         = "none"
      max_buffer_size  = 51.06
      max_file_size    = "100"
      max_size         = "1000"
      mode             = "smart"
      path             = "my_path"
    }
    pq_enabled = true
    preprocess = {
      args = [
        ""
      ]
      command  = "my_command"
      disabled = true
    }
    send_to_routes         = false
    socket_ending_max_wait = 8.02
    socket_idle_timeout    = 5.47
    socket_max_lifespan    = 9.38
    stale_channel_flush_ms = 8313205.9
    streamtags = [
      ""
    ]
    tls = {
      ca_path             = "my_ca_path"
      cert_path           = "my_cert_path"
      certificate_name    = "my_certificate_name"
      common_name_regex   = "{ \"see\": \"documentation\" }"
      disabled            = true
      max_version         = "TLSv1.1"
      min_version         = "TLSv1.1"
      passphrase          = "my_passphrase"
      priv_key_path       = "my_priv_key_path"
      reject_unauthorized = "{ \"see\": \"documentation\" }"
      request_cert        = false
    }
    type = "tcp"
  }
}

resource "criblio_pack" "source_pack" {
  id           = "pack-with-source"
  group_id     = "default"
  description  = "Pack with source"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"

}

