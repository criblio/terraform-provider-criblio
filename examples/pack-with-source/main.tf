
resource "criblio_pack_source" "my_packsource" {
  group_id = "default"
  pack     = criblio_pack.source_pack.id
  id       = "my_id"

  input_tcp = {
    auth_type = "manual"
    breaker_rulesets = [
      "multiline-json",
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
    max_active_cxn      = 8
    metadata = [
      {
        name  = "my_name"
        value = "\"my_value\""
      }
    ]
    pipeline = "my_pipeline"
    port     = 7592
    pq = {
      commit_frequency = 7
      compress         = "none"
      max_buffer_size  = 51
      max_file_size    = "100 MB"
      max_size         = "1000 MB"
      mode             = "smart"
      path             = "my_path"
    }
    pq_enabled = true
    preprocess = {
      args = [
        "-v",
      ]
      command  = "my_command"
      disabled = true
    }
    send_to_routes         = false
    socket_ending_max_wait = 8
    socket_idle_timeout    = 5
    socket_max_lifespan    = 9
    stale_channel_flush_ms = 1500
    streamtags = [
      "tcp",
    ]
    tls = {
      ca_path             = "my_ca_path"
      cert_path           = "my_cert_path"
      certificate_name    = "my_certificate_name"
      common_name_regex   = ".*"
      disabled            = true
      max_version         = "TLSv1.2"
      min_version         = "TLSv1.2"
      passphrase          = "my_passphrase"
      priv_key_path       = "my_priv_key_path"
      reject_unauthorized = true
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
