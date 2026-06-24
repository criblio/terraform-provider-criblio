package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGroupSystemSettings(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is prohibited by the current license")
	}

	resourceName := "criblio_group_system_settings.hybrid_group_settings"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: groupSystemSettingsConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "my-hybrid-group"),
					resource.TestCheckResourceAttr(resourceName, "api.base_url", "https://leader.example.com:9000"),
					resource.TestCheckResourceAttr(resourceName, "api.protocol", "https"),
					resource.TestCheckResourceAttr(resourceName, "api.ssl.disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "upgrade_settings.upgrade_source", "cdn"),
					resource.TestCheckResourceAttr(resourceName, "workers.count", "4"),
				),
			},
			{
				Config: groupSystemSettingsConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "proxy.use_env_vars", "true"),
					resource.TestCheckResourceAttr(resourceName, "workers.load_throttle_perc", "80"),
				),
			},
			{
				Config:   groupSystemSettingsConfig(true),
				PlanOnly: true,
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        "my-hybrid-group",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "group_id",
				ImportStateVerifyIgnore: []string{
					"api.ssl.passphrase",
					"apps",
					"custom_logo.logo_image",
					"support",
				},
			},
		},
	})
}

func groupSystemSettingsConfig(updated bool) string {
	useEnvVars := "false"
	loadThrottlePerc := "85"
	if updated {
		useEnvVars = "true"
		loadThrottlePerc = "80"
	}
	return `locals {
  streamtags_hybrid = [
    "datacenter1",
    "someothertag",
  ]
}

resource "criblio_group" "hybrid_worker_group" {
  id                   = "my-hybrid-group"
  is_fleet             = false
  max_worker_age       = "2h"
  name                 = "my-hybrid-group"
  on_prem              = true
  product              = "stream"
  provisioned          = false
  streamtags           = local.streamtags_hybrid
  worker_remote_access = false
}

resource "criblio_group_system_settings" "hybrid_group_settings" {
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
    sensitive_fields  = ["password", "apiKey", "secret"]
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
    use_env_vars = ` + useEnvVars + `
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
    load_throttle_perc       = ` + loadThrottlePerc + `
    memory                   = 4096
    minimum                  = 2
    startup_max_conns        = 1024
    startup_throttle_timeout = 10000
    v8_single_thread         = false
  }
}
`
}
