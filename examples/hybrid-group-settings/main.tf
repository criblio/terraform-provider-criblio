locals {
  streamtags_hybrid = [
    "datacenter1",
    "someothertag"
  ]
}
resource "criblio_group" "hybrid_worker_group" {
  id                   = "my-hybrid-group"
  name                 = "my-hybrid-group"
  is_fleet             = false
  on_prem              = true
  product              = "stream"
  streamtags           = local.streamtags_hybrid
  worker_remote_access = false
  max_worker_age       = "2h"
  provisioned          = false
}

resource "criblio_hybrid_group_system_settings" "hybrid_group_settings" {
  group_id = criblio_group.hybrid_worker_group.id

  api = {
    base_url          = "https://leader.example.com:9000"
    disable_api_cache = false
    disabled          = false
    host              = "leader.example.com"
    idle_session_ttl  = 3600
    listen_on_port    = true
    login_rate_limit  = "5/minute"
    port              = 9000
    protocol          = "https"
    scripts           = false
    sensitive_fields = [
      "password",
      "apiKey",
      "secret",
    ]
    ssl = {
      ca_path       = "/etc/ssl/certs/ca-bundle.crt"
      cert_path     = "/etc/cribl/ssl/server.crt"
      disabled      = false
      passphrase    = "changeit"
      priv_key_path = "/etc/cribl/ssl/server.key"
    }
    sso_rate_limit       = "10/minute"
    worker_remote_access = false
  }

  backups = {
    backup_persistence = "local"
    backups_directory  = "/var/cribl/backups"
  }

  custom_logo = {
    enabled          = false
    logo_description = ""
    logo_image       = ""
  }

  pii = {
    enable_pii_detection = true
  }

  proxy = {
    use_env_vars = false
  }

  rollback = {
    rollback_enabled = true
    rollback_retries = 3
    rollback_timeout = 600
  }

  shutdown = {
    drain_timeout = 30
  }

  sni = {
    disable_sni_routing = false
  }

  sockets = {
    directory = "/var/run/cribl"
  }

  system = {
    intercom = false
    upgrade  = "api"
  }

  tls = {
    default_cipher_list = "HIGH:!aNULL:!MD5"
    default_ecdh_curve  = "X25519:P-256"
    max_version         = "TLSv1.3"
    min_version         = "TLSv1.2"
    reject_unauthorized = true
  }

  upgrade_group_settings = {
    is_rolling  = true
    quantity    = 5
    retry_count = 3
    retry_delay = 60
  }

  upgrade_settings = {
    automatic_upgrade_check_period = "24h"
    disable_automatic_upgrade      = false
    enable_legacy_edge_upgrade     = false
    upgrade_source                 = "cdn"
  }

  workers = {
    count                    = 4
    enable_heap_snapshots    = false
    load_throttle_perc       = 85
    memory                   = 4096
    minimum                  = 2
    startup_max_conns        = 1024
    startup_throttle_timeout = 10000
    v8_single_thread         = false
  }
}